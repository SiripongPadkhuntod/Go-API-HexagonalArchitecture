# Go Backend & Hexagonal Architecture

เนื้อหานี้ถูกปรับโครงสร้างใหม่ โดยจะเริ่มจากการอธิบายชิ้นส่วน (Components) และเครื่องมือต่างๆ ที่เป็นพื้นฐานให้ครบถ้วนก่อน จากนั้นจึงนำชิ้นส่วนทั้งหมดมาประกอบร่างกันเป็นสถาปัตยกรรม Hexagonal Architecture ตามแบบฉบับโปรเจคจริงครับ

---

# ส่วนที่ 1: เจาะลึกเครื่องมือและชิ้นส่วนของระบบ (Component Discovery)

## สไลด์ที่ 1: HTTP Edge & Routing (ประตูด่านหน้าของระบบ)

**หัวใจของการรับ Request คือการจัดการ HTTP Server ให้มีประสิทธิภาพ:**

- **Framework (Gin/Fiber):** ทำหน้าที่เป็น Routing Engine ช่วยจัดการ Path ให้ง่ายและเร็วกว่า `net/http` พื้นฐาน
- **Route Grouping & Versioning:** การจัดกลุ่ม API เช่น `/api/v1/users` เพื่อความเรียบร้อยและรองรับการทำเวอร์ชันในอนาคต
- **CORS (Cross-Origin Resource Sharing):** การตั้งค่าความปลอดภัยเบราว์เซอร์ เพื่ออนุญาตว่าใคร (Origin ไหน) สามารถเรียกใช้ API ของเราได้บ้าง และอนุญาตให้ส่ง Method/Header อะไรมาได้บ้าง
- **Graceful Shutdown:** กลไกสำคัญก่อน Server ปิดตัว โดยระบบจะไม่ปิดแอปพลิเคชันทันที แต่จะหยุดรับ Request ใหม่ และ "รอ" ให้ Request ที่ค้างอยู่ทำงานให้เสร็จก่อน ค่อยคืนทรัพยากรและปิดตัวอย่างปลอดภัย

## สไลด์ที่ 2: Middleware & Panic Recovery (ด่านตรวจและตัวกู้ภัย)

**Middleware คือโค้ดที่คั่นกลางระหว่าง Request และ Handler:**

- **Lifecycle & `Next()`:** Middleware สามารถดักการทำงานทั้ง "ก่อน" (เช็ค Auth) และ "หลัง" (คำนวณเวลาทำงาน) โดยใช้คำสั่ง `Next()` เพื่อส่งต่อ Request ไปยังชั้นถัดไป
- **Execution Order:** ลำดับสำคัญมาก เช่น ควรวาง Logger ไว้ตัวแรกสุดเพื่อเก็บ Log ทั้งหมด ตามด้วย Recovery, และ Auth (ถ้า Auth ไม่ผ่านก็จะไม่เข้าไปถึง Business Logic)
- **Panic vs Error:**
  - `Error`: ข้อผิดพลาดที่คาดการณ์ไว้ (เช่น หารหัสผ่านไม่เจอ) เราใช้ `if-else` เช็คและโปรแกรมทำงานต่อได้
  - `Panic`: ข้อผิดพลาดร้ายแรง/บั๊ก (เช่น ชี้ไปที่ Memory ปลอม) ทำให้โปรแกรมปิดตัวลงทันที (Crash)
- **Recovery:** Middleware สำคัญที่ใช้ `defer recover()` ทำหน้าที่ดักจับ Panic ไม่ให้ Server ตาย และแปลงให้เป็น HTTP 500 (Internal Server Error) ตอบกลับอย่างสวยงาม

## สไลด์ที่ 3: Observability (ระบบสังเกตการณ์: Log, Trace, Metric)

**เพื่อให้เรารู้ว่าระบบทำงานปกติ หรือพังที่ตรงไหน:**

