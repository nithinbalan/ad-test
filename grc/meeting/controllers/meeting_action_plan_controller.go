package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingActionPlanController struct {
	service services.MeetingActionPlanServiceInterface
	db      *gorm.DB
}

func NewMeetingActionPlanController(service services.MeetingActionPlanServiceInterface) *MeetingActionPlanController {
	return &MeetingActionPlanController{
		service: service,
	}
}

func (mapc *MeetingActionPlanController) GetAllMeetingActionPlans(c *gin.Context) {
	meetingIDStr := c.Param("id")

	var meetingID uint64
	var err error

	// Check if meetingID is provided in the URL
	if meetingIDStr != "" {
		meetingID, err = strconv.ParseUint(meetingIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		meetingID = 0
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(mapc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"severity_level_ids",
		"responsible_user_ids",
		"action_plan_status_ids",
	}

	validSortFields := map[string]string{
		"title":               "meeting_action_plans.title",
		"meeting_title":       "meetings.title",
		"severity_level_en":   "severity_levels.title_en",
		"severity_level_ar":   "severity_levels.title_ar",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"target_date":         "meeting_action_plans.target_date",
		"percentage":          "meeting_action_plans.percentage",
		"status_en":           " action_plan_statuses.title_en",
		"status_ar":           " action_plan_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(mapc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(mapc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := meetingDTOs.MeetingActionPlanListQueryParamsDTO{
		Sort:           sort,
		SortBy:         dbField,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
		Export:         export,
	}
	meetingActionPlans, err := mapc.service.GetAll(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(mapc.db, c, err.Error())
			return
		}

		excelFile, err := mapc.service.ExportMeetingActionPlan(meetingActionPlans.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(mapc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Meeting action plan.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, meetingActionPlans)
}

func (mapc *MeetingActionPlanController) GetMeetingActionPlanDetails(c *gin.Context) {
	meetingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingActionPlanID, err := strconv.Atoi(c.Param("meetingActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingActionPlan, err := mapc.service.GetFullDetailsByID(uint(meetingID), uint(meetingActionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, meetingActionPlan)
}

func (mapc *MeetingActionPlanController) CreateMeetingActionPlan(c *gin.Context) {
	var MeetingActionPlanFormDTO meetingDTOs.MeetingActionPlanFormDTO

	meetingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&MeetingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mapc.db, c, err.Error())
		return
	}

	id, err := mapc.service.Create(uint(meetingID), currentUserID, &MeetingActionPlanFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (mapc *MeetingActionPlanController) UpdateMeetingActionPlan(c *gin.Context) {
	var MeetingActionPlanFormDTO meetingDTOs.MeetingActionPlanFormDTO

	meetingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingActionPlanID, err := strconv.Atoi(c.Param("meetingActionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&MeetingActionPlanFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mapc.db, c, err.Error())
		return
	}

	err = mapc.service.Update(uint(meetingID), uint(meetingActionPlanID), currentUserID, &MeetingActionPlanFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (mapc *MeetingActionPlanController) DeleteMeetingActionPlans(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	meetingID, err := strconv.Atoi(c.Param("id"))
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
		utility.SendInternalErrResp(mapc.db, c, err.Error())
		return
	}

	failedIDs, err := mapc.service.Delete(uint(meetingID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (mapc *MeetingActionPlanController) GetMeetingActionPlanSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := mapc.service.GetMeetingActionPlanSummary(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, summary)
}
func (mapc *MeetingActionPlanController) GetMeetingActionPlanExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	mapc.GetAllMeetingActionPlans(c)
}
