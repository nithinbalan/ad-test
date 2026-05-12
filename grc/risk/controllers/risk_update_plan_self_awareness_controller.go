package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskUpdatePlanSelfAwarenessController struct {
	service services.RiskUpdatePlanSelfAwarenessServiceInterface
	db      *gorm.DB
}

func NewRiskUpdatePlanSelfAwarenessController(service services.RiskUpdatePlanSelfAwarenessServiceInterface) *RiskUpdatePlanSelfAwarenessController {
	return &RiskUpdatePlanSelfAwarenessController{
		service: service,
	}
}

// Self-Awareness Training Document Operations

func (rrupc *RiskUpdatePlanSelfAwarenessController) GetSelfAwarenessTrainingDocumentPreview(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrupc.service.GetSelfAwarenessTrainingDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (rrupc *RiskUpdatePlanSelfAwarenessController) DownloadSelfAwarenessTrainingDocument(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := rrupc.service.GetSelfAwarenessTrainingDocumentPath(uint(planID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

func (rrupc *RiskUpdatePlanSelfAwarenessController) GetSelfAwarenessTrainingDocuments(c *gin.Context) {
	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	documents, err := rrupc.service.GetSelfAwarenessTrainingDocuments(uint(planID))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(rrupc.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func (rrupc *RiskUpdatePlanSelfAwarenessController) CreateSelfAwarenessTrainingDocuments(c *gin.Context) {
	var formDTO riskDTOs.CommonDocumentFormDTO

	planIDStr := c.Param("id")
	planID, err := strconv.ParseUint(planIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.CreateSelfAwarenessTrainingDocuments(uint(planID), currentUserID, formDTO.Documents)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documents updated successfully"})
}

func (rrupc *RiskUpdatePlanSelfAwarenessController) MarkSelfAwarenessTrainingAsNotApplicable(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO riskDTOs.MarkTimelineNotApplicableDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrupc.db, c, err.Error())
		return
	}

	err = rrupc.service.MarkSelfAwarenessTrainingAsNotApplicable(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrupc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
