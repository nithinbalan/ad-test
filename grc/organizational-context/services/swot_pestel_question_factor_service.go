package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	swotPestelQuestionFactorRepositories "jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SwotPestelQuestionFactorServiceInterface interface {
	GetAll(questionId uint, params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error)
	GetAllFactors(params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.SwotPestelQuestionFactorDetailDTO, error)
	GetByID(id uint) (*models.SwotPestelQuestionFactor, error)
	Create(questionId uint, currentUserID uint, organizationalContextDTOs *organizationalContextDTOs.SwotPestelQuestionFactorFormDTO) (uint, error)
	Update(id uint, currentUserID uint, swotPestelQuestionFactor *organizationalContextDTOs.SwotPestelQuestionFactorFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SwotPestelQuestionFactorService struct {
	repo           swotPestelQuestionFactorRepositories.SwotPestelQuestionFactorRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewSwotPestelQuestionFactorService(
	repo swotPestelQuestionFactorRepositories.SwotPestelQuestionFactorRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) SwotPestelQuestionFactorServiceInterface {
	return &SwotPestelQuestionFactorService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (spqfs *SwotPestelQuestionFactorService) GetAll(questionId uint, params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error) {
	return spqfs.repo.GetAll(questionId, params)
}

func (spqfs *SwotPestelQuestionFactorService) GetAllFactors(params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error) {
	return spqfs.repo.GetAllFactors(params)
}

func (spqfs SwotPestelQuestionFactorService) GetByID(id uint) (*models.SwotPestelQuestionFactor, error) {
	swotPestelQuestionFactor, err := spqfs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return swotPestelQuestionFactor, nil
}

func (spqfs SwotPestelQuestionFactorService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.SwotPestelQuestionFactorDetailDTO, error) {
	// Initialize your DTO
	var swotPestelQuestionFactorDTO organizationalContextDTOs.SwotPestelQuestionFactorDetailDTO

	swotPestelQuestionFactor, err := spqfs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(swotPestelQuestionFactor, &swotPestelQuestionFactorDTO)
	if err != nil {
		return nil, err
	}

	swotPestelQuestionFactorDTO.CreatedAt = swotPestelQuestionFactor.CreatedAt
	swotPestelQuestionFactorDTO.UpdatedAt = swotPestelQuestionFactor.UpdatedAt

	return &swotPestelQuestionFactorDTO, nil
}

func (spqfs *SwotPestelQuestionFactorService) Create(questionId uint, currentUserID uint, swotPestelQuestionFactorFormDTO *organizationalContextDTOs.SwotPestelQuestionFactorFormDTO) (uint, error) {

	tx := spqfs.repo.BeginTransaction()

	swotPestelQuestionFactor := &models.SwotPestelQuestionFactor{
		SwotPestelQuestionID: questionId,
		SwotPestelCategoryID: swotPestelQuestionFactorFormDTO.SwotPestelCategoryID,
		Description:          *swotPestelQuestionFactorFormDTO.Description,
		CreatedBYID:          currentUserID,
	}
	id, err := spqfs.repo.Create(swotPestelQuestionFactor, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingSwotPestelQuestionFactor, err := spqfs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch swot pestel question factor", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = spqfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SwotPestelQuestionFactor{}.TableName(),
			id,
			existingSwotPestelQuestionFactor,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SwotPestelQuestionFactor{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (spqfs *SwotPestelQuestionFactorService) Update(id uint, currentUserID uint, swotPestelQuestionFactorFormDTO *organizationalContextDTOs.SwotPestelQuestionFactorFormDTO) error {
	tx := spqfs.repo.BeginTransaction()

	existingSwotPestelQuestionFactor, err := spqfs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldSwotPestelQuestionFactor, err := spqfs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingSwotPestelQuestionFactor.SwotPestelCategoryID = swotPestelQuestionFactorFormDTO.SwotPestelCategoryID
	existingSwotPestelQuestionFactor.Description = *swotPestelQuestionFactorFormDTO.Description
	existingSwotPestelQuestionFactor.UpdatedBYID = &currentUserID

	err = spqfs.repo.Update(existingSwotPestelQuestionFactor, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingSwotPestelQuestionFactor, err := spqfs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"swot_pestel_question": oldSwotPestelQuestionFactor.SwotPestelQuestion.TitleEn,
			"swot_pestel_category": oldSwotPestelQuestionFactor.SwotPestelCategory.TitleEn,
			"description":          oldSwotPestelQuestionFactor.Description,
		}

		newData := map[string]interface{}{
			"swot_pestel_question": existingSwotPestelQuestionFactor.SwotPestelQuestion.TitleEn,
			"swot_pestel_category": existingSwotPestelQuestionFactor.SwotPestelCategory.TitleEn,
			"description":          existingSwotPestelQuestionFactor.Description,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := spqfs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SwotPestelQuestionFactor{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.SwotPestelQuestionFactor{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (spqfs *SwotPestelQuestionFactorService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := spqfs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedSwotPestelQuestionFactor []map[string]interface{}

	for _, id := range ids {
		swotPestelQuestionFactor, err := spqfs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedSwotPestelQuestionFactor = append(deletedSwotPestelQuestionFactor, map[string]interface{}{
			"id":          id,
			"description": swotPestelQuestionFactor.Description,
		})

		// Delete the business
		if err := spqfs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := spqfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.SwotPestelQuestionFactor{}.TableName(),
				successfulIds,
				deletedSwotPestelQuestionFactor,
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
