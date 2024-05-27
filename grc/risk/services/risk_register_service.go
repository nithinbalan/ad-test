package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"path/filepath"
	"strconv"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskRegisterServiceInterface interface {
	GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error)
	GetByID(id uint) (*models.RiskRegister, error)
	Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, risk *riskDTOs.RiskRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type RiskRegisterService struct {
	repo               riskRepositories.RiskRegisterRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	riskCategorySvc    RiskCategoryServiceInterface
	riskSourceSvc      RiskSourceServiceInterface
	userSvc            organizationServices.UserServiceInterface
	riskStatusSvc      RiskStatusServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	// controlSvc               organizationServices.ControlServiceInterface
	// complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
}

func NewRiskRegisterService(
	repo riskRepositories.RiskRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	riskSourceSvc RiskSourceServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	// controlSvc organizationServices.ControlServiceInterface,
	// complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
) RiskRegisterServiceInterface {
	return &RiskRegisterService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		riskCategorySvc:    riskCategorySvc,
		riskSourceSvc:      riskSourceSvc,
		userSvc:            userSvc,
		riskStatusSvc:      riskStatusSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		// controlSvc:               controlSvc,
		// complianceRequirementSvc: complianceRequirementSvc,
	}
}

func (rrs *RiskRegisterService) GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	return rrs.repo.GetAll(params)
}

func (rrs RiskRegisterService) GetByID(id uint) (*models.RiskRegister, error) {
	risk, err := rrs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return risk, nil
}

