package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type SimpleHandler struct {
	w    *os.File
	opts *slog.HandlerOptions
}

func NewSimpleHandler(w *os.File, opts *slog.HandlerOptions) *SimpleHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &SimpleHandler{w: w, opts: opts}
}

func (h *SimpleHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return lvl >= minLevel
}

func (h *SimpleHandler) Handle(_ context.Context, r slog.Record) error {
	ts := r.Time.Format(time.DateTime + ".000")

	levelStr := strings.ToUpper(r.Level.String())
	levelStr = fmt.Sprintf("%-5s", levelStr)

	line := fmt.Sprintf("[%s] %s %s", ts, levelStr, r.Message)

	r.Attrs(func(a slog.Attr) bool {
		val := a.Value.Any()

		if errVal, ok := val.(error); ok {
			line += fmt.Sprintf(" Trace: %s", errVal.Error())
		} else {
			line += fmt.Sprintf(" %s=%v", a.Key, val)
		}

		return true
	})

	_, err := fmt.Fprintln(h.w, line)
	return err
}

func (h *SimpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *SimpleHandler) WithGroup(name string) slog.Handler       { return h }
