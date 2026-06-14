package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	Nik        string         `gorm:"size:20;uniqueIndex;not null" json:"nik"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	Email      string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone      string         `gorm:"size:20" json:"phone"`
	Department string         `gorm:"size:50;not null" json:"department"`
	Position   string         `gorm:"size:50;not null" json:"position"`
	Salary     float64        `gorm:"not null" json:"salary"`
	JoinDate   time.Time      `gorm:"not null" json:"join_date"`
	IsActive   bool           `gorm:"not null" json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateEmployeeRequest struct {
	Nik        string  `json:"nik" binding:"required,min=5,max=20"`
	Name       string  `json:"name" binding:"required,min=2,max=100"`
	Email      string  `json:"email" binding:"required,email"`
	Phone      string  `json:"phone" binding:"required"`
	Department string  `json:"department" binding:"required,min=2,max=50"`
	Position   string  `json:"position" binding:"required,min=2,max=50"`
	Salary     float64 `json:"salary" binding:"required,gt=0"`
	JoinDate   string  `json:"join_date" binding:"required"` // Format: YYYY-MM-DD
}

type UpdateEmployeeRequest struct {
	Nik        string  `json:"nik" binding:"omitempty,min=5,max=20"`
	Name       string  `json:"name" binding:"omitempty,min=2,max=100"`
	Email      string  `json:"email" binding:"omitempty,email"`
	Phone      string  `json:"phone" binding:"omitempty"`
	Department string  `json:"department" binding:"omitempty,min=2,max=50"`
	Position   string  `json:"position" binding:"omitempty,min=2,max=50"`
	Salary     float64 `json:"salary" binding:"omitempty,gt=0"`
	JoinDate   string  `json:"join_date" binding:"omitempty"` // Format: YYYY-MM-DD
	IsActive   *bool   `json:"is_active" binding:"omitempty"`
}

type EmployeeQueryParams struct {
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
	Search     string `form:"search"`
	Department string `form:"department"`
	IsActive   *bool  `form:"is_active"`
}
