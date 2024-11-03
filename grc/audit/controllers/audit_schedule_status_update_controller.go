package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditScheduleStatusUpdateController struct {
	service services.AuditScheduleStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewAuditScheduleStatusUpdateController(service services.AuditScheduleStatusUpdateServiceInterface) *AuditScheduleStatusUpdateController {
	return &AuditScheduleStatusUpdateController{
		service: service,
	}
}

func (assuc *AuditScheduleStatusUpdateController) GetAllAuditScheduleStatusUpdates(c *gin.Context) {
	auditScheduleIDStr := c.Param("scheduleId")
	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	statusUpdates, err := assuc.service.GetAll(uint(auditScheduleID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(assuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, statusUpdates)
}

func (assuc *AuditScheduleStatusUpdateController) CreateAuditScheduleStatusUpdate(c *gin.Context) {
	var formDTO auditDTOs.AuditScheduleStatusUpdateFormDTO
	auditScheduleIDStr := c.Param("scheduleId")
	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(assuc.db, c, err.Error())
		return
	}

	createdID, err := assuc.service.Create(uint(auditScheduleID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(assuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (assuc *AuditScheduleStatusUpdateController) UpdateAuditScheduleStatusUpdate(c *gin.Context) {
	var formDTO auditDTOs.AuditScheduleStatusUpdateFormDTO
	auditScheduleIDStr := c.Param("scheduleId")
	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(assuc.db, c, err.Error())
		return
	}

	err = assuc.service.Update(uint(auditScheduleID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(assuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (assuc *AuditScheduleStatusUpdateController) GetLatestAuditScheduleStatusUpdate(c *gin.Context) {
	auditScheduleIDStr := c.Param("scheduleId")
	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	latestUpdate, err := assuc.service.GetLatestAuditScheduleStatusUpdate(uint(auditScheduleID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(assuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, latestUpdate)
}

func (iasuc *AuditScheduleStatusUpdateController) GetAuditScheduleStatusUpdateDocumentPreview(c *gin.Context) {

	auditIDStr := c.Param("auditId")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditScheduleIDStr := c.Param("scheduleId")

	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditScheduleStatusUpdateDocumentPath, err := iasuc.service.GetDocumentPath(uint(auditID), uint(auditScheduleID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iasuc.db, c, err.Error())
		}
		return
	}
	c.File(auditScheduleStatusUpdateDocumentPath)
}

func (iasuc *AuditScheduleStatusUpdateController) DownloadAuditScheduleStatusUpdateDocument(c *gin.Context) {
	auditIDStr := c.Param("auditId")

	auditID, err := strconv.ParseUint(auditIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditScheduleIDStr := c.Param("scheduleId")

	auditScheduleID, err := strconv.ParseUint(auditScheduleIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditScheduleStatusUpdateDocumentPath, err := iasuc.service.GetDocumentPath(uint(auditID), uint(auditScheduleID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iasuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(auditScheduleStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(auditScheduleStatusUpdateDocumentPath)
}
