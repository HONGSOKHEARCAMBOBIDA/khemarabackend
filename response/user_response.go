package response

import "mysql/model"

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserResponseUpdate struct {
	UserName     string             `json:"username" gorm:"column:username"`
	BranchID     int                `json:"branch_id"`
	BranchName   string             `json:"branch_name"`
	RoleID       int                `json:"role_id"`
	RoleName     string             `json:"role_name"`
	ManageBranch int                `json:"manage_branch"`
	Userpart     []model.UserPart   `json:"user_part" gorm:"-"`
	UserBranch   []model.UserBranch `json:"user_branch" gorm:"-"`
}
