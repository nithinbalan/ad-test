package controllers

import (
	"jethur/grc/meeting/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingTypeController struct {
	service services.MeetingTypeServiceInterface
	db      *gorm.DB
}

func NewMeetingTypeController(service services.MeetingTypeServiceInterface) *MeetingTypeController {
	return &MeetingTypeController{
		service: service,
	}
}

func (mtc *MeetingTypeController) GetAllMeetingTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	meetingTypes, err := mtc.service.GetAll(sort, sortBy, keyword, status)

	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, meetingTypes)
}

func (mtc *MeetingTypeController) GetMeetingTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingType, err := mtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, meetingType)
}

func (mtc *MeetingTypeController) CreateMeetingType(c *gin.Context) {
	var meetingType *models.MeetingType

	if err := c.ShouldBindJSON(&meetingType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	createdType, err := mtc.service.Create(currentUserID, meetingType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdType)
}

func (mtc *MeetingTypeController) UpdateMeetingType(c *gin.Context) {
	var updatedType *models.MeetingType

	meetingTypeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedType); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	updatedType, err = mtc.service.Update(uint(meetingTypeID), currentUserID, updatedType)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mtc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedType)
}

func (mtc *MeetingTypeController) DeleteMeetingTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	failedIDs, err := mtc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (mtc *MeetingTypeController) ActivateMeetingTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	err = mtc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (mtc *MeetingTypeController) DeactivateMeetingTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	err = mtc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(mtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
