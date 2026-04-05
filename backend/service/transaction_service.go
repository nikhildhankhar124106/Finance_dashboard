package service

import (
	"errors"
	"time"

	"backend/domain/models"
	"backend/repository"
	"fmt"
)

type TransactionService interface {
	CreateTransaction(userID uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error)
	UpdateTransaction(id uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error)
	DeleteTransaction(id uint) error
	GetTransactions(userID *uint, category, txType, date, search, sort, order string, page, pageSize int) ([]models.Transaction, int64, int, error)
	ExportTransactions(userID *uint) ([]byte, error)
}

type transactionService struct {
	txRepo   repository.TransactionRepository
	userRepo repository.UserRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, userRepo repository.UserRepository) TransactionService {
	return &transactionService{
		txRepo:   txRepo,
		userRepo: userRepo,
	}
}

func (s *transactionService) CreateTransaction(userID uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error) {
	// Validate user existence natively before creating transaction to prevent foreign key violations
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found, cannot create transaction")
	}

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

	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) UpdateTransaction(id uint, amount float64, txType string, category string, dateStr string, notes string) (*models.Transaction, error) {
	tx, err := s.txRepo.GetByID(id)
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

	if err := s.txRepo.Update(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) DeleteTransaction(id uint) error {
	return s.txRepo.Delete(id)
}

func (s *transactionService) GetTransactions(userID *uint, category, txType, dateStr, search, sort, order string, page, pageSize int) ([]models.Transaction, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return s.txRepo.List(userID, category, txType, dateStr, search, sort, order, page, pageSize)
}

func (s *transactionService) ExportTransactions(userID *uint) ([]byte, error) {
	// Exporting to CSV naturally mapping all fields including financial categorization
	transactions, _, _, err := s.txRepo.List(userID, "", "", "", "", "date", "desc", 1, 10000)
	if err != nil {
		return nil, err
	}

	// Simple CSV construction natively over byte buffers
	csvHeader := "ID,Amount,Type,Category,Date,Notes\n"
	csvContent := csvHeader
	for _, tx := range transactions {
		line := fmt.Sprintf("%d,%.2f,%s,%s,%s,%s\n", 
			tx.ID, tx.Amount, tx.Type, tx.Category, tx.Date.Format("2006-01-02"), tx.Notes)
		csvContent += line
	}

	return []byte(csvContent), nil
}
