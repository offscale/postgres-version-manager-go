//go:build windows

package pvm

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func windowsServiceInstall(args *Args) error {
	var err error
	var m *mgr.Mgr
	if m, err = mgr.Connect(); err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		if err = m.Disconnect(); err != nil {
			panic(err)
		}
	}(m)

	var service *mgr.Service
	if service, err = m.OpenService(args.InstallService.WindowsService.Name); err == nil {
		if err = service.Close(); err != nil {
			return err
		}
		return fmt.Errorf("service %s already exists", args.InstallService.WindowsService.Name)
	}
	pgCtlBinary := filepath.Join(args.BinariesPath, "pg_ctl.exe")
	if service, err = m.CreateService(
		args.InstallService.WindowsService.Name,
		pgCtlBinary,
		mgr.Config{DisplayName: args.InstallService.WindowsService.Description,
			BinaryPathName: fmt.Sprintf(`"%s" runservice -N "%s" -D "%s" -w`,
				pgCtlBinary, args.InstallService.WindowsService.Name, args.DataPath),
		},
		"is", "auto-started",
	); err != nil {
		return err
	}
	defer func(service *mgr.Service) {
		if e := service.Close(); e != nil {
			panic(e)
		}
	}(service)
	if err = eventlog.InstallAsEventCreate(args.InstallService.WindowsService.Name, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
		if e := service.Delete(); e != nil {
			return e
		}
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	fmt.Printf("Created service, to start run:\nnet start %s\n", args.InstallService.WindowsService.Name)
	return nil
}
