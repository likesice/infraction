package infraction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"infraction.mageis.net/internal/api"
	"infraction.mageis.net/internal/config"
	"infraction.mageis.net/internal/data"
	"infraction.mageis.net/migrations"
)

var (
	version     = "1.0.0-alpha"
	serviceName = "infraction"
)

func newServeCommand() *cli.Command {
	return &cli.Command{
		Name:            "serve",
		Usage:           "launch the server daemon",
		UsageText:       "infraction serve <infraction_config_file>",
		Description:     "Launch the Infraction server daemon.",
		Action:          serveAction,
		HideHelpCommand: true,
	}
}

func serveAction(ctx *cli.Context) error {
	if ctx.NArg() != 1 || ctx.Args().First() == "" {
		cli.ShowSubcommandHelpAndExit(ctx, 2)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger = logger.With("version", version, "service", serviceName)
	slog.SetDefault(logger)

	cfg, err := config.Configure(ctx.Args().First())
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse configuration : %q", err))
		return err
	}

	cfg.Db.Dsn = "./data.db"
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("could not connect to database: %q", err))
		return err
	}
	cfg.Db.Db = db

	err = migrations.Migrate(cfg.Db.Db)
	if err != nil {
		logger.Error(fmt.Sprintf("could not execute migrations: %q", err))
		return err
	}

	infractionApi := api.New(logger, cfg, data.NewRepositories(cfg.Db.Db))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Handler:      infractionApi.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Info(fmt.Sprintf("starting %s server on %s", cfg.Env, srv.Addr))

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("encountered error in listen: %s\n", err))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("received shutdown signal. Shutting down ...")

	endCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(endCtx); err != nil {
		logger.Error("encountered error while shutting down server:", err)
	}

	select {
	case <-endCtx.Done():
	}
	logger.Info("server shutdown")
	return nil
}
func openDB(cfg *config.Config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("sqlite3", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
