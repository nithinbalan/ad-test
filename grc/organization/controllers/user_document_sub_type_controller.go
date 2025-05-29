package controllers

import (
	processDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserDocumentSubTypeController struct {
	service services.UserDocumentSubTypeServiceInterface
	db      *gorm.DB
}

func NewUserDocumentSubTypeController(service services.UserDocumentSubTypeServiceInterface) *UserDocumentSubTypeController {
	return &UserDocumentSubTypeController{
		service: service,
	}
}

func (udstc *UserDocumentSubTypeController) GetAllUserDocumentSubTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{
		"user_document_type_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := processDTOs.UserDocumentSubTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	result, err := udstc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (udstc *UserDocumentSubTypeController) GetUserDocumentSubTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := udstc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (udstc *UserDocumentSubTypeController) CreateUserDocumentSubType(c *gin.Context) {
	var dto *processDTOs.UserDocumentSubTypeFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}
	id, err := udstc.service.Create(currentUserID, dto)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udstc.db, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (udstc *UserDocumentSubTypeController) UpdateUserDocumentSubType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *processDTOs.UserDocumentSubTypeFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	err = udstc.service.Update(uint(id), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udstc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (udstc *UserDocumentSubTypeController) DeleteUserDocumentSubTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	failedIDs, err := udstc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (udstc *UserDocumentSubTypeController) ActivateUserDocumentSubTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	err = udstc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (udstc *UserDocumentSubTypeController) DeactivateUserDocumentSubTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	err = udstc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udstc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
