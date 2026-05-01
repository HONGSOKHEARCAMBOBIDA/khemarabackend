package model

type Employee struct {
	ID             int    `json:"id"`
	Code           string `json:"code" gorm:"column:code"`
	NameEn         string `json:"name_en"`
	NameKh         string `json:"name_kh"`
	NationalID     string `json:"national_id_number" gorm:"column:national_id_number"`
	Gender         int    `json:"gender"`
	PositionID     int    `json:"position_id"`
	HireDate       string `json:"hire_date"`
	PromoteDate    string `json:"promote_date"`
	IsPromote      bool   `json:"is_promote"`
	EmployeeTypeID int    `json:"employee_type_id"`
	Isactive       bool   `json:"is_active" gorm:"column:is_active"`
	OfficeID       int    `json:"office_id"`
	CreateBy       int    `json:"create_by"`
	UpdateBy       int    `json:"update_by"`
}
