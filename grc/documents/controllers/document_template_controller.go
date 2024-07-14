package controllers

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DocumentTemplateController struct {
	service services.DocumentTemplateServiceInterface
	db      *gorm.DB
}

func NewDocumentTemplateController(service services.DocumentTemplateServiceInterface) *DocumentTemplateController {
	return &DocumentTemplateController{
		service: service,
	}
}

func (dtc *DocumentTemplateController) GetAllDocumentTemplates(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentTemplates, err := dtc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentTemplates)
}

func (dtc *DocumentTemplateController) GetDocumentTemplateDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentTemplate, err := dtc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentTemplate)
}

func (dtc *DocumentTemplateController) CreateDocumentTemplate(c *gin.Context) {
	var documentTemplateFormDTO documentDTOs.DocumentTemplateFormDTO
	if err := c.ShouldBindJSON(&documentTemplateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	createdId, err := dtc.service.Create(currentUserID, &documentTemplateFormDTO)
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

	c.JSON(http.StatusCreated, createdId)
}

func (dtc *DocumentTemplateController) UpdateDocumentTemplate(c *gin.Context) {
	var documentTemplateFormDTO documentDTOs.DocumentTemplateFormDTO

	templateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&documentTemplateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dtc.db, c, err.Error())
		return
	}

	err = dtc.service.Update(uint(templateID), uint(currentUserID), &documentTemplateFormDTO)
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

	c.Status(http.StatusOK)
}

func (dtc *DocumentTemplateController) DeleteDocumentTemplates(c *gin.Context) {
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

func (dtc *DocumentTemplateController) ActivateDocumentTemplates(c *gin.Context) {
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

func (dtc *DocumentTemplateController) DeactivateDocumentTemplates(c *gin.Context) {
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
