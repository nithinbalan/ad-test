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

type MeetingAgendaController struct {
	service services.MeetingAgendaServiceInterface
	db      *gorm.DB
}

func NewMeetingAgendaController(service services.MeetingAgendaServiceInterface) *MeetingAgendaController {
	return &MeetingAgendaController{
		service: service,
	}
}

func (mac *MeetingAgendaController) GetAllMeetingAgendas(c *gin.Context) {
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

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := meetingDTOs.MeetingAgendaListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}

	meetingAgendas, err := mac.service.GetAll(uint(meetingID), currentUserID, &queryParams)

	if err != nil {
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetingAgendas)
}

func (mac *MeetingAgendaController) GetMeetingAgendaDetails(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("agendaID")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingAgenda, err := mac.service.GetFullDetailsByID(uint(meetingID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetingAgenda)
}

func (mac *MeetingAgendaController) CreateMeetingAgenda(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var meetingAgendaFormDTO meetingDTOs.MeetingAgendaFormDTO

	if err := c.ShouldBindJSON(&meetingAgendaFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}

	id, err := mac.service.Create(uint(meetingID), currentUserID, &meetingAgendaFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (mac *MeetingAgendaController) UpdateMeetingAgenda(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var meetingAgendaFormDTO meetingDTOs.MeetingAgendaFormDTO

	meetingAgendaID, err := strconv.Atoi(c.Param("agendaID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&meetingAgendaFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}

	err = mac.service.Update(uint(meetingID), uint(meetingAgendaID), currentUserID, &meetingAgendaFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mac.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (mac *MeetingAgendaController) DeleteMeetingAgendas(c *gin.Context) {
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
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}

	failedIDs, err := mac.service.Delete(uint(meetingID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
