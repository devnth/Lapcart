package models

import (
	"time"
)

// admin schema for admin table
type Admin struct {
	Admin_ID     int       `json:"user_id"`
	First_Name   string    `json:"first_name"`
	Last_Name    string    `json:"last_name"`
	Password     string    `json:"password"`
	Email        string    `json:"email"`
	Phone_Number int       `json:"phone_number"`
	CreatedAt    time.Time `json:"created_at"`
}
