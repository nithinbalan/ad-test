package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type OrganizationUnitBusinessContextServiceInterface interface {
	GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitBusinessContextListQueryParamsDTO) (*organizationDTOs.OrganizationUnitBusinessContextListResponseDTO, error)
	GetByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error)
	GetFullDetailsByID(departmentID uint) (*organizationDTOs.OrganizationUnitBusinessContextDetailDTO, error)
	Create(departmentID uint, currentUserID uint, contextForm *organizationDTOs.OrganizationUnitBusinessContextFormDTO) (uint, error)
	Update(departmentID uint, currentUserID uint, contextForm *organizationDTOs.OrganizationUnitBusinessContextFormDTO) error
	Delete(departmentID uint, currentUserID uint) error
}

type OrganizationUnitBusinessContextService struct {
	repo           repositories.OrganizationUnitBusinessContextRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewOrganizationUnitBusinessContextService(
	repo repositories.OrganizationUnitBusinessContextRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) OrganizationUnitBusinessContextServiceInterface {
	return &OrganizationUnitBusinessContextService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *OrganizationUnitBusinessContextService) GetAll(departmentID uint, params *organizationDTOs.OrganizationUnitBusinessContextListQueryParamsDTO) (*organizationDTOs.OrganizationUnitBusinessContextListResponseDTO, error) {
	response, err := s.repo.GetAll(departmentID, params)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *OrganizationUnitBusinessContextService) GetByID(departmentID uint) (*models.OrganizationUnitBusinessContext, error) {
	context, err := s.repo.GetByID(departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, departmentID)
		}
		return nil, err
	}

	return context, nil
}

func (s *OrganizationUnitBusinessContextService) GetFullDetailsByID(departmentID uint) (*organizationDTOs.OrganizationUnitBusinessContextDetailDTO, error) {
	context, err := s.repo.GetFullDetailsByID(departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, departmentID)
		}
		return nil, err
	}

	var contextDetailDTO organizationDTOs.OrganizationUnitBusinessContextDetailDTO
	err = mapstructure.Decode(context, &contextDetailDTO)
	if err != nil {
		return nil, err
	}

	return &contextDetailDTO, nil
}

func (s *OrganizationUnitBusinessContextService) Create(departmentID uint, currentUserID uint, contextForm *organizationDTOs.OrganizationUnitBusinessContextFormDTO) (uint, error) {
	exists, err := s.repo.ExistsByDepartmentID(departmentID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, utility.NewValidationError(fmt.Errorf("business context already exists for this department"), "unique", "department_id")
	}

	now := time.Now()
	context := &models.OrganizationUnitBusinessContext{
		DepartmentID:               departmentID,
		InternalAuditNote:          contextForm.InternalAuditNote,
		ExternalContext:            contextForm.ExternalContext,
		AnnualObjectiveDescription: contextForm.AnnualObjectiveDescription,
		CreatedBYID:                currentUserID,
		CreatedAt:                  now,
	}

	createdID, err := s.repo.Create(context)
	if err != nil {
		return 0, err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.OrganizationUnitBusinessContext{}.TableName(),
			createdID,
			nil,
			currentUserID,
			"organization", "organization_unit", "create_business_context", nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.OrganizationUnitBusinessContext{}.TableName(),
				"recordID":  createdID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return createdID, nil
}

func (s *OrganizationUnitBusinessContextService) Update(departmentID uint, currentUserID uint, contextForm *organizationDTOs.OrganizationUnitBusinessContextFormDTO) error {
	_, err := s.repo.GetByID(departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, departmentID)
		}
		return err
	}

	now := time.Now()
	updatedContext := &models.OrganizationUnitBusinessContext{
		InternalAuditNote:          contextForm.InternalAuditNote,
		ExternalContext:            contextForm.ExternalContext,
		AnnualObjectiveDescription: contextForm.AnnualObjectiveDescription,
		UpdatedBYID:                &currentUserID,
		UpdatedAt:                  &now,
	}

	err = s.repo.Update(departmentID, updatedContext)
	if err != nil {
		return err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.OrganizationUnitBusinessContext{}.TableName(),
			departmentID,
			nil,
			currentUserID,
			"organization", "organization_unit", "update_business_context", nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.OrganizationUnitBusinessContext{}.TableName(),
				"recordID":  departmentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (s *OrganizationUnitBusinessContextService) Delete(departmentID uint, currentUserID uint) error {
	_, err := s.repo.GetByID(departmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, departmentID)
		}
		return err
	}

	err = s.repo.Delete(departmentID)
	if err != nil {
		return err
	}

	go func() {
		err := s.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.OrganizationUnitBusinessContext{}.TableName(),
			departmentID,
			nil,
			currentUserID,
			"organization", "organization_unit", "delete_business_context", nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.OrganizationUnitBusinessContext{}.TableName(),
				"recordID":  departmentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}
