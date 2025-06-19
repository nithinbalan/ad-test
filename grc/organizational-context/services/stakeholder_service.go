package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	stakeholderRepositories "jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderServiceInterface interface {
	GetAll(params *organizationalContextDTOs.StakeholderListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderDetailDTO, error)
	GetByID(id uint) (*models.Stakeholder, error)
	Create(currentUserID uint, organizationalContextDTOs *organizationalContextDTOs.StakeholderFormDTO) (uint, error)
	Update(id uint, currentUserID uint, stakeholder *organizationalContextDTOs.StakeholderFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StakeholderService struct {
	repo           stakeholderRepositories.StakeholderRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewStakeholderService(
	repo stakeholderRepositories.StakeholderRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) StakeholderServiceInterface {
	return &StakeholderService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (ss *StakeholderService) GetAll(params *organizationalContextDTOs.StakeholderListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderListDTO, error) {
	return ss.repo.GetAll(params)
}

func (ss StakeholderService) GetByID(id uint) (*models.Stakeholder, error) {
	stakeholder, err := ss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return stakeholder, nil
}

func (ss StakeholderService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderDetailDTO, error) {
	// Initialize your DTO
	var stakeholderDTO organizationalContextDTOs.StakeholderDetailDTO

	stakeholder, err := ss.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(stakeholder, &stakeholderDTO)
	if err != nil {
		return nil, err
	}

	stakeholderDTO.CreatedAt = stakeholder.CreatedAt
	stakeholderDTO.UpdatedAt = stakeholder.UpdatedAt

	return &stakeholderDTO, nil
}

func (ss *StakeholderService) Create(currentUserID uint, stakeholderFormDTO *organizationalContextDTOs.StakeholderFormDTO) (uint, error) {

	tx := ss.repo.BeginTransaction()

	stakeholder := &models.Stakeholder{
		TitleEn:               stakeholderFormDTO.TitleEn,
		TitleAr:               stakeholderFormDTO.TitleAr,
		TypeID:                stakeholderFormDTO.TypeID,
		StakeholderCategoryID: stakeholderFormDTO.StakeholderCategoryID,
		InfluenceLevelID:      stakeholderFormDTO.InfluenceLevelID,
		ImpactLevelID:         stakeholderFormDTO.ImpactLevelID,
		Role:                  stakeholderFormDTO.Role,
		ContactInformation:    stakeholderFormDTO.ContactInformation,
		PrimaryUnitID:         stakeholderFormDTO.OrganizationUnitID,
		ThirdPartyRegisterID:  stakeholderFormDTO.ThirdPartyRegisterID,
		ClientID:              stakeholderFormDTO.ClientID,
		UserID:                stakeholderFormDTO.UserID,
		CreatedBYID:           currentUserID,
	}
	id, err := ss.repo.Create(stakeholder, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStakeholder, err := ss.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch swot pestel question factor", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = ss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Stakeholder{}.TableName(),
			id,
			existingStakeholder,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Stakeholder{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (ss *StakeholderService) Update(id uint, currentUserID uint, stakeholderFormDTO *organizationalContextDTOs.StakeholderFormDTO) error {
	tx := ss.repo.BeginTransaction()

	existingStakeholder, err := ss.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStakeholder, err := ss.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStakeholder.TitleEn = stakeholderFormDTO.TitleEn
	existingStakeholder.TitleAr = stakeholderFormDTO.TitleAr
	existingStakeholder.TypeID = stakeholderFormDTO.TypeID
	existingStakeholder.StakeholderCategoryID = stakeholderFormDTO.StakeholderCategoryID
	existingStakeholder.InfluenceLevelID = stakeholderFormDTO.InfluenceLevelID
	existingStakeholder.ImpactLevelID = stakeholderFormDTO.ImpactLevelID
	existingStakeholder.Role = stakeholderFormDTO.Role
	existingStakeholder.ContactInformation = stakeholderFormDTO.ContactInformation
	existingStakeholder.PrimaryUnitID = stakeholderFormDTO.OrganizationUnitID
	existingStakeholder.ThirdPartyRegisterID = stakeholderFormDTO.ThirdPartyRegisterID
	existingStakeholder.ClientID = stakeholderFormDTO.ClientID
	existingStakeholder.UserID = stakeholderFormDTO.UserID
	existingStakeholder.UpdatedBYID = &currentUserID

	err = ss.repo.Update(existingStakeholder, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStakeholder, err := ss.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title_en": oldStakeholder.TitleEn,
			"title_ar": oldStakeholder.TitleAr,
		}

		newData := map[string]interface{}{
			"title_en": existingStakeholder.TitleEn,
			"title_ar": existingStakeholder.TitleAr,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ss.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Stakeholder{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.Stakeholder{}.TableName(),
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

func (ss *StakeholderService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ss.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStakeholder []map[string]interface{}

	for _, id := range ids {
		stakeholder, err := ss.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStakeholder = append(deletedStakeholder, map[string]interface{}{
			"id":       id,
			"title_en": stakeholder.TitleEn,
		})

		// Delete the business
		if err := ss.repo.Delete(id, tx); err != nil {
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
			err := ss.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.Stakeholder{}.TableName(),
				successfulIds,
				deletedStakeholder,
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
