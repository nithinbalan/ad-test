package controllers

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaeRecoveryActionPlanStatusUpdateController struct {
	service services.TaeRecoveryActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewTaeRecoveryActionPlanStatusUpdateController(service services.TaeRecoveryActionPlanStatusUpdateServiceInterface) *TaeRecoveryActionPlanStatusUpdateController {
	return &TaeRecoveryActionPlanStatusUpdateController{
		service: service,
	}
}

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) GetAllTaeRecoveryActionPlanStatusUpdates(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
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

	slaStatuses, err := rapsuc.service.GetAll(uint(testAndExerciseID), uint(actionPlanID), sort, sortBy, keyword)

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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) CreateTaeRecoveryActionPlanStatusUpdate(c *gin.Context) {
	var taeRecoveryStatusUpdateFormDTO businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO

	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
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

	createdID, err := rapsuc.service.Create(uint(testAndExerciseID), uint(actionPlanID), currentUserID, &taeRecoveryStatusUpdateFormDTO)

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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) UpdateTaeRecoveryActionPlanStatusUpdate(c *gin.Context) {
	var taeRecoveryStatusUpdateFormDTO businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO

	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
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

	err = rapsuc.service.Update(uint(testAndExerciseID), uint(actionPlanID), uint(id), currentUserID, &taeRecoveryStatusUpdateFormDTO)
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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) GetLatestTaeRecoveryActionPlanStatusUpdate(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
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

	slaActionPlanStatusUpdate, err := rapsuc.service.GetLatestTaeRecoveryActionPlanStatusUpdate(uint(testAndExerciseID), uint(actionPlanID))
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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) GetTaeRecoveryActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
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

	slaActionPlanStatusUpdateDocumentPath, err := rapsuc.service.GetDocumentPath(uint(testAndExerciseID), uint(actionPlanID), uint(id), token)
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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) DownloadTaeRecoveryActionPlanStatusUpdateDocument(c *gin.Context) {
	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 64)
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

	slaActionPlanStatusUpdateDocumentPath, err := rapsuc.service.GetDocumentPath(uint(testAndExerciseID), uint(actionPlanID), uint(id), token)
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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) ApproveTaeRecoveryActionPlan(c *gin.Context) {
	var recoveryReviewUserUpdateFormDTO businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO

	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
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

	err = rapsuc.service.Approve(uint(testAndExerciseID), uint(actionPlanID), currentUserID, &recoveryReviewUserUpdateFormDTO)
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

func (rapsuc *TaeRecoveryActionPlanStatusUpdateController) RevertTaeRecoveryActionPlan(c *gin.Context) {
	var recoveryReviewUserUpdateFormDTO businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO

	testAndExerciseIDStr := c.Param("testAndExerciseId")

	testAndExerciseID, err := strconv.ParseUint(testAndExerciseIDStr, 10, 32)
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

	err = rapsuc.service.Revert(uint(testAndExerciseID), uint(actionPlanID), currentUserID, &recoveryReviewUserUpdateFormDTO)
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
