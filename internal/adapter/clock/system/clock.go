package system

import (
	"time"

	"hexagonalarchitecture/internal/core/port"
)

var _ port.Clock = (*Clock)(nil) //_ port.Clock คือ interface ของ clock ที่ถูกกำหนดไว้ใน package port
//ส่วน (*Clock)(nil) คือ pointer ของ Clock ที่ถูกกำหนดไว้ใน package clock

type Clock struct{}

func NewClock() *Clock { // NewClock() คือ function ที่ใช้สำหรับสร้าง instance ของ Clock
	return &Clock{} // return &Clock() คือการคืนค่า instance ของ Clock
}

func (c *Clock) Now() time.Time { // Now() คือ function ที่ใช้สำหรับรับค่าเวลาปัจจุบัน
	return time.Now().UTC() // return time.Now().UTC() คือการคืนค่าเวลาปัจจุบันในรูปแบบ UTC
}
