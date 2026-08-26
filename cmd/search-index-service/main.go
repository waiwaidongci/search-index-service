package main

import (
	"context"
	"github.com/ali/go-0821/search-index-service/internal/adapter/cache"
	"github.com/ali/go-0821/search-index-service/internal/adapter/http"
	"github.com/ali/go-0821/search-index-service/internal/application"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/config"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/search-index-service/internal/infrastructure/repository"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	templates := application.NewService(repository.NewMemory(), cache.NewMemory())
	srv := httpadapter.New(templates, logging.New(), metrics.New())
	if err := httpadapter.Serve(ctx, cfg.HTTPAddr, srv.Handler()); err != nil {
		panic(err)
	}
}
