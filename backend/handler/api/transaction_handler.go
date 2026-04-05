package api

import (
	"net/http"
	"strconv"

	_ "backend/domain/models"
	"backend/pkg/apperrors"
	"backend/service"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	txService service.TransactionService
}

func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		txService: txService,
	}
}

// CreateTransactionInput represents the expected payload for creating a transaction
type CreateTransactionInput struct {
	Amount   float64 `json:"amount" binding:"required,gt=0" example:"150.00"`
	Type     string  `json:"type" binding:"required,oneof=Income Expense" example:"Expense"`
	Category string  `json:"category" binding:"required" example:"Groceries"`
	Date     string  `json:"date" binding:"required" example:"2026-05-01"` // Format: YYYY-MM-DD
	Notes    string  `json:"notes" example:"Weekly groceries shopping"`
}

// UpdateTransactionInput represents the expected payload for updating a transaction completely or partially
type UpdateTransactionInput struct {
	Amount   float64 `json:"amount" binding:"omitempty,gt=0" example:"170.00"`
	Type     string  `json:"type" binding:"omitempty,oneof=Income Expense" example:"Expense"`
	Category string  `json:"category" example:"Groceries"`
	Date     string  `json:"date" example:"2026-05-02"`
	Notes    string  `json:"notes" example:"Updated weekly groceries amount"`
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Creates a new financial transaction under the active user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transaction body CreateTransactionInput true "Transaction payload"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} apperrors.AppError
// @Failure 401 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input CreateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apperrors.NewValidationError(err.Error()))
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.Error(apperrors.NewUnauthorizedError("Unauthorized context"))
		return
	}
	userID := userIDVal.(uint)

	tx, err := h.txService.CreateTransaction(userID, input.Amount, input.Type, input.Category, input.Date, input.Notes)
	if err != nil {
		c.Error(apperrors.NewInternalError(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// GetTransactions godoc
// @Summary Fetch a list of transactions
// @Description Returns transaction paginated list optionally filtered by category, type, and date
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Limits per page" default(10)
// @Param category query string false "Filter by category"
// @Param type query string false "Filter by type (Income/Expense)"
// @Param date query string false "Filter by exact date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /v1/transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	category := c.Query("category")
	txType := c.Query("type")
	dateStr := c.Query("date")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	var userID *uint
	roleVal, exists := c.Get("role")
	if !exists {
		c.Error(apperrors.NewUnauthorizedError("Role not set"))
		return
	}
	
	if roleVal.(string) != "Admin" {
		uid, uidExists := c.Get("user_id")
		if !uidExists {
			c.Error(apperrors.NewUnauthorizedError("UserID not set"))
			return
		}
		uidCast := uid.(uint)
		userID = &uidCast
	}

	transactions, total, totalPages, err := h.txService.GetTransactions(userID, category, txType, dateStr, search, sort, order, page, pageSize)
	if err != nil {
		c.Error(apperrors.NewInternalError("Failed to fetch transactions"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":        transactions,
		"total":       total,
		"total_pages": totalPages,
		"page":        page,
		"limit":       pageSize,
	})
}

// ExportTransactions godoc
// @Summary Export transactions to CSV
// @Description Generates a CSV file of the user's financial transactions
// @Tags transactions
// @Produce text/csv
// @Security BearerAuth
// @Success 200 {string} string "CSV data"
// @Failure 401 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /v1/transactions/export [get]
func (h *TransactionHandler) ExportTransactions(c *gin.Context) {
	var userID *uint
	roleVal, exists := c.Get("role")
	if !exists {
		c.Error(apperrors.NewUnauthorizedError("Role not set"))
		return
	}

	if roleVal.(string) != "Admin" {
		uid, _ := c.Get("user_id")
		uidCast := uid.(uint)
		userID = &uidCast
	}

	data, err := h.txService.ExportTransactions(userID)
	if err != nil {
		c.Error(apperrors.NewInternalError("Failed to export transactions"))
		return
	}

	c.Header("Content-Disposition", "attachment; filename=transactions_export.csv")
	c.Data(http.StatusOK, "text/csv", data)
}

// UpdateTransaction godoc
// @Summary Update transaction
// @Description Modifies an existing transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Param transaction body UpdateTransactionInput true "Updated values"
// @Success 200 {object} models.Transaction
// @Failure 400 {object} apperrors.AppError
// @Failure 404 {object} apperrors.AppError
// @Router /v1/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(apperrors.NewValidationError("Invalid transaction ID parameter"))
		return
	}

	var input UpdateTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	tx, err := h.txService.UpdateTransaction(uint(id), input.Amount, input.Type, input.Category, input.Date, input.Notes)
	if err != nil {
		c.Error(apperrors.NewNotFoundError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, tx)
}

// DeleteTransaction godoc
// @Summary Delete transaction
// @Description Performs hard deletion of transaction record
// @Tags transactions
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} apperrors.AppError
// @Failure 500 {object} apperrors.AppError
// @Router /v1/transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.Error(apperrors.NewValidationError("Invalid transaction ID parameter"))
		return
	}

	err = h.txService.DeleteTransaction(uint(id))
	if err != nil {
		c.Error(apperrors.NewInternalError("Failed to delete transaction"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
