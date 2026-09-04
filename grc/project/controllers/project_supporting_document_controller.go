package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"jethur/grc/project/dtos"
	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSupportingDocumentController struct {
	service services.ProjectSupportingDocumentServiceInterface
	db      *gorm.DB
}

func NewProjectSupportingDocumentController(service services.ProjectSupportingDocumentServiceInterface, db *gorm.DB) *ProjectSupportingDocumentController {
	return &ProjectSupportingDocumentController{
		service: service,
		db:      db,
	}
}

func (psdc *ProjectSupportingDocumentController) GetAllSupportingDocuments(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "id")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"created_by_ids",
		"updated_by_ids",
	}

	validSortFields := map[string]string{
		"id":         "project_supporting_documents.id",
		"title_en":   "project_supporting_documents.title_en",
		"title_ar":   "project_supporting_documents.title_ar",
		"created_at": "project_supporting_documents.created_at",
		"updated_at": "project_supporting_documents.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(psdc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(psdc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.ProjectSupportingDocumentListQueryParamsDTO{
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

	supportingDocs, err := psdc.service.GetAll(uint(projectID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, supportingDocs)
}

func (psdc *ProjectSupportingDocumentController) GetSupportingDocumentDetails(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	supportingDocIDStr := c.Param("supportingDocumentId")
	supportingDocID, err := strconv.ParseUint(supportingDocIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	supportingDoc, err := psdc.service.GetFullDetailsByID(uint(projectID), uint(supportingDocID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, supportingDoc)
}

func (psdc *ProjectSupportingDocumentController) CreateSupportingDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSupportingDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psdc.db, c, err.Error())
		return
	}

	id, err := psdc.service.Create(uint(projectID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (psdc *ProjectSupportingDocumentController) UpdateSupportingDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	supportingDocIDStr := c.Param("supportingDocumentId")
	supportingDocID, err := strconv.ParseUint(supportingDocIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ProjectSupportingDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psdc.db, c, err.Error())
		return
	}

	err = psdc.service.Update(uint(projectID), uint(supportingDocID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(psdc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (psdc *ProjectSupportingDocumentController) DeleteSupportingDocuments(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psdc.db, c, err.Error())
		return
	}

	failedIDs, err := psdc.service.Delete(uint(projectID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(psdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (psdc *ProjectSupportingDocumentController) DownloadDocument(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	supportingDocIDStr := c.Param("supportingDocumentId")
	supportingDocID, err := strconv.ParseUint(supportingDocIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	filePath, err := psdc.service.GetDocumentPath(uint(projectID), uint(supportingDocID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(psdc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}

func (psdc *ProjectSupportingDocumentController) GetDocumentPreview(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	supportingDocIDStr := c.Param("supportingDocumentId")
	supportingDocID, err := strconv.ParseUint(supportingDocIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentIDStr := c.Param("documentId")
	documentID, err := strconv.ParseUint(documentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	filePath, err := psdc.service.GetDocumentPath(uint(projectID), uint(supportingDocID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(psdc.db, c, err.Error())
		}
		return
	}

	c.File(filePath)
}
