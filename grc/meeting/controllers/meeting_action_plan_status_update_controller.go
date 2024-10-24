package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingActionPlanStatusUpdateController struct {
	service services.MeetingActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewMeetingActionPlanStatusUpdateController(service services.MeetingActionPlanStatusUpdateServiceInterface) *MeetingActionPlanStatusUpdateController {
	return &MeetingActionPlanStatusUpdateController{
		service: service,
	}
}

func (mapuc *MeetingActionPlanStatusUpdateController) GetAllMeetingActionPlanStatusUpdates(c *gin.Context) {
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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	meetingActionPlanStatusUpdates, err := mapuc.service.GetAll(uint(meetingID), uint(meetingActionPlanID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetingActionPlanStatusUpdates)
}

func (mapuc *MeetingActionPlanStatusUpdateController) CreateMeetingActionPlanStatusUpdate(c *gin.Context) {
	var meetingActionPlanStatusUpdateFormDTO meetingDTOs.MeetingActionPlanStatusUpdateFormDTO

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

	if err := c.ShouldBindJSON(&meetingActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mapuc.db, c, err.Error())
		return
	}

	createdID, err := mapuc.service.Create(uint(meetingID), uint(meetingActionPlanID), currentUserID, &meetingActionPlanStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (mapuc *MeetingActionPlanStatusUpdateController) UpdateMeetingActionPlanStatusUpdate(c *gin.Context) {
	var meetingActionPlanStatusUpdateFormDTO meetingDTOs.MeetingActionPlanStatusUpdateFormDTO

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

	meetingActionPlanUpdateID, err := strconv.Atoi(c.Param("meetingActionPlanUpdateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&meetingActionPlanStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(mapuc.db, c, err.Error())
		return
	}

	err = mapuc.service.Update(uint(meetingID), uint(meetingActionPlanID), uint(meetingActionPlanUpdateID), currentUserID, &meetingActionPlanStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (mapuc *MeetingActionPlanStatusUpdateController) GetMeetingActionPlanStatusUpdateDetails(c *gin.Context) {
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

	meetingActionPlanUpdateID, err := strconv.Atoi(c.Param("meetingActionPlanUpdateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingTest, err := mapuc.service.GetByID(uint(meetingID), uint(meetingActionPlanID), uint(meetingActionPlanUpdateID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, meetingTest)
}

func (mapuc *MeetingActionPlanStatusUpdateController) GetLatestMeetingActionPlanStatusUpdate(c *gin.Context) {
	meetingIDStr := c.Param("id")

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("meetingActionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	meetingActionPlanStatusUpdate, err := mapuc.service.GetLatestMeetingActionPlanStatusUpdate(uint(meetingID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, meetingActionPlanStatusUpdate)
}

func (mapuc *MeetingActionPlanStatusUpdateController) GetMeetingActionPlanDocumentPreview(c *gin.Context) {
	// Parse the meeting ID from the URL
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse the meeting action plan ID from the URL
	meetingActionPlanIDStr := c.Param("meetingActionPlanId")
	meetingActionPlanID, err := strconv.ParseUint(meetingActionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Parse the meeting action plan update ID from the URL
	meetingActionPlanUpdateIDStr := c.Param("meetingActionPlanUpdateId")
	meetingActionPlanUpdateID, err := strconv.ParseUint(meetingActionPlanUpdateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Extract the token query parameter from the request
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Call the service layer to get the document path
	documentPath, err := mapuc.service.GetDocumentPath(uint(meetingID), uint(meetingActionPlanID), uint(meetingActionPlanUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, err.Error())
			return
		}
	}

	// Serve the file as a response
	c.File(documentPath)
}

func (mapuc *MeetingActionPlanStatusUpdateController) DownloadMeetingActionPlanDocument(c *gin.Context) {
	// meeting ID
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// meeting action plan ID
	meetingActionPlanIDStr := c.Param("meetingActionPlanId")
	meetingActionPlanID, err := strconv.ParseUint(meetingActionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// meeting action plan update ID
	meetingActionPlanUpdateIDStr := c.Param("meetingActionPlanUpdateId")
	meetingActionPlanUpdateID, err := strconv.ParseUint(meetingActionPlanUpdateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Extract the token
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get the document path
	documentPath, err := mapuc.service.GetDocumentPath(uint(meetingID), uint(meetingActionPlanID), uint(meetingActionPlanUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(mapuc.db, c, err.Error())
			return
		}
	}

	// Get the file name
	fileName := filepath.Base(documentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
