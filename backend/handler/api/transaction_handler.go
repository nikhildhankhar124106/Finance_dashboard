package api

import (
	"net/http"
	"strconv"

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

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input struct {
		Amount   float64 `json:"amount" binding:"required,gt=0"`
		Type     string  `json:"type" binding:"required,oneof=Income Expense"`
		Category string  `json:"category" binding:"required"`
		Date     string  `json:"date" binding:"required"` // Format: YYYY-MM-DD
		Notes    string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming the admin/user ID creating this is extracted from token
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
		return
	}
	userID := userIDVal.(uint)

	tx, err := h.txService.CreateTransaction(userID, input.Amount, input.Type, input.Category, input.Date, input.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	category := c.Query("category")
	txType := c.Query("type")
	dateStr := c.Query("date")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var userID *uint
	// In a real application, you might filter automatically if they are a viewer, vs an admin seeing everything.
	// For this example, viewers only see their own transactions, admins see all unless asked.
	roleVal, _ := c.Get("role")
	if roleVal.(string) != "Admin" {
		uid := c.MustGet("user_id").(uint)
		userID = &uid
	}

	transactions, total, err := h.txService.GetTransactions(userID, category, txType, dateStr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  transactions,
		"total": total,
		"page":  page,
		"limit": pageSize,
	})
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var input struct {
		Amount   float64 `json:"amount" binding:"omitempty,gt=0"`
		Type     string  `json:"type" binding:"omitempty,oneof=Income Expense"`
		Category string  `json:"category"`
		Date     string  `json:"date"`
		Notes    string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.txService.UpdateTransaction(uint(id), input.Amount, input.Type, input.Category, input.Date, input.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = h.txService.DeleteTransaction(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
