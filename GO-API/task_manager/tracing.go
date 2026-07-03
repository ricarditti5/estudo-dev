package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint() /*<-- formata o JSON de forma legível no terminal (sem isto, sai tudo numa linha só, ilegível).*/) //cria o exporter.
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		/*cria o motor, configurado com:
		- WithBatcher(exporter) — diz ao provider para agrupar spans e enviá-los ao exporter em lotes (mais eficiente do que um a um).
		- WithResource(...) — metadados sobre o teu serviço. semconv.ServiceName("task-manager") faz com que todos os traces gerados fiquem identificados como vindos desta aplicação — importante quando tiveres vários serviços a exportar para o mesmo sítio, no futuro.*/
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("task-manager"),
		)),
	)

	otel.SetTracerProvider(tp)
	/*regista este provider como o global da aplicação. Isto permite que, em qualquer sítio do código, chames otel.Tracer("nome") sem teres de passar o tp manualmente por todo o lado (parecido com a ideia do logger, mas o OpenTelemetry tem este mecanismo global embutido).*/

	return tp, nil
}
