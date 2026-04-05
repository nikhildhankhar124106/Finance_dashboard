package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin   UserRole = "Admin"
	RoleAnalyst UserRole = "Analyst"
	RoleViewer  UserRole = "Viewer"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Name         string         `gorm:"not null" json:"name"`
	Role         UserRole       `gorm:"type:varchar(20);default:'Viewer';not null" json:"role"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	Transactions []Transaction  `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
