package service

import (
	"backend/repository"
)

type AnalyticsService interface {
	GetSummary(userID *uint) (*repository.SummaryResult, error)
	GetCategoryBreakdown(userID *uint) ([]repository.CategoryBreakdown, error)
	GetMonthlyTrends(userID *uint) ([]repository.MonthlyTrend, error)
}

type analyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(repo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{repo: repo}
}

func (s *analyticsService) GetSummary(userID *uint) (*repository.SummaryResult, error) {
	return s.repo.GetSummary(userID)
}

func (s *analyticsService) GetCategoryBreakdown(userID *uint) ([]repository.CategoryBreakdown, error) {
	return s.repo.GetCategoryBreakdown(userID)
}

func (s *analyticsService) GetMonthlyTrends(userID *uint) ([]repository.MonthlyTrend, error) {
	return s.repo.GetMonthlyTrends(userID)
}
