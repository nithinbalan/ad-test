package services

import (
	auditInsightsRepositories "jethur/grc/audit-insights/repositories"
	auditFindingDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsDraftFindingActionPlanServiceInterface interface {
	GetAll(findingID uint, currentUserID uint, params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditDraftFindingActionPlan, error)
	GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditDraftFindingActionPlanDetailDTO, error)
	Create(findingID uint, currentUserID uint, AuditInsightsDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, auditFinding *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) error
	Delete(findingID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type AuditInsightsDraftFindingActionPlanService struct {
	repo                           auditInsightsRepositories.AuditInsightsDraftFindingActionPlanRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	auditDraftFindingSvc           auditServices.AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc     auditServices.AuditDraftFindingStatusServiceInterface
	auditDraftFindingActionPlanSvc auditServices.AuditDraftFindingActionPlanServiceInterface
	userSvc                        organizationServices.UserServiceInterface
	auditDraftReportSvc              auditServices.AuditDraftReportServiceInterface
	auditDraftReportStatusSvc        auditServices.AuditDraftReportStatusServiceInterface
}

func NewAuditInsightsDraftFindingActionPlanService(
	repo auditInsightsRepositories.AuditInsightsDraftFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDraftFindingSvc auditServices.AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc auditServices.AuditDraftFindingStatusServiceInterface,
	auditDraftFindingActionPlanSvc auditServices.AuditDraftFindingActionPlanServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditDraftReportSvc              auditServices.AuditDraftReportServiceInterface,
	auditDraftReportStatusSvc        auditServices.AuditDraftReportStatusServiceInterface,
) AuditInsightsDraftFindingActionPlanServiceInterface {
	return &AuditInsightsDraftFindingActionPlanService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		auditDraftFindingSvc:           auditDraftFindingSvc,
		auditDraftFindingStatusSvc:     auditDraftFindingStatusSvc,
		auditDraftFindingActionPlanSvc: auditDraftFindingActionPlanSvc,
		userSvc:                        userSvc,
		auditDraftReportSvc:                        auditDraftReportSvc,
		auditDraftReportStatusSvc:                        auditDraftReportStatusSvc,
	}
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) GetAll(findingID uint, currentUserID uint, params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error) {
	orgUnitID, err := aidfaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidfaps.repo.GetAll(findingID, orgUnitID, params)
}

func (aidfaps AuditInsightsDraftFindingActionPlanService) GetByID(findingID uint, id uint) (*models.AuditDraftFindingActionPlan, error) {
	auditDraftFindingActionPlan, err := aidfaps.auditDraftFindingActionPlanSvc.GetByID(findingID, id)
	if err != nil {
		return nil, err
	}
	return auditDraftFindingActionPlan, nil
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditDraftFindingActionPlanDetailDTO, error) {
	auditDraftFindingActionPlan, err := aidfaps.auditDraftFindingActionPlanSvc.GetFullDetailsByID(findingID,id)
	if err != nil {
		return nil, err
	}
	return auditDraftFindingActionPlan, nil
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) Create(findingID uint, currentUserID uint, auditInsightsDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) (uint, error) {
	if err := aidfaps.validateAuditInsightsDraftFindingActionPlanDTO(findingID,currentUserID); err != nil {
		return 0, err
	}
	id, err := aidfaps.auditDraftFindingActionPlanSvc.Create(findingID, currentUserID, auditInsightsDraftFindingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) Update(findingID uint, id uint, currentUserID uint, auditInsightsDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) error {
	if err := aidfaps.validateAuditInsightsDraftFindingActionPlanDTO(findingID,currentUserID); err != nil {
		return err
	}
	err := aidfaps.auditDraftFindingActionPlanSvc.Update(findingID, id, currentUserID, auditInsightsDraftFindingActionPlanFormDTO)
	if err != nil {
		return err
	}
	return nil
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) Delete(findingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if err := aidfaps.validateAuditInsightsDraftFindingActionPlanDTO(findingID,currentUserID); err != nil {
		return nil, err
	}
	failedIds, err := aidfaps.auditDraftFindingActionPlanSvc.Delete(findingID, ids, currentUserID)
	if err != nil {
		return nil, err
	}
	return failedIds, nil
}

func (aidfaps *AuditInsightsDraftFindingActionPlanService) validateAuditInsightsDraftFindingActionPlanDTO(findingID uint, currentUserID uint) error{
	existingAuditDraftFinding, err := aidfaps.auditDraftFindingSvc.GetFullDetailsByID(findingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, findingID)
	}
	orgUnitID, err := aidfaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}

	if *existingAuditDraftFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("Action not allowed:  Only users from the responsible unit can perform this action.")
	}

	// Retrieve the AuditDraftReport using AuditID from AuditDraftFinding.AuditChecklist.Audit
	auditID :=existingAuditDraftFinding.AuditID
	auditDraftReport, err := aidfaps.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}

	// Check if AuditDraftReportStatus is 'InDepartmentReview'
	inReviewStatus, err := aidfaps.auditDraftReportStatusSvc.GetByType("in-department-review")
	if err != nil {
		return err
	}
	if auditDraftReport.AuditDraftReportStatusID != inReviewStatus.ID {
		return utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'In Department Review' status can be perform this action.")
	}
	return nil
}