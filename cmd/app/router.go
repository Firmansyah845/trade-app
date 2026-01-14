package app

import (
	"awesomeProjectCr/internal/database"
	"awesomeProjectCr/internal/handler"
	"awesomeProjectCr/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm/module/apmchi/v2"
)

func newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(chiMiddleware.RequestID)

	router.Use(middleware.Recover())

	router.Use(apmchi.Middleware())

	return router

}

func router() *chi.Mux {
	router := newRouter()

	router.Get("/ping", handler.Ping)

	dbConnection := database.DBConnection

	postgreDB := dbConnection[database.PostgresDb]

	handlerR := handler.NewHandler(postgreDB)

	router.Get("/health-check", handlerR.HealthCheck)

	router.Group(func(s chi.Router) {
		s.Route("/api/v1", func(rApi chi.Router) {

			rApi.Post("/order", handlerR.CreateOrder)
		})
	})

	return router

}
