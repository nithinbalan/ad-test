package services

import (
	"fmt"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	complianceModels "jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	activityLogModels "jethur/grc/logs/models"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
)

type SLARegisterServiceInterface interface {
	GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, sla *slaDTOs.SLARegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type SLARegisterService struct {
	repo               repositories.SLARegisterRepositoryInterface
	slaCategorySvc     SLACategoryServiceInterface
	severityLevelSvc   masterServices.SeverityLevelServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewSLARegisterService(
	repo repositories.SLARegisterRepositoryInterface,
	slaCategorySvc SLACategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) SLARegisterServiceInterface {
	return &SLARegisterService{
		repo:               repo,
		slaCategorySvc:     slaCategorySvc,
		severityLevelSvc:   severityLevelSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (srs *SLARegisterService) GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error) {
	return srs.repo.GetAll(params)
}

func (srs SLARegisterService) GetByID(id uint) (*models.SLARegister, error) {
	return srs.repo.GetByID(id)
}

func (srs *SLARegisterService) Create(currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error) {
	// Validation and error handling code omitted for brevity

	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
		return 0, err
	}
	referenceNumber, err := srs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByStatus("active")
	if err != nil {
		return 0, err
	}
	// Create SLA register
	slaRegister := &models.SLARegister{
		SLACategoryID:     slaRegisterFormDTO.SLACategoryID,
		SeverityLevelID:   slaRegisterFormDTO.SeverityLevelID,
		ResponsibleUnitID: slaRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID: slaRegisterFormDTO.ResponsibleUserID,
		LifecycleStatusID: lifecycleStatus.ID,
		ReferenceNumber:   referenceNumber,
		TitleEn:           slaRegisterFormDTO.TitleEn,
		TitleAr:           slaRegisterFormDTO.TitleAr,
		Description:       slaRegisterFormDTO.Description,
		Client:            slaRegisterFormDTO.Client,
		CreatedBy:         currentUserID,
	}
	id, err := srs.repo.Create(slaRegister)
	if err != nil {
		return 0, err
	}
	// Create SLA documents
	var documents []*models.SLADocument
	for _, docDTO := range slaRegisterFormDTO.Documents {
		doc := &models.SLADocument{
			SLARegisterID: slaRegister.ID, // Use the ID of the created SLA register
			Version:       docDTO.Version,
			IssueDate:     dates.IssueDate,
			ExpiryDate:    dates.ExpiryDate,
			SA1Date:       dates.SA1Date,
			SA2Date:       dates.SA2Date,
			Comment:       docDTO.Comment,
			Language:      docDTO.Language,
			IsLatest:      true,
			URL:           docDTO.URL,
			Name:          docDTO.Name,
			Size:          docDTO.Size,
			Token:         docDTO.Token,
			Type:          docDTO.Type,
			Icon:          docDTO.Icon,
			CreatedBy:     currentUserID,
		}
		documents = append(documents, doc)
	}

	// Create SLA documents in a single transaction
	if err := srs.repo.AssociateDocumentsWithSLARegister(documents); err != nil {
		return 0, err
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.SLARegister{}.TableName(),
			id,
			slaRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.SLARegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (srs *SLARegisterService) Update(id uint, currentUserID uint, slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) error {

	existingSLARegister, err := srs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return err
	}

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
		return err
	}

	existingSLARegister.SLACategoryID = slaRegisterFormDTO.SLACategoryID
	existingSLARegister.ResponsibleUnitID = slaRegisterFormDTO.ResponsibleUnitID
	existingSLARegister.ResponsibleUserID = slaRegisterFormDTO.ResponsibleUserID
	existingSLARegister.SeverityLevelID = slaRegisterFormDTO.SeverityLevelID
	existingSLARegister.LifecycleStatusID = slaRegisterFormDTO.LifecycleStatusID
	existingSLARegister.TitleEn = slaRegisterFormDTO.TitleEn
	existingSLARegister.TitleAr = slaRegisterFormDTO.TitleAr
	existingSLARegister.Description = slaRegisterFormDTO.Description
	existingSLARegister.Client = slaRegisterFormDTO.Client
	existingSLARegister.UpdatedBy = currentUserID

	err = srs.repo.Update(existingSLARegister)

	if err != nil {
		return err
	}
	// Update SLA document details
	for _, docDTO := range slaRegisterFormDTO.Documents {
		// Get the existing SLA document by ID
		existingSLADocument, err := srs.repo.GetSLADocumentByID(*docDTO.ID)
		if err != nil {
			return err
		}

		// Update SLA document details
		existingSLADocument.Version = docDTO.Version
		existingSLADocument.IssueDate = dates.IssueDate
		existingSLADocument.ExpiryDate = dates.ExpiryDate
		existingSLADocument.SA1Date = dates.SA1Date
		existingSLADocument.SA2Date = dates.SA2Date
		existingSLADocument.Comment = docDTO.Comment
		existingSLADocument.Language = docDTO.Language
		existingSLADocument.URL = docDTO.URL
		existingSLADocument.Name = docDTO.Name
		existingSLADocument.Size = docDTO.Size
		existingSLADocument.Token = docDTO.Token
		existingSLADocument.Type = docDTO.Type
		existingSLADocument.Icon = docDTO.Icon
		existingSLADocument.UpdatedBy = currentUserID

		// Update SLA document in the repository by passing it as a slice
		err = srs.repo.AssociateDocumentsWithSLARegister([]*complianceModels.SLADocument{existingSLADocument})
		if err != nil {
			return err
		}
	}
	return nil
}

func (srs *SLARegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := srs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = srs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				activityLogModels.DeleteAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (srs *SLARegisterService) validateSLARegisterFormDTO(slaRegisterFormDTO *slaDTOs.SLARegisterFormDTO) error {
	if _, err := srs.slaCategorySvc.GetByID(slaRegisterFormDTO.SLACategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}

	if slaRegisterFormDTO.LifecycleStatusID != 0 {
		if _, err := srs.lifecycleStatusSvc.GetByID(slaRegisterFormDTO.LifecycleStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
		}
	}
	if slaRegisterFormDTO.SeverityLevelID != 0 {
		if _, err := srs.severityLevelSvc.GetByID(slaRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (srs *SLARegisterService) parseSLARegisterDates(dto *slaDTOs.SLARegisterFormDTO) (*slaDTOs.ParsedSLADates, error) {
	parsed := &slaDTOs.ParsedSLADates{}
	var err error
	if dto.Documents != nil {
		for _, doc := range dto.Documents {
			// Validate issue_date
			if doc.IssueDate != "" {
				parsed.IssueDate, err = utility.ParseDate(doc.IssueDate)
				if err != nil {
					return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
				}
			}

			// Validate expiry_date
			if doc.ExpiryDate != "" {
				parsed.ExpiryDate, err = utility.ParseDate(doc.ExpiryDate)
				if err != nil {
					return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
				}
			}
			// Validate SA1Date
			if doc.SA1Date != "" {
				parsed.SA1Date, err = utility.ParseDate(doc.SA1Date)
				if err != nil {
					return nil, utility.NewValidationError(fmt.Errorf("invalid sa1 date"), "invalid_date", "sa1 date")
				}
			}
			// Validate SA2Date
			if doc.SA2Date != "" {
				parsed.SA2Date, err = utility.ParseDate(doc.SA2Date)
				if err != nil {
					return nil, utility.NewValidationError(fmt.Errorf("invalid sa2 date"), "invalid_date", "sa2 date")
				}
			}
		}
	}
	return parsed, nil
}
