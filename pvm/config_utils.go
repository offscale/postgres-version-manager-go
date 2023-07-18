package pvm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// SaveConfig will only save the config if the file is nonexistent or contains different values
func SaveConfig(args Args) error {
	var err error
	var jsonBytes []byte
	if jsonBytes, err = json.MarshalIndent(args.ConfigStruct, "", ""); err != nil {
		return err
	}
	if _, err = os.Stat(args.ConfigFile); err == nil {
		// Check if
		var jsonFile *os.File
		if jsonFile, err = os.Open(args.ConfigFile); err != nil {
			return err
		}
		jsonParser := json.NewDecoder(jsonFile)
		var readConfig ConfigStruct
		if err = jsonParser.Decode(&readConfig); err != nil {
			return err
		}
		if !cmp.Equal(readConfig, args.ConfigStruct, cmpopts.IgnoreFields(ConfigStruct{}, "PostgresVersion", "ConfigFile", "NoConfigRw")) {
			log.Fatal("TODO: Implement replacement of just this version of PostgreSQL's version pvm-config")
		}
	} else if errors.Is(err, os.ErrNotExist) {
		jsonBytes = append(append([]byte(fmt.Sprintf("{\"%s\": ", args.PostgresVersion)), jsonBytes...), []byte("}")...)
	} else {
		return err
	}
	return os.WriteFile(args.ConfigFile, jsonBytes, 0600)
}
