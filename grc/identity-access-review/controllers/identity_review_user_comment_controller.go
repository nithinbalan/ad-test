package controllers

import (
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityReviewUserCommentController struct {
	service services.IdentityReviewUserCommentServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewUserCommentController(service services.IdentityReviewUserCommentServiceInterface) *IdentityReviewUserCommentController {
	return &IdentityReviewUserCommentController{
		service: service,
	}
}

func (iruscc *IdentityReviewUserCommentController) GetAll(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")

	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := iruscc.service.GetAll(uint(identityReviewUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(iruscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (iruscc *IdentityReviewUserCommentController) GetDetails(c *gin.Context) {
	identityReviewUserIDStr := c.Param("id")

	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	comment, err := iruscc.service.GetByID(uint(identityReviewUserID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(iruscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (iruscc *IdentityReviewUserCommentController) Create(c *gin.Context) {
	var form dtos.IdentityReviewUserCommentFormDTO

	identityReviewUserIDStr := c.Param("id")

	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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
		utility.SendInternalErrResp(iruscc.db, c, err.Error())
		return
	}

	id, err := iruscc.service.Create(uint(identityReviewUserID), currentUserID, &form)
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
			utility.SendInternalErrResp(iruscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (iruscc *IdentityReviewUserCommentController) Update(c *gin.Context) {
	var form dtos.IdentityReviewUserCommentFormDTO

	identityReviewUserIDStr := c.Param("id")

	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iruscc.db, c, err.Error())
		return
	}

	err = iruscc.service.Update(uint(identityReviewUserID), uint(commentID), currentUserID, &form)
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
			utility.SendInternalErrResp(iruscc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iruscc *IdentityReviewUserCommentController) Delete(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	identityReviewUserIDStr := c.Param("id")

	identityReviewUserID, err := strconv.ParseUint(identityReviewUserIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iruscc.db, c, err.Error())
		return
	}

	failedIDs, err := iruscc.service.Delete(uint(identityReviewUserID), request.IDs, currentUserID)
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
			utility.SendInternalErrResp(iruscc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
