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

type AuditScheduleController struct {
	service services.AuditScheduleServiceInterface
	db      *gorm.DB
}

func NewAuditScheduleController(service services.AuditScheduleServiceInterface) *AuditScheduleController {
	return &AuditScheduleController{
		service: service,
	}
}

func (asc *AuditScheduleController) GetAllAuditSchedules(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"audit_ids",
		"audit_activity_type_ids",
		"schedule_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	targetDate := c.Query("target_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := dtos.AuditScheduleListQueryParamsDTO{
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

	auditSchedules, err := asc.service.GetAll(uint(auditID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditSchedules)
}

func (asc *AuditScheduleController) GetAuditScheduleDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditSchedule, err := asc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(asc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditSchedule)
}

func (asc *AuditScheduleController) CreateAuditSchedule(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditScheduleFormDTO dtos.AuditScheduleFormDTO

	if err := c.ShouldBindJSON(&auditScheduleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	id, err := asc.service.Create(uint(auditID), currentUserID, &auditScheduleFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(asc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (asc *AuditScheduleController) UpdateAuditSchedule(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditScheduleFormDTO dtos.AuditScheduleFormDTO

	auditScheduleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditScheduleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	err = asc.service.Update(uint(auditScheduleID), uint(auditID), currentUserID, &auditScheduleFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(asc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (asc *AuditScheduleController) DeleteAuditSchedules(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	failedIDs, err := asc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(asc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (asc *AuditScheduleController) GetScheduleSummary(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := asc.service.GetScheduleSummary(uint(auditId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(asc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
