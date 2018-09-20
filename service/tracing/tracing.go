package tracing

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var tracerIsSet = false

// New sets opentracing.GlobalTracer() to tracer created from function options
// returns Closer, which is used to close the tracker
func New(serviceName, hostPort string) io.Closer {
	log.Printf("Creating new tracer %s on host %s", serviceName, hostPort)

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hostPort,
		},
	}

	tracer, closer, err := cfg.New(
		serviceName,
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		// TODO: Could probably set tracer to send data even if there is no connection to host
		log.Printf("Error initializing tracker: %v", err)
		return MockCloser{}
	}

	opentracing.SetGlobalTracer(tracer)
	tracerIsSet = true

	log.Printf("New tracer %s created on host %s", serviceName, hostPort)

	return closer
}

func Middleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tracerIsSet {
			var sp opentracing.Span
			spanName := fmt.Sprintf("Handler %s %s", r.Method, r.URL.Path)
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))

			if err != nil {
				// If for whatever reason we can't join, go ahead an start a new root span.
				sp = opentracing.StartSpan(spanName)
				log.Printf("TRACE NOT FOUND, %v", err)
			} else {
				sp = opentracing.StartSpan(spanName, opentracing.ChildOf(wireContext))
				log.Printf("TRACE FOUND")
			}
			defer sp.Finish()

			sp.Tracer().Inject(
				sp.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))

			r = r.WithContext(context.WithValue(r.Context(), "header", r.Header))
		} else {
			log.Printf("Tracer is not set")
		}
		h.ServeHTTP(w, r)
	}
}

// TraceFunctionSpan creates new span and then executes provided function.
// If opentracing.GlobalTracer() is not set, then no span is reported nor created
// If possible tracer is extracted from request.
// If you have no request handy, pass in nil
func TraceFunctionSpan(name string, header http.Header, f func() error) error {

	// Create new span if tracer is set
	var sp opentracing.Span
	if tracerIsSet {
		sp = getSpan(name, header)

		defer sp.Finish()
	} else {
		log.Printf("Tracer is not set")
	}

	// Execute function
	err := f()
	if err != nil && tracerIsSet {
		ext.Error.Set(sp, true)
		sp.LogEventWithPayload(fmt.Sprintf("Error"), err)
	}
	return err
}

func getSpan(name string, header http.Header) opentracing.Span {
	var out opentracing.Span
	log.Printf("Creating new span for %s", name)

	// if header is present try to extract traceID and use it
	// if there is no header create new span
	if header != nil {

		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(header))

		if err != nil {
			// If for whatever reason we can't join, go ahead an start a new root span.
			out = opentracing.StartSpan(name)
			log.Printf("Trace not found, creating new span, %v", err)

		} else {
			out = opentracing.StartSpan(name, opentracing.ChildOf(wireContext))
			log.Printf("Trace found, attaching to it")
		}
		out.Context().ForeachBaggageItem(func(k, v string) bool {
			log.Print("k:" + k)
			log.Print("v:" + v)
			return true
		})

	} else {
		out = opentracing.StartSpan(name)
		log.Printf("Trace not found, creating new span")

	}

	return out
}
