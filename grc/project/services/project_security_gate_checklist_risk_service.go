package services

import (
	"fmt"
	"log"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectSecurityGateChecklistRiskServiceInterface interface {
	Create(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistRiskFormDTO) (uint, error)
	Update(securityGateID uint, checklistID uint, riskRegisterID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistRiskFormDTO) error
	Finalize(securityGateID uint, checklistID uint, riskRegisterID uint, currentUserID uint) error
}

type ProjectSecurityGateChecklistRiskService struct {
	repo            repositories.ProjectSecurityGateChecklistRiskRepositoryInterface
	checklistSvc    ProjectSecurityGateChecklistServiceInterface
	userSvc         organizationSvc.UserServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
	activityLog     logServices.ActivityLogServiceInterface
}

func NewProjectSecurityGateChecklistRiskService(
	repo repositories.ProjectSecurityGateChecklistRiskRepositoryInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
) ProjectSecurityGateChecklistRiskServiceInterface {
	return &ProjectSecurityGateChecklistRiskService{
		repo:            repo,
		checklistSvc:    checklistSvc,
		userSvc:         userSvc,
		notificationSvc: notificationSvc,
		activityLog:     activityLog,
	}
}

func (psgcrs *ProjectSecurityGateChecklistRiskService) validateChecklistForRisk(securityGateID uint, checklistID uint, tx utility.DatabaseTransaction) error {
	checklist, err := psgcrs.checklistSvc.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return err
	}

	// An exception carries no response status, so it can never qualify.
	if checklist.ResponseStatusID == nil {
		return utility.NewValidationError(
			fmt.Errorf("a risk can only be raised against a partially implemented or not implemented checklist item"),
			"risk_requires_implementation_gap",
			"checklist",
		)
	}

	responseStatus, err := psgcrs.repo.GetChecklistResponseStatusByID(*checklist.ResponseStatusID, tx)
	if err != nil {
		return err
	}

	if responseStatus.Type != "partially_implemented" && responseStatus.Type != "not_implemented" {
		return utility.NewValidationError(
			fmt.Errorf("a risk can only be raised against a partially implemented or not implemented checklist item"),
			"risk_requires_implementation_gap",
			"checklist",
		)
	}

	return nil
}

func (psgcrs *ProjectSecurityGateChecklistRiskService) Create(securityGateID uint, checklistID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistRiskFormDTO) (id uint, err error) {
	registrationDate, err := time.Parse("2006-01-02", form.RegistrationDate)
	if err != nil {
		return 0, utility.NewValidationError(err, "invalid_date", "registration_date")
	}

	tx := psgcrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if err = psgcrs.validateChecklistForRisk(securityGateID, checklistID, tx); err != nil {
		return 0, err
	}

	if _, err = psgcrs.repo.GetRiskCategoryByID(form.RiskCategoryID, tx); err != nil {
		return 0, err
	}

	referenceNumber, err := psgcrs.repo.GenerateReferenceNumber(tx)
	if err != nil {
		return 0, err
	}

	// A raised risk starts as identified and stays editable until it is finalized
	riskStatus, err := psgcrs.repo.GetRiskStatusByType("identified", tx)
	if err != nil {
		return 0, err
	}

	riskMatrix, err := psgcrs.repo.GetRiskMatrixByType("isms", tx)
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		ReferenceNumber:   &referenceNumber,
		RiskStatusID:      riskStatus.ID,
		RiskCategoryID:    &form.RiskCategoryID,
		RiskMatrixID:      riskMatrix.ID,
		Title:             form.Title,
		RegistrationDate:  &registrationDate,
		RiskCode:          &form.RiskCode,
		DivisionID:        form.DivisionID,
		ResponsibleUnitID: &form.ResponsibleUnitID,
		ResponsibleUserID: &form.ResponsibleUserID,
		Justification:     &form.Justification,
		CreatedBYID:       currentUserID,
		CreatedAt:         time.Now(),
		IsCyberRisk:       true,
	}

	id, err = psgcrs.repo.CreateRiskRegister(riskRegister, tx)
	if err != nil {
		return 0, err
	}

	checklistRisk := &models.ProjectSecurityGateChecklistRisk{
		ChecklistID:    checklistID,
		RiskRegisterID: id,
		CreatedBYID:    currentUserID,
		CreatedAt:      time.Now(),
	}

	if err = psgcrs.repo.CreateChecklistRisk(checklistRisk, tx); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		activityErr := psgcrs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGateChecklistRisk{}.TableName(),
			checklistRisk.ID,
			checklistRisk,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_risk_create", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log activity: %v", activityErr)
		}
	}()

	return id, err
}

