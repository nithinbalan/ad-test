package controllers

import (
	"jethur/grc/audit-insights/services"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanStatusUpdateController struct {
	service services.AuditInsightsFindingActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFindingActionPlanStatusUpdateController(service services.AuditInsightsFindingActionPlanStatusUpdateServiceInterface) *AuditInsightsFindingActionPlanStatusUpdateController {
	return &AuditInsightsFindingActionPlanStatusUpdateController{
		service: service,
	}
}

func (aifsuc *AuditInsightsFindingActionPlanStatusUpdateController) GetAllAuditInsightsFindingActionPlanStatusUpdates(c *gin.Context) {
	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	// Extract `findingId` if available
	findingIDStr := c.Param("findingId")
	var findingID uint
	if findingIDStr != "" {
		// If `findingId` exists, parse it
		findingIDParsed, err := strconv.ParseUint(findingIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
		findingID = uint(findingIDParsed)
	} else {
		// If `findingId` does not exist in route, get it from the service
		findingID, err = aifsuc.service.GetFindingIDByActionPlan(uint(actionPlanID))
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	}
	auditFindingStatuses, err := aifsuc.service.GetAll(findingID, uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, auditFindingStatuses)
}

func (aifsuc *AuditInsightsFindingActionPlanStatusUpdateController) CreateAuditInsightsFindingActionPlanStatusUpdate(c *gin.Context) {
	var auditFindingStatusUpdateFormDTO auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Extract `findingId` if available
	findingIDStr := c.Param("findingId")
	var findingID uint
	if findingIDStr != "" {
		// If `findingId` exists, parse it
		findingIDParsed, err := strconv.ParseUint(findingIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
		findingID = uint(findingIDParsed)
	} else {
		// If `findingId` does not exist in route, get it from the service
		findingID, err = aifsuc.service.GetFindingIDByActionPlan(uint(actionPlanID))
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	}
	if err := c.ShouldBindJSON(&auditFindingStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifsuc.db, c, err.Error())
		return
	}

	createdID, err := aifsuc.service.Create(findingID, uint(actionPlanID), currentUserID, &auditFindingStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (aifsuc *AuditInsightsFindingActionPlanStatusUpdateController) UpdateAuditInsightsFindingActionPlanStatusUpdate(c *gin.Context) {
	var auditFindingStatusUpdateFormDTO auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO

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

	// Extract `findingId` if available
	findingIDStr := c.Param("findingId")
	var findingID uint
	if findingIDStr != "" {
		// If `findingId` exists, parse it
		findingIDParsed, err := strconv.ParseUint(findingIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
		findingID = uint(findingIDParsed)
	} else {
		// If `findingId` does not exist in route, get it from the service
		findingID, err = aifsuc.service.GetFindingIDByActionPlan(uint(actionPlanID))
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	}
	if err := c.ShouldBindJSON(&auditFindingStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(aifsuc.db, c, err.Error())
		return
	}

	err = aifsuc.service.Update(findingID, uint(actionPlanID), uint(id), currentUserID, &auditFindingStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(aifsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (aifsuc *AuditInsightsFindingActionPlanStatusUpdateController) GetLatestAuditInsightsFindingActionPlanStatusUpdate(c *gin.Context) {

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Extract `findingId` if available
	findingIDStr := c.Param("findingId")
	var findingID uint
	if findingIDStr != "" {
		// If `findingId` exists, parse it
		findingIDParsed, err := strconv.ParseUint(findingIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
		findingID = uint(findingIDParsed)
	} else {
		// If `findingId` does not exist in route, get it from the service
		findingID, err = aifsuc.service.GetFindingIDByActionPlan(uint(actionPlanID))
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	}
	auditInsightsFindingActionPlanStatusUpdate, err := aifsuc.service.GetLatestAuditInsightsFindingActionPlanStatusUpdate(findingID, uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(aifsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, auditInsightsFindingActionPlanStatusUpdate)
}

func (aifsuc *AuditInsightsFindingActionPlanStatusUpdateController) DownloadAuditInsightsFindingActionPlanStatusUpdateDocument(c *gin.Context) {
	actionPlanIDStr := c.Param("actionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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
	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := aifsuc.service.GetDocumentPath(uint(findingID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(aifsuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}
