package services

import (
	"errors"
	"fmt"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type EscalationServiceInterface interface {
	GetAll(params dtos.EscalationListQueryParamsDTO) ([]*models.Escalation, error)
	GetByID(id uint) (*models.Escalation, error)
	GetFullDetailsByID(id uint) (*dtos.EscalationDetailsDTO, error)

	GetAllEnabledModules() ([]*dtos.ModuleListDTO, error)
	GetAllEscalationsByModuleID(moduleID uint) ([]*dtos.GroupedEscalationDTO, error)
}

// Service and implementations
type EscalationService struct {
	repo repositories.EscalationRepositoryInterface
}

func NewEscalationService(repo repositories.EscalationRepositoryInterface) EscalationServiceInterface {
	return &EscalationService{
		repo: repo,
	}
}

func (rs *EscalationService) GetAll(params dtos.EscalationListQueryParamsDTO) ([]*models.Escalation, error) {
	return rs.repo.GetAll(params)
}

func (rs EscalationService) GetByID(id uint) (*models.Escalation, error) {
	escalation, err := rs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return escalation, nil

}

func (rrs *EscalationService) GetFullDetailsByID(id uint) (*dtos.EscalationDetailsDTO, error) {
	var escalationDetailDTO dtos.EscalationDetailsDTO

	escalation, err := rrs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	fmt.Println("23")

	err = mapstructure.Decode(escalation, &escalationDetailDTO)
	if err != nil {
		return nil, err
	}

	return &escalationDetailDTO, nil
}

func (rs *EscalationService) GetAllEnabledModules() ([]*dtos.ModuleListDTO, error) {
	return rs.repo.GetAllEnabledModules()
}

func (rs *EscalationService) GetAllEscalationsByModuleID(moduleID uint) ([]*dtos.GroupedEscalationDTO, error) {
	escalations, err := rs.repo.GetEscalationByModuleID(moduleID)
	if err != nil {
		return nil, err
	}

	escalationDTOs := []*dtos.GroupedEscalationDTO{}
	for _, escalation := range escalations {

		escalationDTO := &dtos.GroupedEscalationDTO{
			ID:               escalation.ID,
			TitleEn:          escalation.TitleEn,
			TitleAr:          escalation.TitleAr,
			EscalationLevels: []dtos.EscalationLevelDetailDTO{},
		}

		for _, rule := range escalation.EscalationLevels {

			escalationDTO.EscalationLevels = append(escalationDTO.EscalationLevels, dtos.EscalationLevelDetailDTO{
				ID:                 rule.ID,
				NotificationMethod: rule.NotificationMethod,
				EscalationTo:       rule.EscalationTo,
				DayOffset:          rule.DayOffset,
				CreatedAt:          rule.CreatedAt,
			})
		}

		escalationDTOs = append(escalationDTOs, escalationDTO)
	}
	return escalationDTOs, err
}
