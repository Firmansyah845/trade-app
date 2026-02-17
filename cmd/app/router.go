package app

import (
	"time"

	"awesomeProjectCr/internal/database"
	"awesomeProjectCr/internal/handler"
	"awesomeProjectCr/pkg/middleware"

	"github.com/gin-gonic/gin"
)

var rateLimiterStore = middleware.NewRateLimiterStore(
	60,            // 60 requests per minute
	10,            // burst up to 10
	5*time.Minute, // cleanup visitors inactive for 5 minutes
)

func newRouter() *gin.Engine {
	routerGinNew := gin.New()
	routerGinNew.Use(gin.Logger())
	routerGinNew.Use(middleware.RequestID())
	routerGinNew.Use(middleware.Recover())
	routerGinNew.Use(rateLimiterStore.Middleware())

	return routerGinNew
}

func router() *gin.Engine {
	routerInit := newRouter()

	routerInit.GET("/ping", handler.Ping)

	dbConnection := database.DBConnection
	postgreDB := dbConnection[database.PostgresDb]
	handlerR := handler.NewHandler(postgreDB)

	routerInit.GET("/health-check", handlerR.HealthCheck)

	v1 := routerInit.Group("/api/v1")
	{
		v1.POST("/order", handlerR.CreateOrder)
	}
	return routerInit
}
