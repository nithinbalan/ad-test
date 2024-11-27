package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	"jethur/grc/models"
)

type AuditInsightsFindingRootCauseAnalysisServiceInterface interface {
	GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error)
	Create(findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error
	Delete(findingID uint, id uint, currentUserID uint) (uint, error)
}

type AuditInsightsFindingRootCauseAnalysisService struct {
	auditFindingRootCauseAnalysisSvc auditServices.AuditFindingRootCauseAnalysisServiceInterface
}

func NewAuditInsightsFindingRootCauseAnalysisService(
	auditFindingRootCauseAnalysisSvc auditServices.AuditFindingRootCauseAnalysisServiceInterface,
) AuditInsightsFindingRootCauseAnalysisServiceInterface {
	return &AuditInsightsFindingRootCauseAnalysisService{
		auditFindingRootCauseAnalysisSvc: auditFindingRootCauseAnalysisSvc,
	}
}

func (aiafrcas *AuditInsightsFindingRootCauseAnalysisService) GetAll(findingID uint, params *auditDTOs.AuditFindingRootCauseAnalysisListQueryParamsDTO) ([]*auditDTOs.AuditFindingRootCauseAnalysisListDTO, error) {
	auditFindingRootCauseAnalysis, err := aiafrcas.auditFindingRootCauseAnalysisSvc.GetAll(findingID, params)
	if err != nil {
		return nil, err
	}
	return auditFindingRootCauseAnalysis, nil
}

func (aiafrcas *AuditInsightsFindingRootCauseAnalysisService) GetByID(findingID uint, id uint) (*models.AuditFindingRootCauseAnalysis, error) {
	auditFindingRootCauseAnalysis, err := aiafrcas.auditFindingRootCauseAnalysisSvc.GetByID(findingID, id)
	if err != nil {
		return nil, err
	}
	return auditFindingRootCauseAnalysis, nil
}

func (aiafrcas *AuditInsightsFindingRootCauseAnalysisService) Create(findingID uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) (uint, error) {
	auditFindingRootCauseAnalysis, err := aiafrcas.auditFindingRootCauseAnalysisSvc.Create(findingID, currentUserID, rootCauseAnalysisFormDTO)
	if err != nil {
		return 0, err
	}

	return auditFindingRootCauseAnalysis, nil
}

func (aiafrcas *AuditInsightsFindingRootCauseAnalysisService) Update(findingID uint, id uint, currentUserID uint, rootCauseAnalysisFormDTO *auditDTOs.AuditFindingRootCauseAnalysisFormDTO) error {
	return aiafrcas.auditFindingRootCauseAnalysisSvc.Update(findingID, id, currentUserID, rootCauseAnalysisFormDTO)

}

func (aiafrcas *AuditInsightsFindingRootCauseAnalysisService) Delete(findingID uint, id uint, currentUserID uint) (uint, error) {
	auditFindingRootCauseAnalysis, err := aiafrcas.auditFindingRootCauseAnalysisSvc.Delete(findingID, id, currentUserID)
	if err != nil {
		return 0, err
	}

	return auditFindingRootCauseAnalysis, nil
}
