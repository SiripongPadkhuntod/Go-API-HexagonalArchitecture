package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string // contextKey คือ type ที่ใช้สำหรับกำหนด key ใน context

const loggerContextKey contextKey = "logger" // loggerContextKey คือ instance ของ contextKey ที่ใช้สำหรับกำหนด key ใน context

var fallbackLogger = zap.NewNop() // fallbackLogger คือ instance ของ zap.Logger ที่ใช้สำหรับกำหนด key ใน context

func New() (*zap.Logger, error) { // New() คือ function ที่ใช้สำหรับสร้าง instance ของ zap.Logger
	config := zap.NewProductionConfig()                          // NewProductionConfig() คือ function ที่ใช้สำหรับสร้าง instance ของ zap.Logger
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // EncodeTime คือ function ที่ใช้สำหรับกำหนด key ใน context คือ ISO8601TimeEncoder
	return config.Build()                                        // Build() คือ function ที่ใช้สำหรับสร้าง instance ของ zap.Logger
}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok {
		return logger
	}

	return fallbackLogger
}
