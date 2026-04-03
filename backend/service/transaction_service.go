package service

import (
	"errors"
	"time"

	"backend/domain/models"
	"backend/repository"
)

type TransactionService interface {
	CreateTransaction(userID uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error)
	UpdateTransaction(id uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error)
	DeleteTransaction(id uint) error
	GetTransactions(userID *uint, category string, txType string, dateStr string, page int, pageSize int) ([]models.Transaction, int64, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(userID uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	tx := &models.Transaction{
		UserID:   userID,
		Amount:   amount,
		Type:     models.TransactionType(txType),
		Category: category,
		Date:     parsedDate,
		Notes:    notes,
	}

	if err := s.repo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) UpdateTransaction(id uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	if dateStr != "" {
		parsedDate, parseErr := time.Parse("2006-01-02", dateStr)
		if parseErr != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		tx.Date = parsedDate
	}
	
	if amount > 0 {
		tx.Amount = amount
	}
	if txType != "" {
		tx.Type = models.TransactionType(txType)
	}
	if category != "" {
		tx.Category = category
	}
	tx.Notes = notes

	if err := s.repo.Update(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) DeleteTransaction(id uint) error {
	return s.repo.Delete(id)
}

func (s *transactionService) GetTransactions(userID *uint, category string, txType string, dateStr string, page int, pageSize int) ([]models.Transaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.repo.List(userID, category, txType, dateStr, page, pageSize)
}
