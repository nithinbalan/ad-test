package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentReportController struct {
	service services.DocumentReportServiceInterface
	db      *gorm.DB
}

func NewDocumentReportController(service services.DocumentReportServiceInterface) *DocumentReportController {
	return &DocumentReportController{service: service}
}

func (controller *DocumentReportController) GetAllDocumentReports(c *gin.Context) {
	sortOrder := c.DefaultQuery("sort_order", "desc")
	if sortOrder != "asc" && sortOrder != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}
	validSortFields := map[string]string{
		"title": "document_reports.title", "created_at": "document_reports.created_at",
		"document_category": "document_categories.title_en", "document_type": "document_types.title_en",
	}
	sortBy, ok := validSortFields[c.DefaultQuery("sort_by", "created_at")]
	if !ok {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}
	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	result, err := controller.service.GetAll(&dtos.DocumentReportListQueryParamsDTO{
		Sort: sortOrder, SortBy: sortBy, Keyword: c.Query("keyword"), Page: page, PageSize: pageSize,
	})
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (controller *DocumentReportController) GetDocumentReportDetails(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("documentReportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	result, err := controller.service.GetDocumentReportDetails(uint(reportID))
	if err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			utility.SendNotFoundErrResp(c)
			return
		}
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (controller *DocumentReportController) CreateDocumentReport(c *gin.Context) {
	var form dtos.DocumentReportFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	reportID, err := controller.service.CreateDocumentReport(currentUserID, &form)
	if err != nil {
		switch typedErr := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, typedErr)
		default:
			utility.SendInternalErrResp(controller.db, c, err.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Document report generated successfully", "document_report_id": reportID,
	})
}

func (controller *DocumentReportController) DeleteDocumentReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("documentReportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	if err := controller.service.DeleteDocumentReport(uint(reportID), currentUserID); err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			utility.SendNotFoundErrResp(c)
			return
		}
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
