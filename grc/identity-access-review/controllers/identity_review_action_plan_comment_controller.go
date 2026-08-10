package controllers

import (
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

type IdentityReviewActionPlanCommentController struct {
	service services.IdentityReviewActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewActionPlanCommentController(service services.IdentityReviewActionPlanCommentServiceInterface) *IdentityReviewActionPlanCommentController {
	return &IdentityReviewActionPlanCommentController{
		service: service,
	}
}

func (irapcc *IdentityReviewActionPlanCommentController) GetAllComments(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	comments, err := irapcc.service.GetAll(uint(identityReviewUserID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (irapcc *IdentityReviewActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	commentIDStr := c.Param("commentId")

	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comment, err := irapcc.service.GetFullDetailsByID(uint(identityReviewUserID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (irapcc *IdentityReviewActionPlanCommentController) CreateComment(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	var formDTO dtos.IdentityReviewActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapcc.db, c, err.Error())
		return
	}

	id, err := irapcc.service.Create(uint(identityReviewUserID), uint(actionPlanID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (irapcc *IdentityReviewActionPlanCommentController) UpdateComment(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	commentIDStr := c.Param("commentId")

	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.IdentityReviewActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapcc.db, c, err.Error())
		return
	}

	err = irapcc.service.Update(uint(identityReviewUserID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (irapcc *IdentityReviewActionPlanCommentController) DeleteComment(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(irapcc.db, c, err.Error())
		return
	}

	failedIDs, err := irapcc.service.Delete(uint(identityReviewUserID), uint(actionPlanID), request.IDs, currentUserID)
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
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (irapcc *IdentityReviewActionPlanCommentController) GetIdentityReviewActionPlanCommentDocumentPreview(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 64)
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

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanCommentDocumentPath, err := irapcc.service.GetDocumentPath(uint(identityReviewUserID), uint(actionPlanID), uint(commentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	c.File(identityReviewActionPlanCommentDocumentPath)
}

func (irapcc *IdentityReviewActionPlanCommentController) DownloadIdentityReviewActionPlanCommentDocument(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")
	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 64)
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

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewActionPlanCommentDocumentPath, err := irapcc.service.GetDocumentPath(uint(identityReviewUserID), uint(actionPlanID), uint(commentID), token)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(irapcc.db, c, e.Error())
			return
		}
	}

	filename := filepath.Base(identityReviewActionPlanCommentDocumentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(identityReviewActionPlanCommentDocumentPath)
}
