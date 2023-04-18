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
	var wasLatest bool
	var err error
	var userHomeDir string

	if userHomeDir, err = os.UserHomeDir(); err != nil {
		log.Fatal(err)
	}
	pvm.SetDefaultsFromEnvironment(&args, userHomeDir)
	cacheLocation := path.Join(args.VersionManagerRoot, "downloads")

	arg.MustParse(&args)

	if wasLatest, err = pvm.SetDirectories(&args, userHomeDir); err != nil {
		log.Fatal(err)
	}
	if args.PostgresVersion == "latest" {
		log.Fatalln("latest")
	}

	switch {
	case args.Start != nil:
		if err = pvm.StartSubcommand(args, cacheLocation); err != nil {
			log.Fatal(err)
		}
	case args.Stop != nil:
		if err = pvm.StopPostgres(&args.ConfigStruct); err != nil {
			log.Fatal(err)
		}
	case args.Install != nil:
		if err = pvm.InstallSubcommand(args, wasLatest, cacheLocation); err != nil {
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
