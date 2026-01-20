# xenv

`xenv` is a lightweight Go library for populating struct fields from environment variables. It supports default values, required fields, nested structs, and type conversion.

## Features

- **Struct Mapping**: Map environment variables to struct fields using tags.
- **Default Values**: Specify default values for optional environment variables.
- **Required Fields**: Enforce the presence of critical environment variables.
- **Type Support**: Handles `string`, `int`, `bool`, `slice` (comma-separated), and `time.Duration`.
- **Prefixing**: Option to apply a global prefix to environment variable lookups (e.g., `APP_`).
- **Nested Structs**: Recursively populates nested structs.

## Installation

```bash
go get github.com/yosuke/gopkg/xenv
```

## Usage

### Basic Example

Define your configuration struct with `env`, `default`, or `required` tags.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/josuebrunel/gopkg/xenv"
)

type Config struct {
	Host     string        `env:"HOST" default:"localhost"`
	Port     int           `env:"PORT" default:"8080"`
	Debug    bool          `env:"DEBUG"`
	Timeout  time.Duration `env:"TIMEOUT" default:"30s"`
	Tags     []string      `env:"TAGS"` // Comma-separated values
	APIKey   string        `env:"API_KEY" required:"true"`
}

func main() {
	var cfg Config

	// Load from environment variables
	if err := xenv.Load(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Config: %+v\n", cfg)
}
```

### Using Prefixes

You can use `LoadWithOptions` to specify a prefix that is prepended to all environment variable keys defined in the tags.

```go
opts := xenv.Options{Prefix: "APP_"}

// Will look for APP_HOST, APP_PORT, etc.
xenv.LoadWithOptions(&cfg, opts)
```

## Supported Tags

| Tag        | Description                                                                               | Example           |
| :--------- | :---------------------------------------------------------------------------------------- | :---------------- |
| `env`      | The name of the environment variable to lookup.                                           | `env:"PORT"`      |
| `default`  | The value to use if the environment variable is missing or empty.                         | `default:"8080"`  |
| `required` | If set to `true`, returns an error if the variable is missing and no default is provided. | `required:"true"` |