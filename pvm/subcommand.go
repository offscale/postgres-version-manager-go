package pvm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"unicode"
)

func EnvSubcommand(config ConfigStruct) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.Username, config.Password, "localhost", config.Port, config.Database)
}

func InstallSubcommand(args Args, cacheLocation string) error {
	var err error
	postgresVersion := args.PostgresVersion
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.LogsPath); err != nil {
		return err
	}
	return downloadExtractIfNonexistent(postgresVersion, args.BinaryRepositoryURL, cacheLocation, args.VersionManagerRoot, false)()
}

func InstallServiceSubcommand(args Args) error {
	var err error
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.LogsPath); err != nil {
		return err
	}
	if args.InstallService.Systemd != nil {
		systemd :=
			fmt.Sprintf(`[Unit]
Description=PostgreSQL %s database server
Documentation=man:postgres(1)
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
User=postgres
ExecStart=%s/postgres -D %s
ExecReload=/bin/kill -HUP $MAINPID
KillMode=mixed
KillSignal=SIGINT
TimeoutSec=infinity

[Install]
WantedBy=multi-user.target
`, args.PostgresVersion, args.BinariesPath, args.DataPath)
		var f *os.File

		f, err = os.Create(args.InstallService.Systemd.ServiceInstallPath)

		if err != nil {
			return err
		}

		defer func(f *os.File) {
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
		_, err = f.WriteString(systemd)

		return err
	} else {
		return fmt.Errorf("NotImplementedError: systemd is the only service available for install")
	}
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
