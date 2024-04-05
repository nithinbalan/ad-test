package services

import (
	"fmt"
	slaDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/utility"
)

type SLARenewServiceInterface interface {
	GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error)
	Renew(docID uint, currentUserID uint, document *slaDTOs.SLARenewFormDTO) error
	Revert(docID uint, currentUserID uint) error
}

// Service and implementations
type SLARenewService struct {
	repo           repositories.SLARenewRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewSLARenewService(
	repo repositories.SLARenewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) SLARenewServiceInterface {
	return &SLARenewService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (srs *SLARenewService) GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error) {
	return srs.repo.GetAll(id, params)
}

func (srs *SLARenewService) Renew(docID uint, currentUserID uint, slaRenewFormDTO *slaDTOs.SLARenewFormDTO) error {

	// Check if Document is provided
	if slaRenewFormDTO.Document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}

	if err := srs.validateSLARenewDocumentFormDTO(slaRenewFormDTO.Document); err != nil {
		return err
	}

	dates, err := srs.parseSLARenewDates(slaRenewFormDTO)
	if err != nil {
		return err
	}
	// Check if the file is updated
	existingSLADocument, err := srs.repo.GetSLADocumentByID(docID)
	if err != nil {
		return err
	}
	if err := srs.repo.StoreSLADocumentFile(existingSLADocument.SLARegisterID, slaRenewFormDTO.Document.FileName); err != nil {
		return err
	}
	//renew document
	slaDocument := &complianceModels.SLADocument{
		SLARegisterID: existingSLADocument.SLARegisterID,
		Version:       slaRenewFormDTO.Version,
		IssueDate:     dates.IssueDate,
		ExpiryDate:    dates.ExpiryDate,
		SA1Date:       dates.SA1Date,
		SA2Date:       dates.SA2Date,
		Language:      slaRenewFormDTO.Document.Language,
		URL:           slaRenewFormDTO.Document.URL, //URL To be decided
		Icon:          slaRenewFormDTO.Document.Icon,
		Name:          slaRenewFormDTO.Document.FileName,
		Size:          slaRenewFormDTO.Document.Size,
		Token:         slaRenewFormDTO.Document.FileName, // Assuming token is same as fileName
		Type:          slaRenewFormDTO.Document.Type,
		IsLatest:      true,
		CreatedBy:     currentUserID,
	}

	// Save SLA document
	if err := srs.repo.SaveSLADocument(slaDocument); err != nil {
		return err
	}
	// Update is_latest=false for all other documents of the same SLA register
	if err := srs.repo.UpdateOtherSLADocumentsIsLatest(slaDocument.ID, slaDocument.SLARegisterID); err != nil {
		return err
	}

	go func() {
		err := srs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			complianceModels.SLADocument{}.TableName(),
			docID,
			slaRenewFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.UpdateAction,
				"tableName": complianceModels.SLADocument{}.TableName(),
				"recordID":  docID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
func (srs *SLARenewService) Revert(docID uint, currentUserID uint) error {
	// Check if the file is updated
	existingSLADocument, err := srs.repo.GetSLADocumentByID(docID)
	if err != nil {
		return err
	}
	previousSLADocument, err := srs.repo.FindPreviousDocument(docID, existingSLADocument.SLARegisterID)
	if err != nil {
		return err
	}
	previousSLADocument.IsLatest = true
	previousSLADocument.UpdatedBy = currentUserID
	// Save updated SLA document
	if err := srs.repo.SaveSLADocument(previousSLADocument); err != nil {
		return err
	}
	// Update is_latest=false for all other documents of the same SLA register
	if err := srs.repo.UpdateOtherSLADocumentsIsLatest(previousSLADocument.ID, previousSLADocument.SLARegisterID); err != nil {
		return err
	}
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

	// Validate language
	if document.Language == "" {
		return utility.NewValidationError(fmt.Errorf("language is required"), "required", "language")
	}

	// Validate URL
	if document.URL == "" {
		return utility.NewValidationError(fmt.Errorf("URL is required"), "required", "URL")
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
