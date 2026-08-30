# gopkg

A collection of useful Go packages and mini-libraries for building web applications and services.

## Packages

### assert

A minimal, dependency-free, non-fatal assertion library for testing.

```go
import "github.com/josuebrunel/gopkg/assert"

func TestSomething(t *testing.T) {
    assert.Eq(t, 1, 1)
}
```

See [assert/README.md](assert/README.md) for the full API and more examples.

### component

Reusable UI components built with [templ](https://templ.guide/). Includes basic elements like buttons, inputs, forms, and layouts.

### errorsmap

A helper package for managing multiple errors, useful for form validation or aggregating errors.

```go
import "github.com/josuebrunel/gopkg/errorsmap"

func validate() error {
    errs := errorsmap.New()
    errs["username"] = errors.New("invalid username")
    if errs.Nil() {
        return nil
    }
    return errs
}
```

### etr (Echo Templ Renderer)

Helper functions for rendering [templ](https://templ.guide/) components with [Echo](https://echo.labstack.com/), managing context, CSRF tokens, and reverse routing.

```go
import "github.com/josuebrunel/gopkg/etr"

func Handler(c echo.Context) error {
    return etr.Render(c, http.StatusOK, MyComponent(), data)
}
```

### pbc (PocketBase Client)

A wrapper around the PocketBase API, providing a convenient client for interacting with collections, records, and authentication.

```go
import "github.com/josuebrunel/gopkg/pbc"

client := pbc.New("http://127.0.0.1:8090")
resp, err := client.RecordList("posts")
```

### spbautherror

A utility to parse and handle authentication errors, specifically designed for parsing JSON error responses embedded in error strings (e.g., from Supabase or similar services).

### toast

Toast notification helpers for Echo applications. Supports setting notifications via Cookies or HTMX `HX-Trigger` headers.

```go
import "github.com/josuebrunel/gopkg/toast"

func Handler(c echo.Context) error {
    toast.NotifySuccess(c, "Operation successful!")
    return c.String(http.StatusOK, "OK")
}
```

### xenv

A lightweight library for populating struct fields from environment variables.

See [xenv/README.md](xenv/README.md) for full documentation.

### xlog

A wrapper around `log/slog` that provides structured logging with automatic source file and line number annotation.

```go
import "github.com/josuebrunel/gopkg/xlog"

// Optional: Configure the logger (defaults to INFO, JSON, Source enabled)
xlog.Setup(xlog.Config{
    Level:         "DEBUG",         // "DEBUG", "INFO", "WARN", "ERROR"
    Output:        os.Stdout,
    Format:        xlog.FormatText, // xlog.FormatJSON or xlog.FormatText
    DisableSource: false,           // Set to true to hide file/line info
    SourceDepth:   7,               // Default depth for stack trace
    Color:         true,            // Enable color output (forces text format)
})

xlog.Info("Something happened", "key", "value")
```

### xsession

Session management middleware and helpers for Echo, backed by `alexedwards/scs`. Handles user session storage and retrieval.

```go
// Setup middleware
e.Use(xsession.LoadAndSave(sessionManager))

// Get user from context
user := xsession.GetUser(ctx)
```
