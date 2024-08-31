package envelope

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadFromEnv(ptr interface{}) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ptr must be a pointer to a struct")
	}
	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		envName := field.Tag.Get("env")
		if envName == "" {
			continue // Skip if the tag isn't set.
		}

		envValue := os.Getenv(envName)
		if envValue == "" {
			continue // Skip if the environment variable is not set.
		}

		if err := setFieldValue(fieldValue, envValue); err != nil {
			return fmt.Errorf("error while setting field %s from environment variable %s: %w", field.Name, envName, err)
		}
	}

	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse %s as an integer: %w", value, err)
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("unable to parse %s as an unsigned integer: %w", value, err)
		}
		field.SetUint(uintValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unable to parse %s as an boolean: %w", value, err)
		}
		field.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("unable to parse %s as an float: %w", value, err)
		}
		field.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}
