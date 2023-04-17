package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

const NumberOfPostgresVersions uintptr = 135

var PostgresVersions [NumberOfPostgresVersions]embeddedpostgres.PostgresVersion = [NumberOfPostgresVersions]embeddedpostgres.PostgresVersion{
	"9.3.23",
	"9.3.24",
	"9.3.24-1",
	"9.3.25",
	"9.4.18",
	"9.4.19",
	"9.4.19-1",
	"9.4.20",
	"9.4.21",
	"9.4.22",
	"9.4.23",
	"9.4.24",
	"9.4.25",
	"9.4.25-1",
	"9.4.26",
	"9.5.13",
	"9.5.14",
	"9.5.14-1",
	"9.5.15",
	"9.5.16",
	"9.5.17",
	"9.5.18",
	"9.5.19",
	"9.5.20",
	"9.5.20-1",
	"9.5.21",
	"9.5.22",
	"9.5.23",
	"9.5.24",
	"9.5.24-1",
	"9.5.25",
	"9.6.9",
	"9.6.10",
	"9.6.10-1",
	"9.6.11",
	"9.6.12",
	"9.6.13",
	"9.6.14",
	"9.6.15",
	"9.6.16",
	"9.6.16-1",
	"9.6.17",
	"9.6.18",
	"9.6.19",
	"9.6.20",
	"9.6.20-1",
	"9.6.21",
	"9.6.22",
	"9.6.23",
	"9.6.24",
	"10.4.0",
	"10.5.0",
	"10.5.0-1",
	"10.6.0",
	"10.7.0",
	"10.8.0",
	"10.9.0",
	"10.10.0",
	"10.11.0",
	"10.11.0-1",
	"10.12.0",
	"10.13.0",
	"10.14.0",
	"10.15.0",
	"10.15.0-1",
	"10.16.0",
	"10.17.0",
	"10.18.0",
	"10.19.0",
	"10.20.0",
	"10.21.0",
	"10.22.0",
	"10.23.0",
	"11.0.0",
	"11.1.0",
	"11.2.0",
	"11.3.0",
	"11.4.0",
	"11.5.0",
	"11.6.0",
	"11.6.0-1",
	"11.7.0",
	"11.8.0",
	"11.9.0",
	"11.10.0",
	"11.10.0-1",
	"11.11.0",
	"11.12.0",
	"11.13.0",
	"11.14.0",
	"11.15.0",
	"11.16.0",
	"11.17.0",
	"11.18.0",
	"11.19.0",
	"12.0.0",
	"12.1.0",
	"12.1.0-1",
	"12.2.0",
	"12.3.0",
	"12.4.0",
	"12.5.0",
	"12.5.0-1",
	"12.6.0",
	"12.7.0",
	"12.8.0",
	"12.9.0",
	"12.10.0",
	"12.11.0",
	"12.12.0",
	"12.13.0",
	"12.14.0",
	"13.0.0",
	"13.1.0",
	"13.1.0-1",
	"13.2.0",
	"13.3.0",
	"13.4.0",
	"13.5.0",
	"13.6.0",
	"13.7.0",
	"13.8.0",
	"13.9.0",
	"13.10.0",
	"14.0.0",
	"14.1.0",
	"14.2.0",
	"14.3.0",
	"14.4.0",
	"14.5.0",
	"14.6.0",
	"14.7.0",
	"15.0.0",
	"15.1.0",
	"15.2.0"}

// Micro-optimization to avoid looping
func isValidVersion(v string) bool {
	switch v {
	case
		"9.3.23",
		"9.3.24",
		"9.3.24-1",
		"9.3.25",
		"9.4.18",
		"9.4.19",
		"9.4.19-1",
		"9.4.20",
		"9.4.21",
		"9.4.22",
		"9.4.23",
		"9.4.24",
		"9.4.25",
		"9.4.25-1",
		"9.4.26",
		"9.5.13",
		"9.5.14",
		"9.5.14-1",
		"9.5.15",
		"9.5.16",
		"9.5.17",
		"9.5.18",
		"9.5.19",
		"9.5.20",
		"9.5.20-1",
		"9.5.21",
		"9.5.22",
		"9.5.23",
		"9.5.24",
		"9.5.24-1",
		"9.5.25",
		"9.6.9",
		"9.6.10",
		"9.6.10-1",
		"9.6.11",
		"9.6.12",
		"9.6.13",
		"9.6.14",
		"9.6.15",
		"9.6.16",
		"9.6.16-1",
		"9.6.17",
		"9.6.18",
		"9.6.19",
		"9.6.20",
		"9.6.20-1",
		"9.6.21",
		"9.6.22",
		"9.6.23",
		"9.6.24",
		"10.4.0",
		"10.5.0",
		"10.5.0-1",
		"10.6.0",
		"10.7.0",
		"10.8.0",
		"10.9.0",
		"10.10.0",
		"10.11.0",
		"10.11.0-1",
		"10.12.0",
		"10.13.0",
		"10.14.0",
		"10.15.0",
		"10.15.0-1",
		"10.16.0",
		"10.17.0",
		"10.18.0",
		"10.19.0",
		"10.20.0",
		"10.21.0",
		"10.22.0",
		"10.23.0",
		"11.0.0",
		"11.1.0",
		"11.2.0",
		"11.3.0",
		"11.4.0",
		"11.5.0",
		"11.6.0",
		"11.6.0-1",
		"11.7.0",
		"11.8.0",
		"11.9.0",
		"11.10.0",
		"11.10.0-1",
		"11.11.0",
		"11.12.0",
		"11.13.0",
		"11.14.0",
		"11.15.0",
		"11.16.0",
		"11.17.0",
		"11.18.0",
		"11.19.0",
		"12.0.0",
		"12.1.0",
		"12.1.0-1",
		"12.2.0",
		"12.3.0",
		"12.4.0",
		"12.5.0",
		"12.5.0-1",
		"12.6.0",
		"12.7.0",
		"12.8.0",
		"12.9.0",
		"12.10.0",
		"12.11.0",
		"12.12.0",
		"12.13.0",
		"12.14.0",
		"13.0.0",
		"13.1.0",
		"13.1.0-1",
		"13.2.0",
		"13.3.0",
		"13.4.0",
		"13.5.0",
		"13.6.0",
		"13.7.0",
		"13.8.0",
		"13.9.0",
		"13.10.0",
		"14.0.0",
		"14.1.0",
		"14.2.0",
		"14.3.0",
		"14.4.0",
		"14.5.0",
		"14.6.0",
		"14.7.0",
		"15.0.0",
		"15.1.0",
		"15.2.0":
		return true
	}
	return false
}

