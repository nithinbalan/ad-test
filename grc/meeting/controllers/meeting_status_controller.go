package controllers

import (
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingStatusController struct {
	service services.MeetingStatusServiceInterface
	db      *gorm.DB
}

func NewMeetingStatusController(service services.MeetingStatusServiceInterface) *MeetingStatusController {
	return &MeetingStatusController{
		service: service,
	}
}

func (msc *MeetingStatusController) GetAllMeetingStatuses(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	meetingStatuses, err := msc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(msc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetingStatuses)
}

func (msc *MeetingStatusController) GetMeetingStatusDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingStatus, err := msc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, meetingStatus)
}
