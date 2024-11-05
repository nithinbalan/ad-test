package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThreatServiceInterface interface {
	GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error)
	GetByID(id uint) (*models.Threat, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ThreatDetailDTO, error)
	Create(threatForm *dtos.ThreatFormDTO, currentUserID uint) (uint, error)
	Update(id uint, threatForm *dtos.ThreatFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ThreatService struct {
	repo             repositories.ThreatRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	threatTypeSvc    ThreatTypeServiceInterface
	severityLevelSvc masterServices.SeverityLevelServiceInterface
	threatStatusSvc  ThreatStatusServiceInterface
	userSvc          organizationServices.UserServiceInterface
}

func NewThreatService(
	repo repositories.ThreatRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatTypeSvc ThreatTypeServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	threatStatusSvc ThreatStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ThreatServiceInterface {
	return &ThreatService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		threatTypeSvc:    threatTypeSvc,
		severityLevelSvc: severityLevelSvc,
		threatStatusSvc:  threatStatusSvc,
		userSvc:          userSvc,
	}
}

func (ts *ThreatService) GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error) {
	return ts.repo.GetAll(params)
}

func (ts *ThreatService) GetByID(id uint) (*models.Threat, error) {
	threat, err := ts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return threat, nil
}

func (as *ThreatService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ThreatDetailDTO, error) {
	var threatDetailDTO dtos.ThreatDetailDTO

	threat, err := as.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(threat, &threatDetailDTO)
	if err != nil {
		return nil, err
	}

	if threat.DetectionDate != nil {
		detectionDateString := threat.DetectionDate.Format("2006-01-02")
		threatDetailDTO.DetectionDateString = &detectionDateString
	}
	if threat.ResolutionDate != nil {
		resolutionDateString := threat.ResolutionDate.Format("2006-01-02")
		threatDetailDTO.ResolutionDateString = &resolutionDateString
	}

	// Add additional fields as necessary for the details screen
	// For example, setting relationship data or additional fields as required by the DTO

	return &threatDetailDTO, nil
}

func (ts *ThreatService) Create(threatFormDTO *dtos.ThreatFormDTO, currentUserID uint) (uint, error) {
	// Validate the DTO
	if err := ts.validateThreatFormDTO(threatFormDTO); err != nil {
		return 0, err
	}

	// Parse dates from the DTO
	dates, err := ts.parseThreatDates(threatFormDTO)
	if err != nil {
		return 0, err
	}

	// Retrieve the primary organization unit for the responsible user if not provided
	organizationUnitID := threatFormDTO.ResponsibleUnitID
	if organizationUnitID == 0 {
		organizationUnitID, err = ts.userSvc.GetUserPrimaryUnitID(threatFormDTO.ResponsibleUserID)
		if err != nil {
			return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
		}
	}

	// Populate the Threat model with data from the DTO
	threat := &models.Threat{
		TitleEn:           threatFormDTO.TitleEn,
		TitleAr:           threatFormDTO.TitleAr,
		ThreatTypeID:      threatFormDTO.ThreatTypeID,
		SeverityLevelID:   threatFormDTO.SeverityLevelID,
		ThreatStatusID:    threatFormDTO.ThreatStatusID,
		ResponsibleUserID: threatFormDTO.ResponsibleUserID,
		ResponsibleUnitID: organizationUnitID,
		DescriptionEn:     threatFormDTO.DescriptionEn,
		DescriptionAr:     threatFormDTO.DescriptionAr,
		DetectionDate:     dates.DetectionDate,
		ResolutionDate:    dates.ResolutionDate,
		CreatedBYID:       currentUserID,
		CreatedAt:         time.Now(),
	}

	// Save the threat record in the repository
	threatID, err := ts.repo.Create(threat)
	if err != nil {
		return 0, err
	}

	// Log the activity asynchronously
	go func() {
		err := ts.activityLogSvc.LogActivity(
			models.CreateAction,
			threat.TableName(),
			threatID,
			threatFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": threat.TableName(),
				"recordID":  threatID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log threat activity", fields)
		}
	}()

	return threatID, nil
}

func (ts *ThreatService) Update(id uint, threatFormDTO *dtos.ThreatFormDTO, currentUserID uint) error {
	// Retrieve the existing threat record
	existingThreat, err := ts.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate the incoming DTO
	if err := ts.validateThreatFormDTO(threatFormDTO); err != nil {
		return err
	}

	// Parse dates for DetectionDate and ResolutionDate from DTO
	dates, err := ts.parseThreatDates(threatFormDTO)
	if err != nil {
		return err
	}

	// Check if ResponsibleUserID has changed and update ResponsibleUnitID if necessary
	if existingThreat.ResponsibleUserID != threatFormDTO.ResponsibleUserID {
		organizationUnitID, err := ts.userSvc.GetUserPrimaryUnitID(threatFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
		}
		existingThreat.ResponsibleUnitID = organizationUnitID
	}

	// Update the fields with new values from the DTO
	existingThreat.TitleEn = threatFormDTO.TitleEn
	existingThreat.TitleAr = threatFormDTO.TitleAr
	existingThreat.ThreatTypeID = threatFormDTO.ThreatTypeID
	existingThreat.SeverityLevelID = threatFormDTO.SeverityLevelID
	existingThreat.ThreatStatusID = threatFormDTO.ThreatStatusID
	existingThreat.ResponsibleUserID = threatFormDTO.ResponsibleUserID
	existingThreat.DescriptionEn = threatFormDTO.DescriptionEn
	existingThreat.DescriptionAr = threatFormDTO.DescriptionAr
	existingThreat.DetectionDate = dates.DetectionDate
	existingThreat.ResolutionDate = dates.ResolutionDate
	existingThreat.UpdatedBYID = &currentUserID

	// Persist the updated threat record
	return ts.repo.Update(existingThreat)
}

func (as *ThreatService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		err := as.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Threat{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ts *ThreatService) validateThreatFormDTO(threatFormDTO *dtos.ThreatFormDTO) error {
	// Validate Threat Type
	if _, err := ts.threatTypeSvc.GetByID(threatFormDTO.ThreatTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("threat type does not exist"), "not_exist", "threat type")
	}

	// Validate Severity Level
	if _, err := ts.severityLevelSvc.GetByID(threatFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	// Validate Threat Status
	if _, err := ts.threatStatusSvc.GetByID(threatFormDTO.ThreatStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("threat status does not exist"), "not_exist", "threat status")
	}

	// Validate Responsible User
	if _, err := ts.userSvc.GetByID(threatFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	return nil
}

func (ts *ThreatService) parseThreatDates(threatRegisterFormDTO *dtos.ThreatFormDTO) (*dtos.ThreatParsedDates, error) {
	parse := &dtos.ThreatParsedDates{}

	var err error
	if threatRegisterFormDTO.DetectionDate != nil {
		parse.DetectionDate, err = utility.ParseDate(*threatRegisterFormDTO.DetectionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid detection date"), "invalid_date", "detection date")
		}
	}
	if threatRegisterFormDTO.ResolutionDate != nil {
		parse.ResolutionDate, err = utility.ParseDate(*threatRegisterFormDTO.ResolutionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid resolution date"), "invalid_date", "resolution date")
		}
	}
	return parse, nil
}
