package services

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
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
	repo                  repositories.ComplianceReviewRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
}

func NewComplianceReviewService(
	repo repositories.ComplianceReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
) ComplianceReviewServiceInterface {
	return &ComplianceReviewService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		complianceRegisterSvc: complianceRegisterSvc,
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
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
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

	// @review
	// @we have to update next review date of compliance register.``
	// add transaction

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceReview{}.TableName(),
			id,
			complianceReview,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
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

	// Update Compliance Review
	existingComplianceReview.Comment = complianceReviewFormDTO.Comment
	existingComplianceReview.ReviewedBy = currentUserID
	//@review
	// we have to provide updated by and updated at right?

	err = crs.repo.Update(existingComplianceReview)

	if err != nil {
		return err
	}

	return nil
}
