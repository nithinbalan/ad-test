package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/grc/project/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectSecurityGateReportController struct {
	service services.ProjectSecurityGateReportServiceInterface
	db      *gorm.DB
}

func NewProjectSecurityGateReportController(service services.ProjectSecurityGateReportServiceInterface) *ProjectSecurityGateReportController {
	return &ProjectSecurityGateReportController{
		service: service,
	}
}

func (psgrc *ProjectSecurityGateReportController) GetProjectSecurityGateReportExportZip(c *gin.Context) {
	securityGateIDStr := c.Param("id")
	securityGateID, err := strconv.ParseUint(securityGateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	language := c.Query("language")
	if language != "ar" {
		language = "en"
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	zipFilePath, err := psgrc.service.GetProjectSecurityGateReportExportZip(uint(securityGateID), uint(currentUserID), language)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, fmt.Sprintf("failed to export ZIP: %v", err))
			return
		}
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(zipFilePath)))
	c.File(zipFilePath)

	defer os.Remove(zipFilePath)
}

func (psgrc *ProjectSecurityGateReportController) ExportProjectSecurityGateExceptionReportToPDF(c *gin.Context) {
	language := c.Query("language")
	if language != "ar" {
		language = "en"
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(psgrc.db, c, err.Error())
		return
	}

	fileName, err := psgrc.service.ExportProjectSecurityGateExceptionReport(uint(currentUserID), language)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(psgrc.db, c, fmt.Sprintf("Failed to export Project Security Gate Exception Report PDF: %s", err.Error()))
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
