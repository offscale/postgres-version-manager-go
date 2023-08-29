package pvm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SetDefaultsFromEnvironment(args *Args, userHomeDir *string) {
	if args.PostgresVersion == "" {
		args.PostgresVersion = "latest"
	}
	args.VersionManagerRoot = filepath.Join(*userHomeDir, "postgres-version-manager")
	versionedRoot := filepath.Join(args.VersionManagerRoot, args.PostgresVersion)
	args.DataPath = filepath.Join(versionedRoot, "data")
	args.RuntimePath = filepath.Join(versionedRoot, "run")
	args.LogsPath = filepath.Join(versionedRoot, "logs")
	args.BinariesPath = filepath.Join(versionedRoot, "bin")
	args.ConfigFile = filepath.Join(args.VersionManagerRoot, "pvm-config.json")
}

func SetVersionAndDirectories(args *Args) error {
	var wasLatest bool
	if wasLatest = args.PostgresVersion == "latest"; wasLatest {
		if args.NoRemote {
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		} else {
			var err error
			if versionsFromMaven, err = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				return err
			}
			args.PostgresVersion = versionsFromMaven[len(versionsFromMaven)-1]
		}
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

	if !wasLatest && !isValidVersion(args.PostgresVersion) {
		if _, err := os.Stat(args.BinariesPath); err != nil {
			fmt.Printf("Could not find %s\n", args.BinariesPath)
			return fmt.Errorf("iiii invalid/unsupported PostgreSQL version: %s", args.PostgresVersion)
		}
	}

	return nil
}

func PostgresVersionFromLocalOrGlobal(localOptionPostgresVersion, postgresVersion, BinariesPath string) (string, error) {
	if localOptionPostgresVersion != "" {
		if localOptionPostgresVersion != "latest" && !isValidVersion(localOptionPostgresVersion) {
			if _, err := os.Stat(BinariesPath); err != nil {
				return "", fmt.Errorf("invalid/unsupported PostgreSQL version: %s", localOptionPostgresVersion)
			}
		}
		return localOptionPostgresVersion, nil
	}
	return postgresVersion, nil
}

func parseEnvFromArgTag(argTag string) (string, error) {
	for _, word := range strings.FieldsFunc(argTag, func(c rune) bool {
		return c == ','
	}) {
		if strings.HasPrefix(word, "env:") {
			return word[4:], nil
		}
	}
	return "", fmt.Errorf("no env tag found in struct arg tag")
}
