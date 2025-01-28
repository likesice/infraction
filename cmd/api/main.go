package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"infraction.mageis.net/internal/data"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	dsn string
}
type Application struct {
	registry *data.Registry
	logger   *slog.Logger
	config   config
}

func main() {

	app := Application{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
		config: config{dsn: "/mnt/c/Users/Martin/AppData/Local/Temp/data.db"},
	}
	slog.SetDefault(app.logger)
	db, err := openDB(app.config)
	m, err := migrate.New("file://./migrations", "sqlite3://"+app.config.dsn)
	if err != nil {
		panic(err)
	}
	err = m.Steps(0)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
	}

	app.registry = data.NewRegistry(db)
	srv := http.Server{
		Addr:         ":8080",
		IdleTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      app.routes(),
	}

	go func() {
		log.Panic(srv.ListenAndServe())
	}()

	app.logger.Info(fmt.Sprintf("started app at addr: %s", srv.Addr))
	var stopCh = make(chan os.Signal, 2)

	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stopCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		app.logger.Error("error while shutting down app", err)
		return
	}
	app.logger.Info("shut down app")
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
