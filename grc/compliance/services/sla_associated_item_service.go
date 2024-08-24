package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type SLAAssociatedItemServiceInterface interface {
	GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error)
	GetAllDocuments(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemDocumentListResponseDTO, error)
	GetAllIncidents(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemRiskListResponseDTO, error)
	GetAllAssociatedCounts(slaID uint) (*complianceDTOs.SlaAssociatedItemsCountResponseDTO, error)
}

// Service and implementations
type SLAAssociatedItemService struct {
	repo           complianceRepositories.SLAAssociatedItemRepositoryInterface
	slaRegisterSvc SLARegisterServiceInterface
}

func NewSLAAssociatedItemService(
	repo complianceRepositories.SLAAssociatedItemRepositoryInterface,
	slaRegisterSvc SLARegisterServiceInterface,

) SLAAssociatedItemServiceInterface {
	return &SLAAssociatedItemService{
		repo:           repo,
		slaRegisterSvc: slaRegisterSvc,
	}
}

func (sais *SLAAssociatedItemService) GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error) {
	if _, err := sais.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return sais.repo.GetAllProcesses(slaID, params)
}

func (sais *SLAAssociatedItemService) GetAllDocuments(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemDocumentListResponseDTO, error) {
	if _, err := sais.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return sais.repo.GetAllDocuments(slaID, params)
}

func (sais *SLAAssociatedItemService) GetAllIncidents(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemIncidentListResponseDTO, error) {
	if _, err := sais.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return sais.repo.GetAllIncidents(slaID, params)
}

func (sais *SLAAssociatedItemService) GetAllRisks(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemRiskListResponseDTO, error) {
	if _, err := sais.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return sais.repo.GetAllRisks(slaID, params)
}

func (sais *SLAAssociatedItemService) GetAllAssociatedCounts(slaID uint) (*complianceDTOs.SlaAssociatedItemsCountResponseDTO, error) {
	if _, err := sais.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	response, err := sais.repo.GetAllAssociatedCounts(slaID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
