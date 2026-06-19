# 🚀 Backend Go Self-Learning Checklist

## Phase 1: Deep-Dive Component Discovery

เป้าหมาย: ทำความเข้าใจเครื่องมือ, Library และชิ้นส่วนที่ใช้งานใน Core Service อย่างลึกซึ้ง ไม่ใช่แค่เรียกใช้งาน แต่ต้องรู้เบื้องหลัง

### A. HTTP Edge & Routing

- [ ] สร้าง `net/http` พื้นฐาน กับ Framework อย่าง Gin หรือ Fiber (Performance, Routing Engine)
- [ ] ศึกษาความแตกต่างระหว่าง `net/http` พื้นฐาน กับ Framework อย่าง Gin หรือ Fiber (Performance, Routing Engine)
- [ ] เรียนรู้วิธีการทำ Route Grouping และ Versioning (เช่น `/api/v1/...`)
- [ ] จำลองสถานการณ์และตั้งค่า CORS (Cross-Origin Resource Sharing) ให้ถูกต้อง (การจัดการ `Allow-Origins`, `Allow-Methods`, `Allow-Headers`)
- [ ] เรียนรู้วิธีการทำ Graceful Shutdown เพื่อให้ Server รอจนกว่า Request ที่ค้างอยู่จะทำงานเสร็จก่อนปิดตัว

### B. Logger (Structured & Centralized Logging)

- [ ] ศึกษาการ Setup Logger (เช่น `zap` หรือ `logrus`) ให้พ่น Log ออกมาเป็นรูปแบบ JSON
- [ ] ทำความเข้าใจ Log Levels (Debug, Info, Warn, Error, Fatal) และควรใช้ตอนไหน
- [ ] เรียนรู้วิธีการฝัง Logger เข้าไปใน `context.Context` (Context-aware logging) เพื่อให้สามารถดึงมาใช้ใน Layer อื่นๆ ได้
- [ ] เรียนรู้วิธีการใช้งาน `stack trace`
- [ ] ศึกษาการเพิ่มฟิลด์ (Index/Tags) เข้าไปใน Log เช่น `trace_id`, `user_id` เพื่อให้ง่ายต่อการค้นหาในระบบ Centralized Log (Kibana/Datadog)
- [ ] เรียนรู้กลไก Fire log (การยิง Log ออกไปยังปลายทางอื่นๆ นอกเหนือจาก Console เช่น ไฟล์ หรือ Log Aggregator)

### C. Middleware

- [ ] เรียนรู้วงจรชีวิตของ Middleware และการทำงานของคำสั่ง `Next()` (การส่งต่อ Request ไปยัง Layer ถัดไป)
- [ ] สร้าง Custom Middleware พื้นฐาน เช่น Request ID Generator, Authentication (JWT Validator)
- [ ] ศึกษาเรื่อง Execution Order ว่า Middleware ตัวไหนควรอยู่ก่อน-หลัง (เช่น Auth ควรมาก่อน Logic, Logger ควรครอบคลุมทั้งหมด)

### D. Tracer (Distributed Tracing)

- [ ] ทำความเข้าใจคอนเซปต์ของ OpenTelemetry (Trace, Span, Trace ID, Span ID)
- [ ] เรียนรู้วิธีการสร้าง Span ใหม่เมื่อมี Request เข้ามา (Root Span) และสร้าง Child Span เมื่อมีการเรียกใช้ DB หรือ External API
- [ ] ทดลองทำ Context Propagation ส่ง Trace ID ข้าม Service หรือนำไปผูกเข้ากับ Logger เพื่อให้ค้นหาได้ง่ายขึ้น

### E. Panic and Recovery

- [ ] ทำความเข้าใจว่าทำไม Go ถึงเกิด Panic และความแตกต่างระหว่าง Panic กับ Error
- [ ] ศึกษาการทำงานของ `defer recover()`
- [ ] ทดลองสร้าง Recovery Middleware เพื่อดักจับ Panic ที่เกิดจาก Business Logic หรือ Library ไม่ให้แอปพลิเคชันพัง (Crash)
- [ ] จัดการ Response เมื่อเกิด Panic ให้พ่น `500 Internal Server Error` กลับไปอย่างสวยงาม และ Log stack trace เก็บไว้

### F. Outbound REST Adapter

- [ ] สร้าง Custom HTTP Client แทนการใช้ `http.DefaultClient` เสมอ
- [ ] เรียนรู้การตั้งค่า Timeout ทั้งระดับ Client และระดับ Request (`context.WithTimeout`)
- [ ] สร้าง Wrapper สำหรับ Call External API ที่มีการจัดการ HTTP Status Code และอ่าน Body อย่างปลอดภัย
- [ ] เขียนระบบเพื่อ Log ทั้ง Request Payload และ Response Payload (พร้อมเทคนิค Masking ข้อมูล Sensitive เช่น Password/Token)

### G. Circuit Breaker

- [ ] ทำความเข้าใจ State Machine ของ Circuit Breaker (Closed, Open, Half-Open)
- [ ] ทดลอง Implement หรือใช้ Library (เช่น `sony/gobreaker`) ครอบ Outbound REST Adapter
- [ ] จำลองสถานการณ์ที่ External Service ล่ม เพื่อดูว่า Circuit Breaker ตัดวงจรอย่างไร
- [ ] ออกแบบ Fallback Strategy เมื่อ Circuit Breaker เปิด (เช่น คืนค่า Cache หรือ Default Value กลับไปแทนการพ่น Error)

