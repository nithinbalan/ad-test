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

type ComplianceAssessmentRequirementController struct {
	service services.ComplianceAssessmentRequirementServiceInterface
	db      *gorm.DB
}

func NewComplianceAssessmentRequirementController(service services.ComplianceAssessmentRequirementServiceInterface) *ComplianceAssessmentRequirementController {
	return &ComplianceAssessmentRequirementController{
		service: service,
	}
}

func (cac *ComplianceAssessmentRequirementController) GetComplianceAssessmentRequirementDetails(c *gin.Context) {
	assessmentIDStr := c.Param("assessmentId")

	assessmentID, err := strconv.ParseUint(assessmentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	complianceAssessment, err := cac.service.GetFullDetailsByID(uint(assessmentID), uint(complianceRequirementID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, complianceAssessment)
}

func (cac *ComplianceAssessmentRequirementController) UpdateComplianceAssessmentRequirement(c *gin.Context) {
	var complianceAssessmentRequirementUpdateFormDTO complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO

	complianceAssessmentID, err := strconv.Atoi(c.Param("assessmentId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	requirementIDStr := c.Param("requirementId")

	complianceRequirementID, err := strconv.ParseUint(requirementIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&complianceAssessmentRequirementUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(cac.db, c, err.Error())
		return
	}

	err = cac.service.Update(uint(complianceAssessmentID), uint(complianceRequirementID), currentUserID, &complianceAssessmentRequirementUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(cac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
