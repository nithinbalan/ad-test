package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentChangeRequestController struct {
	service services.DocumentChangeRequestServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestController(service services.DocumentChangeRequestServiceInterface) *DocumentChangeRequestController {
	return &DocumentChangeRequestController{
		service: service,
	}
}

func (dcrc *DocumentChangeRequestController) GetAllDocumentChangeRequests(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"document_change_request_type_ids",
		"document_change_request_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentChangeRequestQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}
	documentChangeRequests, err := dcrc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentChangeRequests)
}

func (dcrc *DocumentChangeRequestController) GetDocumentChangeRequestDetails(c *gin.Context) {
	idStr := c.Param("documentChangeRequestId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentChangeRequest, err := dcrc.service.GetFullDetailsByID(uint(id))
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
	c.JSON(http.StatusOK, documentChangeRequest)
}

func (dcrc *DocumentChangeRequestController) CreateDocumentChangeRequest(c *gin.Context) {
	var documentChangeRequestFormDTO documentDTOs.DocumentChangeRequestFormDTO

	if err := c.ShouldBindJSON(&documentChangeRequestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	id, err := dcrc.service.Create(currentUserID, &documentChangeRequestFormDTO)

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

func (dcrc *DocumentChangeRequestController) UpdateDocumentChangeRequest(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var documentChangeRequestFormDTO documentDTOs.DocumentChangeRequestFormDTO
	if err := c.ShouldBindJSON(&documentChangeRequestFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcrc.db, c, err.Error())
		return
	}

	err = dcrc.service.Update(uint(id), currentUserID, &documentChangeRequestFormDTO)
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

// func (dcrc *DocumentChangeRequestController) DeleteDocumentChangeRequest(c *gin.Context) {
// 	var request struct {
// 		IDs []uint `json:"ids"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		utility.SendInvalidErrResp(c, err)
// 		return
// 	}

// 	currentUserID, err := utility.GetCurrentUserID(c)
// 	if err != nil {
// 		utility.SendInternalErrResp(dcrc.db, c, err.Error())
// 		return
// 	}

// 	failedIDs, err := dcrc.service.Delete(request.IDs, currentUserID)
// 	if err != nil {
// 		utility.SendInternalErrResp(dcrc.db, c, err.Error())
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
// }
