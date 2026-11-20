# OpenTelemetry-Logs-Go

[![Go Reference](https://pkg.go.dev/badge/github.com/hyperdxio/opentelemetry-logs-go.svg)](https://pkg.go.dev/github.com/hyperdxio/opentelemetry-logs-go)
[![codecov](https://codecov.io/github/hyperdxio/opentelemetry-logs-go/graph/badge.svg?token=F1NW0R0W75)](https://codecov.io/github/hyperdxio/opentelemetry-logs-go)

> **Deprecated:** `github.com/hyperdxio/opentelemetry-logs-go` is no longer
> maintained. The logs signal in the upstream
> [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/) is now
> stable; use it instead. It sends to ClickStack over OTLP: set
> `OTEL_EXPORTER_OTLP_ENDPOINT` to your ClickStack endpoint and
> `OTEL_EXPORTER_OTLP_HEADERS=authorization=<ingestion key>`.
>
> - Root package and `logs` (Logs Bridge API): [`go.opentelemetry.io/otel/log`](https://pkg.go.dev/go.opentelemetry.io/otel/log)
>   and [`go.opentelemetry.io/otel/log/global`](https://pkg.go.dev/go.opentelemetry.io/otel/log/global)
> - `sdk/logs`: [`go.opentelemetry.io/otel/sdk/log`](https://pkg.go.dev/go.opentelemetry.io/otel/sdk/log)
> - `exporters/otlp/otlplogs`: [`otlploghttp`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp)
>   and [`otlploggrpc`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc)
> - `exporters/stdout/stdoutlogs`: [`stdoutlog`](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/stdout/stdoutlog)
> - `autoconfigure`: [`go.opentelemetry.io/contrib/exporters/autoexport`](https://pkg.go.dev/go.opentelemetry.io/contrib/exporters/autoexport),
>   which picks the exporter from `OTEL_LOGS_EXPORTER`
>
> Upstream log bridges:
> [`otelslog`](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otelslog),
> [`otelzap`](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otelzap),
> [`otellogrus`](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otellogrus)
> and [`otellogr`](https://pkg.go.dev/go.opentelemetry.io/contrib/bridges/otellogr).

OpenTelemetry-Logs-Go is the [Go](https://golang.org) implementation of [OpenTelemetry](https://opentelemetry.io/) Logs.
It provides API to directly send logging data to observability platforms. It is an extension of official
[open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) to support Logs.

## Project Life Cycle

This project was created due log module freeze in
official [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) repository:

```
The Logs signal development is halted for this project while we stablize the Metrics SDK. 
No Logs Pull Requests are currently being accepted.
```

This project will be deprecated once official [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go)
repository Logs module will have status "Stable".

## Project packages

| Packages                         | Description                                                                |
|----------------------------------|----------------------------------------------------------------------------|
| [autoconfigure](./autoconfigure) | Autoconfiguration SDK. Allow to configure log exporters with env variables |
| [sdk](./sdk)                     | Opentelemetry Logs SDK                                                     |
| [exporters/otlp](./exporters)    | OTLP format exporter                                                       |
| [exporters/stdout](./exporters)  | Console exporter                                                           |                                                            

## Quick start

This is an implementation of [Logs Bridge API](https://opentelemetry.io/docs/specs/otel/logs/bridge-api/) and not
intended to use by developers directly. It is provided for logging library authors to build log appenders, which use
this API to bridge between existing logging libraries and the OpenTelemetry log data model.

Example bellow will show how logging library could be instrumented with current API:

```go
package myInstrumentedLogger

import (
	otel "github.com/hyperdxio/opentelemetry-logs-go"
	"github.com/hyperdxio/opentelemetry-logs-go/logs"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

const (
	instrumentationName    = "otel/zap"
	instrumentationVersion = "0.0.1"
)

var (
	logger = otel.GetLoggerProvider().Logger(
		instrumentationName,
		logs.WithInstrumentationVersion(instrumentationVersion),
		logs.WithSchemaURL(semconv.SchemaURL),
	)
)

func (c otlpCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {

	lrc := logs.LogRecordConfig{
		Body: &ent.Message,
		...
	}
	logRecord := logs.NewLogRecord(lrc)
	logger.Emit(logRecord)
}
```

and application initialization code:

```go
package main

import (
	"os"
	"context"
	"github.com/hyperdxio/opentelemetry-logs-go"
	"github.com/hyperdxio/opentelemetry-logs-go/exporters/otlp/otlplogs"
	"github.com/hyperdxio/opentelemetry-logs-go/exporters/otlp/otlplogs/otlplogshttp"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	sdk "github.com/hyperdxio/opentelemetry-logs-go/sdk/logs"
)

func newResource() *resource.Resource {
	host, _ := os.Hostname()
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("otlplogs-example"),
		semconv.ServiceVersion("0.0.1"),
		semconv.HostName(host),
	)
}

func main() {
	ctx := context.Background()

	exporter, _ := otlplogs.NewExporter(ctx, otlplogs.WithClient(otlplogshttp.NewClient()))
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(exporter),
		sdk.WithResource(newResource()),
	)
	otel.SetLoggerProvider(loggerProvider)

	myInstrumentedLogger.Info("Hello OpenTelemetry")
}
```

## References

Logger Bridge API implementations for `zap`, `slog`, `zerolog` and other
loggers can be found in https://github.com/hyperdxio/opentelemetry-go

