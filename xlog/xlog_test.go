package xlog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestSetup(t *testing.T) {
	// Refactored to test Setup and Level logic via output capture
	var buf bytes.Buffer

	tests := []struct {
		name        string
		config      Config
		logFunc     func(string, ...any)
		shouldPrint bool
		msg         string
	}{
		{
			name:        "Debug_Printed_When_Debug",
			config:      Config{Level: "DEBUG", Output: &buf},
			logFunc:     Debug,
			shouldPrint: true,
			msg:         "debug msg",
		},
		{
			name:        "Debug_Not_Printed_When_Info",
			config:      Config{Level: "INFO", Output: &buf},
			logFunc:     Debug,
			shouldPrint: false,
			msg:         "debug hidden",
		},
		{
			name:        "Info_Printed_When_Info",
			config:      Config{Level: "INFO", Output: &buf},
			logFunc:     Info,
			shouldPrint: true,
			msg:         "info msg",
		},
		{
			name:   "String_Case_Insensitive",
			config: Config{Level: "debug", Output: &buf},
			logFunc: func(s string, a ...any) {
				Debug(s, a...)
			},
			shouldPrint: true,
			msg:         "case insensitive",
		},
		{
			name:   "Int_Fallback",
			config: Config{Level: "-4", Output: &buf}, // -4 is Debug in slog
			logFunc: func(s string, a ...any) {
				Debug(s, a...)
			},
			shouldPrint: true,
			msg:         "int fallback",
		},
		{
			name:   "Text_Format",
			config: Config{Level: "INFO", Output: &buf, Format: FormatText},
			logFunc: func(s string, a ...any) {
				Info(s, a...)
			},
			shouldPrint: true,
			msg:         "text_format_test",
		},
		{
			name:   "Disable_Source",
			config: Config{Level: "INFO", Output: &buf, DisableSource: true},
			logFunc: func(s string, a ...any) {
				Info(s, a...)
			},
			shouldPrint: true,
			msg:         "sourceless",
		},
		{
			name:   "Color_Output",
			config: Config{Level: "INFO", Output: &buf, Color: true},
			logFunc: func(s string, a ...any) {
				Info(s, a...)
			},
			shouldPrint: true,
			msg:         "color_output_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			Setup(tt.config)
			tt.logFunc(tt.msg)

			output := buf.String()
			if tt.shouldPrint {
				if !strings.Contains(output, tt.msg) {
					t.Errorf("Expected output to contain %q, got: %q", tt.msg, output)
				}
				// Additional check for Text format structure if needed
				if tt.config.Format == FormatText {
					if !strings.Contains(output, "level=INFO") {
						t.Errorf("Expected text format to contain level=INFO, got: %q", output)
					}
				}
				if tt.config.Color {
					want := ColorGreen + "INFO" + ColorReset
					if !strings.Contains(output, want) {
						t.Errorf("Expected color output to contain %q, got: %q", want, output)
					}
				}
			} else {

				if output != "" && strings.Contains(output, tt.msg) {
					t.Errorf("Expected message %q NOT to be logged, got: %q", tt.msg, output)
				}
			}
		})
	}
}

func TestColorHandlerWithAttrsAndGroup(t *testing.T) {
	var buf bytes.Buffer
	h := NewColorHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	h = h.WithAttrs([]slog.Attr{slog.String("service", "gopkg")}).(*ColorHandler)
	h = h.WithGroup("req").(*ColorHandler)

	logger := slog.New(h)
	logger.Info("hello", "id", 42)

	output := buf.String()
	if want := ColorGray + "service" + "=" + ColorReset + "gopkg"; !strings.Contains(output, want) {
		t.Errorf("expected persistent attr from WithAttrs to survive WithGroup, got: %q", output)
	}
	if want := ColorGray + "req.id" + "=" + ColorReset + "42"; !strings.Contains(output, want) {
		t.Errorf("expected record attr to be namespaced under group, got: %q", output)
	}
}
