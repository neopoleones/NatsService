package rest

import (
	"L0_task/internal/config"
	"L0_task/internal/models"
	"L0_task/internal/repository"
	"L0_task/internal/subscriber"
	"L0_task/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Service struct {
	repo repository.OrdersRepository
	sub  subscriber.Subscriber
	sl   *slog.Logger

	server   *http.Server
	settings *config.Service
}

func (s Service) Utilize() {
	ctx, done := context.WithTimeout(context.Background(), time.Second*4)
	defer done()

	s.sl.Info("Utilizing the service")
	_ = s.repo.Close()
	_ = s.server.Shutdown(ctx)
}

func (s Service) Subscribe(ctx context.Context) error {
	s.sl.Info("Subscriber is ready")

	var wg sync.WaitGroup
	for i := 0; i < 3; i += 1 {
		entry, err := s.sub.Subscribe()
		if err != nil {
			return err
		}
		s.enrollEntry(ctx, entry, &wg)
	}

	wg.Wait()
	s.sl.Warn("All subscribers are dead!")

	return nil
}

func (s Service) Run() error {
	if s.settings == nil {
		return errors.New("settings not applied")
	}

	s.sl.Info("Service is ready", slog.Int("port", s.settings.Port))
	return s.server.ListenAndServe()
}

func (s Service) enrollEntry(ctx context.Context, entry subscriber.Entry, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		raw, err := entry.PullMessage()
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				continue
			}

			s.sl.Warn("Subscriber died", utils.WithErr(err))
			return
		}

		var order models.Order
		if err := json.Unmarshal(raw.Data, &order); err != nil {
			s.sl.Warn("Incorrect order raw data", utils.WithErr(err))
		}

		if err := s.repo.AddOrder(ctx, &order); err != nil {
			s.sl.Warn("Failed to add order", utils.WithErr(err))
		}
	}
}

func NewService(repo repository.OrdersRepository, sub subscriber.Subscriber, sl *slog.Logger, cfg config.Service) Service {
	addr := fmt.Sprintf(":%d", cfg.Port)

	return Service{
		repo:     repo,
		sub:      sub,
		sl:       sl,
		settings: &cfg,
		server:   getRouter(addr, sl, repo).asServer(),
	}
}
