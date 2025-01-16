package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFinalReportAppendixController struct {
	service services.AuditFinalReportAppendixServiceInterface
	db      *gorm.DB
}

func NewAuditFinalReportAppendixController(service services.AuditFinalReportAppendixServiceInterface) *AuditFinalReportAppendixController {
	return &AuditFinalReportAppendixController{
		service: service,
	}
}

func (afrac *AuditFinalReportAppendixController) GetAuditFinalReportAppendixDetails(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	finalReportID, err := strconv.Atoi(c.Param("finalReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFinalReportAppendix, err := afrac.service.GetByID(uint(auditID), uint(finalReportID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFinalReportAppendix)
}
func (afrac *AuditFinalReportAppendixController) CreateAuditFinalReportAppendix(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	finalReportID, err := strconv.Atoi(c.Param("finalReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditFinalReportAppendixFormDTO dtos.AuditFinalReportAppendixFormDTO

	if err := c.ShouldBindJSON(&auditFinalReportAppendixFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrac.db, c, err.Error())
		return
	}

	_, err = afrac.service.Create(uint(auditID), uint(finalReportID), currentUserID, &auditFinalReportAppendixFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
func (afrac *AuditFinalReportAppendixController) DeleteAuditFinalReportAppendix(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	finalReportID, err := strconv.Atoi(c.Param("finalReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrac.db, c, err.Error())
		return
	}

	deletedID, err := afrac.service.Delete(uint(auditID), uint(finalReportID), uint(id), currentUserID)
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
			utility.SendInternalErrResp(afrac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"deleted_id": deletedID})
}

func (afrac *AuditFinalReportAppendixController) DownloadFinalReportAppendixDocument(c *gin.Context) {
	auditIdStr := c.Param("auditId")
	auditId, err := strconv.ParseUint(auditIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	finalReportID, err := strconv.Atoi(c.Param("finalReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	appendixDocumentPath, err := afrac.service.GetDocumentPath(uint(auditId), uint(finalReportID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afrac.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(appendixDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(appendixDocumentPath)
}