func (psgcrs *ProjectSecurityGateChecklistRiskService) Update(securityGateID uint, checklistID uint, riskRegisterID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistRiskFormDTO) (err error) {
	registrationDate, err := time.Parse("2006-01-02", form.RegistrationDate)
	if err != nil {
		return utility.NewValidationError(err, "invalid_date", "registration_date")
	}

	tx := psgcrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgcrs.checklistSvc.GetByID(securityGateID, checklistID, tx); err != nil {
		return err
	}

	if _, err = psgcrs.repo.GetChecklistRiskByRiskID(checklistID, riskRegisterID, tx); err != nil {
		return err
	}

	if _, err = psgcrs.repo.GetRiskCategoryByID(form.RiskCategoryID, tx); err != nil {
		return err
	}

	riskRegister, err := psgcrs.repo.GetRiskRegisterByID(riskRegisterID, tx)
	if err != nil {
		return err
	}

	identifiedStatus, err := psgcrs.repo.GetRiskStatusByType("identified", tx)
	if err != nil {
		return err
	}

	// Finalizing moves the risk off identified, and from that point it is read only.
	if riskRegister.RiskStatusID != identifiedStatus.ID {
		return utility.NewValidationError(
			fmt.Errorf("a finalized risk can no longer be edited"),
			"risk_already_finalized",
			"risk",
		)
	}

	oldData := map[string]interface{}{
		"title":               riskRegister.Title,
		"risk_code":           riskRegister.RiskCode,
		"justification":       riskRegister.Justification,
		"registration_date":   riskRegister.RegistrationDate,
		"risk_category_id":    nil,
		"division_id":         riskRegister.DivisionID,
		"responsible_unit_id": nil,
		"responsible_user_id": nil,
	}
	if riskRegister.RiskCategoryID != nil {
		oldData["risk_category_id"] = *riskRegister.RiskCategoryID
	}
	if riskRegister.ResponsibleUnitID != nil {
		oldData["responsible_unit_id"] = *riskRegister.ResponsibleUnitID
	}
	if riskRegister.ResponsibleUserID != nil {
		oldData["responsible_user_id"] = *riskRegister.ResponsibleUserID
	}

	now := time.Now()
	riskRegister.Title = form.Title
	riskRegister.RiskCode = &form.RiskCode
	riskRegister.RegistrationDate = &registrationDate
	riskRegister.Justification = &form.Justification
	riskRegister.RiskCategoryID = &form.RiskCategoryID
	riskRegister.DivisionID = form.DivisionID
	riskRegister.ResponsibleUnitID = &form.ResponsibleUnitID
	riskRegister.ResponsibleUserID = &form.ResponsibleUserID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	newData := map[string]interface{}{
		"title":               riskRegister.Title,
		"risk_code":           riskRegister.RiskCode,
		"justification":       riskRegister.Justification,
		"registration_date":   riskRegister.RegistrationDate,
		"risk_category_id":    form.RiskCategoryID,
		"division_id":         riskRegister.DivisionID,
		"responsible_unit_id": form.ResponsibleUnitID,
		"responsible_user_id": form.ResponsibleUserID,
	}
	changes := utility.CompareFieldChanges(oldData, newData)

	if err = psgcrs.repo.UpdateRiskRegister(riskRegister, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		if len(changes) == 0 {
			return
		}

		activityErr := psgcrs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGateChecklistRisk{}.TableName(),
			riskRegisterID,
			changes,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_risk_update", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log update activity: %v", activityErr)
		}
	}()

	return err
}

