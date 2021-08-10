package confighelper

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Type int

const (
	ENV Type = iota
	YAML
	// CONF
)

type Config interface {
}

func LoadConfig(from Type, file string, to Config) {
	configDefault(to)

	switch from {
	case ENV:
		configFromEnv(to)
	case YAML:
		configFromYAML(file, to)
	default:
		panic("config type unsupported now")
	}
}

func configDefault(to Config) {
	v := reflect.ValueOf(to).Elem()
	if v.Kind() == reflect.Invalid {
		v = reflect.New(reflect.TypeOf(to).Elem()).Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		defaultValue := t.Field(i).Tag.Get("default")
		if defaultValue != "" {
			switch t.Field(i).Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value, err := strconv.Atoi(defaultValue); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetInt(int64(value))
				}

			case reflect.Float32:
				if value, err := strconv.ParseFloat(defaultValue, 32); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetFloat(float64(value))
				}

			case reflect.Float64:
				if value, err := strconv.ParseFloat(defaultValue, 64); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetFloat(float64(value))
				}

			case reflect.Bool:
				v.FieldByName(t.Field(i).Name).SetBool(defaultValue == "true")

			case reflect.String:
				v.FieldByName(t.Field(i).Name).SetString(defaultValue)

			}
		}
	}
}

func configFromEnv(to Config) {
	v := reflect.ValueOf(to).Elem()
	if v.Kind() == reflect.Invalid {
		v = reflect.New(reflect.TypeOf(to).Elem()).Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		configName := t.Field(i).Tag.Get("config")
		configValue := os.Getenv(configName)
		if configValue != "" {
			switch t.Field(i).Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value, err := strconv.Atoi(configValue); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetInt(int64(value))
				}

			case reflect.Float32:
				if value, err := strconv.ParseFloat(configValue, 32); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetFloat(float64(value))
				}

			case reflect.Float64:
				if value, err := strconv.ParseFloat(configValue, 64); err != nil {
					panic(err.Error())
				} else {
					v.FieldByName(t.Field(i).Name).SetFloat(float64(value))
				}

			case reflect.Bool:
				v.FieldByName(t.Field(i).Name).SetBool(configValue == "true")

			case reflect.String:
				v.FieldByName(t.Field(i).Name).SetString(configValue)

			}
		}
	}
}

func configFromYAML(file string, to Config) {
	if f, err := os.Open(file); err != nil {
		defer f.Close()
		panic(err.Error())
	} else {
		defer f.Close()
		var raw map[string]interface{}
		if err := yaml.NewDecoder(f).Decode(&raw); err != nil {
			panic(err.Error())
		}
		if data, err := json.Marshal(raw); err != nil {
			panic(err.Error())
		} else {
			if err := json.Unmarshal(data, to); err != nil {
				panic(err.Error())
			}
		}

		v := reflect.ValueOf(to).Elem()
		if v.Kind() == reflect.Invalid {
			v = reflect.New(reflect.TypeOf(to).Elem()).Elem()
		}
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			configName := t.Field(i).Tag.Get("config")
			if configValue, ok := raw[configName]; ok && configValue != "" {
				switch t.Field(i).Type.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v.FieldByName(t.Field(i).Name).SetInt(int64(configValue.(int)))

				case reflect.Float32, reflect.Float64:
					v.FieldByName(t.Field(i).Name).SetFloat(configValue.(float64))

				case reflect.Bool:
					v.FieldByName(t.Field(i).Name).SetBool(configValue.(bool))

				case reflect.String:
					v.FieldByName(t.Field(i).Name).SetString(configValue.(string))

				}
			}
		}
	}
}
