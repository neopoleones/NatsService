package rest

import (
	"L0_task/internal/models"
	"L0_task/internal/repository"
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type router struct {
	addr string
	e    *echo.Echo
	sl   *slog.Logger
	repo repository.OrdersRepository
}

func (r router) asServer() *http.Server {
	r.setupEndpoints()

	return &http.Server{
		Addr:         r.addr,
		Handler:      r.e,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
}

func (r router) setupEndpoints() {
	r.e.GET("/", r.statusEndpoint())

	api := r.e.Group("/api", r.loggerMiddleware)
	api.GET("/orders", r.listOrdersEndpoints())
	api.GET("/orders/:id", r.getOrderEndpoint())

}

func (r router) statusEndpoint() func(c echo.Context) error {
	type StatusResponse struct {
		Status string        `json:"status"`
		Uptime time.Duration `json:"uptime"`
	}
	savedTime := time.Now()

	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, StatusResponse{
			Status: "green",
			Uptime: time.Now().Sub(savedTime),
		})
	}
}

func (r router) listOrdersEndpoints() func(c echo.Context) error {
	type OrdersResponse struct {
		Error  string          `json:"error,omitempty"`
		Orders []*models.Order `json:"orders"`
		Status string          `json:"status"`
	}

	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		orders, err := r.repo.GetOrders(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, OrdersResponse{
				Error:  err.Error(),
				Status: "error",
			})
		} else {
			return c.JSON(http.StatusOK, OrdersResponse{
				Orders: orders,
				Status: "success",
			})
		}
	}
}

func (r router) getOrderEndpoint() func(c echo.Context) error {
	type OrderResponse struct {
		Error  string        `json:"error,omitempty"`
		Order  *models.Order `json:"order,omitempty"`
		Status string        `json:"status"`
	}

	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		uid := c.Param("id")
		order, err := r.repo.GetOrderByUID(ctx, uid)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, OrderResponse{
				Error:  err.Error(),
				Status: "error",
			})
		} else {
			return c.JSON(http.StatusOK, OrderResponse{
				Order:  order,
				Status: "success",
			})
		}
	}
}

func (r router) loggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		res := next(c)

		r.sl.Debug("got request", slog.Group(
			"req",
			slog.String("method", c.Request().Method),
			slog.String("url", c.Request().URL.Path),
			slog.String("status", strconv.Itoa(c.Response().Status)),
			slog.String("dur", time.Since(start).String()),
		))
		return res
	}
}

func getRouter(addr string, sl *slog.Logger, repo repository.OrdersRepository) router {
	return router{addr, echo.New(), sl, repo}
}
