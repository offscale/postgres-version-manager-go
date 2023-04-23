package pvm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"unicode"
)

func InstallSubcommand(args Args, cacheLocation string) error {
	var err error
	postgresVersion := args.PostgresVersion
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.LogsPath); err != nil {
		return err
	}
	return downloadExtractIfNonexistent(postgresVersion, args.BinaryRepositoryURL, cacheLocation, args.VersionManagerRoot, false)()
}

func StartSubcommand(args Args, cacheLocation string) error {
	var err error
	fmt.Printf("args.LogsPath: \"%s\"\n", args.LogsPath)
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.RuntimePath, args.LogsPath); err != nil {
		return err
	}
	if err = downloadExtractIfNonexistent(args.PostgresVersion, args.BinaryRepositoryURL, cacheLocation, args.VersionManagerRoot, args.Start.NoInstall)(); err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(args.DataPath, "pg_wal")); errors.Is(err, os.ErrNotExist) {
		if err = defaultInitDatabase(args.BinariesPath, args.RuntimePath, args.DataPath, args.Username, args.Password, args.Locale, os.Stdout); err != nil {
			return err
		}
	}
	if err = startPostgres(&args.ConfigStruct); err != nil {
		return err
	}
	if err = defaultCreateDatabase(args.Port, args.Username, args.Password, args.Database); err != nil {
		return err
	}
	return nil
}

func LsSubcommand(err error, args Args) error {
	var dirs []os.DirEntry
	dirs, err = os.ReadDir(args.VersionManagerRoot)
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		if dir.Name() != "downloads" && dir.IsDir() && unicode.IsDigit(rune(dir.Name()[0])) {
			fmt.Println(dir.Name())
		}
	}
	return err
}

func LsRemoteSubcommand(args Args) error {
	var err error
	if args.NoRemote {
		for _, version := range PostgresVersions {
			fmt.Println(version)
		}
	} else {
		if versionsFromMaven == nil {
			if err, versionsFromMaven = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				return err
			}
		}
		for _, version := range versionsFromMaven {
			fmt.Println(version)
		}
	}
	return nil
}

func EnvSubcommand(config ConfigStruct) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s\n", config.Username, config.Password, "localhost", config.Port, config.Database)
}
