// Package pvm prepare_database.go from github.com/fergusstrange/embedded-postgres and modified
// (so should be considered under same license)

package pvm

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lib/pq"
)

const (
	fmtCloseDBConn = "unable to close database connection: %w"
	fmtAfterError  = "%v happened after error: %w"
)

func defaultInitDatabase(binaryExtractLocation, runtimePath, pgDataDir, username, password, locale string, logger *os.File) error {
	passwordFile, err := createPasswordFile(runtimePath, password)
	if err != nil {
		return err
	}

	args := []string{
		"-A", "password",
		"-U", username,
		"-D", pgDataDir,
		fmt.Sprintf("--pwfile=%s", passwordFile),
	}

	if locale != "" {
		args = append(args, fmt.Sprintf("--locale=%s", locale))
	}

	postgresInitDBBinary := filepath.Join(binaryExtractLocation, "bin", "initdb")
	postgresInitDBProcess := exec.Command(postgresInitDBBinary, args...)
	postgresInitDBProcess.Stderr = logger
	postgresInitDBProcess.Stdout = logger

	if err = postgresInitDBProcess.Run(); err != nil {
		logContent, readLogsErr := readLogsOrTimeout(logger) // we want to preserve the original error
		if readLogsErr != nil {
			logContent = []byte(string(logContent) + " - " + readLogsErr.Error())
		}
		return fmt.Errorf("unable to init database using '%s': %w\n%s", postgresInitDBProcess.String(), err, string(logContent))
	}

	if err = os.Remove(passwordFile); err != nil {
		return fmt.Errorf("unable to remove password file '%v': %w", passwordFile, err)
	}

	return nil
}

func createPasswordFile(runtimePath, password string) (string, error) {
	passwordFileLocation := filepath.Join(runtimePath, "pwfile")
	if err := os.WriteFile(passwordFileLocation, []byte(password), 0600); err != nil {
		return "", fmt.Errorf("unable to write password file to %s", passwordFileLocation)
	}

	return passwordFileLocation, nil
}

func defaultCreateDatabase(port uint32, username, password, database string) (err error) {
	if database == "postgres" {
		return nil
	}

	conn, err := openDatabaseConnection(port, username, password, "postgres")
	if err != nil {
		return errorCustomDatabase(database, err)
	}

	db := sql.OpenDB(conn)
	defer func() {
		err = connectionClose(db, err)
	}()
	var rows *sql.Rows
	if rows, err = db.Query("SELECT * FROM pg_database WHERE datname = $1", database); err != nil {
		return err
	}
	if !rows.Next() {
		fmt.Printf("CREATE DATABASE %s\n", database)
		if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
			return errorCustomDatabase(database, err)
		}
	}

	return nil
}

// connectionClose closes the database connection and handles the error of the function that used the database connection
func connectionClose(db io.Closer, err error) error {
	closeErr := db.Close()
	if closeErr != nil {
		closeErr = fmt.Errorf(fmtCloseDBConn, closeErr)

		if err != nil {
			err = fmt.Errorf(fmtAfterError, closeErr, err)
		} else {
			err = closeErr
		}
	}

	return err
}

/*
func healthCheckDatabaseOrTimeout(config ConfigStruct) error {
	healthCheckSignal := make(chan bool)

	defer close(healthCheckSignal)

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)

	defer cancelFunc()

	go func() {
		for timeout.Err() == nil {
			if err := healthCheckDatabase(config.Port, config.Database, config.Username, config.Password); err != nil {
				continue
			}
			healthCheckSignal <- true

			break
		}
	}()

	select {
	case <-healthCheckSignal:
		return nil
	case <-timeout.Done():
		return errors.New("timed out waiting for database to become available")
	}
}

func healthCheckDatabase(port uint32, database, username, password string) (err error) {
	conn, err := openDatabaseConnection(port, username, password, database)
	if err != nil {
		return err
	}

	db := sql.OpenDB(conn)
	defer func() {
		err = connectionClose(db, err)
	}()

	if _, err := db.Query("SELECT 1"); err != nil {
		return err
	}

	return nil
}
*/

func openDatabaseConnection(port uint32, username string, password string, database string) (*pq.Connector, error) {
	conn, err := pq.NewConnector(fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable",
		port,
		username,
		password,
		database))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func errorCustomDatabase(database string, err error) error {
	return fmt.Errorf("unable to connect to create database with custom name %s with the following error: %s", database, err)
}
