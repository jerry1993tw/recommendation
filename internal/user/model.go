package user

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Password    string    `json:"-"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	VerifyToken string    `json:"verify_token" gorm:"default:null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
