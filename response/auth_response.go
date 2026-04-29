package response

import "mysql/model"

type AuthResponse struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Contact         string             `json:"contact"`
	DeviceName      string             `json:"device_name"`
	IpAddress       string             `json:"ipaddress"`
	Token           string             `josn:"token"`
	RoleID          uint               `json:"role_id"`
	Parts           []UserPartResponse `json:"parts"`
	ManageBranch    int                `json:"manage_branh"`
	Permissions     []model.Permission `json:"permissions"`
	BranchID        int                `json:"branch_id"`
	BranchLatitude  string             `json:"branch_latitude"`
	BranchLongitude string             `json:"branch_longitude"`
	BranchRadius    int                `json:"branch_radius"`
	EmployeeID      int                `json:"employee_id"`
}
