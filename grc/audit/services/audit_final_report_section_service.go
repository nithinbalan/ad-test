package services

import (
	"errors"
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
