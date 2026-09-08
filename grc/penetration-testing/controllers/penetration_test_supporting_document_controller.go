package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	"jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestSupportingDocumentController struct {
	service services.PenetrationTestSupportingDocumentServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestSupportingDocumentController(service services.PenetrationTestSupportingDocumentServiceInterface, db *gorm.DB) *PenetrationTestSupportingDocumentController {
	return &PenetrationTestSupportingDocumentController{
		service: service,
		db:      db,
	}
}

func (ptsdc *PenetrationTestSupportingDocumentController) GetAllSupportingDocuments(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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
		"id":         "penetration_test_supporting_documents.id",
		"title_en":   "penetration_test_supporting_documents.title_en",
		"title_ar":   "penetration_test_supporting_documents.title_ar",
		"created_at": "penetration_test_supporting_documents.created_at",
		"updated_at": "penetration_test_supporting_documents.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(ptsdc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(ptsdc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.PenetrationTestSupportingDocumentListQueryParamsDTO{
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

	supportingDocs, err := ptsdc.service.GetAll(uint(penetrationTestID), &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptsdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, supportingDocs)
}

func (ptsdc *PenetrationTestSupportingDocumentController) GetSupportingDocumentDetails(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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

	supportingDoc, err := ptsdc.service.GetFullDetailsByID(uint(penetrationTestID), uint(supportingDocID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptsdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, supportingDoc)
}

func (ptsdc *PenetrationTestSupportingDocumentController) CreateSupportingDocument(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.PenetrationTestSupportingDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptsdc.db, c, err.Error())
		return
	}

	id, err := ptsdc.service.Create(uint(penetrationTestID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptsdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ptsdc *PenetrationTestSupportingDocumentController) UpdateSupportingDocument(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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

	var form dtos.PenetrationTestSupportingDocumentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptsdc.db, c, err.Error())
		return
	}

	err = ptsdc.service.Update(uint(penetrationTestID), uint(supportingDocID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptsdc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptsdc *PenetrationTestSupportingDocumentController) DeleteSupportingDocuments(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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
		utility.SendInternalErrResp(ptsdc.db, c, err.Error())
		return
	}

	failedIDs, err := ptsdc.service.Delete(uint(penetrationTestID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptsdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ptsdc *PenetrationTestSupportingDocumentController) DownloadDocument(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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

	filePath, err := ptsdc.service.GetDocumentPath(uint(penetrationTestID), uint(supportingDocID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ptsdc.db, c, err.Error())
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

func (ptsdc *PenetrationTestSupportingDocumentController) GetDocumentPreview(c *gin.Context) {
	penetrationTestIDStr := c.Param("id")
	penetrationTestID, err := strconv.ParseUint(penetrationTestIDStr, 10, 64)
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

	filePath, err := ptsdc.service.GetDocumentPath(uint(penetrationTestID), uint(supportingDocID), uint(documentID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ptsdc.db, c, err.Error())
		}
		return
	}

	c.File(filePath)
}
