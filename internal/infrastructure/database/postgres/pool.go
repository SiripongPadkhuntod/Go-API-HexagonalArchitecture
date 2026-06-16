package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) { // NewPool() คือ function ที่ใช้สำหรับสร้าง instance ของ Pool
	config, err := pgxpool.ParseConfig(databaseURL) // ParseConfig() คือ function ที่ใช้สำหรับ parse config จาก databaseURL
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10                      // กำหนดค่า MaxConns คือ จำนวนการเชื่อมต่อสูงสุดที่จะสร้างขึ้น
	config.MinConns = 2                       // กำหนดค่า MinConns คือ จำนวนการเชื่อมต่อต่ำสุดที่จะสร้างขึ้น
	config.MaxConnLifetime = time.Hour        // กำหนดค่า MaxConnLifetime คือ ระยะเวลาสูงสุดที่การเชื่อมต่อจะมีอายุ
	config.MaxConnIdleTime = 30 * time.Minute // กำหนดค่า MaxConnIdleTime คือ ระยะเวลาสูงสุดที่การเชื่อมต่อจะอยู่ในสถานะ idle

	pool, err := pgxpool.NewWithConfig(ctx, config) // NewWithConfig() คือ function ที่ใช้สำหรับสร้าง instance ของ Pool เพื่อเชื่อมต่อฐานข้อมูล
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil { // Ping() คือ function ที่ใช้สำหรับ ping database เพื่อตรวจสอบว่าฐานข้อมูลพร้อมใช้งานหรือไม่
		pool.Close() // Close() คือ function ที่ใช้สำหรับปิด database
		return nil, err
	}

	return pool, nil
}
