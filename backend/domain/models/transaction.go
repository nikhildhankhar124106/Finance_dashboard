package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TypeIncome  TransactionType = "Income"
	TypeExpense TransactionType = "Expense"
)

type Transaction struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	UserID    uint            `gorm:"index;not null" json:"user_id"`
	Amount    float64         `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type      TransactionType `gorm:"type:transaction_type;not null" json:"type"`
	Category  string          `gorm:"type:varchar(100);index;not null" json:"category"`
	Date      time.Time       `gorm:"type:date;index;not null" json:"date"`
	Notes     string          `gorm:"type:text" json:"notes"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}
