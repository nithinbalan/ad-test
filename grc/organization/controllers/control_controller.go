package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlController struct {
	service services.ControlServiceInterface
	db      *gorm.DB
}

func NewControlController(service services.ControlServiceInterface) *ControlController {
	return &ControlController{
		service: service,
	}
}

func (cc *ControlController) GetAllControls(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"control_type_ids",
		"control_category_ids",
		"owner_unit_ids",
		"owner_ids",
		"control_implementation_status_ids",
		"control_efficiency_status_ids",
		"test_responsible_user_ids",
		"lifecycle_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	//Parse last_review_date_from, last_review_date_to, next_review_date_from, next_review_date_to, last_test_date_from, last_test_date_to, next_test_date_from, next_test_date_to, implemented_at_from, implemented_at_to parameters
	lastReviewDateFrom := c.Query("last_review_date_from")
	lastReviewDateTo := c.Query("last_review_date_to")
	nextReviewDateFrom := c.Query("next_review_date_from")
	nextReviewDateTo := c.Query("next_review_date_to")
	lastTestDateFrom := c.Query("last_test_date_from")
	lastTestDateTo := c.Query("last_test_date_to")
	nextTestDateFrom := c.Query("next_test_date_from")
	nextTestDateTo := c.Query("next_test_date_to")
	implementedAtFrom := c.Query("implemented_at_from")
	implementedAtTo := c.Query("implemented_at_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlListQueryParamsDTO{
		Sort:               sort,
		SortBy:             sortBy,
		Keyword:            keyword,
		Status:             status,
		Page:               page,
		PageSize:           pageSize,
		Filters:            filters,
		QuickFilter:        quickFilter,
		LastReviewDateFrom: lastReviewDateFrom,
		LastReviewDateTo:   lastReviewDateTo,
		NextReviewDateFrom: nextReviewDateFrom,
		NextReviewDateTo:   nextReviewDateTo,
		LastTestDateFrom:   lastTestDateFrom,
		LastTestDateTo:     lastTestDateTo,
		NextTestDateFrom:   nextTestDateFrom,
		NextTestDateTo:     nextTestDateTo,
		ImplementedAtFrom:  implementedAtFrom,
		ImplementedAtTo:    implementedAtTo,
	}
	controls, err := cc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, controls)
}

func (crc *ControlController) GetControlDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	control, err := crc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, control)
}

func (cc *ControlController) CreateControl(c *gin.Context) {
	var controlFormDTO organizationDTOs.ControlFormDTO

	if err := c.ShouldBindJSON(&controlFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	id, err := cc.service.Create(currentUserID, &controlFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (cc *ControlController) UpdateControl(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var controlFormDTO organizationDTOs.ControlFormDTO
	if err := c.ShouldBindJSON(&controlFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.Update(uint(controlID), currentUserID, &controlFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (cc *ControlController) DeleteControls(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	failedIDs, err := cc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (cc *ControlController) ArchiveControls(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	successfulIds, err := cc.service.Archive(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}

func (cc *ControlController) RestoreControls(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	successfulIds, err := cc.service.Restore(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfulIds": successfulIds})
}
