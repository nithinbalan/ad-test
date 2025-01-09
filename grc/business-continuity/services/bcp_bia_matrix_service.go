package services

import (
	"errors"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	bcpBiaMatrixRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BcpBiaMatrixServiceInterface interface {
	// GetAll(bcpId uint, params *businessContinuityDTOs.BcpBiaMatrixListQueryParamsDTO) (*businessContinuityDTOs.BcpBiaMatrixListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpBiaMatrixDetailDTO, error)
	Update(bcpId uint, currentUserID uint, bcpProcessFormDTO *businessContinuityDTOs.BcpProcessFormDTO) error
}

// Service and implementations
type BcpBiaMatrixService struct {
	repo                      bcpBiaMatrixRepositories.BcpBiaMatrixRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface
	businessFunctionSvc       BusinessFunctionServiceInterface
}

func NewBcpBiaMatrixService(
	repo bcpBiaMatrixRepositories.BcpBiaMatrixRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface,
	businessFunctionSvc BusinessFunctionServiceInterface,
) BcpBiaMatrixServiceInterface {
	return &BcpBiaMatrixService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		businessContinuityPlanSvc: businessContinuityPlanSvc,
		businessFunctionSvc:       businessFunctionSvc,
	}
}

func (bbms BcpBiaMatrixService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpBiaMatrixDetailDTO, error) {
	// Initialize your DTO
	var businessContinuityPlanDTO businessContinuityDTOs.BcpBiaMatrixDetailDTO

	businessContinuityPlan, err := bbms.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(businessContinuityPlan, &businessContinuityPlanDTO)
	if err != nil {
		return nil, err
	}

	return &businessContinuityPlanDTO, nil
}

func (bbms *BcpBiaMatrixService) Update(bcpId uint, currentUserID uint, bcpProcessFormDTO *businessContinuityDTOs.BcpProcessFormDTO) error {
	tx := bbms.repo.BeginTransaction()

	existingBusinessContinuityPlan, err := bbms.businessContinuityPlanSvc.GetByID(bcpId)
	if err != nil {
		return utility.NewResourceNotFoundError(err, bcpId)
	}

	oldBusinessContinuityPlan, err := bbms.businessContinuityPlanSvc.GetFullDetailsByID(bcpId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBusinessContinuityPlan.Rto = bcpProcessFormDTO.Rto
	existingBusinessContinuityPlan.RtoUnitId = bcpProcessFormDTO.RtoUnitId
	existingBusinessContinuityPlan.Mtpd = bcpProcessFormDTO.Mtpd
	existingBusinessContinuityPlan.MtpdUnitId = bcpProcessFormDTO.MtpdUnitId
	existingBusinessContinuityPlan.Rpo = bcpProcessFormDTO.Rpo
	existingBusinessContinuityPlan.RpoUnitId = bcpProcessFormDTO.RpoUnitId
	existingBusinessContinuityPlan.UpdatedBYID = &currentUserID

	err = bbms.businessContinuityPlanSvc.UpdateRecoveryObjectives(existingBusinessContinuityPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingBusinessFunction, err := bbms.businessFunctionSvc.GetByID(existingBusinessContinuityPlan.BusinessFunctionID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, bcpId)
	}

	existingBusinessFunction.Rto = bcpProcessFormDTO.Rto
	existingBusinessFunction.RtoUnitId = bcpProcessFormDTO.RtoUnitId
	existingBusinessFunction.Mtpd = bcpProcessFormDTO.Mtpd
	existingBusinessFunction.MtpdUnitId = bcpProcessFormDTO.MtpdUnitId
	existingBusinessFunction.Rpo = bcpProcessFormDTO.Rpo
	existingBusinessFunction.RpoUnitId = bcpProcessFormDTO.RpoUnitId
	existingBusinessFunction.UpdatedBYID = &currentUserID

	err = bbms.businessFunctionSvc.UpdateBusinessFunction(existingBusinessFunction, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// var bcpProcesses []*models.BcpProcess
	for _, process := range bcpProcessFormDTO.Processes {
		existingBcpProcess, err := bbms.repo.GetByID(process.ID)
		if err != nil {
			return utility.NewResourceNotFoundError(err, bcpId)
		}

		existingBcpProcess.Rto = process.Rto
		existingBcpProcess.RtoUnitId = process.RtoUnitId
		existingBcpProcess.Mtpd = process.Mtpd
		existingBcpProcess.MtpdUnitId = process.MtpdUnitId
		existingBcpProcess.Rpo = process.Rpo
		existingBcpProcess.RpoUnitId = process.RpoUnitId
		existingBcpProcess.UpdatedBYID = &currentUserID

		err = bbms.repo.Update(existingBcpProcess, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	tx.Commit()
	go func() {
		existingBcp, err := bbms.businessContinuityPlanSvc.GetFullDetailsByID(bcpId)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"rto":       oldBusinessContinuityPlan.Rto,
			"rto_unit":  utility.GetFieldValue(oldBusinessContinuityPlan.RtoUnit),
			"mtpd":      oldBusinessContinuityPlan.Mtpd,
			"mtpd_unit": utility.GetFieldValue(oldBusinessContinuityPlan.MtpdUnit),
			"rpo":       oldBusinessContinuityPlan.Rpo,
			"rpo_unit":  utility.GetFieldValue(oldBusinessContinuityPlan.RpoUnit),
		}

		newData := map[string]interface{}{
			"rto":       existingBcp.Rto,
			"rto_unit":  existingBcp.RtoUnit.TitleEn,
			"mtpd":      existingBcp.Mtpd,
			"mtpd_unit": existingBcp.MtpdUnit.TitleEn,
			"rpo":       existingBcp.Rpo,
			"rpo_unit":  existingBcp.RpoUnit.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bbms.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BusinessContinuityPlan{}.TableName(),
				bcpId,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BusinessContinuityPlan{}.TableName(),
					"recordID":  bcpId,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}
