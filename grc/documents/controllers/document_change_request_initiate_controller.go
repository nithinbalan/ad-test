package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateController struct {
	service services.DocumentChangeRequestInitiateServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestInitiateController(service services.DocumentChangeRequestInitiateServiceInterface) *DocumentChangeRequestInitiateController {
	return &DocumentChangeRequestInitiateController{
		service: service,
	}
}

// Get all document change requests with filtering, sorting, and pagination
func (dcrc *DocumentChangeRequestInitiateController) GetAllDocumentChangeRequestInitiates(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

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

	params := documentDTOs.DocumentChangeRequestInitiateQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	documentChangeRequests, err := dcrc.service.GetAll(&params)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documentChangeRequests)
}

// Get details of a single document change request by ID
func (dcrc *DocumentChangeRequestInitiateController) GetDocumentChangeRequestInitiateDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := dcrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, details)
}

// Create a new document change request
func (dcrc *DocumentChangeRequestInitiateController) CreateDocumentChangeRequestInitiate(c *gin.Context) {
	var formDTO documentDTOs.DocumentChangeRequestInitiateFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	documentID, err := strconv.Atoi(c.Param("document_id"))
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid document_id parameter"))
		return
	}

	id, err := dcrc.service.Create(currentUserID, uint(documentID), &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

// Update an existing document change request
func (dcrc *DocumentChangeRequestInitiateController) UpdateDocumentChangeRequestInitiate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO documentDTOs.DocumentChangeRequestInitiateFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// Delete document change requests by IDs
func (dcrc *DocumentChangeRequestInitiateController) DeleteDocumentChangeRequestInitiate(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	failedIDs, err := dcrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
