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

type ComplianceRequirementStatusUpdateController struct {
	service services.ComplianceRequirementStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewComplianceRequirementStatusUpdateController(service services.ComplianceRequirementStatusUpdateServiceInterface) *ComplianceRequirementStatusUpdateController {
	return &ComplianceRequirementStatusUpdateController{
		service: service,
	}
}

func (crsuc *ComplianceRequirementStatusUpdateController) GetAllComplianceRequirementStatusUpdates(c *gin.Context) {
	complianceIdStr := c.Param("complianceId")

	complianceId, err := strconv.ParseUint(complianceIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIdStr := c.Param("requirementId")

	requirementId, err := strconv.ParseUint(requirementIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	complianceStatuses, err := crsuc.service.GetAll(uint(complianceId), uint(requirementId), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crsuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceStatuses)
}

func (crsuc *ComplianceRequirementStatusUpdateController) CreateComplianceRequirementStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceRequirementStatusUpdateFormDTO

	complianceIdStr := c.Param("complianceId")

	complianceId, err := strconv.ParseUint(complianceIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIdStr := c.Param("requirementId")

	requirementId, err := strconv.ParseUint(requirementIdStr, 10, 32)
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
		utility.SendInternalErrResp(crsuc.db, c, err.Error())
		return
	}

	createdID, err := crsuc.service.Create(uint(complianceId), uint(requirementId), currentUserID, &complianceStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (crsuc *ComplianceRequirementStatusUpdateController) UpdateComplianceRequirementStatusUpdate(c *gin.Context) {
	var complianceStatusUpdateFormDTO complianceDTOs.ComplianceRequirementStatusUpdateFormDTO

	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIdStr := c.Param("requirementId")

	requirementId, err := strconv.ParseUint(requirementIdStr, 10, 32)
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
		utility.SendInternalErrResp(crsuc.db, c, err.Error())
		return
	}

	err = crsuc.service.Update(uint(complianceID), uint(requirementId), uint(id), currentUserID, &complianceStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (csuc *ComplianceRequirementStatusUpdateController) GetLatestComplianceRequirementStatusUpdate(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIdStr := c.Param("requirementId")

	requirementId, err := strconv.ParseUint(requirementIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	complianceRequirementStatusUpdate, err := csuc.service.GetLatestComplianceRequirementStatusUpdate(uint(complianceID), uint(requirementId))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, complianceRequirementStatusUpdate)
}
func (casuc *ComplianceRequirementStatusUpdateController) GetComplianceRequirementStatusUpdateDocumentPreview(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	complianceRequirementStatusUpdateDocumentPath, err := casuc.service.GetDocumentPath(uint(complianceID), uint(requirementID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(casuc.db, c, err.Error())
		}
		return
	}
	c.File(complianceRequirementStatusUpdateDocumentPath)
}

func (casuc *ComplianceRequirementStatusUpdateController) DownloadComplianceRequirementStatusUpdateDocument(c *gin.Context) {
	complianceIDStr := c.Param("complianceId")

	complianceID, err := strconv.ParseUint(complianceIDStr, 10, 64)
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

	complianceRequirementStatusUpdateDocumentPath, err := casuc.service.GetDocumentPath(uint(complianceID), uint(requirementID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(casuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(complianceRequirementStatusUpdateDocumentPath)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(complianceRequirementStatusUpdateDocumentPath)
}
