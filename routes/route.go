package routes

import (
	"mysql/constant/permission"
	"mysql/constant/route"
	"mysql/controller"
	"mysql/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	authcontroller := controller.NewAuthController()
	r.Static("/clientimage", "./public/clientimage")
	r.POST("/login", authcontroller.Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewRole))

	}
}
