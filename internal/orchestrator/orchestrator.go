package orchestrator

import (
	"calculator/internal/orchestrator/handler"
	"calculator/internal/orchestrator/impl/sqlite"
	"calculator/internal/orchestrator/impl/sqlite_expression_storage"
	"calculator/internal/orchestrator/impl/sqlite_task_storage"

	"calculator/internal/orchestrator/use_cases/scheduler"
	"calculator/internal/orchestrator/web"
	"calculator/internal/shared/configs"
	"calculator/pkg/logger"
	"calculator/pkg/metrics/entities"
	"calculator/pkg/metrics/healthz"
	"calculator/proto/calculator/proto"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"calculator/pkg/middlewares"

	"google.golang.org/grpc"
)

var appInfo = &entities.AppInfo{
	Name:         "Expression Calculator",
	BuildVersion: "1.0.0",
	BuildTime:    time.Now().Format(time.RFC3339),
}

// Orchestrator represents the orchestrator.

type App struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	conf       *configs.Config
}

// NewOrchestrator creates a new instance of the Orchestrator.
func New(conf *configs.Config) (*App, error) {
	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	app := new(App)
	app.conf = conf

	db, err := sqlite.NewSQLiteDB("calculator.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite database: %v", err)
	}

	// Create separate storages using the same database
	expressionStorage := sqlite_expression_storage.NewStorage(db)
	taskStorage := sqlite_task_storage.NewTaskPool(db)

	scheduler := scheduler.NewScheduler(expressionStorage, taskStorage, app.conf)

	// Setup HTTP server
	httpHandler := handler.NewHandler(scheduler)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)
	healthz.RegisterRoutes(mux, appInfo)
	web.RegisterRoutes(mux)

	wrappedMux := middlewares.MakeLoggingMiddleware(mux)
	wrappedMux = middlewares.PanicRecoveryMiddleware(wrappedMux)

	app.httpServer = &http.Server{
		Handler:      wrappedMux,
		Addr:         ":" + strconv.Itoa(conf.Server.HttpPort),
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}

	// Setup gRPC server
	app.grpcServer = grpc.NewServer()
	grpcHandler := handler.NewGRPCHandler(scheduler)
	proto.RegisterCalculatorServer(app.grpcServer, grpcHandler)

	return app, nil
}

func (a *App) Run() error {
	// Start HTTP server
	go func() {
		logger.Info("starting http server...")
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(a.conf.Server.GrpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen for gRPC: %v", err)
	}
	logger.Info("starting grpc server...")
	return a.grpcServer.Serve(lis)
}

func (a *App) stop(ctx context.Context) error {
	logger.Info("shutdowning server...")
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}

	// Stop gRPC server
	a.grpcServer.GracefulStop()

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
