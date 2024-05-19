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

type ControlActionPlanController struct {
	service services.ControlActionPlanServiceInterface
	db      *gorm.DB
}

func NewControlActionPlanController(service services.ControlActionPlanServiceInterface) *ControlActionPlanController {
	return &ControlActionPlanController{
		service: service,
	}
}

func (capc *ControlActionPlanController) GetAllControlActionPlans(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
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
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.ControlActionPlanListQueryParamsDTO{
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
	controlActionPlans, err := capc.service.GetAll(uint(controlID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlActionPlans)
}

func (capc *ControlActionPlanController) GetControlActionPlanDetails(c *gin.Context) {
	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlan, err := capc.service.GetFullDetailsByID(uint(controlID), uint(controlActionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, controlActionPlan)
}

func (capc *ControlActionPlanController) CreateControlActionPlan(c *gin.Context) {
	var ControlActionPlanFormDTO organizationDTOs.ControlActionPlanFormDTO

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ControlActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	id, err := capc.service.Create(uint(controlID), currentUserID, &ControlActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (capc *ControlActionPlanController) UpdateControlActionPlan(c *gin.Context) {
	var ControlActionPlanFormDTO organizationDTOs.ControlActionPlanFormDTO

	controlID, err := strconv.Atoi(c.Param("controlId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlActionPlanID, err := strconv.Atoi(c.Param("controlActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&ControlActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	err = capc.service.Update(uint(controlID), uint(controlActionPlanID), currentUserID, &ControlActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (capc *ControlActionPlanController) DeleteControlActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	controlID, err := strconv.Atoi(c.Param("controlId"))
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
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	failedIDs, err := capc.service.Delete(uint(controlID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
