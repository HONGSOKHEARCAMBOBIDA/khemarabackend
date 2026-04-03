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
	positioncontroller := controller.NewPositionController()
	positionlevelcontroller := controller.NewPositionLevelController()
	currencycontroller := controller.NewCurrencyController()
	currencypaircontroller := controller.NewCurrencyPairController()
	exchangeratecontroller := controller.NewExchangeRateController()
	managebranchcontroller := controller.NewManageBranchController()
	provincecontroller := controller.NewProvinceController()
	districtcontroller := controller.NewDistrictController()
	communcecontroller := controller.NewCommunceController()
	villagecontroller := controller.NewVillageController()
	dayofweekcontroller := controller.NewDayOfWeekController()
	officecontroller := controller.NewOfficeController()
	shiftcontroller := controller.NewShiftController()
	shiftsessioncontroller := controller.NewShiftSessionController()
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

		// Position
		auth.GET(route.ViewPosition, middleware.PermissionMiddleware(permission.ViewPosition), positioncontroller.GetAllPosition)
		auth.GET(route.ViewPositionByDepartment, middleware.PermissionMiddleware(permission.ViewPosition), positioncontroller.GetByDepartmentID)
		auth.POST(route.AddPosition, middleware.PermissionMiddleware(permission.AddPosition), positioncontroller.CreatePosition)
		auth.PUT(route.UpdatePosition, middleware.PermissionMiddleware(permission.UpdatePosition), positioncontroller.UpdatePosition)
		auth.PUT(route.ChangeStatusPosition, middleware.PermissionMiddleware(permission.UpdatePosition), positioncontroller.ChangeStatusPosition)

		// PositionLevel
		auth.GET(route.ViewPositionLevel, middleware.PermissionMiddleware(permission.ViewPositionLevel), positionlevelcontroller.GetPositionLevel)
		auth.POST(route.AddPositionLevel, middleware.PermissionMiddleware(permission.AddPositionLevel), positionlevelcontroller.CreatePositionLevel)
		auth.PUT(route.UpdatePositionLevel, middleware.PermissionMiddleware(permission.UpdatePositionLevel), positionlevelcontroller.UpdatePositionLevel)
		auth.PUT(route.ChangeStatusPositionLevel, middleware.PermissionMiddleware(permission.ChangeStatusPositionLevel), positionlevelcontroller.ChangeStatusPositionLevel)

		// Currency
		auth.GET(route.ViewCurrency, middleware.PermissionMiddleware(permission.ViewCurrency), currencycontroller.GetCurrency)
		auth.POST(route.AddCurrency, middleware.PermissionMiddleware(permission.AddCurrency), currencycontroller.CreateCurrency)
		auth.PUT(route.UpdateCurrency, middleware.PermissionMiddleware(permission.UpdateCurrency), currencycontroller.UpdateCurrency)
		auth.PUT(route.ChangeStatusCurrency, middleware.PermissionMiddleware(permission.ChangeStatusCurrency), currencycontroller.ChangeStatusCurrency)

		// CurrencyPair
		auth.GET(route.ViewCurrencyPair, middleware.PermissionMiddleware(permission.ViewCurrencyPair), currencypaircontroller.GetCurrencypair)
		auth.POST(route.AddCurrencyPair, middleware.PermissionMiddleware(permission.AddCurrencyPair), currencypaircontroller.CreateCurrencyPair)
		auth.PUT(route.UpdateCurrencyPair, middleware.PermissionMiddleware(permission.UpdateCurrencyPair), currencypaircontroller.UpdateCurrencyPaire)
		auth.PUT(route.ChangeStatusCurrencyPair, middleware.PermissionMiddleware(permission.ChangeStatusCurrencyPair), currencypaircontroller.ChangeStatusCurrencyPair)

		// ExchangeRate
		auth.GET(route.ViewExchangeRate, middleware.PermissionMiddleware(permission.ViewExchangeRate), exchangeratecontroller.GetExchangeRate)
		auth.POST(route.AddExchangeRate, middleware.PermissionMiddleware(permission.AddExchangeRate), exchangeratecontroller.CreateExchangeRate)
		auth.PUT(route.UpdateExchangeRate, middleware.PermissionMiddleware(permission.UpdateExchangeRate), exchangeratecontroller.UpdateExchangeRate)
		auth.PUT(route.ChangeStatusExchangeRate, middleware.PermissionMiddleware(permission.ChangeStatusExchangeRate), exchangeratecontroller.ChangeStatusExchangeRate)

		// ManageBranch
		auth.GET(route.ViewManageBranch, middleware.PermissionMiddleware(permission.ViewManageBranch), managebranchcontroller.GetManageBranch)

		// Address
		auth.GET(route.ViewProvince, middleware.PermissionMiddleware(permission.ViewProvince), provincecontroller.GetProvince)
		auth.GET(route.ViewDistrict, middleware.PermissionMiddleware(permission.ViewDistrict), districtcontroller.GetDistrict)
		auth.GET(route.ViewCommunce, middleware.PermissionMiddleware(permission.ViewCommunce), communcecontroller.GetCommunce)
		auth.GET(route.ViewVillage, middleware.PermissionMiddleware(permission.ViewVillage), villagecontroller.GetVillage)

		// Dayofweek
		auth.GET(route.ViewDayofweek, middleware.PermissionMiddleware(permission.ViewDayofweek), dayofweekcontroller.GetDayOfWeek)

		// Office
		auth.GET(route.ViewOffice, middleware.PermissionMiddleware(permission.ViewOffice), officecontroller.GetAllOffice)

		// Shift
		auth.GET(route.ViewAllShift, middleware.PermissionMiddleware(permission.ViewShift), shiftcontroller.GetAllShift)
		auth.GET(route.ViewShfitByBranchID, middleware.PermissionMiddleware(permission.ViewShift), shiftcontroller.GetByBranchID)
		auth.POST(route.AddShift, middleware.PermissionMiddleware(permission.AddShift), shiftcontroller.CreateShift)
		auth.PUT(route.UpdateShift, middleware.PermissionMiddleware(permission.UpdateShift), shiftcontroller.UpdateShift)
		auth.PUT(route.ChangeStatusShift, middleware.PermissionMiddleware(permission.ChangeStatusShift), shiftcontroller.ChangeStatusShift)

		// ShiftSession
		auth.GET(route.ViewAllShiftSession, middleware.PermissionMiddleware(permission.ViewShift), shiftsessioncontroller.GetAllShiftSession)
		auth.GET(route.ViewShiftSessionByShiftID, middleware.PermissionMiddleware(permission.ViewShift), shiftsessioncontroller.GetByShiftID)
		auth.POST(route.AddShiftSession, middleware.PermissionMiddleware(permission.AddShift), shiftsessioncontroller.CreateShiftSession)
		auth.PUT(route.UpdateShiftSession, middleware.PermissionMiddleware(permission.UpdateShift), shiftsessioncontroller.UpdateShiftSession)
		auth.PUT(route.ChangeStatusShiftSession, middleware.PermissionMiddleware(permission.ChangeStatusShift), shiftsessioncontroller.ChangeStatusShiftSession)
	}
}
