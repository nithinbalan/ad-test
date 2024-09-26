package controllers

import (
	"errors"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditDraftReportCommentController struct {
	service services.AuditDraftReportCommentServiceInterface
	db      *gorm.DB
}

func NewAuditDraftReportCommentController(service services.AuditDraftReportCommentServiceInterface) *AuditDraftReportCommentController {
	return &AuditDraftReportCommentController{
		service: service,
	}
}

func (adrcc *AuditDraftReportCommentController) GetAllAuditDraftReportComments(c *gin.Context) {
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
	sortBy := c.DefaultQuery("sort_by", "comment")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audit_draft_report_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, "Invalid page_size parameter")
		return
	}
	validSortFields := map[string]string{
		"comment": "audit_draft_report_comments.comment",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := dtos.AuditDraftReportCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	auditDraftReportComments, err := adrcc.service.GetAll(uint(auditID), uint(auditDraftReportID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}

	// Otherwise, return the data as JSON
	c.JSON(http.StatusOK, auditDraftReportComments)
}

func (adrcc *AuditDraftReportCommentController) GetAuditDraftReportCommentDetails(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	_, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("commentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditDraftFinding, err := adrcc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(adrcc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditDraftFinding)
}

func (adrcc *AuditDraftReportCommentController) CreateAuditDraftReportComment(c *gin.Context) {
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
	var auditDraftReportCommentFormDTO dtos.AuditDraftReportCommentFormDTO

	if err := c.ShouldBindJSON(&auditDraftReportCommentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}

	id, err := adrcc.service.Create(uint(auditID), uint(auditDraftReportID), currentUserID, &auditDraftReportCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (adrcc *AuditDraftReportCommentController) UpdateAuditDraftReportComment(c *gin.Context) {
	var AuditDraftReportCommentFormDTO dtos.AuditDraftReportCommentFormDTO
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
	idStr := c.Param("commentId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&AuditDraftReportCommentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}

	err = adrcc.service.Update(uint(auditID), uint(auditDraftReportID), uint(id), currentUserID, &AuditDraftReportCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (adrcc *AuditDraftReportCommentController) DeleteAuditDraftReportComments(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

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

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adrcc.db, c, err.Error())
		return
	}

	failedIDs, err := adrcc.service.Delete(uint(auditID), uint(auditDraftReportID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
