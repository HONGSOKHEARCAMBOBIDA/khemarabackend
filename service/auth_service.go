package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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
	RefreshToken(input request.RefreshTokenRequest, c *gin.Context) (*response.AuthResponse, error)
	//  Logout(refreshToken string) error
	GetSessions(userID uint) ([]model.Session, error)
	RevokeSession(sessionID uint, userID uint) error
	RevokeAllSessions(userID uint) error
	Register(id int, input request.RegisterRequest, c *gin.Context) error
	GetUserByBranch(id int) ([]response.UserResponse, error)
	UpdateUser(id int, input request.UserRequestUpdate) error
	ChangePassword(id int, input request.NewPassword) error
	GetUserByID(id int) ([]response.UserResponseUpdate, error)
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
	var user model.User
	if err := s.db.
		Where("(contact = ? OR email = ? OR username = ?) AND is_active = ?",
			input.UserName, input.UserName, input.UserName, 1).
		First(&user).Error; err != nil {
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ ឬ អ្នកប្រើប្រាស់ត្រូវបានបិទគណនី")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	); err != nil {
		utils.Redis.Incr(utils.Ctx, key)
		utils.Redis.Expire(utils.Ctx, key, 10*time.Minute)
		return nil, errors.New("ព័ត៌មានមិនត្រឹមត្រូវ")
	}
	utils.Redis.Del(utils.Ctx, key)

	var userparts []response.UserPartResponse
	if err := s.db.Table("user_parts up").
		Select("up.id AS id, p.id AS part_id, p.name AS part_name, p.display_name AS part_display_name").
		Joins("JOIN parts p ON p.id = up.part_id").
		Where("up.user_id = ?", user.ID).
		Scan(&userparts).Error; err != nil {
		return nil, err
	}

	var permissions []model.Permission
	if err := s.db.Table("permissions p").
		Select("p.id AS id, p.name AS name").
		Joins("JOIN role_has_permissions rhp ON rhp.permission_id = p.id").
		Where("rhp.role_id = ? AND p.name IN ?", user.RoleID, []string{
			"update.user", "add.user", "edit.salary", "add.salary",
			"add.loan", "add.payroll", "change.shift.pattern",
			"change.day.off", "approve.leave",
		}).
		Scan(&permissions).Error; err != nil {
		return nil, err
	}

	accessExpiry := time.Now().Add(1 * time.Minute)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"contact": user.Contact,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenStr, err := accessToken.SignedString(utils.Jwtkey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshTokenStr := hex.EncodeToString(refreshTokenBytes)
	tokenPrefix := refreshTokenStr[:16]
	hashedRefresh, err := bcrypt.GenerateFromPassword([]byte(refreshTokenStr), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	var sessionCount int64
	s.db.Model(&model.Session{}).
		Where("user_id = ? AND expires_at > ?", user.ID, time.Now()).
		Count(&sessionCount)

	if sessionCount >= 5 {
		s.db.Where("user_id = ? AND expires_at > ?", user.ID, time.Now()).
			Order("created_at ASC").
			Limit(1).
			Delete(&model.Session{})
	}

	session := model.Session{
		UserID:       uint(user.ID),
		RefreshToken: string(hashedRefresh),
		TokenPrefix:  tokenPrefix,
		DeviceName:   deviceName,
		IPAddress:    ipAddress,
		LastActive:   time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// 11. Build and return response
	resp := &response.AuthResponse{
		ID:           user.ID,
		Name:         user.UserName,
		Contact:      user.Contact,
		DeviceName:   deviceName,
		IpAddress:    ipAddress,
		RoleID:       uint(user.RoleID),
		Parts:        userparts,
		ManageBranch: user.ManageBranch,
		Permissions:  permissions,
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}

	return resp, nil
}

func (s *authservice) RefreshToken(input request.RefreshTokenRequest, c *gin.Context) (*response.AuthResponse, error) {
	if len(input.RefreshToken) < 16 {
		return nil, errors.New("Invalid refresh token")
	}
	prefix := input.RefreshToken[:16]
	var session model.Session
	err := s.db.Where("token_prefix = ? AND expires_at > ?",
		prefix, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), []byte(input.RefreshToken)); err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}

	newRefreshBytes := make([]byte, 32)
	rand.Read(newRefreshBytes)
	newRefreshStr := hex.EncodeToString(newRefreshBytes)
	newHash, _ := bcrypt.GenerateFromPassword([]byte(newRefreshStr), bcrypt.DefaultCost)
	newPrefix := newRefreshStr[:16]

	s.db.Model(&session).Updates(model.Session{
		RefreshToken: string(newHash),
		TokenPrefix:  newPrefix,
		LastActive:   time.Now(),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
		IPAddress:    c.ClientIP(),
	})

	// 4. Issue new access token
	var user model.User
	s.db.First(&user, session.UserID)

	accessExpiry := time.Now().Add(1 * time.Minute)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role_id": user.RoleID,
		"exp":     accessExpiry.Unix(),
	}
	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(utils.Jwtkey)

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshStr,
		ExpiresIn:    int64(accessExpiry.Unix()),
	}, nil
}

