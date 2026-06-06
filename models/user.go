package models
import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID uint 					`gorm:"primaryKey" json:"id"`
	Name string 				`gorm:"size:100;not null" json:"name"`
	Email string 				`gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string 			`gorm:"not null" json:"-"`
	CreatedAt time.Time 		`json:"created_at"`
	UpdatedAt time.Time 		`json:"updated_at"`
	DeletedAt gorm.DeletedAt 	`gorm:"index" json:"-"`
}

type RegisterUserRequest struct{
	Name string			`json:"name" binding:"required,min=2"`
	Email string		`json:"email" binding:"required,email"`
	Password string		`json:"password" binding:"required,min=8"`
}

type LoginRequest struct{
	Email string		`json:"email" binding:"required,email"`
	Password string		`json:"password" binding:"required"`
}

type AuthResponse struct{
	Token string `json:"token"`
	User User `json:"user"`
}