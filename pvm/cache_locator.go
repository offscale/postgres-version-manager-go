package pvm

import (
	"fmt"
	"os"
	"path/filepath"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func defaultCacheLocator(cacheDirectory string, versionStrategy embeddedpostgres.VersionStrategy) embeddedpostgres.CacheLocator {
	return func() (string, bool) {
		operatingSystem, architecture, version := versionStrategy()
		cacheLocation := filepath.Join(cacheDirectory,
			fmt.Sprintf("embedded-postgres-binaries-%s-%s-%s.txz",
				operatingSystem,
				architecture,
				version))

		info, err := os.Stat(cacheLocation)

		if err != nil {
			return cacheLocation, os.IsExist(err) && !info.IsDir()
		}

		return cacheLocation, !info.IsDir()
	}
}
