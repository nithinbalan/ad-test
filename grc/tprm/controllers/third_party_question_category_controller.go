package controllers

import (
	"jethur/grc/models"
	"jethur/grc/tprm/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyQuestionCategoryController struct {
	service services.ThirdPartyQuestionCategoryServiceInterface
	db      *gorm.DB
}

func NewThirdPartyQuestionCategoryController(service services.ThirdPartyQuestionCategoryServiceInterface) *ThirdPartyQuestionCategoryController {
	return &ThirdPartyQuestionCategoryController{
		service: service,
	}
}

func (tpqcc *ThirdPartyQuestionCategoryController) GetAllThirdPartyQuestionCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	tprmCategories, err := tpqcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, tprmCategories)
}

func (tpqcc *ThirdPartyQuestionCategoryController) GetThirdPartyQuestionCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	tprmType, err := tpqcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, tprmType)
}

func (tpqcc *ThirdPartyQuestionCategoryController) CreateThirdPartyQuestionCategory(c *gin.Context) {
	var tprmType *models.ThirdPartyQuestionCategory

	if err := c.ShouldBindJSON(&tprmType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	createdType, err := tpqcc.service.Create(currentUserID, tprmType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpqcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (tpqcc *ThirdPartyQuestionCategoryController) UpdateThirdPartyQuestionCategory(c *gin.Context) {
	var updatedType *models.ThirdPartyQuestionCategory

	tprmTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	updatedType, err = tpqcc.service.Update(uint(tprmTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpqcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (tpqcc *ThirdPartyQuestionCategoryController) DeleteThirdPartyQuestionCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	failedIDs, err := tpqcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (tpqcc *ThirdPartyQuestionCategoryController) ActivateThirdPartyQuestionCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	err = tpqcc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (tpqcc *ThirdPartyQuestionCategoryController) DeactivateThirdPartyQuestionCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	err = tpqcc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tpqcc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
