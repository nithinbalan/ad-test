package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PotentialRiskController struct {
	service services.PotentialRiskServiceInterface
	db      *gorm.DB
}

func NewPotentialRiskController(service services.PotentialRiskServiceInterface) *PotentialRiskController {
	return &PotentialRiskController{
		service: service,
	}
}

func (prc *PotentialRiskController) GetAllPotentialRisks(c *gin.Context) {

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"severity_level_ids",
		"potential_risk_status_ids",
	}

	validSortFields := map[string]string{
		"title":             "potential_risks.title",
		"description":       "potential_risks.description",
		"risk_category_en":  "risk_categories.title_en",
		"risk_category_ar":  "risk_categories.title_ar",
		"severity_level_en": "severity_levels.title_en",
		"severity_level_ar": "severity_levels.title_ar",
		"status_en":         "potential_risk_statuses.title_en",
		"status_ar":         "potential_risk_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.PotentialRiskListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}
	potentialRisks, err := prc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}

		excelFile, err := prc.service.ExportPotentialRisk(potentialRisks.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Potential risk.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, potentialRisks)
}

func (prc *PotentialRiskController) GetPotentialRiskDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	potentialRisk, err := prc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, potentialRisk)
}

func (prc *PotentialRiskController) CreatePotentialRisk(c *gin.Context) {

	var potentialRiskFormDTO riskDTOs.PotentialRiskFormDTO

	if err := c.ShouldBindJSON(&potentialRiskFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	id, err := prc.service.Create(currentUserID, &potentialRiskFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (prc *PotentialRiskController) UpdatePotentialRisk(c *gin.Context) {

	potentialRiskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var potentialRiskFormDTO riskDTOs.PotentialRiskFormDTO
	if err := c.ShouldBindJSON(&potentialRiskFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.Update(uint(potentialRiskID), currentUserID, &potentialRiskFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (prc *PotentialRiskController) DeletePotentialRisks(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	failedIDs, err := prc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
