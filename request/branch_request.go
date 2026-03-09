package request

type BranchRequestCreate struct {
	Name      string `json:"name" bind:"required"`
	Latitude  string `json:"latitude" bind:"required"`
	Longitude string `json:"longitude" bind:"required"`
	Radius    int    `json:"radius" bind:"required"`
}

type BranchRequesUpdate struct {
	Name      *string `json:"name"`
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
	Radius    *int    `json:"radius"`
}
