package port

import (
	"context"
	"time"
)

// Clock defines a port for time-related operations
type Clock interface {
	Now() time.Time
}

// IDGenerator defines a port for generating unique IDs
type IDGenerator interface {
	NewID() string
}

// Logger defines a port for application logging
type Logger interface {
	Info(msg string, args ...any)  // ใช้สำหรับแสดงข้อมูลทั่วไป
	Error(msg string, args ...any) // Error ใช้แสดงข้อมูลข้อผิดพลาด
	Fatal(msg string, args ...any) // Fatal ใช้แสดงข้อมูลข้อผิดพลาดร้ายแรง และหยุดการทำงานของแอปพลิเคชัน
}

// StoragePort defines a port for object storage operations
type StoragePort interface {
	UploadImage(ctx context.Context, bucketName, objectName string, data []byte, contentType string) (string, error)
}
