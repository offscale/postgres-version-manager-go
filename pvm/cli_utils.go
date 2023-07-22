package pvm

import (
	"fmt"
	"path/filepath"
)

func SetDefaultsFromEnvironment(args *Args, userHomeDir *string) {
	args.PostgresVersion = "latest"
	args.VersionManagerRoot = filepath.Join(*userHomeDir, "postgres-version-manager")
	versionedRoot := filepath.Join(args.VersionManagerRoot, args.PostgresVersion)
	args.DataPath = filepath.Join(versionedRoot, "data")
	args.RuntimePath = filepath.Join(versionedRoot, "run")
	args.LogsPath = filepath.Join(versionedRoot, "logs")
	args.BinariesPath = filepath.Join(versionedRoot, "bin")
	args.ConfigFile = filepath.Join(args.VersionManagerRoot, "pvm-config.json")
}

func SetVersionAndDirectories(args *Args) error {
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
	latestVersionedRoot := filepath.Join(args.VersionManagerRoot, "latest")
	specificVersionedRoot := filepath.Join(args.VersionManagerRoot, args.PostgresVersion)
	if filepath.Join(latestVersionedRoot, "bin") == args.BinariesPath {
		args.BinariesPath = filepath.Join(specificVersionedRoot, "bin")
	}
	if filepath.Join(latestVersionedRoot, "data") == args.DataPath {
		args.DataPath = filepath.Join(specificVersionedRoot, "data")
	}
	if filepath.Join(latestVersionedRoot, "run") == args.RuntimePath {
		args.RuntimePath = filepath.Join(specificVersionedRoot, "run")
	}
	if filepath.Join(latestVersionedRoot, "logs") == args.LogsPath {
		args.LogsPath = filepath.Join(specificVersionedRoot, "logs")
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
