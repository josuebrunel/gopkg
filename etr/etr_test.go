package etr

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josuebrunel/gopkg/assert"
	"github.com/labstack/echo/v5"
)

type MockComponent struct {
	Data []byte
}

func (m MockComponent) Render(ctx context.Context, w io.Writer) error {
	_, err := w.Write(m.Data)
	return err
}

func TestRender(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockComponent := MockComponent{Data: []byte("<div>Hello, World!</div>")}

	err := Render(c, http.StatusOK, mockComponent, nil)
	assert.Eq(t, err, nil)
	assert.Eq(t, http.StatusOK, rec.Code)
	assert.Eq(t, "<div>Hello, World!</div>", rec.Body.String())

	// Test with data and CSRF
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	var (
		csrf = "test_csrf_token"
		data = "some_data"
	)
	c.Set("csrf", "test_csrf_token")

	mockComponentWithData := MockComponent{Data: []byte("Data: " + data + ", CSRF: " + csrf)}

	err = Render(c, http.StatusOK, mockComponentWithData, "some_data")
	assert.Eq(t, err, nil)
	assert.Eq(t, http.StatusOK, rec.Code)
	assert.Eq(t, "Data: some_data, CSRF: test_csrf_token", rec.Body.String())
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, xc, map[string]any{
		"key": "value",
	})

	value := Get[string](ctx, "key")
	assert.Eq(t, "value", value)
}

func TestReverse(t *testing.T) {
	e := echo.New()
	e.AddRoute(echo.Route{
		Name: "getUser", Method: http.MethodGet,
		Path: "/users/:id", Handler: func(c *echo.Context) error { return nil },
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, xc, map[string]any{
		"reverse": e.Router().Routes().Reverse,
	})

	path := Reverse(ctx, "getUser", "123")
	assert.Eq(t, "/users/123", path)

	// Test with non-existent route
	path = Reverse(ctx, "nonExistentRoute")
	assert.Eq(t, "", path)
}

func TestReverseX(t *testing.T) {
	e := echo.New()
	// e.GET("/products/:name", func(c *echo.Context) error { return nil }).Name = "getProduct"
	e.AddRoute(echo.Route{
		Name: "getProduct", Method: http.MethodGet,
		Path: "/products/:name", Handler: func(c *echo.Context) error { return nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	path := ReverseX(c, "getProduct", "apple")
	assert.Eq(t, "/products/apple", path)

	// Test with non-existent route
	path = ReverseX(c, "nonExistentProduct")
	assert.Eq(t, "", path)
}

func TestWithQS(t *testing.T) {
	url_ := "/search"
	qs := map[string]string{
		"q":      "test",
		"filter": "active",
	}

	result := WithQS(url_, qs)
	assert.Eq(t, result, "/search?filter=active&q=test")

	// Test with existing query parameters
	url_ = "/search?page=1"
	qs = map[string]string{
		"q": "test",
	}
	result = WithQS(url_, qs)
	assert.Eq(t, result, "/search?page=1&q=test")

	// Test with empty query string map
	url_ = "/search"
	qs = map[string]string{}
	result = WithQS(url_, qs)
	assert.Eq(t, "/search", result)
}
