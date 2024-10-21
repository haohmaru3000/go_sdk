package jaeger

import (
	"context"
	"flag"
	"fmt"

	"github.com/haohmaru3000/go_sdk/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type jaeger struct {
	logger            logger.Logger
	stopChan          chan bool
	processName       string
	sampleTraceRating float64
	agentURI          string
	port              int
	stdTracingEnabled bool
	context           context.Context
	tracerProvider    *sdktrace.TracerProvider
}

func NewJaeger(processName string) *jaeger {
	return &jaeger{
		processName: processName,
		stopChan:    make(chan bool),
	}
}

func (j *jaeger) Name() string {
	return "jaeger"
}

func (j *jaeger) GetPrefix() string {
	return j.Name()
}

func (j *jaeger) Get() interface{} {
	return j
}

func (j *jaeger) InitFlags() {
	flag.Float64Var(
		&j.sampleTraceRating,
		"jaeger-trace-sample-rate",
		1.0,
		"sample rating for remote tracing from OpenSensus: 0.0 -> 1.0 (default is 1.0)",
	)

	flag.StringVar(
		&j.agentURI,
		"jaeger-agent-uri",
		"",
		"jaeger agent URI to receive tracing data directly",
	)

	flag.IntVar(
		&j.port,
		"jaeger-agent-port",
		4318,
		"jaeger agent URI to receive tracing data directly",
	)

	flag.BoolVar(
		&j.stdTracingEnabled,
		"jaeger-std-enabled",
		false,
		"enable tracing export to std (default is false)",
	)
}

func (j *jaeger) Configure() error {
	j.logger = logger.GetCurrent().GetLogger(j.Name())
	return nil
}

func (j *jaeger) Run() error {
	if err := j.Configure(); err != nil {
		return err
	}

	return j.connectToJaegerAgent()
}

func (j *jaeger) Stop() <-chan bool {
	go func() {
		if !j.isEnabled() {
			j.stopChan <- true
			return
		}
		j.tracerProvider.Shutdown(j.context)
		j.stopChan <- true
		j.logger.Infoln("shut down Tracer-Provider")
	}()

	return j.stopChan
}

func (j *jaeger) isEnabled() bool {
	return j.agentURI != ""
}

func (j *jaeger) getSampler() sdktrace.Sampler {
	if j.sampleTraceRating >= 1 {
		return sdktrace.AlwaysSample()
	} else {
		return sdktrace.TraceIDRatioBased(j.sampleTraceRating)
	}
}

func (j *jaeger) getResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(j.processName),
		semconv.ServiceVersionKey.String("1.0.0"),
	)
}

func (j *jaeger) connectToJaegerAgent() error {
	ctx := context.Background()

	if !j.isEnabled() {
		return nil
	}

	url := fmt.Sprintf("%s:%d", j.agentURI, j.port)
	j.logger.Infof("connecting to Jaeger Agent on %s...", url)

	je, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(url),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(j.getSampler()),
		sdktrace.WithBatcher(je), // Set je as our 'Trace Exporter'
		sdktrace.WithResource(j.getResource()),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Trace view for console
	// if j.stdTracingEnabled {
	// 	// Register stats and trace exporters to export
	// 	// the collected data.
	// 	view.RegisterExporter(&PrintExporter{})

	// 	// Register the views to collect server request count.
	// 	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
	// 		j.logger.Errorf("jaeger error: %s", err.Error())
	// 	}
	// }

	j.context = ctx
	j.tracerProvider = tracerProvider

	return nil
}
