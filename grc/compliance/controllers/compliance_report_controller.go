package controllers

import (
	"errors"
	"fmt"
	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceReportController struct {
	service services.ComplianceReportServiceInterface
	db      *gorm.DB
}

func NewComplianceReportController(service services.ComplianceReportServiceInterface) *ComplianceReportController {
	return &ComplianceReportController{
		service: service,
	}
}

func (crc *ComplianceReportController) GetAllComplianceReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "compliance_reports.title",
		"start_date": "compliance_reports.start_date",
		"end_date":   "compliance_reports.end_date",
		"created_at": "compliance_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := dtos.ComplianceReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := crc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (crc *ComplianceReportController) CreateComplianceReport(c *gin.Context) {
	var request dtos.ComplianceReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	complianceReportID, err := crc.service.Create(currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Compliance report created successfully",
		"id":      complianceReportID,
	})
}

func (crc *ComplianceReportController) GetComplianceReportDetails(c *gin.Context) {
	complianceReportID, err := strconv.Atoi(c.Param("complianceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceReport, err := crc.service.GetComplianceReportDetails(uint(complianceReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceReport)
}

func (crc *ComplianceReportController) UpdateComplianceReportSection(c *gin.Context) {
	complianceReportID, err := strconv.Atoi(c.Param("complianceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.ComplianceReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.UpdateSection(uint(complianceReportID), uint(sectionID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crc *ComplianceReportController) DeleteComplianceReport(c *gin.Context) {
	complianceReportID, err := strconv.Atoi(c.Param("complianceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.DeleteComplianceReport(uint(complianceReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (crc *ComplianceReportController) ExportComplianceReportToPDF(c *gin.Context) {
	complianceReportID, err := strconv.Atoi(c.Param("complianceReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	// Delegate to service to handle PDF generation
	fileName, err := crc.service.ExportComplianceReportToPDF(uint(complianceReportID), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(crc.db, c, fmt.Sprintf("Failed to export PDF: %s", err.Error()))
			return
		}
	}

	// Set headers to ensure the file is served with the correct file name
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	// Serve the PDF file
	c.File(fileName)
	// Cleanup after the response has been sent
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second) // Add a small delay to ensure file is no longer in use
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}
