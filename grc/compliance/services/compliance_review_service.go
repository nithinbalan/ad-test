package services

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ComplianceReviewServiceInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error)
	Create(complianceID uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, ComplianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) error
}

// Service and implementations
type ComplianceReviewService struct {
	repo                      repositories.ComplianceReviewRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	complianceRegisterSvc     ComplianceRegisterServiceInterface
	reviewFrequencySvc        masterServices.ReviewFrequencyServiceInterface
	complianceReviewStatusSvc ComplianceReviewStatusServiceInterface
}

func NewComplianceReviewService(
	repo repositories.ComplianceReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	complianceReviewStatusSvc ComplianceReviewStatusServiceInterface,
) ComplianceReviewServiceInterface {
	return &ComplianceReviewService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		complianceRegisterSvc:     complianceRegisterSvc,
		reviewFrequencySvc:        reviewFrequencySvc,
		complianceReviewStatusSvc: complianceReviewStatusSvc,
	}
}

func (crs *ComplianceReviewService) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	return crs.repo.GetAll(complianceID, sort, sortBy, keyword)
}

func (crs ComplianceReviewService) GetByID(complianceID uint, id uint) (*models.ComplianceReview, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceReview, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceReview, nil
}

func (crs *ComplianceReviewService) Create(complianceID uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) (uint, error) {
	complianceRegister, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return 0, err
	}

	complianceReview := &models.ComplianceReview{
		ComplianceRegisterID: complianceID,
		Comment:              complianceReviewFormDTO.Comment,
		ReviewedBy:           currentUserID,
		ReviewedAt:           time.Now(),
	}

	id, err := crs.repo.Create(complianceReview)
	if err != nil {
		return 0, err
	}

	if complianceRegister.ReviewFrequencyID != nil {
		reviewFrequency, err := crs.reviewFrequencySvc.GetByID(*complianceRegister.ReviewFrequencyID)
		if err != nil {
			return 0, err
		}
		nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)
		complianceRegister.NextReviewDate = &nextReviewDate
	}

	complianceReviewStatus, err := crs.complianceReviewStatusSvc.GetByType("reviewed")
	if err != nil {
		return 0, err
	}

	currentTime := time.Now()

	complianceRegister.LastReviewedDate = &currentTime
	complianceRegister.ComplianceReviewStatusID = &complianceReviewStatus.ID

	if err := crs.complianceRegisterSvc.ReviewStatusUpdate(complianceRegister); err != nil {
		return 0, err
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceReview{}.TableName(),
			id,
			complianceReview,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceReview{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceReviewService) Update(complianceID uint, id uint, currentUserID uint, complianceReviewFormDTO *complianceDTOs.ComplianceReviewFormDTO) error {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceReview, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	existingComplianceReview.Comment = complianceReviewFormDTO.Comment
	existingComplianceReview.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceReview)

	if err != nil {
		return err
	}

	return nil
}
