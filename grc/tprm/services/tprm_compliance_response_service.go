package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TPRMComplianceResponseServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TPRMComplianceResponse, error)
	GetByID(id uint) (*models.TPRMComplianceResponse, error)
	GetFullDetailsByID(id uint) (*dtos.TPRMComplianceResponseDetailDTO, error)
	Create(currentUserID uint, complianceResponseFormDTO *dtos.TPRMComplianceResponseFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.TPRMComplianceResponseFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type TPRMComplianceResponseService struct {
	repo           repositories.TPRMComplianceResponseRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewTPRMComplianceResponseService(repo repositories.TPRMComplianceResponseRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) TPRMComplianceResponseServiceInterface {
	return &TPRMComplianceResponseService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (tcrs *TPRMComplianceResponseService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TPRMComplianceResponse, error) {
	return tcrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (tcrs TPRMComplianceResponseService) GetByID(id uint) (*models.TPRMComplianceResponse, error) {
	complianceResponse, err := tcrs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return complianceResponse, nil

}

func (tcrs *TPRMComplianceResponseService) GetFullDetailsByID(id uint) (*dtos.TPRMComplianceResponseDetailDTO, error) {
	// Initialize DTO
	var complianceResponseDTO dtos.TPRMComplianceResponseDetailDTO

	// Get from repo
	complianceResponse, err := tcrs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map to DTO
	err = mapstructure.Decode(complianceResponse, &complianceResponseDTO)
	if err != nil {
		return nil, err
	}
	complianceResponseDTO.CreatedAt = complianceResponse.CreatedAt
	complianceResponseDTO.UpdatedAt = complianceResponse.UpdatedAt

	return &complianceResponseDTO, nil
}

func (tcrs *TPRMComplianceResponseService) Create(currentUserID uint, dto *dtos.TPRMComplianceResponseFormDTO) (uint, error) {
	tx := tcrs.repo.BeginTransaction()

	// Check for unique title
	if exists, err := tcrs.repo.ExistsByTitle(dto.TitleEn, dto.TitleAr); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	if exists, err := tcrs.repo.ExistsByScore(dto.Score); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the score must be unique"), "unique", "score")
	}

	// Build model
	model := &models.TPRMComplianceResponse{
		TitleEn:     dto.TitleEn,
		TitleAr:     dto.TitleAr,
		Score:       dto.Score,
		ColorCode:   &dto.ColorCode,
		CreatedBYID: currentUserID,
	}

	// Create
	id, err := tcrs.repo.Create(model, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	// Async activity log
	go func() {
		fullDetails, err := tcrs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details", map[string]interface{}{"error": err.Error()})
			return
		}

		err = tcrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TPRMComplianceResponse{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TPRMComplianceResponse{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (tcrs *TPRMComplianceResponseService) Update(id uint, currentUserID uint, dto *dtos.TPRMComplianceResponseFormDTO) error {
	existing, err := tcrs.GetByID(id)
	if err != nil {
		return err
	}

	// Check title uniqueness if changed
	if dto.TitleEn != existing.TitleEn || dto.TitleAr != existing.TitleAr {
		if exists, err := tcrs.repo.ExistsByTitleExcludingID(id, dto.TitleEn, dto.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	if dto.Score != existing.Score {
		if exists, err := tcrs.repo.ExistsByScoreExcludingID(id, dto.Score); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the score must be unique"), "unique", "score")
		}
	}

	// Update fields
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.Score = dto.Score
	existing.ColorCode = &dto.ColorCode
	existing.UpdatedBYID = &currentUserID

	// Save
	err = tcrs.repo.Update(existing)
	if err != nil {
		return err
	}

	// Log activity asynchronously
	go func() {
		fullDetails, err := tcrs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch strategic objective details for activity log", map[string]interface{}{"error": err.Error()})
			return
		}

		err = tcrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.TPRMComplianceResponse{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log update activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.TPRMComplianceResponse{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (tcrs *TPRMComplianceResponseService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := tcrs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = tcrs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := tcrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.TPRMComplianceResponse{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (tcrs *TPRMComplianceResponseService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := tcrs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tcrs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.TPRMComplianceResponse{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (tcrs *TPRMComplianceResponseService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := tcrs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := tcrs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.TPRMComplianceResponse{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
