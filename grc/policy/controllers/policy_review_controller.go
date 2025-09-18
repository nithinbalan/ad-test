package controllers

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyReviewController struct {
	service services.DocumentReviewServiceInterface
	db      *gorm.DB
}

func NewPolicyReviewController(service services.DocumentReviewServiceInterface) *PolicyReviewController {
	return &PolicyReviewController{
		service: service,
	}
}

func (prc *PolicyReviewController) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := prc.service.GetComments(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (prc *PolicyReviewController) AddComment(c *gin.Context) {
	var documentReviewCommentFormDTO documentDTOs.DocumentReviewCommentFormDTO
	if err := c.ShouldBindJSON(&documentReviewCommentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.AddComment(uint(id), currentUserID, &documentReviewCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
