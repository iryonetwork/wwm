package tracing

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
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
		log.Printf("Error initializing tracker: %v", err)
		return MockCloser{}
	}

	opentracing.SetGlobalTracer(tracer)
	tracerIsSet = true

	return closer
}

func Middleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tracerIsSet {
			var sp opentracing.Span
			spanName := fmt.Sprintf("HTTP %s", r.URL.Path)
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
		} else {
			log.Printf("Tracer is not set")
		}
		h.ServeHTTP(w, r)
	}
}
