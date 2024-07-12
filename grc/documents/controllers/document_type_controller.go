package controllers

import (
	"jethur/grc/documents/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentTypeController struct {
	service services.DocumentTypeServiceInterface
	db      *gorm.DB
}

func NewDocumentTypeController(service services.DocumentTypeServiceInterface) *DocumentTypeController {
	return &DocumentTypeController{
		service: service,
	}
}

func (dtc *DocumentTypeController) GetAllDocumentTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	typeFilter := c.Query("type") // Added to handle type filter

	documentTypes, err := dtc.service.GetAll(sort, sortBy, keyword, typeFilter, status)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentTypes)
}

func (dtc *DocumentTypeController) GetDocumentTypeDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentType, err := dtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentType)
}

func (dtc *DocumentTypeController) CreateDocumentType(c *gin.Context) {
	var documentType *models.DocumentType
	if err := c.ShouldBindJSON(&documentType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	createdDocumentType, err := dtc.service.Create(currentUserID, documentType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdDocumentType)
}

func (dtc *DocumentTypeController) UpdateDocumentType(c *gin.Context) {
	var updatedDocumentType *models.DocumentType

	documentTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedDocumentType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	updatedDocumentType, err = dtc.service.Update(uint(documentTypeID), currentUserID, updatedDocumentType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedDocumentType)
}

func (dtc *DocumentTypeController) DeleteDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	failedIDs, err := dtc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (dtc *DocumentTypeController) ActivateDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	err = dtc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (dtc *DocumentTypeController) DeactivateDocumentTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	err = dtc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
