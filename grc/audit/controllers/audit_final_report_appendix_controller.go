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
	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFinalReportAppendix, err := afrac.service.GetByID(uint(auditID), uint(finalReportID), uint(sectionID), uint(id))
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
	sectionID, err := strconv.Atoi(c.Param("sectionId"))
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

	_, err = afrac.service.Create(uint(auditID), uint(finalReportID), uint(sectionID), currentUserID, &auditFinalReportAppendixFormDTO)
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
func (afrac *AuditFinalReportAppendixController) UpdateAuditFinalReportAppendix(c *gin.Context) {
	auditID, err := strconv.Atoi(c.Param("auditId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var auditFinalReportAppendixFormDTO dtos.AuditFinalReportAppendixFormDTO

	finalReportID, err := strconv.Atoi(c.Param("finalReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFinalReportAppendixFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrac.db, c, err.Error())
		return
	}

	err = afrac.service.Update(uint(auditID), uint(finalReportID), uint(sectionID), uint(id), currentUserID, &auditFinalReportAppendixFormDTO)
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
	sectionID, err := strconv.Atoi(c.Param("sectionId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	var request struct {
		IDs []uint `json:"ids"`
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afrac.db, c, err.Error())
		return
	}

	failedID, err := afrac.service.Delete(uint(auditID), uint(finalReportID), uint(sectionID), request.IDs, currentUserID)
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

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
