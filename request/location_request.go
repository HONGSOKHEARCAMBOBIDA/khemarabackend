package request

type LocationRequest struct {
	BranchID   int    `form:"branch_id"`
	EmployeeID int    `form:"employee_id"`
	Latitude   string `form:"latitude"`
	Longitude  string `form:"longitude"`
	Note       string `form:"note"`
}
