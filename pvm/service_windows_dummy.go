//go:build !windows

package pvm

import (
	"fmt"
	"runtime"
)

func windowsServiceInstall(args *Args) error {
	return fmt.Errorf("windowsServiceInstall is only available on Windows, not %s", runtime.GOOS)
}
