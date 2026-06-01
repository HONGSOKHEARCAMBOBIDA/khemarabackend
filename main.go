package main

import (
	"mysql/config"
	"mysql/routes"
	"mysql/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	config.ConnectDatabase()

	// Create Gin router
	r := gin.Default()
	// if want to protect file size
	// r.MaxMultipartMemory = 8 << 20
	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://7ml45f42-5173.asse.devtunnels.ms"}, // your frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(utils.SecurityHeaders())
	// Set up routes
	routes.SetupRoutes(r)

	// Start server
	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
