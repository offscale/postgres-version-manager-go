package main

/*
TODO: fix root config reflection into subcommands
TODO: add `start` and `stop` subcommands
TODO: add positional arg to override version in `install`, `start`, and `stop`
*/

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/mkideal/cli"
)

const NumberOfPostgresVersions uintptr = 135

var PostgresVersions [NumberOfPostgresVersions]embeddedpostgres.PostgresVersion = [NumberOfPostgresVersions]embeddedpostgres.PostgresVersion{
	"9.3.23",
	"9.3.24",
	"9.3.24-1",
	"9.3.25",
	"9.4.18",
	"9.4.19",
	"9.4.19-1",
	"9.4.20",
	"9.4.21",
	"9.4.22",
	"9.4.23",
	"9.4.24",
	"9.4.25",
	"9.4.25-1",
	"9.4.26",
	"9.5.13",
	"9.5.14",
	"9.5.14-1",
	"9.5.15",
	"9.5.16",
	"9.5.17",
	"9.5.18",
	"9.5.19",
	"9.5.20",
	"9.5.20-1",
	"9.5.21",
	"9.5.22",
	"9.5.23",
	"9.5.24",
	"9.5.24-1",
	"9.5.25",
	"9.6.9",
	"9.6.10",
	"9.6.10-1",
	"9.6.11",
	"9.6.12",
	"9.6.13",
	"9.6.14",
	"9.6.15",
	"9.6.16",
	"9.6.16-1",
	"9.6.17",
	"9.6.18",
	"9.6.19",
	"9.6.20",
	"9.6.20-1",
	"9.6.21",
	"9.6.22",
	"9.6.23",
	"9.6.24",
	"10.4.0",
	"10.5.0",
	"10.5.0-1",
	"10.6.0",
	"10.7.0",
	"10.8.0",
	"10.9.0",
	"10.10.0",
	"10.11.0",
	"10.11.0-1",
	"10.12.0",
	"10.13.0",
	"10.14.0",
	"10.15.0",
	"10.15.0-1",
	"10.16.0",
	"10.17.0",
	"10.18.0",
	"10.19.0",
	"10.20.0",
	"10.21.0",
	"10.22.0",
	"10.23.0",
	"11.0.0",
	"11.1.0",
	"11.2.0",
	"11.3.0",
	"11.4.0",
	"11.5.0",
	"11.6.0",
	"11.6.0-1",
	"11.7.0",
	"11.8.0",
	"11.9.0",
	"11.10.0",
	"11.10.0-1",
	"11.11.0",
	"11.12.0",
	"11.13.0",
	"11.14.0",
	"11.15.0",
	"11.16.0",
	"11.17.0",
	"11.18.0",
	"11.19.0",
	"12.0.0",
	"12.1.0",
	"12.1.0-1",
	"12.2.0",
	"12.3.0",
	"12.4.0",
	"12.5.0",
	"12.5.0-1",
	"12.6.0",
	"12.7.0",
	"12.8.0",
	"12.9.0",
	"12.10.0",
	"12.11.0",
	"12.12.0",
	"12.13.0",
	"12.14.0",
	"13.0.0",
	"13.1.0",
	"13.1.0-1",
	"13.2.0",
	"13.3.0",
	"13.4.0",
	"13.5.0",
	"13.6.0",
	"13.7.0",
	"13.8.0",
	"13.9.0",
	"13.10.0",
	"14.0.0",
	"14.1.0",
	"14.2.0",
	"14.3.0",
	"14.4.0",
	"14.5.0",
	"14.6.0",
	"14.7.0",
	"15.0.0",
	"15.1.0",
	"15.2.0"}

