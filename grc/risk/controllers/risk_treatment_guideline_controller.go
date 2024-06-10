package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskTreatmentGuidelineController struct {
	service services.RiskTreatmentGuidelineServiceInterface
	db      *gorm.DB
}

func NewRiskTreatmentGuidelineController(service services.RiskTreatmentGuidelineServiceInterface) *RiskTreatmentGuidelineController {
	return &RiskTreatmentGuidelineController{
		service: service,
	}
}

func (rtgc *RiskTreatmentGuidelineController) GetAllRiskTreatmentGuidelines(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_level_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskTreatmentGuidelineListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}
	guidelines, err := rtgc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, guidelines)
}

func (rtgc *RiskTreatmentGuidelineController) GetRiskTreatmentGuidelineDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	guideline, err := rtgc.service.GetByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rtgc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, guideline)
}

func (rtgc *RiskTreatmentGuidelineController) CreateRiskTreatmentGuideline(c *gin.Context) {
	var guidelineDTO riskDTOs.RiskTreatmentGuidelineFormDTO
	if err := c.ShouldBindJSON(&guidelineDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	id, err := rtgc.service.Create(currentUserID, &guidelineDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rtgc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rtgc *RiskTreatmentGuidelineController) UpdateRiskTreatmentGuideline(c *gin.Context) {
	guidelineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var guidelineDTO riskDTOs.RiskTreatmentGuidelineFormDTO
	if err := c.ShouldBindJSON(&guidelineDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	err = rtgc.service.Update(uint(guidelineID), currentUserID, &guidelineDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rtgc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rtgc *RiskTreatmentGuidelineController) DeleteRiskTreatmentGuidelines(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	failedIDs, err := rtgc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rtgc *RiskTreatmentGuidelineController) ActivateRiskTreatmentGuidelines(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	err = rtgc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rtgc *RiskTreatmentGuidelineController) DeactivateRiskTreatmentGuidelines(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	err = rtgc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rtgc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
