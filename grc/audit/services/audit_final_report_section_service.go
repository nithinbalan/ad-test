package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFinalReportSectionServiceInterface interface {
	GetByID(id uint) (*models.AuditFinalReportSection, error)
	Update(auditID uint, id uint, currentUserID uint, auditFinalReportSectionFormDTO *auditDTOs.AuditFinalReportSectionFormDTO) error
	Delete(auditID uint, id uint, currentUserID uint) (uint, error)
}

type AuditFinalReportSectionService struct {
	repo            repositories.AuditFinalReportSectionRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	auditSvc        AuditServiceInterface
	auditFindingSvc AuditFindingServiceInterface
}

func NewAuditFinalReportSectionService(
	repo repositories.AuditFinalReportSectionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
) AuditFinalReportSectionServiceInterface {
	return &AuditFinalReportSectionService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		auditSvc:        auditSvc,
		auditFindingSvc: auditFindingSvc,
	}
}

func (afrss *AuditFinalReportSectionService) GetByID(id uint) (*models.AuditFinalReportSection, error) {
	auditFinalReportSection, err := afrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinalReportSection, nil
}

func (afrss *AuditFinalReportSectionService) Update(auditID uint, id uint, currentUserID uint, auditFinalReportSectionFormDTO *auditDTOs.AuditFinalReportSectionFormDTO) error {
	tx := afrss.repo.BeginTransaction()
	if _, err := afrss.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	existingAuditFinalReportSection, err := afrss.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditFinalReportSection.TitleEn = auditFinalReportSectionFormDTO.TitleEn
	existingAuditFinalReportSection.TitleAr = auditFinalReportSectionFormDTO.TitleAr
	existingAuditFinalReportSection.DescriptionEn = auditFinalReportSectionFormDTO.DescriptionEn
	existingAuditFinalReportSection.DescriptionAr = auditFinalReportSectionFormDTO.DescriptionAr
	existingAuditFinalReportSection.UpdatedBYID = &currentUserID

	err = afrss.repo.Update(existingAuditFinalReportSection, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := afrss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditFinalReportSection{}.TableName(),
			id,
			auditFinalReportSectionFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditFinalReportSection{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (afrss *AuditFinalReportSectionService) Delete(auditID uint, id uint, currentUserID uint) (uint, error) {
	// Begin a transaction
	tx := afrss.repo.BeginTransaction()
	// Check if the Audit exists
	if _, err := afrss.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}

	var failedId uint

	// Check if the section to be deleted exists
	sectionToDelete, err := afrss.GetByID(id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if sectionToDelete.AuditFinalReportSectionID == nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Deletion failed: Main sections cannot be deleted. Please delete the subsections first.")
	}

	// Delete the section
	if err := afrss.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}

	// Reorder the remaining sub sections after deletion
	// Fetch all remaining sub sections for this audit
	if sectionToDelete.AuditFinalReportSectionID != nil {
		parentID := *sectionToDelete.AuditFinalReportSectionID
		subSections, err := afrss.repo.GetAllSubSectionsByParentID(parentID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		parent, err := afrss.GetByID(parentID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		baseOrder := *parent.Order // Get the base order from the parent

		// Iterate through subSections and update their order
		for i, section := range subSections {
			// Set new order based on parent's order + index increment
			newOrder := baseOrder + float64(i+1)*0.1 // This ensures distinct ordering
			section.Order = float64Ptr(newOrder) // Use the helper function to create a pointer
			if err := afrss.repo.Update(&section, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	tx.Commit()

	// Log successful delete activity asynchronously
	go func() {
		err := afrss.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditFinalReportSection{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
