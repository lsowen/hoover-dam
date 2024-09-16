package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/lsowen/hoover-dam/pkg/api/service"
	"github.com/lsowen/hoover-dam/pkg/config"
	"github.com/lsowen/hoover-dam/pkg/db"
)

func Serve(ctx context.Context, cfg config.Config) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(AuthMiddleware(cfg))

	database, err := db.NewDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}

	server := NewAPIService(database)
	handler := service.HandlerWithOptions(server, service.ChiServerOptions{
		BaseURL:    "/api/v1",
		BaseRouter: r,
	})
	return handler, nil
}
