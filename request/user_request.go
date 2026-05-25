package request

type UserRequestUpdate struct {
	BranchID     int    `json:"branch_id"`
	RoleID       int    `json:"role_id"`
	ManageBranch int    `json:"manage_branch"`
	PartIDs      []int  `json:"part_ids"`
	BranchIDs    *[]int `json:"branch_ids"`
}

type NewPassword struct {
	NewPassword string `json:"newpassword"`
}
