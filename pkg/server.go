package pkg

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	jaegerLog "github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/pkg/jaeger/log"
	tracingNethttp "github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerPrometheus "github.com/uber/jaeger-lib/metrics/prometheus"

	mongodbTracer "github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/pkg/mongodb/tracer"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/public"
)

type key int

const (
	requestIDKey key = 0
)

var (
	healthy int32
)

type indexValues struct {
	Count    int
	Hostname string
}

// RunServer serves the site and handles server health and metrics
func RunServer(listenAddr string, persistence bool, mongodbURI string, mongodbDatabase string, mongodbCollection string, mongodbDocumentID string, jaegerServiceName string) {

	logger := log.New()

	metricsFactory := jaegerPrometheus.New()
	cfg, err := jaegerConfig.FromEnv()
	if err != nil {
		logger.Fatal(err)
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = jaegerServiceName
	}
	tracer, tracerClose, err := cfg.NewTracer(
		jaegerConfig.Metrics(metricsFactory),
		jaegerConfig.Logger(jaegerLog.NewLogrusAdapter(logger)),
	)
	if err != nil {
		logger.Fatal(err)
	}
	opentracing.SetGlobalTracer(tracer)
	defer tracerClose.Close()

	funcMap := template.FuncMap{
		"str":  strconv.Itoa,
		"stoa": stoa,
	}
	indexTemplate := template.Must(template.New("index").Funcs(funcMap).Parse(public.IndexTpl))

	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	hostname, err := os.Hostname()
	if err != nil {
		log.WithError(err).Fatal()
	}

	counter := counterHandler{
		Active: persistence,
	}
	if persistence {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mtr := mongodbTracer.NewTracer()
		monitor := &event.CommandMonitor{
			Started:   mtr.HandleStartedEvent,
			Succeeded: mtr.HandleSucceededEvent,
			Failed:    mtr.HandleFailedEvent,
		}
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI).SetMonitor(monitor))
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				log.Fatal(err)
			}
		}()
		counter.mongodbClient = client
		counter.mongodbDatabase = mongodbDatabase
		counter.mongodbCollection = mongodbCollection
		counter.mongodbDocumentID = mongodbDocumentID

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Connected to MongoDB at %s", mongodbURI)
	}

	router := http.NewServeMux()
	router.Handle("/", tracingNethttp.Middleware(tracer, std.Handler("", mdlw, index(indexTemplate, counter, hostname))))
	router.Handle("/healthz", healthz())
	router.Handle("/metrics", promhttp.Handler())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	logWriter := logger.Writer()
	defer logWriter.Close()

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     stdlog.New(logWriter, "", 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Println("Server stopped")
}

func index(template *template.Template, counterHandle counterHandler, hostname string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		indexData := indexValues{
			Hostname: hostname,
		}
		if counterHandle.Active {
			indexData.Count = counterHandle.Content(r.Context()).Count
		}

		err := template.Execute(w, indexData)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			_, err := w.Write([]byte("OK"))
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func stoa(s string) []rune {
	return []rune(s)
}
