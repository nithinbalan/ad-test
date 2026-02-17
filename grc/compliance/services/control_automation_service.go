package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
)

type ControlAutomationServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ControlAutomationResponseDTO, error)
	GetByID(id uint) (*complianceDTOs.ControlAutomationResponseDTO, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ControlAutomationService struct {
	repo           repositories.ControlAutomationRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewControlAutomationService(repo repositories.ControlAutomationRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ControlAutomationServiceInterface {
	return &ControlAutomationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (cas *ControlAutomationService) GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ControlAutomationResponseDTO, error) {
	list, err := cas.repo.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}

	result := make([]*complianceDTOs.ControlAutomationResponseDTO, 0, len(list))
	for _, m := range list {
		result = append(result, mapControlAutomationToResponseDTO(m))
	}
	return result, nil
}

func (cas *ControlAutomationService) GetByID(id uint) (*complianceDTOs.ControlAutomationResponseDTO, error) {
	controlAutomation, err := cas.repo.GetByID(id)
	if err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}
	return mapControlAutomationToResponseDTO(controlAutomation), nil
}

func mapUserToUserMinimalDTO(user *models.User) organizationDTOs.UserMinimalDTO {
	if user == nil {
		return organizationDTOs.UserMinimalDTO{}
	}
	return organizationDTOs.UserMinimalDTO{
		ID:               int(user.ID),
		FirstNameEN:      user.FirstNameEn,
		LastNameEN:       user.LastNameEn,
		FirstNameAR:      user.FirstNameAr,
		LastNameAR:       user.LastNameAr,
		ProfileImagePath: user.ProfileImagePath,
		Email:            user.Email,
		Initial:          user.Initial,
		Color:            user.Color,
	}
}

func mapControlAutomationToResponseDTO(m *models.ControlAutomation) *complianceDTOs.ControlAutomationResponseDTO {
	if m == nil {
		return nil
	}

	dto := &complianceDTOs.ControlAutomationResponseDTO{
		ID:          m.ID,
		TitleEn:     m.TitleEn,
		TitleAr:     m.TitleAr,
		Type:        m.Type,
		Score:       m.Score,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		CreatedByID: m.CreatedBYID,
	}

	if m.CreatedByUser != nil {
		createdBy := mapUserToUserMinimalDTO(m.CreatedByUser)
		dto.CreatedByUser = &createdBy
	}

	return dto
}

func (cas *ControlAutomationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := cas.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cas.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ControlAutomation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
			}
		}()
	}

	return err
}

func (cas *ControlAutomationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := cas.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := cas.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ControlAutomation{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
			}
		}()
	}

	return err
}
