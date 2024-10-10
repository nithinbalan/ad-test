package services

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"time"

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

	// Fetch the full details of the audit plan
	auditPlanDetails, err := aps.auditPlanSvc.GetFullDetailsByID(auditPlanId)
	if err != nil {
		return err
	}

	// Load the user details for FromUserID (team lead or assigned user)
	fromUser, err := aps.userSvc.GetFullDetailsByID(*auditPlanDetails.TeamLeadID)
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
		// Generate the subject and body in both English and Arabic
		subjectEn := fmt.Sprintf("Audit Commencement Letter - %s", auditPlanDetails.Title)
		subjectAr := fmt.Sprintf("خطاب بدء التدقيق - %s", auditPlanDetails.Title)
		bodyEn := generateLetterBodyEn(auditableItems, auditPlanDetails, fromUser)
		bodyAr := generateLetterBodyAr(auditableItems, auditPlanDetails, fromUser)

		// Create the commencement letter entry
		commencementLetter := &models.AuditPlanCommencementLetter{
			AuditPlanID:       auditPlanId,
			ResponsibleUnitID: unitID,
			SubjectEn:         subjectEn,
			SubjectAr:         subjectAr,
			BodyEn:            bodyEn,
			BodyAr:            bodyAr,
			FromUserID:        *auditPlanDetails.TeamLeadID,
			CreatedByID:       currentUserID,
			CreatedAt:         time.Now(),
		}

		if _, err := aps.repo.Create(commencementLetter, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func generateLetterBodyEn(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO, fromUser *organizationDTOs.UserDetailsDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleEn != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleEn) // Use English title if available
		}

		// Get the responsible unit's name (TitleEn)
		if item.AuditableItem.ResponsibleUnit.TitleEn != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleEn
		}
	}

	// Team lead information
	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameEN, auditPlan.TeamLead.LastNameEN)

	// From User details for the signature
	fromUserName := fmt.Sprintf("%s %s", fromUser.FirstNameEn, fromUser.LastNameEn)
	fromUserDesignation := fromUser.Designation

	// Generate the body of the letter
	body := fmt.Sprintf(`
Dear %s,

This is to inform you that JethurGRC will commence an audit of your department for the period %s to %s.

The audit scope will include:
%s

Please prepare relevant documents for review. Our team, led by %s, will contact you to schedule an initial meeting.

Your cooperation is appreciated. For any questions, please contact %s.

Thank you,
%s
%s
`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName, fromUserName, fromUserDesignation.TitleEn)

	return body
}

func generateLetterBodyAr(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO, fromUser *organizationDTOs.UserDetailsDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleAr != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleAr) // Use Arabic title if available
		}

		// Get the responsible unit's name (TitleAr)
		if item.AuditableItem.ResponsibleUnit.TitleAr != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleAr
		}
	}

	// Team lead information
	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameAR, auditPlan.TeamLead.LastNameAR)

	// From User details for the signature
	fromUserName := fmt.Sprintf("%s %s", fromUser.FirstNameAr, fromUser.LastNameAr)
	fromUserDesignation := fromUser.Designation

	// Generate the body of the letter in Arabic
	body := fmt.Sprintf(`
عزيزي %s،

نود إبلاغكم أن JethurGRC سيبدأ عملية التدقيق على إدارتكم خلال الفترة من %s إلى %s.

سيشمل نطاق التدقيق:
%s

يرجى إعداد المستندات ذات الصلة للمراجعة. سيقوم فريقنا، بقيادة %s، بالتواصل معكم لتحديد موعد الاجتماع الأولي.

نقدر تعاونكم. لأي استفسارات، يرجى التواصل مع %s.

شكرًا لكم،
%s
%s
`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName, fromUserName, fromUserDesignation.TitleAr)

	return body
}
