package response

type AuthResponse struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	Contact      string             `json:"contact"`
	Token        string             `josn:"token"`
	RoleID       uint               `json:"role_id"`
	Parts        []UserPartResponse `json:"parts"`
	ManageBranch int                `json:"manage_branh"`
}
