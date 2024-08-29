package controllers

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditFindingActionPlanStatusUpdateController struct {
	service services.AuditFindingActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewAuditFindingActionPlanStatusUpdateController(service services.AuditFindingActionPlanStatusUpdateServiceInterface) *AuditFindingActionPlanStatusUpdateController {
	return &AuditFindingActionPlanStatusUpdateController{
		service: service,
	}
}

func (afsuc *AuditFindingActionPlanStatusUpdateController) GetAllAuditFindingActionPlanStatusUpdates(c *gin.Context) {
	auditFindingIDStr := c.Param("findingId")

	auditFindingID, err := strconv.ParseUint(auditFindingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	fmt.Println(auditFindingID)

	auditFindingStatuses, err := afsuc.service.GetAll(uint(auditFindingID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFindingStatuses)
}

func (afsuc *AuditFindingActionPlanStatusUpdateController) CreateAuditFindingActionPlanStatusUpdate(c *gin.Context) {
	var auditFindingStatusUpdateFormDTO auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO

	auditFindingIDStr := c.Param("findingId")

	auditFindingID, err := strconv.ParseUint(auditFindingIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&auditFindingStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afsuc.db, c, err.Error())
		return
	}

	createdID, err := afsuc.service.Create(uint(auditFindingID), uint(actionPlanID), currentUserID, &auditFindingStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (afsuc *AuditFindingActionPlanStatusUpdateController) UpdateAuditFindingActionPlanStatusUpdate(c *gin.Context) {
	var auditFindingStatusUpdateFormDTO auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO

	auditFindingIDStr := c.Param("findingId")

	auditFindingID, err := strconv.ParseUint(auditFindingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&auditFindingStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afsuc.db, c, err.Error())
		return
	}

	err = afsuc.service.Update(uint(auditFindingID), uint(actionPlanID), uint(id), currentUserID, &auditFindingStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afsuc *AuditFindingActionPlanStatusUpdateController) GetLatestAuditFindingActionPlanStatusUpdate(c *gin.Context) {
	auditFindingIDStr := c.Param("findingId")

	auditFindingID, err := strconv.ParseUint(auditFindingIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	auditFindingActionPlanStatusUpdate, err := afsuc.service.GetLatestAuditFindingActionPlanStatusUpdate(uint(auditFindingID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, auditFindingActionPlanStatusUpdate)
}
