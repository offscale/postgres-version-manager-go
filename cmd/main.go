package cmd

import (
	"fmt"
	"log"

	"github.com/alexflint/go-arg"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	pvm "github.com/offscale/postgres-version-manager-go/pvm"
)

var versionsFromMaven []string = nil

// ConfigStruct originally from github.com/fergusstrange/embedded-postgres@v1.20.0/config.go
type ConfigStruct struct {
	PostgresVersion     embeddedpostgres.PostgresVersion `arg:"--postgres-version,env:POSTGRES_VERSION" default:"latest"`
	Port                uint32                           `arg:"-p,env:PGPORT" default:"5432"`
	Database            string                           `arg:"-d,env:POSTGRES_DATABASE" default:"database"`
	Username            string                           `arg:"-u,env:POSTGRES_USERNAME" default:"username"`
	Password            string                           `arg:"env:POSTGRES_PASSWORD" default:"password"`
	VersionManagerRoot  string                           `arg:"env:VERSION_MANAGER_ROOT"`
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

type LsCmd struct {
}

type LsRemoteCmd struct {
}

type EnvCmd struct {
}

type args struct {
	ConfigStruct
	NoRemote bool         `arg:"--no-remote" default:"false" help:"Disable HTTPS calls for everything except 'install'"`
	Env      *EnvCmd      `arg:"subcommand:env" help:"Print out database connection string"`
	Start    *StartCmd    `arg:"subcommand:start" help:"Start specified PostgreSQL server"`
	Stop     *StopCmd     `arg:"subcommand:stop" help:"Stop specific (running) PostgreSQL server"`
	Install  *InstallCmd  `arg:"subcommand:install" help:"Install specified PostgreSQL version"`
	Ls       *LsCmd       `arg:"subcommand:ls" help:"List what versions of PostgreSQL are installed"`
	LsRemote *LsRemoteCmd `arg:"subcommand:ls-remote" help:"List what versions of PostgreSQL are available"`
}

func (args) Description() string {
	return "PostgreSQL version manager"
}

func main() {
	var args args
	var wasLatest bool
	var err error
	var cacheLocation string

	if cacheLocation, err = pvm.setDefaultsFromEnvironment(&args); err != nil {
		log.Fatal(err)
	}

	arg.MustParse(&args)

	if wasLatest, err = pvm.setDirectories(&args); err != nil {
		log.Fatal(err)
	}
	if args.PostgresVersion == "latest" {
		log.Fatalln("latest")
	}

	switch {
	case args.Start != nil:
		if err = pvm.startSubcommand(args, cacheLocation); err != nil {
			log.Fatal(err)
		}
	case args.Stop != nil:
		if err = pvm.stopPostgres(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	case args.Install != nil:
		if err = pvm.installSubcommand(args, wasLatest, cacheLocation); err != nil {
			log.Fatal(err)
		}
	case args.Ls != nil:
		if err = pvm.LsSubcommand(err, args); err != nil {
			log.Fatal(err)
		}
	case args.LsRemote != nil:
		if err = pvm.LsRemoteSubcommand(args); err != nil {
			log.Fatal(err)
		}
	case args.Env != nil:
		fmt.Printf("postgresql://%s:%s@%s:%d/%s\n", args.Username, args.Password, "localhost", args.Port, args.Database)
	default:
		log.Fatal("missing subcommand")
	}
}
