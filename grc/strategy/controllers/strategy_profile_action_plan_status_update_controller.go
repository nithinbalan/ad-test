package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategyProfileActionPlanStatusUpdateController struct {
	service services.StrategyProfileActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileActionPlanStatusUpdateController(service services.StrategyProfileActionPlanStatusUpdateServiceInterface) *StrategyProfileActionPlanStatusUpdateController {
	return &StrategyProfileActionPlanStatusUpdateController{
		service: service,
	}
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) GetAllStrategyProfileActionPlanStatusUpdates(c *gin.Context) {
	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 32)
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

	slaStatuses, err := rapsuc.service.GetAll(uint(strategyProfileIdID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, slaStatuses)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) CreateStrategyProfileActionPlanStatusUpdate(c *gin.Context) {
	var taeRecoveryStatusUpdateFormDTO strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO

	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&taeRecoveryStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		return
	}

	createdID, err := rapsuc.service.Create(uint(strategyProfileIdID), uint(actionPlanID), currentUserID, &taeRecoveryStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) UpdateStrategyProfileActionPlanStatusUpdate(c *gin.Context) {
	var taeRecoveryStatusUpdateFormDTO strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO

	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&taeRecoveryStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		return
	}

	err = rapsuc.service.Update(uint(strategyProfileIdID), uint(actionPlanID), uint(id), currentUserID, &taeRecoveryStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) GetLatestStrategyProfileActionPlanStatusUpdate(c *gin.Context) {
	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 64)
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

	slaActionPlanStatusUpdate, err := rapsuc.service.GetLatestStrategyProfileActionPlanStatusUpdate(uint(strategyProfileIdID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, slaActionPlanStatusUpdate)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) GetStrategyProfileActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 64)
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

	slaActionPlanStatusUpdateDocumentPath, err := rapsuc.service.GetDocumentPath(uint(strategyProfileIdID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		}
		return
	}
	c.File(slaActionPlanStatusUpdateDocumentPath)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) DownloadStrategyProfileActionPlanStatusUpdateDocument(c *gin.Context) {
	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 64)
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

	slaActionPlanStatusUpdateDocumentPath, err := rapsuc.service.GetDocumentPath(uint(strategyProfileIdID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(slaActionPlanStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(slaActionPlanStatusUpdateDocumentPath)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) ApproveStrategyProfileActionPlan(c *gin.Context) {
	var recoveryReviewUserUpdateFormDTO strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO

	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&recoveryReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		return
	}

	err = rapsuc.service.Approve(uint(strategyProfileIdID), uint(actionPlanID), currentUserID, &recoveryReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rapsuc *StrategyProfileActionPlanStatusUpdateController) RevertStrategyProfileActionPlan(c *gin.Context) {
	var recoveryReviewUserUpdateFormDTO strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO

	strategyProfileIdIDStr := c.Param("strategyProfileId")

	strategyProfileIdID, err := strconv.ParseUint(strategyProfileIdIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&recoveryReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rapsuc.db, c, err.Error())
		return
	}

	err = rapsuc.service.Revert(uint(strategyProfileIdID), uint(actionPlanID), currentUserID, &recoveryReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(rapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
