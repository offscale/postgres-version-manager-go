package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexflint/go-arg"

	"postgres-version-manager-go/pvm"
)

func main() {
	var args pvm.Args
	var err error
	var userHomeDir string

	if userHomeDir, err = os.UserHomeDir(); err != nil {
		log.Fatal(err)
	}
	pvm.SetDefaultsFromEnvironment(&args, userHomeDir)
	cacheLocation := path.Join(args.VersionManagerRoot, "downloads")

	arg.MustParse(&args)

	// Logic to prioritise POSTGRES_VERSION as a positional
	if args.PostgresVersion, err = pvm.PostgresVersionFromLocalOrGlobal(func() string {
		switch {
		case args.Install != nil:
			return args.Install.PostgresVersion
		case args.Start != nil:
			return args.Start.PostgresVersion
		case args.Stop != nil:
			return args.Stop.PostgresVersion
		default:
			return args.PostgresVersion
		}
	}(), args.PostgresVersion); err != nil {
		log.Fatal(err)
	}
	if err = pvm.SetVersionAndDirectories(&args, userHomeDir); err != nil {
		log.Fatal(err)
	}
	if args.PostgresVersion == "latest" {
		log.Fatalln("latest")
	}

	switch {
	case args.Env != nil:
		fmt.Println(pvm.EnvSubcommand(args.ConfigStruct))
	case args.GetDataPath != nil:
		fmt.Println(args.ConfigStruct.DataPath)
	case args.Install != nil:
		if err = pvm.InstallSubcommand(args, cacheLocation); err != nil {
			log.Fatal(err)
		}
	case args.InstallService != nil:
		if err = pvm.InstallServiceSubcommand(args); err != nil {
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
	case args.Start != nil:
		if err = pvm.StartSubcommand(args, cacheLocation); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started PostgreSQL server, access via:\n%s\n", pvm.EnvSubcommand(args.ConfigStruct))
	case args.Stop != nil:
		if err = pvm.StopPostgres(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("missing subcommand")
	}
}
