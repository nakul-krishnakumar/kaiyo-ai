package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id,omitempty" mapstructure:"id"`
	Email         string    `json:"email,omitempty" mapstructure:"email"`
	Password      string    `json:"password,omitempty" mapstructure:"password"`
	Name          string    `json:"name,omitempty" mapstructure:"name"`
	/* 
		ID -- primary key
		Password - password hash
	*/

	GoogleID      string    `json:"google_id,omitempty" mapstructure:"google_id"`
	TwitterID     string    `json:"twitter_id,omitempty" mapstructure:"twitter_id"`

	EmailVerified bool      `json:"email_verified,omitempty" mapstructure:"email_verified"`
	IsActive      bool      `json:"is_active,omitempty" mapstructure:"is_active"` 
	/* 
		IsActive -- for soft deletion 
	*/

	CreatedAt     time.Time `json:"created_at,omitempty" mapstructure:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" mapstructure:"updated_at"`
	LastLoginAt   time.Time `json:"last_login_at,omitempty" mapstructure:"last_login_at"`
}
