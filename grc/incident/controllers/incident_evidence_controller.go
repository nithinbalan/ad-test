package controllers

import (
	"errors"
	"jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentEvidenceController struct {
	service services.IncidentEvidenceServiceInterface
	db      *gorm.DB
}

func NewIncidentEvidenceController(service services.IncidentEvidenceServiceInterface) *IncidentEvidenceController {
	return &IncidentEvidenceController{
		service: service,
	}
}

func (iec *IncidentEvidenceController) GetAllIncidentEvidences(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")
	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid incident id"))
		return
	}

	keyword := c.Query("keyword")

	result, err := iec.service.GetAll(uint(incidentID), keyword)
	if err != nil {
		utility.SendInternalErrResp(iec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (iec *IncidentEvidenceController) GetIncidentEvidenceDetails(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evidenceID, err := strconv.ParseUint(c.Param("evidenceId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := iec.service.GetFullDetailsByID(uint(incidentID), uint(evidenceID))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (iec *IncidentEvidenceController) CreateIncidentEvidence(c *gin.Context) {
	var form *dtos.IncidentEvidenceFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid incident id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iec.db, c, err.Error())
		return
	}

	ids, err := iec.service.Create(uint(incidentID), form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"ids": ids})
}

func (iec *IncidentEvidenceController) UpdateIncidentEvidence(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.IncidentEvidenceFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	evidenceID, err := strconv.ParseUint(c.Param("evidenceId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iec.db, c, err.Error())
		return
	}

	err = iec.service.Update(uint(incidentID), uint(evidenceID), &form, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iec.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iec *IncidentEvidenceController) DeleteIncidentEvidences(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
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
		utility.SendInternalErrResp(iec.db, c, err.Error())
		return
	}

	failedIDs, err := iec.service.Delete(uint(incidentID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(iec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iec *IncidentEvidenceController) GetIncidentEvidenceDocumentPreview(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evidenceID, err := strconv.ParseUint(c.Param("evidenceId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	filePath, err := iec.service.GetIncidentEvidencePath(uint(incidentID), uint(evidenceID), token, 0)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iec.db, c, err.Error())
		}
		return
	}

	c.File(filePath)
}

func (iec *IncidentEvidenceController) DownloadIncidentEvidence(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evidenceID, err := strconv.ParseUint(c.Param("evidenceId"), 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	// This is an unauthenticated route (accessed via token), so the current user
	// may not be available. History is only recorded when a user id is present.
	currentUserID, _ := utility.GetCurrentUserID(c)

	filePath, err := iec.service.GetIncidentEvidencePath(uint(incidentID), uint(evidenceID), token, currentUserID)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iec.db, c, err.Error())
		}
		return
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

func (iec *IncidentEvidenceController) GetIncidentEvidenceDownloadHistories(c *gin.Context) {
	incidentID, err := strconv.ParseUint(c.Param("incidentID"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	evidenceID, err := strconv.ParseUint(c.Param("evidenceId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := iec.service.GetDownloadHistories(uint(incidentID), uint(evidenceID))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(iec.db, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
