package orchestrator

import (
	"calculator/internal/orchestrator/handler"
	"calculator/internal/orchestrator/impl/memory_expression_storage"
	"calculator/internal/orchestrator/impl/memory_task_storage"
	"calculator/internal/orchestrator/use_cases/scheduler"
	"calculator/internal/orchestrator/web"
	"calculator/internal/shared/configs"
	"calculator/pkg/logger"
	"calculator/pkg/metrics/entities"
	"calculator/pkg/metrics/healthz"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"calculator/pkg/middlewares"
)

var appInfo = &entities.AppInfo{
	Name:         "Expression Calculator",
	BuildVersion: "1.0.0",
	BuildTime:    time.Now().Format(time.RFC3339),
}

// Orchestrator represents the orchestrator.

type App struct {
	server *http.Server
	conf   *configs.Config
}

// NewOrchestrator creates a new instance of the Orchestrator.
func New(conf *configs.Config) (*App, error) {
	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	var err error

	app := new(App)

	logger.Info("setting TZ ...")
	if err = os.Setenv("TZ", "UTC"); err != nil {
		logger.Error("failed to set UTC timezone", err)
		return nil, err
	}

	app.conf = conf
	app.server = &http.Server{
		Handler:      app.Router(),
		Addr:         ":" + strconv.Itoa(conf.Server.Port),
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}

	return app, nil
}

// Router returns the router for the orchestrator.
func (o *App) Router() http.Handler {

	mux := http.NewServeMux()
	storage := memory_expression_storage.NewStorage()
	task_pool := memory_task_storage.NewTaskPool()
	sheduler := scheduler.NewScheduler(storage, task_pool, o.conf)

	handler := handler.NewHandler(sheduler)
	handler.RegisterRoutes(mux)

	healthz.RegisterRoutes(mux, appInfo)

	web.RegisterRoutes(mux)

	result := middlewares.MakeLoggingMiddleware(mux)
	result = middlewares.PanicRecoveryMiddleware(result)

	return result
}

func (a *App) Run() error {
	logger.Info("starting http server...")
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server was stop with err: %w", err)
	}
	logger.Info("server was stop")

	return nil
}

func (a *App) stop(ctx context.Context) error {
	logger.Info("shutdowning server...")
	err := a.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}
	logger.Info("server was shutdown")
	return nil
}

func (a *App) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 30 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			logger.Error("graceful shutdown timed out... forcing exit")
			os.Exit(1)
		}
	}()

	err := a.stop(shutdownCtx)
	if err != nil {
		logger.Error("graceful shutdown timed out... forcing exit")
		os.Exit(1)
	}
	serverStopCtx()
	shutdownStopCtx()
}
