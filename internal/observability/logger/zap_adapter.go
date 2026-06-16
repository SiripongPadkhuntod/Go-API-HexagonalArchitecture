package logger

import (
	"go.uber.org/zap"

	"hexagonalarchitecture/internal/core/port"
)

var _ port.Logger = (*ZapAdapter)(nil)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{logger: logger}
}

func (l *ZapAdapter) Info(msg string, args ...any) {
	l.logger.Info(msg, fields(args...)...)
}

func (l *ZapAdapter) Error(msg string, args ...any) {
	l.logger.Error(msg, fields(args...)...)
}

func (l *ZapAdapter) Fatal(msg string, args ...any) {
	l.logger.Fatal(msg, fields(args...)...)
}

func fields(args ...any) []zap.Field {
	zapFields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			zapFields = append(zapFields, zap.Any("arg", args[i]))
			continue
		}
		if i+1 >= len(args) {
			zapFields = append(zapFields, zap.Any(key, nil))
			continue
		}
		zapFields = append(zapFields, zap.Any(key, args[i+1]))
	}

	return zapFields
}
