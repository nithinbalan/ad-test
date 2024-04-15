package services

import (
	"fmt"
	"jethur/config"
	slaDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
	"log"
	"strconv"
)

type SLARenewServiceInterface interface {
	GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error)
	Renew(slaId uint, currentUserID uint, document *slaDTOs.SLARenewFormDTO) error
	Revert(slaId uint, currentUserID uint) error
}

// Service and implementations
type SLARenewService struct {
	repo               repositories.SLARenewRepositoryInterface
	slaRegRepo         repositories.SLARegisterRepositoryInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
}

func NewSLARenewService(
	repo repositories.SLARenewRepositoryInterface,
	slaRegRepo repositories.SLARegisterRepositoryInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) SLARenewServiceInterface {
	return &SLARenewService{
		repo:               repo,
		slaRegRepo:         slaRegRepo,
		lifecycleStatusSvc: lifecycleStatusSvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (srs *SLARenewService) GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error) {
	return srs.repo.GetAll(id, params)
}

func (srs *SLARenewService) Renew(slaId uint, currentUserID uint, slaRenewFormDTO *slaDTOs.SLARenewFormDTO) error {
	tx := srs.repo.BeginTransaction()

	//@review
	//we can move this checking in to validateSLARenewDocumentFormDTO function.

	if err := srs.validateSLARenewDocumentFormDTO(slaRenewFormDTO.Document); err != nil {
		return err
	}

	dates, err := srs.parseSLARenewDates(slaRenewFormDTO)
	if err != nil {
		return err
	}
	slaRegister, err := srs.slaRegRepo.GetByID(slaId)
	if err != nil {
		return err
	}
	// Check if the SLA status is not archived
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}

	if slaRegister.LifecycleStatusID == lifecycleStatus.ID {
		return utility.NewForbiddenError("Cannot renew SLA with archived status")
	}
	existingSLADocument, err := srs.slaRegRepo.GetLatestSLADocumentBySLARegisterID(slaId)
	if err != nil {
		return err
	}

	if err := srs.slaRegRepo.StoreSLADocumentFile(slaId, slaRenewFormDTO.Document.FileName); err != nil {
		tx.Rollback()
		return err
	}

	//encrypt FileName for token
	dataToEncrypt := []byte(slaRenewFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/sla-registers/%d/%s", config.GetEnv("API_URL", "http://localhost:8080"), existingSLADocument.SLARegisterID, encryptedFilename)

	// Get the current version of the SLA document
	currentVersionStr := existingSLADocument.Version

	// Parse the current version as a float
	currentVersionFloat, err := strconv.ParseFloat(currentVersionStr, 64)
	if err != nil {
		log.Printf("Failed to parse current version string as float: %v", err)
		return err
	}

	// Convert the float version to an integer (truncate the decimal part)
	currentVersion := int(currentVersionFloat)
	// Increment the version by one
	newVersion := currentVersion + 1.0
	// Convert the new version back to a string
	newVersionStr := strconv.Itoa(newVersion)

	//renew document
	slaDocument := &complianceModels.SLADocument{
		SLARegisterID: existingSLADocument.SLARegisterID,
		Version:       newVersionStr, //@review, Need to increment version
		Language:      "en",
		IssueDate:     dates.IssueDate,
		ExpiryDate:    dates.ExpiryDate,
		SA1Date:       dates.SA1Date,
		SA2Date:       dates.SA2Date,
		URL:           url,
		Comment:       slaRenewFormDTO.Comment,
		Icon:          slaRenewFormDTO.Document.Icon,
		Name:          slaRenewFormDTO.Document.FileName,
		Size:          slaRenewFormDTO.Document.Size,
		Token:         encryptedFilename,
		Type:          slaRenewFormDTO.Document.Type,
		IsLatest:      true,
		CreatedBy:     currentUserID,
		// @add comment
	}

	// Save SLA document
	if err := srs.slaRegRepo.SaveSLADocument(slaDocument, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := srs.repo.SetOlderSLADocsToNotLatest(slaDocument.ID, slaDocument.SLARegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			complianceModels.SLADocument{}.TableName(),
			slaId,
			slaRenewFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": complianceModels.SLADocument{}.TableName(),
				"recordID":  slaId,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
func (srs *SLARenewService) Revert(slaId uint, currentUserID uint) error {
	tx := srs.repo.BeginTransaction()

	//@review
	//get sla register and validate status. (archiave check)
	//get latest document
	//delete file from storage
	//delete latest dcoument
	//updte previous one as is latest true

	slaRegister, err := srs.slaRegRepo.GetByID(slaId)
	if err != nil {
		return err
	}
	// Check if the SLA status is not archived
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}
	if slaRegister.LifecycleStatusID == lifecycleStatus.ID {
		return utility.NewForbiddenError("Cannot revert SLA with archived status")
	}

	// Check if the file is updated
	existingSLADocument, err := srs.slaRegRepo.GetLatestSLADocumentBySLARegisterID(slaId)
	if err != nil {
		return err
	}
	if err := srs.slaRegRepo.RemoveSLADocumentFile(slaId, existingSLADocument.Name); err != nil {
		tx.Rollback()
		return err
	}

	previousSLADocument, err := srs.repo.FindPreviousDocument(existingSLADocument.ID, existingSLADocument.SLARegisterID)
	if err != nil {
		return err
	}
	// Delete existing SLA document
	if err := srs.repo.DeleteSLADocumentByID(existingSLADocument.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	previousSLADocument.IsLatest = true
	previousSLADocument.UpdatedBy = currentUserID
	// Update SLA document with IsLatest true
	if err := srs.slaRegRepo.SaveSLADocument(previousSLADocument, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			complianceModels.SLADocument{}.TableName(),
			previousSLADocument.ID,
			previousSLADocument,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": complianceModels.SLADocument{}.TableName(),
				"recordID":  previousSLADocument.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
func (srs *SLARenewService) parseSLARenewDates(slaRenewForm *slaDTOs.SLARenewFormDTO) (*slaDTOs.ParsedSLARenewDates, error) {
	parsed := &slaDTOs.ParsedSLARenewDates{}
	var err error

	// Validate issue_date
	if slaRenewForm.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(slaRenewForm.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	// Validate expiry_date
	if slaRenewForm.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(slaRenewForm.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	// Validate SA1Date
	if slaRenewForm.SA1Date != "" {
		parsed.SA1Date, err = utility.ParseDate(slaRenewForm.SA1Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa1 date"), "invalid_date", "sa1 date")
		}
	}

	// Validate SA2Date
	if slaRenewForm.SA2Date != "" {
		parsed.SA2Date, err = utility.ParseDate(slaRenewForm.SA2Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa2 date"), "invalid_date", "sa2 date")
		}
	}

	return parsed, nil
}

func (srs *SLARenewService) validateSLARenewDocumentFormDTO(document *slaDTOs.SLARenewDocumentFormDTO) error {
	// Validate document
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
	// Validate icon
	if document.Icon == "" {
		return utility.NewValidationError(fmt.Errorf("icon is required"), "required", "icon")
	}

	// Validate file type (PDF)
	if document.Type != ".pdf" {
		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
	}

	// @todo
	// Validate file size (if needed)
	// For example:
	// if document.Size > MaxFileSize {
	//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
	// }

	return nil
}
