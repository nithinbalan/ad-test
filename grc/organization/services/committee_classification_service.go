package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"
)

// Interface
type CommitteeClassificationServiceInterface interface {
	GetAll(params *dtos.CommitteeClassificationListQueryParamsDTO) ([]*dtos.CommitteeClassificationListDTO, error)
	GetByID(id uint) (*models.CommitteeClassification, error)
	Create(currentUserID uint, ccDTO *dtos.CommitteeClassificationFormDTO) (uint, error)
	Update(id uint, currentUserID uint, ccDTO *dtos.CommitteeClassificationFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Implementation
type CommitteeClassificationService struct {
	repo           repositories.CommitteeClassificationRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewCommitteeClassificationService(repo repositories.CommitteeClassificationRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) CommitteeClassificationServiceInterface {
	return &CommitteeClassificationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ccs *CommitteeClassificationService) GetAll(params *dtos.CommitteeClassificationListQueryParamsDTO) ([]*dtos.CommitteeClassificationListDTO, error) {
	return ccs.repo.GetAll(params)
}

func (ccs *CommitteeClassificationService) GetByID(id uint) (*models.CommitteeClassification, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *CommitteeClassificationService) Create(currentUserID uint, ccDTO *dtos.CommitteeClassificationFormDTO) (uint, error) {
	if err := ccs.validateCommitteeClassificationFormDTO(ccDTO, nil); err != nil {
		return 0, err
	}

	committeeClassification := &models.CommitteeClassification{
		TitleEn:       ccDTO.TitleEn,
		TitleAr:       ccDTO.TitleAr,
		DescriptionEn: ccDTO.DescriptionEn,
		DescriptionAr: ccDTO.DescriptionAr,
		Status:        ccDTO.Status,
	}

	id, err := ccs.repo.Create(committeeClassification)
	if err != nil {
		return 0, err
	}

	// Log activity
	go func() {
		createdClassification, err := ccs.repo.GetByID(id)
		if err != nil {
			return
		}
		err = ccs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CommitteeClassification{}.TableName(),
			id,
			createdClassification,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (ccs *CommitteeClassificationService) Update(id uint, currentUserID uint, ccDTO *dtos.CommitteeClassificationFormDTO) error {
	existingCC, err := ccs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := ccs.validateCommitteeClassificationFormDTO(ccDTO, &id); err != nil {
		return err
	}

	updatedClassification := &models.CommitteeClassification{
		TitleEn:       ccDTO.TitleEn,
		TitleAr:       ccDTO.TitleAr,
		DescriptionEn: ccDTO.DescriptionEn,
		DescriptionAr: ccDTO.DescriptionAr,
		Status:        ccDTO.Status,
	}

	existingCC.TitleEn = ccDTO.TitleEn
	existingCC.TitleAr = ccDTO.TitleAr
	existingCC.DescriptionEn = ccDTO.DescriptionEn
	existingCC.DescriptionAr = ccDTO.DescriptionAr
	existingCC.Status = ccDTO.Status

	err = ccs.repo.Update(existingCC)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingCC, updatedClassification, fieldsToCompare)
	if err != nil {
		log.Println("compare error:", err)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("json error:", err)
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CommitteeClassification{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return nil
}

func (ccs *CommitteeClassificationService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = ccs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CommitteeClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activity: %v", err)
			}
		}()
	}

	return failedIds, nil
}

func (ccs *CommitteeClassificationService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.CommitteeClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ccs *CommitteeClassificationService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ccs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.CommitteeClassification{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				log.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ccs *CommitteeClassificationService) validateCommitteeClassificationFormDTO(ccDTO *dtos.CommitteeClassificationFormDTO, excludeID *uint) error {
	// Check for duplicate titles
	if excludeID != nil {
		exists, err := ccs.repo.ExistsByTitleExcludingID(*excludeID, ccDTO.TitleEn, ccDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("committee classification with this title already exists"), "duplicate", "title")
		}
	} else {
		exists, err := ccs.repo.ExistsByTitle(ccDTO.TitleEn, ccDTO.TitleAr)
		if err != nil {
			return err
		}
		if exists {
			return utility.NewValidationError(fmt.Errorf("committee classification with this title already exists"), "duplicate", "title")
		}
	}

	return nil
}
