package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementActionPlanCommentController struct {
	service services.ComplianceAssessmentRequirementActionPlanCommentServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementActionPlanCommentController(service services.ComplianceAssessmentRequirementActionPlanCommentServiceInterface) *ComplianceAssessmentRequirementActionPlanCommentController {
	return &ComplianceAssessmentRequirementActionPlanCommentController{
		service: service,
	}
}

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) GetAllComments(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	comments, err := capc.service.GetAll(uint(assessmentID), uint(requirementID), uint(actionPlanID))
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

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) GetCommentDetails(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	comment, err := capc.service.GetFullDetailsByID(uint(assessmentID), uint(requirementID), uint(actionPlanID), uint(commentID))
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

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) CreateComment(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	var formDTO complianceDTOs.ComplianceActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}
	id, err := capc.service.Create(uint(assessmentID), uint(requirementID), uint(actionPlanID), currentUserID, &formDTO)
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

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) UpdateComment(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	var formDTO complianceDTOs.ComplianceActionPlanCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	err = capc.service.Update(uint(assessmentID), uint(requirementID), uint(actionPlanID), uint(commentID), currentUserID, &formDTO)
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

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) DeleteComment(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")
	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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
		utility.SendInternalErrResp(capc.db, c, err.Error())
		return
	}

	failedIDs, err := capc.service.Delete(uint(assessmentID), uint(requirementID), uint(actionPlanID), request.IDs, currentUserID)
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

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) GetComplianceAssessmentRequirementActionPlanCommentDocumentPreview(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	complianceAssessmentRequirementActionPlanCommentDocumentPath, err := capc.service.GetDocumentPath(uint(assessmentID), uint(requirementID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(capc.db, c, err.Error())
		}
		return
	}
	c.File(complianceAssessmentRequirementActionPlanCommentDocumentPath)
}

func (capc *ComplianceAssessmentRequirementActionPlanCommentController) DownloadComplianceAssessmentRequirementActionPlanCommentDocument(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	requirementIDStr := c.Param("requirementId")
	requirementID, err := strconv.ParseUint(requirementIDStr, 10, 32)
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

	complianceAssessmentRequirementActionPlanCommentDocumentPath, err := capc.service.GetDocumentPath(uint(assessmentID), uint(requirementID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(capc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(complianceAssessmentRequirementActionPlanCommentDocumentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(complianceAssessmentRequirementActionPlanCommentDocumentPath)
}
