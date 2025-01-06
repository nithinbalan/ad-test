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

type MeetingController struct {
	service services.MeetingServiceInterface
	db      *gorm.DB
}

func NewMeetingController(service services.MeetingServiceInterface) *MeetingController {
	return &MeetingController{
		service: service,
	}
}

func (mc *MeetingController) GetAllMeetings(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"meeting_type_ids",
		"meeting_frequency_ids",
		"organizer_ids",
		"meeting_status_ids",
	}

	validSortFields := map[string]string{
		"reference_number":  "meetings.reference_number",
		"title":             "meetings.title",
		"meeting_type_en":   "meeting_types.title_en",
		"meeting_type_ar":   "meeting_types.title_ar",
		"date":              "meetings.date",
		"time":              "meetings.start_time",
		"location":          "meetings.location",
		"organizer_user_en": "users.first_name_en",
		"organizer_user_ar": "users.first_name_ar",
		"status_en":         " meeting_statuses.title_en",
		"status_ar":         " meeting_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := meetingDTOs.MeetingListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		Export:      export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	meetings, err := mc.service.GetAll(currentUserID, &queryParams)

	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(mc.db, c, err.Error())
			return
		}

		excelFile, err := mc.service.ExportMeeting(meetings.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(mc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Meetings.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (mc *MeetingController) GetMeetingDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meeting, err := mc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meeting)
}

func (mc *MeetingController) CreateMeeting(c *gin.Context) {
	var meetingFormDTO meetingDTOs.MeetingFormDTO

	if err := c.ShouldBindJSON(&meetingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	id, err := mc.service.Create(currentUserID, &meetingFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (mc *MeetingController) UpdateMeeting(c *gin.Context) {
	var meetingFormDTO meetingDTOs.MeetingFormDTO

	meetingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&meetingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	err = mc.service.Update(uint(meetingID), currentUserID, &meetingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (mc *MeetingController) DeleteMeetings(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	failedIDs, err := mc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (mc *MeetingController) GetAllCounts(c *gin.Context) {
	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	associatedCounts, err := mc.service.GetAllCounts(uint(meetingID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, associatedCounts)
}

func (mc *MeetingController) GetMeetingEffectiveness(c *gin.Context) {
	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	effectiveness, err := mc.service.GetMeetingEffectiveness(uint(meetingID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, effectiveness)
}

func (mc *MeetingController) GetMeetingExport(c *gin.Context) {
	c.Request.URL.RawQuery += "&export=true"

	mc.GetAllMeetings(c)
}
