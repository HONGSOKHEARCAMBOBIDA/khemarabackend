package request

type LocationRequest struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Note      string `json:"note"`
}
