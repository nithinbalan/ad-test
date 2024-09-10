package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditMeetingController struct {
	service services.AuditMeetingServiceInterface
	db      *gorm.DB
}

func NewAuditMeetingController(service services.AuditMeetingServiceInterface) *AuditMeetingController {
	return &AuditMeetingController{
		service: service,
	}
}

func (amc *AuditMeetingController) GetAllAuditMeetings(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"audits_ids",
		"created_by_ids",
		"updated_by_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := auditDTOs.AuditMeetingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	idStr := c.Param("auditId")

	auditId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditMeetings, err := amc.service.GetAll(uint(auditId), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, auditMeetings)
}

func (amc *AuditMeetingController) GetAuditMeetingDetails(c *gin.Context) {
	idStr := c.Param("meetingId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditMeeting, err := amc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(amc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditMeeting)
}

func (amc *AuditMeetingController) CreateAuditMeeting(c *gin.Context) {
	var auditMeetingFormDTO auditDTOs.AuditMeetingFormDTO

	if err := c.ShouldBindJSON(&auditMeetingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	id, err := amc.service.Create(currentUserID, &auditMeetingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(amc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (amc *AuditMeetingController) UpdateAuditMeeting(c *gin.Context) {
	var auditMeetingFormDTO auditDTOs.AuditMeetingFormDTO

	meetingID, err := strconv.Atoi(c.Param("meetingId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditMeetingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	err = amc.service.Update(uint(meetingID), currentUserID, &auditMeetingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(amc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (amc *AuditMeetingController) DeleteAuditMeetings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	failedIDs, err := amc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(amc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
