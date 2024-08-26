package controllers

import (
	"jethur/grc/audit/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingCategoryController struct {
	service services.AuditFindingCategoryServiceInterface
	db      *gorm.DB
}

func NewAuditFindingCategoryController(service services.AuditFindingCategoryServiceInterface) *AuditFindingCategoryController {
	return &AuditFindingCategoryController{
		service: service,
	}
}

func (afcc *AuditFindingCategoryController) GetAllAuditFindingCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	auditFindingCategories, err := afcc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(afcc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditFindingCategories)
}

func (afcc *AuditFindingCategoryController) GetAuditFindingCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingCategory, err := afcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, auditFindingCategory)
}

func (afcc *AuditFindingCategoryController) UpdateAuditFindingCategory(c *gin.Context) {
	var updatedCategory *models.AuditFindingCategory

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
		utility.SendInternalErrResp(afcc.db, c, err.Error())
		return
	}

	updatedCategory, err = afcc.service.Update(uint(categoryID), currentUserID, updatedCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedCategory)
}

func (afcc *AuditFindingCategoryController) DeleteAuditFindingCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afcc.db, c, err.Error())
		return
	}

	failedIDs, err := afcc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(afcc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
