package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// мидлвар для трейсов
func Tracing(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)
	return func(c *gin.Context) {
		// извлекаем контекст из запроса
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// серверный спан
		spanName := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())

		// устанавливаем его в контекст
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				semconv.HTTPRouteKey.String(c.FullPath()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// контекст передаем в джин-конекст
		c.Request = c.Request.WithContext(ctx)

		// идем дальше
		c.Next()
	}
}
