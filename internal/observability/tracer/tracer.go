package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

//sdktrace.TracerProvider คือ interface ที่ใช้สำหรับสร้าง instances ของ trace providers  (ตัวจัดการข้อมูล trace)
//stdouttrace.Exporter คือ interface ที่ใช้สำหรับส่งออก (export) ข้อมูล trace ไปยังปลายทาง เช่น console, file, หรือ collector
//stdouttrace.WithPrettyPrint() คือ function ที่ใช้สำหรับกำหนดรูปแบบการแสดงผลของ trace ให้เป็นแบบ pretty-print (จัดรูปแบบให้อ่านง่าย)
//sdktrace.New() คือ function ที่ใช้สำหรับสร้าง instance ของ stdouttrace.Exporter
//sdktrace.WithBatcher(exporter) คือ function ที่ใช้สำหรับกำหนด key ใน context คือ stdouttrace.Exporter
//sdktrace.WithSampler(sdktrace.AlwaysSample()) คือ function ที่ใช้สำหรับกำหนด key ใน context คือ stdouttrace.Exporter
//otel.SetTracerProvider() คือ function ที่ใช้สำหรับกำหนด key ใน context คือ stdouttrace.Exporter
//otel.SetTextMapPropagator() คือ function ที่ใช้สำหรับกำหนด key ใน context คือ stdouttrace.Exporter

func NewProvider() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	provider := sdktrace.NewTracerProvider(  
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return provider, nil
}

func Shutdown(ctx context.Context, provider *sdktrace.TracerProvider) error {
	if provider == nil {
		return nil
	}

	return provider.Shutdown(ctx)
}
