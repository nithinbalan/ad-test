package controllers

import (
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MyDashboardController struct {
	service organizationServices.MyDashboardServiceInterface
	db      *gorm.DB
}

func NewMyDashboardController(service organizationServices.MyDashboardServiceInterface, db *gorm.DB) *MyDashboardController {
	return &MyDashboardController{
		service: service,
		db:      db,
	}
}

func (dc *MyDashboardController) GetDashboardSummary(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	summary, err := dc.service.GetDashboardSummary(currentUserID)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (dc *MyDashboardController) GetKPIs(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	kpis, err := dc.service.GetKPIs(currentUserID)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, kpis)
}

func (dc *MyDashboardController) GetUserTasks(c *gin.Context) {
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

	userTasks, err := dc.service.GetUserTasks(currentUserID, limit)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, userTasks)

}

func (dc *MyDashboardController) GetOverdueTasks(c *gin.Context) {
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
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, overdueTasks)
}

func (dc *MyDashboardController) GetAssignedRisks(c *gin.Context) {
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
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, assignedRisks)
}

func (dc *MyDashboardController) GetRecentActivity(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dc.db, c, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	recentActivity, err := dc.service.GetRecentActivity(currentUserID, limit)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, recentActivity)
}
