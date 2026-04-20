package response

type UserBranchResponse struct {
	ID         int    `json:"id"`
	BranchID   int    `json:"branch_id"`
	BranchName string `json:"branch_name"`
}
