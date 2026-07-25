package xlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
)

// Color constants
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGray   = "\033[90m"
)

// attrGroup holds the attributes attached (via WithAttrs) while a given
// WithGroup scope was active, so they can be namespaced correctly on output.
type attrGroup struct {
	name  string // empty for the root scope
	attrs []slog.Attr
}

// ColorHandler is a slog.Handler that writes human-readable, color-coded log lines.
type ColorHandler struct {
	opts   slog.HandlerOptions
	w      io.Writer
	mu     *sync.Mutex
	groups []attrGroup
}

// NewColorHandler creates a ColorHandler writing to w.
func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		w:      w,
		opts:   *opts,
		mu:     &sync.Mutex{},
		groups: []attrGroup{{}},
	}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)

	// Time
	if !r.Time.IsZero() {
		buf = append(buf, ColorGray...)
		buf = r.Time.AppendFormat(buf, "15:04:05.000")
		buf = append(buf, ColorReset...)
		buf = append(buf, ' ')
	}

	// Level
	levelColor := ColorReset
	switch r.Level {
	case slog.LevelDebug:
		levelColor = ColorGray
	case slog.LevelInfo:
		levelColor = ColorGreen
	case slog.LevelWarn:
		levelColor = ColorYellow
	case slog.LevelError:
		levelColor = ColorRed
	}
	buf = append(buf, levelColor...)
	buf = append(buf, r.Level.String()...)
	buf = append(buf, ColorReset...)
	buf = append(buf, ' ')

	// Source
	if h.opts.AddSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			buf = append(buf, ColorGray...)
			buf = append(buf, fmt.Sprintf("%s:%d", f.File, f.Line)...)
			buf = append(buf, ColorReset...)
			buf = append(buf, ' ')
		}
	}

	// Message
	buf = append(buf, r.Message...)

	// Attributes accumulated via WithAttrs, namespaced by the WithGroup scope
	// that was active when they were added.
	prefix := ""
	for _, g := range h.groups {
		if g.name != "" {
			prefix = joinKey(prefix, g.name)
		}
		for _, a := range g.attrs {
			buf = h.appendAttr(buf, joinKey(prefix, a.Key), a.Value)
		}
	}
	// The record's own attributes belong to the innermost active group.
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, joinKey(prefix, a.Key), a.Value)
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf)
	return err
}

func (h *ColorHandler) appendAttr(buf []byte, key string, v slog.Value) []byte {
	buf = append(buf, ' ')
	buf = append(buf, ColorGray...)
	buf = append(buf, key...)
	buf = append(buf, '=')
	buf = append(buf, ColorReset...)
	buf = append(buf, fmt.Sprint(v.Resolve().Any())...)
	return buf
}

func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// WithAttrs returns a new handler with attrs attached to the current group scope.
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	newGroups := make([]attrGroup, len(h.groups))
	copy(newGroups, h.groups)
	last := newGroups[len(newGroups)-1]
	last.attrs = append(append([]slog.Attr{}, last.attrs...), attrs...)
	newGroups[len(newGroups)-1] = last
	return &ColorHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		groups: newGroups,
	}
}

// WithGroup returns a new handler that namespaces subsequently added attributes under name.
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroups := make([]attrGroup, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = attrGroup{name: name}
	return &ColorHandler{
		opts:   h.opts,
		w:      h.w,
		mu:     h.mu,
		groups: newGroups,
	}
}
