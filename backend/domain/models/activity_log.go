package models

import (
	"time"
)

type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Action    string    `json:"action"` // e.g., "CREATE", "UPDATE", "DELETE"
	Resource  string    `json:"resource"` // e.g., "Transaction", "User"
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
