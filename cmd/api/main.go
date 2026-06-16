package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	_ "hexagonalarchitecture/docs"
	httpadapter "hexagonalarchitecture/internal/adapter/inbound/http"
	clockadapter "hexagonalarchitecture/internal/adapter/outbound/clock/system"
	"hexagonalarchitecture/internal/adapter/outbound/event/httpclient"
	"hexagonalarchitecture/internal/adapter/outbound/event/noop"
	idadapter "hexagonalarchitecture/internal/adapter/outbound/id/uuid"
	"hexagonalarchitecture/internal/adapter/outbound/repository/postgres"
	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/core/service"
	"hexagonalarchitecture/internal/infrastructure/config"
	databasepostgres "hexagonalarchitecture/internal/infrastructure/database/postgres"
	observabilitylogger "hexagonalarchitecture/internal/infrastructure/observability/logger"
	"hexagonalarchitecture/internal/infrastructure/observability/tracer"
)

// @title Hexagonal Architecture CRUD API
// @version 1.0
// @description A CRUD REST API built with Gin, PostgreSQL, Docker, and Hexagonal Architecture.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	logger, err := observabilitylogger.New() // สร้าง instance ของ logger (instance หมายถึง การสร้าง object จาก class)
	if err != nil {
		panic(err) // panic() ใช้สำหรับหยุดการทำงานของโปรแกรมทันทีเมื่
	}
	zap.ReplaceGlobals(logger)                             // แทนที่ logger ทั่วโลกด้วย logger ที่สร้างขึ้น
	defer logger.Sync()                                    // ปิด logger เมื่อโปรแกรมทำงานเสร็จสิ้น
	appLogger := observabilitylogger.NewZapAdapter(logger) // สร้าง instance ของ logger adapter

	tracerProvider, err := tracer.NewProvider() // สร้าง instance ของ tracer provider
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // สร้าง context เพื่อรอการปิด tracer
		defer cancel()                                                          // ปิด context เมื่อโปรแกรมทำงานเสร็จสิ้น
		if err := tracer.Shutdown(ctx, tracerProvider); err != nil {            // ปิด tracer
			logger.Error("failed to shutdown tracer", zap.Error(err))
		}
	}()

	metricsRegistry := prometheus.NewRegistry()  // สร้าง instance ของ metrics registry
	httpadapter.RegisterMetrics(metricsRegistry) // ลงทะเบียน metrics

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // สร้าง context เพื่อรอการปิด
	defer cancel()                                                           // ปิด context เมื่อโปรแกรมทำงานเสร็จสิ้น

	dbPool, err := databasepostgres.NewPool(ctx, cfg.Database.URL()) // สร้าง instance ของ database pool
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err)) // แสดงข้อผิดพลาดและหยุดการทำงานของแอปพลิเคชัน
	}
	defer dbPool.Close() // ปิด database pool เมื่อโปรแกรมทำงานเสร็จสิ้น

	userRepo := postgres.NewUserRepository(dbPool) // สร้าง instance ของ user repository

	outboundClient := newOutboundAPIClient(cfg, appLogger)         // สร้าง instance ของ outbound client
	idGenerator := idadapter.NewGenerator()                        // สร้าง instance ของ id generator
	clock := clockadapter.NewClock()                               // สร้าง instance ของ clock
	userService := service.NewUserService(service.UserServiceDeps{ // สร้าง instance ของ user service
		Repo:      userRepo,       // ส่ง user repository ไปยัง user service
		Publisher: outboundClient, // ส่ง outbound client ไปยัง user service
		Logger:    appLogger,      // ส่ง logger ไปยัง user service
		IDs:       idGenerator,    // ส่ง id generator ไปยัง user service
		Clock:     clock,          // ส่ง clock ไปยัง user service
	})

	r := httpadapter.New(userService, appLogger, otel.Tracer("hexagonalarchitecture-api"), metricsRegistry) // สร้าง instance ของ http adapter
	server := &http.Server{                                                                                 // สร้าง instance ของ http server
		Addr:    cfg.ServerAddress(), // รับค่า addr จาก config
		Handler: r,                   // รับค่า handler จาก http adapter
	}

	logger.Info("server is running", zap.String("address", cfg.ServerAddress())) // แสดงข้อความว่า server กำลังทำงาน
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { // ตรวจสอบว่ามีข้อผิดพลาดหรือไม่ และไม่ใช่ข้อผิดพลาดจากการปิด server
			logger.Fatal("failed to run server", zap.Error(err)) // แสดงข้อผิดพลาดและหยุดการทำงานของแอปพลิเคชัน
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM) // สร้าง context เพื่อรอการปิด server
	defer stop()                                                                                     // ปิด context เมื่อโปรแกรมทำงานเสร็จสิ้น

	<-shutdownCtx.Done() // รอการปิด server
	stop()               // ปิดการรับสัญญาณ

	logger.Info("shutdown signal received") // แสดงข้อความว่าได้รับสัญญาณปิด

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second) // สร้าง context เพื่อรอการปิด server
	defer cancel()                                                          // ปิด context เมื่อโปรแกรมทำงานเสร็จสิ้น

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("failed to shutdown server gracefully", zap.Error(err)) // แสดงข้อผิดพลาดและหยุดการทำงานของแอปพลิเคชัน
	}

	logger.Info("server stopped gracefully") // แสดงข้อความว่า server ปิดตัวลงแล้ว
}

func newOutboundAPIClient(cfg config.Config, logger port.Logger) port.UserEventPublisher { // สร้าง instance ของ outbound client
	if cfg.OutboundAPI.BaseURL == "" { // ตรวจสอบว่ามี BaseURL หรือไม่
		return noop.NewClient() // ถ้าไม่มี ให้คืนค่า noop.NewClient() เพื่อไม่ต้องติดต่อกับ API อื่น
	}

	client, err := httpclient.New(httpclient.Config{ // สร้าง instance ของ http client
		BaseURL: cfg.OutboundAPI.BaseURL, // กำหนดค่า BaseURL จาก config
	})
	if err != nil { // ตรวจสอบว่ามีข้อผิดพลาดหรือไม่
		logger.Fatal("failed to create outbound API client", "error", err) // แสดงข้อผิดพลาดและหยุดการทำงานของแอปพลิเคชัน
	}

	return client // คืนค่า outbound client
}
