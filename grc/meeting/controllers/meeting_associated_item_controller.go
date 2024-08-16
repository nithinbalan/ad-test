package controllers

import (
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingAssociatedItemController struct {
	service services.MeetingAssociatedItemServiceInterface
	db      *gorm.DB
}

func NewMeetingAssociatedItemController(service services.MeetingAssociatedItemServiceInterface) *MeetingAssociatedItemController {
	return &MeetingAssociatedItemController{
		service: service,
	}
}

func (maic *MeetingAssociatedItemController) GetAllMeetingControls(c *gin.Context) {

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

	queryParams := meetingDTOs.MeetingAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	meetings, err := maic.service.GetAllControls(uint(meetingID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(maic.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingComplianceRequirements(c *gin.Context) {

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

	queryParams := meetingDTOs.MeetingAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	meetings, err := maic.service.GetAllComplianceRequirements(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingRisks(c *gin.Context) {

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

	queryParams := meetingDTOs.MeetingAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	meetings, err := maic.service.GetAllRisks(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func (maic *MeetingAssociatedItemController) GetAllMeetingProcesses(c *gin.Context) {

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

	queryParams := meetingDTOs.MeetingAssociatedItemListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}
	meetings, err := maic.service.GetAllProcesses(uint(meetingID), &queryParams)

	if err != nil {
		utility.SendInternalErrResp(maic.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetings)
}
