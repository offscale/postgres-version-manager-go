package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	cli "github.com/urfave/cli/v2"
)

var POSTGRES_VERSIONS [117]embeddedpostgres.PostgresVersion = [117]embeddedpostgres.PostgresVersion{
	"15.2.0",
	"15.1.0",
	"15.0.0",
	"14.7.0",
	"14.6.0",
	"14.5.0",
	"14.4.0",
	"14.3.0",
	"14.2.0",
	"14.1.0",
	"14.0.0",
	"13.10.0",
	"13.9.0",
	"13.8.0",
	"13.7.0",
	"13.6.0",
	"13.5.0",
	"13.4.0",
	"13.3.0",
	"13.2.0",
	"13.1.0-1",
	"13.1.0",
	"13.0.0",
	"12.14.0",
	"12.13.0",
	"12.12.0",
	"12.11.0",
	"12.10.0",
	"12.9.0",
	"12.8.0",
	"12.7.0",
	"12.6.0",
	"12.5.0-1",
	"12.5.0",
	"12.4.0",
	"12.3.0",
	"12.2.0",
	"12.1.0-1",
	"12.1.0",
	"12.0.0",
	"11.19.0",
	"11.18.0",
	"11.17.0",
	"11.16.0",
	"11.15.0",
	"11.14.0",
	"11.13.0",
	"11.12.0",
	"11.11.0",
	"11.10.0-1",
	"11.10.0",
	"11.9.0",
	"11.8.0",
	"11.7.0",
	"11.6.0-1",
	"11.6.0",
	"11.5.0",
	"11.4.0",
	"11.3.0",
	"11.2.0",
	"11.1.0",
	"11.0.0",
	"10.23.0",
	"10.22.0",
	"10.21.0",
	"10.20.0",
	"10.19.0",
	"10.18.0",
	"10.17.0",
	"10.16.0",
	"10.15.0-1",
	"10.15.0",
	"10.14.0",
	"10.13.0",
	"10.12.0",
	"10.11.0-1",
	"10.11.0",
	"10.10.0",
	"10.9.0",
	"10.8.0",
	"10.7.0",
	"10.6.0",
	"10.5.0-1",
	"10.5.0",
	"10.4.0",
	"9.6.24",
	"9.5.25",
	"9.5.24-1",
	"9.5.24",
	"9.5.23",
	"9.5.22",
	"9.5.21",
	"9.5.20-1",
	"9.5.20",
	"9.5.19",
	"9.5.18",
	"9.5.17",
	"9.5.16",
	"9.5.15",
	"9.5.14-1",
	"9.5.14",
	"9.5.13",
	"9.4.26",
	"9.4.25-1",
	"9.4.25",
	"9.4.24",
	"9.4.23",
	"9.4.22",
	"9.4.21",
	"9.4.20",
	"9.4.19-1",
	"9.4.19",
	"9.4.18",
	"9.3.25",
	"9.3.24-1",
	"9.3.24",
	"9.3.23"}

