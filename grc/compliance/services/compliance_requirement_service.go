package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRequirementServiceInterface interface {
	GetAll(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error)
	GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
	ExportComplianceRequirements(data []*complianceDTOs.ComplianceRequirementListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type ComplianceRequirementService struct {
	repo                             complianceRepositories.ComplianceRequirementRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	complianceRegisterSvc            ComplianceRegisterServiceInterface
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
	severityLevelSvc                 masterServices.SeverityLevelServiceInterface
	controlSvc                       organizationServices.ControlServiceInterface
	strategicObjectiveSvc            organizationServices.StrategicObjectiveServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
	userSettingSvc                   organizationServices.UserSettingServiceInterface
}

func NewComplianceRequirementService(
	repo complianceRepositories.ComplianceRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	strategicObjectiveSvc organizationServices.StrategicObjectiveServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) ComplianceRequirementServiceInterface {
	return &ComplianceRequirementService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		complianceRegisterSvc:            complianceRegisterSvc,
		complianceRequirementCategorySvc: complianceRequirementCategorySvc,
		complianceStatusSvc:              complianceStatusSvc,
		severityLevelSvc:                 severityLevelSvc,
		controlSvc:                       controlSvc,
		strategicObjectiveSvc:            strategicObjectiveSvc,
		userSvc:                          userSvc,
		notificationSvc:                  notificationSvc,
		userSettingSvc:                   userSettingSvc,
	}
}

func (crs *ComplianceRequirementService) GetAll(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	if complianceID != 0 {
		if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
			return nil, err
		}
	}
	return crs.repo.GetAll(complianceID, params)
}

func (crs *ComplianceRequirementService) GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	return crs.repo.GetRequirementsByComplianceID(complianceID, params)
}

func (crs ComplianceRequirementService) GetByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceRequirement, err := crs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error) {
	complianceRequirement, err := crs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var complianceRequirementDTO complianceDTOs.ComplianceRequirementDetailDTO

	complianceRequirement, err := crs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to DTO
	err = mapstructure.Decode(complianceRequirement, &complianceRequirementDTO)
	if err != nil {
		return nil, err
	}

	if complianceRequirement.Deadline != nil {
		deadline := complianceRequirement.Deadline.Format("2006-01-02")
		complianceRequirementDTO.DeadlineString = &deadline
	}

	return &complianceRequirementDTO, nil
}

