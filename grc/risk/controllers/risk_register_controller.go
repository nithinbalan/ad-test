package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterController struct {
	service services.RiskRegisterServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterController(service services.RiskRegisterServiceInterface) *RiskRegisterController {
	return &RiskRegisterController{
		service: service,
	}
}

func (rrc *RiskRegisterController) GetAllRiskRegisters(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	valueStr := c.Query("value")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"risk_category_ids",
		"risk_source_ids",
		"responsible_unit_ids",
		"responsible_user_ids",
		"risk_status_ids",
		"review_frequency_ids",
		"inherent_risk_rating_ids",
		"residual_risk_rating_ids",
		"current_likelihood_ids",
		"current_impact_ids",
		"current_risk_rating_ids",
		"risk_response_strategy_ids",
	}

	validSortFields := map[string]string{
		"reference_number":          "risk_register.reference_number",
		"title":                     "risk_register.title",
		"description":               "risk_register.description",
		"risk_category_en":          "risk_categories.title_en",
		"risk_category_ar":          "risk_categories.title_ar",
		"organization_unit_en":      "organization_units.title_en",
		"organization_unit_ar":      "organization_units.title_ar",
		"responsible_user_en":       "users.first_name_en",
		"responsible_user_ar":       "users.first_name_ar",
		"risk_source_en":            "risk_sources.title_en",
		"risk_source_ar":            "risk_sources.title_ar",
		"risk_status_en":            "risk_statuses.title_en",
		"risk_status_ar":            "risk_statuses.title_ar",
		"inherent_likelihood_en":    "inherent_likelihood.title_en",
		"inherent_likelihood_ar":    "inherent_likelihood.title_ar",
		"inherent_impact_en":        "inherent_impact.title_en",
		"inherent_impact_ar":        "inherent_impact.title_ar",
		"inherent_risk_rating_en":   "inherent_risk_rating.title_en",
		"inherent_risk_rating_ar":   "inherent_risk_rating.title_ar",
		"risk_response_strategy_en": "risk_response_strategies.title_en",
		"risk_response_strategy_ar": "risk_response_strategies.title_ar",
		"residual_likelihood_en":    "residual_likelihood.title_en",
		"residual_likelihood_ar":    "residual_likelihood.title_ar",
		"residual_impact_en":        "residual_impact.title_en",
		"residual_impact_ar":        "residual_impact.title_ar",
		"residual_risk_rating_en":   "residual_risk_rating.title_en",
		"residual_risk_rating_ar":   "residual_risk_rating.title_ar",
		"reported_at":               "risk_register.reported_at",
		"review_frequency_en":       "review_frequencies.title_en",
		"review_frequency_ar":       "review_frequencies.title_ar",
		"next_review_date":          "risk_register.next_review_date",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	reportedAtFrom := c.Query("reported_at")
	reportedAtTo := c.Query("reported_at")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, "Invalid page_size parameter")
		return
	}

	var value int
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, "Invalid value parameter")
			return
		}
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskRegisterListQueryParamsDTO{
		Sort:               sort,
		SortBy:             dbField,
		Keyword:            keyword,
		Status:             status,
		Page:               page,
		PageSize:           pageSize,
		Value:              value,
		Filters:            filters,
		QuickFilter:        quickFilter,
		ReportedAtFrom:     reportedAtFrom,
		ReportedAtTo:       reportedAtTo,
		NextReviewDateFrom: nextReviewDateFrom,
		NextReviewDateTo:   nextReviewDateTo,
		Export:             export,
		DateFrom:           dateFrom,
		DateTo:             dateTo,
	}
	risks, err := rrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}

		excelFile, err := rrc.service.ExportRiskRegister(risks.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(rrc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Risk register.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, risks)
}

func (rrc *RiskRegisterController) GetRiskRegisterDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	risk, err := rrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, risk)
}

func (rrc *RiskRegisterController) CreateRiskRegister(c *gin.Context) {
	var riskRegisterFormDTO riskDTOs.RiskRegisterFormDTO

	if err := c.ShouldBindJSON(&riskRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	id, err := rrc.service.Create(currentUserID, &riskRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rrc *RiskRegisterController) UpdateRiskRegister(c *gin.Context) {
	var riskRegisterFormDTO riskDTOs.RiskRegisterFormDTO

	riskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskRegisterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Update(uint(riskID), currentUserID, &riskRegisterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *RiskRegisterController) DeleteRiskRegisters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	failedIDs, err := rrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rrc *RiskRegisterController) UpdateResponseStrategy(c *gin.Context) {
	var riskResponseStrategyFormDTO riskDTOs.RiskResponseStrategyFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskResponseStrategyFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.UpdateResponseStrategy(uint(riskID), currentUserID, &riskResponseStrategyFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *RiskRegisterController) DeleteResponseStrategy(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.DeleteResponseStrategy(uint(riskID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrc *RiskRegisterController) GetControlPerformance(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlPerformance, err := rrc.service.GetControlPerformance(uint(riskID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, controlPerformance)
}

func (rrc *RiskRegisterController) GetRiskRegisterDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (rrc *RiskRegisterController) DownloadRiskRegisterDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}

func (rrc *RiskRegisterController) GetRiskRegisterExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	rrc.GetAllRiskRegisters(c)
}
