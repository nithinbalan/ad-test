package controllers

import (
	"errors"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftReportReviewController struct {
	service services.AuditDraftReportReviewServiceInterface
	db      *gorm.DB
}

func NewAuditDraftReportReviewController(service services.AuditDraftReportReviewServiceInterface) *AuditDraftReportReviewController {
	return &AuditDraftReportReviewController{
		service: service,
	}
}

func (adrrc *AuditDraftReportReviewController) GetAllAuditDraftReportReviews(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "id")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_ids",
		"workflow_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	validSortFields := map[string]string{
		"id":      "audit_report_reviews.id",
		"comment": "audit_report_reviews.comment",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := dtos.AuditReportReviewListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditDraftReportComments, err := adrrc.service.GetAll(uint(auditID), uint(auditDraftReportID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, auditDraftReportComments)
}


func (adrrc *AuditDraftReportReviewController) SubmitAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowSubmitFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Submit(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrrc *AuditDraftReportReviewController) ApproveAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowApproveFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Approve(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrrc *AuditDraftReportReviewController) RevertAuditDraftReport(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	auditDraftReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrrc.db, c, err.Error())
		return
	}
	var auditReportWorkflowFormDTO dtos.AuditReportWorkflowApproveFormDTO

	if err := c.ShouldBindJSON(&auditReportWorkflowFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = adrrc.service.Revert(uint(auditID), uint(auditDraftReportID), currentUserID, &auditReportWorkflowFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(adrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

