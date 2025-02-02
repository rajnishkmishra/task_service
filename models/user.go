package models

import "time"

type User struct {
	ID          uint64    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LoginTime   time.Time `json:"login_time" gorm:"default:NULL"`
}
