package main

import (
	"L0_task/internal/app"
	"L0_task/internal/config"
	"L0_task/internal/repository/dual"
	"L0_task/internal/service/rest"
	"L0_task/internal/subscriber/nats"
	lg "L0_task/pkg/logger"
	"L0_task/pkg/utils"
	"context"
	"log/slog"
)

func main() {

	cfg := config.MustLoad()
	logger := lg.GetLogger(cfg.Env)
	logger.Info("Starting the service", slog.String("env", cfg.Env))

	storage, err := dual.NewOrdersRepository(context.Background(), &cfg.Postgres)
	if err != nil {
		logger.Error("Failed to initialize storage", utils.WithErr(err))
		return
	}

	sub, err := nats.NewSubscriber(&cfg.Subscriber)
	if err != nil {
		logger.Error("Failed to initialize subscriber", utils.WithErr(err))
		return
	}

	service := rest.NewService(storage, sub, logger, cfg.Service)

	application := app.With(service)
	application.MustStart(context.Background())
}
