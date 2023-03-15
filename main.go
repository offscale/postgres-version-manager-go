package main

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
	cli "github.com/urfave/cli/v2"
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

func main() {
	var versionsFromMaven []string = nil

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "install",
				Usage:       "Specific postgres version",
				Value:       "latest",
				DefaultText: "latest",
				Action: func(ctx *cli.Context, v string) error {
					var postgresVersion embeddedpostgres.PostgresVersion
					if v == "latest" {
						if versionsFromMaven == nil {
							versionsFromMaven = getVersionsFromMaven()
						}
						postgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
						goto validVersion
					} else if isValidVersion(v) {
						postgresVersion = embeddedpostgres.PostgresVersion(v)
						goto validVersion
					}

					return fmt.Errorf("invalid version: %s\n", v)
				validVersion:
					return postgresConfigThen(postgresVersion, false, false)
				},
			},
			&cli.StringFlag{
				Name:        "start",
				Usage:       "Specific postgres version",
				Value:       "latest",
				DefaultText: "latest",
				Action: func(ctx *cli.Context, v string) error {
					var postgresVersion embeddedpostgres.PostgresVersion
					if v == "latest" {
						if versionsFromMaven == nil {
							versionsFromMaven = getVersionsFromMaven()
						}
						postgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
						goto validVersion
					} else if isValidVersion(v) {
						postgresVersion = embeddedpostgres.PostgresVersion(v)
						goto validVersion
					}
					return fmt.Errorf("invalid version: %s\n", v)
				validVersion:
					return postgresConfigThen(postgresVersion, true, false)
				},
			},
			&cli.StringFlag{
				Name:        "stop",
				Usage:       "Specific postgres version",
				Value:       "latest",
				DefaultText: "latest",
				Action: func(ctx *cli.Context, v string) error {
					var postgresVersion embeddedpostgres.PostgresVersion
					if v == "latest" {
						if versionsFromMaven == nil {
							versionsFromMaven = getVersionsFromMaven()
						}
						postgresVersion = embeddedpostgres.PostgresVersion(versionsFromMaven[len(versionsFromMaven)-1])
						goto validVersion
					} else if isValidVersion(v) {
						postgresVersion = embeddedpostgres.PostgresVersion(v)
						goto validVersion
					}
					return fmt.Errorf("invalid version: %s\n", v)
				validVersion:
					return postgresConfigThen(postgresVersion, false, true)
				},
			},
			&cli.BoolFlag{
				Name:  "ls-remote",
				Usage: "list available versionsFromMaven",
				Action: func(ctx *cli.Context, v bool) error {
					if v {
						for _, version := range PostgresVersions {
							fmt.Println(version)
						}
					}
					return nil
				},
			},
			&cli.UintFlag{
				Name:  "port",
				Usage: "specify port to start/stop PostgreSQL on",
				Value: 5432,
				// EnvVars: string[]{"PORT"},
				Action: nil,
			},
		},
		Description: "PostgreSQL version manager",
		HideHelp:    false, // I have no idea why `--help` is always printing
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
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

func postgresConfigThen(version embeddedpostgres.PostgresVersion, start bool, stop bool) error {
	const database string = "postgres"
	const username string = "postgres"
	const password string = "postgres"
	const port uint32 = 5432
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	var binariesPath string = path.Join(userConfigDir, "postgres-version-manager-go", string(version))
	var dataPath string = path.Join(binariesPath, "data")
	var runtimePath string = path.Join(binariesPath, "run")
	if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dataPath, os.ModePerm); err != nil {
			return err
		}
		fmt.Printf("No issue making dataPath = \"%s\"\n", dataPath)
	}
	const locale string = "en_US.UTF-8"
	const binaryRepositoryURL string = "https://repo1.maven.org/maven2"
	const startTimeout time.Duration = 15 * time.Second // nanosecond
	var logger io.Writer = os.Stdout
	var embeddedPostgres *embeddedpostgres.EmbeddedPostgres = embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().Version(version).Port(port).Database(database).Username(username).Password(password).RuntimePath(runtimePath).DataPath(dataPath).BinariesPath(binariesPath).Locale(locale).BinaryRepositoryURL(binaryRepositoryURL).StartTimeout(startTimeout).Logger(logger),
	)
	if start {
		if err := embeddedPostgres.Start(); err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("RDBMS_URI=\"postgresql://%s:%s@%s:%d/%s\"\n\n", username, password, "localhost", port, database)
		}
	} else if stop {
		if err := embeddedPostgres.Stop(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
