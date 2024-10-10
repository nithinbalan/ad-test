package services

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
)

type AuditPlanCommencementLetterServiceInterface interface {
	Create(auditPlanId uint, currentUserID uint) error
}

type AuditPlanCommencementLetterService struct {
	repo           repositories.AuditPlanCommencementLetterRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditPlanSvc   AuditPlanServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewAuditPlanCommencementLetterService(
	repo repositories.AuditPlanCommencementLetterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AuditPlanCommencementLetterServiceInterface {
	return &AuditPlanCommencementLetterService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditPlanSvc:   auditPlanSvc,
		userSvc:        userSvc,
	}
}

func (aps *AuditPlanCommencementLetterService) Create(auditPlanId uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()

	auditPlanDetails, err := aps.auditPlanSvc.GetFullDetailsByID(auditPlanId)
	if err != nil {
		return err
	}

	// Load the user details for FromUserID
	fromUser, err := aps.userSvc.GetFullDetailsByID(*auditPlanDetails.TeamLeadID) // Assuming fromUserID is the team lead or passed correctly
	if err != nil {
		tx.Rollback()
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
		body := generateLetterBody(auditableItems, auditPlanDetails, fromUser)

		commencementLetter := &models.AuditPlanCommencementLetter{
			AuditPlanID:       auditPlanId,
			ResponsibleUnitID: unitID,
			Subject:           subject,
			Body:              body,
			FromUserID:        *auditPlanDetails.TeamLeadID,
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

func generateLetterBody(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO, fromUser *organizationDTOs.UserDetailsDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleEn != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleEn) // Use English title if available
		} else if item.AuditableItem.TitleAr != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleAr) // Fall back to Arabic title if necessary
		}

		// Get the responsible unit's name (TitleEn or TitleAr)
		if item.AuditableItem.ResponsibleUnit.TitleEn != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleEn
		} else {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleAr
		}
	}

	// Team lead information
	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameEN, auditPlan.TeamLead.LastNameEN)
	teamLeadEmail := auditPlan.TeamLead.Email
	teamLeadPhone := auditPlan.TeamLead.MobileNumber

	// Conditionally format contact information based on the presence of a phone number
	contactInfo := fmt.Sprintf("%s at %s", teamLeadName, teamLeadEmail)
	if teamLeadPhone != nil {
		contactInfo += fmt.Sprintf(" or %s", *teamLeadPhone)
	}

	// From User details for the signature
	fromUserName := fmt.Sprintf("%s %s", fromUser.FirstNameEn, fromUser.LastNameEn)
	fromUserDesignation := fromUser.Designation // Assuming designation is preloaded and available

	// Generate the body of the letter
	body := fmt.Sprintf(`
Dear %s,

This is to inform you that JethurGRC will commence an audit of your organization for the period %s to %s. The audit will assess compliance, internal controls, and financial accuracy.

The audit scope will include:
%s

Please prepare relevant documents for review. Our team, led by %s, will contact you to schedule an initial meeting.

Your cooperation is appreciated. For any questions, please contact %s.

Thank you,
%s
%s
`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, contactInfo, fromUserName, fromUserDesignation.TitleEn)

	return body
}
