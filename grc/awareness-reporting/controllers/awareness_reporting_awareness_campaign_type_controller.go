package controllers

import (
	"jethur/grc/awareness-reporting/dtos"
	"jethur/grc/awareness-reporting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessCampaignTypeController struct {
	service services.AwarenessCampaignTypeServiceInterface
	db      *gorm.DB
}

func NewAwarenessCampaignTypeController(service services.AwarenessCampaignTypeServiceInterface) *AwarenessCampaignTypeController {
	return &AwarenessCampaignTypeController{
		service: service,
	}
}

func (actc *AwarenessCampaignTypeController) GetAllAwarenessCampaignTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	awarenessCampaignTypes, err := actc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, awarenessCampaignTypes)
}

func (actc *AwarenessCampaignTypeController) GetAwarenessCampaignTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	awarenessCampaignType, err := actc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, awarenessCampaignType)
}

func (actc *AwarenessCampaignTypeController) CreateAwarenessCampaignType(c *gin.Context) {
	var form dtos.AwarenessCampaignTypeFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	createdType, err := actc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(actc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (actc *AwarenessCampaignTypeController) UpdateAwarenessCampaignType(c *gin.Context) {
	var form dtos.AwarenessCampaignTypeFormDTO

	awarenessCampaignTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	updatedType, err := actc.service.Update(uint(awarenessCampaignTypeID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(actc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (actc *AwarenessCampaignTypeController) DeleteAwarenessCampaignTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	failedIDs, err := actc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (actc *AwarenessCampaignTypeController) ActivateAwarenessCampaignTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	err = actc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (actc *AwarenessCampaignTypeController) DeactivateAwarenessCampaignTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	err = actc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(actc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}