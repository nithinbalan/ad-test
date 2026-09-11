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

type ThreatIntelligenceRiskLevelController struct {
	service services.ThreatIntelligenceRiskLevelServiceInterface
	db      *gorm.DB
}

func NewThreatIntelligenceRiskLevelController(service services.ThreatIntelligenceRiskLevelServiceInterface) *ThreatIntelligenceRiskLevelController {
	return &ThreatIntelligenceRiskLevelController{
		service: service,
	}
}

func (tirlc *ThreatIntelligenceRiskLevelController) GetAllThreatIntelligenceRiskLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	allowedParams := []string{
		"type",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	riskLevels, err := tirlc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskLevels)
}

func (tirlc *ThreatIntelligenceRiskLevelController) GetThreatIntelligenceRiskLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskLevel, err := tirlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, riskLevel)
}

func (tirlc *ThreatIntelligenceRiskLevelController) CreateThreatIntelligenceRiskLevel(c *gin.Context) {
	var riskLevel *tiDTOs.ThreatIntelligenceRiskLevelFormDTO

	if err := c.ShouldBindJSON(&riskLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	id, err := tirlc.service.Create(currentUserID, riskLevel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tirlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (tirlc *ThreatIntelligenceRiskLevelController) UpdateThreatIntelligenceRiskLevel(c *gin.Context) {
	riskLevelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedRiskLevel *tiDTOs.ThreatIntelligenceRiskLevelFormDTO
	if err := c.ShouldBindJSON(&updatedRiskLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	err = tirlc.service.Update(uint(riskLevelID), currentUserID, updatedRiskLevel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tirlc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tirlc *ThreatIntelligenceRiskLevelController) DeleteThreatIntelligenceRiskLevels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	failedIDs, err := tirlc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tirlc *ThreatIntelligenceRiskLevelController) ActivateThreatIntelligenceRiskLevel(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	err = tirlc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tirlc *ThreatIntelligenceRiskLevelController) DeactivateThreatIntelligenceRiskLevel(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	err = tirlc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tirlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
