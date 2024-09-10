package controllers

import (
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingMinuteStatusController struct {
	service services.MeetingMinuteStatusServiceInterface
	db      *gorm.DB
}

func NewMeetingMinuteStatusController(service services.MeetingMinuteStatusServiceInterface) *MeetingMinuteStatusController {
	return &MeetingMinuteStatusController{
		service: service,
	}
}

func (mmsc *MeetingMinuteStatusController) GetAllMeetingMinuteStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	meetingMinuteStatuses, err := mmsc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(mmsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetingMinuteStatuses)
}

func (mmsc *MeetingMinuteStatusController) GetMeetingMinuteStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingMinuteStatus, err := mmsc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, meetingMinuteStatus)
}
