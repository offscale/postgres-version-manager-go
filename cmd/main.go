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
		case args.Reload != nil:
			return args.Reload.PostgresVersion
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
	if !args.NoConfigRead && args.Ls == nil && args.LsRemote == nil {
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
	var res string

	switch {
	case args.Env != nil:
		res = pvm.EnvSubcommand(&args.ConfigStruct)
	case args.GetPath != nil:
		res, err = pvm.GetPathSubcommand(&args.GetPath.DirectoryToFind, &args.ConfigStruct)
	case args.Install != nil:
		err = pvm.InstallSubcommand(&args, &cacheLocation)
	case args.InstallService != nil:
		err = pvm.InstallServiceSubcommand(&args)
	case args.Ls != nil:
		err = pvm.LsSubcommand(&args)
	case args.LsRemote != nil:
		err = pvm.LsRemoteSubcommand(&args)
	case args.Ping != nil:
		err = pvm.PingSubcommand(&args.ConfigStruct)
	case args.Reload != nil:
		err = pvm.ReloadPostgres(&args.ConfigStruct)
	case args.Start != nil:
		err = pvm.StartSubcommand(&args, &cacheLocation)
		res = fmt.Sprintf("Started PostgreSQL server, access via:\n%s\n", pvm.EnvSubcommand(&args.ConfigStruct))
	case args.Stop != nil:
		err = pvm.StopPostgres(&args.ConfigStruct)
	case args.Uri != nil:
		res = pvm.UriSubcommand(&args.ConfigStruct)
	default:
		err = fmt.Errorf("missing subcommand, use `--help` to see which subcommand are available")
	}
	if err != nil {
		log.Fatal(err)
	} else if res != "" {
		fmt.Print(res)
	}

	if !args.NoConfigWrite && !nonConfigAlteringSubcommand {
		if err = pvm.SaveConfig(&args); err != nil {
			log.Fatal(err)
		}
	}
}
