package services

import (
	"fmt"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
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

	if err := srs.validateSLARenewDocumentFormDTO(&slaRenewFormDTO.Document); err != nil {
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

	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}

	if slaRegister.LifecycleStatusID == lifecycleStatus.ID {
		return utility.NewForbiddenError("Cannot renew SLA with archived status")
	}
	existingSLADocument, err := srs.repo.GetLatestSLADocumentBySLARegisterID(slaId)
	if err != nil {
		return err
	}
	//Encrypt FileName
	dataToEncrypt := []byte(slaRenewFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return err
	}
	url := utility.GenerateDocumentURL("sla-register", strconv.FormatUint(uint64(existingSLADocument.SLARegisterID), 10), encryptedFilename)

	currentVersionStr := existingSLADocument.Version
	currentVersionFloat, err := strconv.ParseFloat(currentVersionStr, 64)
	if err != nil {
		return err
	}

	currentVersion := int(currentVersionFloat)
	newVersion := currentVersion + 1
	newVersionStr := strconv.Itoa(newVersion)

	//renew document
	slaDocument := &models.SLADocument{
		SLARegisterID: existingSLADocument.SLARegisterID,
		Version:       newVersionStr, //@review, Need to increment version
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
		UpdatedBy:     currentUserID,
	}

	if err := srs.slaRegRepo.CreateDocument(slaDocument, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := srs.repo.SetOlderSLADocsToNotLatest(slaDocument.ID, slaDocument.SLARegisterID, tx); err != nil {
		tx.Rollback()
		return err
	}
	// Store SLA Document File using utility.StoreFile
	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaId), 10), "agreements")

	if err := utility.StoreFile(slaDocument.Name, storagePath); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		err := srs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SLADocument{}.TableName(),
			slaId,
			slaRenewFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SLADocument{}.TableName(),
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

	slaRegister, err := srs.slaRegRepo.GetByID(slaId)
	if err != nil {
		return err
	}
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}
	if slaRegister.LifecycleStatusID == lifecycleStatus.ID {
		return utility.NewForbiddenError("Cannot revert SLA with archived status")
	}

	existingSLADocument, err := srs.repo.GetLatestSLADocumentBySLARegisterID(slaId)
	if err != nil {
		return err
	}

	previousSLADocument, err := srs.repo.FindPreviousDocument(existingSLADocument.ID, existingSLADocument.SLARegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := srs.repo.DeleteDocumentByID(existingSLADocument.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	previousSLADocument.IsLatest = true
	previousSLADocument.UpdatedBy = currentUserID
	if err := srs.slaRegRepo.CreateDocument(previousSLADocument, tx); err != nil {
		tx.Rollback()
		return err
	}
	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaId), 10), "agreements")
	if err := utility.RemoveFile(existingSLADocument.Name, storagePath); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := srs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SLADocument{}.TableName(),
			previousSLADocument.ID,
			previousSLADocument,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SLADocument{}.TableName(),
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
