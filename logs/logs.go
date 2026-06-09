package logs

import (
	"context"
	"log/slog"
	"os"
)

func InitLogger() {
	stderrHandler := NewSimpleHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	stdoutHandler := NewSimpleHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	//need to pass in higher level handlers first!
	handler := newMultiHandler(stderrHandler, stdoutHandler)

	slog.SetDefault(slog.New(handler))
}

func newMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{handlers}
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(_ context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(context.Background(), level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(_ context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(context.Background(), record.Level) {
			// Make sure to clone the record before reusing it
			rec := record
			if err := h.Handle(context.Background(), rec); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{hs}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithGroup(name)
	}
	return &multiHandler{hs}
}
