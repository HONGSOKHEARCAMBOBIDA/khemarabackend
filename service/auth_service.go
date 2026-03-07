package service

import (
	"errors"
	"mysql/config"
	"mysql/model"
	"mysql/request"
	"mysql/response"
	"mysql/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(input request.AuthRequest) (*response.AuthResponse, error)
}

type authservice struct {
	db *gorm.DB
}

func NewAuthService() AuthService {
	return &authservice{
		db: config.DB,
	}
}

func (s *authservice) Login(input request.AuthRequest) (*response.AuthResponse, error) {

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
		Token:        tokenStr,
		RoleID:       uint(user.RoleID),
		Parts:        userparts,
		ManageBranch: user.ManageBranch,
	}

	return resp, nil
}
