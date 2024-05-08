package controllers

import (
	"jethur/grc/models"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskSourceController struct {
	service services.RiskSourceServiceInterface
	db      *gorm.DB
}

func NewRiskSourceController(service services.RiskSourceServiceInterface) *RiskSourceController {
	return &RiskSourceController{
		service: service,
	}
}

func (rsc *RiskSourceController) GetAllRiskSources(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	riskSources, err := rsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskSources)
}

func (rsc *RiskSourceController) GetRiskSourceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskSource, err := rsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskSource)
}

func (rsc *RiskSourceController) CreateRiskSource(c *gin.Context) {

	var source *models.RiskSource

	if err := c.ShouldBindJSON(&source); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	createdSource, err := rsc.service.Create(currentUserID, source)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSource)
}

func (rsc *RiskSourceController) UpdateRiskSource(c *gin.Context) {

	sourceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedSource *models.RiskSource
	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	updatedSource, err = rsc.service.Update(uint(sourceID), currentUserID, updatedSource)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSource)
}

func (rsc *RiskSourceController) DeleteRiskSources(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	failedIDs, err := rsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rsc *RiskSourceController) ActivateRiskSources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	err = rsc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rsc *RiskSourceController) DeactivateRiskSources(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	err = rsc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
