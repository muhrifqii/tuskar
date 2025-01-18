package domain

import "time"

type User struct {
	ID       int64  `json:"-" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"a_password"`

	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}

type (
	AuthnRequest struct {
		Username string `json:"username" validate:"required,x_username_or_email"`
		Password string `json:"password" validate:"required"`
	}

	AuthnResponse struct {
		AccessToken           string    `json:"access_token"`
		AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
		RefreshToken          string    `json:"-"`
		RefreshTokenExpiresAt time.Time `json:"-"`
	}
)
