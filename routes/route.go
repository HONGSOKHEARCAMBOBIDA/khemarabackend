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
	employeetypecontroller := controller.NewEmployeeTypeController()
	educationlevelcontroller := controller.NewEducationLevelController()
	branchcontroller := controller.NewBranchController()
	departmentcontroller := controller.NewDepartmentController()
	r.Static("/clientimage", "./public/clientimage")
	r.POST("/login", authcontroller.Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{ //Role
		auth.GET(route.ViewRole, middleware.PermissionMiddleware(permission.ViewRole), rolecontroller.GetRole)
		auth.POST(route.AddRole, middleware.PermissionMiddleware(permission.AddRole), rolecontroller.CreateRole)
		auth.PUT(route.EditRole, middleware.PermissionMiddleware(permission.EditRole), rolecontroller.UpdateRole)
		auth.PUT(route.ChangeStatusRole, middleware.PermissionMiddleware(permission.ChangeStatusRole), rolecontroller.ChangeStatusRole)
		// Rolehaspermission
		auth.POST(route.AddPermissionTORole, middleware.PermissionMiddleware(permission.AddPermissionTORole), rolehaspermissioncontroller.CreateRoleHasPermission)
		auth.DELETE(route.RemovePermissionFromRole, middleware.PermissionMiddleware(permission.RemovePermissionFromRole), rolehaspermissioncontroller.DeleteRoleHasPermission)
		auth.GET(route.ViewRoleHasPermission, middleware.PermissionMiddleware(permission.ViewRoleHasPermission), rolehaspermissioncontroller.GetRoleHasPermission)

		// Employee Type
		auth.GET(route.ViewEmployeeType, middleware.PermissionMiddleware(permission.ViewEmployeeType), employeetypecontroller.GetEmployeeType)
		auth.POST(route.AddEmployeeType, middleware.PermissionMiddleware(permission.AddEmployeeType), employeetypecontroller.CreateEmployeeType)
		auth.PUT(route.UpdateEmployeeType, middleware.PermissionMiddleware(permission.UpdateEmployeeType), employeetypecontroller.UpdateEmployeeType)
		auth.PUT(route.ChangeStatusEmployeeType, middleware.PermissionMiddleware(permission.ChangeStatusEmployeeType), employeetypecontroller.ChangeStatusEmployeeType)

		// Education Level
		auth.GET(route.ViewEducationLevel, middleware.PermissionMiddleware(permission.ViewEducationLevel), educationlevelcontroller.GetEducationLevel)
		auth.POST(route.AddEducationLevel, middleware.PermissionMiddleware(permission.AddEducationLevel), educationlevelcontroller.CreateEducationLevel)
		auth.PUT(route.UpdateEducationLevel, middleware.PermissionMiddleware(permission.UpdateEducationLevel), educationlevelcontroller.UpdateEducationLevel)
		auth.PUT(route.ChangeStatusEducationLevel, middleware.PermissionMiddleware(permission.ChangeStatusEducationLevel), educationlevelcontroller.ChangeStatusEducationLevel)

		// Branch
		auth.GET(route.ViewBranch, middleware.PermissionMiddleware(permission.ViewBranch), branchcontroller.GetBranch)
		auth.POST(route.AddBranch, middleware.PermissionMiddleware(permission.AddBranch), branchcontroller.CreateBranch)
		auth.PUT(route.UpdateBranch, middleware.PermissionMiddleware(permission.UpdateBranch), branchcontroller.UpdateBranch)
		auth.PUT(route.ChangeStatusBranch, middleware.PermissionMiddleware(permission.ChangeStatusBranch), branchcontroller.ChangeStatusBranch)

		// Department
		auth.GET(route.ViewDepartment, middleware.PermissionMiddleware(permission.ViewDepartment), departmentcontroller.GetDepartment)
		auth.POST(route.AddDepartment, middleware.PermissionMiddleware(permission.AddDepartment), departmentcontroller.CreateDepartment)
		auth.PUT(route.UpdateDepartment, middleware.PermissionMiddleware(permission.UpdateDepartment), departmentcontroller.UpdateDepartment)
		auth.PUT(route.ChangeStatusDepartment, middleware.PermissionMiddleware(permission.ChangeStatusDepartment), departmentcontroller.ChangeStatusDepartment)
	}
}