// Micro-optimization to avoid looping
func isValidVersion(v string) bool {
	switch v {
	case
		"9.3.23",
		"9.3.24",
		"9.3.24-1",
		"9.3.25",
		"9.4.18",
		"9.4.19",
		"9.4.19-1",
		"9.4.20",
		"9.4.21",
		"9.4.22",
		"9.4.23",
		"9.4.24",
		"9.4.25",
		"9.4.25-1",
		"9.4.26",
		"9.5.13",
		"9.5.14",
		"9.5.14-1",
		"9.5.15",
		"9.5.16",
		"9.5.17",
		"9.5.18",
		"9.5.19",
		"9.5.20",
		"9.5.20-1",
		"9.5.21",
		"9.5.22",
		"9.5.23",
		"9.5.24",
		"9.5.24-1",
		"9.5.25",
		"9.6.9",
		"9.6.10",
		"9.6.10-1",
		"9.6.11",
		"9.6.12",
		"9.6.13",
		"9.6.14",
		"9.6.15",
		"9.6.16",
		"9.6.16-1",
		"9.6.17",
		"9.6.18",
		"9.6.19",
		"9.6.20",
		"9.6.20-1",
		"9.6.21",
		"9.6.22",
		"9.6.23",
		"9.6.24",
		"10.4.0",
		"10.5.0",
		"10.5.0-1",
		"10.6.0",
		"10.7.0",
		"10.8.0",
		"10.9.0",
		"10.10.0",
		"10.11.0",
		"10.11.0-1",
		"10.12.0",
		"10.13.0",
		"10.14.0",
		"10.15.0",
		"10.15.0-1",
		"10.16.0",
		"10.17.0",
		"10.18.0",
		"10.19.0",
		"10.20.0",
		"10.21.0",
		"10.22.0",
		"10.23.0",
		"11.0.0",
		"11.1.0",
		"11.2.0",
		"11.3.0",
		"11.4.0",
		"11.5.0",
		"11.6.0",
		"11.6.0-1",
		"11.7.0",
		"11.8.0",
		"11.9.0",
		"11.10.0",
		"11.10.0-1",
		"11.11.0",
		"11.12.0",
		"11.13.0",
		"11.14.0",
		"11.15.0",
		"11.16.0",
		"11.17.0",
		"11.18.0",
		"11.19.0",
		"12.0.0",
		"12.1.0",
		"12.1.0-1",
		"12.2.0",
		"12.3.0",
		"12.4.0",
		"12.5.0",
		"12.5.0-1",
		"12.6.0",
		"12.7.0",
		"12.8.0",
		"12.9.0",
		"12.10.0",
		"12.11.0",
		"12.12.0",
		"12.13.0",
		"12.14.0",
		"13.0.0",
		"13.1.0",
		"13.1.0-1",
		"13.2.0",
		"13.3.0",
		"13.4.0",
		"13.5.0",
		"13.6.0",
		"13.7.0",
		"13.8.0",
		"13.9.0",
		"13.10.0",
		"14.0.0",
		"14.1.0",
		"14.2.0",
		"14.3.0",
		"14.4.0",
		"14.5.0",
		"14.6.0",
		"14.7.0",
		"15.0.0",
		"15.1.0",
		"15.2.0":
		return true
	}
	return false
}

var versionsFromMaven []string = nil

/*
<metadata>
  <groupId>io.zonky.test.postgres</groupId>
  <artifactId>embedded-postgres-binaries-bom</artifactId>
  <versioning>
    <latest>15.2.0</latest>
    <release>15.2.0</release>
	<versions>

*/

type Metadata struct {
	XMLName    xml.Name `xml:"metadata"`
	Text       string   `xml:",chardata"`
	GroupId    string   `xml:"groupId"`
	ArtifactId string   `xml:"artifactId"`
	Versioning struct {
		Text     string `xml:",chardata"`
		Latest   string `xml:"latest"`
		Release  string `xml:"release"`
		Versions struct {
			Text    string   `xml:",chardata"`
			Version []string `xml:"version"`
		} `xml:"versions"`
		LastUpdated string `xml:"lastUpdated"`
	} `xml:"versioning"`
}

// ConfigCmd taken from github.com/fergusstrange/embedded-postgres@v1.20.0/config.go
type ConfigCmd struct {
	Version             string `cli:"version" dft:"latest"`
	Port                uint32 `cli:"port" dft:"5432"`
	Database            string `cli:"database" dft:"database"`
	Username            string `cli:"username" dft:"username"`
	Password            string `cli:"password" dft:"password"`
	RuntimePath         string `cli:"runtime-path"`
	DataPath            string `cli:"data-path"`
	BinariesPath        string `cli:"binaries-path"`
	Locale              string `cli:"locale" dft:"en_US.UTF-8"`
	BinaryRepositoryURL string `cli:"binary-repository-url" dft:"https://repo1.maven.org/maven2"`
}