func (crs *ComplianceRequirementService) Create(complianceID uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	compliance, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return 0, err
	}

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	// Create or retrieve the compliance clause
	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             complianceRequirementFormDTO.ComplianceClauseNo,
		TitleEn:              complianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Calculate the score
	score, err := crs.calculateComplianceRequirementScore(complianceRequirementFormDTO.ComplianceStatusID, complianceRequirementFormDTO.SeverityLevelID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceRequirement := &models.ComplianceRequirement{
		ComplianceRequirementCategoryID: complianceRequirementFormDTO.ComplianceRequirementCategoryID,
		ComplianceStatusID:              complianceRequirementFormDTO.ComplianceStatusID,
		ResponsibleUserID:               complianceRequirementFormDTO.ResponsibleUserID,
		SeverityLevelID:                 complianceRequirementFormDTO.SeverityLevelID,
		TitleEn:                         complianceRequirementFormDTO.TitleEn,
		TitleAr:                         complianceRequirementFormDTO.TitleAr,
		Deadline:                        dates.Deadline,
		CreatedByID:                     currentUserID,
		ComplianceClauseID:              clause.ID,
		Score:                           score,
	}

	id, err := crs.repo.Create(complianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceStatusID := complianceRequirementFormDTO.ComplianceStatusID
	complianceRequirementUpdate := &models.ComplianceRequirementStatusUpdate{
		ComplianceRequirementID: id,
		ComplianceStatusID:      complianceStatusID,
		Comment:                 "Initial compliance status",
		CreatedByID:             currentUserID,
	}

	_, err = crs.repo.CreateInitialComplianceRequirementStatusUpdate(complianceRequirementUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	averageScore, err := crs.repo.CalculateAverageScore(complianceID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update the ComplianceRegister with the new total score
	err = crs.complianceRegisterSvc.UpdateComplianceTotalScore(complianceID, averageScore, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Add Controls
	var controls []*models.ComplianceControl
	for _, controlID := range complianceRequirementFormDTO.ControlIDs {
		control := &models.ComplianceControl{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: &id,
			ControlID:               *controlID,
			CreatedByID:             currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add Strategic Objectives
	var strategicObjectives []*models.ComplianceRequirementStrategicObjective
	for _, strategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		strategicObjective := &models.ComplianceRequirementStrategicObjective{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: id,
			StrategicObjectiveID:    *strategicObjectiveID,
			CreatedByID:             currentUserID,
		}
		strategicObjectives = append(strategicObjectives, strategicObjective)
	}

	if len(strategicObjectives) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectives, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.ComplianceRequirementAssociatedDocument
	for _, documentID := range complianceRequirementFormDTO.DocumentIds {
		document := &models.ComplianceRequirementAssociatedDocument{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: id,
			DocumentID:              *documentID,
			CreatedByID:             currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := crs.repo.CreateComplianceRequirementDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {

		existingComplianceRequirement, err := crs.GetFullDetailsByID(complianceID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			existingComplianceRequirement,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRequirement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		urlParam := fmt.Sprintf(`{"compliance_register_id":%d,"requirement_id":%d}`, complianceID, id)
		queryParam := ""
		entity := "compliance-requirement"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               complianceRequirementFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Requirement Responsiblity Assigned",
			HeadingAr:            "تم تعيين مسؤولية متطلبات الامتثال",
			TitleEn:              complianceRequirementFormDTO.TitleEn,
			TitleAr:              complianceRequirementFormDTO.TitleAr,
			ReferenceNumber:      compliance.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := crs.userSvc.GetByID(complianceRequirementFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		assignedBy, err := crs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		complianceRequirementCategory, err := crs.complianceRequirementCategorySvc.GetByID(complianceRequirementFormDTO.ComplianceRequirementCategoryID)
		if err != nil {
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Assigned", complianceRequirementFormDTO.TitleEn)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                   responsibleUser.FirstNameEn,
			"LastNameEn":                    responsibleUser.LastNameEn,
			"ComplianceReferenceNumber":     compliance.ReferenceNumber,
			"ComplianceTitleEn":             compliance.TitleEn,
			"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
			"RequirementTitleEn":            complianceRequirementFormDTO.TitleEn,
			"ComplianceClauseTitle":         complianceRequirementFormDTO.ComplianceClauseTitle,
			"Deadline":                      dates.Deadline.Format("2006-01-02"),
			"AssignedByFirstNameEn":         assignedBy.FirstNameEn,
			"AssignedByLastNameEn":          assignedBy.LastNameEn,
			"AssignedOn":                    time.Now().Format("2006-01-02"),
			"RedirectURL":                   fmt.Sprintf("%scompliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_requirement_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		// send notification and mail
		err = crs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_requirement")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (crs *ComplianceRequirementService) Update(complianceID uint, id uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	tx := crs.repo.BeginTransaction()

	compliance, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return err
	}

	existingComplianceRequirement, err := crs.GetByID(complianceID, id, tx)
	if err != nil {
		return err
	}

	oldComplianceRequirement, err := crs.GetFullDetailsByID(complianceID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldResponsibleUserID := existingComplianceRequirement.ResponsibleUserID

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return err
	}

	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             complianceRequirementFormDTO.ComplianceClauseNo,
		TitleEn:              complianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	score, err := crs.calculateComplianceRequirementScore(complianceRequirementFormDTO.ComplianceStatusID, complianceRequirementFormDTO.SeverityLevelID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Compliance Requirement
	existingComplianceRequirement.ComplianceRequirementCategoryID = complianceRequirementFormDTO.ComplianceRequirementCategoryID
	existingComplianceRequirement.ComplianceStatusID = complianceRequirementFormDTO.ComplianceStatusID
	existingComplianceRequirement.ResponsibleUserID = complianceRequirementFormDTO.ResponsibleUserID
	existingComplianceRequirement.SeverityLevelID = complianceRequirementFormDTO.SeverityLevelID
	existingComplianceRequirement.TitleEn = complianceRequirementFormDTO.TitleEn
	existingComplianceRequirement.TitleAr = complianceRequirementFormDTO.TitleAr
	existingComplianceRequirement.Deadline = dates.Deadline
	existingComplianceRequirement.Score = score
	existingComplianceRequirement.UpdatedByID = &currentUserID
	existingComplianceRequirement.ComplianceClauseID = clause.ID

	err = crs.repo.Update(existingComplianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Controls
	complianceRequirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirement.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirementFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ComplianceControl
	for _, newControlID := range complianceRequirementFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new compliance control record
			control := &models.ComplianceControl{
				ComplianceRegisterID:    complianceID,
				ComplianceRequirementID: &id,
				ControlID:               *newControlID,
				CreatedByID:             currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range complianceRequirement.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Strategic Objectives
	existingStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirement.StrategicObjectives {
		existingStrategicObjectiveIDs[obj.StrategicObjectiveID] = true
	}

	// Collect new strategic objective IDs
	newStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		newStrategicObjectiveIDs[*obj] = true
	}

	// Add new strategic objectives
	var strategicObjectivesToAdd []*models.ComplianceRequirementStrategicObjective
	for _, newStrategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		if _, exists := existingStrategicObjectiveIDs[*newStrategicObjectiveID]; !exists {
			// Create a new compliance strategic objective record
			strategicObjective := &models.ComplianceRequirementStrategicObjective{
				ComplianceRequirementID: id,
				StrategicObjectiveID:    *newStrategicObjectiveID,
				CreatedByID:             currentUserID,
			}
			strategicObjectivesToAdd = append(strategicObjectivesToAdd, strategicObjective)
		}
	}

	// Save the new strategic objectives
	if len(strategicObjectivesToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectivesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove strategic objectives not in the new list
	var deleteStrategicObjectiveIDs []*uint

	for _, existingStrategicObjective := range complianceRequirement.StrategicObjectives {
		if _, exists := newStrategicObjectiveIDs[existingStrategicObjective.StrategicObjectiveID]; !exists {
			// Append strategic objective to delete list
			deleteStrategicObjectiveIDs = append(deleteStrategicObjectiveIDs, existingStrategicObjective.ID)
		}
	}

	// Delete strategic objectives not in the new list
	if len(deleteStrategicObjectiveIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range complianceRequirement.Documents {
		existingDocumentIDs[doc.DocumentID] = true
	}

	// Collect new document IDs
	newDocumentIDs := make(map[uint]bool)
	for _, doc := range complianceRequirementFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	// Add new documents
	var documentsToAdd []*models.ComplianceRequirementAssociatedDocument
	for _, newDocumentID := range complianceRequirementFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			// Create a new compliance requirement document record
			document := &models.ComplianceRequirementAssociatedDocument{
				ComplianceRegisterID:    complianceID,
				ComplianceRequirementID: id,
				DocumentID:              *newDocumentID,
				CreatedByID:             currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	// Save the new documents
	if len(documentsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove documents not in the new list
	var deleteDocumentIDs []*uint

	for _, existingDocument := range complianceRequirement.Documents {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			// Append document to delete list
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	averageScore, err := crs.repo.CalculateAverageScore(complianceID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = crs.complianceRegisterSvc.UpdateComplianceTotalScore(complianceID, averageScore, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete documents not in the new list
	if len(deleteDocumentIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingRequirement, err := crs.GetFullDetailsByID(complianceID, id)
		if err != nil {
			return
		}
		// Determine control
		controlsAdded := []uint{}
		controlsRemoved := []uint{}
		for id := range newControlIDs {
			if !existingControlIDs[id] {
				controlsAdded = append(controlsAdded, id)
			}
		}
		for id := range existingControlIDs {
			if !newControlIDs[id] {
				controlsRemoved = append(controlsRemoved, id)
			}
		}

		// Determine strategic objectives
		strategicObjectivesAdded := []uint{}
		strategicObjectivesRemoved := []uint{}
		for id := range newStrategicObjectiveIDs {
			if !existingStrategicObjectiveIDs[id] {
				strategicObjectivesAdded = append(strategicObjectivesAdded, id)
			}
		}
		for id := range existingStrategicObjectiveIDs {
			if !newStrategicObjectiveIDs[id] {
				strategicObjectivesRemoved = append(strategicObjectivesRemoved, id)
			}
		}

		// Determine strategic objectives
		documentsAdded := []uint{}
		documentssRemoved := []uint{}
		for id := range newDocumentIDs {
			if !existingDocumentIDs[id] {
				documentsAdded = append(documentsAdded, id)
			}
		}
		for id := range existingDocumentIDs {
			if !newDocumentIDs[id] {
				documentssRemoved = append(documentssRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(controlsAdded) > 0 || len(controlsRemoved) > 0 {
			associatedItems["control_ids"] = map[string]interface{}{
				"added":   controlsAdded,
				"removed": controlsRemoved,
			}
		}

		if len(strategicObjectivesAdded) > 0 || len(strategicObjectivesRemoved) > 0 {
			associatedItems["strategic_objective_ids"] = map[string]interface{}{
				"added":   strategicObjectivesAdded,
				"removed": strategicObjectivesRemoved,
			}
		}

		if len(documentsAdded) > 0 || len(documentssRemoved) > 0 {
			associatedItems["document_ids"] = map[string]interface{}{
				"added":   documentsAdded,
				"removed": documentssRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"compliance_requirement_category": oldComplianceRequirement.ComplianceRequirementCategory.TitleEn,
			"compliance_status":               oldComplianceRequirement.ComplianceStatus.TitleEn,
			"responsible_user":                fmt.Sprintf("%s %s", oldComplianceRequirement.ResponsibleUser.FirstNameEN, oldComplianceRequirement.ResponsibleUser.LastNameEN),
			"severity_level":                  oldComplianceRequirement.SeverityLevel.TitleEn,
			"title":                           oldComplianceRequirement.TitleEn,
			"dead_line":                       oldComplianceRequirement.DeadlineString,
			"compliance_clause":               oldComplianceRequirement.ComplianceClause.Title,
		}

		newData := map[string]interface{}{
			"compliance_requirement_category": existingRequirement.ComplianceRequirementCategory.TitleEn,
			"compliance_status":               existingRequirement.ComplianceStatus.TitleEn,
			"responsible_user":                fmt.Sprintf("%s %s", existingRequirement.ResponsibleUser.FirstNameEN, existingRequirement.ResponsibleUser.LastNameEN),
			"severity_level":                  existingRequirement.SeverityLevel.TitleEn,
			"title":                           existingRequirement.TitleEn,
			"dead_line":                       existingRequirement.DeadlineString,
			"compliance_clause":               existingRequirement.ComplianceClause.Title,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := crs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceRequirement{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceRequirement{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
		// Send notifications if the responsible user has changed
		if oldResponsibleUserID != existingComplianceRequirement.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"compliance_register_id":%d,"requirement_id":%d}`, complianceID, id)
			queryParam := ""
			entity := "compliance-requirement"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Requirement Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية متطلبات الامتثال",
				TitleEn:              existingComplianceRequirement.TitleEn,
				ReferenceNumber:      compliance.ReferenceNumber,
				TitleAr:              existingComplianceRequirement.TitleEn,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               complianceRequirementFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Requirement Responsibility Assigned",
				HeadingAr:            "تم تعيين المسؤولية لمتطلبات الامتثال",
				TitleEn:              existingComplianceRequirement.TitleEn,
				TitleAr:              existingComplianceRequirement.TitleEn,
				ReferenceNumber:      compliance.ReferenceNumber,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := crs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := crs.userSvc.GetByID(complianceRequirementFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := crs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			complianceRequirementCategory, err := crs.complianceRequirementCategorySvc.GetByID(complianceRequirementFormDTO.ComplianceRequirementCategoryID)
			if err != nil {
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Removed", existingComplianceRequirement.TitleEn)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":                   oldResponsibleUser.FirstNameEn,
				"LastNameEn":                    oldResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber":     compliance.ReferenceNumber,
				"ComplianceTitle":               compliance.TitleEn,
				"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
				"RequirementTitleEn":            existingComplianceRequirement.TitleEn,
				"ComplianceClauseTitle":         complianceClause.TitleEn,
				"Deadline":                      dates.Deadline.Format("2006-01-02"),
				"UpdatedByFirstNameEn":          updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":           updatedByUser.LastNameEn,
				"AssignedOn":                    time.Now().Format("2006-01-02"),
				"RedirectURL":                   fmt.Sprintf("%scompliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "compliance/templates/email/compliance_requirement_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Assigned", existingComplianceRequirement.TitleEn)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":                   newResponsibleUser.FirstNameEn,
				"LastNameEn":                    newResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber":     compliance.ReferenceNumber,
				"ComplianceTitle":               compliance.TitleEn,
				"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
				"RequirementTitleEn":            existingComplianceRequirement.TitleEn,
				"ComplianceClauseTitle":         complianceClause.TitleEn,
				"Deadline":                      dates.Deadline.Format("2006-01-02"),
				"UpdatedByFirstNameEn":          updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":           updatedByUser.LastNameEn,
				"AssignedOn":                    time.Now().Format("2006-01-02"),
				"RedirectURL":                   fmt.Sprintf("%scompliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "compliance/templates/email/compliance_requirement_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = crs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_requirement")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = crs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "compliance_requirement")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil
}

func (crs *ComplianceRequirementService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	compliance, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceRequirement []map[string]interface{}

	var notifications []struct {
		OwnerID          *uint
		Title            string
		RequirementTitle string
		ReferenceNumber  string
	}

	// Track clause IDs for deletion
	clauseIDsToDelete := make(map[uint]bool)

	for _, id := range ids {
		requirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID          *uint
			Title            string
			RequirementTitle string
			ReferenceNumber  string
		}{
			OwnerID:          &requirement.ResponsibleUserID,
			Title:            compliance.TitleEn,
			RequirementTitle: requirement.TitleEn,
			ReferenceNumber:  compliance.ReferenceNumber,
		})

		// Check if controls exist
		if len(requirement.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range requirement.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Check if strategic objectives exist
		if len(requirement.StrategicObjectives) > 0 {
			var deleteStrategicObjectiveIds []*uint
			for _, obj := range requirement.StrategicObjectives {
				deleteStrategicObjectiveIds = append(deleteStrategicObjectiveIds, obj.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(requirement.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, document := range requirement.Documents {
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		complianceRequirementStatusUpdates, err := crs.repo.GetStatusUpdateByRequirementID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if len(complianceRequirementStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range complianceRequirementStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if err := crs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
			// Track the clause IDs
			clauseIDsToDelete[uint(requirement.ComplianceClause.ID)] = true
		}

		deletedComplianceRequirement = append(deletedComplianceRequirement, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": complianceID,
			"title":                  requirement.TitleEn,
		})

	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Delete clauses if all their requirements are deleted
	for clauseID := range clauseIDsToDelete {
		count, err := crs.repo.CountRequirementsByClauseID(clauseID, tx)
		if err != nil {
			failedIds = append(failedIds, clauseID)
			fmt.Printf("Error counting requirements for clause %d: %v\n", clauseID, err)
			continue
		} else if count == 0 {
			if err := crs.repo.DeleteClause(clauseID, tx); err != nil {
				failedIds = append(failedIds, clauseID)
				// Log error, but proceed with other deletions
				fmt.Printf("Error deleting clause %d: %v\n", clauseID, err)
				continue
			}
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Log successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirement{}.TableName(),
				successfulIds,
				deletedComplianceRequirement,
				currentUserID,
				nil,
			)
			if err != nil {
				// Log error, but proceed as it might be recoverable later
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
			// Send DB and email notifications for each successfully deleted compliance register
			for _, notification := range notifications {
				urlParam := ""
				entity := "compliance-requirement"
				action := "delete"

				// Prepare DB notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Compliance Requirement Deleted",
					HeadingAr:            "تم حذف متطلبات الامتثال",
					TitleEn:              notification.RequirementTitle,
					TitleAr:              notification.RequirementTitle,
					ReferenceNumber:      compliance.ReferenceNumber,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Fetch owner details (responsible user)
				ownerUser, err := crs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Compliance Requirement/%s: Deleted", notification.RequirementTitle)
				templateRelativePath := "compliance/templates/email/compliance_requirement_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":               ownerUser.FirstNameEn,
					"LastNameEn":                ownerUser.LastNameEn,
					"ComplianceReferenceNumber": notification.ReferenceNumber,
					"ComplianceTitle":           notification.Title,
					"RequirementTitleEn":        notification.RequirementTitle,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = crs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "compliance_requirement")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}()
	}

	tx.Commit()

	return failedIds, nil
}

func (crs *ComplianceRequirementService) StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingComplianceRequirement, err := crs.GetByID(complianceID, id)
	if err != nil {
		return err
	}

	existingComplianceRequirement.ComplianceStatusID = complianceStatusID
	existingComplianceRequirement.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingComplianceRequirement, tx)
	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRequirementService) calculateComplianceRequirementScore(complianceStatusID uint, severityLevelID uint, tx ...utility.DatabaseTransaction) (int, error) {
	complianceStatus, err := crs.complianceStatusSvc.GetByID(complianceStatusID)
	if err != nil {
		return 0, err
	}

	severityLevel, err := crs.severityLevelSvc.GetByID(severityLevelID)
	if err != nil {
		return 0, err
	}

	maxSeverityWeight, err := crs.severityLevelSvc.GetMaxSeverityWeight(tx...)
	if err != nil {
		return 0, err
	}

	maxBaseScore, err := crs.complianceStatusSvc.GetMaxBaseScore()
	if err != nil {
		return 0, err
	}

	maxPossibleScore := maxSeverityWeight * maxBaseScore

	// Calculate the actual score
	rawScore := complianceStatus.BaseScore * severityLevel.Weight

	// Normalize the score to be out of 100
	normalizedScore := (float64(rawScore) / float64(maxPossibleScore)) * 100

	return int(normalizedScore), nil
}

func (crs *ComplianceRequirementService) validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	if _, err := crs.complianceRequirementCategorySvc.GetByID(ComplianceRequirementFormDTO.ComplianceRequirementCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance requirement category does not exist"), "not_exist", "compliance requirement category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(ComplianceRequirementFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if _, err := crs.severityLevelSvc.GetByID(ComplianceRequirementFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	controlIds := ComplianceRequirementFormDTO.ControlIDs

	controlIDs := make(map[uint]bool)

	for _, controlId := range controlIds {
		if controlIDs[*controlId] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", controlId), "duplicate", "control")
		}

		if _, err := crs.controlSvc.GetByID(*controlId); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlId), "not_exist", "control")
		}

		controlIDs[*controlId] = true
	}

	strategicObjectiveIds := ComplianceRequirementFormDTO.StrategicObjectiveIDs

	strategicObjectiveIDs := make(map[uint]bool)

	for _, strategicObjectiveId := range strategicObjectiveIds {
		if strategicObjectiveIDs[*strategicObjectiveId] {
			return utility.NewValidationError(fmt.Errorf("duplicate strategic objective ID: %d", strategicObjectiveId), "duplicate", "strategic objective")
		}

		if _, err := crs.strategicObjectiveSvc.GetByID(*strategicObjectiveId); err != nil {
			return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveId), "not_exist", "strategic objective")
		}

		strategicObjectiveIDs[*strategicObjectiveId] = true
	}

	return nil
}

func (crs *ComplianceRequirementService) parseComplianceRequirementDates(dto *complianceDTOs.ComplianceRequirementFormDTO) (*complianceDTOs.ComplianceRequirementParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRequirementParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	return parsed, nil
}

func (srs *ComplianceRequirementService) ExportComplianceRequirements(data []*complianceDTOs.ComplianceRequirementListDTO, currentUserID uint) (*os.File, error) {
	// Retrieve user settings for language preference
	userSetting, err := srs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Clause No",
			"Clause",
			"Title (AR)",
			"Category (AR)",
			"Severity Level (AR)",
			"Responsible User (AR)",
			"Deadline",
			"Compliance Status (AR)",
		}
	} else {
		headers = []string{
			"SL No.",
			"Clause No",
			"Clause",
			"Title (EN)",
			"Category (EN)",
			"Severity Level (EN)",
			"Responsible User (EN)",
			"Deadline",
			"Compliance Status (EN)",
		}
	}

	// Convert data to [][]interface{} based on language preference
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				item.ComplianceClauseNo,
				item.ComplianceClauseTitleAr,
				item.TitleAr,
				item.ComplianceRequirementCategoryTitleAr,
				item.SeverityLevelTitleAr,
				item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr,
				item.Deadline,
				item.ComplianceStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				item.ComplianceClauseNo,
				item.ComplianceClauseTitleEn,
				item.TitleEn,
				item.ComplianceRequirementCategoryTitleEn,
				item.SeverityLevelTitleEn,
				item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn,
				item.Deadline,
				item.ComplianceStatusTitleEn,
			}
		}
	}
	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
