package controllers

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFinalReportSectionController struct {
	service services.AuditFinalReportSectionServiceInterface
	db      *gorm.DB
}

func NewAuditFinalReportSectionController(service services.AuditFinalReportSectionServiceInterface) *AuditFinalReportSectionController {
	return &AuditFinalReportSectionController{
		service: service,
	}
}

func (afrsc *AuditFinalReportSectionController) GetAuditFinalReportSectionDetails(c *gin.Context) {
	_, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sectionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFinalReportSection, err := afrsc.service.GetByID(uint(sectionID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afrsc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFinalReportSection)
}

func (afrsc *AuditFinalReportSectionController) UpdateAuditFinalReportSection(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditFinalReportSectionFormDTO dtos.AuditFinalReportSectionFormDTO

	sectionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFinalReportSectionFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrsc.db, c, err.Error())
		return
	}

	err = afrsc.service.Update(uint(auditID), uint(sectionID), currentUserID, &auditFinalReportSectionFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afrsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afrsc *AuditFinalReportSectionController) DeleteAuditFinalReportSection(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrsc.db, c, err.Error())
		return
	}

	failedID, err := afrsc.service.Delete(uint(auditID), uint(id), currentUserID)
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
			utility.SendInternalErrResp(afrsc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
