package pvm

import (
	"runtime"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func downloadExtractIfNonexistent(postgresVersion embeddedpostgres.PostgresVersion, binaryRepositoryURL, cacheLocation string, versionManagerRoot string) embeddedpostgres.RemoteFetchStrategy {
	versionStrategy := defaultVersionStrategy(
		postgresVersion,
		runtime.GOOS,
		runtime.GOARCH,
		linuxMachineName,
		shouldUseAlpineLinuxBuild,
	)
	cacheLocator := defaultCacheLocator(cacheLocation, versionStrategy)
	return defaultRemoteFetchStrategy(binaryRepositoryURL, versionStrategy, cacheLocator, versionManagerRoot)
}
