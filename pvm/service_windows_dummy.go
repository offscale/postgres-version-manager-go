//go:build !windows

package pvm

import (
	"fmt"
)

func windowsServiceInstall(args *Args) error {
	return fmt.Errorf("windowsServiceInstall is only available on Windows")
}
