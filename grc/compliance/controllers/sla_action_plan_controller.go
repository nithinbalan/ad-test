package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLAActionPlanController struct {
	service services.SLAActionPlanServiceInterface
	db      *gorm.DB
}

func NewSLAActionPlanController(service services.SLAActionPlanServiceInterface) *SLAActionPlanController {
	return &SLAActionPlanController{
		service: service,
	}
}

func (sapc *SLAActionPlanController) GetAllSLAActionPlans(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
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
		"sla_check_question_ids",
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
		utility.SendInternalErrResp(sapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := complianceDTOs.SLAActionPlanListQueryParamsDTO{
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
	slaActionPlans, err := sapc.service.GetAll(uint(slaID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slaActionPlans)
}

func (sapc *SLAActionPlanController) GetSLAActionPlanDetails(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
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

	slaActionPlan, err := sapc.service.GetFullDetailsByID(uint(slaID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, slaActionPlan)
}

func (sapc *SLAActionPlanController) CreateSLAActionPlan(c *gin.Context) {
	var slaActionPlanFormDTO complianceDTOs.SLAActionPlanFormDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&slaActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	id, err := sapc.service.Create(uint(slaID), currentUserID, &slaActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (sapc *SLAActionPlanController) UpdateSLAActionPlan(c *gin.Context) {
	var slaActionPlanFormDTO complianceDTOs.SLAActionPlanFormDTO

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&slaActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	err = sapc.service.Update(uint(slaID), uint(id), currentUserID, &slaActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sapc *SLAActionPlanController) DeleteSLAActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
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
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	failedIDs, err := sapc.service.Delete(uint(slaID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sapc *SLAActionPlanController) GetSLAActionPlanCountByStatus(c *gin.Context) {
	slaIDStr := c.Param("slaId")

	slaID, err := strconv.ParseUint(slaIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanCounts, err := sapc.service.GetSLAActionPlanCountByStatus(uint(slaID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, actionPlanCounts)
}
