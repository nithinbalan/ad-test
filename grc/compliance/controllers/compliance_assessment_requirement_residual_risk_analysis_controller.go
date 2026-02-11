package controllers

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementResidualRiskAnalysisController struct {
	service services.ComplianceAssessmentRequirementResidualRiskAnalysisServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementResidualRiskAnalysisController(service services.ComplianceAssessmentRequirementResidualRiskAnalysisServiceInterface) *ComplianceAssessmentRequirementResidualRiskAnalysisController {
	return &ComplianceAssessmentRequirementResidualRiskAnalysisController{
		service: service,
	}
}

func (carrac *ComplianceAssessmentRequirementResidualRiskAnalysisController) CreateResidualRiskAnalysis(c *gin.Context) {
	responsibleUnitIDStr := c.Param("responsibleUnitId")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "responsible_unit_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	var riskAnalysisFormDTO []*complianceDTOs.ComplianceAssessmentRequirementResidualRiskFormDTO
	if err := c.ShouldBindJSON(&riskAnalysisFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = carrac.service.CreateResidualRiskAnalysis(uint(responsibleUnitID), currentUserID, riskAnalysisFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(carrac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (carrac *ComplianceAssessmentRequirementResidualRiskAnalysisController) DeleteResidualRiskAnalysis(c *gin.Context) {
	responsibleUnitIDStr := c.Param("responsibleUnitId")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "responsible_unit_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	err = carrac.service.DeleteResidualRiskAnalysis(uint(responsibleUnitID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// GetAllResidualControls retrieves all residual controls for a responsible unit with pagination
func (carrac *ComplianceAssessmentRequirementResidualRiskAnalysisController) GetAllResidualControls(c *gin.Context) {
	responsibleUnitIDStr := c.Param("responsibleUnitId")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "responsible_unit_id"))
		return
	}

	// Parse query parameters
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	queryParams := &complianceDTOs.ComplianceAssessmentResidualControlListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		QuickFilter: quickFilter,
	}

	response, err := carrac.service.GetAllByResponsibleUnitID(uint(responsibleUnitID), queryParams)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateResidualControls creates residual control mappings for a responsible unit
func (carrac *ComplianceAssessmentRequirementResidualRiskAnalysisController) CreateResidualControls(c *gin.Context) {
	responsibleUnitIDStr := c.Param("responsibleUnitId")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "responsible_unit_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	var residualControlFormDTO []*complianceDTOs.ComplianceAssessmentResidualControlFormDTO
	if err := c.ShouldBindJSON(&residualControlFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = carrac.service.CreateControl(uint(responsibleUnitID), currentUserID, residualControlFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(carrac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

// DeleteResidualControls deletes all residual control mappings for a responsible unit
func (carrac *ComplianceAssessmentRequirementResidualRiskAnalysisController) DeleteResidualControls(c *gin.Context) {
	responsibleUnitIDStr := c.Param("responsibleUnitId")
	responsibleUnitID, err := strconv.ParseUint(responsibleUnitIDStr, 10, 32)
	if err != nil {
		utility.SendInvalidErrResp(c, utility.NewValidationError(err, "invalid", "responsible_unit_id"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	err = carrac.service.DeleteControl(uint(responsibleUnitID), currentUserID)
	if err != nil {
		utility.SendInternalErrResp(carrac.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
