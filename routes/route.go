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
	rolecontroller := controller.NewRoleController()
	rolehaspermissioncontroller := controller.NewRoleHasPermissionController()
	r.Static("/clientimage", "./public/clientimage")
	r.POST("/login", authcontroller.Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewRole), rolecontroller.GetRole)
		auth.POST(route.AddRole, middleware.PermissionMiddleware(permission.AddRole), rolecontroller.CreateRole)
		auth.PUT(route.EditRole, middleware.PermissionMiddleware(permission.EditRole), rolecontroller.UpdateRole)
		auth.PUT(route.ChangeStatusRole, middleware.PermissionMiddleware(permission.ChangeStatusRole), rolecontroller.ChangeStatusRole)
		auth.POST(route.AddPermissionTORole, middleware.PermissionMiddleware(permission.AddPermissionTORole), rolehaspermissioncontroller.CreateRoleHasPermission)
		auth.DELETE(route.RemovePermissionFromRole, middleware.PermissionMiddleware(permission.RemovePermissionFromRole), rolehaspermissioncontroller.DeleteRoleHasPermission)
		auth.GET(route.ViewRoleHasPermission, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.GetRoleHasPermission)
	}
}
