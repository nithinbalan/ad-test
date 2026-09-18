package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateReportController struct {
	service services.ProjectSecurityGateReportServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateReportController(service services.ProjectSecurityGateReportServiceInterface) *ProjectSecurityGateReportController {
	return &ProjectSecurityGateReportController{
		service: service,
	}
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateReportExportZip(c *gin.Context) {
	securityGateIDStr := c.Param("id")
	securityGateID, err := strconv.ParseUint(securityGateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	language := c.Query("language")
	if language != "ar" {
		language = "en"
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	zipFilePath, err := psgrc.service.GetProjectSecurityGateReportExportZip(uint(securityGateID), uint(currentUserID), language)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(zipFilePath)))
	c.File(zipFilePath)

	defer os.Remove(zipFilePath)
}

func (psgrc *ProjectSecurityGateReportController) CreateProjectSecurityGateExceptionReport(c *gin.Context) {
	var form projectDTOs.ProjectSecurityGateExceptionReportFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	reportID, err := psgrc.service.CreateSecurityGateExceptionReport(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: reportID,
	})
}

func (psgrc *ProjectSecurityGateReportController) GetAllProjectSecurityGateExceptionReports(c *gin.Context) {
	queryParams, err := psgrc.parseReportListQueryParams(c, "project_security_gate_exception_reports")
	if err != nil {
		return
	}

	reports, err := psgrc.service.GetAllSecurityGateExceptionReports(queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateExceptionReportDetails(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := psgrc.service.GetSecurityGateExceptionReportDetails(uint(reportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

func (psgrc *ProjectSecurityGateReportController) ExportProjectSecurityGateExceptionReportToPDF(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	language := c.Query("language")
	if language != "ar" {
		language = "en"
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	fileName, err := psgrc.service.ExportSecurityGateExceptionReport(uint(reportID), uint(currentUserID), language)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, fmt.Sprintf("Failed to export Project Security Gate Exception Report PDF: %s", err.Error()))
			return
		}
	}

	psgrc.sendPDF(c, fileName)
}

func (psgrc *ProjectSecurityGateReportController) DeleteProjectSecurityGateExceptionReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.DeleteSecurityGateExceptionReport(uint(reportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgrc *ProjectSecurityGateReportController) CreateProjectSecurityGateReport(c *gin.Context) {
	var form projectDTOs.ProjectSecurityGateReportFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	reportID, err := psgrc.service.CreateSecurityGateReport(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: reportID,
	})
}

func (psgrc *ProjectSecurityGateReportController) GetAllProjectSecurityGateReports(c *gin.Context) {
	queryParams, err := psgrc.parseReportListQueryParams(c, "project_security_gate_reports")
	if err != nil {
		return
	}

	reports, err := psgrc.service.GetAllSecurityGateReports(queryParams)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateReportDetails(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := psgrc.service.GetSecurityGateReportDetails(uint(reportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

func (psgrc *ProjectSecurityGateReportController) ExportProjectSecurityGateReportToPDF(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	language := c.Query("language")
	if language != "ar" {
		language = "en"
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	fileName, err := psgrc.service.ExportSecurityGateReport(uint(reportID), uint(currentUserID), language)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, fmt.Sprintf("Failed to export Project Security Gate Report PDF: %s", err.Error()))
			return
		}
	}

	psgrc.sendPDF(c, fileName)
}

func (psgrc *ProjectSecurityGateReportController) DeleteProjectSecurityGateReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.DeleteSecurityGateReport(uint(reportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Both report lists share the same columns
func (psgrc *ProjectSecurityGateReportController) parseReportListQueryParams(c *gin.Context, tableName string) (*projectDTOs.ProjectSecurityGateReportListQueryParamsDTO, error) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	if sort != "asc" && sort != "desc" {
		err := errors.New("invalid 'sort_order' parameter")
		utility.SendInvalidErrResp(c, err)
		return nil, err
	}

	validSortFields := map[string]string{
		"title":                      tableName + ".title",
		"created_at":                 tableName + ".created_at",
		"created_user_first_name_en": "created_user.first_name_en",
		"created_user_first_name_ar": "created_user.first_name_ar",
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return nil, err
	}

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, "Invalid page parameter")
		return nil, err
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, "Invalid page_size parameter")
		return nil, err
	}

	allowedParams := []string{
		"created_by_ids",
	}

	return &projectDTOs.ProjectSecurityGateReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Filters:  utility.ParseQueryFilters(allowedParams, c.Request.URL.Query()),
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (psgrc *ProjectSecurityGateReportController) sendPDF(c *gin.Context, fileName string) {
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	c.File(fileName)
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second)
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	section, err := psgrc.service.GetSecurityGateReportSection(reportID, sectionID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, section)
}

func (psgrc *ProjectSecurityGateReportController) UpdateProjectSecurityGateReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	var form projectDTOs.ProjectSecurityGateReportSectionFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.UpdateSecurityGateReportSection(reportID, sectionID, &form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgrc *ProjectSecurityGateReportController) DeleteProjectSecurityGateReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.DeleteSecurityGateReportSection(reportID, sectionID, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateExceptionReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	section, err := psgrc.service.GetSecurityGateExceptionReportSection(reportID, sectionID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, section)
}

func (psgrc *ProjectSecurityGateReportController) UpdateProjectSecurityGateExceptionReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	var form projectDTOs.ProjectSecurityGateReportSectionFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.UpdateSecurityGateExceptionReportSection(reportID, sectionID, &form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psgrc *ProjectSecurityGateReportController) DeleteProjectSecurityGateExceptionReportSection(c *gin.Context) {
	reportID, sectionID, ok := psgrc.parseReportSectionParams(c)
	if !ok {
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	err = psgrc.service.DeleteSecurityGateExceptionReportSection(reportID, sectionID, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Both report types address sections the same way
func (psgrc *ProjectSecurityGateReportController) parseReportSectionParams(c *gin.Context) (uint, uint, bool) {
	reportID, err := strconv.ParseUint(c.Param("reportId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return 0, 0, false
	}

	sectionID, err := strconv.ParseUint(c.Param("sectionId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return 0, 0, false
	}

	return uint(reportID), uint(sectionID), true
}
