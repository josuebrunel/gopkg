package toast

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestConstructors(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(string) Toast
		msg       string
		wantKind  string
		wantTitle string
	}{
		{"Success", Success, "ok", KindSuccess, TitleSuccess},
		{"Error", Error, "fail", KindError, TitleError},
		{"Warning", Warning, "warn", KindWarning, TitleWarning},
		{"Info", Info, "info", KindInfo, TitleInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.msg)
			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", got.Kind, tt.wantKind)
			}
			if got.Title != tt.wantTitle {
				t.Errorf("Title = %v, want %v", got.Title, tt.wantTitle)
			}
			if got.Message != tt.msg {
				t.Errorf("Message = %v, want %v", got.Message, tt.msg)
			}
		})
	}
}

func TestNotify_HTMX(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	Notify(c, KindSuccess, TitleSuccess, "test message")

	got := rec.Header().Get(HeaderKey)
	if got == "" {
		t.Fatal("Header HX-Trigger not set")
	}

	var te ToastEvent
	if err := json.Unmarshal([]byte(got), &te); err != nil {
		t.Fatalf("Failed to unmarshal header value: %v", err)
	}

	if te.Event.Kind != KindSuccess {
		t.Errorf("Kind = %v, want %v", te.Event.Kind, KindSuccess)
	}
	if te.Event.Message != "test message" {
		t.Errorf("Message = %v, want %v", te.Event.Message, "test message")
	}
}

func TestNotify_Cookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No HX-Request header
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	Notify(c, KindError, TitleError, "cookie message")

	res := rec.Result()
	cookies := res.Cookies()

	var toastCookie *http.Cookie
	for _, ck := range cookies {
		if ck.Name == Name {
			toastCookie = ck
			break
		}
	}

	if toastCookie == nil {
		t.Fatalf("Cookie %s not found", Name)
	}

	decodedBytes, err := base64.RawURLEncoding.DecodeString(toastCookie.Value)
	if err != nil {
		t.Fatalf("Failed to decode cookie value: %v", err)
	}

	var tst Toast
	if err := json.Unmarshal(decodedBytes, &tst); err != nil {
		t.Fatalf("Failed to unmarshal cookie json: %v", err)
	}

	if tst.Kind != KindError {
		t.Errorf("Kind = %v, want %v", tst.Kind, KindError)
	}
	if tst.Message != "cookie message" {
		t.Errorf("Message = %v, want %v", tst.Message, "cookie message")
	}
}

func TestUnmarshallToast(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	expected := Toast{
		Kind:    KindInfo,
		Title:   TitleInfo,
		Message: "incoming",
	}
	te := ToastEvent{Event: expected}
	b, _ := json.Marshal(te)

	req.Header.Set(HeaderKey, string(b))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	got := UnmarshallToast(c)

	if got.Kind != expected.Kind {
		t.Errorf("Kind = %v, want %v", got.Kind, expected.Kind)
	}
	if got.Message != expected.Message {
		t.Errorf("Message = %v, want %v", got.Message, expected.Message)
	}
}

func TestUnmarshallToast_Empty(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	got := UnmarshallToast(c)

	if got != ToastNull {
		t.Errorf("Expected ToastNull, got %+v", got)
	}
}

func TestNotifyWrappers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	NotifySuccess(c, "success")
	if rec.Header().Get(HeaderKey) == "" {
		t.Error("NotifySuccess failed to set header")
	}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	NotifyError(c, "error")
	if rec.Header().Get(HeaderKey) == "" {
		t.Error("NotifyError failed to set header")
	}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	NotifyWarning(c, "warning")
	if rec.Header().Get(HeaderKey) == "" {
		t.Error("NotifyWarning failed to set header")
	}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	NotifyInfo(c, "info")
	if rec.Header().Get(HeaderKey) == "" {
		t.Error("NotifyInfo failed to set header")
	}
}
