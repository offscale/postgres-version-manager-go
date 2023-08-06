// Package pvm remote_fetch.go from github.com/fergusstrange/embedded-postgres and modified
// (so should be considered under same license; except for `extractArchiveToDir`)

package pvm

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xi2/xz"
)

// RemoteFetchStrategy provides a strategy to fetch a Postgres binary so that it is available for use.
type RemoteFetchStrategy func() error

//nolint:funlen
func defaultRemoteFetchStrategy(remoteFetchHost string, versionStrategy VersionStrategy, cacheLocator CacheLocator, versionManagerRoot string, noExtract bool, noInstall bool) RemoteFetchStrategy {
	return func() error {
		if _, existsAlready := cacheLocator(); existsAlready {
			return nil
		}
		operatingSystem, architecture, version := versionStrategy()
		if noInstall {
			return fmt.Errorf("noInstall toggled and installation not found for %s-%s-%s",
				operatingSystem,
				architecture,
				version)
		}

		jarDownloadURL := fmt.Sprintf("%s/io/zonky/test/postgres/embedded-postgres-binaries-%s-%s/%s/embedded-postgres-binaries-%s-%s-%s.jar",
			remoteFetchHost,
			operatingSystem,
			architecture,
			version,
			operatingSystem,
			architecture,
			version)

		jarDownloadResponse, err := http.Get(jarDownloadURL)
		if err != nil {
			return fmt.Errorf("unable to connect to %s", remoteFetchHost)
		}

		defer closeBody(jarDownloadResponse)()

		if jarDownloadResponse.StatusCode != http.StatusOK {
			return fmt.Errorf("no version found matching %s", version)
		}

		jarBodyBytes, err := io.ReadAll(jarDownloadResponse.Body)
		if err != nil {
			return errorFetchingPostgres(err)
		}

		shaDownloadURL := fmt.Sprintf("%s.sha256", jarDownloadURL)
		shaDownloadResponse, err := http.Get(shaDownloadURL)

		defer closeBody(shaDownloadResponse)()

		if err == nil && shaDownloadResponse.StatusCode == http.StatusOK {
			if shaBodyBytes, err := io.ReadAll(shaDownloadResponse.Body); err == nil {
				jarChecksum := sha256.Sum256(jarBodyBytes)
				if !bytes.Equal(shaBodyBytes, []byte(hex.EncodeToString(jarChecksum[:]))) {
					return errors.New("downloaded checksums do not match")
				}
			}
		}

		if noExtract {
			return nil
		}
		// Extract out the archive to the cache dir:
		// e.g., 'embedded-postgres-binaries-linux-amd64-15.1.0.txz' to "$HOME"'/postgres-version-manager/downloads/'
		var tarFilePath string
		if tarFilePath, err = decompressResponse(jarBodyBytes, jarDownloadResponse.ContentLength, cacheLocator, jarDownloadURL); err != nil {
			return err
		}

		destinationFolderPath := path.Join(versionManagerRoot, version)
		if _, err = os.Stat(path.Join(destinationFolderPath, "bin")); err == nil {
			return nil
		}
		return extractArchiveToDir(tarFilePath, destinationFolderPath)
	}
}

func extractArchiveToDir(archiveFilePath string, destinationFolderPath string) error {
	var err error
	var f *os.File

	if f, err = os.Open(archiveFilePath); err != nil {
		return err
	}

	var r *xz.Reader
	if r, err = xz.NewReader(f, 0); err != nil {
		return err
	}
	tr := tar.NewReader(r)
	alreadyClosed := false
	defer func(f *os.File) {
		if !alreadyClosed {
			if err = f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}(f)
	for {
		var hdr *tar.Header
		if hdr, err = tr.Next(); err != nil {
			if err == io.EOF {
				alreadyClosed = true
				return f.Close()
			}
			return err
		}
		name := path.Join(destinationFolderPath, hdr.Name)
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(name, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			var w *os.File
			if err = os.MkdirAll(path.Dir(name), hdr.FileInfo().Mode()); err != nil {
				return err
			}
			if w, err = os.Create(name); err != nil {
				return err
			}
			if _, err = io.Copy(w, tr); err != nil {
				return err
			}
			if err = w.Close(); err != nil {
				return err
			}
			if err = os.Chmod(name, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeLink:
			fmt.Println("TypeLink: " + name)
		case tar.TypeSymlink:
			targetPath := path.Join(path.Dir(name), hdr.Linkname)
			// To solve https://nvd.nist.gov/vuln/detail/CVE-2020-27833
			if !strings.HasPrefix(targetPath, destinationFolderPath) {
				return fmt.Errorf("invalid symlink %q -> %q", name, hdr.Linkname)
			}
			if err = os.Symlink(hdr.Linkname, name); err != nil {
				return err
			}
		default:
			fmt.Printf("doing nothing with: %s of %d\n", name, hdr.Typeflag)
		}
	}
}

func closeBody(resp *http.Response) func() {
	return func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func decompressResponse(bodyBytes []byte, contentLength int64, cacheLocator CacheLocator, downloadURL string) (string, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), contentLength)
	if err != nil {
		return "", errorFetchingPostgres(err)
	}

	for _, file := range zipReader.File {
		if !file.FileHeader.FileInfo().IsDir() && strings.HasSuffix(file.FileHeader.Name, ".txz") {
			archiveReader, err := file.Open()
			if err != nil {
				return "", errorExtractingPostgres(err)
			}

			archiveBytes, err := io.ReadAll(archiveReader)
			if err != nil {
				return "", errorExtractingPostgres(err)
			}

			cacheLocation, _ := cacheLocator()

			if err := os.MkdirAll(filepath.Dir(cacheLocation), 0755); err != nil {
				return "", errorExtractingPostgres(err)
			}

			if err := os.WriteFile(cacheLocation, archiveBytes, file.FileHeader.Mode()); err != nil {
				return "", errorExtractingPostgres(err)
			}

			return cacheLocation, nil
		}
	}

	return "", fmt.Errorf("error fetching postgres: cannot find binary in archive retrieved from %s", downloadURL)
}

func errorExtractingPostgres(err error) error {
	return fmt.Errorf("unable to extract postgres archive: %s", err)
}

func errorFetchingPostgres(err error) error {
	return fmt.Errorf("error fetching postgres: %s", err)
}
