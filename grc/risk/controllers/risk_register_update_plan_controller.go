package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanController struct {
	service services.RiskRegisterUpdatePlanServiceInterface
	db      *gorm.DB
}

func NewRiskRegisterUpdatePlanController(service services.RiskRegisterUpdatePlanServiceInterface) *RiskRegisterUpdatePlanController {
	return &RiskRegisterUpdatePlanController{
		service: service,
	}
}

func (rrupc *RiskRegisterUpdatePlanController) GetAllRiskRegisterUpdatePlan(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	status := c.Query("status")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	// Allowed filters
	allowedParams := []string{
		"division_ids",
		"business_unit_ids",
		"risk_register_update_plan_status_ids",
		"risk_champion_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"created_at":             "risk_register_update_plans.created_at",
		"updated_at":             "risk_register_update_plans.updated_at",
		"division_title_en":      "div.title_en",
		"division_title_ar":      "div.title_ar",
		"business_unit_title_en": "bu.title_en",
		"business_unit_title_ar": "bu.title_ar",
		"status":                 "status.title_en",
		"risk_champion_title_en": "champion.first_name_en",
		"risk_champion_title_ar": "champion.first_name_ar",
		"start_date":             "risk_register_update_plans.start_date",
		"end_date":               "risk_register_update_plans.end_date",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	response, err := rrupc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (rrupc *RiskRegisterUpdatePlanController) GetRiskRegisterUpdatePlanDetailsByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	plan, err := rrupc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, plan)
}

func (rrupc *RiskRegisterUpdatePlanController) CreateRiskRegisterUpdatePlan(c *gin.Context) {
	var formDTO riskDTOs.RiskRegisterUpdatePlanFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	planID, err := rrupc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: planID})
}

func (rrupc *RiskRegisterUpdatePlanController) UpdateRiskRegisterUpdatePlan(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.RiskRegisterUpdatePlanFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskRegisterUpdatePlanController) DeleteRiskRegisterUpdatePlan(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	// Call service
	failedIDs, err := rrupc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}

func (rrupc *RiskRegisterUpdatePlanController) ForceDeleteRiskRegisterUpdatePlan(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	// Bind JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get current user
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	// Call service
	failedIDs, err := rrupc.service.ForceDelete(request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}
