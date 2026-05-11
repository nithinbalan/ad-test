package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanRiskTeamAssessmentController struct {
	service services.RiskUpdatePlanRiskTeamAssessmentServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanRiskTeamAssessmentController(service services.RiskUpdatePlanRiskTeamAssessmentServiceInterface) *RiskUpdatePlanRiskTeamAssessmentController {
	return &RiskUpdatePlanRiskTeamAssessmentController{
		service: service,
	}
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) GetRiskTeamAssessmentUser(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	user, err := rrupc.service.GetRiskTeamAssessmentUser(uint(planID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, user)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) CreateRiskTeamAssessment(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateRiskTeamAssessment(uint(planID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) MarkRiskTeamAssessmentAsNotApplicable(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = rrupc.service.MarkRiskTeamAssessmentAsNotApplicable(uint(planID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) GetRiskTeamAssessmentDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := rrupc.service.GetRiskTeamAssessmentDocuments(uint(planID))
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) CreateRiskTeamAssessmentDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	var documentsWrapper riskDTOs.CommonDocumentFormDTO
	if err := c.ShouldBindJSON(&documentsWrapper); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = rrupc.service.CreateRiskTeamAssessmentDocuments(uint(planID), currentUserID, documentsWrapper.Documents)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) GetRiskTeamAssessmentDocumentPreview(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rrupc.service.GetRiskTeamAssessmentDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
			return
		}
	}

	c.File(filePath)
}

func (rrupc *RiskUpdatePlanRiskTeamAssessmentController) DownloadRiskTeamAssessmentDocument(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendInvalidErrResp(c, errors.New("token is required"))
		return
	}

	filePath, err := rrupc.service.GetRiskTeamAssessmentDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
			return
		}
	}

	filename := filepath.Base(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}
