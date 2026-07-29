package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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
	return &DocumentReportController{
		service: service,
	}
}

func (drc *DocumentReportController) GetAllDocumentReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"title":             "document_reports.title",
		"created_at":        "document_reports.created_at",
		"document_category": "document_categories.title_en",
		"document_type":     "document_types.title_en",
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

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.DocumentReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	reports, err := drc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (drc *DocumentReportController) GetDocumentReportDetails(c *gin.Context) {
	reportID, err := strconv.Atoi(c.Param("documentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := drc.service.GetDocumentReportDetails(uint(reportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (drc *DocumentReportController) CreateDocumentReport(c *gin.Context) {
	var request dtos.DocumentReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	reportID, err := drc.service.Create(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Document report generated successfully", "document_report_id": reportID,
	})
}

func (drc *DocumentReportController) DeleteDocumentReport(c *gin.Context) {
	reportID, err := strconv.Atoi(c.Param("documentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	err = drc.service.DeleteDocumentReport(uint(reportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (drc *DocumentReportController) GetDocumentReportExportZip(c *gin.Context) {
	reportIDString := c.Param("documentReportId")
	reportID, err := strconv.ParseUint(reportIDString, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	zipFilePath, err := drc.service.GetDocumentReportExportZip(uint(reportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(drc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=document_report_%d.zip", reportID))
	c.File(zipFilePath)

	defer os.Remove(zipFilePath)
}
