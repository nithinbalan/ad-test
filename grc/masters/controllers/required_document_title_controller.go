package controllers

import (
	"jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RequiredDocumentTitleController struct {
	service services.RequiredDocumentTitleServiceInterface
	db      *gorm.DB
}

func NewRequiredDocumentTitleController(service services.RequiredDocumentTitleServiceInterface) *RequiredDocumentTitleController {
	return &RequiredDocumentTitleController{
		service: service,
	}
}

func (rdtc *RequiredDocumentTitleController) GetAllRequiredDocumentTitles(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentTitles, err := rdtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentTitles)
}

func (rdtc *RequiredDocumentTitleController) GetRequiredDocumentTitleDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentTitle, err := rdtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentTitle)
}

func (rdtc *RequiredDocumentTitleController) CreateRequiredDocumentTitle(c *gin.Context) {
	var document *models.RequiredDocumentTitle

	if err := c.ShouldBindJSON(&document); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	createdDocument, err := rdtc.service.Create(currentUserID, document)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rdtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdDocument)
}

func (rdtc *RequiredDocumentTitleController) UpdateRequiredDocumentTitle(c *gin.Context) {
	var updatedDocument *models.RequiredDocumentTitle

	documentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedDocument); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	updatedDocument, err = rdtc.service.Update(uint(documentID), currentUserID, updatedDocument)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rdtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedDocument)
}

func (rdtc *RequiredDocumentTitleController) DeleteRequiredDocumentTitles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	failedIDs, err := rdtc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (rdtc *RequiredDocumentTitleController) ActivateRequiredDocumentTitles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	err = rdtc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (rdtc *RequiredDocumentTitleController) DeactivateRequiredDocumentTitles(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	err = rdtc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rdtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
