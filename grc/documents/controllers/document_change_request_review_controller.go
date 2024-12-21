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

type DocumentChangeRequestReviewController struct {
	service services.DocumentChangeRequestReviewServiceInterface
	db      *gorm.DB
}

func NewDocumentChangeRequestReviewController(service services.DocumentChangeRequestReviewServiceInterface) *DocumentChangeRequestReviewController {
	return &DocumentChangeRequestReviewController{
		service: service,
	}
}

func (dcrrc *DocumentChangeRequestReviewController) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := dcrrc.service.GetComments(uint(id))
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

func (dcrrc *DocumentChangeRequestReviewController) AddComment(c *gin.Context) {
	var documentCRReviewCommentFormDTO documentDTOs.DocumentCRReviewCommentFormDTO
	if err := c.ShouldBindJSON(&documentCRReviewCommentFormDTO); err != nil {
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
		utility.SendInternalErrResp(dcrrc.db, c, err.Error())
		return
	}

	err = dcrrc.service.AddComment(uint(id), currentUserID, &documentCRReviewCommentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(dcrrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
