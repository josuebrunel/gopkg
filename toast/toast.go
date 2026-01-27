// Package toast provides helper functions for sending toast notifications in Echo applications.
package toast

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
)

const (
	KindSuccess  = "success"
	KindError    = "error"
	KindWarning  = "warning"
	KindInfo     = "info"
	TitleSuccess = "Success!"
	TitleError   = "Error!"
	TitleWarning = "Warning!"
	TitleInfo    = "Info!"
	HeaderKey    = "HX-Trigger"
	Name         = "toast"
)

var (
	ToastNull = Toast{}
)

// Toast represents a toast notification message.
type Toast struct {
	Kind    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ToastEvent wraps a Toast for HTMX events.
type ToastEvent struct {
	Event Toast `json:"toast"`
}

// Notify sends a toast notification.
// If the request is an HTMX request, it sets the HX-Trigger header.
// Otherwise, it sets a cookie.
func Notify(c *echo.Context, kind string, title string, msg string) {
	t := ToastEvent{
		Event: Toast{
			Kind:    kind,
			Title:   title,
			Message: msg,
		},
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		notification := marshallFrom(t)
		c.Response().Header().Set(HeaderKey, notification)
		return
	}
	cookie := new(http.Cookie)
	cookie.Name = Name
	cookie.Value = base64.RawURLEncoding.EncodeToString([]byte(marshallFrom(t.Event)))
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(3 * time.Second)
	c.SetCookie(cookie)
}

// Success creates a success toast.
func Success(msg string) Toast {
	return Toast{KindSuccess, TitleSuccess, msg}
}

// Error creates an error toast.
func Error(msg string) Toast {
	return Toast{KindError, TitleError, msg}
}

// Warning creates a warning toast.
func Warning(msg string) Toast {
	return Toast{KindWarning, TitleWarning, msg}
}

// Info creates an info toast.
func Info(msg string) Toast {
	return Toast{KindInfo, TitleInfo, msg}
}

// NotifySuccess sends a success notification.
func NotifySuccess(c *echo.Context, msg string) {
	Notify(c, KindSuccess, TitleSuccess, msg)
}

// NotifyError sends an error notification.
func NotifyError(c *echo.Context, msg string) {
	Notify(c, KindError, TitleError, msg)
}

// NotifyWarning sends a warning notification.
func NotifyWarning(c *echo.Context, msg string) {
	Notify(c, KindWarning, TitleWarning, msg)
}

// NotifyInfo sends an info notification.
func NotifyInfo(c *echo.Context, msg string) {
	Notify(c, KindInfo, TitleInfo, msg)
}

func marshallFrom[T any](t T) string {
	d, err := json.Marshal(t)
	if err != nil {
		xlog.Error("Error marshalling toast", "toast", t, "error", err)
		return ""
	}
	return string(d)
}

// UnmarshallToast extracts a toast from the HTMX header if present.
func UnmarshallToast(c *echo.Context) Toast {
	var t ToastEvent
	if value := c.Request().Header.Get(HeaderKey); value != "" {
		if err := json.Unmarshal([]byte(value), &t); err != nil {
			xlog.Error("Error unmarshalling toast", "toast", value, "error", err)
			return ToastNull
		}
	}
	return t.Event
}
