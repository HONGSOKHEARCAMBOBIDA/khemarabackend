package model

type UserBranch struct {
	ID       uint `json:"id"`
	UserID   uint `json:"user_id"`
	BranchID uint `json:"branch_ids"`
}
