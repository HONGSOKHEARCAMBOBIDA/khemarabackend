package model

type ApproveWorkflow struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StepOrder int    `json:"step_order"`
	RoleName  string `json:"role_name"`
}
