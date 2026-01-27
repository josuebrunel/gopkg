// Package xsession provides session management middleware and helpers for Echo, backed by alexedwards/scs.
package xsession

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// SessionName is the key used for storing user session data.
const SessionName = "xtoken"

// XUser represents a user in the session.
type XUser struct {
	ID    string            `json:"id"`
	Token string            `json:"token"`
	Email string            `json:"email"`
	Data  map[string]string `json:"data"`
}

// NewXUser creates a new XUser.
func NewXUser(id, token, email string, data map[string]string) XUser {
	if len(data) == 0 {
		data = make(map[string]string)
	}
	return XUser{
		ID: id, Token: token, Email: email,
		Data: data,
	}
}

// SetData sets a key-value pair in the user's data.
func (x *XUser) SetData(key, value string) {
	x.Data[key] = value
}

// GetData retrieves a value from the user's data by key.
func (x *XUser) GetData(key string) string {
	return x.Data[key]
}

// DeleteData removes a key from the user's data.
func (x *XUser) DeleteData(key string) {
	delete(x.Data, key)
}

// SessionConfig defines the configuration for the session middleware.
type SessionConfig struct {
	Skipper        middleware.Skipper
	SessionManager *scs.SessionManager
}

// DefaultSessionConfig is the default configuration for the session middleware.
var (
	DefaultSessionConfig = SessionConfig{
		Skipper: middleware.DefaultSkipper,
	}
	// SessionManager is the global session manager instance.
	SessionManager *scs.SessionManager
)

func init() {
	SessionManager = scs.New()
	SessionManager.Lifetime = 1 * time.Hour
}

// LoadAndSave is a middleware that loads and saves session data.
func LoadAndSave(sessionManager *scs.SessionManager) echo.MiddlewareFunc {
	c := DefaultSessionConfig
	c.SessionManager = sessionManager

	return LoadAndSaveWithConfig(c)
}

// LoadAndSaveWithConfig is a middleware that loads and saves session data with a custom configuration.
func LoadAndSaveWithConfig(config SessionConfig) echo.MiddlewareFunc {

	if config.Skipper == nil {
		config.Skipper = DefaultSessionConfig.Skipper
	}

	if config.SessionManager == nil {
		panic("Session middleware requires a session manager")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			ctx := c.Request().Context()

			var token string
			cookie, err := c.Cookie(config.SessionManager.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = config.SessionManager.Load(ctx, token)
			if err != nil {
				return err
			}

			c.SetRequest(c.Request().WithContext(ctx))

			resp, err := echo.UnwrapResponse(c.Response())
			if err != nil {
				resp.Before(func() {
					if config.SessionManager.Status(ctx) != scs.Unmodified {
						responseCookie := &http.Cookie{
							Name:     config.SessionManager.Cookie.Name,
							Path:     config.SessionManager.Cookie.Path,
							Domain:   config.SessionManager.Cookie.Domain,
							Secure:   config.SessionManager.Cookie.Secure,
							HttpOnly: config.SessionManager.Cookie.HttpOnly,
							SameSite: config.SessionManager.Cookie.SameSite,
						}

						switch config.SessionManager.Status(ctx) {
						case scs.Modified:
							token, _, err := config.SessionManager.Commit(ctx)
							if err != nil {
								panic(err)
							}

							responseCookie.Value = token

						case scs.Destroyed:
							responseCookie.Expires = time.Unix(1, 0)
							responseCookie.MaxAge = -1
						}

						c.SetCookie(responseCookie)
						addHeaderIfMissing(c.Response(), "Cache-Control", `no-cache="Set-Cookie"`)
						addHeaderIfMissing(c.Response(), "Vary", "Cookie")
					}
				})
			}

			return next(c)
		}
	}
}

func addHeaderIfMissing(w http.ResponseWriter, key, value string) {
	for _, h := range w.Header()[key] {
		if h == value {
			return
		}
	}
	w.Header().Add(key, value)
}

// LoginRequired is a middleware that requires the user to be logged in.
// If the user is not logged in, it redirects to the specified path.
func LoginRequired(redirectPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			sess := GetUser(c.Request().Context())
			if strings.EqualFold(sess.Token, "") {
				return c.Redirect(http.StatusFound, redirectPath)
			}
			return next(c)
		}
	}
}

// LoginRequiredFunc is a middleware that requires the user to be logged in.
// If the user is not logged in, it executes the provided function.
func LoginRequiredFunc(fn echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			sess := GetUser(c.Request().Context())
			if strings.EqualFold(sess.Token, "") {
				return fn(c)
			}
			return next(c)
		}
	}
}

// Get retrieves a value from the session context.
func Get[T any](c context.Context, name string) T {
	var r T
	if v, ok := SessionManager.Get(c, name).(T); ok {
		r = v
	}
	return r
}

// Set sets a value in the session context.
func Set(c context.Context, name string, value any) {
	SessionManager.Put(c, name, value)
}

// Delete removes a value from the session context.
func Delete(c context.Context, name string) {
	SessionManager.Remove(c, name)
}

// SetUser stores the user in the session.
func SetUser(c context.Context, u XUser) {
	b, _ := json.Marshal(u)
	Set(c, SessionName, b)
}

// GetUser retrieves the user from the session.
func GetUser(c context.Context) XUser {
	b := Get[[]byte](c, SessionName)
	var u XUser
	json.Unmarshal(b, &u)
	return u
}

// DeleteUser removes the user from the session.
func DeleteUser(c context.Context) {
	Delete(c, SessionName)
}

// IsAuthenticated checks if the user is authenticated (i.e., has an ID).
func IsAuthenticated(c context.Context) bool {
	return GetUser(c).ID != ""
}
