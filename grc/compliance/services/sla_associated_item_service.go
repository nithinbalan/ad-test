package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type SLAAssociatedItemServiceInterface interface {
	GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error)
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
