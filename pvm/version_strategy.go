// Package pvm version_strategy.go from github.com/fergusstrange/embedded-postgres and modified
// (so should be considered under same license)
package pvm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func defaultVersionStrategy(postgresVersion embeddedpostgres.PostgresVersion, goos, arch string, linuxMachineName func() string, shouldUseAlpineLinuxBuild func() bool) embeddedpostgres.VersionStrategy {
	return func() (string, string, embeddedpostgres.PostgresVersion) {
		goos := goos
		arch := arch

		if goos == "linux" {
			// the zonkyio/embedded-postgres-binaries project produces
			// arm binaries with the following name schema:
			// 32bit: arm32v6 / arm32v7
			// 64bit (aarch64): arm64v8
			if arch == "arm64" {
				arch += "v8"
			} else if arch == "arm" {
				machineName := linuxMachineName()
				if strings.HasPrefix(machineName, "armv7") {
					arch += "32v7"
				} else if strings.HasPrefix(machineName, "armv6") {
					arch += "32v6"
				}
			}

			if shouldUseAlpineLinuxBuild() {
				arch += "-alpine"
			}
		}

		// postgres below version 14.2 is not available for macOS on arm
		if goos == "darwin" && arch == "arm64" {
			var majorVer, minorVer int
			if _, err := fmt.Sscanf(string(postgresVersion), "%d.%d", &majorVer, &minorVer); err == nil &&
				(majorVer < 14 || (majorVer == 14 && minorVer < 2)) {
				arch = "amd64"
			} else {
				arch += "v8"
			}
		}

		return goos, arch, postgresVersion
	}
}

func linuxMachineName() string {
	var uname string

	if output, err := exec.Command("uname", "-m").Output(); err == nil {
		uname = string(output)
	}

	return uname
}

func shouldUseAlpineLinuxBuild() bool {
	_, err := os.Stat("/etc/alpine-release")
	return err == nil
}
