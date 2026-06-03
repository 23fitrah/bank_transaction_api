package dto

type User struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