type LsCmd struct {
	NoRemote bool `cli:"no-remote" dft:"true"`
}

type InstallStruct struct {
	cli.Helper
}

var installCommand = &cli.Command{
	Name: "install",
	Desc: "install specified version of PostgreSQL",
	Argv: func() interface{} { return new(InstallStruct) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*InstallStruct)
		ctx.String("[install] argv.database = %s\n", argv)
		return postgresConfigThen(config, false, false)
	},
}

var lsCommand = &cli.Command{
	Name: "ls-remote",
	Desc: "list available versionsFromMaven",
	Argv: func() interface{} { return new(LsCmd) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*LsCmd)
		if argv.NoRemote {
			for _, version := range PostgresVersions {
				fmt.Println(version)
			}
			return nil
		}
		if versionsFromMaven == nil {
			versionsFromMaven = getVersionsFromMaven()
		}
		for _, version := range versionsFromMaven {
			fmt.Println(version)
		}
		return nil
	},
}

type rootCli struct {
	cli.Helper
	ConfigCmd
}

type LsCli struct {
	cli.Helper
	LsCmd
}

var config = embeddedpostgres.DefaultConfig()

var root = &cli.Command{
	Argv: func() interface{} { return new(rootCli) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootCli)

		var binariesPath string
		var dataPath string
		var runtimePath string
		var postgresVersion string
		if argv.Version == "latest" {
			if versionsFromMaven == nil {
				versionsFromMaven = getVersionsFromMaven()
			}
			postgresVersion = versionsFromMaven[len(versionsFromMaven)-1]
		}
		if argv.BinariesPath == "" {
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				return err
			}
			binariesPath = path.Join(userConfigDir, "postgres-version-manager-go", postgresVersion)
		} else {
			binariesPath = argv.BinariesPath
		}
		if argv.DataPath == "" {
			dataPath = path.Join(binariesPath, "data")
		} else {
			dataPath = argv.DataPath
		}
		if argv.DataPath == "" {
			runtimePath = path.Join(binariesPath, "run")
		} else {
			runtimePath = argv.RuntimePath
		}
		if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
			if err = os.MkdirAll(dataPath, os.ModePerm); err != nil {
				return err
			}
			fmt.Printf("No issue making dataPath = \"%s\"\n", dataPath)
		}
		fmt.Printf("runtimePath: \"%s\"\n", runtimePath)

		config = config.Database(argv.Database).Username(argv.Username).Password(argv.Password).Port(argv.Port).BinariesPath(binariesPath).DataPath(dataPath).RuntimePath(runtimePath).Version(embeddedpostgres.PostgresVersion(postgresVersion))
		return nil
	},
}

func main() {
	if err := cli.Root(root,
		cli.Tree(installCommand),
		cli.Tree(lsCommand),
	).Run(os.Args[1:]); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}

func getVersionsFromMaven() []string {
	var metadata Metadata
	{
		resp, err := http.Get("https://repo1.maven.org/maven2/io/zonky/test/postgres/embedded-postgres-binaries-bom/maven-metadata.xml")
		if err != nil {
			panic(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err := xml.Unmarshal(body, &metadata); err != nil {
			panic(err)
		}
	}
	return metadata.Versioning.Versions.Version
}

func postgresConfigThen(config embeddedpostgres.Config, start bool, stop bool) error {
	const startTimeout time.Duration = 15 * time.Second // nanosecond
	var logger io.Writer = os.Stdout
	var embeddedPostgres *embeddedpostgres.EmbeddedPostgres = embeddedpostgres.NewDatabase(
		config.StartTimeout(startTimeout).Logger(logger),
	)
	if start {
		if err := embeddedPostgres.Start(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("RDBMS_URI=\"postgresql://%s:%s@%s:%d/%s\"\n\n", "TODO", "password", "localhost", "", "")
		}
	} else if stop {
		if err := embeddedPostgres.Stop(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
