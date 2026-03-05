package model

type User struct {
	ID           int    `json:"id"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Contact      string `json:"contact"`
	BranchID     int    `json:"branch_id"`
	RoleID       int    `json:"role_id"`
	EmployeeID   int    `json:"employee_id"`
	Isactive     bool   `json:"is_active"`
	ManageBranch int    `json:"manage_branch"`
}
