package services

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"

	"jethur/grc/models"
)

type AuditPlanCommencementLetterServiceInterface interface {
	Create(auditPlanId uint, currentUserID uint) error
}

type AuditPlanCommencementLetterService struct {
	repo           repositories.AuditPlanCommencementLetterRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditPlanSvc   AuditPlanServiceInterface
}

func NewAuditPlanCommencementLetterService(
	repo repositories.AuditPlanCommencementLetterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
) AuditPlanCommencementLetterServiceInterface {
	return &AuditPlanCommencementLetterService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditPlanSvc:   auditPlanSvc,
	}
}

func (aps *AuditPlanCommencementLetterService) Create(auditPlanId uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()

	auditPlanDetails, err := aps.auditPlanSvc.GetFullDetailsByID(auditPlanId)
	if err != nil {
		return err
	}
	// Group Auditable Items by ResponsibleUnitID
	unitToItemsMap := make(map[uint][]*auditDTOs.AuditPlanAuditableItemDTO)
	for _, item := range auditPlanDetails.AuditableItems {
		if item.AuditableItem.ResponsibleUnitID != nil {
			unitToItemsMap[*item.AuditableItem.ResponsibleUnitID] = append(unitToItemsMap[*item.AuditableItem.ResponsibleUnitID], item)
		}
	}

	for unitID, auditableItems := range unitToItemsMap {

		subject := fmt.Sprintf("Audit Commencement Letter - %s", auditPlanDetails.Title)
		body := generateLetterBody(auditableItems, auditPlanDetails)

		commencementLetter := &models.AuditPlanCommencementLetter{
			AuditPlanID:       auditPlanId,
			ResponsibleUnitID: unitID,
			Subject:           subject,
			Body:              body,
			FromUserID:        currentUserID,
			CreatedByID:       currentUserID,
		}

		if _, err := aps.repo.Create(commencementLetter, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func generateLetterBody(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO) string {
	auditScope := ""
	for _, item := range auditableItems {
		if item.AuditableItem.TitleEn != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleEn) // Use English title if available
		} else if item.AuditableItem.TitleAr != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleAr) // Fall back to Arabic title if necessary
		}
	}

	body := fmt.Sprintf(`
Dear %s,

This is to inform you that JethurGRC will commence an audit of your organization for the period %s to %s. The audit will assess compliance, internal controls, and financial accuracy.

The audit scope will include:
%s

Please prepare relevant documents for review. Our team will contact you to schedule an initial meeting.

Your cooperation is appreciated. For any questions, please contact us.

Thank you,
%s
`, "Responsible Unit Name", auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, "Audit Department Name")

	return body
}
