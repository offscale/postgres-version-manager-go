package pvm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

// SaveConfig will only perform write operations if the file is nonexistent or contains different values
func SaveConfig(args Args) error {
	var err error
	var jsonBytes []byte

	if jsonBytes, err = json.Marshal(args.ConfigStruct); err != nil {
		return err
	}
	if _, err = os.Stat(args.ConfigFile); err == nil {
		var jsonFile *os.File
		if jsonFile, err = os.Open(args.ConfigFile); err != nil {
			return err
		}
		jsonParser := json.NewDecoder(jsonFile)
		var configs ConfigStructs
		if err = jsonParser.Decode(&configs); err != nil {
			return err
		}
		var readConfig *ConfigStruct
		var readConfigPosition int
		for idx, configStruct := range configs {
			if configStruct.PostgresVersion == args.PostgresVersion {
				readConfig = &configStruct
				readConfigPosition = idx
			}
		}
		if readConfig == nil {
			configs = append(configs, args.ConfigStruct)
		} else if !reflect.DeepEqual(*readConfig, args.ConfigStruct) {
			configs[readConfigPosition] = args.ConfigStruct
		} else {
			return err
		}
		if jsonBytes, err = json.Marshal(configs); err != nil {
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		jsonBytes = append(append([]byte(fmt.Sprintf("{\"%s\": ", args.PostgresVersion)), jsonBytes...), []byte("}")...)
	} else {
		return err
	}
	return os.WriteFile(args.ConfigFile, jsonBytes, 0600)
}

func (configs *ConfigStructs) UnmarshalJSON(bytes []byte) error {
	configMap := make(map[string]ConfigStruct)
	if err := json.Unmarshal(bytes, &configMap); err != nil {
		return err
	}
	*configs = make(ConfigStructs, 0, len(configMap))
	for version, c := range configMap {
		c.PostgresVersion = version
		*configs = append(*configs, c)
	}
	return nil
}

func (configs ConfigStructs) MarshalJSON() ([]byte, error) {
	if configs == nil {
		return []byte("null"), nil
	}
	var jsonBody []byte = []byte("{")
	configsLen := len(configs)
	for i, config := range configs {
		var err error
		var configBytes []byte
		if configBytes, err = json.Marshal(config); err != nil {
			return nil, err
		}
		versionKey := fmt.Sprintf("\"%s\":", config.PostgresVersion)
		jsonBody = append(append(jsonBody, []byte(versionKey)...), configBytes...)
		if i != configsLen-1 {
			jsonBody = append(jsonBody, []byte(",")...)
		}
	}
	jsonBody = append(jsonBody, []byte("}")...)
	return jsonBody, nil
}
