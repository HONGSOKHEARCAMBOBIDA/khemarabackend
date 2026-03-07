package response

type PermissionWithAssignedRole struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Group       string `json:"group"`
	Short       int    `json:"short"`
	Assigned    bool   `json:"assigned"`
}
