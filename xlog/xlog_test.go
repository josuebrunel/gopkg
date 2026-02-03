package xlog

import (
	"bytes"
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
			} else {

				if output != "" && strings.Contains(output, tt.msg) {
					t.Errorf("Expected message %q NOT to be logged, got: %q", tt.msg, output)
				}
			}
		})
	}
}
