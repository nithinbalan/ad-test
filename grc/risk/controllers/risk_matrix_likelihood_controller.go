package controllers

import (
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskMatrixLikelihoodController struct {
	service services.RiskMatrixLikelihoodServiceInterface
	db      *gorm.DB
}

func NewRiskMatrixLikelihoodController(service services.RiskMatrixLikelihoodServiceInterface) *RiskMatrixLikelihoodController {
	return &RiskMatrixLikelihoodController{
		service: service,
	}
}

func (rmlc *RiskMatrixLikelihoodController) GetAllRiskMatrixLikelihoods(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	likelihoods, err := rmlc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, likelihoods)
}

func (rmlc *RiskMatrixLikelihoodController) GetRiskMatrixLikelihoodDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskLikelihood, err := rmlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskLikelihood)
}

func (rmlc *RiskMatrixLikelihoodController) CreateRiskMatrixLikelihood(c *gin.Context) {

	var likelihood *riskDTOs.RiskMatrixLikelihoodFormDTO

	if err := c.ShouldBindJSON(&likelihood); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	createdLikelihood, err := rmlc.service.Create(currentUserID, likelihood)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdLikelihood)
}

func (rmlc *RiskMatrixLikelihoodController) UpdateRiskMatrixLikelihood(c *gin.Context) {

	likelihoodID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedLikelihood *riskDTOs.RiskMatrixLikelihoodFormDTO
	if err := c.ShouldBindJSON(&updatedLikelihood); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	err = rmlc.service.Update(uint(likelihoodID), currentUserID, updatedLikelihood)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmlc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rmlc *RiskMatrixLikelihoodController) DeleteRiskMatrixLikelihoods(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	failedIDs, err := rmlc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rmlc *RiskMatrixLikelihoodController) ActivateRiskMatrixLikelihood(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	err = rmlc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rmlc *RiskMatrixLikelihoodController) DeactivateRiskMatrixLikelihood(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	err = rmlc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
