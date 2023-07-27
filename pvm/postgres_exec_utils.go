// Package pvm these functions from github.com/fergusstrange/embedded-postgres and modified
// (so should be considered under same license)
package pvm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ReloadPostgres(config *ConfigStruct) error {
	var err error
	postgresBinary := filepath.Join(config.BinariesPath, "pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "reload",
		"-D", config.DataPath, "-s")

	var syncedLog *syncedLogger

	if syncedLog, err = newSyncedLogger(config.LogsPath, os.Stdout); err != nil {
		return err
	}

	postgresProcess.Stdout = syncedLog.file
	postgresProcess.Stderr = syncedLog.file

	if err = postgresProcess.Run(); err != nil {
		_ = syncedLog.flush()
		logContent, _ := readLogsOrTimeout(syncedLog.file)

		return fmt.Errorf("could not reload postgres using %s:\n%s", postgresProcess.String(), string(logContent))
	}

	return nil
}

func startPostgres(config *ConfigStruct) error {
	postgresBinary := filepath.Join(config.BinariesPath, "pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "start", "-w",
		"-D", config.DataPath,
		"-o", fmt.Sprintf(`"-p %d"`, config.Port))

	syncedLog, err := newSyncedLogger(config.LogsPath, os.Stdout)
	if err != nil {
		return err
	}

	postgresProcess.Stdout = syncedLog.file
	postgresProcess.Stderr = syncedLog.file

	if err := postgresProcess.Run(); err != nil {
		_ = syncedLog.flush()
		logContent, _ := readLogsOrTimeout(syncedLog.file)

		return fmt.Errorf("could not start postgres using %s:\n%s", postgresProcess.String(), string(logContent))
	}

	return nil
}

func StopPostgres(config *ConfigStruct) error {
	postgresBinary := filepath.Join(config.BinariesPath, "pg_ctl")
	postgresProcess := exec.Command(postgresBinary, "stop", "-w",
		"-D", config.DataPath,
		"-o", fmt.Sprintf(`"-p %d"`, config.Port))

	syncedLog, err := newSyncedLogger(config.LogsPath, os.Stdout)
	if err != nil {
		return err
	}

	postgresProcess.Stdout = syncedLog.file
	postgresProcess.Stderr = syncedLog.file

	if err := postgresProcess.Run(); err != nil {
		_ = syncedLog.flush()
		logContent, _ := readLogsOrTimeout(syncedLog.file)

		return fmt.Errorf("could not stop postgres using %s:\n%s", postgresProcess.String(), string(logContent))
	}

	return nil
}
