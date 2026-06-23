
package services

import (
	"errors"
	"fmt"

	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	awarenessRepositories "jethur/grc/awareness-reporting/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AwarenessReportServiceInterface interface {
	GetAll(params *awarenessDTOs.AwarenessReportListQueryParamsDTO) (*awarenessDTOs.AwarenessReportListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessReportDetailResponseDTO, error)
	GetByID(id uint) (*models.AwarenessReport, error)
	Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessReportFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessReportFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AwarenessReportService struct {
	repo                      awarenessRepositories.AwarenessReportRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	awarenessCampaignTypeSvc  AwarenessCampaignTypeServiceInterface
	awarenessReportSubCatSvc  AwarenessReportSubCategoryServiceInterface
	severityLevelSvc          masterServices.SeverityLevelServiceInterface
	awarenessHandlingStatusSvc AwarenessHandlingStatusServiceInterface
	userSvc                   organizationServices.UserServiceInterface
}

func NewAwarenessReportService(
	repo awarenessRepositories.AwarenessReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	awarenessCampaignTypeSvc AwarenessCampaignTypeServiceInterface,
	awarenessReportSubCatSvc AwarenessReportSubCategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	awarenessHandlingStatusSvc AwarenessHandlingStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AwarenessReportServiceInterface {
	return &AwarenessReportService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		awarenessCampaignTypeSvc:   awarenessCampaignTypeSvc,
		awarenessReportSubCatSvc:   awarenessReportSubCatSvc,
		severityLevelSvc:           severityLevelSvc,
		awarenessHandlingStatusSvc: awarenessHandlingStatusSvc,
		userSvc:                    userSvc,
	}
}

func (ars *AwarenessReportService) GetAll(params *awarenessDTOs.AwarenessReportListQueryParamsDTO) (*awarenessDTOs.AwarenessReportListResponseDTO, error) {
	return ars.repo.GetAll(params)
}

func (ars *AwarenessReportService) GetByID(id uint) (*models.AwarenessReport, error) {
	report, err := ars.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return report, nil
}

func (ars *AwarenessReportService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessReportDetailResponseDTO, error) {
	var reportDTO awarenessDTOs.AwarenessReportDetailResponseDTO

	report, err := ars.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(report, &reportDTO); err != nil {
		return nil, err
	}

	reportDTO.CreatedAt = report.CreatedAt
	reportDTO.UpdatedAt = report.UpdatedAt

	return &reportDTO, nil
}

func (ars *AwarenessReportService) Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessReportFormDTO) (uint, error) {
	tx := ars.repo.BeginTransaction()

	if err := ars.validateAwarenessReportFormDTO(formDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := ars.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	awarenessReport := &models.AwarenessReport{
		ReferenceNumber:           referenceNumber,
		AwarenessReportTypeID:     formDTO.AwarenessReportTypeID,
		SubCategoryID:             formDTO.SubCategoryID,
		ReportSourceChannel:       formDTO.ReportSourceChannel,
		SeverityID:                formDTO.SeverityID,
		ReportDate:                formDTO.ReportDate,
		ResultedFromCampaign:      formDTO.ResultedFromCampaign,
		AwarenessHandlingStatusID: formDTO.AwarenessHandlingStatusID,
		HandlingOwnerID:           formDTO.HandlingOwnerID,
		ClosureDate:               formDTO.ClosureDate,
		Result:                    formDTO.Result,
		CreatedByID:               currentUserID,
	}

	id, err := ars.repo.Create(awarenessReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingReport, err := ars.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch awareness report", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ars.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessReport{}.TableName(),
			id,
			existingReport,
			currentUserID,
			"awareness", "awareness_report", "report", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AwarenessReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (ars *AwarenessReportService) Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessReportFormDTO) error {
	tx := ars.repo.BeginTransaction()

	existingReport, err := ars.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldReport, err := ars.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ars.validateAwarenessReportFormDTO(formDTO); err != nil {
		return err
	}

	existingReport.AwarenessReportTypeID = formDTO.AwarenessReportTypeID
	existingReport.SubCategoryID = formDTO.SubCategoryID
	existingReport.ReportSourceChannel = formDTO.ReportSourceChannel
	existingReport.SeverityID = formDTO.SeverityID
	existingReport.ReportDate = formDTO.ReportDate
	existingReport.ResultedFromCampaign = formDTO.ResultedFromCampaign
	existingReport.AwarenessHandlingStatusID = formDTO.AwarenessHandlingStatusID
	existingReport.HandlingOwnerID = formDTO.HandlingOwnerID
	existingReport.ClosureDate = formDTO.ClosureDate
	existingReport.Result = formDTO.Result
	existingReport.UpdatedByID = &currentUserID

	if err := ars.repo.Update(existingReport, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingFullReport, err := ars.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"awareness_report_type":      oldReport.AwarenessReportType.TitleEn,
			"report_source_channel":      utility.GetStringValue(oldReport.ReportSourceChannel),
			"resulted_from_campaign":     oldReport.ResultedFromCampaign,
			"report_date":                oldReport.ReportDate,
			"closure_date":               oldReport.ClosureDate,
			"result":                     utility.GetStringValue(oldReport.Result),
		}

		newData := map[string]interface{}{
			"awareness_report_type":      existingFullReport.AwarenessReportType.TitleEn,
			"report_source_channel":      utility.GetStringValue(existingFullReport.ReportSourceChannel),
			"resulted_from_campaign":     existingFullReport.ResultedFromCampaign,
			"report_date":                existingFullReport.ReportDate,
			"closure_date":               existingFullReport.ClosureDate,
			"result":                     utility.GetStringValue(existingFullReport.Result),
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := ars.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AwarenessReport{}.TableName(),
				id,
				changedData,
				currentUserID,
				"awareness", "awareness_report", "report", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AwarenessReport{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (ars *AwarenessReportService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ars.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedReports []map[string]interface{}

	for _, id := range ids {
		report, err := ars.repo.GetByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedReports = append(deletedReports, map[string]interface{}{
			"id":               id,
			"reference_number": report.ReferenceNumber,
		})

		if err := ars.repo.Delete(id, tx); err != nil {
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

	if len(successfulIds) > 0 {
		go func() {
			err := ars.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessReport{}.TableName(),
				successfulIds,
				deletedReports,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (ars *AwarenessReportService) validateAwarenessReportFormDTO(formDTO *awarenessDTOs.AwarenessReportFormDTO) error {
	if _, err := ars.awarenessCampaignTypeSvc.GetByID(formDTO.AwarenessReportTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("awareness report type does not exist"), "not_exist", "awareness report type")
	}

	if formDTO.SubCategoryID != nil {
		if _, err := ars.awarenessReportSubCatSvc.GetByID(*formDTO.SubCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("sub category does not exist"), "not_exist", "sub category")
		}
	}

	if formDTO.SeverityID != nil {
		if _, err := ars.severityLevelSvc.GetByID(*formDTO.SeverityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if formDTO.AwarenessHandlingStatusID != nil {
		if _, err := ars.awarenessHandlingStatusSvc.GetByID(*formDTO.AwarenessHandlingStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("awareness handling status does not exist"), "not_exist", "awareness handling status")
		}
	}

	if formDTO.HandlingOwnerID != nil {
		if _, err := ars.userSvc.GetByID(*formDTO.HandlingOwnerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("handling owner does not exist"), "not_exist", "handling owner")
		}
	}

	return nil
}