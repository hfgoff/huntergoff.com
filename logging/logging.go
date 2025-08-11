// Package logging was ripped from https://djwong.net/2025/05/28/cool-go-slog-tricks.html
package logging

import (
	"context"
	"log/slog"
	"sync"
)

var (
	globalHandler slog.Handler
	pendingLogs   []slog.Record
	pendingLogsMu sync.Mutex
)

func Provide(handler slog.Handler) {
	pendingLogsMu.Lock()
	for _, record := range pendingLogs {
		_ = handler.Handle(context.Background(), record)
	}
	pendingLogs = nil
	pendingLogsMu.Unlock()
	globalHandler = handler
}

func New(name string) *slog.Logger {
	handler := &placeholder{
		attrs: []slog.Attr{
			slog.String("instrument", name),
		},
	}
	return slog.New(handler)
}

type placeholder struct {
	attrs  []slog.Attr
	groups []string

	once    sync.Once
	handler slog.Handler
}

func (h *placeholder) init() {
	h.once.Do(func() {
		handler := globalHandler
		for _, group := range h.groups {
			handler = handler.WithGroup(group)
		}
		h.handler = handler.WithAttrs(h.attrs)
	})
}

func (h *placeholder) Enabled(ctx context.Context, level slog.Level) bool {
	if globalHandler == nil {
		return true
	}

	return globalHandler.Enabled(ctx, level)
}

func (h *placeholder) Handle(ctx context.Context, record slog.Record) error {
	if globalHandler == nil {
		pendingLogsMu.Lock()
		pendingLogs = append(pendingLogs, record)
		pendingLogsMu.Unlock()
		return nil
	}

	h.init()
	return h.handler.Handle(ctx, record)
}

func (h *placeholder) WithAttrs(attrs []slog.Attr) slog.Handler {
	if globalHandler != nil {
		h.init()
		return h.handler.WithAttrs(attrs)
	}

	return &placeholder{
		attrs:  append(append([]slog.Attr{}, attrs...), h.attrs...),
		groups: append([]string{}, h.groups...),
	}
}

func (h *placeholder) WithGroup(name string) slog.Handler {
	if globalHandler != nil {
		h.init()
		return h.handler.WithGroup(name)
	}

	return &placeholder{
		attrs:  append([]slog.Attr{}, h.attrs...),
		groups: append(append([]string{}, h.groups...), name),
	}
}

// Error returns a slog.Attr for logging error messages.
func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
