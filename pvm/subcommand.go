package pvm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"unicode"
)

func EnvSubcommand(config *ConfigStruct) string {
	val := reflect.ValueOf(*config)
	typ := reflect.TypeOf(*config)
	var envStr string = ""
	for i := 0; i < val.NumField(); i++ {
		valField := val.Field(i)
		typField := typ.Field(i)
		argTag := typField.Tag.Get("arg")
		envName, err := parseEnvFromArgTag(argTag)
		if err != nil {
			envName = typField.Name
		}
		envStr = fmt.Sprintf("%s%s=%c%v%c\n", envStr, envName, CmdQuoteChar, valField.Interface(), CmdQuoteChar)
	}
	return envStr
}

func GetPathSubcommand(directoryToFind *string, config *ConfigStruct) (string, error) {
	switch *directoryToFind {
	case "bin":
		return config.BinariesPath, nil
	case "data":
		return config.DataPath, nil
	case "log":
		return config.LogsPath, nil
	case "runtime":
		return config.RuntimePath, nil
	default:
		return "", fmt.Errorf("usupported path \"%s\"; choose one of: bin, data, log, runtime", directoryToFind)
	}
}

func InstallSubcommand(args *Args, cacheLocation *string) error {
	var err error
	postgresVersion := args.PostgresVersion
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.LogsPath); err != nil {
		return err
	}
	return downloadExtractIfNonexistent(postgresVersion, args.BinaryRepositoryURL, *cacheLocation, args.VersionManagerRoot, false)()
}

func InstallServiceSubcommand(args *Args) error {
	var err error
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.LogsPath); err != nil {
		return err
	}
	if args.InstallService.Systemd != nil {
		systemd :=
			fmt.Sprintf(`[Unit]
Description=PostgreSQL %s database server
After=network.target

[Service]
Type=forking

User=postgres
Group=postgres

OOMScoreAdjust=-1000
Environment=PG_OOM_ADJUST_FILE=/proc/self/oom_score_adj
Environment=PG_OOM_ADJUST_VALUE=0

Environment=PGSTARTTIMEOUT=270

Environment=PGDATA=%s
Environment=PGPORT=%d


ExecStart=%s/pg_ctl start -D ${PGDATA} -s -w -t ${PGSTARTTIMEOUT}
ExecStop=%s/pg_ctl stop -D ${PGDATA} -s -m fast
ExecReload=%s/pg_ctl reload -D ${PGDATA} -s

TimeoutSec=300

[Install]
WantedBy=multi-user.target
`, args.PostgresVersion, args.DataPath, args.Port, args.BinariesPath, args.BinariesPath, args.BinariesPath)
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

func LsSubcommand(args *Args) error {
	var dirs []os.DirEntry
	var err error
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

func LsRemoteSubcommand(args *Args) error {
	var err error
	if args.NoRemote {
		for _, version := range PostgresVersions {
			fmt.Println(version)
		}
	} else {
		if versionsFromMaven == nil {
			if versionsFromMaven, err = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				return err
			}
		}
		for _, version := range versionsFromMaven {
			fmt.Println(version)
		}
	}
	return nil
}

func PingSubcommand(config *ConfigStruct) error {
	conn, err := openDatabaseConnection(config.Port, config.Username, config.Password, config.Database)
	if err != nil {
		return err
	}

	db := sql.OpenDB(conn)
	defer func() {
		if err = connectionClose(db, err); err != nil {
			panic(err)
		}
	}()

	if _, err := db.Query("SELECT 1"); err != nil {
		return err
	}

	return nil
}

func StartSubcommand(args *Args, cacheLocation *string) error {
	var err error
	if err = ensureDirsExist(args.VersionManagerRoot, args.DataPath, args.RuntimePath, args.LogsPath); err != nil {
		return err
	}
	if err = downloadExtractIfNonexistent(args.PostgresVersion, args.BinaryRepositoryURL, *cacheLocation, args.VersionManagerRoot, args.Start.NoInstall)(); err != nil {
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
	return defaultCreateDatabase(args.Port, args.Username, args.Password, args.Database)
}

func UriSubcommand(config *ConfigStruct) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", config.Username, config.Password, "localhost", config.Port, config.Database)
}
