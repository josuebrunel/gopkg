package xenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Options allows for global configuration like prefixes
type Options struct {
	Prefix string
}

func Load(container any) error {
	return LoadWithOptions(container, Options{})
}

// LoadOptions populates a struct using environment variables with an optional prefix.
func LoadWithOptions(container any, opts Options) error {
	val := reflect.ValueOf(container)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("container must be a struct or a pointer to a struct")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		// Recursive call for nested structs
		if field.Kind() == reflect.Struct && field.Type() != reflect.TypeOf(time.Time{}) {
			if err := LoadWithOptions(field.Addr().Interface(), opts); err != nil {
				return err
			}
			continue
		}

		envKey := structField.Tag.Get("env")
		if envKey == "" {
			continue
		}

		// Apply Prefix
		fullKey := opts.Prefix + envKey

		defaultVal := structField.Tag.Get("default")
		required, _ := strconv.ParseBool(structField.Tag.Get("required"))

		valueToSet, exists := os.LookupEnv(fullKey)
		if (!exists || valueToSet == "") && defaultVal != "" {
			valueToSet = defaultVal
		}

		if valueToSet == "" && required {
			return fmt.Errorf("missing required environment variable: %s", fullKey)
		}
		if valueToSet == "" {
			continue
		}

		if err := setField(field, valueToSet); err != nil {
			return fmt.Errorf("field %s: %v", structField.Name, err)
		}
	}
	return nil
}

func setField(field reflect.Value, value string) error {
	// Handle time.Duration specifically
	if field.Type() == reflect.TypeOf(time.Second) {
		d, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(d))
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		iv, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(iv))
	case reflect.Bool:
		bv, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bv)
	case reflect.Slice:
		parts := strings.Split(value, ",")
		newSlice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
		for i, part := range parts {
			if err := setField(newSlice.Index(i), strings.TrimSpace(part)); err != nil {
				return err
			}
		}
		field.Set(newSlice)
	}
	return nil
}
