package main

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/lib/pq"
)

const (
	fmtCloseDBConn = "unable to close database connection: %w"
	fmtAfterError  = "%v happened after error: %w"
)

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

func createDatabase(port uint32, username, password, database string) (err error) {
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

	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
		return errorCustomDatabase(database, err)
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