// func (s *authservice) Logout(refreshToken string) error {
//     // same bcrypt search loop as RefreshToken()
//     // then: s.db.Model(matched).Update("is_revoked", true)
// }

func (s *authservice) GetSessions(userID uint) ([]model.Session, error) {
	var sessions []model.Session
	s.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		userID, false, time.Now()).
		Order("last_active DESC").Find(&sessions)
	return sessions, nil
}

func (s *authservice) RevokeSession(sessionID uint, userID uint) error {
	return s.db.Model(&model.Session{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Update("is_revoked", true).Error
}

func (s *authservice) RevokeAllSessions(userID uint) error {
	return s.db.Model(&model.Session{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (s *authservice) GetUserByBranch(id int) ([]response.UserResponse, error) {
	var user []response.UserResponse
	db := s.db.Table("users u").
		Select(`
		e.id AS id,
		e.name_kh AS name
	`).
		Joins("LEFT JOIN employees e ON e.id = u.employee_id").Where("u.branch_id = ?", id)
	if err := db.Order("e.id DESC").Scan(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authservice) Register(id int, input request.RegisterRequest, c *gin.Context) error {

	var uploadedFiles []string
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			helper.DeleteFiles(uploadedFiles)
		}
	}()

	profilePath, err := helper.SaveImage(c, "profile_image", "public/profileimage")
	if err != nil {
		tx.Rollback()
		return err
	}

	uploadedFiles = append(uploadedFiles, profilePath)
	qrcodePath, err := helper.SaveImage(c, "qr_code_bank_account", "public/qrcodeimage")
	if err != nil {
		tx.Rollback()
		return err
	}

	uploadedFiles = append(uploadedFiles, qrcodePath)
	educationPaths, err := helper.SaveImages(c, "education_image", "public/educationimage")
	if err != nil {
		tx.Rollback()
		return err
	}

	uploadedFiles = append(uploadedFiles, educationPaths...)

	employee := model.Employee{
		Code:           utils.GenerateEmployeeCode(),
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
		helper.DeleteFiles(uploadedFiles)
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
		helper.DeleteFiles(uploadedFiles)
		return err
	}

	count := len(input.EducationLevelID)

	for i := 0; i < count; i++ {
		var imagePath string
		if i < len(educationPaths) {
			imagePath = educationPaths[i]
		}

		var majorField string
		if i < len(input.MajorField) {
			majorField = input.MajorField[i]
		}

		var startDate string
		if i < len(input.StartDateEducation) {
			startDate = input.StartDateEducation[i]
		}

		var endDate string
		if i < len(input.EndDateEducation) {
			endDate = input.EndDateEducation[i]
		}

		var note string
		if i < len(input.NoteEducation) {
			note = input.NoteEducation[i]
		}

		employeeEducation := model.EmployeeEducation{
			EmployeeID:       employee.ID,
			EducationLevelID: input.EducationLevelID[i],
			MajorField:       majorField,
			StartDate:        startDate,
			EndDate:          endDate,
			Note:             note,
			Image:            imagePath,
			CreateBy:         id,
		}

		if err := tx.Create(&employeeEducation).Error; err != nil {
			tx.Rollback()
			helper.DeleteFiles(uploadedFiles)
			return err
		}
	}

	if len(input.CompanyName) > 0 {

		for i := 0; i < len(input.CompanyName); i++ {

			var positiontitl string
			if i < len(input.PositionTitle) {
				positiontitl = input.PositionTitle[i]
			}

			var StartDate string
			if i < len(input.StartDate) {
				StartDate = input.StartDate[i]
			}

			var EndDate string
			if i < len(input.EndDate) {
				EndDate = input.EndDate[i]
			}

			var Jd string
			if i < len(input.JobDescription) {
				Jd = input.JobDescription[i]
			}

			employeeworkexperience := model.EmployeeWorkExperience{
				EmployeeID:     employee.ID,
				CompanyName:    input.CompanyName[i],
				PositionTitle:  positiontitl,
				StartDate:      StartDate,
				EndDate:        EndDate,
				JobDescription: Jd,
				CreateBy:       id,
			}

			if err := tx.Create(&employeeworkexperience).Error; err != nil {
				tx.Rollback()
				helper.DeleteFiles(uploadedFiles)
				return err
			}
		}
	}

	salary := model.Salary{
		EmployeeID:    employee.ID,
		BaseSalary:    input.BaseSalary,
		WorkDay:       input.WorkDay,
		DailyRate:     input.DailyRate,
		EffectiveDate: time.Now().Format("2006-01-02 15:04:05"),
		ExpireDate:    nil,
		CurrencyID:    input.CurrencyID,
		Isactive:      true,
	}

	if err := tx.Create(&salary).Error; err != nil {
		tx.Rollback()
		helper.DeleteFiles(uploadedFiles)
		return err
	}

	username := strings.ReplaceAll(strings.ToLower(input.NameEn), "         ", "")
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
		helper.DeleteFiles(uploadedFiles)
		return err
	}

	if len(input.PartIDs) > 0 {
		for _, part := range input.PartIDs {
			if err := tx.Create(&model.UserPart{
				UserID: int(user.ID),
				PartID: uint(part),
			}).Error; err != nil {
				tx.Rollback()
				helper.DeleteFiles(uploadedFiles)
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
				helper.DeleteFiles(uploadedFiles)
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
				helper.DeleteFiles(uploadedFiles)
				return err
			}
		}
	} else {
		return errors.New("DayOfWeeks and IsDayOff must be same length")
	}

	if err := tx.Commit().Error; err != nil {
		helper.DeleteFiles(uploadedFiles)
		return err
	}
	return nil

}

func (s *authservice) UpdateUser(id int, input request.UserRequestUpdate) error {
	log.Printf("[UpdateUser] id=%d input=%+v", id, input)
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := tx.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"branch_id":     input.BranchID,
		"role_id":       input.RoleID,
		"manage_branch": input.ManageBranch,
	})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if err := tx.Where("user_id = ?", id).Delete(&model.UserPart{}).Error; err != nil {
		tx.Rollback()
		return err

	}

	for _, partID := range input.PartIDs {
		if err := tx.Create(&model.UserPart{
			UserID: id,
			PartID: uint(partID),
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("user_id = ?", id).Delete(&model.UserBranch{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if input.ManageBranch == 2 {
		for _, branchID := range *input.BranchIDs {
			if err := tx.Create(&model.UserBranch{
				UserID:   uint(id),
				BranchID: uint(branchID),
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (s *authservice) ChangePassword(id int, input request.NewPassword) error {

	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with id %d not found", id)
		}
		return err
	}
	hash := utils.HasPassword(input.NewPassword)
	if err := s.db.Model(&user).Update("password", hash).Error; err != nil {
		return err
	}

	return nil
}

func (s *authservice) GetUserByID(id int) ([]response.UserResponseUpdate, error) {
	var users []response.UserResponseUpdate

	if err := s.db.Table("users u").
		Select(`
            u.id AS id,
            u.username AS username,
            b.id AS branch_id,
            b.name AS branch_name,
            r.id AS role_id,
            r.name AS role_name,
            u.manage_branch AS manage_branch
        `).
		Joins("LEFT JOIN branches b ON b.id = u.branch_id").
		Joins("LEFT JOIN roles r ON r.id = u.role_id").
		Where("u.id = ?", id).
		Scan(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	for i := range users {
		var userParts []model.UserPart
		if err := s.db.Table("user_parts up").
			Select("up.id AS id, up.user_id AS user_id, up.part_id AS part_id").
			Where("up.user_id = ?", id).
			Scan(&userParts).Error; err != nil {
			return nil, err
		}
		users[i].Userpart = userParts

		var userBranches []model.UserBranch
		if err := s.db.Table("user_branches ub").
			Select("ub.id AS id, ub.user_id AS user_id, ub.branch_id AS branch_id").
			Where("ub.user_id = ?", id).
			Scan(&userBranches).Error; err != nil {
			return nil, err
		}
		users[i].UserBranch = userBranches
	}

	return users, nil
}
