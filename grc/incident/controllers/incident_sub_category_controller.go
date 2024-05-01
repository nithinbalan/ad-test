package controllers

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"

	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentSubCategoryController struct {
	service services.IncidentSubCategoryServiceInterface
	db      *gorm.DB
}

func NewIncidentSubCategoryController(service services.IncidentSubCategoryServiceInterface) *IncidentSubCategoryController {
	return &IncidentSubCategoryController{
		service: service,
	}
}

func (iscc *IncidentSubCategoryController) GetAllIncidentSubCategories(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	// Define an array of allowed parameter names
	allowedParams := []string{
		"incident_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := incidentDTOs.IncidentSubCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	incidentSubCategories, err := iscc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, incidentSubCategories)
}

func (iscc *IncidentSubCategoryController) GetIncidentSubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentSubCategory, err := iscc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, incidentSubCategory)
}

func (iscc *IncidentSubCategoryController) CreateIncidentSubCategory(c *gin.Context) {
	var subCategory *incidentDTOs.IncidentSubCategoryFormDTO

	if err := c.ShouldBindJSON(&subCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}
	id, err := iscc.service.Create(currentUserID, subCategory)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iscc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (iscc *IncidentSubCategoryController) UpdateIncidentSubCategory(c *gin.Context) {
	subCategoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedIncidentSubCategory *incidentDTOs.IncidentSubCategoryFormDTO
	if err := c.ShouldBindJSON(&updatedIncidentSubCategory); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	err = iscc.service.Update(uint(subCategoryID), currentUserID, updatedIncidentSubCategory)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iscc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
func (iscc *IncidentSubCategoryController) DeleteIncidentSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	failedIDs, err := iscc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iscc *IncidentSubCategoryController) ActivateIncidentSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	err = iscc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (iscc *IncidentSubCategoryController) DeactivateIncidentSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	err = iscc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iscc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
