package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewPool(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10) // ตั้งค่าจำนวน connection สูงสุดที่เปิดได้
	db.SetMaxIdleConns(2) // ตั้งค่าจำนวน connection ที่ไม่ได้ใช้งาน
	db.SetConnMaxLifetime(time.Hour) // ตั้งค่าอายุของ connection
	db.SetConnMaxIdleTime(30 * time.Minute) // ตั้งค่าอายุของ connection ที่ไม่ได้ใช้งาน

	if err := db.PingContext(ctx); err != nil { // PingContext() ใช้สำหรับตรวจสอบว่า connection ใช้งานได้หรือไม่
		db.Close() // Close() ใช้สำหรับปิด connection
		return nil, err // Return error ทันทีหาก connection ใช้งานไม่ได้
	}

	return db, nil // Return db, nil 
}
