package port

type Logger interface {
	Info(msg string, args ...any)  //ใช้สำหรับแสดงข้อมูลทั่วไป
	Error(msg string, args ...any) // Error ใช้แสดงข้อมูลข้อผิดพลาด
	Fatal(msg string, args ...any) // Fatal ใช้แสดงข้อมูลข้อผิดพลาดร้ายแรง และหยุดการทำงานของแอปพลิเคชัน
}