1. **Logger (`zap`):** เก็บข้อมูลแบบ Structured (JSON) มีระดับความสำคัญ (Log Levels: Debug, Info, Warn, Error, Fatal) และควรทำ Context-aware (การฝัง Logger ลงใน `context` เพื่อให้พิมพ์ Log นำร่องพร้อมค่าต่างๆ เช่น `user_id` อัตโนมัติ) รวมถึงเทคนิค Fire Log ไปเก็บที่ส่วนกลาง (เช่น Kibana)
2. **Distributed Tracing (OpenTelemetry):** การตามรอย Request ที่กระโดดไปมา
   - **Trace ID:** เลขบัตรประชาชนของ 1 Request
   - **Span ID (Root & Child):** ระยะเวลาทำงานแต่ละจุด (เช่น Root Span คือเริ่มรับ Request -> Child Span คือจังหวะที่แวะไป Query DB)
   - **Context Propagation:** การส่ง Trace ID แปะไปกับ Context ข้ามไปยัง Service อื่น
3. **Metrics (Prometheus):** เก็บสถิติเชิงตัวเลข
   - `Counter`: ใช้นับจำนวน (เช่น ยอด Request รวมทั้งหมด `HttpRequestTotal`, จำนวน Error)
   - `Histogram`: ใช้จับเวลา (เช่น `HttpLatency` เพื่อหาค่า P95, P99 ว่าระบบโดยรวมตอบสนองเร็วแค่ไหน)
   - `Gauge`: ค่าที่ขึ้นและลงได้ (เช่น Memory ที่ใช้อยู่)

## สไลด์ที่ 4: การเชื่อมต่อภายนอกและความทนทาน (Outbound & Resilience)

**การไปคุยกับระบบอื่น (เช่น 3rd Party API) ต้องมีเกราะป้องกัน:**

- **Custom HTTP Client:** ไม่ควรใช้ Default Client ควรสร้างเองและกำหนด **Timeout** เพื่อไม่ให้ระบบเรารอจนค้างตาย
- **Wrapper & Masking:** การห่อหุ้มคำสั่งเรียก API ไว้ เพื่อดักจัดการ HTTP Status ต่างๆ อย่างรัดกุม พร้อมกับทำการ Masking (เซ็นเซอร์) ข้อมูลสำคัญเช่น Password ใน Log
- **Circuit Breaker:** ป้องกัน "ปฏิกิริยาลูกโซ่" (Cascading Failure) เมื่อระบบปลายทางพัง:
  - `Closed`: ปล่อย Request ผ่านปกติ
  - `Open`: ปลายทางล่มบ่อยถึงเกณฑ์ เบรกเกอร์ "ตัด" ทันที คืนค่า Error กลับเลยโดยไม่ต้องรอ Timeout
  - `Half-Open`: แอบส่ง Request ไปทดสอบดู ถ้าหายพังก็กลับไป Closed
  - **Fallback Strategy:** เมื่อวงจรตัด (Open) อาจคืนค่า Cache หรือข้อมูล Default ให้ผู้ใช้ไปก่อนแทนที่จะโชว์ Error

## สไลด์ที่ 5: ฐานข้อมูล (Database & Connection Pool)

**เทคนิคการทำงานกับ Database (`sqlx` หรือ Gorm):**

- **Connection Pool:** บ่อเก็บ Connection ของ DB เพื่อหยิบไปใช้และคืน โดยต้องตั้งค่าให้เหมาะสม เช่น `SetMaxOpenConns` (จำนวนเชื่อมต่อสูงสุด), `SetMaxIdleConns`, `SetConnMaxLifetime` (อายุการเชื่อมต่อ)
- **Pool (`DB`) vs Client (`Conn`):**
  - `Pool`: ใช้ใน 99% ของเวลาทั้งหมด เพื่อให้ระบบสลับ Connection กันใช้อัตโนมัติ
  - `Client/Single Conn`: จอง Connection ไว้ใช้คนเดียว มักใช้กับการ Lock ข้อมูลแบบเจาะจง (`SELECT ... FOR UPDATE`)

