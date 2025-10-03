package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"jethur/grc/meeting/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeetingReportController struct {
	service services.MeetingReportServiceInterface
	db      *gorm.DB
}

func NewMeetingReportController(service services.MeetingReportServiceInterface) *MeetingReportController {
	return &MeetingReportController{
		service: service,
	}
}

func (irc *MeetingReportController) ExportMeetingReportToPDF(c *gin.Context) {
	meetingIDStr := c.Param("id")
	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irc.db, c, err.Error())
		return
	}

	// Delegate to service to handle PDF generation
	fileName, err := irc.service.ExportMeetingReport(uint(meetingID), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irc.db, c, fmt.Sprintf("Failed to export Meeting Report PDF: %s", err.Error()))
			return
		}
	}

	// Set headers to ensure the file is served with the correct file name
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	// Serve the PDF file
	c.File(fileName)
	// Cleanup after the response has been sent
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second) // Add a small delay to ensure file is no longer in use
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}
