/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"context"
	"errors"
	"fmt"
	cmlog "log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	// health "github.com/heptiolabs/healthcheck"

	"github.com/adrianpk/poslan/internal/amazon"
	"github.com/adrianpk/poslan/internal/config"
	c "github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/auth"
	"github.com/go-kit/kit/log"
	"github.com/heptiolabs/healthcheck"
	health "github.com/heptiolabs/healthcheck"
	zipkin "github.com/openzipkin/zipkin-go"
	reporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var (
	serviceName        = "poslanAuthentication"
	serviceHostPort    = "localhost:8000"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v2/spans"
)

// checkSigTerm listens to sigterm events.
func checkSigTerm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cmlog.Printf("[ERROR] service interrupted.")
	cancel()
}

func makeService(ctx context.Context, cfg *c.Config, log log.Logger) *service {
	return &service{
		name:   "Poslan",
		ctx:    ctx,
		cfg:    cfg,
		logger: log,
		auth:   auth.Server{Logger: log},
	}
}

// Init a service instance.
func (svc *service) Init() (s Service, err error) {
	svc.Disable()
	// FIX: Readiness & liveness checks temporarily disabled.
	svc.health = health.NewHandler()
	svc.health.AddReadinessCheck("ready", svc.ReadinessCheck())
	svc.health.AddLivenessCheck("heap-threshold", svc.HeapLivenessCheck(10))
	svc.health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(25))

	ok1 := initAmazon(svc)
	// ok2 := initSesgrid(s)

	if !<-ok1 {
		return nil, fmt.Errorf("Cannot initialize '%s' service", svc.name)
	}

	s = addLogging(svc, svc.logger)
	// s = addTracing(svc)
	s = addInstrumentation(svc, svc.logger)
	s = addAuthentication(svc, svc.logger, svc.auth)

	return s, nil
}

func initAmazon(svc *service) chan bool {
	ok := make(chan bool)
	go func() {
		defer close(ok)
		p, err := amazon.Init(svc.ctx, svc.cfg, svc.logger)
		if err != nil {
			svc.logger.Log(
				"level", config.LogLevel.Error,
				"package", "main",
				"method", "initAmazon",
				"message", "Cannot initialize Amazon SES client.",
				"error", err.Error(),
			)
			ok <- false
			return
		}
		svc.mux.Lock()
		svc.providers = append(svc.providers, p)
		svc.mux.Unlock()
		ok <- true
	}()
	return ok
}

// Middleware
func addLogging(svc Service, logger log.Logger) Service {
	if loggingOn {
		return loggingMiddleware{
			logger: logger,
			next:   svc}
	}
	return svc
}

func addInstrumentation(svc Service, logger log.Logger) Service {
	if instrumentationOn {
		m := instrumentationMeters()
		return instrumentationMiddleware{
			logger:         logger,
			requestCount:   m.ReqCount,
			requestLatency: m.ReqLatency,
			countResult:    m.CountResult,
			next:           svc,
		}
	}
	return svc
}

func addAuthentication(svc Service, logger log.Logger, auth auth.SecServer) Service {
	return authenticationMiddleware{
		logger: svc.Logger(),
		auth:   auth,
		next:   svc,
	}
}

func makeLogger() log.Logger {
	w := log.NewSyncWriter(os.Stdout)
	logger := log.NewLogfmtLogger(w)
	logger.Log("level", c.LogLevel.Info, "message", "Config Logger started.")
	return logger
}

func makeTracer() (*zipkin.Tracer, error) {
	r := reporter.NewReporter(zipkinHTTPEndpoint)
	ze, err := zipkin.NewEndpoint(serviceName, serviceHostPort)
	if err != nil {
		return nil, err
	}
	return zipkin.NewTracer(r, zipkin.WithLocalEndpoint(ze))
}

// Start the service.
func (svc *service) Start() {
	go svc.checkCancel()
	svc.StartProviders()
}

func (svc *service) checkCancel() {
	<-svc.ctx.Done()
	svc.StopProviders()
}

// StarMailer is used in service startup
// to start each configured provider.
func (svc *service) StartProviders() {
	for _, m := range svc.providers {
		m.Start()
	}
}

// StarMailer is used in service stop
// to stop each configured provider.
func (svc *service) StopProviders() {
	for _, m := range svc.providers {
		m.Stop()
	}
}

// IsReady is a readiness test for the service.
func (svc *service) ReadinessCheck() healthcheck.Check {
	return func() error {

		if !svc.ready {
			msg := fmt.Sprintf("%s service is not ready!", svc.name)
			svc.logger.Log("level", c.LogLevel.Warn, "message", msg)
			return errors.New(msg)
		}

		msg := fmt.Sprintf("%s service is ready", svc.name)
		svc.logger.Log("level", c.LogLevel.Info, "message", msg)

		return nil
	}
}

// HeapLivenessCheck is a heap allocation liveness test for the service.
func (svc *service) HeapLivenessCheck(maxMb uint64) healthcheck.Check {
	return func() error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		mb := toMb(m.Alloc)
		r := mb > maxMb

		if r {
			msg := fmt.Sprintf("%s is not in healthy state.", svc.name)
			return errors.New(msg)
		}

		return nil
	}
}

func toMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func checkError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 && msg[0] != "" {
			cmlog.Println("level", c.LogLevel.Error, "message", msg[0])
		}
		cmlog.Println("level", c.LogLevel.Error, "message", err.Error())
		os.Exit(1)
	}
}
