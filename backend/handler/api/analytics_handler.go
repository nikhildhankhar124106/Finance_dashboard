package api

import (
	"net/http"

	"backend/service"

	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsSvc service.AnalyticsService
}

func NewAnalyticsHandler(analyticsSvc service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsSvc: analyticsSvc,
	}
}

// resolveFilter extract userID if it's not a master admin requesting aggregated overview.
func resolveFilter(c *gin.Context) *uint {
	var userID *uint
	roleVal, _ := c.Get("role")
	
	// Admins and Analysts can view all summaries if no specific scope limit is passed. 
	// For actual production software we might want explicit scopes passed in query parameters.
	// For now, if role is Analyst or Admin, view all, else scoped to identity.
	if roleVal.(string) == "Viewer" {
		uid := c.MustGet("user_id").(uint)
		userID = &uid
	}
	return userID
}

// GetSummary godoc
// @Summary Fetch dashboard summary
// @Description Returns financial summary metrics (total income, total expense, balance) for the dashboard
// @Tags analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /v1/summary [get]
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	userID := resolveFilter(c)

	res, err := h.analyticsSvc.GetSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetCategoryBreakdown godoc
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /v1/category-breakdown [get]
func (h *AnalyticsHandler) GetCategoryBreakdown(c *gin.Context) {
	userID := resolveFilter(c)

	res, err := h.analyticsSvc.GetCategoryBreakdown(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category breakdown"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetMonthlyTrends godoc
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /v1/monthly-trends [get]
func (h *AnalyticsHandler) GetMonthlyTrends(c *gin.Context) {
	userID := resolveFilter(c)

	res, err := h.analyticsSvc.GetMonthlyTrends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monthly trends"})
		return
	}

	c.JSON(http.StatusOK, res)
}
