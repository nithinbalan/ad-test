package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
)

type ControlConsistencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ControlConsistencyResponseDTO, error)
	GetByID(id uint) (*complianceDTOs.ControlConsistencyResponseDTO, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type ControlConsistencyService struct {
	repo           repositories.ControlConsistencyRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewControlConsistencyService(repo repositories.ControlConsistencyRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) ControlConsistencyServiceInterface {
	return &ControlConsistencyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *ControlConsistencyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceDTOs.ControlConsistencyResponseDTO, error) {
	list, err := ccs.repo.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		return nil, err
	}
	result := make([]*complianceDTOs.ControlConsistencyResponseDTO, 0, len(list))
	for _, m := range list {
		result = append(result, mapControlConsistencyToResponseDTO(m))
	}
	return result, nil
}

func (ccs *ControlConsistencyService) GetByID(id uint) (*complianceDTOs.ControlConsistencyResponseDTO, error) {
	controlConsistency, err := ccs.repo.GetByID(id)
	if err != nil {
		if _, ok := err.(*utility.ResourceNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}
	return mapControlConsistencyToResponseDTO(controlConsistency), nil
}

func mapControlConsistencyUserToUserMinimalDTO(user *models.User) organizationDTOs.UserMinimalDTO {
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

func mapControlConsistencyToResponseDTO(m *models.ControlConsistency) *complianceDTOs.ControlConsistencyResponseDTO {
	if m == nil {
		return nil
	}
	dto := &complianceDTOs.ControlConsistencyResponseDTO{
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
		createdBy := mapControlConsistencyUserToUserMinimalDTO(m.CreatedByUser)
		dto.CreatedByUser = &createdBy
	}
	return dto
}

func (ccs *ControlConsistencyService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.ControlConsistency{}.TableName(),
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

func (ccs *ControlConsistencyService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.ControlConsistency{}.TableName(),
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


