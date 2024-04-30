package controllers

import (
	"jethur/grc/incident/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentCategoryController struct {
	service services.IncidentCategoryServiceInterface
	db      *gorm.DB
}

func NewIncidentCategoryController(service services.IncidentCategoryServiceInterface) *IncidentCategoryController {
	return &IncidentCategoryController{
		service: service,
	}
}

func (icc *IncidentCategoryController) GetAllIncidentCategories(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	incidentCategories, err := icc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentCategories)
}

func (icc *IncidentCategoryController) GetIncidentCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentCategory, err := icc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentCategory)
}

func (icc *IncidentCategoryController) CreateIncidentCategory(c *gin.Context) {
	var category *models.IncidentCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	createdCategory, err := icc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (icc *IncidentCategoryController) UpdateIncidentCategory(c *gin.Context) {
	var updatedCategory *models.IncidentCategory

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
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	updatedCategory, err = icc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (icc *IncidentCategoryController) DeleteIncidentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	failedIDs, err := icc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (icc *IncidentCategoryController) ActivateIncidentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (icc *IncidentCategoryController) DeactivateIncidentCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