## สไลด์ที่ 6: Goroutine & Concurrency (การทำงานแบบคู่ขนาน)

**Goroutine คืออะไร?**
- คือการจำลอง Thread ขนาดเล็ก (Lightweight Thread) ที่ถูกจัดการโดย Go Runtime เอง ไม่ได้ผูกมัดกับ OS Thread แบบ 1:1
- กิน Memory น้อยมาก (เริ่มต้นแค่ ~2KB ต่อ 1 Goroutine) ทำให้แอปพลิเคชันสามารถสร้าง Goroutines หลักหมื่นหรือหลักแสนตัวทำงานพร้อมกันได้สบายๆ
- **ประโยชน์:** ช่วยให้ระบบรับ Request ได้เยอะขึ้น และเหมาะกับงานที่ต้อง "รอ" (I/O Bound) เช่น เรียก API ภายนอก หรือคุยกับ Database โดยไม่ไปบล็อกการทำงานหลัก

**การใช้งานและตัวอย่างโค้ด:**
เพียงแค่เติมคำว่า `go` นำหน้าการเรียกฟังก์ชัน ฟังก์ชันนั้นก็จะถูกแยกไปทำงานแบบคู่ขนาน (Background) ทันที

```go
func SendEmail(userID int) {
    // จำลองการทำงานที่ใช้เวลา 2 วินาที
    time.Sleep(2 * time.Second)
    fmt.Println("Email sent to", userID)
}

// 1. การเรียกใช้แบบปกติ (รอนาน 2 วินาที)
SendEmail(1)

// 2. การเรียกใช้แบบ Goroutine (ไม่รอ ข้ามไปทำบรรทัดต่อไปทันที)
go SendEmail(1)
```

**ข้อควรระวังเรื่อง Logging และ Context ใน Goroutine:**
เมื่อเราแตก Goroutine ออกมาทำงานแบบ Background จาก HTTP Request เราต้องระวัง 2 เรื่องคือ: 
1. `Context` ดั้งเดิมของ HTTP มักจะถูกยกเลิก (Cancel) ไปแล้วเมื่อตอบ Response จบ
2. หากไม่ส่ง Logger หรือ Context ที่ถูกต้องเข้าไป `Trace ID` ใน Log จะหายไป ทำให้ตามรอยไม่ได้

**ตัวอย่างการส่ง Context และ Logger เข้าไปใน Goroutine:**
```go
func ProcessOrderHandler(c *gin.Context, logger *zap.Logger) {
    // 1. ใช้ context.WithoutCancel เพื่อโคลน Context ออกมา (ป้องกัน Context โดนตัดเมื่อ Response จบ)
    bgCtx := context.WithoutCancel(c.Request.Context())
    
    // 2. สร้าง Goroutine พร้อมโยน Context และ Logger เข้าไปเป็น Parameter
    go func(ctx context.Context, log *zap.Logger) {
        log.Info("Starting background order processing...")
        // ทำงานจำลอง
        time.Sleep(2 * time.Second)
        log.Info("Background processing done!")
    }(bgCtx, logger)
    
    // รีบตอบกลับ Client ทันทีโดยไม่ต้องรองาน Background เสร็จ
    c.JSON(200, gin.H{"status": "processing_in_background"})
}
```

---

# ส่วนที่ 2: สถาปัตยกรรมและการประกอบร่าง (Architecture & Assembly)

## สไลด์ที่ 7: Software Architecture คืออะไร?

**สถาปัตยกรรมคือ "แปลนบ้าน" ของการเขียนโปรแกรม ช่วยให้ระบบโตได้ ไม่พังง่าย และแก้ไขสะดวก:**

