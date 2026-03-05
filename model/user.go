package model

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Isactive bool   `json:"is_active"`
}
