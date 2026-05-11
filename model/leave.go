package model

type Leave struct {
	ID            int    `json:"id"`
	EmployeeID    int    `json:"employee_id" gorm:"column:employee_id"`
	LeaveTypeID   int    `json:"leave_type_id" gorm:"column:leave_type_id"`
	StartDate     string `json:"start_date" gorm:"column:start_date"`
	EndDate       string `json:"end_date" gorm:"column:end_date"`
	BackDate      string `json:"back_date" gorm:"column:back_date"`
	Description   string `json:"description" gorm:"column:description"`
	StatusLeaveID int    `json:"status_leave_id" gorm:"column:status_leave_id"`
	ApproveByID   int    `json:"approve_by_id" gorm:"column:approve_by_id"`
	BranchID      int    `json:"branch_id" gorm:"column:branch_id"`
	CreateBy      int    `json:"create_by" gorm:"column:create_by"`
}
