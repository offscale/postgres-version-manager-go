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
	Port                uint32                           `arg:"-p,env:POSTGRES_PORT" default:"5432"`
	Database            string                           `arg:"-d,env:POSTGRES_DATABASE" default:"database"`
	Username            string                           `arg:"-u,env:POSTGRES_USERNAME" default:"username"`
	Password            string                           `arg:"env:POSTGRES_PASSWORD" default:"password"`
	RuntimePath         string                           `arg:"--runtime-path,env:RUNTIME_PATH"`
	DataPath            string                           `arg:"--data-path,env:DATA_PATH"`
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

type args struct {
	ConfigStruct
	NoRemote bool         `arg:"--no-remote" default:"false"`
	Env      *EnvCmd      `arg:"subcommand:env"`
	Start    *StartCmd    `arg:"subcommand:start"`
	Stop     *StopCmd     `arg:"subcommand:stop"`
	Install  *InstallCmd  `arg:"subcommand:install"`
	LsRemote *LsRemoteCmd `arg:"subcommand:ls-remote"`
}

func (args) Description() string {
	return "PostgreSQL version manager"
}

func main() {
	var args args
	args.PostgresVersion = "latest"
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	args.BinariesPath = path.Join(userConfigDir, "postgres-version-manager-go", string(args.PostgresVersion))
	args.DataPath = path.Join(args.BinariesPath, "data")
	args.RuntimePath = path.Join(args.BinariesPath, "run")

	arg.MustParse(&args)

	// Need to know what "latest" is if not doing `--ls-remote` and "latest" is PostgresVersion
	if args.PostgresVersion == "latest" {
		if args.NoRemote {
			var err error
			if err, versionsFromMaven = getVersionsFromMaven(args.BinaryRepositoryURL); err != nil {
				log.Fatal(err)
			}
			args.PostgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
		} else {
			args.PostgresVersion = PostgresVersions[NumberOfPostgresVersions-1]
		}
		if path.Join(userConfigDir, "postgres-version-manager-go", "latest") == args.BinariesPath {
			args.BinariesPath = path.Join(userConfigDir, "postgres-version-manager-go", string(args.PostgresVersion))
			args.DataPath = path.Join(args.BinariesPath, "data")
			args.RuntimePath = path.Join(args.BinariesPath, "run")
		}
	} else if !isValidVersion(string(args.PostgresVersion)) {
		log.Fatalf("invalid/unsupported PostgreSQL version: %s", args.PostgresVersion)
	}

	var config = embeddedpostgres.DefaultConfig().Database(args.Database).Username(args.Username).Password(args.Password).Port(args.Port).BinariesPath(args.BinariesPath).DataPath(args.DataPath).RuntimePath(args.RuntimePath).Version(args.PostgresVersion)

	switch {
	case args.Start != nil:
		if err := embeddedpostgres.NewDatabase(config).Start(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("postgresql://%s:%s@%s:%d/%s\n", args.Username, args.Password, "localhost", args.Port, args.Database)
	case args.Stop != nil:
		if err := embeddedpostgres.NewDatabase(config).Stop(); err != nil {
			log.Fatal(err)
		}
	case args.Install != nil:
		embeddedpostgres.NewDatabase(config)
	case args.LsRemote != nil:
		if args.NoRemote {
			for _, version := range PostgresVersions {
				fmt.Println(version)
			}
		} else {
			var err error
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
	}
}
