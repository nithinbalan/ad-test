package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
)

type ComplianceAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error)
}

// Service and implementations
type ComplianceAssociatedItemService struct {
	repo                  complianceRepositories.ComplianceAssociatedItemRepositoryInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
}

func NewComplianceAssociatedItemService(
	repo complianceRepositories.ComplianceAssociatedItemRepositoryInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
) ComplianceAssociatedItemServiceInterface {
	return &ComplianceAssociatedItemService{
		repo:                  repo,
		complianceRegisterSvc: complianceRegisterSvc,
	}
}

func (cais *ComplianceAssociatedItemService) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllControls(complianceID, params)
}
