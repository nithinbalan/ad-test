package controllers

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestPlanController struct {
	service services.PenetrationTestPlanServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestPlanController(service services.PenetrationTestPlanServiceInterface) *PenetrationTestPlanController {
	return &PenetrationTestPlanController{
		service: service,
	}
}

func (ptpc *PenetrationTestPlanController) GetAllPenetrationTestPlans(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	penetrationTestPlans, err := ptpc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, penetrationTestPlans)
}

func (ptpc *PenetrationTestPlanController) GetPenetrationTestPlanDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	penetrationTestPlan, err := ptpc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, penetrationTestPlan)
}

func (ptpc *PenetrationTestPlanController) CreatePenetrationTestPlan(c *gin.Context) {
	var penetrationTestPlan *models.PenetrationTestPlan

	if err := c.ShouldBindJSON(&penetrationTestPlan); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	created, err := ptpc.service.Create(currentUserID, penetrationTestPlan)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, created)
}

func (ptpc *PenetrationTestPlanController) UpdatePenetrationTestPlan(c *gin.Context) {
	var updated *models.PenetrationTestPlan

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updated); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	updated, err = ptpc.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptpc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updated)
}

func (ptpc *PenetrationTestPlanController) DeletePenetrationTestPlans(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	failedIDs, err := ptpc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ptpc *PenetrationTestPlanController) ActivatePenetrationTestPlans(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	if err := ptpc.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ptpc *PenetrationTestPlanController) DeactivatePenetrationTestPlans(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	if err := ptpc.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ptpc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
