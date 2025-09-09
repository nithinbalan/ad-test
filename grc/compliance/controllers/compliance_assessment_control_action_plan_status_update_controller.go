package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlActionPlanStatusUpdateController struct {
	service services.ComplianceAssessmentControlActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentControlActionPlanStatusUpdateController(service services.ComplianceAssessmentControlActionPlanStatusUpdateServiceInterface) *ComplianceAssessmentControlActionPlanStatusUpdateController {
	return &ComplianceAssessmentControlActionPlanStatusUpdateController{
		service: service,
	}
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) GetAllComplianceAssessmentControlActionPlanStatusUpdates(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := cacapsuc.service.GetAll(uint(assessmentID), uint(controlID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) CreateComplianceAssessmentControlActionPlanStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO

	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&complianceStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		return
	}

	createdID, err := cacapsuc.service.Create(uint(assessmentID), uint(controlID), uint(actionPlanID), currentUserID, &complianceStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) UpdateComplianceAssessmentControlActionPlanStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO

	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		return
	}

	err = cacapsuc.service.Update(uint(assessmentID), uint(controlID), uint(actionPlanID), uint(id), currentUserID, &complianceStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) GetLatestComplianceAssessmentControlActionPlanStatusUpdate(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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

	complianceActionPlanStatusUpdate, err := cacapsuc.service.GetLatestComplianceAssessmentControlActionPlanStatusUpdate(uint(assessmentID), uint(controlID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceActionPlanStatusUpdate)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) GetComplianceAssessmentControlActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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

	idStr := c.Param("id")
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

	complianceActionPlanStatusUpdateDocumentPath, err := cacapsuc.service.GetDocumentPath(uint(assessmentID), uint(controlID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		}
		return
	}
	c.File(complianceActionPlanStatusUpdateDocumentPath)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) DownloadComplianceAssessmentControlActionPlanStatusUpdateDocument(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
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
	idStr := c.Param("id")
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

	complianceActionPlanStatusUpdateDocumentPath, err := cacapsuc.service.GetDocumentPath(uint(assessmentID), uint(controlID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(complianceActionPlanStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(complianceActionPlanStatusUpdateDocumentPath)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) ApproveComplianceAssessmentControlActionPlan(c *gin.Context) {
	var complianceStatusReviewUserUpdateFormDTO complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO

	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		return
	}

	err = cacapsuc.service.Approve(uint(assessmentID), uint(controlID), uint(complianceActionPlanID), currentUserID, &complianceStatusReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (cacapsuc *ComplianceAssessmentControlActionPlanStatusUpdateController) RevertComplianceAssessmentControlActionPlan(c *gin.Context) {
	var complianceStatusReviewUserUpdateFormDTO complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO

	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceActionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceStatusReviewUserUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cacapsuc.db, c, err.Error())
		return
	}

	err = cacapsuc.service.Revert(uint(assessmentID), uint(controlID), uint(complianceActionPlanID), currentUserID, &complianceStatusReviewUserUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(cacapsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
