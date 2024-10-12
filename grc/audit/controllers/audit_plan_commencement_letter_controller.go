package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditPlanCommencementLetterController struct {
	service services.AuditPlanCommencementLetterServiceInterface
	db      *gorm.DB
}

func NewAuditPlanCommencementLetterController(service services.AuditPlanCommencementLetterServiceInterface) *AuditPlanCommencementLetterController {
	return &AuditPlanCommencementLetterController{
		service: service,
	}
}

func (apclc *AuditPlanCommencementLetterController) GetAllAuditPlanCommencementLetters(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	keyword := c.Query("keyword")

	commencementLetters, err := apclc.service.GetAll(uint(auditPlanID), keyword)
	if err != nil {
		utility.SendInternalErrResp(apclc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, commencementLetters)
}

func (apclc *AuditPlanCommencementLetterController) GetAuditPlanCommencementLetterDetails(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commencementLetterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commencementLetter, err := apclc.service.GetFullDetailsByID(uint(auditPlanID), uint(commencementLetterID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apclc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, commencementLetter)
}

func (apc *AuditPlanCommencementLetterController) CreateAuditPlanCommencementLetters(c *gin.Context) {

	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	err = apc.service.Create(uint(auditPlanID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apclc *AuditPlanCommencementLetterController) UpdateAuditPlanCommencementLetter(c *gin.Context) {
	auditPlanID, err := strconv.Atoi(c.Param("planId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commencementLetterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var commencementLetterFormDTO auditDTOs.AuditPlanCommencementLetterFormDTO

	if err := c.ShouldBindJSON(&commencementLetterFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apclc.db, c, err.Error())
		return
	}

	err = apclc.service.Update(uint(auditPlanID), uint(commencementLetterID), currentUserID, &commencementLetterFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apclc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apclc *AuditPlanCommencementLetterController) DeleteAuditPlanCommencementLetter(c *gin.Context) {
	auditPlanIDStr := c.Param("planId")

	auditPlanID, err := strconv.ParseUint(auditPlanIDStr, 10, 64)
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
		utility.SendInternalErrResp(apclc.db, c, err.Error())
		return
	}

	failedID, err := apclc.service.Delete(uint(auditPlanID), uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(apclc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
