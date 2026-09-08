package controllers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectReportController struct {
	service services.ProjectReportServiceInterface
	db      *gorm.DB
}

func NewProjectReportController(service services.ProjectReportServiceInterface) *ProjectReportController {
	return &ProjectReportController{
		service: service,
	}
}

func (prc *ProjectReportController) ExportProjectReportToPDF(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	if projectIDStr == "" {
		projectIDStr = c.Param("id")
	}
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	fileName, err := prc.service.ExportProjectReport(uint(projectID), uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(prc.db, c, fmt.Sprintf("Failed to export Project Report PDF: %s", err.Error()))
			return
		}
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", "application/pdf")

	c.File(fileName)
	c.Writer.Flush()

	go func() {
		time.Sleep(2 * time.Second)
		err := os.Remove(fileName)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", fileName, err)
		} else {
			fmt.Printf("Successfully deleted file %s\n", fileName)
		}
	}()
}