// Micro-optimization to avoid looping
func isValidVersion(v string) bool {
	switch v {
	case
		"15.2.0",
		"15.1.0",
		"15.0.0",
		"14.7.0",
		"14.6.0",
		"14.5.0",
		"14.4.0",
		"14.3.0",
		"14.2.0",
		"14.1.0",
		"14.0.0",
		"13.10.0",
		"13.9.0",
		"13.8.0",
		"13.7.0",
		"13.6.0",
		"13.5.0",
		"13.4.0",
		"13.3.0",
		"13.2.0",
		"13.1.0-1",
		"13.1.0",
		"13.0.0",
		"12.14.0",
		"12.13.0",
		"12.12.0",
		"12.11.0",
		"12.10.0",
		"12.9.0",
		"12.8.0",
		"12.7.0",
		"12.6.0",
		"12.5.0-1",
		"12.5.0",
		"12.4.0",
		"12.3.0",
		"12.2.0",
		"12.1.0-1",
		"12.1.0",
		"12.0.0",
		"11.19.0",
		"11.18.0",
		"11.17.0",
		"11.16.0",
		"11.15.0",
		"11.14.0",
		"11.13.0",
		"11.12.0",
		"11.11.0",
		"11.10.0-1",
		"11.10.0",
		"11.9.0",
		"11.8.0",
		"11.7.0",
		"11.6.0-1",
		"11.6.0",
		"11.5.0",
		"11.4.0",
		"11.3.0",
		"11.2.0",
		"11.1.0",
		"11.0.0",
		"10.23.0",
		"10.22.0",
		"10.21.0",
		"10.20.0",
		"10.19.0",
		"10.18.0",
		"10.17.0",
		"10.16.0",
		"10.15.0-1",
		"10.15.0",
		"10.14.0",
		"10.13.0",
		"10.12.0",
		"10.11.0-1",
		"10.11.0",
		"10.10.0",
		"10.9.0",
		"10.8.0",
		"10.7.0",
		"10.6.0",
		"10.5.0-1",
		"10.5.0",
		"10.4.0",
		"9.6.24",
		"9.5.25",
		"9.5.24-1",
		"9.5.24",
		"9.5.23",
		"9.5.22",
		"9.5.21",
		"9.5.20-1",
		"9.5.20",
		"9.5.19",
		"9.5.18",
		"9.5.17",
		"9.5.16",
		"9.5.15",
		"9.5.14-1",
		"9.5.14",
		"9.5.13",
		"9.4.26",
		"9.4.25-1",
		"9.4.25",
		"9.4.24",
		"9.4.23",
		"9.4.22",
		"9.4.21",
		"9.4.20",
		"9.4.19-1",
		"9.4.19",
		"9.4.18",
		"9.3.25",
		"9.3.24-1",
		"9.3.24",
		"9.3.23":
		return true
	}
	return false
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "install",
				Usage:       "Specific postgres version",
				Value:       "latest",
				DefaultText: "latest",
				Action: func(ctx *cli.Context, v string) error {
					if v == "latest" {
						var version embeddedpostgres.PostgresVersion = POSTGRES_VERSIONS[0]
						fmt.Printf("latest::installing: %s\n", version)
						return config_then(version, false)
					} else if isValidVersion(v) {
						var version embeddedpostgres.PostgresVersion = embeddedpostgres.PostgresVersion(v)
						fmt.Printf("        installing: %s\n", version)
						return config_then(version, false)
					}
					return fmt.Errorf("invalid version: %s\n", v)
				},
			},
			&cli.StringFlag{
				Name:        "start",
				Usage:       "Specific postgres version",
				Value:       "latest",
				DefaultText: "latest",
				Action: func(ctx *cli.Context, v string) error {
					if v == "latest" {
						var version embeddedpostgres.PostgresVersion = POSTGRES_VERSIONS[0]
						fmt.Printf("latest::starting: %s\n", version)
						return config_then(version, true)
					} else if isValidVersion(v) {
						var version embeddedpostgres.PostgresVersion = embeddedpostgres.PostgresVersion(v)
						fmt.Printf("        starting: %s\n", version)
						return config_then(version, true)
					}
					return fmt.Errorf("invalid version: %s\n", v)
				},
			},
			&cli.BoolFlag{
				Name:  "ls-remote",
				Usage: "list available versions",
				Action: func(ctx *cli.Context, v bool) error {
					if v {
						for _, version := range POSTGRES_VERSIONS {
							fmt.Println(version)
						}
					}
					return nil
				},
			},
		},
		HideHelp: false, // I have no idea why `--help` is always printing
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	//install()
}

func config_then(version embeddedpostgres.PostgresVersion, start bool) error {
	const port uint32 = 5433
	const database string = "postgres"
	const username string = "postgres"
	const password string = "postgres"
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	var versionBasePath string = path.Join(userConfigDir, "postgres-version-manager-go", string(version))
	var dataPath string = path.Join(versionBasePath, "data")
	var runtimePath string = dataPath
	var binariesPath string = path.Join(versionBasePath, "bin")
	if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			return err
		}
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
			fmt.Printf("RDBMS_URI=\"postgresql://%s:%s@%s:%d/%s\n\"", username, password, "localhost", port, database)
		}
	}

	//if stop {
	//	if err := embeddedPostgres.Stop(); err != nil {
	//		log.Fatal(err)
	//	}
	//}
	return nil
}
