package uuid

import (
	"github.com/google/uuid"

	"hexagonalarchitecture/internal/core/port"
)

var _ port.IDGenerator = (*Generator)(nil) // IDGenerator คือ interface ที่กำหนดไว้ใน core/port และ Generator คือ struct ที่ implements interface นี้

type Generator struct{}

func NewGenerator() *Generator { // NewGenerator() คือ function ที่ใช้สำหรับสร้าง instance ของ Generator
	return &Generator{} // NewGenerator() คืนค่า instance ของ Generator
}

func (g *Generator) NewID() string { // NewID() คือ function ที่ใช้สำหรับสร้าง instance ของ Generator
	return uuid.New().String() // NewID() คืนค่า instance ของ Generator ที่เป็น string โดยใช้ uuid.New() สร้าง UUID และ .String() แปลงเป็น string
}

// NewID with Prefix
func (g *Generator) NewIDWithPrefix(prefix string) string { // NewIDWithPrefix() คือ function ที่ใช้สำหรับสร้าง instance ของ Generator ที่มี prefix
	return prefix + uuid.New().String() // NewIDWithPrefix() คืนค่า instance ของ Generator ที่เป็น string โดยใช้ uuid.New() สร้าง UUID และ .String() แปลงเป็น string
}