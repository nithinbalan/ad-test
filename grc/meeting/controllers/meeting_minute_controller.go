package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingMinuteController struct {
	service services.MeetingMinuteServiceInterface
	db      *gorm.DB
}

func NewMeetingMinuteController(service services.MeetingMinuteServiceInterface) *MeetingMinuteController {
	return &MeetingMinuteController{
		service: service,
	}
}

func (mmc *MeetingMinuteController) GetAllMeetingMinutes(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
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
		"meeting_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := meetingDTOs.MeetingMinuteListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}

	meetingMinutes, err := mmc.service.GetAll(uint(meetingID), currentUserID, &queryParams)

	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetingMinutes)
}

func (mmc *MeetingMinuteController) GetMeetingMinuteDetails(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("meetingMinuteID")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingMinute, err := mmc.service.GetFullDetailsByID(uint(meetingID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mmc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetingMinute)
}

func (mmc *MeetingMinuteController) CreateMeetingMinute(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var meetingMinuteFormDTO meetingDTOs.MeetingMinuteFormDTO

	if err := c.ShouldBindJSON(&meetingMinuteFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}

	id, err := mmc.service.Create(uint(meetingID), currentUserID, &meetingMinuteFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mmc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (mmc *MeetingMinuteController) UpdateMeetingMinute(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var meetingMinuteFormDTO meetingDTOs.MeetingMinuteFormDTO

	meetingMinuteID, err := strconv.Atoi(c.Param("meetingMinuteID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&meetingMinuteFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}

	err = mmc.service.Update(uint(meetingID), uint(meetingMinuteID), currentUserID, &meetingMinuteFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mmc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (mmc *MeetingMinuteController) DeleteMeetingMinutes(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}

	failedIDs, err := mmc.service.Delete(uint(meetingID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
