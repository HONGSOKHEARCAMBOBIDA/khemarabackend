package request

type RoleHasPermissionRequestCreate struct {
	RoleID        uint  `json:"role_id"`
	PermissionIDs []int `json:"permission_ids"`
}

type RoleHasPermissionRequestDelete struct {
	RoleID        uint  `json:"role_id"`
	PermissionIDs []int `json:"permission_ids"`
}
