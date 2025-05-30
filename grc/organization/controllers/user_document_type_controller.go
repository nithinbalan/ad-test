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

type UserDocumentTypeController struct {
	service services.UserDocumentTypeServiceInterface
	db      *gorm.DB
}

func NewUserDocumentTypeController(service services.UserDocumentTypeServiceInterface) *UserDocumentTypeController {
	return &UserDocumentTypeController{
		service: service,
	}
}

func (udtc *UserDocumentTypeController) GetAllUserDocumentTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	queryParams := processDTOs.UserDocumentTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
	}
	result, err := udtc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (udtc *UserDocumentTypeController) GetUserDocumentTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	item, err := udtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (udtc *UserDocumentTypeController) CreateUserDocumentType(c *gin.Context) {
	var dto *processDTOs.UserDocumentTypeFormDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}
	id, err := udtc.service.Create(currentUserID, dto)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udtc.db, c, e.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, id)
}

func (udtc *UserDocumentTypeController) UpdateUserDocumentType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *processDTOs.UserDocumentTypeFormDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	err = udtc.service.Update(uint(id), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(udtc.db, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (udtc *UserDocumentTypeController) DeleteUserDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	failedIDs, err := udtc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (udtc *UserDocumentTypeController) ActivateUserDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	err = udtc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (udtc *UserDocumentTypeController) DeactivateUserDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	err = udtc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(udtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
