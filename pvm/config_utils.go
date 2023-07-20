package pvm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
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
		defer func(jsonFile *os.File) {
			err := jsonFile.Close()
			if err != nil {
				panic(err)
			}
		}(jsonFile)

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

func GetConfigFromFileOfConfigs(args Args) (*ConfigStruct, error) {
	var err error
	if _, err = os.Stat(args.ConfigFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var jsonFile *os.File
	if jsonFile, err = os.Open(args.ConfigFile); err != nil {
		return nil, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			panic(err)
		}
	}(jsonFile)
	jsonParser := json.NewDecoder(jsonFile)
	var configs ConfigStructs
	err = jsonParser.Decode(&configs)
	if err != nil {
		return nil, err
	}
	for _, configStruct := range configs {
		if configStruct.PostgresVersion == args.PostgresVersion {
			return &configStruct, err
		}
	}
	return nil, nil
}

func FieldAndValueWhenNonDefaultValue(configStruct ConfigStruct) map[string]interface{} {
	val := reflect.ValueOf(configStruct)
	fieldToValue := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		defaultTag := field.Tag.Get("default")
		fieldValue := fmt.Sprintf("%v", val.Field(i).Interface())

		if defaultTag != fieldValue && defaultTag != "" {
			fieldToValue[field.Name] = fieldValue
		}
	}

	return fieldToValue
}

func SetField(obj interface{}, name string, value interface{}) error {
	// TODO: Handle - Int  Int8  Int16  Int32  Uint  Uint8  Uint16  Uint32  Uintptr  Float32  Float64  Complex64
	// TODO: Handle - Complex128  Array  Chan  Func  Map  Pointer  Slice  Struct  UnsafePointer
	// Current support: interface, string, int64, uint64
	var err error
	if obj == nil {
		return errors.New("object is nil")
	}

	structValue := reflect.ValueOf(obj)
	if structValue.Kind() != reflect.Ptr || !structValue.Elem().IsValid() {
		return fmt.Errorf("object is not a pointer to valid struct")
	}

	structValue = structValue.Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("no such field: %s in struct", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)
	kind := structFieldValue.Kind()

	if kind >= reflect.Int && kind <= reflect.Int64 {
		var i int64
		var ok bool
		i, ok = value.(int64)
		if !ok {
			if i, err = strconv.ParseInt(value.(string), 0, 64); err != nil {
				return err
			}
		}
		if structFieldValue.OverflowInt(i) {
			return fmt.Errorf("provided int value is not valid for field type [%v]", structFieldValue.Type().Kind())
		}

		structFieldValue.SetInt(i)
		return nil
	}

	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		var i uint64
		var ok bool
		i, ok = value.(uint64)
		if !ok {
			if i, err = strconv.ParseUint(value.(string), 0, 64); err != nil {
				return err
			}
		}
		if structFieldValue.OverflowUint(i) {
			return fmt.Errorf("provided uint value is not valid for field type [%v]", structFieldValue.Type().Kind())
		}

		structFieldValue.SetUint(i)
		return nil
	}

	if structFieldValue.Type() != val.Type() {
		return fmt.Errorf("provided value type [%v] didn't match obj field type [%v]", val.Type(), structFieldValue.Type())
	}

	structFieldValue.Set(val)
	return nil
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
