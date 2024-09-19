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

type DocumentAcknowledgeUserController struct {
	service services.DocumentAcknowledgeUserServiceInterface
	db      *gorm.DB
}

func NewDocumentAcknowledgeUserController(service services.DocumentAcknowledgeUserServiceInterface) *DocumentAcknowledgeUserController {
	return &DocumentAcknowledgeUserController{
		service: service,
	}
}

func (dauc *DocumentAcknowledgeUserController) GetAllDocumentAcknowledgeUsers(c *gin.Context) {
	idStr := c.Param("id")
	documentId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"user_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(dauc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(dauc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := documentDTOs.DocumentAcknowledgeListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	users, err := dauc.service.GetAll(uint(documentId), &queryParams)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (dauc *DocumentAcknowledgeUserController) CreateDocumentAcknowledgeUser(c *gin.Context) {
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dauc.db, c, err.Error())
		return
	}

	id, err := dauc.service.Create(uint(documentId), currentUserID)
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
			utility.SendInternalErrResp(dauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (dauc *DocumentAcknowledgeUserController) DeleteDocumentAcknowledgeUser(c *gin.Context) {
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dauc.db, c, err.Error())
		return
	}

	failedID, err := dauc.service.Delete(uint(documentId), currentUserID)
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
			utility.SendInternalErrResp(dauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
