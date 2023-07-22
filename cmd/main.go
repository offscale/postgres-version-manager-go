package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/alexflint/go-arg"
	"github.com/offscale/postgres-version-manager-go/pvm"
)

func main() {
	var args pvm.Args
	var err error
	var userHomeDir string

	if userHomeDir, err = os.UserHomeDir(); err != nil {
		log.Fatal(err)
	}
	pvm.SetDefaultsFromEnvironment(&args, &userHomeDir)
	cacheLocation := path.Join(args.VersionManagerRoot, "downloads")

	arg.MustParse(&args)
	var fieldToNonDefaultValue map[string]interface{} = pvm.FieldAndValueWhenNonDefaultValue(&args.ConfigStruct)

	// Logic to prioritise POSTGRES_VERSION as a positional
	if args.PostgresVersion, err = pvm.PostgresVersionFromLocalOrGlobal(func() string {
		switch {
		case args.Install != nil:
			return args.Install.PostgresVersion
		case args.Ping != nil:
			return args.Ping.PostgresVersion
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
	if err = pvm.SetVersionAndDirectories(&args); err != nil {
		log.Fatal(err)
	}
	if args.PostgresVersion == "latest" {
		log.Fatalln("latest")
	}

	var nonConfigAlteringSubcommand bool = args.Env != nil || args.Ls != nil || args.LsRemote != nil || args.Stop != nil
	if !args.NoConfigRead && (!nonConfigAlteringSubcommand || args.Env != nil || args.Stop != nil) {
		var configStruct *pvm.ConfigStruct
		if configStruct, err = pvm.GetConfigFromFileOfConfigs(&args); err != nil {
			log.Fatal(err)
		}
		if configStruct != nil {
			for field, value := range fieldToNonDefaultValue {
				if err = pvm.SetField(configStruct, field, value); err != nil {
					log.Fatal(err)
				}
			}
			args.ConfigStruct = *configStruct
		}
	}

	switch {
	case args.Env != nil:
		fmt.Println(pvm.EnvSubcommand(&args.ConfigStruct))
	case args.GetPath != nil:
		var pathLocation string
		if pathLocation, err = pvm.GetPathSubcommand(&args.GetPath.DirectoryToFind, &args.ConfigStruct); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(pathLocation)
		}
	case args.Install != nil:
		if err = pvm.InstallSubcommand(&args, &cacheLocation); err != nil {
			log.Fatal(err)
		}
	case args.InstallService != nil:
		if err = pvm.InstallServiceSubcommand(&args); err != nil {
			log.Fatal(err)
		}
	case args.Ls != nil:
		if err = pvm.LsSubcommand(&args); err != nil {
			log.Fatal(err)
		}
	case args.LsRemote != nil:
		if err = pvm.LsRemoteSubcommand(&args); err != nil {
			log.Fatal(err)
		}
	case args.Ping != nil:
		if err = pvm.PingSubcommand(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	case args.Start != nil:
		if err = pvm.StartSubcommand(&args, &cacheLocation); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started PostgreSQL server, access via:\n%s\n", pvm.EnvSubcommand(&args.ConfigStruct))
	case args.Stop != nil:
		fmt.Printf("args: %v\n", args)
		if err = pvm.StopPostgres(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("missing subcommand, use `--help` to see which subcommand are available")
	}

	if !args.NoConfigWrite && !nonConfigAlteringSubcommand {
		if err = pvm.SaveConfig(&args); err != nil {
			log.Fatal(err)
		}
	}
}
