package service

import (
	"errors"
	"fmt"
	"mysql/config"
	"mysql/helper"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error)
	Register(id int, input request.RegisterRequest, c *gin.Context) error
}

type authservice struct {
	db *gorm.DB
}

func NewAuthService() AuthService {
	return &authservice{
		db: config.DB,
	}
}

func (s *authservice) Login(input request.AuthRequest, c *gin.Context) (*response.AuthResponse, error) {
	deviceName := c.Request.UserAgent()
	ipAddress := c.ClientIP()
	key := "login_attempt:" + input.UserName
	attempts, _ := utils.Redis.Get(utils.Ctx, key).Int()
	if attempts >= 5 {
		return nil, errors.New("អ្នកព្យាយាមចូលច្រើនពេក សូមព្យាយាមម្តងទៀតក្រោយ 10 នាទី")
	}
	// 1. Find user
	var user model.User
	if err := s.db.
		Where("(contact = ? OR email = ? OR username = ?) AND is_active = ?",
			input.UserName, input.UserName, input.UserName, 1).
		First(&user).Error; err != nil {
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ ឬ អ្នកប្រើប្រាស់ត្រូវបានបិទគណនី")
	}

	// 2. Check password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {

		utils.Redis.Incr(utils.Ctx, key)
		utils.Redis.Expire(utils.Ctx, key, 10*time.Minute)

		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ")
	}
	utils.Redis.Del(utils.Ctx, key)

	// 4. Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	var userparts []response.UserPartResponse
	if err := s.db.Table("user_parts up").
		Select("up.id AS id,p.id AS part_id,p.name AS part_name,p.display_name AS part_display_name").
		Joins("JOIN parts p ON p.id = up.part_id").
		Where("up.user_id =?", user.ID).Scan(&userparts).Error; err != nil {
		return nil, err
	}

	var permissions []model.Permission

	if err := s.db.Table("permissions p").Select("p.id AS id,p.name AS name,p.display_name AS display_name,p.group_name AS group_name,p.short_name AS short_name").
		Joins("JOIN role_has_permissions rhp ON rhp.permission_id = p.id").Where("rhp.role_id =?", user.RoleID).Scan(&permissions).Error; err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"contact": user.Contact,
		"role_id": user.RoleID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(utils.Jwtkey)
	if err != nil {
		return nil, err
	}

	// 5. Build response
	resp := &response.AuthResponse{
		ID:           user.ID,
		Name:         user.UserName,
		Contact:      user.Contact,
		DeviceName:   deviceName,
		IpAddress:    ipAddress,
		Token:        tokenStr,
		RoleID:       uint(user.RoleID),
		Parts:        userparts,
		ManageBranch: user.ManageBranch,
		Permissions:  permissions,
	}

	return resp, nil
}

func (s *authservice) Register(id int, input request.RegisterRequest, c *gin.Context) error {

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()

		}
	}()

	profilePath, err := helper.SaveImage(c, "profile_image", "public/profile")
	if err != nil {
		tx.Rollback()
		return err
	}

	qrcodePath, err := helper.SaveImage(c, "qrcode_image", "public/qrcode")
	if err != nil {
		tx.Rollback()
		return err
	}

	educationPath, err := helper.SaveImage(c, "education_image", "public/education")
	if err != nil {
		tx.Rollback()
		return err
	}

	employee := model.Employee{
		NameEn:         input.NameEn,
		NameKh:         input.NameKh,
		NationalID:     input.NationalID,
		Gender:         input.Gender,
		PositionID:     input.PositionID,
		HireDate:       input.HireDate,
		PromoteDate:    input.PromoteDate,
		IsPromote:      false,
		EmployeeTypeID: input.EmployeeTypeID,
		Isactive:       true,
		OfficeID:       input.OfficeID,
		CreateBy:       id,
	}

	if err := tx.Create(&employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	employeeprofile := model.EmployeeProfile{
		EmployeeID:              employee.ID,
		PositionLevelID:         input.PositionLevelID,
		DoB:                     input.DoB,
		VillageIdOfBirth:        input.VillageIdOfBirth,
		MaterialStatus:          input.MaterialStatus,
		ProfileImage:            profilePath,
		VillageIDCurrentAddress: input.VillageIDCurrentAddress,
		FamilyPhone:             input.FamilyPhone,
		BankName:                input.BankName,
		BankAccountNumber:       input.BankAccountNumber,
		QrCodeBankAccount:       qrcodePath,
		Note:                    "",
		ReportoID:               input.ReportoID,
		WifeName:                input.WifeName,
		HusBanName:              input.HusBanName,
		SonNumber:               input.SonNumber,
		DaughterNumber:          input.DaughterNumber,
		CreateBy:                id,
	}

	if err := tx.Create(&employeeprofile).Error; err != nil {
		tx.Rollback()
		return err
	}

	employeeeducation := model.EmployeeEducation{
		EmployeeID:       employee.ID,
		EducationLevelID: input.EducationLevelID,
		MajorField:       input.MajorField,
		StartDate:        input.StartDate,
		EndDate:          input.EndDate,
		Note:             "",
		Image:            educationPath,
		CreateBy:         id,
	}

	if err := tx.Create(&employeeeducation).Error; err != nil {
		tx.Rollback()
		return err
	}

	employeeworkexperience := model.EmployeeWorkExperience{
		EmployeeID:     employee.ID,
		CompanyName:    input.CompanyName,
		PositionTitle:  input.PositionTitle,
		StartDate:      input.StartDateEducation,
		EndDate:        input.EndDateEducation,
		JobDescription: input.JobDescription,
		CreateBy:       id,
	}

	if err := tx.Create(&employeeworkexperience).Error; err != nil {
		tx.Rollback()
		return err
	}

	salary := model.Salary{
		EmployeeID:    employee.ID,
		BaseSalary:    input.BaseSalary,
		WorkDay:       input.WorkDay,
		DailyRate:     input.DailyRate,
		EffectiveDate: time.Now(),
		ExpireDate:    time.Now(),
		CurrencyID:    input.CurrencyID,
		Isactive:      true,
	}

	if err := tx.Create(&salary).Error; err != nil {
		tx.Rollback()
		return err
	}

	username := strings.ToLower(input.NameEn)
	email := fmt.Sprintf("%s168@gmail.com", username)
	password := utils.HasPassword("123456")
	user := model.User{
		UserName:     username,
		Email:        email,
		Password:     password,
		Contact:      input.Contact,
		BranchID:     input.BranchID,
		RoleID:       input.RoleID,
		EmployeeID:   employee.ID,
		Isactive:     true,
		ManageBranch: input.ManageBranch,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(input.PartIDs) > 0 {
		for _, part := range input.PartIDs {
			if err := tx.Create(&model.UserPart{
				UserID: int(user.ID),
				PartID: uint(part),
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(input.BranchIDs) > 0 {
		for _, branch := range input.BranchIDs {
			if err := tx.Create(&model.UserBranch{
				UserID:   uint(user.ID),
				BranchID: uint(branch),
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(input.Dayofweeks) > 0 && len(input.Dayofweeks) == len(input.Isdayoff) {
		for i, day := range input.Dayofweeks {
			shift := model.ShiftPattern{
				EmployeeID:  employee.ID,
				DayOfWeekID: day,
				ShiftID:     input.ShiftID,
				Isdayoff:    input.Isdayoff[i],
			}

			if err := tx.Create(&shift).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		return errors.New("DayOfWeeks and IsDayOff must be same length")
	}

	return tx.Commit().Error

}
