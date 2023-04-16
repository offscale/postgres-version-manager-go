package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexflint/go-arg"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

var versionsFromMaven []string = nil

// ConfigStruct taken from github.com/offscale/embedded-postgres@v1.20.0/config.go
type ConfigStruct struct {
	PostgresVersion     embeddedpostgres.PostgresVersion `arg:"--postgres-version,env:POSTGRES_VERSION" default:"latest"`
	Port                uint32                           `arg:"-p,env:PGPORT" default:"5432"`
	Database            string                           `arg:"-d,env:POSTGRES_DATABASE" default:"database"`
	Username            string                           `arg:"-u,env:POSTGRES_USERNAME" default:"username"`
	Password            string                           `arg:"env:POSTGRES_PASSWORD" default:"password"`
	RuntimePath         string                           `arg:"--runtime-path,env:RUNTIME_PATH"`
	DataPath            string                           `arg:"--data-path,env:PGDATA"`
	BinariesPath        string                           `arg:"--binary-path,env:BINARY_PATH"`
	Locale              string                           `default:"en_US.UTF-8,env:LANGUAGE"`
	BinaryRepositoryURL string                           `arg:"--binary-repository-url,env:BINARY_REPOSITORY_URL" default:"https://repo1.maven.org/maven2"`
}

type StartCmd struct {
	Pid int `arg:"--pid" default:"-1" help:"If PID provided and exists, will stop that process."`
}

type StopCmd struct {
	Pid int `arg:"--pid" default:"-1"`
}

type InstallCmd struct {
	PostgresVersion embeddedpostgres.PostgresVersion `arg:"positional" default:""`
}

type LsRemoteCmd struct {
}

type EnvCmd struct {
}

func (InstallCmd) Description() string {
	return "this program does this and that"
}

type args struct {
	ConfigStruct
	NoRemote bool         `arg:"--no-remote" default:"false" help:"Disable HTTPS calls for everything except 'install'"`
	Env      *EnvCmd      `arg:"subcommand:env" help:"Print out database connection string"`
	Start    *StartCmd    `arg:"subcommand:start" help:"Start specified PostgreSQL server"`
	Stop     *StopCmd     `arg:"subcommand:stop" help:"Stop specific (running) PostgreSQL server"`
	Install  *InstallCmd  `arg:"subcommand:install" help:"Install specified PostgreSQL version"`
	LsRemote *LsRemoteCmd `arg:"subcommand:ls-remote" help:"List what versions of PostgreSQL are available"`
}

func (args) Description() string {
	return "PostgreSQL version manager"
}

func main() {
	var args args
	var err error
	var userConfigDir string

	wasLatest := false
	args.PostgresVersion = "latest"
	userConfigDir, err = os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	args.BinariesPath = path.Join(userConfigDir, "postgres-version-manager-go", string(args.PostgresVersion))
	args.DataPath = path.Join(args.BinariesPath, "data")
	args.RuntimePath = path.Join(args.BinariesPath, "run")

	arg.MustParse(&args)

	// Need to know what "latest" is if not doing `--ls-remote` and "latest" is PostgresVersion
	if args.PostgresVersion == "latest" {
		wasLatest = true
		if args.NoRemote {
			if err, versionsFromMaven = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				log.Fatal(err)
			}
			args.PostgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
		} else {
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		}
		latestBinariesPath := path.Join(userConfigDir, "postgres-version-manager-go", "latest")
		if latestBinariesPath == args.BinariesPath {
			args.BinariesPath = path.Join(userConfigDir, "postgres-version-manager-go", string(args.PostgresVersion))
			if path.Join(latestBinariesPath, "data") == args.DataPath {
				args.DataPath = path.Join(args.BinariesPath, "data")
			}
			if path.Join(latestBinariesPath, "run") == args.RuntimePath {
				args.RuntimePath = path.Join(args.BinariesPath, "run")
			}
		}
	} else if !isValidVersion(string(args.PostgresVersion)) {
		log.Fatalf("invalid/unsupported PostgreSQL version: %s\n", args.PostgresVersion)
	}

	var config = embeddedpostgres.DefaultConfig().Database(args.Database).Username(args.Username).Password(args.Password).Port(args.Port).BinariesPath(args.BinariesPath).DataPath(args.DataPath).RuntimePath(args.RuntimePath).Version(args.PostgresVersion).Locale(args.Locale)

	switch {
	case args.Start != nil:
		embeddedpostgres.NewDatabase(config)
		if err = startPostgres(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	case args.Stop != nil:
		if err = embeddedpostgres.NewDatabase(config).Stop(); err != nil {
			log.Fatal(err)
		}
	case args.Install != nil:
		if args.Install.PostgresVersion != "" && (args.Install.PostgresVersion != "latest" || !wasLatest) {
			if !isValidVersion(string(args.Install.PostgresVersion)) {
				log.Fatalf("invalid/unsupported PostgreSQL version: %s\n", args.Install.PostgresVersion)
			}
			config = config.Version(args.Install.PostgresVersion)
		}
		embeddedpostgres.NewDatabase(config)
	case args.LsRemote != nil:
		if args.NoRemote {
			for _, version := range PostgresVersions {
				fmt.Println(version)
			}
		} else {
			if versionsFromMaven == nil {
				if err, versionsFromMaven = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
					log.Fatal(err)
				}
			}
			for _, version := range versionsFromMaven {
				fmt.Println(version)
			}
		}
	case args.Env != nil:
		fmt.Printf("postgresql://%s:%s@%s:%d/%s\n", args.Username, args.Password, "localhost", args.Port, args.Database)
	default:
		log.Fatal("missing subcommand")
	}
}
