package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	"jethur/utility"
)

type AuditInsightsFindingActionPlanStatusUpdateServiceInterface interface {
	GetAll(findingID uint, actionPlanID uint, sort string, sortBy string, keyword string, tx ...utility.DatabaseTransaction) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	Create(findingID uint, actionPlanID uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
	Update(findingID uint, actionPlanID uint, id uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error
	GetLatestAuditInsightsFindingActionPlanStatusUpdate(findingID uint, actionPlanID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	GetFindingIDByActionPlan(actionPlanID uint) (uint, error)
	GetDocumentPath(findingID uint, actionPlanID uint, id uint, token string, tx ...utility.DatabaseTransaction) (string, error)
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

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetAll(findingID uint, actionPlanID uint, sort string, sortBy string, keyword string, tx ...utility.DatabaseTransaction) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.GetAll(findingID, actionPlanID, sort, sortBy, keyword)
	if err != nil {
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) Create(findingID uint, actionPlanID uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.Create(findingID, actionPlanID, currentUserID, actionPlanStatusUpdateFormDTO)
	if err != nil {
		return 0, err
	}

	return auditFindingActionPlanStatusUpdate, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) Update(findingID uint, actionPlanID uint, id uint, currentUserID uint, actionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error {
	return aifapsus.auditFindingActionPlanStatusUpdateSvc.Update(findingID, actionPlanID, id, currentUserID, actionPlanStatusUpdateFormDTO)
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetLatestAuditInsightsFindingActionPlanStatusUpdate(findingID uint, actionPlanID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	auditFindingActionPlanStatusUpdate, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.GetLatestAuditFindingActionPlanStatusUpdate(findingID, actionPlanID)
	if err != nil {
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetFindingIDByActionPlan(actionPlanID uint) (uint, error) {
	auditInsightsFindingActionPlan, err := aifapsus.auditInsightsFindingActionPlanSvc.GetFullDetailsByID(actionPlanID)
	if err != nil {
		return 0, err
	}
	return auditInsightsFindingActionPlan.AuditFinding.ID, nil
}

func (aifapsus *AuditInsightsFindingActionPlanStatusUpdateService) GetDocumentPath(findingID uint, actionPlanID uint, id uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	auditDocumentRequest, err := aifapsus.auditFindingActionPlanStatusUpdateSvc.GetDocumentPath(findingID, actionPlanID, id, token)
	if err != nil {
		return "", err
	}

	return auditDocumentRequest, nil
}
