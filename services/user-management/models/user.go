package models

import (
	"time"
)

type User struct {
	ID              int       `json:"id"`
	AlpacaAccountID string    `json:"alpaca_account_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
