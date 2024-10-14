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

type DocumentReviewController struct {
	service services.DocumentReviewServiceInterface
	db      *gorm.DB
}

func NewDocumentReviewController(service services.DocumentReviewServiceInterface) *DocumentReviewController {
	return &DocumentReviewController{
		service: service,
	}
}

func (drc *DocumentReviewController) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := drc.service.GetComments(uint(id))
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

func (drc *DocumentReviewController) AddComment(c *gin.Context) {
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
		utility.SendInternalErrResp(drc.db, c, err.Error())
		return
	}

	err = drc.service.AddComment(uint(id), currentUserID, &documentReviewCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(drc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
