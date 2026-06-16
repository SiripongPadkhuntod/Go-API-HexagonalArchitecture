package http

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/core/usecase"
)

const requestIDHeader = "X-Request-ID"

func RequestContextMiddleware(logger port.Logger, tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Writer.Header().Set(requestIDHeader, requestID)

		ctx := otel.GetTextMapPropagator().Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		ctx, span := tracer.Start(ctx, "http "+c.Request.Method+" "+route)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", route),
			attribute.String("http.target", c.Request.URL.Path),
			attribute.String("request_id", requestID),
		)

		logArgs := requestLogArgs(ctx, requestID, c.Request.Method, c.Request.URL.Path, route)
		c.Request = c.Request.WithContext(ctx)
		logger.Info("request started", logArgs...)

		c.Next()

		statusCode := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))
		if statusCode >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		}

		logger.Info("request completed", append(logArgs,
			"status_code", statusCode,
			"latency", time.Since(start),
		)...)
	}
}

func RecoveryMiddleware(logger port.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered", append(requestLogArgs(
					c.Request.Context(),
					c.GetHeader(requestIDHeader),
					c.Request.Method,
					c.Request.URL.Path,
					c.FullPath(),
				),
					"panic", recovered,
					"stack", string(debug.Stack()),
				)...)

				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					newErrorResponse(usecase.ERROR_CODE_INTERNAL_SERVER_ERROR, usecase.ERROR_MESSAGE_INTERNAL_SERVER_ERROR),
				)
			}
		}()

		c.Next()
	}
}

func requestLogArgs(ctx context.Context, requestID, method, path, route string) []any {
	args := []any{
		"request_id", requestID,
		"method", method,
		"path", path,
		"route", route,
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return args
	}

	return append(args,
		"trace_id", spanContext.TraceID().String(),
		"span_id", spanContext.SpanID().String(),
	)
}
