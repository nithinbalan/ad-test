package controllers

import (
	"errors"
	"fmt"
	"net/http"

	// "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentReportController struct {
	service services.IncidentReportServiceInterface
	db      *gorm.DB
}

func NewIncidentReportController(service services.IncidentReportServiceInterface) *IncidentReportController {
	return &IncidentReportController{
		service: service,
	}
}

func (irc *IncidentReportController) GetAllIncidentReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "incident_reports.title",
		"start_date": "incident_reports.start_date",
		"end_date":   "incident_reports.end_date",
		"created_at": "incident_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := incidentDTOs.IncidentReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := irc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (irc *IncidentReportController) ExportIncidentReportToPDF(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")
	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	// Delegate to service to handle PDF generation
	fileName, err := irc.service.ExportIncidentReport(uint(incidentID), uint(currentUserID))
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
			utility.SendInternalErrResp(irc.db, c, fmt.Sprintf("Failed to export Incident Report PDF: %s", err.Error()))
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

func (irc *IncidentReportController) CreateIncidentReport(c *gin.Context) {
	var request incidentDTOs.IncidentReportRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	reportID, err := irc.service.Create(currentUserID, request)
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
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: reportID,
	})
}

func (irc *IncidentReportController) GetIncidentReportDetails(c *gin.Context) {
	incidentReportIDStr := c.Param("incidentReportId")
	incidentReportID, err := strconv.ParseUint(incidentReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := irc.service.GetIncidentReportDetails(uint(incidentReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

func (irc *IncidentReportController) DeleteIncidentReports(c *gin.Context) {
	incidentReportID, err := strconv.Atoi(c.Param("incidentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	err = irc.service.DeleteIncidentReport(uint(incidentReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (irc *IncidentReportController) GetIncidentReportSection(c *gin.Context) {
	incidentReportID, err := strconv.Atoi(c.Param("incidentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := irc.service.GetSectionByID(uint(incidentReportID), uint(sectionID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (irc *IncidentReportController) UpdateIncidentReportSection(c *gin.Context) {
	incidentReportID, err := strconv.Atoi(c.Param("incidentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO incidentDTOs.IncidentReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	err = irc.service.UpdateSection(uint(incidentReportID), uint(sectionID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}
	// Dummy response payload for testing the API contract
	c.JSON(http.StatusOK, gin.H{
		"id":                 sectionID,
		"incident_report_id": incidentReportID,
		"type":               "details",
		"title_en":           formDTO.TitleEn,
		"title_ar":           formDTO.TitleAr,
		"description_en":     formDTO.DescriptionEn,
		"description_ar":     formDTO.DescriptionAr,
		"order_num":          1,
	})
}

func (irc *IncidentReportController) DeleteIncidentReportSection(c *gin.Context) {
	incidentReportID, err := strconv.Atoi(c.Param("incidentReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	err = irc.service.DeleteIncidentReportSection(uint(incidentReportID), uint(sectionID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, e.Error())
			return
		}
	}

	// Dummy response payload for testing the API contract
	c.JSON(http.StatusOK, gin.H{
		"message":            "Incident report section deleted successfully",
		"incident_report_id": incidentReportID,
		"section_id":         sectionID,
	})
}

func (irc *IncidentReportController) GetIncidentReportExportZip(c *gin.Context) {
	incidentReportIDStr := c.Param("incidentReportId")
	incidentReportID, err := strconv.ParseUint(incidentReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	// Call the service to export and get ZIP file path
	zipFilePath, err := irc.service.GetIncidentReportExportZip(uint(incidentReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(irc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}
	// Set headers and return the file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=incident_report_%d.zip", incidentReportID))
	c.File(zipFilePath)

	// Optionally delete the file after serving it
	defer os.Remove(zipFilePath)
}
