package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	BcpCommunicationPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BcpCommunicationPlanServiceInterface interface {
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpCommunicationPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpCommunicationPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpCommunicationPlanDetailDTO, error)
	GetByID(id uint) (*models.BcpCommunicationPlan, error)
	Create(bcpId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BcpCommunicationPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, BcpCommunicationPlan *businessContinuityDTOs.BcpCommunicationPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BcpCommunicationPlanService struct {
	repo           BcpCommunicationPlanRepositories.BcpCommunicationPlanRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewBcpCommunicationPlanService(
	repo BcpCommunicationPlanRepositories.BcpCommunicationPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) BcpCommunicationPlanServiceInterface {
	return &BcpCommunicationPlanService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (bcps *BcpCommunicationPlanService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpCommunicationPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpCommunicationPlanListResponseDTO, error) {
	return bcps.repo.GetAll(bcpId, params)
}

func (bcps BcpCommunicationPlanService) GetByID(id uint) (*models.BcpCommunicationPlan, error) {
	BcpCommunicationPlan, err := bcps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return BcpCommunicationPlan, nil
}

func (bcps BcpCommunicationPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpCommunicationPlanDetailDTO, error) {
	// Initialize your DTO
	var BcpCommunicationPlanDTO businessContinuityDTOs.BcpCommunicationPlanDetailDTO

	BcpCommunicationPlan, err := bcps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(BcpCommunicationPlan, &BcpCommunicationPlanDTO)
	if err != nil {
		return nil, err
	}

	BcpCommunicationPlanDTO.CreatedAt = BcpCommunicationPlan.CreatedAt
	BcpCommunicationPlanDTO.UpdatedAt = BcpCommunicationPlan.UpdatedAt

	return &BcpCommunicationPlanDTO, nil
}

func (bcps *BcpCommunicationPlanService) Create(bcpId uint, currentUserID uint, BcpCommunicationPlanFormDTO *businessContinuityDTOs.BcpCommunicationPlanFormDTO) (uint, error) {
	tx := bcps.repo.BeginTransaction()

	responsibleUnitID, err := bcps.userSvc.GetUserPrimaryUnitID(uint(BcpCommunicationPlanFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	fmt.Println("12")
	BcpCommunicationPlan := &models.BcpCommunicationPlan{
		BcpPhaseID:               BcpCommunicationPlanFormDTO.BcpPhaseID,
		BcpProcessID:             BcpCommunicationPlanFormDTO.BcpProcessID,
		BusinessContinuityPlanID: bcpId,
		ResponsibleUserID:        BcpCommunicationPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        responsibleUnitID,
		Title:                    BcpCommunicationPlanFormDTO.Title,
		Frequency:                BcpCommunicationPlanFormDTO.Frequency,
		Channel:                  BcpCommunicationPlanFormDTO.Channel,
		Recipient:                BcpCommunicationPlanFormDTO.Recipient,
		CreatedBYID:              currentUserID,
	}

	id, err := bcps.repo.Create(BcpCommunicationPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingBcpCommunicationPlan, err := bcps.GetFullDetailsByID(BcpCommunicationPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch bcp communication plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = bcps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BcpCommunicationPlan{}.TableName(),
			id,
			existingBcpCommunicationPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BcpCommunicationPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (bcps *BcpCommunicationPlanService) Update(id uint, currentUserID uint, BcpCommunicationPlanFormDTO *businessContinuityDTOs.BcpCommunicationPlanFormDTO) error {
	tx := bcps.repo.BeginTransaction()

	existingBcpCommunicationPlan, err := bcps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBcpCommunicationPlan, err := bcps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingBcpCommunicationPlan.ResponsibleUserID != BcpCommunicationPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := bcps.userSvc.GetUserPrimaryUnitID(existingBcpCommunicationPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingBcpCommunicationPlan.ResponsibleUnitID = responsibleUnitID
	}

	existingBcpCommunicationPlan.BcpPhaseID = BcpCommunicationPlanFormDTO.BcpPhaseID
	existingBcpCommunicationPlan.BcpProcessID = BcpCommunicationPlanFormDTO.BcpProcessID
	existingBcpCommunicationPlan.ResponsibleUserID = BcpCommunicationPlanFormDTO.ResponsibleUserID
	existingBcpCommunicationPlan.Title = BcpCommunicationPlanFormDTO.Title
	existingBcpCommunicationPlan.Frequency = BcpCommunicationPlanFormDTO.Frequency
	existingBcpCommunicationPlan.Channel = BcpCommunicationPlanFormDTO.Channel
	existingBcpCommunicationPlan.Recipient = BcpCommunicationPlanFormDTO.Recipient
	existingBcpCommunicationPlan.UpdatedBYID = &currentUserID

	err = bcps.repo.Update(existingBcpCommunicationPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingBcpCommunicationPlan, err := bcps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", oldBcpCommunicationPlan.ResponsibleUser.FirstNameEN, oldBcpCommunicationPlan.ResponsibleUser.LastNameEN),
			"title":            oldBcpCommunicationPlan.Title,
			"bcp_phase":        oldBcpCommunicationPlan.BcpPhase.TitleEn,
			"bcp_process":      oldBcpCommunicationPlan.BcpProcess.Process.TitleEn,
			"frequency":        oldBcpCommunicationPlan.Frequency,
			"channel":          oldBcpCommunicationPlan.Channel,
			"recipient":        oldBcpCommunicationPlan.Recipient,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingBcpCommunicationPlan.ResponsibleUser.FirstNameEN, existingBcpCommunicationPlan.ResponsibleUser.LastNameEN),
			"title":            existingBcpCommunicationPlan.Title,
			"bcp_phase":        existingBcpCommunicationPlan.BcpPhase.TitleEn,
			"bcp_process":      existingBcpCommunicationPlan.BcpProcess.Process.TitleEn,
			"frequency":        existingBcpCommunicationPlan.Frequency,
			"channel":          existingBcpCommunicationPlan.Channel,
			"recipient":        existingBcpCommunicationPlan.Recipient,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bcps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BcpCommunicationPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BcpCommunicationPlan{}.TableName(),
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

func (bcps *BcpCommunicationPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := bcps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBcpCommunicationPlan []map[string]interface{}

	for _, id := range ids {
		BcpCommunicationPlan, err := bcps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBcpCommunicationPlan = append(deletedBcpCommunicationPlan, map[string]interface{}{
			"id":    id,
			"title": BcpCommunicationPlan.Title,
		})

		// Delete the business
		if err := bcps.repo.Delete(id, tx); err != nil {
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
			err := bcps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BcpCommunicationPlan{}.TableName(),
				successfulIds,
				deletedBcpCommunicationPlan,
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
