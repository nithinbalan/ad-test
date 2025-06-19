package controllers

import (
	"jethur/grc/models"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectiveCategoryController struct {
	service services.StrategicObjectiveCategoryServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveCategoryController(service services.StrategicObjectiveCategoryServiceInterface) *StrategicObjectiveCategoryController {
	return &StrategicObjectiveCategoryController{
		service: service,
	}
}

func (soc *StrategicObjectiveCategoryController) GetAllStrategicObjectiveCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	strategicObjectiveCategories, err := soc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, strategicObjectiveCategories)
}

func (soc *StrategicObjectiveCategoryController) GetStrategicObjectiveCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjectiveCategory, err := soc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, strategicObjectiveCategory)
}

func (soc *StrategicObjectiveCategoryController) CreateStrategicObjectiveCategory(c *gin.Context) {
	var category *models.StrategicObjectiveCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	createdCategory, err := soc.service.Create(currentUserID, category)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdCategory)
}

func (soc *StrategicObjectiveCategoryController) UpdateStrategicObjectiveCategory(c *gin.Context) {
	var updatedCategory *models.StrategicObjectiveCategory

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
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	updatedCategory, err = soc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (soc *StrategicObjectiveCategoryController) DeleteStrategicObjectiveCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	failedIDs, err := soc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (soc *StrategicObjectiveCategoryController) ActivateStrategicObjectiveCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicObjectiveCategoryController) DeactivateStrategicObjectiveCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
