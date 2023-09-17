package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"testFIO/cmd/config"
	"testFIO/cmd/loggers"
	"testFIO/internal/app/api"
	"testFIO/internal/app/kafka"
	"testFIO/internal/handlers"
	"testFIO/internal/storage"
	"testFIO/internal/storage/repositories"
	"testFIO/pkg/client/postgresql"
)

type App struct {
	router *chi.Mux
	cfg    config.ServerConfig
	logger *loggers.Logger
}

func NewServerApp() *App {
	cfg := config.ServerConfigInit()
	router := chi.NewRouter()
	logger := loggers.NewLogger()

	return &App{
		router: router,
		cfg:    *cfg,
		logger: logger,
	}
}

func (a *App) Run() {

	//определение БД
	var store storage.Storage

	var err error
	//инициализация БД
	if len(a.cfg.DatabaseDSN) != 0 {
		client, err := postgresql.NewClient(context.Background(), 5, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
		a.logger.LogInfo("Client PostgreSQL:", "", "start client")
		store, err = repositories.NewPGSStore(client, &a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
		a.logger.LogInfo("BD PostgreSQL:", a.cfg.DatabaseDSN, "start PostrgreSQL")
	} else {
		store, err = repositories.NewRepository(&a.cfg, a.logger)
		if err != nil {
			a.logger.LogErr(err, "")
			os.Exit(1)
		}
	}

	service := storage.NewService(store)
	a.logger.LogInfo("Service:", "", "start service")

	external := api.NewExternalApi(a.cfg, *a.logger)
	consumer, err := kafka.NewConsumer(*service, external, a.cfg, *a.logger)
	go func() {
		a.logger.LogInfo("Consumer:", "", "start consumer")
		if err = consumer.Read(); err != nil {
			a.logger.LogErr(err, "")
		}
	}()

	handler := handlers.NewHandler(*service, *a.logger, a.cfg)
	//регистрация хендлера
	handler.Register(a.router)
	a.logger.LogInfo("Handlers:", "", "register handlers")

	srv := http.Server{
		Addr:    a.cfg.Addr,
		Handler: a.router,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.LogErr(err, "server not started")
		}
	}()
	a.logger.LogInfo("server is listen:", a.cfg.Addr, "start server")

	//gracefullshutdown
	<-done

	a.logger.LogInfo("", "", "server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		a.logger.LogErr(err, "Server Shutdown Failed")
	}
	a.logger.LogInfo("", "", "Server Exited Properly")

}
