package model

type PayrollApproval struct {
	ID         int    `json:"id"`
	PayrollID  int    `json:"payroll_id" gorm:"column:payroll_id"`
	ApproveBy  int    `json:"approved_by" gorm:"column:approved_by"`
	Status     string `json:"status" gorm:"column:status"`
	Comment    string `json:"comment" gorm:"column:comment"`
	ActionDate string `json:"action_date" gorm:"column:action_date"`
	StepOrder  int    `json:"step_order" gorm:"column:step_order"`
}
