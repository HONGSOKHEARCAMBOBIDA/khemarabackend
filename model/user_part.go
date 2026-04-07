package model

type UserPart struct {
	ID     uint `json:"id"`
	UserID int  `json:"user_id"`
	PartID uint `json:"part_id"`
}
