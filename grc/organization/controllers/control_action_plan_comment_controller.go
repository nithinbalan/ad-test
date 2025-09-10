package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlActionPlanCommentController struct {
	service services.ControlActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewControlActionPlanCommentController(service services.ControlActionPlanCommentServiceInterface) *ControlActionPlanCommentController {
	return &ControlActionPlanCommentController{
		service: service,
	}
}

func (capc *ControlActionPlanCommentController) GetAllComments(c *gin.Context) {
	controlIDStr := c.Param("controlId")
	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("controlActionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comments, err := capc.service.GetAll(uint(controlID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (capc *ControlActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	controlIDStr := c.Param("controlId")
	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("controlActionPlanId")
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

	comment, err := capc.service.GetFullDetailsByID(uint(controlID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (capc *ControlActionPlanCommentController) CreateComment(c *gin.Context) {
	controlIDStr := c.Param("controlId")
	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("controlActionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO controlDTOs.ControlActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	id, err := capc.service.Create(uint(controlID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (capc *ControlActionPlanCommentController) UpdateComment(c *gin.Context) {
	controlIDStr := c.Param("controlId")
	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("controlActionPlanId")
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

	var formDTO controlDTOs.ControlActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	err = capc.service.Update(uint(controlID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (capc *ControlActionPlanCommentController) DeleteComment(c *gin.Context) {
	controlIDStr := c.Param("controlId")
	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanIDStr := c.Param("controlActionPlanId")
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	failedIDs, err := capc.service.Delete(uint(controlID), uint(actionPlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(capc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (capc *ControlActionPlanCommentController) GetControlActionPlanCommentDocumentPreview(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("controlActionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("commentId")
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

	controlActionPlanCommentDocumentPath, err := capc.service.GetDocumentPath(uint(controlID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(capc.db, c, err.Error())
		}
		return
	}
	c.File(controlActionPlanCommentDocumentPath)
}

func (capc *ControlActionPlanCommentController) DownloadControlActionPlanCommentDocument(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("controlActionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	idStr := c.Param("commentId")
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

	controlActionPlanCommentDocumentPath, err := capc.service.GetDocumentPath(uint(controlID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(capc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(controlActionPlanCommentDocumentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(controlActionPlanCommentDocumentPath)
}
