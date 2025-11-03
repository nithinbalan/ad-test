package controllers

import (
	"errors"
	"fmt"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/services"
	"jethur/utility"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyReportController struct {
	service services.PolicyReportServiceInterface
	db      *gorm.DB
}

func NewPolicyReportController(
	service services.PolicyReportServiceInterface,
) *PolicyReportController {
	return &PolicyReportController{
		service: service,
	}
}

func (orc *PolicyReportController) GetAllPolicyReports(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	validSortFields := map[string]string{
		"title":      "policy_reports.title",
		"start_date": "policy_reports.start_date",
		"end_date":   "policy_reports.end_date",
		"created_at": "policy_reports.created_at",
	}

	filters := utility.ParseQueryFilters([]string{}, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := policyDTOs.PolicyReportListQueryParamsDTO{
		Sort:     sort,
		SortBy:   dbField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	reports, err := orc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(orc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (prc *PolicyReportController) CreatePolicyReport(c *gin.Context) {
	var request policyDTOs.PolicyReportFormDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	policyReportID, err := prc.service.Create(currentUserID, &request)
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

	c.JSON(http.StatusCreated, gin.H{
		"message":          "Policy report generated successfully",
		"policy_report_id": policyReportID,
	})
}

func (prc *PolicyReportController) DeletePolicyReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.Delete(uint(id), currentUserID)
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

func (prc *PolicyReportController) GetPolicySummaryReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	report, err := prc.service.GetPolicySummaryReport(uint(id))
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
	c.JSON(http.StatusOK, report)
}

func (prc *PolicyReportController) ExportPolicySummaryReportToPDF(c *gin.Context) {
	IDStr := c.Param("id")
	ID, err := strconv.ParseUint(IDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	// Delegate to service to handle PDF generation
	fileName, err := prc.service.GetPolicySummaryReportExport(uint(ID), uint(currentUserID))
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
			utility.SendInternalErrResp(prc.db, c, fmt.Sprintf("Failed to export Policy Summary Report PDF: %s", err.Error()))
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
