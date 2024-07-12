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

type DocumentCategoryController struct {
	service services.DocumentCategoryServiceInterface
	db      *gorm.DB
}

func NewDocumentCategoryController(service services.DocumentCategoryServiceInterface) *DocumentCategoryController {
	return &DocumentCategoryController{
		service: service,
	}
}

func (dcc *DocumentCategoryController) GetAllDocumentCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	documentCategories, err := dcc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, documentCategories)
}

func (dcc *DocumentCategoryController) GetDocumentCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentCategory, err := dcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, documentCategory)
}

func (dcc *DocumentCategoryController) CreateDocumentCategory(c *gin.Context) {
	var category *models.DocumentCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	createdCategory, err := dcc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (dcc *DocumentCategoryController) UpdateDocumentCategory(c *gin.Context) {
	var updatedCategory *models.DocumentCategory

	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	updatedCategory, err = dcc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(dcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (dcc *DocumentCategoryController) DeleteDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	failedIDs, err := dcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (dcc *DocumentCategoryController) ActivateDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	err = dcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (dcc *DocumentCategoryController) DeactivateDocumentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	err = dcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(dcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
