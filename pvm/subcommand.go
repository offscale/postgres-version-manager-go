package pvm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"
)

func InstallSubcommand(args Args, wasLatest bool, cacheLocation string) error {
	var err error
	postgresVersion := args.PostgresVersion
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath); err != nil {
		return err
	}
	if args.Install.PostgresVersion != "" && (args.Install.PostgresVersion != "latest" || !wasLatest) {
		if !isValidVersion(string(args.Install.PostgresVersion)) {
			return errors.New(fmt.Sprintf("invalid/unsupported PostgreSQL version: %s\n", args.Install.PostgresVersion))
		}
		postgresVersion = args.Install.PostgresVersion
	}
	return downloadExtractIfNonexistent(postgresVersion, args.BinaryRepositoryURL, cacheLocation, args.VersionManagerRoot)()
}

func StartSubcommand(args Args, cacheLocation string) error {
	var err error
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.RuntimePath); err != nil {
		return err
	}
	if err = downloadExtractIfNonexistent(args.PostgresVersion, args.BinaryRepositoryURL, cacheLocation, args.VersionManagerRoot)(); err != nil {
		return err
	}
	if err = defaultInitDatabase(args.BinariesPath, args.RuntimePath, args.DataPath, args.Username, args.Password, args.Locale, os.Stdout); err != nil {
		return err
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
			fmt.Println(dir)
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
