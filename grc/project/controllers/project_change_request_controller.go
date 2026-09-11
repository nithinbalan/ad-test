package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectChangeRequestController struct {
	service services.ProjectChangeRequestServiceInterface
	db      *gorm.DB
}

func NewProjectChangeRequestController(service services.ProjectChangeRequestServiceInterface) *ProjectChangeRequestController {
	return &ProjectChangeRequestController{
		service: service,
	}
}

func (pcrc *ProjectChangeRequestController) GetAllProjectChangeRequests(c *gin.Context) {
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
		"workflow_status_ids",
		"requested_by_user_ids",
		"created_by_ids",
		"updated_by_ids",
	}

	validSortFields := map[string]string{
		"id":                 "project_change_requests.id",
		"reference_number":   "project_change_requests.reference_number",
		"title":              "project_change_requests.title",
		"requested_on":       "project_change_requests.requested_on",
		"workflow_status_id": "project_change_requests.workflow_status_id",
		"created_at":         "project_change_requests.created_at",
		"updated_at":         "project_change_requests.updated_at",
		"impact":             "cr_impacts.impact",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(pcrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(pcrc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := projectDTOs.ProjectChangeRequestListQueryParamsDTO{
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

	changeRequests, err := pcrc.service.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, changeRequests)
}

func (pcrc *ProjectChangeRequestController) GetProjectChangeRequestSummary(c *gin.Context) {
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

	summary, err := pcrc.service.GetSummary(uint(projectID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}

func (pcrc *ProjectChangeRequestController) GetProjectChangeRequestDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("changeRequestId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequest, err := pcrc.service.GetFullDetailsByID(uint(projectID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, changeRequest)
}

func (pcrc *ProjectChangeRequestController) CreateProjectChangeRequest(c *gin.Context) {
	var formDTO projectDTOs.ProjectChangeRequestFormDTO
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrc.db, c, err.Error())
		return
	}

	id, err := pcrc.service.Create(uint(projectID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pcrc *ProjectChangeRequestController) UpdateProjectChangeRequest(c *gin.Context) {
	var formDTO projectDTOs.ProjectChangeRequestFormDTO
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestID, err := strconv.Atoi(c.Param("changeRequestId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pcrc.db, c, err.Error())
		return
	}

	err = pcrc.service.Update(uint(projectID), uint(changeRequestID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (pcrc *ProjectChangeRequestController) DeleteProjectChangeRequests(c *gin.Context) {
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
		utility.SendInternalErrResp(pcrc.db, c, err.Error())
		return
	}

	failedIDs, err := pcrc.service.Delete(uint(projectID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(pcrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (pcrc *ProjectChangeRequestController) DownloadDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	documentPath, err := pcrc.service.GetDocumentPath(uint(projectID), uint(changeRequestID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
		default:
			utility.SendInternalErrResp(pcrc.db, c, err.Error())
		}
		return
	}

	c.FileAttachment(documentPath, documentPath)
}

func (pcrc *ProjectChangeRequestController) GetDocumentPreview(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	changeRequestIDStr := c.Param("changeRequestId")
	changeRequestID, err := strconv.ParseUint(changeRequestIDStr, 10, 32)
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

	documentPath, err := pcrc.service.GetDocumentPath(uint(projectID), uint(changeRequestID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
		default:
			utility.SendInternalErrResp(pcrc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}
