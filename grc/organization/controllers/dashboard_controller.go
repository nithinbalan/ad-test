package controllers

import (
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct {
	service organizationServices.DashboardServiceInterface
	db      *gorm.DB
}

func NewDashboardController(service organizationServices.DashboardServiceInterface, db *gorm.DB) *DashboardController {
	return &DashboardController{
		service: service,
		db:      db,
	}
}

func (dc *DashboardController) GetDashboardSummary(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	summary, err := dc.service.GetDashboardSummary(currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Unable to fetch dashboard summary")
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (dc *DashboardController) GetUserTasks(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	if limit > 20 {
		limit = 20
	}

	userTasks, err := dc.service.GetUserTasks(currentUserID, limit)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Unable to fetch user tasks")
		return
	}

	c.JSON(http.StatusOK, userTasks)

}

func (dc *DashboardController) GetOverdueTasks(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	overdueTasks, err := dc.service.GetOverdueTasks(currentUserID, limit)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Unable to fetch overdue tasks")
		return
	}

	c.JSON(http.StatusOK, overdueTasks)
}

func (dc *DashboardController) GetAssignedRisks(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	assignedRisks, err := dc.service.GetAssignedRisks(currentUserID, limit)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Unable to fetch assigned risks")
		return
	}

	c.JSON(http.StatusOK, assignedRisks)
}

func (dc *DashboardController) GetRecentActivity(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendUnauthorizedErrResp(c, "Unauthorized")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	recentActivity, err := dc.service.GetRecentActivity(currentUserID, limit)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, "Unable to fetch recent activity")
		return
	}

	c.JSON(http.StatusOK, recentActivity)
}
