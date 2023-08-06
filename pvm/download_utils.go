package pvm

import (
	"runtime"
)

func downloadExtractIfNonexistent(postgresVersion string, binaryRepositoryURL, cacheLocation string, versionManagerRoot string, noExtract bool, noInstall bool) RemoteFetchStrategy {
	versionStrategy := defaultVersionStrategy(
		postgresVersion,
		runtime.GOOS,
		runtime.GOARCH,
		linuxMachineName,
		shouldUseAlpineLinuxBuild,
	)
	cacheLocator := defaultCacheLocator(cacheLocation, versionStrategy)
	return defaultRemoteFetchStrategy(binaryRepositoryURL, versionStrategy, cacheLocator, versionManagerRoot, noExtract, noInstall)
}
