package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/grc/threat-intelligence/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatIntelligencePlannedAdvisoryController struct {
	service services.ThreatIntelligencePlannedAdvisoryServiceInterface
	db      *gorm.DB
}

func NewThreatIntelligencePlannedAdvisoryController(service services.ThreatIntelligencePlannedAdvisoryServiceInterface) *ThreatIntelligencePlannedAdvisoryController {
	return &ThreatIntelligencePlannedAdvisoryController{
		service: service,
	}
}

func (pac *ThreatIntelligencePlannedAdvisoryController) GetAllThreatIntelligencePlannedAdvisories(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(pac.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"advisory_source_ids",
		"case_classification_ids",
		"advisory_type_ids",
		"importance_ids",
		"security_classification_ids",
		"advisory_status_ids",
		"internal_status_ids",
		"advisory_recipient_ids",
		"organization_unit_ids",
		"threat_severity_ids",
		"likelihood_ids",
		"risk_level_ids",
		"risk_level_value_ids",
		"processing_time",
	}

	validSortFields := map[string]string{
		"reference_number":                 "pa.reference_number",
		"case_number":                      "pa.case_number",
		"title_en":                         "pa.title_en",
		"title_ar":                         "pa.title_ar",
		"receipt_date":                     "pa.receipt_date",
		"case_link":                        "pa.case_link",
		"threat_type":                      "pa.threat_type",
		"processing_time":                  "pa.processing_time",
		"advisory_source_en":               "planned_advisory_sources.title_en",
		"advisory_source_ar":               "planned_advisory_sources.title_ar",
		"case_classification_en":           "planned_advisory_classifications.title_en",
		"case_classification_ar":           "planned_advisory_classifications.title_ar",
		"advisory_type_en":                 "planned_advisory_types.title_en",
		"advisory_type_ar":                 "planned_advisory_types.title_ar",
		"importance_en":                    "planned_advisory_importances.title_en",
		"importance_ar":                    "planned_advisory_importances.title_ar",
		"security_classification_en":       "planned_advisory_security_classifications.title_en",
		"security_classification_ar":       "planned_advisory_security_classifications.title_ar",
		"advisory_status_en":               "planned_advisory_statuses.title_en",
		"advisory_status_ar":               "planned_advisory_statuses.title_ar",
		"advisory_recipient_first_name_en": "advisory_recipient.first_name_en",
		"responsible_organization_unit_en": "responsible_unit.title_en",
		"responsible_organization_unit_ar": "responsible_unit.title_ar",
		"created_at":                       "pa.created_at",
		"updated_at":                       "pa.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	receiptDateFrom := c.Query("receipt_date_from")
	receiptDateTo := c.Query("receipt_date_to")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := tiDTOs.PlannedAdvisoryListQueryParamsDTO{
		Sort:            sort,
		SortBy:          dbField,
		Keyword:         keyword,
		Status:          status,
		Page:            page,
		PageSize:        pageSize,
		Filters:         filters,
		QuickFilter:     quickFilter,
		ReceiptDateFrom: receiptDateFrom,
		ReceiptDateTo:   receiptDateTo,
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		Export:          export,
	}

	advisories, err := pac.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(pac.db, c, err.Error())
			return
		}

		excelFile, err := pac.service.ExportPlannedAdvisory(advisories.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(pac.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Planned advisories.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, advisories)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) GetThreatIntelligencePlannedAdvisoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	advisory, err := pac.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, advisory)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) CreateThreatIntelligencePlannedAdvisory(c *gin.Context) {
	var formDTO tiDTOs.PlannedAdvisoryFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	id, err := pac.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pac *ThreatIntelligencePlannedAdvisoryController) UpdateThreatIntelligencePlannedAdvisory(c *gin.Context) {
	var formDTO tiDTOs.PlannedAdvisoryFormDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	err = pac.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) UpdateThreatIntelligencePlannedAdvisoryStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var statusUpdateDTO tiDTOs.PlannedAdvisoryStatusUpdateFormDTO
	if err := c.ShouldBindJSON(&statusUpdateDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	if err := pac.service.UpdateStatus(uint(id), currentUserID, &statusUpdateDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) DeleteThreatIntelligencePlannedAdvisories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	failedIDs, err := pac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pac *ThreatIntelligencePlannedAdvisoryController) GetThreatIntelligencePlannedAdvisoryDocumentPreview(c *gin.Context) {
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

	documentPath, err := pac.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
		}
		return
	}

	c.File(documentPath)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) DownloadThreatIntelligencePlannedAdvisoryDocument(c *gin.Context) {
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

	documentPath, err := pac.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(pac.db, c, e.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

func (pac *ThreatIntelligencePlannedAdvisoryController) ExportThreatIntelligencePlannedAdvisories(c *gin.Context) {
	c.Set("export", true)

	pac.GetAllThreatIntelligencePlannedAdvisories(c)
}
