package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	stakeholderDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderActionPlanStatusUpdateController struct {
	service services.StakeholderActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewStakeholderActionPlanStatusUpdateController(service services.StakeholderActionPlanStatusUpdateServiceInterface) *StakeholderActionPlanStatusUpdateController {
	return &StakeholderActionPlanStatusUpdateController{
		service: service,
	}
}

func (casuc *StakeholderActionPlanStatusUpdateController) GetAllStakeholderActionPlanStatusUpdates(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
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

	stakeholderStatuses, err := casuc.service.GetAll(uint(stakeholderID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, stakeholderStatuses)
}

func (casuc *StakeholderActionPlanStatusUpdateController) CreateStakeholderActionPlanStatusUpdate(c *gin.Context) {
	var stakeholderStatusUpdateFormDTO stakeholderDTOs.StakeholderActionPlanStatusUpdateFormDTO

	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&stakeholderStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	createdID, err := casuc.service.Create(uint(stakeholderID), uint(actionPlanID), currentUserID, &stakeholderStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (casuc *StakeholderActionPlanStatusUpdateController) UpdateStakeholderActionPlanStatusUpdate(c *gin.Context) {
	var stakeholderStatusUpdateFormDTO stakeholderDTOs.StakeholderActionPlanStatusUpdateFormDTO

	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 64)
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

	if err := c.ShouldBindJSON(&stakeholderStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	err = casuc.service.Update(uint(stakeholderID), uint(actionPlanID), uint(id), currentUserID, &stakeholderStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (casuc *StakeholderActionPlanStatusUpdateController) GetLatestStakeholderActionPlanStatusUpdate(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 64)
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

	stakeholderActionPlanStatusUpdate, err := casuc.service.GetLatestStakeholderActionPlanStatusUpdate(uint(stakeholderID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, stakeholderActionPlanStatusUpdate)
}

func (casuc *StakeholderActionPlanStatusUpdateController) GetStakeholderActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 64)
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

	stakeholderActionPlanStatusUpdateDocumentPath, err := casuc.service.GetDocumentPath(uint(stakeholderID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(casuc.db, c, err.Error())
		}
		return
	}
	c.File(stakeholderActionPlanStatusUpdateDocumentPath)
}

func (casuc *StakeholderActionPlanStatusUpdateController) DownloadStakeholderActionPlanStatusUpdateDocument(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 64)
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

	stakeholderActionPlanStatusUpdateDocumentPath, err := casuc.service.GetDocumentPath(uint(stakeholderID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(casuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(stakeholderActionPlanStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(stakeholderActionPlanStatusUpdateDocumentPath)
}

func (casuc *StakeholderActionPlanStatusUpdateController) ApproveStakeholderActionPlan(c *gin.Context) {
	var stakeholderStatusReviewUserUpdateFormDTO stakeholderDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO

	stakeholderID, err := strconv.Atoi(c.Param("stakeholderId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stakeholderActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&stakeholderStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	err = casuc.service.Approve(uint(stakeholderID), uint(stakeholderActionPlanID), currentUserID, &stakeholderStatusReviewUserUpdateFormDTO)
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
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (casuc *StakeholderActionPlanStatusUpdateController) RevertStakeholderActionPlan(c *gin.Context) {
	var stakeholderStatusReviewUserUpdateFormDTO stakeholderDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO

	stakeholderID, err := strconv.Atoi(c.Param("stakeholderId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	stakeholderActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&stakeholderStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(casuc.db, c, err.Error())
		return
	}

	err = casuc.service.Revert(uint(stakeholderID), uint(stakeholderActionPlanID), currentUserID, &stakeholderStatusReviewUserUpdateFormDTO)
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
			utility.SendInternalErrResp(casuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
