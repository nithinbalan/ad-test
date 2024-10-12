package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"

	"jethur/grc/models"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanCommencementLetterServiceInterface interface {
	GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error)
	GetByID(auditPlanID uint, letterID uint) (*models.AuditPlanCommencementLetter, error)
	GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanCommencementLetterDetailDTO, error)
	Create(auditPlanId uint, currentUserID uint) error
	Update(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterFormDTO *auditDTOs.AuditPlanCommencementLetterFormDTO) error
	Delete(auditPlanID uint, letterID uint, currentUserID uint) (uint, error)
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

func (apcls *AuditPlanCommencementLetterService) GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error) {
	return apcls.repo.GetAll(auditPlanID, keyword)
}

func (apcls *AuditPlanCommencementLetterService) GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanCommencementLetterDetailDTO, error) {
	var commencementLetterDTO auditDTOs.AuditPlanCommencementLetterDetailDTO

	// Fetch full details of the commencement letter by ID
	commencementLetter, err := apcls.repo.GetFullDetailsByID(auditPlanID, letterID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, letterID)
		}
		return nil, err
	}

	// Decode the model into DTO
	err = mapstructure.Decode(commencementLetter, &commencementLetterDTO)
	if err != nil {
		return nil, err
	}

	return &commencementLetterDTO, nil
}

func (aps *AuditPlanCommencementLetterService) GetByID(auditPlanID uint, letterID uint) (*models.AuditPlanCommencementLetter, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}

	commencementLetter, err := aps.repo.GetByID(letterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, letterID)
		}
		return nil, err
	}

	// Return the found commencement letter
	return commencementLetter, nil
}

func (aps *AuditPlanCommencementLetterService) Create(auditPlanId uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()

	auditPlanDetails, err := aps.auditPlanSvc.GetFullDetailsByID(auditPlanId)
	if err != nil {
		return err
	}

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
		subjectEn := fmt.Sprintf("Audit Commencement Letter - %s", auditPlanDetails.Title)
		subjectAr := fmt.Sprintf("خطاب بدء التدقيق - %s", auditPlanDetails.Title)
		bodyEn := generateLetterBodyEn(auditableItems, auditPlanDetails, fromUser)
		bodyAr := generateLetterBodyAr(auditableItems, auditPlanDetails, fromUser)

		commencementLetter := &models.AuditPlanCommencementLetter{
			AuditPlanID:       auditPlanId,
			ResponsibleUnitID: unitID,
			SubjectEn:         subjectEn,
			SubjectAr:         subjectAr,
			BodyEn:            bodyEn,
			BodyAr:            bodyAr,
			FromUserID:        *auditPlanDetails.TeamLeadID,
			CreatedBYID:       currentUserID,
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

func (aps *AuditPlanCommencementLetterService) Update(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterFormDTO *auditDTOs.AuditPlanCommencementLetterFormDTO) error {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}

	existingCommencementLetter, err := aps.GetByID(auditPlanID, letterID)
	if err != nil {
		return err
	}

	// Update the fields of the commencement letter
	existingCommencementLetter.SubjectEn = commencementLetterFormDTO.SubjectEn
	existingCommencementLetter.SubjectAr = commencementLetterFormDTO.SubjectAr
	existingCommencementLetter.BodyEn = commencementLetterFormDTO.BodyEn
	existingCommencementLetter.BodyAr = commencementLetterFormDTO.BodyAr
	existingCommencementLetter.FromUserID = commencementLetterFormDTO.FromUserID
	existingCommencementLetter.UpdatedBYID = &currentUserID

	// Call the repository to update the commencement letter in the database
	err = aps.repo.Update(existingCommencementLetter)
	if err != nil {
		return err
	}

	return nil
}

func (aps *AuditPlanCommencementLetterService) Delete(auditPlanID uint, letterID uint, currentUserID uint) (uint, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}

	var failedID uint

	_, err := aps.GetByID(auditPlanID, letterID)
	if err != nil {
		return 0, err
	}

	if err := aps.repo.Delete(letterID); err != nil {
		failedID = letterID
		return failedID, err
	}

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditPlanCommencementLetter{}.TableName(),
			letterID,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	// Return nil for no failure
	return failedID, nil
}

func generateLetterBodyEn(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO, fromUser *organizationDTOs.UserDetailsDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleEn != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleEn) // Use English title if available
		}

		if item.AuditableItem.ResponsibleUnit.TitleEn != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleEn
		}
	}

	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameEN, auditPlan.TeamLead.LastNameEN)

	body := fmt.Sprintf(`
Dear %s,

This is to inform you that JethurGRC will commence an audit of your department for the period %s to %s.

The audit scope will include:
%s

Please prepare relevant documents for review. Our team, led by %s, will contact you to schedule an initial meeting.

Your cooperation is appreciated. For any questions, please contact %s.

Thank you,
`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName)

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

		if item.AuditableItem.ResponsibleUnit.TitleAr != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleAr
		}
	}

	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameAR, auditPlan.TeamLead.LastNameAR)

	body := fmt.Sprintf(`
عزيزي %s،

نود إبلاغكم أن JethurGRC سيبدأ عملية التدقيق على إدارتكم خلال الفترة من %s إلى %s.

سيشمل نطاق التدقيق:
%s

يرجى إعداد المستندات ذات الصلة للمراجعة. سيقوم فريقنا، بقيادة %s، بالتواصل معكم لتحديد موعد الاجتماع الأولي.

نقدر تعاونكم. لأي استفسارات، يرجى التواصل مع %s.

شكرًا لكم،

`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName)

	return body
}
