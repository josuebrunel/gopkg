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

type Toast struct {
	Kind    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ToastEvent struct {
	Event Toast `json:"toast"`
}

func Notify(c echo.Context, kind string, title string, msg string) {
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

func Success(msg string) Toast {
	return Toast{KindSuccess, TitleSuccess, msg}
}

func Error(msg string) Toast {
	return Toast{KindError, TitleError, msg}
}

func Warning(msg string) Toast {
	return Toast{KindWarning, TitleWarning, msg}
}

func Info(msg string) Toast {
	return Toast{KindInfo, TitleInfo, msg}
}

func NotifySuccess(c echo.Context, msg string) {
	Notify(c, KindSuccess, TitleSuccess, msg)
}

func NotifyError(c echo.Context, msg string) {
	Notify(c, KindError, TitleError, msg)
}

func NotifyWarning(c echo.Context, msg string) {
	Notify(c, KindWarning, TitleWarning, msg)
}

func NotifyInfo(c echo.Context, msg string) {
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

func UnmarshallToast(c echo.Context) Toast {
	var t ToastEvent
	if value := c.Request().Header.Get(HeaderKey); value != "" {
		if err := json.Unmarshal([]byte(value), &t); err != nil {
			xlog.Error("Error unmarshalling toast", "toast", value, "error", err)
			return ToastNull
		}
	}
	return t.Event
}
