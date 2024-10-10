package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
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
