package user

import "time"

type User struct {
	Id        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"updated_at,omitempty"`
	Role      string    `json:"role,omitempty"`
	Token     *string   `json:"token"`
}
