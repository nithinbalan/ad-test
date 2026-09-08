package controllers

import (
	"errors"
	"net/http"
	"strconv"

	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectDecisionLogController struct {
	service services.ProjectDecisionLogServiceInterface
	db      *gorm.DB
}

func NewProjectDecisionLogController(service services.ProjectDecisionLogServiceInterface) *ProjectDecisionLogController {
	return &ProjectDecisionLogController{
		service: service,
	}
}

func (pdlc *ProjectDecisionLogController) GetAllDecisionLogs(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"project_ids",
		"requested_by_user_ids",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"id":                   "project_decision_logs.id",
		"reference_number":     "project_decision_logs.reference_number",
		"decision":             "project_decision_logs.decision",
		"decision_date":        "project_decision_logs.decision_date",
		"requested_by_user_id": "project_decision_logs.requested_by_user_id",
		"created_at":           "project_decision_logs.created_at",
		"updated_at":           "project_decision_logs.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pdlc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pdlc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := projectDTOs.ProjectDecisionLogListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	decisionLogs, err := pdlc.service.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, decisionLogs)
}

func (pdlc *ProjectDecisionLogController) GetProjectDecisionLogSummary(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	var projectID uint64
	var err error

	if projectIDStr != "" {
		projectID, err = strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	}

	summary, err := pdlc.service.GetSummary(uint(projectID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (pdlc *ProjectDecisionLogController) GetDecisionLogDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLog, err := pdlc.service.GetFullDetailsByID(uint(projectID), uint(decisionLogID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, decisionLog)
}

func (pdlc *ProjectDecisionLogController) CreateDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectDecisionLogFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlc.db, c, err.Error())
		return
	}

	id, err := pdlc.service.Create(uint(projectID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (pdlc *ProjectDecisionLogController) UpdateDecisionLog(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form projectDTOs.ProjectDecisionLogFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pdlc.db, c, err.Error())
		return
	}

	if err := pdlc.service.Update(uint(projectID), uint(decisionLogID), currentUserID, &form); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project decision log updated successfully"})
}

func (pdlc *ProjectDecisionLogController) DeleteDecisionLogs(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
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
		utility.SendInternalErrResp(pdlc.db, c, err.Error())
		return
	}

	failedIDs, err := pdlc.service.Delete(uint(projectID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pdlc *ProjectDecisionLogController) DownloadDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := pdlc.service.GetDocumentPath(uint(projectID), uint(decisionLogID), uint(documentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(pdlc.db, c, e.Error())
		}
		return
	}

	c.FileAttachment(documentPath, documentPath)
}

func (pdlc *ProjectDecisionLogController) GetDocumentPreview(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	decisionLogIDStr := c.Param("decisionLogId")
	decisionLogID, err := strconv.ParseUint(decisionLogIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := pdlc.service.GetDocumentPath(uint(projectID), uint(decisionLogID), uint(documentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}

	c.File(documentPath)
}