func (rrs RiskRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error) {
	// Initialize your DTO
	var riskDTO riskDTOs.RiskRegisterDetailDTO

	risk, err := rrs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(risk, &riskDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskDTO, nil
}

func (rrs *RiskRegisterService) Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

	if err := rrs.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := rrs.parseRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := rrs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		ReferenceNumber:   referenceNumber,
		RiskCategoryID:    riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:      riskRegisterFormDTO.RiskSourceID,
		OwnerUnitID:       riskRegisterFormDTO.OwnerUnitID,
		OwnerID:           riskRegisterFormDTO.OwnerID,
		RiskStatusID:      riskRegisterFormDTO.RiskStatusID,
		Title:             riskRegisterFormDTO.Title,
		Description:       riskRegisterFormDTO.Description,
		ReportedAt:        dates.ReportedAt,
		ReviewFrequencyID: riskRegisterFormDTO.ReviewFrequencyID,
		NextReviewDate:    dates.NextReviewDate,
		Comments:          riskRegisterFormDTO.Comments,
		CreatedBy:         currentUserID,
	}

	id, err := rrs.repo.Create(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// // Save ComplianceRequirements
	// var complianceRequirements []*models.RiskComplianceRequirement
	// for _, req := range riskRegisterFormDTO.ComplianceRequirements {
	// 	requirement := &models.RiskComplianceRequirement{
	// 		RiskRegisterID:          id,
	// 		ComplianceRequirementID: *req,
	// 		CreatedBy:               currentUserID,
	// 	}
	// 	complianceRequirements = append(complianceRequirements, requirement)
	// }

	// if len(complianceRequirements) > 0 {
	// 	if err := rrs.repo.SaveRiskComplianceRequirements(complianceRequirements, tx); err != nil {
	// 		tx.Rollback()
	// 		return 0, err
	// 	}
	// }
	// // Save Controls
	// var controls []*models.RiskControl
	// for _, ctrl := range riskRegisterFormDTO.Controls {
	// 	control := &models.RiskControl{
	// 		RiskRegisterID: id,
	// 		ControlID:      *ctrl,
	// 		CreatedBy:      currentUserID,
	// 	}
	// 	controls = append(controls, control)
	// }

	// if len(controls) > 0 {
	// 	if err := rrs.repo.SaveRiskControls(controls, tx); err != nil {
	// 		tx.Rollback()
	// 		return 0, err
	// 	}
	// }

	var documents []*models.RiskDocument

	// Iterate over each document in the array
	for _, doc := range riskRegisterFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), encryptedFilename)

		// Create a new RiskDocument instance
		document := &models.RiskDocument{
			RiskRegisterID: id,
			URL:            url,
			Icon:           doc.Icon,
			Name:           doc.FileName,
			Size:           doc.Size,
			Token:          encryptedFilename,
			OriginalName:   doc.OriginalName,
			Type:           doc.Type,
			CreatedBy:      currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the RiskDocument
		if err := rrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			riskRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rrs *RiskRegisterService) Update(id uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	tx := rrs.repo.BeginTransaction()

	if err := rrs.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return err
	}

	existingRiskRegister, err := rrs.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := rrs.parseRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return err
	}

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.OwnerUnitID = riskRegisterFormDTO.OwnerUnitID
	existingRiskRegister.OwnerID = riskRegisterFormDTO.OwnerID
	existingRiskRegister.RiskStatusID = riskRegisterFormDTO.RiskStatusID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReportedAt = dates.ReportedAt
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.NextReviewDate = dates.ReportedAt
	existingRiskRegister.Comments = riskRegisterFormDTO.Comments
	existingRiskRegister.UpdatedBy = &currentUserID

	if err := rrs.repo.Update(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	// // Save ComplianceRequirements
	// var complianceRequirements []*models.RiskComplianceRequirement
	// for _, req := range risk.ComplianceRequirements {
	// 	requirement := &models.RiskComplianceRequirement{
	// 		RiskRegisterID:          id,
	// 		ComplianceRequirementID: *req,
	// 		CreatedBy:               currentUserID,
	// 	}
	// 	complianceRequirements = append(complianceRequirements, requirement)
	// }

	// if len(complianceRequirements) > 0 {
	// 	if err := rrs.repo.SaveRiskComplianceRequirements(complianceRequirements, tx); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	// Save Controls
	// var controls []*models.RiskControl
	// for _, ctrl := range risk.Controls {
	// 	control := &models.RiskControl{
	// 		RiskRegisterID: id,
	// 		ControlID:      *ctrl,
	// 		CreatedBy:      currentUserID,
	// 	}
	// 	controls = append(controls, control)
	// }

	// if len(controls) > 0 {
	// 	if err := rrs.repo.SaveRiskControls(controls, tx); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	// Save Documents
	// documents update
	var documents []*models.RiskDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskDocument

	// Iterate over each document in the array
	for _, doc := range riskRegisterFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("risk-register", strconv.FormatUint(uint64(id), 10), encryptedFilename)
			// Create a new RiskDocument instance
			document := &models.RiskDocument{
				RiskRegisterID: id,
				URL:            url,
				Icon:           doc.Icon,
				Name:           doc.FileName,
				Size:           doc.Size,
				Token:          encryptedFilename,
				OriginalName:   doc.OriginalName,
				Type:           doc.Type,
				CreatedBy:      currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := rrs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := rrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(existingRiskRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(existingRiskRegister.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			id,
			riskRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (rrs *RiskRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := rrs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*riskDTOs.RiskDocumentDTO)

	for _, id := range ids {
		risk, err := rrs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// // Check if controls exist
		// if len(risk.Controls) > 0 {
		// 	var deleteControlIds []*uint
		// 	for _, ctrl := range risk.Controls {
		// 		deleteControlIds = append(deleteControlIds, ctrl.ID)
		// 	}
		// 	if err := rrs.repo.DeleteRiskControls(deleteControlIds, tx); err != nil {

		// 		tx.Rollback()
		// 		return nil, err
		// 	}
		// }

		// // Delete each compliance requirement if requirements exist
		// if len(risk.ComplianceRequirements) > 0 {
		// 	var deleteRequirementIds []*uint
		// 	for _, req := range risk.ComplianceRequirements {
		// 		deleteRequirementIds = append(deleteRequirementIds, req.ID)
		// 	}
		// 	if err := rrs.repo.DeleteRiskComplianceRequirements(deleteRequirementIds, tx); err != nil {
		// 		tx.Rollback()
		// 		return nil, err
		// 	}

		// }

		// Delete each risk document if documents exist
		if len(risk.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range risk.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := rrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Delete the risk
		if err := rrs.repo.Delete(id, tx); err != nil {
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

	if len(successfulIds) > 0 {
		for riskID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
				}
			}
			destinationPath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := rrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskRegister{}.TableName(),
				successfulIds,
				nil,
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

func (rrs *RiskRegisterService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	if _, err := rrs.riskCategorySvc.GetByID(riskRegisterFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}

	if _, err := rrs.riskSourceSvc.GetByID(riskRegisterFormDTO.RiskSourceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk source does not exist"), "not_exist", "risk source")
	}

	//@todo org unit validation

	if _, err := rrs.userSvc.GetByID(riskRegisterFormDTO.OwnerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner does not exist"), "not_exist", "owner")
	}

	if _, err := rrs.riskStatusSvc.GetByID(riskRegisterFormDTO.RiskStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk status does not exist"), "not_exist", "risk status")
	}

	if riskRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := rrs.reviewFrequencySvc.GetByID(*riskRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	return nil
}

func (rrs *RiskRegisterService) parseRiskRegisterDates(dto *riskDTOs.RiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterParsedDates{}

	var err error
	if *dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(*dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported at"), "invalid_date", "reported at")
		}
	}

	if *dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}
