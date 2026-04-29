package request

type LocationRequest struct {
	BranchID        int    `form:"branch_id"`
	BranchLatitude  string `form:"branch_latitude"`
	BranchLongitude string `form:"branch_longitude"`
	BranchRadius    int    `form:"branch_radius"`
	EmployeeID      int    `form:"employee_id"`
	Latitude        string `form:"latitude"`
	Longitude       string `form:"longitude"`
	Note            string `form:"note"`
}
