package pvm

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/fergusstrange/embedded-postgres"
)

func setDirectories(args *args) (bool, error) {
	var userHomeDir string
	var err error
	wasLatest := false

	// Need to know what "latest" is if not doing `--ls-remote` and "latest" is PostgresVersion
	if args.PostgresVersion == "latest" {
		wasLatest = true
		if args.NoRemote {
			if err, versionsFromMaven = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				return false, err
			}
			args.PostgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
		} else {
			fmt.Println("changing")
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		}
	} else if !isValidVersion(string(args.PostgresVersion)) {
		return false, errors.New(fmt.Sprintf("invalid/unsupported PostgreSQL version: %s\n", args.PostgresVersion))
	}

	// If not provided, use $HOME/postgres-version-manager-go/version
	latestBinariesPath := path.Join(args.VersionManagerRoot, "latest")
	if latestBinariesPath == args.BinariesPath {
		args.BinariesPath = path.Join(userHomeDir, "postgres-version-manager-go", string(args.PostgresVersion))
		if path.Join(latestBinariesPath, "data") == args.DataPath {
			args.DataPath = path.Join(args.BinariesPath, "data")
		}
		if path.Join(latestBinariesPath, "run") == args.RuntimePath {
			args.RuntimePath = path.Join(args.BinariesPath, "run")
		}
		fmt.Printf("args.BinariesPath: %s\nargs.DataPath: %s\nargs.RuntimePath: %s\n", args.BinariesPath, args.DataPath, args.RuntimePath)
	} else {
		fmt.Printf("[else] latestBinariesPath != args.BinariesPath, \"%s\" != \"%s\"\n[else] args.BinariesPath: %s\n[else] args.DataPath: %s\n[else] args.RuntimePath: %s\n", latestBinariesPath, args.BinariesPath, args.BinariesPath, args.DataPath, args.RuntimePath)
	}
	return wasLatest, err
}

func setDefaultsFromEnvironment(args *args) (string, error) {
	args.PostgresVersion = "latest"
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	args.VersionManagerRoot = path.Join(userHomeDir, "postgres-version-manager")
	args.BinariesPath = path.Join(args.VersionManagerRoot, string(args.PostgresVersion))
	args.DataPath = path.Join(args.BinariesPath, "data")
	args.RuntimePath = path.Join(args.BinariesPath, "run")

	cacheLocation := path.Join(args.VersionManagerRoot, "downloads")
	return cacheLocation, nil
}
