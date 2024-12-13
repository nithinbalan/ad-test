package controllers

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BusinessImpactAssessmentController struct {
	service services.BusinessImpactAssessmentServiceInterface
	db      *gorm.DB
}

func NewBusinessImpactAssessmentController(service services.BusinessImpactAssessmentServiceInterface) *BusinessImpactAssessmentController {
	return &BusinessImpactAssessmentController{
		service: service,
	}
}

func (biac *BusinessImpactAssessmentController) GetAllBusinessImpactAssessments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "id")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Convert the export parameter to a boolean
	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"process_ids",
		"asset_register_ids",
		"risk_level_ids",
	}
	validSortFields := map[string]string{
		"id": "business_impact_assessments.id",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, "Invalid page_size parameter")
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}
	businessImpactAssessments, err := biac.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, businessImpactAssessments)
}

func (biac *BusinessImpactAssessmentController) UpdateBusinessImpactAssessment(c *gin.Context) {
	var BusinessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO

	if err := c.ShouldBindJSON(&BusinessImpactAssessmentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(biac.db, c, err.Error())
		return
	}

	err = biac.service.Update(BusinessImpactAssessmentFormDTO, currentUserID)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(biac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}
