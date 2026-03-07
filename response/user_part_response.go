package response

type UserPartResponse struct {
	ID              int    `json:"id"`
	PartID          int    `json:"part_id"`
	PartName        string `json:"part_name"`
	PartDisplayName string `json:"part_display_name"`
}
