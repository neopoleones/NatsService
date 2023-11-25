package app

import (
	"L0_task/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	srv service.Service
}

func (app App) waitAndUtilize(ch <-chan os.Signal, done context.CancelFunc) {
	defer app.srv.Utilize()
	defer done()

	sig := <-ch
	fmt.Printf("Got %v...", sig)
}

func (app App) MustStart(ctx context.Context) {
	// Graceful shutdown
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	ctx, done := context.WithCancel(ctx)

	// Run subscriber
	go func() {
		if err := app.srv.Subscribe(ctx); err != nil {
			panic(err)
		}
	}()

	// Run service
	go func() {
		if err := app.srv.Run(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
			fmt.Println("Bye")
		}
	}()

	app.waitAndUtilize(ch, done)
}

func With(srv service.Service) *App {
	return &App{srv: srv}
}
