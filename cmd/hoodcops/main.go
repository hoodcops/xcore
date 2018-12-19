package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

const (
	development = "development"
	production  = "production"
	driver      = "mysql"
)

var env = struct {
	Port                     int           `envconfig:"PORT" required:"true"`
	Environment              string        `envconfig:"ENVIRONMENT" default:"development"`
	ServiceDSN               string        `envconfig:"SERVICE_DSN" required:"true"`
	DbConnMaxLife            time.Duration `envconfig:"DB_CONN_MAX_LIFE" default:"14400s"`
	DbMaxIdleConns           int           `envconfig:"DB_MAX_IDLE_CONNS" default:"50"`
	DbMaxOpenConns           int           `envconfig:"DB_MAX_OPEN_CONNS" default:"100"`
	City                     string        `envconfig:"CITY" required:"true"`
	Locale                   string        `envconfig:"LOCALE" default:"en"`
	TwilioVerificationAPIKey string        `envconfig:"TWILIO_VERIFICATION_API_KEY" required:"true"`
}{}

func init() {
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalf("failed loading env vars : %v", err)
	}
}

func initLogger(environment string) (*zap.Logger, error) {
	if environment == production {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}

func main() {
	logger, err := initLogger(env.Environment)
	if err != nil {
		log.Fatalf("failed initializing logger : %v", err)
	}

	if env.Environment == development {
		logger.Info("loaded env vars successfully", zap.Any("configs", env))
	}

	dbConn, err := sqlx.Open(driver, env.ServiceDSN)
	if err != nil {
		logger.Fatal("failed initializing db connection", zap.Error(err))
	}

	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		logger.Fatal("failed pinging database", zap.Error(err))
	}

	dbConn.SetConnMaxLifetime(env.DbConnMaxLife)
	dbConn.SetMaxIdleConns(env.DbMaxIdleConns)
	dbConn.SetMaxOpenConns(env.DbMaxOpenConns)

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		logger.Fatal("failed binding to port", zap.Int("port", env.Port))
	}
	defer listener.Close()

	server := http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Handler:           nil,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	connsClosed := make(chan struct{})
	go func() {
		defer close(connsClosed)

		recv := <-sigs
		logger.Info("received signal, shutting down", zap.Any("signal", recv.String))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("failed shutting down server", zap.Error(err))
		}
	}()

	url := fmt.Sprintf("http://%s", listener.Addr())
	logger.Info("server listening on ", zap.String("url", url))

	if err = server.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			logger.Fatal("failed starting server", zap.Error(err))
		}
	}

	<-connsClosed
	logger.Info("server shutdown successfully")
}
