package services

import (
	"fmt"
	"log"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type ProjectSecurityGateServiceInterface interface {
	SetChecklistService(svc ProjectSecurityGateChecklistServiceInterface)
	GetAll(params *dtos.ProjectSecurityGateListQueryParamsDTO) (*dtos.ProjectSecurityGateListResponseDTO, error)
	GetSummary() (*dtos.ProjectSecurityGateSummaryDTO, error)
	GetByID(id uint) (*models.ProjectSecurityGate, error)
	GetFullDetailsByID(id uint) (*dtos.ProjectSecurityGateDetailDTO, error)
	Create(currentUserID uint, form *dtos.ProjectSecurityGateFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ProjectSecurityGateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GoLive(id uint, currentUserID uint) error
}

type ProjectSecurityGateService struct {
	repo                 repositories.ProjectSecurityGateRepositoryInterface
	checklistSvc         ProjectSecurityGateChecklistServiceInterface
	stageSvc             ProjectSecurityGateStageServiceInterface
	classificationSvc    ProjectSecurityGateClassificationServiceInterface
	checklistTemplateSvc ProjectSecurityGateChecklistTemplateServiceInterface
	readinessStatusSvc   ProjectSecurityGateReadinessStatusServiceInterface
	userSvc              organizationServices.UserServiceInterface
	activityLog          logServices.ActivityLogServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
}

func NewProjectSecurityGateService(
	repo repositories.ProjectSecurityGateRepositoryInterface,
	stageSvc ProjectSecurityGateStageServiceInterface,
	classificationSvc ProjectSecurityGateClassificationServiceInterface,
	checklistTemplateSvc ProjectSecurityGateChecklistTemplateServiceInterface,
	readinessStatusSvc ProjectSecurityGateReadinessStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectSecurityGateServiceInterface {
	return &ProjectSecurityGateService{
		repo:                 repo,
		stageSvc:             stageSvc,
		classificationSvc:    classificationSvc,
		checklistTemplateSvc: checklistTemplateSvc,
		readinessStatusSvc:   readinessStatusSvc,
		userSvc:              userSvc,
		activityLog:          activityLog,
		notificationSvc:      notificationSvc,
	}
}

func (psgs *ProjectSecurityGateService) SetChecklistService(svc ProjectSecurityGateChecklistServiceInterface) {
	psgs.checklistSvc = svc
}

func (psgs *ProjectSecurityGateService) GetAll(params *dtos.ProjectSecurityGateListQueryParamsDTO) (*dtos.ProjectSecurityGateListResponseDTO, error) {
	return psgs.repo.GetAll(params)
}

func (psgs *ProjectSecurityGateService) GetSummary() (*dtos.ProjectSecurityGateSummaryDTO, error) {
	return psgs.repo.GetSummary()
}

func (psgs *ProjectSecurityGateService) GetByID(id uint) (*models.ProjectSecurityGate, error) {
	return psgs.repo.GetByID(id)
}

func (psgs *ProjectSecurityGateService) GetFullDetailsByID(id uint) (*dtos.ProjectSecurityGateDetailDTO, error) {
	projectSecurityGate, err := psgs.repo.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	var detail dtos.ProjectSecurityGateDetailDTO
	if err := mapstructure.Decode(projectSecurityGate, &detail); err != nil {
		return nil, err
	}

	detail.CreatedAt = projectSecurityGate.CreatedAt
	detail.UpdatedAt = projectSecurityGate.UpdatedAt

	detail.Classifications = make([]*dtos.ProjectSecurityGateClassificationDTO, 0, len(projectSecurityGate.Classifications))
	for _, mapping := range projectSecurityGate.Classifications {
		if mapping.Classification == nil {
			continue
		}

		var classificationDTO dtos.ProjectSecurityGateClassificationDTO
		if err := mapstructure.Decode(mapping.Classification, &classificationDTO); err != nil {
			return nil, err
		}
		detail.Classifications = append(detail.Classifications, &classificationDTO)
	}

	return &detail, nil
}

func (psgs *ProjectSecurityGateService) Create(currentUserID uint, form *dtos.ProjectSecurityGateFormDTO) (id uint, err error) {
	tx := psgs.repo.BeginTransaction()

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

	stage, err := psgs.stageSvc.GetByID(form.StageID, tx)
	if err != nil {
		return 0, err
	}
	if _, err = psgs.userSvc.GetByIDWithTx(form.BusinessOwnerID, tx); err != nil {
		return 0, err
	}
	if err = psgs.validateClassificationIDs(form.ClassificationIDs, tx); err != nil {
		return 0, err
	}

	template, err := psgs.checklistTemplateSvc.GetByStageID(form.StageID, tx)
	if err != nil {
		return 0, err
	}

	latestPublishedVersion, err := psgs.checklistTemplateSvc.GetLatestPublishedVersionByTemplateID(template.ID, tx)
	if err != nil {
		return 0, err
	}
	if latestPublishedVersion == nil {
		return 0, utility.NewValidationError(
			fmt.Errorf("no published checklist template version exists for this stage"),
			"no_published_version",
			"template version",
		)
	}

	draftReadinessStatus, err := psgs.readinessStatusSvc.GetByType("draft", tx)
	if err != nil {
		return 0, err
	}

	projectSecurityGate := &models.ProjectSecurityGate{
		ProjectTitleEn:    form.ProjectTitleEn,
		ProjectTitleAr:    form.ProjectTitleAr,
		StageID:           form.StageID,
		TemplateVersionID: latestPublishedVersion.ID,
		BusinessOwnerID:   form.BusinessOwnerID,
		ReadinessStatusID: draftReadinessStatus.ID,
		Notes:             form.Notes,
		CreatedBYID:       currentUserID,
		CreatedAt:         time.Now(),
	}

	if err = psgs.repo.Create(projectSecurityGate, tx); err != nil {
		return 0, err
	}

	checklists := make([]*models.ProjectSecurityGateChecklist, 0, len(latestPublishedVersion.Requirements))
	for _, requirement := range latestPublishedVersion.Requirements {
		checklists = append(checklists, &models.ProjectSecurityGateChecklist{
			SecurityGateID: projectSecurityGate.ID,
			RequirementID:  requirement.ID,
		})
	}

	if err = psgs.checklistSvc.CreateChecklists(checklists, tx); err != nil {
		return 0, err
	}

	if err = psgs.repo.ReplaceClassificationMappings(projectSecurityGate.ID, form.ClassificationIDs, tx); err != nil {
		return 0, err
	}

	id = projectSecurityGate.ID

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		activityErr := psgs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGate{}.TableName(),
			id,
			projectSecurityGate,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_create", nil,
		)
		if activityErr != nil {
			log.Printf("Failed to log activity: %v", activityErr)
		}

		psgs.notifyBusinessOwner(id, projectSecurityGate, stage, currentUserID)
	}()

	return id, err
}

func (psgs *ProjectSecurityGateService) notifyBusinessOwner(id uint, securityGate *models.ProjectSecurityGate, stage *models.ProjectSecurityGateStage, currentUserID uint) {
	businessOwner, err := psgs.userSvc.GetByID(securityGate.BusinessOwnerID)
	if err != nil {
		utility.LogError("Failed to fetch business owner for security gate notification", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	createdBy, err := psgs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch creator for security gate notification", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	entity := "project-security-gate"
	action := "business-owner-assign"
	appURL := utility.GetAppURL()

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{securityGate.BusinessOwnerID},
		ModuleIdentifier:     "project",
		NotificationType:     string(commonDTOs.AssignmentType),
		HeadingEn:            "Project Security Gate Assigned",
		HeadingAr:            "تم تعيين بوابة أمان المشروع",
		TitleEn:              securityGate.ProjectTitleEn,
		TitleAr:              securityGate.ProjectTitleAr,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":          businessOwner.FirstNameEn,
		"LastNameEn":           businessOwner.LastNameEn,
		"ProjectTitleEn":       securityGate.ProjectTitleEn,
		"ProjectTitleAr":       securityGate.ProjectTitleAr,
		"Stage":                stage.TitleEn,
		"CreatedByFirstNameEn": createdBy.FirstNameEn,
		"CreatedByLastNameEn":  createdBy.LastNameEn,
		"AssignedOn":           time.Now().Format("2006-01-02"),
		"RedirectURL":          fmt.Sprintf("%s/projects/security-gate/%d", appURL, id),
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{businessOwner.Email},
		Subject:              fmt.Sprintf("Project Security Gate: %s - Business Owner Assigned", securityGate.ProjectTitleEn),
		TemplateRelativePath: "project/templates/email/project_security_gate_create_business_owner_mail.html",
		EmailData:            emailData,
	}

	if err := psgs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_security_gate"); err != nil {
		utility.LogError("Failed to send security gate assignment notification", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (psgs *ProjectSecurityGateService) GoLive(id uint, currentUserID uint) (err error) {
	tx := psgs.repo.BeginTransaction()

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

	securityGate, err := psgs.repo.GetByID(id, tx)
	if err != nil {
		return err
	}

	goLiveStatus, err := psgs.readinessStatusSvc.GetByType("go_live", tx)
	if err != nil {
		return err
	}

	if securityGate.ReadinessStatusID == goLiveStatus.ID {
		err = utility.NewValidationError(
			fmt.Errorf("security gate is already go live"),
			"security_gate_already_go_live",
			"readiness_status",
		)
		return err
	}

	checklists, err := psgs.checklistSvc.GetAllBySecurityGateID(id, tx)
	if err != nil {
		return err
	}

	var hasUnanswered, hasPendingException, hasIncompleteStatus bool

	for _, checklist := range checklists {
		if checklist.ResponseStatusID == nil && !checklist.IsException {
			hasUnanswered = true
			continue
		}

		if checklist.IsException {
			if checklist.ApprovedBYID == nil {
				hasPendingException = true
			}
			continue
		}

		if checklist.Requirement == nil || !checklist.Requirement.IsCritical {
			continue
		}

		if checklist.ResponseStatus == nil ||
			(checklist.ResponseStatus.Type != "implemented" && checklist.ResponseStatus.Type != "not_applicable") {
			hasIncompleteStatus = true
		}
	}

	if len(checklists) == 0 || hasUnanswered {
		err = utility.NewValidationError(
			fmt.Errorf("all checklist items must be answered before go live"),
			"security_gate_checklist_not_answered",
			"checklist",
		)
		return err
	}

	if hasIncompleteStatus {
		err = utility.NewValidationError(
			fmt.Errorf("all checklist items must be implemented or not applicable before go live"),
			"security_gate_checklist_not_implemented",
			"checklist",
		)
		return err
	}

	if hasPendingException {
		err = utility.NewValidationError(
			fmt.Errorf("all checklist exceptions must be approved before go live"),
			"security_gate_exception_not_approved",
			"checklist exception",
		)
		return err
	}

	previousSecurityGate := *securityGate

	now := time.Now()
	securityGate.ReadinessStatusID = goLiveStatus.ID
	securityGate.UpdatedBYID = &currentUserID
	securityGate.UpdatedAt = &now

	if err = psgs.repo.Update(securityGate, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"ReadinessStatusID"}
		changes, err := utility.CompareFields(&previousSecurityGate, securityGate, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
		}

		changedDataJSON, err := utility.ChangesToJSON(changes)
		if err != nil {
			log.Println("error :", err.Error())
		}

		err = psgs.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGate{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"project", "project_security_gate", "project_security_gate_go_live", nil,
		)
		if err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return err
}

func (psgs *ProjectSecurityGateService) Update(id uint, currentUserID uint, form *dtos.ProjectSecurityGateFormDTO) (err error) {
	tx := psgs.repo.BeginTransaction()

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

	existing, err := psgs.repo.GetByID(id, tx)
	if err != nil {
		return err
	}

	if existing.StageID != form.StageID {
		return utility.NewValidationError(
			fmt.Errorf("stage cannot be changed"),
			"stage_change_not_allowed",
			"stage_id",
		)
	}

	if _, err = psgs.stageSvc.GetByID(form.StageID, tx); err != nil {
		return err
	}
	if _, err = psgs.userSvc.GetByIDWithTx(form.BusinessOwnerID, tx); err != nil {
		return err
	}
	if err = psgs.validateClassificationIDs(form.ClassificationIDs, tx); err != nil {
		return err
	}

	oldData := map[string]interface{}{
		"project_title_en":  existing.ProjectTitleEn,
		"project_title_ar":  existing.ProjectTitleAr,
		"stage_id":          existing.StageID,
		"business_owner_id": existing.BusinessOwnerID,
		"notes":             existing.Notes,
	}

	now := time.Now()
	existing.ProjectTitleEn = form.ProjectTitleEn
	existing.ProjectTitleAr = form.ProjectTitleAr
	existing.StageID = form.StageID
	existing.BusinessOwnerID = form.BusinessOwnerID
	existing.Notes = form.Notes
	existing.UpdatedBYID = &currentUserID
	existing.UpdatedAt = &now

	if err = psgs.repo.Update(existing, tx); err != nil {
		return err
	}

	if err = psgs.repo.ReplaceClassificationMappings(id, form.ClassificationIDs, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"project_title_en":  form.ProjectTitleEn,
			"project_title_ar":  form.ProjectTitleAr,
			"stage_id":          form.StageID,
			"business_owner_id": form.BusinessOwnerID,
			"notes":             form.Notes,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			activityErr := psgs.activityLog.LogActivity(
				models.UpdateAction,
				models.ProjectSecurityGate{}.TableName(),
				id,
				changes,
				currentUserID,
				"project", "project_security_gate", "project_security_gate_update", nil,
			)
			if activityErr != nil {
				log.Printf("Failed to log update activity: %v", activityErr)
			}
		}
	}()

	return err
}

func (psgs *ProjectSecurityGateService) Delete(ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := psgs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}
	}()

	var successfulIDs []uint
	var deletedSecurityGates []map[string]interface{}

	for _, id := range ids {
		existing, fetchErr := psgs.repo.GetByID(id, tx)
		if fetchErr != nil {
			err = fetchErr
			return failedIDs, err
		}

		if err = psgs.repo.DeleteClassificationMappingsBySecurityGateID(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err = psgs.checklistSvc.DeleteChecklistsBySecurityGateID(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err = psgs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedSecurityGates = append(deletedSecurityGates, map[string]interface{}{
			"id":               id,
			"project_title_en": existing.ProjectTitleEn,
			"project_title_ar": existing.ProjectTitleAr,
		})
	}

	if len(failedIDs) > 0 {
		return failedIDs, err
	}

	if err = tx.Commit(); err != nil {
		return failedIDs, err
	}

	if len(successfulIDs) > 0 {
		go func() {
			activityErr := psgs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectSecurityGate{}.TableName(),
				successfulIDs,
				deletedSecurityGates,
				currentUserID,
				nil,
			)
			if activityErr != nil {
				log.Printf("Failed to log delete activity: %v", activityErr)
			}
		}()
	}

	return failedIDs, nil
}

func (psgs *ProjectSecurityGateService) validateClassificationIDs(classificationIDs []uint, tx ...utility.DatabaseTransaction) error {
	for _, classificationID := range classificationIDs {
		if _, err := psgs.classificationSvc.GetByID(classificationID, tx...); err != nil {
			return err
		}
	}
	return nil
}
