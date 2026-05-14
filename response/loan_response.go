package response

type LoanResponse struct {
	ID               int                `json:"id"`
	Code             string             `json:"code"`
	BranchID         int                `json:"branch_id"`
	BranchName       string             `json:"branch_name"`
	EmployeeID       int                `json:"employee_id"`
	EmployeeName     string             `json:"employee_name"`
	EmployeeGender   int                `json:"employee_gender"`
	EmployeeDob      string             `json:"employee_dob" gorm:"column:employee_dob"`
	EmployeeContact  string             `json:"employee_contact"`
	EmployeeCode     string             `json:"employee_code"`
	LoanAmount       string             `json:"loan_amount"`
	CurrencyID       int                `json:"currency_id"`
	CurrencyCode     string             `json:"currency_code"`
	ApproveDate      string             `json:"approve_date"`
	LoanStartDate    string             `json:"loan_start_date"`
	LoanEndDate      string             `json:"loan_end_date"`
	NumberofLoan     int                `json:"number_of_loan" gorm:"column:number_of_loan"`
	ApproveByID      int                `json:"approve_by_id"`
	ApproveByName    string             `json:"approve_by_name"`
	LoanPurpose      string             `json:"loan_purpose"`
	Status           int                `json:"loan_status" gorm:"column:loan_status"`
	LoanDuration     int                `json:"loan_duration"`
	ScheduleResponse []ScheduleResponse `json:"schedule" gorm:"-"`
}
