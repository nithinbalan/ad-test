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

type RiskMatrixImpactGuidelineController struct {
	service services.RiskMatrixImpactGuidelineServiceInterface
	db      *gorm.DB
}

func NewRiskMatrixImpactGuidelineController(service services.RiskMatrixImpactGuidelineServiceInterface) *RiskMatrixImpactGuidelineController {
	return &RiskMatrixImpactGuidelineController{
		service: service,
	}
}

func (rmigc *RiskMatrixImpactGuidelineController) GetAllRiskMatrixImpactGuidelines(c *gin.Context) {

	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_impact_category_ids",
		"risk_matrix_impact_id",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskMatrixImpactGuidelineListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	riskMatrixImpactGuidelines, err := rmigc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskMatrixImpactGuidelines)
}

func (rmigc *RiskMatrixImpactGuidelineController) GetRiskMatrixImpactGuidelineDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskMatrixImpactGuideline, err := rmigc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rmigc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskMatrixImpactGuideline)
}

func (rmigc *RiskMatrixImpactGuidelineController) CreateRiskMatrixImpactGuideline(c *gin.Context) {

	var riskMatrixImpactGuidelineFormDTO riskDTOs.RiskMatrixImpactGuidelineFormDTO
	if err := c.ShouldBindJSON(&riskMatrixImpactGuidelineFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, err.Error())
		return
	}
	
	ids , err := rmigc.service.Create(currentUserID, &riskMatrixImpactGuidelineFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmigc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.SuccessfulIDsResponse{SuccessfulIds: ids})
}

func (rmigc *RiskMatrixImpactGuidelineController) UpdateRiskMatrixImpactGuideline(c *gin.Context) {

	riskMatrixImpactGuidelineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var riskMatrixImpactGuidelineUpdateFormDTO riskDTOs.RiskMatrixImpactGuidelineUpdateFormDTO
	if err := c.ShouldBindJSON(&riskMatrixImpactGuidelineUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, err.Error())
		return
	}

	err = rmigc.service.Update(uint(riskMatrixImpactGuidelineID), currentUserID, &riskMatrixImpactGuidelineUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rmigc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rmigc *RiskMatrixImpactGuidelineController) DeleteRiskMatrixImpactGuidelines(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, err.Error())
		return
	}

	failedIDs, err := rmigc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rmigc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
