package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportActionPlanCommentController struct {
	service services.ViolationReportActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewViolationReportActionPlanCommentController(service services.ViolationReportActionPlanCommentServiceInterface) *ViolationReportActionPlanCommentController {
	return &ViolationReportActionPlanCommentController{
		service: service,
	}
}

func (vrapc *ViolationReportActionPlanCommentController) GetAllComments(c *gin.Context) {
	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
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

	comments, err := vrapc.service.GetAll(uint(violationReportID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (vrapc *ViolationReportActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
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

	comment, err := vrapc.service.GetFullDetailsByID(uint(violationReportID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (vrapc *ViolationReportActionPlanCommentController) CreateComment(c *gin.Context) {
	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
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

	var formDTO violationReportDTOs.ViolationReportActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrapc.db, c, err.Error())
		return
	}

	id, err := vrapc.service.Create(uint(violationReportID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (vrapc *ViolationReportActionPlanCommentController) UpdateComment(c *gin.Context) {
	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
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

	var formDTO violationReportDTOs.ViolationReportActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrapc.db, c, err.Error())
		return
	}

	err = vrapc.service.Update(uint(violationReportID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vrapc *ViolationReportActionPlanCommentController) DeleteComment(c *gin.Context) {
	violationReportIDStr := c.Param("id")
	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
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
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(vrapc.db, c, err.Error())
		return
	}

	failedIDs, err := vrapc.service.Delete(uint(violationReportID), uint(actionPlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (vrapc *ViolationReportActionPlanCommentController) GetViolationReportActionPlanCommentDocumentPreview(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
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

	path, err := vrapc.service.GetDocumentPath(uint(violationReportID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrapc.db, c, err.Error())
		}
		return
	}
	c.File(path)
}

func (vrapc *ViolationReportActionPlanCommentController) DownloadViolationReportActionPlanCommentDocument(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
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

	path, err := vrapc.service.GetDocumentPath(uint(violationReportID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrapc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(path)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(path)
}
