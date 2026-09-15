package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestReportController struct {
	service services.PenetrationTestReportServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestReportController(service services.PenetrationTestReportServiceInterface) *PenetrationTestReportController {
	return &PenetrationTestReportController{
		service: service,
	}
}

func (ptrc *PenetrationTestReportController) GetAllPenetrationTestReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "penetration_test_reports.title",
		"start_date": "penetration_test_reports.start_date",
		"end_date":   "penetration_test_reports.end_date",
		"created_at": "penetration_test_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}

	if page < 1 || pageSize < 1 {
		utility.SendInvalidErrResp(c, errors.New("page and page_size must be greater than zero"))
		return
	}

	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := penetrationTestDTOs.PenetrationTestReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := ptrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (ptrc *PenetrationTestReportController) CreatePenetrationTestReport(c *gin.Context) {
	var request penetrationTestDTOs.PenetrationTestReportRequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	reportID, err := ptrc.service.Create(currentUserID, request)
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
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{
		ID: reportID,
	})
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestReportDetails(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := ptrc.service.GetPenetrationTestReportDetails(uint(penetrationTestReportID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

func (ptrc *PenetrationTestReportController) DeletePenetrationTestReports(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	err = ptrc.service.DeletePenetrationTestReport(uint(penetrationTestReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (ptrc *PenetrationTestReportController) GetPenetrationTestReportSection(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
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

	details, err := ptrc.service.GetSectionByID(uint(penetrationTestReportID), uint(sectionID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (ptrc *PenetrationTestReportController) UpdatePenetrationTestReportSection(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
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

	var formDTO penetrationTestDTOs.PenetrationTestReportSectionFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	err = ptrc.service.UpdateSection(uint(penetrationTestReportID), uint(sectionID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (ptrc *PenetrationTestReportController) DeletePenetrationTestReportSection(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
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
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	err = ptrc.service.DeletePenetrationTestReportSection(uint(penetrationTestReportID), uint(sectionID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptrc *PenetrationTestReportController) ExportPenetrationTestReportToPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	fileName, err := ptrc.service.ExportPenetrationTestReport(uint(id), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, fmt.Sprintf("Failed to export PenetrationTest Report PDF: %s", err.Error()))
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

func (ptrc *PenetrationTestReportController) GetPenetrationTestReportExportZip(c *gin.Context) {
	penetrationTestReportIDStr := c.Param("penetrationTestReportId")
	penetrationTestReportID, err := strconv.ParseUint(penetrationTestReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptrc.db, c, err.Error())
		return
	}

	zipFilePath, err := ptrc.service.GetPenetrationTestReportExportZip(uint(penetrationTestReportID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptrc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=penetration_test_report_%d.zip", penetrationTestReportID))
	c.File(zipFilePath)

	defer os.Remove(zipFilePath)
}
