package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type SLAAssociatedItemServiceInterface interface {
	GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error)
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
