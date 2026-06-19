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

	//custom logger
	// config := zap.Config{
	// 	Level:             zap.NewAtomicLevelAt(zap.InfoLevel), // InfoLevel คือ ระดับของ log
	// 	Development:       false, // false คือ production mode
	// 	DisableCaller:     false, // false คือ ไม่ต้อง disable caller
	// 	DisableStacktrace: false, // false คือ ไม่ต้อง disable stacktrace
	// 	Encoding:          "json", // json คือ รูปแบบของ log
	// 	EncoderConfig: zapcore.EncoderConfig{
	// 		MessageKey:     "message", // กำหนดชื่อ key สำหรับข้อความ
	// 		LevelKey:       "level", // กำหนดชื่อ key สำหรับระดับของ log
	// 		TimeKey:        "timestamp", // กำหนดชื่อ key สำหรับเวลา
	// 		NameKey:        "logger", // กำหนดชื่อ key สำหรับชื่อของ logger
	// 		CallerKey:      "caller", // กำหนดชื่อ key สำหรับ caller
	// 		FunctionKey:    "function", // กำหนดชื่อ key สำหรับ function
	// 		StacktraceKey:  "stacktrace", // กำหนดชื่อ key สำหรับ stacktrace
	// 		LineEnding:     zapcore.DefaultLineEnding, // กำหนด line ending
	// 		EncodeLevel:    zapcore.CapitalLevelEncoder, // กำหนดการ encode level
	// 		EncodeCaller:   zapcore.ShortCallerEncoder, // กำหนดการ encode caller
	// 		EncodeTime:     zapcore.ISO8601TimeEncoder, // กำหนดการ encode time
	// 		EncodeDuration: zapcore.SecondsDurationEncoder, // กำหนดการ encode duration
	// 	},
	// 	OutputPaths: []string{"stdout"}, // กำหนด output paths
	// 	ErrorOutputPaths: []string{"stderr"}, // กำหนด error output paths
	// }
	
	config := zap.NewProductionConfig()                          // NewProductionConfig() คือ function ที่ใช้สำหรับสร้าง instance ของ zap.Logger
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // EncodeTime คือ function ที่ใช้สำหรับกำหนด key ใน context คือ ISO8601TimeEncoder
	return config.Build()                                        // Build() คือ function ที่ใช้สำหรับสร้าง instance ของ zap.Logger
}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok { // ok คือตัวแปรที่ใช้สำหรับตรวจสอบว่า loggerContextKey มีค่าอยู่ใน ctx หรือไม่
		return logger
	}

	return fallbackLogger
}
