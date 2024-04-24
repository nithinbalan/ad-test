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

type RiskMatrixImpactController struct {
	service services.RiskMatrixImpactServiceInterface
	db      *gorm.DB
}

func NewRiskMatrixImpactController(service services.RiskMatrixImpactServiceInterface) *RiskMatrixImpactController {
	return &RiskMatrixImpactController{
		service: service,
	}
}

func (rmic *RiskMatrixImpactController) GetAllRiskMatrixImpacts(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	impacts, err := rmic.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, impacts)
}

func (rmic *RiskMatrixImpactController) GetRiskMatrixImpactDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskImpact, err := rmic.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskImpact)
}

func (rmic *RiskMatrixImpactController) CreateRiskMatrixImpact(c *gin.Context) {

	var impact *riskDTOs.RiskMatrixImpactFormDTO

	if err := c.ShouldBindJSON(&impact); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	id, err := rmic.service.Create(currentUserID, impact)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (rmic *RiskMatrixImpactController) UpdateRiskMatrixImpact(c *gin.Context) {

	impactID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedImpact *riskDTOs.RiskMatrixImpactFormDTO
	if err := c.ShouldBindJSON(&updatedImpact); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	err = rmic.service.Update(uint(impactID), currentUserID, updatedImpact)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rmic *RiskMatrixImpactController) DeleteRiskMatrixImpacts(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	failedIDs, err := rmic.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rmic *RiskMatrixImpactController) ActivateRiskMatrixImpact(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	err = rmic.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rmic *RiskMatrixImpactController) DeactivateRiskMatrixImpact(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	err = rmic.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmic.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
