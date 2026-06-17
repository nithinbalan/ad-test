package services

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type ComplianceEvidenceSubmissionServiceInterface interface {
	GetAllAssessments(params *complianceDTOs.ComplianceEvidenceSubmissionListQueryParamsDTO, currentUserID uint) (*complianceDTOs.ComplianceEvidenceSubmissionListResponseDTO, error)
	GetRequirementAssessmentFullDetailsByID(id uint) (*complianceDTOs.ComplianceEvidenceSubmissionRequirementDetailDTO, error)
}

type ComplianceEvidenceSubmissionService struct {
	repo    complianceRepositories.ComplianceEvidenceSubmissionRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewComplianceEvidenceSubmissionService(
	repo complianceRepositories.ComplianceEvidenceSubmissionRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
) ComplianceEvidenceSubmissionServiceInterface {
	return &ComplianceEvidenceSubmissionService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (cpas *ComplianceEvidenceSubmissionService) GetAllAssessments(params *complianceDTOs.ComplianceEvidenceSubmissionListQueryParamsDTO, currentUserID uint) (*complianceDTOs.ComplianceEvidenceSubmissionListResponseDTO, error) {
	isSuperAdmin, err := cpas.userSvc.IsSuperAdmin(currentUserID)
	if err != nil {
		return nil, err
	}

	if isSuperAdmin {
		params.IsSuperAdmin = true
	}

	return cpas.repo.GetAll(params, currentUserID)
}

func (cpas *ComplianceEvidenceSubmissionService) GetRequirementAssessmentFullDetailsByID(id uint) (*complianceDTOs.ComplianceEvidenceSubmissionRequirementDetailDTO, error) {
	responsibleUnit, err := cpas.repo.GetRequirementAssessmentFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var dto complianceDTOs.ComplianceEvidenceSubmissionRequirementDetailDTO

	if err := mapstructure.Decode(responsibleUnit, &dto); err != nil {
		return nil, err
	}

	dto.CompletedAt = responsibleUnit.CompletedAt
	dto.CreatedAt = responsibleUnit.CreatedAt
	dto.UpdatedAt = responsibleUnit.UpdatedAt

	if responsibleUnit.ComplianceAssessmentRequirement != nil && responsibleUnit.ComplianceAssessmentRequirement.ComplianceRequirement != nil {
		if err := mapstructure.Decode(responsibleUnit.ComplianceAssessmentRequirement.ComplianceRequirement, &dto.ComplianceRequirement); err != nil {
			return nil, err
		}

		if responsibleUnit.ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause != nil &&
			responsibleUnit.ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister != nil {
			if err := mapstructure.Decode(responsibleUnit.ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister, &dto.ComplianceRegister); err != nil {
				return nil, err
			}
		}

		if responsibleUnit.ComplianceAssessmentRequirement.ComplianceAssessment != nil {
			if err := mapstructure.Decode(responsibleUnit.ComplianceAssessmentRequirement.ComplianceAssessment, &dto.ComplianceAssessment); err != nil {
				return nil, err
			}
		}
	}

	for i, c := range dto.RequirementComments {
		if i >= len(responsibleUnit.RequirementComments) {
			break
		}
		c.CreatedAt = responsibleUnit.RequirementComments[i].CreatedAt
		dto.RequirementComments[i] = c
	}

	for i, title := range dto.RequiredDocumentTitles {
		if i >= len(responsibleUnit.RequiredDocumentTitles) {
			break
		}
		if responsibleUnit.RequiredDocumentTitles[i].Documents == nil {
			continue
		}
		if dto.RequiredDocumentTitles[i].Documents == nil {
			continue
		}
		for j, doc := range title.Documents {
			if j >= len(*responsibleUnit.RequiredDocumentTitles[i].Documents) {
				break
			}
			if j >= len(dto.RequiredDocumentTitles[i].Documents) {
				break
			}
			dto.RequiredDocumentTitles[i].Documents[j] = doc
		}
	}

	for i, actionPlan := range dto.ActionPlans {
		if i >= len(responsibleUnit.ActionPlans) {
			break
		}
		if responsibleUnit.ActionPlans[i].TargetDate != nil {
			formattedTargetDate := responsibleUnit.ActionPlans[i].TargetDate.Format("2006-01-02")
			actionPlan.TargetDateString = &formattedTargetDate
			dto.ActionPlans[i] = actionPlan
		}
	}

	return &dto, nil
}
