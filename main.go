package main

import (
	"log"
	"mysql/config"
	"mysql/model"
	"mysql/routes"
	"mysql/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.ConnectDatabase()
	utils.InitMetrics()
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			result := config.DB.Where("expires_at < ? ", time.Now()).
				Delete(&model.Session{})
			log.Printf("Session cleanup: removed %d expired/revoked sessions", result.RowsAffected)
		}
	}()

	r := gin.Default()
	r.Use(utils.MetricsMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(utils.SecurityHeaders())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	routes.SetupRoutes(r)

	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
