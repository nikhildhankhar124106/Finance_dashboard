package repository

import (
	"backend/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	Update(tx *models.Transaction) error
	Delete(id uint) error
	GetByID(id uint) (*models.Transaction, error)
	List(userID *uint, category, txType, date, search, sort, order string, page, pageSize int) ([]models.Transaction, int64, int, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.First(&tx, id).Error
	return &tx, err
}

func (r *transactionRepository) List(userID *uint, category, txType, date, search, sort, order string, page, pageSize int) ([]models.Transaction, int64, int, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	if date != "" {
		// Fixed Date matching logic accurately natively over PG timestamps
		query = query.Where("date = ?", date)
	}

	if search != "" {
		// Global search across category and notes mapping natively
		searchPattern := "%" + search + "%"
		query = query.Where("category ILIKE ? OR notes ILIKE ?", searchPattern, searchPattern)
	}

	// Dynamic Sorting
	allowedSortFields := map[string]string{
		"amount":   "amount",
		"date":     "date",
		"category": "category",
	}
	sortField, exists := allowedSortFields[sort]
	if !exists {
		sortField = "date" // Default sort
	}
	if order != "asc" {
		order = "desc" // Default order
	}
	query = query.Order(fmt.Sprintf("%s %s", sortField, order))

	// Get total count for pagination info
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&transactions).Error
	
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return transactions, total, totalPages, err
}
