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

type RiskAppetiteStatementController struct {
	service services.RiskAppetiteStatementServiceInterface
	db      *gorm.DB
}

func NewRiskAppetiteStatementController(service services.RiskAppetiteStatementServiceInterface) *RiskAppetiteStatementController {
	return &RiskAppetiteStatementController{
		service: service,
	}
}

func (rasc *RiskAppetiteStatementController) GetAllRiskAppetiteStatements(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	// byCategory := c.Query("by_category")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"severity_level_ids",
		"potential_risk_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskAppetiteStatementListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	var riskAppetiteStatements interface{}
	riskAppetiteStatements, err = rasc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskAppetiteStatements)
}

func (rasc *RiskAppetiteStatementController) GetAllRiskAppetiteStatementByCategory(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	byCategory := c.Query("by_category")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"severity_level_ids",
		"potential_risk_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.RiskAppetiteStatementListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	var riskAppetiteStatements interface{}
	if byCategory != "" {
		riskAppetiteStatements, err = rasc.service.GetAllByCategories(&queryParams)
	}

	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, riskAppetiteStatements)
}

func (rasc *RiskAppetiteStatementController) GetRiskAppetiteStatementDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	riskAppetiteStatement, err := rasc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rasc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskAppetiteStatement)
}

func (rasc *RiskAppetiteStatementController) CreateRiskAppetiteStatement(c *gin.Context) {

	var riskAppetiteStatementForm riskDTOs.RiskAppetiteStatementFormDTO

	if err := c.ShouldBindJSON(&riskAppetiteStatementForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	id, err := rasc.service.Create(currentUserID, &riskAppetiteStatementForm)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rasc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rasc *RiskAppetiteStatementController) UpdateRiskAppetiteStatement(c *gin.Context) {

	riskAppetiteStatementID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var riskAppetiteStatementForm riskDTOs.RiskAppetiteStatementFormDTO
	if err := c.ShouldBindJSON(&riskAppetiteStatementForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	err = rasc.service.Update(uint(riskAppetiteStatementID), currentUserID, &riskAppetiteStatementForm)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rasc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rasc *RiskAppetiteStatementController) DeleteRiskAppetiteStatements(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	failedIDs, err := rasc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rasc *RiskAppetiteStatementController) ActivateRiskAppetiteStatements(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	err = rasc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rasc *RiskAppetiteStatementController) DeactivateRiskAppetiteStatements(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	err = rasc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rasc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
