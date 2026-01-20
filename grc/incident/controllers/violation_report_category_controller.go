package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportCategoryController struct {
	service services.ViolationReportCategoryServiceInterface
	db      *gorm.DB
}

func NewViolationReportCategoryController(service services.ViolationReportCategoryServiceInterface) *ViolationReportCategoryController {
	return &ViolationReportCategoryController{
		service: service,
	}
}

func (vrc *ViolationReportCategoryController) GetAllViolationReportCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	// Define an array of allowed parameter names
	allowedParams := []string{
		"department_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.ViolationReportCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	violationReportCategories, err := vrc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportCategories)
}

func (vrc *ViolationReportCategoryController) GetViolationReportCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportCategory, err := vrc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, violationReportCategory)
}

func (vrc *ViolationReportCategoryController) CreateViolationReportCategory(c *gin.Context) {
	var formDTO *dtos.ViolationReportCategoryFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	createdViolationReportCategory, err := vrc.service.Create(currentUserID, formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdViolationReportCategory)
}

func (vrc *ViolationReportCategoryController) UpdateViolationReportCategory(c *gin.Context) {
	var formDTO dtos.ViolationReportCategoryFormDTO

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	err = vrc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vrc *ViolationReportCategoryController) DeleteViolationReportCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	failedIDs, err := vrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (vrc *ViolationReportCategoryController) ActivateViolationReportCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	err = vrc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (vrc *ViolationReportCategoryController) DeactivateViolationReportCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	err = vrc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
