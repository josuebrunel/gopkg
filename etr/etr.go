// Package etr (Echo Templ Renderer) provides helper functions for rendering templ components with Echo.
package etr

import (
	"context"
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
)

type (
	xcontextkey string
	// QS is a type alias for a query string map.
	QS          = map[string]string
)

var xc xcontextkey = "xcontext"

// Render renders a templ component with the given status code and data.
// It injects context values like request, url, reverse function, and csrf token.
func Render(c *echo.Context, status int, tpl templ.Component, data any) error {
	c.Response().WriteHeader(status)

	var csrf string
	if v := c.Get("csrf"); v != nil {
		csrf = v.(string)
	}
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, xc, map[string]any{
		"request": c.Request(),
		"url":     c.Request().URL.String(),
		"reverse": c.Echo().Router().Routes().Reverse,
		"csrf":    csrf,
		"data":    data,
	})

	err := tpl.Render(ctx, c.Response())
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to render response template")
	}

	return nil
}

// Get retrieves a value from the context injected by Render.
func Get[T any](ctx context.Context, key string) T {
	var cx map[string]any
	if v := ctx.Value(xc); v != nil {
		cx = v.(map[string]any)
	}
	var r T
	if v, ok := cx[key]; ok {
		r = v.(T)
	}
	return r
}

// Reverse generates a URL for a named route using the reverse function in the context.
func Reverse(cx context.Context, name string, values ...any) string {
	reverse := Get[func(string, ...any) (string, error)](cx, "reverse")
	path, err := reverse(name, values...)
	if err != nil {
		xlog.Error("failed to reverse route", "name", name, "values", values, "error", err)
	}
	return path
}

// ReverseX generates a URL for a named route using the Echo context directly.
func ReverseX(c *echo.Context, name string, values ...any) string {
	path, err := c.Echo().Router().Routes().Reverse(name, values...)
	if err != nil {
		xlog.Error("error while reversing route", "name", name, "values", values, "error", err)
	}
	return path
}

// WithQS appends query string parameters to a URL.
func WithQS(url_ string, qs map[string]string) string {
	u, err := url.Parse(url_)
	if err != nil {
		xlog.Error("failed to parse url", "url", url_, "error", err)
		return url_
	}
	q := u.Query()
	for k, v := range qs {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// GetCSRF retrieves the CSRF token from the context.
func GetCSRF(ctx context.Context) string {
	return Get[string](ctx, "csrf")
}
