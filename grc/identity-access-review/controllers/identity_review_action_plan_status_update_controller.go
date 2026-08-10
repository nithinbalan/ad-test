package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityReviewActionPlanStatusUpdateController struct {
	service services.IdentityReviewActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewActionPlanStatusUpdateController(service services.IdentityReviewActionPlanStatusUpdateServiceInterface) *IdentityReviewActionPlanStatusUpdateController {
	return &IdentityReviewActionPlanStatusUpdateController{
		service: service,
	}
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) GetAllIdentityReviewActionPlanStatusUpdates(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "desc")
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	validSortFields := map[string]string{
		"created_at": "identity_review_action_plan_status_updates.created_at",
		"updated_at": "identity_review_action_plan_status_updates.updated_at",
		"percentage": "identity_review_action_plan_status_updates.percentage",
		"status_en":  "action_plan_statuses.title_en",
		"status_ar":  "action_plan_statuses.title_ar",
	}

	sortBy, valid := validSortFields[c.DefaultQuery("sort_by", "created_at")]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	result, err := irapsuc.service.GetAll(uint(userID), uint(actionPlanID), sort, sortBy, c.Query("keyword"))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) GetIdentityReviewActionPlanStatusUpdateDetails(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdateIDStr := c.Param("statusUpdateId")

	id, err := strconv.ParseUint(statusUpdateIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := irapsuc.service.GetByID(uint(userID), uint(actionPlanID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) GetLatestIdentityReviewActionPlanStatusUpdate(c *gin.Context) {
	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := irapsuc.service.GetLatestIdentityReviewActionPlanStatusUpdate(uint(userID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) CreateIdentityReviewActionPlanStatusUpdate(c *gin.Context) {
	var form dtos.IdentityReviewActionPlanStatusUpdateFormDTO

	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapsuc.db, c, err.Error())
		return
	}

	id, err := irapsuc.service.Create(uint(userID), uint(actionPlanID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) UpdateIdentityReviewActionPlanStatusUpdate(c *gin.Context) {
	var form dtos.IdentityReviewActionPlanStatusUpdateFormDTO

	idStr := c.Param("id")

	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdateIDStr := c.Param("statusUpdateId")

	id, err := strconv.ParseUint(statusUpdateIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapsuc.db, c, err.Error())
		return
	}

	err = irapsuc.service.Update(uint(userID), uint(actionPlanID), uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) ApproveIdentityReviewActionPlan(c *gin.Context) {
	var identityReviewActionPlanReviewFormDTO dtos.IdentityReviewActionPlanReviewFormDTO

	identityReviewUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&identityReviewActionPlanReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapsuc.db, c, err.Error())
		return
	}

	err = irapsuc.service.Approve(uint(identityReviewUserID), uint(identityReviewActionPlanID), currentUserID, &identityReviewActionPlanReviewFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) RevertIdentityReviewActionPlan(c *gin.Context) {
	var identityReviewActionPlanReviewFormDTO dtos.IdentityReviewActionPlanReviewFormDTO

	identityReviewUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&identityReviewActionPlanReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapsuc.db, c, err.Error())
		return
	}

	err = irapsuc.service.Revert(uint(identityReviewUserID), uint(identityReviewActionPlanID), currentUserID, &identityReviewActionPlanReviewFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) GetIdentityReviewActionPlanDocumentPreview(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanIDStr := c.Param("actionPlanId")
	identityReviewActionPlanID, err := strconv.ParseUint(identityReviewActionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanStatusUpdateIDStr := c.Param("statusUpdateId")
	identityReviewActionPlanStatusUpdateID, err := strconv.ParseUint(identityReviewActionPlanStatusUpdateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := irapsuc.service.GetDocumentPath(uint(identityReviewUserID), uint(identityReviewActionPlanID), uint(identityReviewActionPlanStatusUpdateID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	c.File(documentPath)
}

func (irapsuc *IdentityReviewActionPlanStatusUpdateController) DownloadIdentityReviewActionPlanDocument(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanIDStr := c.Param("actionPlanId")
	identityReviewActionPlanID, err := strconv.ParseUint(identityReviewActionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanStatusUpdateIDStr := c.Param("statusUpdateId")
	identityReviewActionPlanStatusUpdateID, err := strconv.ParseUint(identityReviewActionPlanStatusUpdateIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := irapsuc.service.GetDocumentPath(uint(identityReviewUserID), uint(identityReviewActionPlanID), uint(identityReviewActionPlanStatusUpdateID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapsuc.db, c, e.Error())
			return
		}
	}

	fileName := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