### H. Metrics (`pkg/http-adapter/metric.go`)

- [ ] ทำความเข้าใจชนิดของ Prometheus Metrics หลักๆ (Counter, Histogram, Gauge)
- [ ] สร้างและเก็บค่า **HttpRequestTotal** (Counter) - นับจำนวน Request ทั้งหมดที่เข้ามา (แยกตาม Route/Status Code)
- [ ] สร้างและเก็บค่า **HttpLatency** (Histogram) - จับเวลาการประมวลผลของแต่ละ Request เพื่อนำไปหา P95, P99
- [ ] สร้างและเก็บค่า **HttpErrors** (Counter) - นับจำนวน Error ที่เกิดขึ้น
- [ ] นำ Metrics ไปฝังใน Middleware เพื่อให้ระบบเก็บค่าโดยอัตโนมัติ และสร้าง `/metrics` endpoint สำหรับให้ Prometheus มาดึงข้อมูล (Scrape)

### I. Database (PostgreSQL & sqlx)

- [ ] เรียนรู้วิธีการตั้งค่า Connection Pool ที่ดี (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)
- [ ] เข้าใจความแตกต่างระหว่าง `sql.DB` (Pool) และการทำ Transaction ด้วย `*sqlx.Tx`
- [ ] **Pool (`*sqlx.DB`) vs Client (`*sqlx.Conn`):**
- ศึกษาว่า **Pool** ใช้เมื่อต้องการให้ระบบจัดการ Connection อัตโนมัติ (ใช้งาน 99% ของเวลาทั้งหมด, ดึงจาก Pool มาใช้แล้วคืนกลับ)
- ศึกษาว่า **Client/Single Conn** ใช้เมื่อต้องการจอง Connection นั้นไว้ใช้งานแบบ Exclusive (เช่น การทำ Database Locks (`SELECT ... FOR UPDATE`), LISTEN/NOTIFY ใน Postgres, หรือการทำงานที่ต้องผูกติดกับ Session ID ของ Connection นั้นๆ)

---

## Phase 2: Architecture & Assembly (The Monolith Go Template)

เป้าหมาย: ประกอบทุกความรู้จาก Phase 1 เข้าด้วยกัน สร้างเป็น Monolith Repository ที่ได้มาตรฐาน Hexagonal/Clean Architecture และใช้ Dependency Injection (DI)

### 1. Project Setup & Skeleton

- [ ] จัดโครงสร้างโฟลเดอร์ตาม Standard Go Project Layout (เช่น `cmd/`, `internal/`, `pkg/`)
- [ ] เตรียมไฟล์ `go.mod` และกำหนดชื่อ Module ให้ถูกต้อง

### 2. Core Domain & Interfaces (The Hexagon Center)

- [ ] สร้างโฟลเดอร์ `internal/core/domain` หรือ `internal/entity` เพื่อเก็บ Struct ข้อมูลหลัก
- [ ] สร้างไฟล์กำหนด Ports (Interfaces) สำหรับ Repository (Driven Port)
- [ ] สร้างไฟล์กำหนด Ports (Interfaces) สำหรับ Usecase/Service (Driving Port)

### 3. Driven Adapters (Infrastructure)

- [ ] สร้าง Adapter สำหรับเชื่อมต่อ Database Postgres (`internal/adapter/repository/postgres`) พร้อม Implement Interface ที่กำหนดไว้
- [ ] สร้าง Adapter สำหรับ Outbound API ที่ฝัง Circuit Breaker เข้าไปเรียบร้อยแล้ว (`internal/adapter/outbound/...`)

### 4. Core Logic (Usecase/Service)

- [ ] สร้าง Business Logic (Usecase) โดยให้ `struct` รับ Dependency เข้ามาผ่าน Interface ของ Repository และ Outbound Adapter (Constructor Injection)
- [ ] เขียน Unit Test สำหรับ Usecase โดยใช้ Mock ของ Repository เพื่อพิสูจน์ว่า Logic ไม่ผูกติดกับ Database

### 5. Driving Adapters (Delivery/Transport)

- [ ] สร้าง HTTP Handler โดยใช้ Gin หรือ Fiber (`internal/adapter/handler/http`)
- [ ] นำ Middleware, Logger, Tracer และ Metrics จาก Phase 1 มาประกอบเข้ากับ Router

### 6. Dependency Injection (DI) & Wiring (`main.go`)

- [ ] สร้างไฟล์ Config เพื่อโหลดค่า Environment Variables (เช่น DB URL, Port)
- [ ] เริ่มกระบวนการ Wiring ใน `cmd/main.go` ตามลำดับ:
- Init Logger & Tracer
- Init Database Pool
- Inject DB เข้า Postgres Repository
- Inject Repository เข้า Usecase
- Inject Usecase เข้า HTTP Handler

- [ ] นำ HTTP Handler ไปผูกกับ Router (Gin/Fiber)

### 7. Final Polish

- [ ] ฝังระบบ Graceful Shutdown ลงใน `main.go`
- [ ] รันระบบ ยิง Request ทดสอบ ตั้งแต่ Router -> Middleware -> Usecase -> Database เพื่อตรวจสอบว่า Log, Metric และ Trace ไหลต่อกันสมบูรณ์ครบทุก Layer

---