func (psgcrs *ProjectSecurityGateChecklistRiskService) Finalize(securityGateID uint, checklistID uint, riskRegisterID uint, currentUserID uint) (err error) {
	tx := psgcrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgcrs.checklistSvc.GetByID(securityGateID, checklistID, tx); err != nil {
		return err
	}

	if _, err = psgcrs.repo.GetChecklistRiskByRiskID(checklistID, riskRegisterID, tx); err != nil {
		return err
	}

	riskRegister, err := psgcrs.repo.GetRiskRegisterByID(riskRegisterID, tx)
	if err != nil {
		return err
	}

	identifiedStatus, err := psgcrs.repo.GetRiskStatusByType("identified", tx)
	if err != nil {
		return err
	}

	if riskRegister.RiskStatusID != identifiedStatus.ID {
		return utility.NewValidationError(
			fmt.Errorf("the risk has already been finalized"),
			"risk_already_finalized",
			"risk",
		)
	}

	openStatus, err := psgcrs.repo.GetRiskStatusByType("open", tx)
	if err != nil {
		return err
	}

	previous := *riskRegister

	now := time.Now()
	riskRegister.RiskStatusID = openStatus.ID
	riskRegister.UpdatedBYID = &currentUserID
	riskRegister.UpdatedAt = &now

	if err = psgcrs.repo.UpdateRiskRegister(riskRegister, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"RiskStatusID"}
		changes, err := utility.CompareFields(&previous, riskRegister, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
		}

		changedDataJSON, err := utility.ChangesToJSON(changes)
		if err != nil {
			log.Println("error :", err.Error())
		}

		activityErr := psgcrs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGateChecklistRisk{}.TableName(),
			riskRegisterID,
			changedDataJSON,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_checklist_risk_finalize", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log finalize activity: %v", activityErr)
		}

		psgcrs.notifyResponsibleUser(riskRegister, currentUserID)
	}()

	return err
}

func (psgcrs *ProjectSecurityGateChecklistRiskService) notifyResponsibleUser(riskRegister *models.RiskRegister, currentUserID uint) {
	if riskRegister.ResponsibleUserID == nil {
		return
	}

	responsibleUser, err := psgcrs.userSvc.GetByID(*riskRegister.ResponsibleUserID)
	if err != nil {
		utility.LogError("Failed to fetch responsible user", map[string]interface{}{"error": err.Error()})
		return
	}

	createdByUser, err := psgcrs.userSvc.GetByID(riskRegister.CreatedBYID)
	if err != nil {
		utility.LogError("Failed to fetch risk creator", map[string]interface{}{"error": err.Error()})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskRegister.ID)
	entity := "asset-risk-register"
	action := "responsible-user-assign"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*riskRegister.ResponsibleUserID},
		ModuleIdentifier:     "asset_risk",
		NotificationType:     string(commonDTOs.AssignmentType),
		HeadingEn:            "Project Security Gate Risk Responsibility Assigned",
		HeadingAr:            "تم تعيين مسؤولية مخاطر بوابة أمان المشروع",
		ReferenceNumber:      utility.GetStringValue(riskRegister.ReferenceNumber),
		TitleEn:              riskRegister.Title,
		TitleAr:              riskRegister.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	emailData := map[string]interface{}{
		"FirstNameEn":          responsibleUser.FirstNameEn,
		"LastNameEn":           responsibleUser.LastNameEn,
		"RiskReferenceNumber":  utility.GetStringValue(riskRegister.ReferenceNumber),
		"RiskTitle":            riskRegister.Title,
		"RiskJustification":    utility.GetStringValue(riskRegister.Justification),
		"CreatedByFirstNameEn": createdByUser.FirstNameEn,
		"CreatedByLastNameEn":  createdByUser.LastNameEn,
		"AssignedOn":           time.Now().Format("2006-01-02"),
		"RedirectURL":          fmt.Sprintf("%s/asset/risks/%d/overview", appURL, riskRegister.ID),
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{responsibleUser.Email},
		Subject:              fmt.Sprintf("Project Security Gate Risk: %s - Responsibility Assigned", riskRegister.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_register_create_responsible_user_mail.html",
		EmailData:            emailData,
	}

	if err := psgcrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
		utility.LogError("Failed to send responsible user assign notification", map[string]interface{}{"error": err.Error()})
	}
}
