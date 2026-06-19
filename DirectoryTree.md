# 📂 โครงสร้างโปรเจค (Project Structure)

โปรเจคนี้พัฒนาด้วยแนวคิด **Hexagonal Architecture (Ports and Adapters)** โดยมีการแบ่งชั้นของโค้ดอย่างชัดเจน เพื่อให้ระบบมีความยืดหยุ่น ดูแลรักษาง่าย และแยก Business Logic ออกจาก Framework และ Database ภายนอกอย่างเด็ดขาด

## 🌳 Directory Tree (ภาพรวม)

```text
.
├── cmd/
│   └── api/                  # จุดเริ่มต้นของแอปพลิเคชัน (Entry point) 
│       └── main.go           # ไฟล์รันเซิร์ฟเวอร์ ประกอบส่วนต่างๆ (Dependency Injection)
├── db/                       # สคริปต์และไฟล์ตั้งต้นของฐานข้อมูล
│   ├── mysql/init/           # สคริปต์ SQL สำหรับ MySQL
│   └── postgres/init/        # สคริปต์ SQL สำหรับ PostgreSQL
├── internal/                 # โค้ดหลักของโปรเจกต์ (จำกัดสิทธิ์นำไป import โดยตรง)
│   ├── adapter/              # Adapters: ตัวแปลงข้อมูลระหว่าง Core กับโลกภายนอก
│   │   ├── inbound/          # 📥 ฝั่งรับ Request (Driving Adapters)
│   │   │   ├── grpc/         # ตัวจัดการ gRPC
│   │   │   └── http/         # ตัวจัดการ HTTP REST API (Handlers, Router, Middleware)
│   │   └── outbound/         # 📤 ฝั่งคุยกับระบบอื่น (Driven Adapters)
│   │       ├── clock/        # ตัวจัดการเรื่องเวลา
│   │       ├── event/        # ตัวจัดการยิง Event ภายนอก
│   │       ├── id/           # ตัวจัดการสร้าง ID (เช่น UUID)
│   │       ├── repository/   # ตัวจัดการ Database (MySQL, Postgres)
│   │       └── storage/      # การเชื่อมต่อที่เก็บไฟล์ (เช่น MinIO)
│   ├── core/                 # 🧠 Core / Business Logic (หัวใจหลัก ไม่มี Library ภายนอก)
│   │   ├── domain/           # Entities และ Business Rules หลัก
│   │   ├── port/             # Interfaces กำหนดข้อตกลง (Contracts) กับภายนอก
│   │   └── service/          # Use Cases เรียกใช้ Domain ผ่านทาง Port
│   └── infrastructure/       # โครงสร้างพื้นฐานทางเทคนิค
│       ├── config/           # โหลด Configuration (เช่น จากไฟล์ .env)
│       ├── database/         # จัดการ Connection Pool
│       └── observability/    # จัดการการสังเกตการณ์ (Logging ด้วย Zap, Tracing)
├── pkg/                      # 📦 โค้ดส่วนกลางที่ให้โปรเจกต์อื่นนำไปใช้ได้
├── docker-compose.yml        # ไฟล์ตั้งค่า Docker สำหรับ Development
├── go.mod / go.sum           # ไฟล์จัดการ Dependencies ของ Go
└── README.md                 # เอกสารอธิบายโปรเจกต์
```

## 📖 รายละเอียดแต่ละส่วน (Details)

### 🚀 `cmd/api/main.go`
นี่คือจุด **Entry Point** หรือศูนย์กลางควบคุมการทำงาน ตอนรันแอปจะเริ่มจากที่นี่ หน้าที่หลักคือโหลดคอนฟิก, เชื่อมต่อฐานข้อมูล, และทำ **Dependency Injection** (การหยิบ Adapter และ Service ต่างๆ มาเสียบเข้าหากัน) ก่อนจะเปิดพอร์ตรับ Request

### 🔒 `internal/`
โฟลเดอร์สำหรับโค้ดหลักของโปรเจคที่ไม่ต้องการให้โปรเจคภายนอกเรียกใช้งาน (`import`) โดยตรง
- **`core/` (หัวใจของแอป):** เก็บ Business Logic และ Entity ต่างๆ โดยชั้นนี้จะ **ห้ามมีโค้ดที่ไปผูกติดกับ Database หรือ Web Framework เด็ดขาด** 
  - `domain/`: โครงสร้างข้อมูล (Structs) และเงื่อนไขเฉพาะของตัวข้อมูล
  - `port/`: สร้าง Interface ขึ้นมาเพื่อบอกว่า Core ต้องการคุยกับภายนอกด้วยคำสั่งอะไรบ้าง
  - `service/`: หรือ Usecase เก็บเงื่อนไขทางธุรกิจ (เช่น สมัครสมาชิก, เช็คยอดเงิน)
- **`adapter/` (ตัวแปลงภาษา):** ส่วนที่เชื่อมโลกภายนอกเข้ากับ Core
  - `inbound/`: ฝั่ง **"รับ"** ข้อมูลเข้ามา เช่น HTTP, gRPC (Gin / Fiber จะอยู่ที่นี่)
  - `outbound/`: ฝั่ง **"ส่ง"** ข้อมูลออกไป เช่น เขียนลง Database, ไปเซฟไฟล์ที่ MinIO, สร้าง UUID
- **`infrastructure/` (โครงสร้างพื้นฐาน):** จัดการเรื่องเชิงเทคนิคลึกๆ เช่น การสร้าง Database Connection Pool, การทำ Logs (Zap) และ Tracing (OpenTelemetry)

### 📦 `pkg/`
ถ้ามีโค้ดบางอย่างที่เป็น Utility กลาง (เช่น ตัวช่วยจัดการ String, วันที่) และสามารถให้โปรเจคอื่นในบริษัทนำไป `import` ใช้ได้ด้วย จะนำมาใส่ไว้ในโฟลเดอร์นี้

### 🗄️ `db/`
เก็บสคริปต์ Database Migration หรือไฟล์ `.sql` สำหรับตอนรันขึ้นมาครั้งแรก 

### 🐳 `docker-compose.yml`
ใช้สำหรับการรัน Environment ที่จำเป็นทั้งหมด (เช่น ฐานข้อมูลจำลอง, Redis) บนเครื่องนักพัฒนา (Local Development) แบบครบจบในคำสั่งเดียว
