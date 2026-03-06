// Package xenv provides functionality to load environment variables into struct fields.
package xenv

import (
	"bufio"
	"fmt"
	"io"
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

// Load populates a struct using environment variables.
func Load(container any) error {
	return LoadWithOptions(container, Options{})
}

// LoadWithOptions populates a struct using environment variables with an optional prefix.
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

// LoadEnvFile loads environment variables from the given file path into the
// process environment and then populates container using those variables.
// Lines starting with '#' and blank lines are ignored.
// Values may include an inline comment separated by " #".
func LoadEnvFile(path string, container any) error {
	return LoadEnvFileWithOptions(path, container, Options{})
}

// LoadEnvFileWithOptions is like LoadEnvFile but accepts Options (e.g. a prefix).
func LoadEnvFileWithOptions(path string, container any, opts Options) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("xenv: open env file: %w", err)
	}
	defer f.Close()

	if err := parseEnvFile(f); err != nil {
		return err
	}
	return LoadWithOptions(container, opts)
}

// parseEnvFile reads key=value pairs from r and sets them in the environment.
func parseEnvFile(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		// strip inline comment (" #…") before quote removal
		if idx := strings.Index(value, " #"); idx != -1 {
			value = strings.TrimSpace(value[:idx])
		}
		// strip surrounding quotes
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("xenv: setenv %s: %w", key, err)
		}
	}
	return scanner.Err()
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
