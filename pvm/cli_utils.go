package pvm

import (
	"fmt"
	"path"
)

func SetDefaultsFromEnvironment(args *Args, userHomeDir string) {
	args.PostgresVersion = "latest"
	args.VersionManagerRoot = path.Join(userHomeDir, "postgres-version-manager")
	args.BinariesPath = path.Join(args.VersionManagerRoot, args.PostgresVersion)
	args.DataPath = path.Join(args.BinariesPath, "data")
	args.RuntimePath = path.Join(args.BinariesPath, "run")
	args.LogsPath = path.Join(args.BinariesPath, "logs")
	args.ConfigFile = path.Join(args.VersionManagerRoot, "pvm-config.json")
}

func SetVersionAndDirectories(args *Args, userHomeDir string) error {
	if args.PostgresVersion == "latest" {
		if args.NoRemote {
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		} else {
			var err error
			if versionsFromMaven, err = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				return err
			}
			args.PostgresVersion = versionsFromMaven[len(versionsFromMaven)-1]
		}
	} else if !isValidVersion(args.PostgresVersion) {
		return fmt.Errorf("invalid/unsupported PostgreSQL version: %s", args.PostgresVersion)
	}

	// If not provided, use $HOME/postgres-version-manager-go/$POSTGRES_VERSION/
	latestBinariesPath := path.Join(args.VersionManagerRoot, "latest")
	if latestBinariesPath == args.BinariesPath {
		args.BinariesPath = path.Join(userHomeDir, "postgres-version-manager", args.PostgresVersion)
		if path.Join(latestBinariesPath, "data") == args.DataPath {
			args.DataPath = path.Join(args.BinariesPath, "data")
		}
		if path.Join(latestBinariesPath, "run") == args.BinariesPath {
			args.RuntimePath = path.Join(args.BinariesPath, "run")
		}
		if path.Join(latestBinariesPath, "logs") == args.LogsPath {
			args.LogsPath = path.Join(args.BinariesPath, "logs")
		}
	}
	return nil
}

func PostgresVersionFromLocalOrGlobal(localOptionPostgresVersion string, postgresVersion string) (string, error) {
	if localOptionPostgresVersion != "" {
		if localOptionPostgresVersion != "latest" && !isValidVersion(localOptionPostgresVersion) {
			return "", fmt.Errorf("invalid/unsupported PostgreSQL version: %s", localOptionPostgresVersion)
		}
		return localOptionPostgresVersion, nil
	}
	return postgresVersion, nil
}
