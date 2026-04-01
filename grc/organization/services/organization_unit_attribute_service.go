package services

import (
	"errors"
	"time"

	services "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	repositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type OrganizationUnitAttributeServiceInterface interface {
	GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitAttributeListQueryParamsDTO) (*organizationDTOs.OrganizationUnitAttributeListResponseDTO, error)
	GetByID(id uint) (*models.OrganizationUnitAttribute, error)
	GetFullDetailsByID(id uint) (*organizationDTOs.OrganizationUnitAttributeDetailDTO, error)
	Create(departmentID uint, currentUserID uint, attributeForm *organizationDTOs.OrganizationUnitAttributeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, attributeForm *organizationDTOs.OrganizationUnitAttributeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type OrganizationUnitAttributeService struct {
	repo           repositories.OrganizationUnitAttributeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewOrganizationUnitAttributeService(
	repo repositories.OrganizationUnitAttributeRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) OrganizationUnitAttributeServiceInterface {
	return &OrganizationUnitAttributeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *OrganizationUnitAttributeService) GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitAttributeListQueryParamsDTO) (*organizationDTOs.OrganizationUnitAttributeListResponseDTO, error) {
	response, err := s.repo.GetAll(departmentID, params)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrganizationUnitAttributeService) GetByID(id uint) (*models.OrganizationUnitAttribute, error) {
	attribute, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return attribute, nil
}

func (s *OrganizationUnitAttributeService) GetFullDetailsByID(id uint) (*organizationDTOs.OrganizationUnitAttributeDetailDTO, error) {
	attribute, err := s.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var attributeDetailDTO organizationDTOs.OrganizationUnitAttributeDetailDTO
	err = mapstructure.Decode(attribute, &attributeDetailDTO)
	if err != nil {
		return nil, err
	}

	return &attributeDetailDTO, nil
}

func (s *OrganizationUnitAttributeService) Create(departmentID uint, currentUserID uint, attributeForm *organizationDTOs.OrganizationUnitAttributeFormDTO) (uint, error) {
	now := time.Now()
	attribute := &models.OrganizationUnitAttribute{
		DepartmentID: departmentID,
		Title:        attributeForm.Title,
		Description:  attributeForm.Description,
		CreatedBYID:  currentUserID,
		CreatedAt:    now,
	}

	createdID, err := s.repo.Create(attribute)
	if err != nil {
		return 0, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.OrganizationUnitAttribute{}.TableName(),
			createdID,
			nil,
			currentUserID,
			"organization", "organization_unit", "create_attribute", nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.OrganizationUnitAttribute{}.TableName(),
				"recordID":  createdID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return createdID, nil
}

func (s *OrganizationUnitAttributeService) Update(id uint, currentUserID uint, attributeForm *organizationDTOs.OrganizationUnitAttributeFormDTO) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	now := time.Now()
	updatedAttribute := &models.OrganizationUnitAttribute{
		Title:       attributeForm.Title,
		Description: attributeForm.Description,
		UpdatedBYID: &currentUserID,
		UpdatedAt:   &now,
	}

	err = s.repo.Update(id, updatedAttribute)
	if err != nil {
		return err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.OrganizationUnitAttribute{}.TableName(),
			id,
			nil,
			currentUserID,
			"organization", "organization_unit", "update_attribute", nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.OrganizationUnitAttribute{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (s *OrganizationUnitAttributeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		_, err := s.repo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIds = append(failedIds, id)
				continue
			}
			failedIds = append(failedIds, id)
			continue
		}
		successfulIds = append(successfulIds, id)
	}

	if len(successfulIds) > 0 {
		err := s.repo.Delete(successfulIds)
		if err != nil {
			return failedIds, err
		}

		go func() {
			for _, id := range successfulIds {
				err := s.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.OrganizationUnitAttribute{}.TableName(),
					id,
					nil,
					currentUserID,
					"organization", "organization_unit", "delete_attribute", nil,
				)

				if err != nil {
					utility.LogError("Failed to log activity", map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.OrganizationUnitAttribute{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					})
				}
			}
		}()
	}

	return failedIds, nil
}
