package controllers

import (
	"errors"
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
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, "Invalid export parameter")
		return
	}

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

	validSortFields := map[string]string{
		"reference_number":         "controls.reference_number",
		"code":                     "controls.code",
		"title_en":                 "controls.title_en",
		"title_ar":                 "controls.title_ar",
		"description":              "controls.description",
		"objective":                "controls.objective",
		"control_category_en":      "control_categories.title_en",
		"control_category_ar":      "control_categories.title_ar",
		"control_type_en":          "control_types.title_en",
		"control_type_ar":          "control_types.title_ar",
		"owner_en":                 "owner.first_name_en",
		"owner_ar":                 "owner.first_name_ar",
		"review_frequency_en":      "review_frequencies.title_en",
		"review_frequency_ar":      "review_frequencies.title_ar",
		"implementation_status_en": "control_implementation_statuses.title_en",
		"implementation_status_ar": "control_implementation_statuses.title_ar",
		"implementation_date":      "controls.implemented_at",
		"efficiency_status_en":     "control_efficiency_statuses.title_en",
		"efficiency_status_ar":     "control_efficiency_statuses.title_ar",
		"last_test_date":           "controls.last_test_date",
		"next_test_date":           "controls.next_test_date",
		"test_frequency_en":        "test_frequencies.title_en",
		"test_frequency_ar":        "test_frequencies.title_ar",
		"last_review_note":         "controls.last_review_note",
		"next_review_date":         "controls.next_review_date",
		"status_en":                "lifecycle_statuses.title_en",
		"status_ar":                "lifecycle_statuses.title_ar",
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

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.ControlListQueryParamsDTO{
		Sort:               sort,
		SortBy:             dbField,
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
		Export:             export,
	}

	controls, err := cc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(cc.db, c, err.Error())
			return
		}

		excelFile, err := cc.service.ExportControls(controls.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(cc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Controls.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, controls)
}

func (crc *ControlController) GetControlDetails(c *gin.Context) {
	idStr := c.Param("controlId")

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
	controlID, err := strconv.Atoi(c.Param("controlId"))
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

func (cc *ControlController) GetControlsExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	cc.GetAllControls(c)
}

func (cc *ControlController) MarkControlAsApplicable(c *gin.Context) {

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.MarkControlAsApplicable(uint(controlID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (cc *ControlController) MarkControlAsNotApplicable(c *gin.Context) {

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cc.db, c, err.Error())
		return
	}

	err = cc.service.MarkControlAsNotApplicable(uint(controlID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
