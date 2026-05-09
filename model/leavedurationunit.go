package model

type LeaveDurationUnit struct {
	ID        int    `json:"id"`
	Code      string `json:"code" gorm:"column:code"`
	NameEn    string `json:"name_en" gorm:"column:name_en"`
	NameKh    string `json:"name_km" gorm:"column:name_km"`
	ToMinute  int    `json:"to_minutes" gorm:"column:to_minutes"`
	SortOrder int    `json:"sort_order" gorm:"column:sort_order"`
}
