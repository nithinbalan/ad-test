package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectReportController struct {
	service services.ProjectReportServiceInterface
	db      *gorm.DB
}

func NewProjectReportController(service services.ProjectReportServiceInterface) *ProjectReportController {
	return &ProjectReportController{
		service: service,
	}
}

func (prc *ProjectReportController) GetAllProjectReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "project_reports.title",
		"start_date": "project_reports.start_date",
		"end_date":   "project_reports.end_date",
		"created_at": "project_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

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

	queryParams := projectDTOs.ProjectReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := prc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (prc *ProjectReportController) ExportProjectReportToPDF(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	if projectIDStr == "" {
		projectIDStr = c.Param("id")
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	fileName, err := prc.service.ExportProjectReport(uint(projectID), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(prc.db, c, fmt.Sprintf("Failed to export Project Report PDF: %s", err.Error()))
			return
		}
	}

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

func (prc *ProjectReportController) CreateProjectReport(c *gin.Context) {
	var request projectDTOs.ProjectReportRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	reportID, err := prc.service.Create(currentUserID, request)
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
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: reportID,
	})
}

func (prc *ProjectReportController) GetProjectReportDetails(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := prc.service.GetProjectReportDetails(uint(projectReportID))
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

	c.JSON(http.StatusOK, details)
}

func (prc *ProjectReportController) DeleteProjectReports(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.DeleteProjectReport(uint(projectReportID), currentUserID)
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

func (prc *ProjectReportController) GetProjectReportSection(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionIDStr := c.Param("sectionId")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := prc.service.GetProjectReportSection(uint(projectReportID), uint(sectionID))
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
	c.JSON(http.StatusOK, details)
}

func (prc *ProjectReportController) UpdateProjectReportSection(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionIDStr := c.Param("sectionId")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO projectDTOs.ProjectReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.UpdateProjectReportSection(uint(projectReportID), uint(sectionID), formDTO, currentUserID)
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

func (prc *ProjectReportController) DeleteProjectReportSection(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionIDStr := c.Param("sectionId")
	sectionID, err := strconv.ParseUint(sectionIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.DeleteProjectReportSection(uint(projectReportID), uint(sectionID), currentUserID)
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

func (prc *ProjectReportController) GetProjectReportExportZip(c *gin.Context) {
	projectReportIDStr := c.Param("projectReportId")
	projectReportID, err := strconv.ParseUint(projectReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	zipFilePath, err := prc.service.GetProjectReportExportZip(uint(projectReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=project_report_%d.zip", projectReportID))
	c.File(zipFilePath)

	defer os.Remove(zipFilePath)
}
