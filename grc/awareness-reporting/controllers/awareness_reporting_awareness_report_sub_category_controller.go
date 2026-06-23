package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessReportSubCategoryController struct {
	service services.AwarenessReportSubCategoryServiceInterface
	db      *gorm.DB
}

func NewAwarenessReportSubCategoryController(
	service services.AwarenessReportSubCategoryServiceInterface,
) *AwarenessReportSubCategoryController {
	return &AwarenessReportSubCategoryController{
		service: service,
	}
}

func (arscc *AwarenessReportSubCategoryController) GetAllAwarenessReportSubCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	subCategories, err := arscc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, subCategories)
}

func (arscc *AwarenessReportSubCategoryController) GetAwarenessReportSubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	subCategory, err := arscc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, subCategory)
}

func (arscc *AwarenessReportSubCategoryController) CreateAwarenessReportSubCategory(c *gin.Context) {
	var form models.AwarenessReportSubCategory

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	createdSubCategory, err := arscc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdSubCategory)
}

func (arscc *AwarenessReportSubCategoryController) UpdateAwarenessReportSubCategory(c *gin.Context) {
	var form models.AwarenessReportSubCategory

	subCategoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	updatedSubCategory, err := arscc.service.Update(
		uint(subCategoryID),
		currentUserID,
		&form,
	)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedSubCategory)
}

func (arscc *AwarenessReportSubCategoryController) DeleteAwarenessReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	failedIDs, err := arscc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}

func (arscc *AwarenessReportSubCategoryController) ActivateAwarenessReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	err = arscc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (arscc *AwarenessReportSubCategoryController) DeactivateAwarenessReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	err = arscc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}