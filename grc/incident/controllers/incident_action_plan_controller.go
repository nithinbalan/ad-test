package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentActionPlanController struct {
	service services.IncidentActionPlanServiceInterface
	db      *gorm.DB
}

func NewIncidentActionPlanController(service services.IncidentActionPlanServiceInterface) *IncidentActionPlanController {
	return &IncidentActionPlanController{
		service: service,
	}
}

func (iapc *IncidentActionPlanController) GetAllIncidentActionPlans(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	var incidentID uint64
	var err error

	// Check if incidentID is provided in the URL
	if incidentIDStr != "" {
		incidentID, err = strconv.ParseUint(incidentIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		incidentID = 0
	}
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
		"incident_register_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(iapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(iapc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := incidentDTOs.IncidentActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
	}
	incidentActionPlans, err := iapc.service.GetAll(uint(incidentID), &queryParams)
	
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, incidentActionPlans)
}

func (iapc *IncidentActionPlanController) GetIncidentActionPlanDetails(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentActionPlan, err := iapc.service.GetFullDetailsByID(uint(incidentID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, incidentActionPlan)
}

func (iapc *IncidentActionPlanController) CreateIncidentActionPlan(c *gin.Context) {
	var IncidentActionPlanFormDTO incidentDTOs.IncidentActionPlanFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iapc.db, c, err.Error())
		return
	}

	id, err := iapc.service.Create(uint(incidentID), currentUserID, &IncidentActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (iapc *IncidentActionPlanController) UpdateIncidentActionPlan(c *gin.Context) {
	var IncidentActionPlanFormDTO incidentDTOs.IncidentActionPlanFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iapc.db, c, err.Error())
		return
	}

	err = iapc.service.Update(uint(incidentID), uint(id), currentUserID, &IncidentActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iapc *IncidentActionPlanController) DeleteIncidentActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iapc.db, c, err.Error())
		return
	}

	failedIDs, err := iapc.service.Delete(uint(incidentID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
