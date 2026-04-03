package repository

import (
	"gorm.io/gorm"
)

type SummaryResult struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

type CategoryBreakdown struct {
	Category string  `json:"category"`
	Type     string  `json:"type"`
	Total    float64 `json:"total"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type AnalyticsRepository interface {
	GetSummary(userID *uint) (*SummaryResult, error)
	GetCategoryBreakdown(userID *uint) ([]CategoryBreakdown, error)
	GetMonthlyTrends(userID *uint) ([]MonthlyTrend, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

func (r *analyticsRepository) GetSummary(userID *uint) (*SummaryResult, error) {
	var res SummaryResult
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'Income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'Expense' THEN amount ELSE 0 END), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN type = 'Income' THEN amount ELSE -amount END), 0) AS net_balance
		FROM transactions
	`
	
	if userID != nil {
		query += ` WHERE user_id = ?`
		err := r.db.Raw(query, *userID).Scan(&res).Error
		return &res, err
	}

	err := r.db.Raw(query).Scan(&res).Error
	return &res, err
}

func (r *analyticsRepository) GetCategoryBreakdown(userID *uint) ([]CategoryBreakdown, error) {
	var res []CategoryBreakdown
	query := `
		SELECT category, type, COALESCE(SUM(amount), 0) as total
		FROM transactions
	`
	
	if userID != nil {
		query += ` WHERE user_id = ? `
	}
	
	query += ` GROUP BY category, type ORDER BY total DESC`

	var err error
	if userID != nil {
		err = r.db.Raw(query, *userID).Scan(&res).Error
	} else {
		err = r.db.Raw(query).Scan(&res).Error
	}
	
	return res, err
}

func (r *analyticsRepository) GetMonthlyTrends(userID *uint) ([]MonthlyTrend, error) {
	var res []MonthlyTrend
	query := `
		SELECT 
			TO_CHAR(DATE_TRUNC('month', date), 'YYYY-MM') as month,
			COALESCE(SUM(CASE WHEN type = 'Income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'Expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
	`

	if userID != nil {
		query += ` WHERE user_id = ? `
	}
	
	query += ` GROUP BY DATE_TRUNC('month', date) ORDER BY month ASC`

	var err error
	if userID != nil {
		err = r.db.Raw(query, *userID).Scan(&res).Error
	} else {
		err = r.db.Raw(query).Scan(&res).Error
	}

	return res, err
}
