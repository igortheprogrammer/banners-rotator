package main

import (
	"banners-rotator/internal/bandit"
	"banners-rotator/internal/config"
	"banners-rotator/internal/logger"
	"banners-rotator/internal/rmq"
	"banners-rotator/internal/rotator"
	internalgrpc "banners-rotator/internal/server/grpc"
	sqlstorage "banners-rotator/internal/storage/sql"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.dev.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewAppConfig(configFile)
	if err != nil {
		fmt.Printf("Critical app error: %v", err)
		os.Exit(1)
	}

	logg, err := logger.NewLogger(cfg.Logger.Level, []string{"stdout"})
	if err != nil {
		fmt.Printf("Critical app error: %v", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logg.Info("getting storage...")
	s, err := getStorage(ctx, cfg)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	conn, err := amqp.Dial(cfg.Rmq.Uri)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	p := rmq.NewRMQProducer(cfg.Rmq.Name, conn)
	err = p.Connect()
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	b := bandit.NewBandit()
	app := rotator.NewApp(s, p, b)
	srv := internalgrpc.NewRPCServer(logg, app, cfg.Api.Host, cfg.Api.Port)

	go func() {
		<-ctx.Done()
		srv.Stop()
	}()

	logg.Info("rotator is running...")

	if err := srv.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		os.Exit(1)
	}

	defer cancel()
}

func getStorage(ctx context.Context, cfg *config.AppConfig) (rotator.Storage, error) {
	storage, err := sqlstorage.NewStorage(ctx, cfg.Storage.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("get storage -> %w", err)
	}

	err = storage.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("get storage -> %w", err)
	}

	return storage, nil
}
