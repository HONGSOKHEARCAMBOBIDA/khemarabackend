package request

type AuthRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
