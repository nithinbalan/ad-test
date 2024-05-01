package controllers

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ChecklistController struct {
	service services.ChecklistServiceInterface
	db      *gorm.DB
}

func NewChecklistController(service services.ChecklistServiceInterface) *ChecklistController {
	return &ChecklistController{
		service: service,
	}
}

func (cc *ChecklistController) GetAllChecklists(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	checklists, err := cc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, checklists)
}

func (cc *ChecklistController) GetChecklistDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	checklist, err := cc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, checklist)
}

func (cc *ChecklistController) CreateChecklist(c *gin.Context) {

	var checklist *models.Checklist

	if err := c.ShouldBindJSON(&checklist); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	createdChecklist, err := cc.service.Create(currentUserID, checklist)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdChecklist)
}

func (cc *ChecklistController) UpdateChecklist(c *gin.Context) {

	checklistID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedChecklist *models.Checklist
	if err := c.ShouldBindJSON(&updatedChecklist); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	updatedChecklist, err = cc.service.Update(uint(checklistID), currentUserID, updatedChecklist)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedChecklist)
}

func (cc *ChecklistController) DeleteChecklists(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	failedIDs, err := cc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cc *ChecklistController) ActivateChecklists(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *ChecklistController) DeactivateChecklists(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
