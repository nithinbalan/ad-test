package services

import (
	"errors"
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

func (es *EscalationService) GetAll(params dtos.EscalationListQueryParamsDTO) ([]*models.Escalation, error) {
	return es.repo.GetAll(params)
}

func (es EscalationService) GetByID(id uint) (*models.Escalation, error) {
	escalation, err := es.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return escalation, nil

}

func (res *EscalationService) GetFullDetailsByID(id uint) (*dtos.EscalationDetailsDTO, error) {
	var escalationDetailDTO dtos.EscalationDetailsDTO

	escalation, err := res.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(escalation, &escalationDetailDTO)
	if err != nil {
		return nil, err
	}

	return &escalationDetailDTO, nil
}

func (es *EscalationService) GetAllEnabledModules() ([]*dtos.ModuleListDTO, error) {
	return es.repo.GetAllEnabledModules()
}

func (es *EscalationService) GetAllEscalationsByModuleID(moduleID uint) ([]*dtos.GroupedEscalationDTO, error) {
	escalations, err := es.repo.GetEscalationByModuleID(moduleID)
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