- **MVC (Model-View-Controller):** แยกส่วน Database, การแสดงผล, และตรรกะควบคุม เหมาะกับโปรเจคไม่ใหญ่มาก แต่มักเกิดปัญหา Fat Controller เมื่อตรรกะเยอะขึ้น
- **Clean Architecture:** แบ่งเป็นชั้นวงกลมชัดเจน กฎคือ "วงนอกเรียกใช้วงในเท่านั้น" โครงสร้างแข็งแกร่งมากแต่อาจซับซ้อนไปสำหรับบางโปรเจค
- **Hexagonal Architecture (Ports and Adapters):** เน้นปฏิบัติ แบ่งเป็นแค่ **Core (ตรรกะธุรกิจหลัก)** และ **Adapters (ระบบภายนอก)** สื่อสารกันผ่าน **Ports (Interfaces)**

## สไลด์ที่ 8: โครงสร้างโปรเจค Hexagonal ของเรา (Project Layout)

- **`cmd/api/main.go`:** ตัว Entry Point ที่รันแอปพลิเคชัน
- **`internal/`:** โค้ดหลักที่ไม่อนุญาตให้โปรเจคอื่นดึงไปใช้
  - **`core/` (ศูนย์กลาง):**
    - `domain/`: โครงสร้างข้อมูลหลัก
    - `port/`: Interfaces ที่กำหนดวิธีคุยกับคนอื่น
    - `service/`: Use case หรือ Business Logic โดย **ห้าม** ยุ่งกับ Library หรือ DB โดยตรง
  - **`adapter/` (ตัวเชื่อมต่อ):**
    - `inbound/`: รับ Request (เช่น HTTP Handler)
    - `outbound/`: คุยกับภายนอก (เช่น Database Repository, หรือ API Client)
  - **`infrastructure/`:** จัดการเรื่องเทคนิค (โหลด Config, สร้าง DB Pool)

## สไลด์ที่ 9: Dependency Injection (การประกอบร่างที่ `main.go`)

**`main.go` เปรียบเสมือนโรงงานประกอบชิ้นส่วน โดยมีลำดับดังนี้:**

1. โหลด Config และตั้งค่า Observability (Logger, Tracer, Metrics)
2. สร้าง Database Connection Pool
3. **Inject (ฉีด)** Database เข้าไปใน Repository (Outbound Adapter)
4. **Inject** Repository เข้าไปใน Usecase/Service (Core)
5. **Inject** Usecase/Service เข้าไปใน HTTP Handler (Inbound Adapter)
6. นำ Handler ไปผูกกับ Router (Gin) พร้อมคลุมด้วย Middleware ต่างๆ
7. สั่งรัน Server พร้อมระบบ Graceful Shutdown

## สไลด์ที่ 10: การไหลของข้อมูล (Full Flow Request)

**เมื่อระบบประกอบเสร็จสมบูรณ์ 1 Request จะเดินทางดังนี้:**

1. Request วิ่งเข้ามาติด **Middleware** (บันทึก Log, จับเวลา, เช็ค Auth)
2. วิ่งเข้า **Router** ชี้เป้าไปที่ **HTTP Handler (Inbound Adapter)**
3. Handler แปลงข้อมูลและส่งให้ **Core Service** ผ่าน **Inbound Port**
4. Service คำนวณตรรกะธุรกิจ และสั่งงานผ่าน **Outbound Port**
5. **Repository (Outbound Adapter)** รับคำสั่งไป Query ข้อมูลจาก Database
6. (หากใช้ External API: ก็จะวิ่งผ่าน Circuit Breaker เพื่อป้องกันความเสี่ยง)
7. ผลลัพธ์ถูกส่งย้อนกลับมาจนถึง Handler, จัดรูปแบบ Response ออกไปให้ Client
8. Middleware ชั้นนอกสุด บันทึกเวลาตอบกลับและ Status Code จบวงจรอย่างสมบูรณ์
