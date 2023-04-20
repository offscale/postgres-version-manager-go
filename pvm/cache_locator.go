// Package pvm cache_locator.go from github.com/fergusstrange/embedded-postgres and modified
// (so should be considered under same license)

package pvm

import (
	"fmt"
	"os"
	"path/filepath"
)

// CacheLocator retrieves the location of the Postgres binary cache returning it to location.
// The result of whether this cache is present will be returned to exists.
type CacheLocator func() (location string, exists bool)

func defaultCacheLocator(cacheDirectory string, versionStrategy VersionStrategy) CacheLocator {
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
