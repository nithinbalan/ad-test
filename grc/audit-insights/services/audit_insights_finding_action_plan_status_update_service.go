package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	"jethur/utility"
)

type AuditInsightsFindingActionPlanStatusUpdateServiceInterface interface {
	GetAll(actionPlanID uint, sort string, sortBy string, keyword string, tx ...utility.DatabaseTransaction) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	Create(actionPlanID uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
	Update(actionPlanID uint, id uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error
	GetLatestAuditInsightsFindingActionPlanStatusUpdate(actionPlanID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
}

type AuditInsightsFindingActionPlanStatusUpdateService struct {
	auditInsightsFindingActionPlanSvc     AuditInsightsFindingActionPlanServiceInterface
	auditFindingActionPlanStatusUpdateSvc auditServices.AuditFindingActionPlanStatusUpdateServiceInterface
}

func NewAuditInsightsFindingActionPlanStatusUpdateService(
	auditInsightsFindingActionPlanSvc AuditInsightsFindingActionPlanServiceInterface,
	auditFindingActionPlanStatusUpdateSvc auditServices.AuditFindingActionPlanStatusUpdateServiceInterface,
) AuditInsightsFindingActionPlanStatusUpdateServiceInterface {
	return &AuditInsightsFindingActionPlanStatusUpdateService{
		auditInsightsFindingActionPlanSvc:     auditInsightsFindingActionPlanSvc,
		auditFindingActionPlanStatusUpdateSvc: auditFindingActionPlanStatusUpdateSvc,
	}
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetAll(actionPlanID uint, sort string, sortBy string, keyword string, tx ...utility.DatabaseTransaction) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	auditInsightsFindingActionPlan, err := aifapsus.auditInsightsFindingActionPlanSvc.GetFullDetailsByID(actionPlanID, tx...)
	if err != nil {
		return nil, err
	}
	findingID := auditInsightsFindingActionPlan.AuditFinding.ID
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.GetAll(findingID, actionPlanID, sort, sortBy, keyword)
	if err != nil {
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) Create(actionPlanID uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
	auditInsightsFindingActionPlan, err := aifapsus.auditInsightsFindingActionPlanSvc.GetFullDetailsByID(actionPlanID, tx...)
	if err != nil {
		return 0, err
	}
	findingID := auditInsightsFindingActionPlan.AuditFinding.ID
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.Create(findingID, actionPlanID, currentUserID, actionPlanStatusUpdateFormDTO)
	if err != nil {
		return 0, err
	}

	return auditFindingActionPlanStatusUpdate, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) Update(actionPlanID uint, id uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error {
	auditInsightsFindingActionPlan, err := aifapsus.auditInsightsFindingActionPlanSvc.GetFullDetailsByID(actionPlanID, tx...)
	if err != nil {
		return err
	}
	findingID := auditInsightsFindingActionPlan.AuditFinding.ID
	return aifapsus.auditFindingActionPlanStatusUpdateSvc.Update(findingID, actionPlanID, id, currentUserID, actionPlanStatusUpdateFormDTO)

}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetLatestAuditInsightsFindingActionPlanStatusUpdate(actionPlanID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	auditInsightsFindingActionPlan, err := aifapsus.auditInsightsFindingActionPlanSvc.GetFullDetailsByID(actionPlanID, tx...)
	if err != nil {
		return nil, err
	}
	findingID := auditInsightsFindingActionPlan.AuditFinding.ID
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.GetLatestAuditFindingActionPlanStatusUpdate(findingID, actionPlanID)
	if err != nil {
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}
