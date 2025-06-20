package controllers

import (
	dtos "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDocumentController struct {
	service services.UserDocumentServiceInterface
	db      *gorm.DB
}

func NewUserDocumentController(service services.UserDocumentServiceInterface) *UserDocumentController {
	return &UserDocumentController{
		service: service,
	}
}

func (udc *UserDocumentController) GetAllUserDocuments(c *gin.Context) {
	UserIDStr := c.Param("userId")

	UserID, err := strconv.ParseUint(UserIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	allowedParams := []string{
		"user_document_type_ids",
		"user_document_sub_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.UserDocumentListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Filters: filters,
	}
	result, err := udc.service.GetAll(uint(UserID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (udc *UserDocumentController) GetUserDocumentDetails(c *gin.Context) {
	UserIDStr := c.Param("userId")

	UserID, err := strconv.ParseUint(UserIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	item, err := udc.service.GetFullDetailsByID(uint(UserID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(udc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, item)
}

func (udc *UserDocumentController) CreateUserDocument(c *gin.Context) {
	UserIDStr := c.Param("userId")

	UserID, err := strconv.ParseUint(UserIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *dtos.UserDocumentFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}
	id, err := udc.service.Create(uint(UserID), currentUserID, dto)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udc.db, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (udc *UserDocumentController) UpdateUserDocument(c *gin.Context) {

	UserIDStr := c.Param("userId")

	UserID, err := strconv.ParseUint(UserIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *dtos.UserDocumentFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}
	err = udc.service.Update(uint(UserID), uint(id), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (udc *UserDocumentController) DeleteUserDocuments(c *gin.Context) {

	UserIDStr := c.Param("userId")

	UserID, err := strconv.ParseUint(UserIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}

	failedIDs, err := udc.service.Delete(uint(UserID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (udc *UserDocumentController) GetUserDocumentPreview(c *gin.Context) {
	userIDStr := c.Param("userId")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	userDocumentPath, err := udc.service.GetDocumentPath(uint(userID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(udc.db, c, err.Error())
		}
		return
	}
	c.File(userDocumentPath)
}

func (udc *UserDocumentController) DownloadUserDocument(c *gin.Context) {
	userIDStr := c.Param("userId")

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	userDocumentPath, err := udc.service.GetDocumentPath(uint(userID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(udc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(userDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(userDocumentPath)
}
