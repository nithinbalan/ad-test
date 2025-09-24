package controllers

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanCommentController struct {
	service services.AuditFindingActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewAuditInsightsFindingActionPlanCommentController(service services.AuditFindingActionPlanCommentServiceInterface) *AuditInsightsFindingActionPlanCommentController {
	return &AuditInsightsFindingActionPlanCommentController{
		service: service,
	}
}

func (afapc *AuditInsightsFindingActionPlanCommentController) GetAllComments(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	comments, err := afapc.service.GetAll(uint(findingID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (afapc *AuditInsightsFindingActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	comment, err := afapc.service.GetFullDetailsByID(uint(findingID), uint(actionPlanID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (afapc *AuditInsightsFindingActionPlanCommentController) CreateComment(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	var formDTO auditDTOs.AuditFindingActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	id, err := afapc.service.Create(uint(findingID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (afapc *AuditInsightsFindingActionPlanCommentController) UpdateComment(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	var formDTO auditDTOs.AuditFindingActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	err = afapc.service.Update(uint(findingID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (afapc *AuditInsightsFindingActionPlanCommentController) DeleteComment(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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
		utility.SendInternalErrResp(afapc.db, c, err.Error())
		return
	}

	failedIDs, err := afapc.service.Delete(uint(findingID), uint(actionPlanID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(afapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (afapc *AuditInsightsFindingActionPlanCommentController) GetAuditFindingActionPlanCommentDocumentPreview(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	complianceActionPlanCommentDocumentPath, err := afapc.service.GetDocumentPath(uint(findingID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afapc.db, c, err.Error())
		}
		return
	}
	c.File(complianceActionPlanCommentDocumentPath)
}

func (afapc *AuditInsightsFindingActionPlanCommentController) DownloadAuditFindingActionPlanCommentDocument(c *gin.Context) {
	findingIDStr := c.Param("findingId")
	findingID, err := strconv.ParseUint(findingIDStr, 10, 32)
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

	complianceActionPlanCommentDocumentPath, err := afapc.service.GetDocumentPath(uint(findingID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(afapc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(complianceActionPlanCommentDocumentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(complianceActionPlanCommentDocumentPath)
}
