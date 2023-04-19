package pvm

import (
	"errors"
	"fmt"
	"path"

	"github.com/fergusstrange/embedded-postgres"
)

func SetDirectories(args *Args, userHomeDir string) (bool, error) {
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
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		}
	} else if !isValidVersion(string(args.PostgresVersion)) {
		return false, errors.New(fmt.Sprintf("invalid/unsupported PostgreSQL version: %s\n", args.PostgresVersion))
	}

	// If not provided, use $HOME/postgres-version-manager-go/version
	latestBinariesPath := path.Join(args.VersionManagerRoot, "latest")
	if latestBinariesPath == args.BinariesPath {
		args.BinariesPath = path.Join(userHomeDir, "postgres-version-manager", string(args.PostgresVersion))
		if path.Join(latestBinariesPath, "data") == args.DataPath {
			args.DataPath = path.Join(args.BinariesPath, "data")
		}
		if path.Join(latestBinariesPath, "run") == args.RuntimePath {
			args.RuntimePath = path.Join(args.BinariesPath, "run")
		}
	}
	return wasLatest, err
}

func SetDefaultsFromEnvironment(args *Args, userHomeDir string) {
	args.PostgresVersion = "latest"
	args.VersionManagerRoot = path.Join(userHomeDir, "postgres-version-manager")
	args.BinariesPath = path.Join(args.VersionManagerRoot, string(args.PostgresVersion))
	args.DataPath = path.Join(args.BinariesPath, "data")
	args.RuntimePath = path.Join(args.BinariesPath, "run")
}
