package pvm

import (
	"runtime"
)

func downloadExtractIfNonexistent(postgresVersion string, binaryRepositoryURL, cacheLocation string, versionManagerRoot string, noInstall bool) RemoteFetchStrategy {
	versionStrategy := defaultVersionStrategy(
		postgresVersion,
		runtime.GOOS,
		runtime.GOARCH,
		linuxMachineName,
		shouldUseAlpineLinuxBuild,
	)
	cacheLocator := defaultCacheLocator(cacheLocation, versionStrategy)
	return defaultRemoteFetchStrategy(binaryRepositoryURL, versionStrategy, cacheLocator, versionManagerRoot, noInstall)
}
