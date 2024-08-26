package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanScheduleController struct {
	service services.AuditPlanScheduleServiceInterface
	db      *gorm.DB
}

func NewAuditPlanScheduleController(service services.AuditPlanScheduleServiceInterface) *AuditPlanScheduleController {
	return &AuditPlanScheduleController{
		service: service,
	}
}

func (apsc *AuditPlanScheduleController) GetAllAuditPlanSchedules(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_plan_ids",
		"audit_activity_type_ids",
		"schedule_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	targetDate := c.Query("target_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditPlanScheduleListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		TargetDate:  targetDate,
	}

	auditPlanSchedules, err := apsc.service.GetAll(uint(auditPlanID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditPlanSchedules)
}

func (apsc *AuditPlanScheduleController) GetAuditPlanScheduleDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditPlanSchedule, err := apsc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apsc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditPlanSchedule)
}

func (apsc *AuditPlanScheduleController) CreateAuditPlanSchedule(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanScheduleFormDTO dtos.AuditPlanScheduleFormDTO

	if err := c.ShouldBindJSON(&auditPlanScheduleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}

	id, err := apsc.service.Create(uint(auditPlanID), currentUserID, &auditPlanScheduleFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (apsc *AuditPlanScheduleController) UpdateAuditPlanSchedule(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditPlanScheduleFormDTO dtos.AuditPlanScheduleFormDTO

	auditPlanScheduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditPlanScheduleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}

	err = apsc.service.Update(uint(auditPlanScheduleID), uint(auditPlanID), currentUserID, &auditPlanScheduleFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apsc *AuditPlanScheduleController) DeleteAuditPlanSchedules(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}

	failedIDs, err := apsc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
