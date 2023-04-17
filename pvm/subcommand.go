package pvm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func installSubcommand(args args, wasLatest bool, cacheLocation string) error {
	var err error
	var postgresVersion embeddedpostgres.PostgresVersion = args.PostgresVersion
	if err = ensureDirsExist(args.ConfigStruct.VersionManagerRoot, args.ConfigStruct.DataPath); err != nil {
		return err
	}
	if args.Install.PostgresVersion != "" && (args.Install.PostgresVersion != "latest" || !wasLatest) {
		if !isValidVersion(string(args.Install.PostgresVersion)) {
			return errors.New(fmt.Sprintf("invalid/unsupported PostgreSQL version: %s\n", args.Install.PostgresVersion))
		}
		postgresVersion = args.Install.PostgresVersion
	}
	return downloadExtractIfNonexistent(postgresVersion, args.ConfigStruct.BinaryRepositoryURL, cacheLocation)()
}

func startSubcommand(args args, cacheLocation string) error {
	var err error
	if err = ensureDirsExist(args.ConfigStruct.VersionManagerRoot, args.ConfigStruct.DataPath); err != nil {
		return err
	}
	if err = downloadExtractIfNonexistent(args.PostgresVersion, args.ConfigStruct.BinaryRepositoryURL, cacheLocation)(); err != nil {
		return err
	}
	if err = startPostgres(&args.ConfigStruct); err != nil {
		return err
	}
	return createDatabase(args.Port, args.Username, args.Password, args.Database)
}

func LsSubcommand(err error, args args) error {
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

func LsRemoteSubcommand(args args) error {
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
