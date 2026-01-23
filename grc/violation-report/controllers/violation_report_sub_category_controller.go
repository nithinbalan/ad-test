package controllers

import (
	"net/http"
	"strconv"

	"jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportSubCategoryController struct {
	service services.ViolationReportSubCategoryServiceInterface
	db      *gorm.DB
}

func NewViolationReportSubCategoryController(service services.ViolationReportSubCategoryServiceInterface) *ViolationReportSubCategoryController {
	return &ViolationReportSubCategoryController{
		service: service,
	}
}

func (vrsc *ViolationReportSubCategoryController) GetAllViolationReportSubCategories(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	// Define an array of allowed parameter names
	allowedParams := []string{
		"violation_report_category_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := dtos.ViolationReportSubCategoryListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	violationReportSubCategories, err := vrsc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, violationReportSubCategories)
}

func (vrsc *ViolationReportSubCategoryController) GetViolationReportSubCategoryDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	violationReportSubCategory, err := vrsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, violationReportSubCategory)
}

func (vrsc *ViolationReportSubCategoryController) CreateViolationReportSubCategory(c *gin.Context) {
	var formDTO *dtos.ViolationReportSubCategoryFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	createdViolationReportSubCategory, err := vrsc.service.Create(currentUserID, formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdViolationReportSubCategory)
}

func (vrsc *ViolationReportSubCategoryController) UpdateViolationReportSubCategory(c *gin.Context) {
	var formDTO dtos.ViolationReportSubCategoryFormDTO

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
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	err = vrsc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vrsc *ViolationReportSubCategoryController) DeleteViolationReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	failedIDs, err := vrsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (vrsc *ViolationReportSubCategoryController) ActivateViolationReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	err = vrsc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (vrsc *ViolationReportSubCategoryController) DeactivateViolationReportSubCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	err = vrsc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(vrsc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
