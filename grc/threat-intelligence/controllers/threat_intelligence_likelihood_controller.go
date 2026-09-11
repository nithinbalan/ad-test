package controllers

import (
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatIntelligenceLikelihoodController struct {
	service services.ThreatIntelligenceLikelihoodServiceInterface
	db      *gorm.DB
}

func NewThreatIntelligenceLikelihoodController(service services.ThreatIntelligenceLikelihoodServiceInterface) *ThreatIntelligenceLikelihoodController {
	return &ThreatIntelligenceLikelihoodController{
		service: service,
	}
}

func (tilc *ThreatIntelligenceLikelihoodController) GetAllThreatIntelligenceLikelihoods(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	likelihoods, err := tilc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, likelihoods)
}

func (tilc *ThreatIntelligenceLikelihoodController) GetThreatIntelligenceLikelihoodDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	likelihood, err := tilc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, likelihood)
}

func (tilc *ThreatIntelligenceLikelihoodController) CreateThreatIntelligenceLikelihood(c *gin.Context) {
	var likelihood *tiDTOs.ThreatIntelligenceLikelihoodFormDTO

	if err := c.ShouldBindJSON(&likelihood); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	id, err := tilc.service.Create(currentUserID, likelihood)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tilc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (tilc *ThreatIntelligenceLikelihoodController) UpdateThreatIntelligenceLikelihood(c *gin.Context) {
	likelihoodID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedLikelihood *tiDTOs.ThreatIntelligenceLikelihoodFormDTO
	if err := c.ShouldBindJSON(&updatedLikelihood); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	err = tilc.service.Update(uint(likelihoodID), currentUserID, updatedLikelihood)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tilc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tilc *ThreatIntelligenceLikelihoodController) DeleteThreatIntelligenceLikelihoods(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	failedIDs, err := tilc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tilc *ThreatIntelligenceLikelihoodController) ActivateThreatIntelligenceLikelihood(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	err = tilc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tilc *ThreatIntelligenceLikelihoodController) DeactivateThreatIntelligenceLikelihood(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	err = tilc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tilc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