/*
<metadata>
  <groupId>io.zonky.test.postgres</groupId>
  <artifactId>embedded-postgres-binaries-bom</artifactId>
  <versioning>
    <latest>15.2.0</latest>
    <release>15.2.0</release>
	<versions>

*/

type Metadata struct {
	XMLName    xml.Name `xml:"metadata"`
	Text       string   `xml:",chardata"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Versioning struct {
		Text     string `xml:",chardata"`
		Latest   string `xml:"latest"`
		Release  string `xml:"release"`
		Versions struct {
			Text    string   `xml:",chardata"`
			Version []string `xml:"version"`
		} `xml:"versions"`
		LastUpdated string `xml:"lastUpdated"`
	} `xml:"versioning"`
}

func getVersionsFromMaven(binaryRepositoryURL string) (error, []string) {
	var metadata Metadata
	{
		resp, err := http.Get(binaryRepositoryURL + "/io/zonky/test/postgres/embedded-postgres-binaries-bom/maven-metadata.xml")
		if err != nil {
			return err, nil
		}
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				panic(err)
			}
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err := xml.Unmarshal(body, &metadata); err != nil {
			return err, nil
		}
	}
	return nil, metadata.Versioning.Versions.Version
}

// Originally from embeddedpostgres so should be under its license
func startPostgres(config *ConfigStruct) error {
	postgresBinary := filepath.Join(config.BinariesPath, "bin", "pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "start", "-w",
		"-D", config.DataPath,
		"-o", fmt.Sprintf(`"-p %d"`, config.Port))

	syncedLog, err := newSyncedLogger(config.DataPath, os.Stdout)
	if err != nil {
		return err
	}

	postgresProcess.Stdout = syncedLog.file
	postgresProcess.Stderr = syncedLog.file

	if err := postgresProcess.Run(); err != nil {
		_ = syncedLog.flush()
		logContent, _ := readLogsOrTimeout(syncedLog.file)

		return fmt.Errorf("could not start postgres using %s:\n%s", postgresProcess.String(), string(logContent))
	}

	return nil
}

func stopPostgres(config *ConfigStruct) error {
	postgresBinary := filepath.Join(config.BinariesPath, "bin", "pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "stop", "-w",
		"-D", config.DataPath,
		"-o", fmt.Sprintf(`"-p %d"`, config.Port))

	syncedLog, err := newSyncedLogger(config.DataPath, os.Stdout)
	if err != nil {
		return err
	}

	postgresProcess.Stdout = syncedLog.file
	postgresProcess.Stderr = syncedLog.file

	if err := postgresProcess.Run(); err != nil {
		_ = syncedLog.flush()
		logContent, _ := readLogsOrTimeout(syncedLog.file)

		return fmt.Errorf("could not stop postgres using %s:\n%s", postgresProcess.String(), string(logContent))
	}

	return nil
}

func defaultCacheLocator(cacheDirectory string, versionStrategy embeddedpostgres.VersionStrategy) embeddedpostgres.CacheLocator {
	return func() (string, bool) {
		operatingSystem, architecture, version := versionStrategy()
		cacheLocation := filepath.Join(cacheDirectory,
			fmt.Sprintf("embedded-postgres-binaries-%s-%s-%s.txz",
				operatingSystem,
				architecture,
				version))

		info, err := os.Stat(cacheLocation)

		if err != nil {
			return cacheLocation, os.IsExist(err) && !info.IsDir()
		}

		return cacheLocation, !info.IsDir()
	}
}

func downloadExtractIfNonexistent(postgresVersion embeddedpostgres.PostgresVersion, binaryRepositoryURL, cacheLocation string) {
	versionStrategy := defaultVersionStrategy(
		postgresVersion,
		runtime.GOOS,
		runtime.GOARCH,
		linuxMachineName,
		shouldUseAlpineLinuxBuild,
	)
	cacheLocator := defaultCacheLocator(cacheLocation, versionStrategy)
	remoteFetchStrategy := defaultRemoteFetchStrategy(binaryRepositoryURL, versionStrategy, cacheLocator)
	fmt.Printf("remoteFetchStrategy(): \"%s\"\n", remoteFetchStrategy())
}

func ensureDirsExist(dirs ...string) error {
	for _, d := range dirs {
		if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(d, fs.ModeDir); err != nil {
				return err
			}
		}
	}
	return nil
}
