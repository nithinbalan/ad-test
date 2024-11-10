package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	assetServices "jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"

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
	UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO) error
	DeleteResponseStrategy(riskID uint, currentUserID uint) error
	GetControlPerformance(riskID uint) (map[string]float64, error)
	UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction, evaluationType string) error
	ClearRiskRating(riskID uint, updateDTO *riskDTOs.RiskRegisterRiskRatingUpdateDTO, tx utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string) (string, error)

	ExportRiskRegister(data []*riskDTOs.RiskRegisterListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type RiskRegisterService struct {
	repo                     riskRepositories.RiskRegisterRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	riskCategorySvc          RiskCategoryServiceInterface
	riskSourceSvc            RiskSourceServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	riskStatusSvc            RiskStatusServiceInterface
	reviewFrequencySvc       masterServices.ReviewFrequencyServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	strategicObjectiveSvc    organizationServices.StrategicObjectiveServiceInterface
	slaRegisterSvc           complianceServices.SLARegisterServiceInterface
	riskMatrixSvc            RiskMatrixServiceInterface
	riskTimelineSvc          RiskTimelineServiceInterface
	userSettingSvc           organizationServices.UserSettingServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
	assetRegisterSvc         assetServices.AssetRegisterServiceInterface
}

func NewRiskRegisterService(
	repo riskRepositories.RiskRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	riskSourceSvc RiskSourceServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	strategicObjectiveSvc organizationServices.StrategicObjectiveServiceInterface,
	slaRegisterSvc complianceServices.SLARegisterServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	assetRegisterSvc assetServices.AssetRegisterServiceInterface,
) RiskRegisterServiceInterface {
	return &RiskRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		riskCategorySvc:          riskCategorySvc,
		riskSourceSvc:            riskSourceSvc,
		userSvc:                  userSvc,
		riskStatusSvc:            riskStatusSvc,
		reviewFrequencySvc:       reviewFrequencySvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		processSvc:               processSvc,
		strategicObjectiveSvc:    strategicObjectiveSvc,
		slaRegisterSvc:           slaRegisterSvc,
		riskMatrixSvc:            riskMatrixSvc,
		riskTimelineSvc:          riskTimelineSvc,
		userSettingSvc:           userSettingSvc,
		notificationSvc:          notificationSvc,
		assetRegisterSvc:         assetRegisterSvc,
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
	var riskDTO riskDTOs.RiskRegisterDetailDTO

	risk, err := rrs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(risk, &riskDTO)
	if err != nil {
		return nil, err
	}

	if risk.ReportedAt != nil {
		reportedAt := risk.ReportedAt.Format("2006-01-02")
		riskDTO.ReportedAtString = &reportedAt
	}

	if risk.NextReviewDate != nil {
		nextReviewDate := risk.NextReviewDate.Format("2006-01-02")
		riskDTO.NextReviewDateString = &nextReviewDate
	}

	riskDTO.CreatedAt = risk.CreatedAt
	riskDTO.UpdatedAt = risk.UpdatedAt

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

	responsibleUnitID, err := rrs.userSvc.GetUserPrimaryUnitID(uint(riskRegisterFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
	}

	riskMatrix, err := rrs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	riskStatus, err := rrs.riskStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		ReferenceNumber:   referenceNumber,
		RiskCategoryID:    riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:      riskRegisterFormDTO.RiskSourceID,
		RiskMatrixID:      riskMatrix.ID,
		ResponsibleUnitID: responsibleUnitID,
		ResponsibleUserID: riskRegisterFormDTO.ResponsibleUserID,
		RiskStatusID:      riskStatus.ID,
		Title:             riskRegisterFormDTO.Title,
		Description:       riskRegisterFormDTO.Description,
		ReportedAt:        dates.ReportedAt,
		ReviewFrequencyID: riskRegisterFormDTO.ReviewFrequencyID,
		NextReviewDate:    dates.NextReviewDate,
		LastReviewNote:    riskRegisterFormDTO.LastReviewNote,
		CreatedBYID:       currentUserID,
	}

	id, err := rrs.repo.Create(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save ComplianceRequirements
	var complianceRequirements []*models.RiskComplianceRequirement
	for _, requirementID := range riskRegisterFormDTO.ComplianceRequirementIDs {
		complianceRequirement, err := rrs.complianceRequirementSvc.GetByComplianceRequirementID(*requirementID)
		if err != nil {
			return 0, err
		}
		requirement := &models.RiskComplianceRequirement{
			RiskRegisterID:          id,
			ComplianceRegisterID:    complianceRequirement.ComplianceClause.ComplianceRegisterID,
			ComplianceRequirementID: *requirementID,
			CreatedByID:             currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := rrs.repo.CreateRiskComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Controls
	var controls []*models.RiskControl
	for _, controlID := range riskRegisterFormDTO.ControlIDs {
		control := &models.RiskControl{
			RiskRegisterID: id,
			ControlID:      *controlID,
			CreatedByID:    currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := rrs.repo.CreateRiskControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Processes
	var processes []*models.RiskProcess
	for _, processID := range riskRegisterFormDTO.ProcessIDs {
		process := &models.RiskProcess{
			RiskRegisterID: id,
			ProcessID:      *processID,
			CreatedByID:    currentUserID,
		}
		processes = append(processes, process)
	}

	if len(processes) > 0 {
		if err := rrs.repo.CreateRiskProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add Strategic Objectives
	var strategicObjectives []*models.RiskStrategicObjective
	for _, strategicObjectiveID := range riskRegisterFormDTO.StrategicObjectiveIDs {
		strategicObjective := &models.RiskStrategicObjective{
			RiskRegisterID:       id,
			StrategicObjectiveID: *strategicObjectiveID,
			CreatedByID:          currentUserID,
		}
		strategicObjectives = append(strategicObjectives, strategicObjective)
	}

	if len(strategicObjectives) > 0 {
		if err := rrs.repo.CreateRiskStrategicObjectives(strategicObjectives, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add SLA Register
	var slaRegisters []*models.RiskSLA
	for _, slaRegisterID := range riskRegisterFormDTO.SLARegisterIDs {
		slaRegister := &models.RiskSLA{
			RiskRegisterID: id,
			SLARegisterID:  *slaRegisterID,
			CreatedByID:    currentUserID,
		}
		slaRegisters = append(slaRegisters, slaRegister)
	}

	if len(slaRegisters) > 0 {
		if err := rrs.repo.CreateRiskSLAs(slaRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var associatedDocs []*models.RiskAssociatedDocument
	for _, documentID := range riskRegisterFormDTO.DocumentIds {
		document := &models.RiskAssociatedDocument{
			RiskRegisterID: id,
			DocumentID:     *documentID,
			CreatedByID:    currentUserID,
		}
		associatedDocs = append(associatedDocs, document)
	}

	if len(associatedDocs) > 0 {
		if err := rrs.repo.CreateRiskDocuments(associatedDocs, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var assets []*models.RiskAsset
	for _, ast := range riskRegisterFormDTO.AssetIds {
		asset := &models.RiskAsset{
			RiskRegisterID:  id,
			AssetRegisterID: *ast,
			CreatedByID:     currentUserID,
		}
		assets = append(assets, asset)
	}

	if len(assets) > 0 {
		if err := rrs.repo.CreateRiskAssets(assets, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add timeline
	// todo
	timelineTypes, err := rrs.riskTimelineSvc.GetAllTypes()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get timeline statuses
	completedStatus, err := rrs.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	pendingStatus, err := rrs.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var timelineFormDTOs []*riskDTOs.RiskTimelineFormDTO

	for _, timelineType := range timelineTypes {
		var (
			statusID      uint
			performedByID *uint
			currentTime   *time.Time
		)

		if timelineType.Type == "identified" {
			statusID = completedStatus.ID
			performedByID = &currentUserID
			now := time.Now()
			currentTime = &now
		} else {
			statusID = pendingStatus.ID
			performedByID = nil
			currentTime = nil
		}

		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			RiskRegisterID:     id,
			RiskTimelineTypeID: timelineType.ID,
			TimelineStatusID:   statusID,
			PerformedByID:      performedByID,
			PerformedAt:        currentTime,
		}
		timelineFormDTOs = append(timelineFormDTOs, timelineFormDTO)
	}
	if err := rrs.riskTimelineSvc.Create(timelineFormDTOs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

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
		url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), "documents")

		// Create a new RiskDocument instance
		document := &models.RiskDocument{
			RiskRegisterID: id,
			URL:            url,
			Name:           doc.FileName,
			Size:           doc.Size,
			Token:          encryptedFilename,
			OriginalName:   doc.OriginalName,
			Type:           doc.Type,
			CreatedByID:    currentUserID,
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

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "risk-register"
		action := "responsible-user-assign"
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               riskRegisterFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Register Responsibility Assigned",
			HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
			TitleEn:              riskRegisterFormDTO.Title,
			ReferenceNumber:      referenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}
		responsibleUser, err := rrs.userSvc.GetByID(riskRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		createdBy, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", riskRegister.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":          responsibleUser.FirstNameEn,
			"LastNameEn":           responsibleUser.LastNameEn,
			"RiskReferenceNumber":  riskRegister.ReferenceNumber,
			"RiskTitle":            riskRegister.Title,
			"RiskDescription":      riskRegister.Description,
			"CreatedByFirstNameEn": createdBy.FirstNameEn,
			"CreatedByLastNameEn":  createdBy.LastNameEn,
			"AssignedOn":           time.Now().Format("2006-01-02"),
			"RedirectURL":          fmt.Sprintf("%srisk/risk-register/%d/overview", appURL, id),
		}
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_register_create_responsible_user_mail.html",
			EmailData:            emailData,
		}
		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
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

	oldResponsibleUserID := existingRiskRegister.ResponsibleUserID

	if existingRiskRegister.ResponsibleUserID != riskRegisterFormDTO.ResponsibleUserID {
		responsibleUnitID, err := rrs.userSvc.GetUserPrimaryUnitID(uint(riskRegisterFormDTO.ResponsibleUserID))
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingRiskRegister.ResponsibleUnitID = responsibleUnitID
	}

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.ResponsibleUserID = riskRegisterFormDTO.ResponsibleUserID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReportedAt = dates.ReportedAt
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.NextReviewDate = dates.NextReviewDate
	existingRiskRegister.LastReviewNote = riskRegisterFormDTO.LastReviewNote
	existingRiskRegister.UpdatedBYID = &currentUserID

	if err := rrs.repo.Update(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range existingRisk.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	newControlIDs := make(map[uint]bool)
	for _, ctrl := range riskRegisterFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	var controlsToAdd []*models.RiskControl
	for _, newControlID := range riskRegisterFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			control := &models.RiskControl{
				RiskRegisterID: id,
				ControlID:      *newControlID,
				CreatedByID:    currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	if len(controlsToAdd) > 0 {
		if err := rrs.repo.CreateRiskControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteControlIDs []*uint

	for _, existingControl := range existingRisk.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	if len(deleteControlIDs) > 0 {
		if err := rrs.repo.DeleteRiskControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingComplianceRequirementIDs := make(map[uint]bool)
	for _, req := range existingRisk.ComplianceRequirements {
		existingComplianceRequirementIDs[req.ComplianceRequirementID] = true
	}

	// Collect new compliance requirement IDs
	newComplianceRequirementIDs := make(map[uint]bool)
	for _, req := range riskRegisterFormDTO.ComplianceRequirementIDs {
		newComplianceRequirementIDs[*req] = true
	}

	// Add new compliance requirements
	var complianceRequirementsToAdd []*models.RiskComplianceRequirement
	for _, newComplianceRequirementID := range riskRegisterFormDTO.ComplianceRequirementIDs {
		if _, exists := existingComplianceRequirementIDs[*newComplianceRequirementID]; !exists {
			// Create a new risk compliance requirement record
			req, err := rrs.complianceRequirementSvc.GetByComplianceRequirementID(*newComplianceRequirementID)
			if err != nil {
				return err
			}
			complianceRequirement := &models.RiskComplianceRequirement{
				RiskRegisterID:          id,
				ComplianceRegisterID:    req.ComplianceClause.ComplianceRegisterID,
				ComplianceRequirementID: *newComplianceRequirementID,
				CreatedByID:             currentUserID,
			}
			complianceRequirementsToAdd = append(complianceRequirementsToAdd, complianceRequirement)
		}
	}

	// Save the new compliance requirements
	if len(complianceRequirementsToAdd) > 0 {
		if err := rrs.repo.CreateRiskComplianceRequirements(complianceRequirementsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove compliance requirements not in the new list
	var deleteComplianceRequirementIDs []*uint

	for _, existingComplianceRequirement := range existingRisk.ComplianceRequirements {
		if _, exists := newComplianceRequirementIDs[existingComplianceRequirement.ComplianceRequirementID]; !exists {
			// Append compliance requirement to delete list
			deleteComplianceRequirementIDs = append(deleteComplianceRequirementIDs, existingComplianceRequirement.ID)
		}
	}

	if len(deleteComplianceRequirementIDs) > 0 {
		if err := rrs.repo.DeleteRiskComplianceRequirements(deleteComplianceRequirementIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingProcessIDs := make(map[uint]bool)
	for _, proc := range existingRisk.Processes {
		existingProcessIDs[proc.ProcessID] = true
	}

	// Collect new process IDs
	newProcessIDs := make(map[uint]bool)
	for _, proc := range riskRegisterFormDTO.ProcessIDs {
		newProcessIDs[*proc] = true
	}

	// Add new processes
	var processesToAdd []*models.RiskProcess
	for _, newProcessID := range riskRegisterFormDTO.ProcessIDs {
		if _, exists := existingProcessIDs[*newProcessID]; !exists {
			// Create a new risk process record
			process := &models.RiskProcess{
				RiskRegisterID: id,
				ProcessID:      *newProcessID,
				CreatedByID:    currentUserID,
			}
			processesToAdd = append(processesToAdd, process)
		}
	}

	// Save the new processes
	if len(processesToAdd) > 0 {
		if err := rrs.repo.CreateRiskProcesses(processesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove processes not in the new list
	var deleteProcessIDs []*uint

	for _, existingProcess := range existingRisk.Processes {
		if _, exists := newProcessIDs[existingProcess.ProcessID]; !exists {
			// Append process to delete list
			deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
		}
	}

	// Delete processes not in the new list
	if len(deleteProcessIDs) > 0 {
		if err := rrs.repo.DeleteRiskProcesses(deleteProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingSLARegisterIDs := make(map[uint]bool)
	for _, sla := range existingRisk.SLARegisters {
		existingSLARegisterIDs[sla.SLARegisterID] = true
	}

	// Collect new slaess IDs
	newSLARegisterIDs := make(map[uint]bool)
	for _, sla := range riskRegisterFormDTO.SLARegisterIDs {
		newSLARegisterIDs[*sla] = true
	}

	// Add new processes
	var slaRegistersToAdd []*models.RiskSLA
	for _, newSLARegisterID := range riskRegisterFormDTO.SLARegisterIDs {
		if _, exists := existingSLARegisterIDs[*newSLARegisterID]; !exists {
			// Create a new risk slaRegister record
			slaRegister := &models.RiskSLA{
				RiskRegisterID: id,
				SLARegisterID:  *newSLARegisterID,
				CreatedByID:    currentUserID,
			}
			slaRegistersToAdd = append(slaRegistersToAdd, slaRegister)
		}
	}

	// Save the new slaRegisters
	if len(slaRegistersToAdd) > 0 {
		if err := rrs.repo.CreateRiskSLAs(slaRegistersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove slaRegisters not in the new list
	var deleteSLARegisterIDs []*uint

	for _, existingSLARegister := range existingRisk.SLARegisters {
		if _, exists := newSLARegisterIDs[existingSLARegister.SLARegisterID]; !exists {
			// Append slaRegister to delete list
			deleteSLARegisterIDs = append(deleteSLARegisterIDs, existingSLARegister.ID)
		}
	}

	// Delete slaRegisters not in the new list
	if len(deleteSLARegisterIDs) > 0 {
		if err := rrs.repo.DeleteRiskSLAs(deleteSLARegisterIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Update Strategic Objectives
	existingStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range existingRisk.StrategicObjectives {
		existingStrategicObjectiveIDs[obj.StrategicObjectiveID] = true
	}

	// Collect new strategic objective IDs
	newStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range riskRegisterFormDTO.StrategicObjectiveIDs {
		newStrategicObjectiveIDs[*obj] = true
	}

	// Add new strategic objectives
	var strategicObjectivesToAdd []*models.RiskStrategicObjective
	for _, newStrategicObjectiveID := range riskRegisterFormDTO.StrategicObjectiveIDs {
		if _, exists := existingStrategicObjectiveIDs[*newStrategicObjectiveID]; !exists {
			// Create a new risk strategic objective record
			strategicObjective := &models.RiskStrategicObjective{
				RiskRegisterID:       id,
				StrategicObjectiveID: *newStrategicObjectiveID,
				CreatedByID:          currentUserID,
			}
			strategicObjectivesToAdd = append(strategicObjectivesToAdd, strategicObjective)
		}
	}

	// Save the new strategic objectives
	if len(strategicObjectivesToAdd) > 0 {
		if err := rrs.repo.CreateRiskStrategicObjectives(strategicObjectivesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove strategic objectives not in the new list
	var deleteStrategicObjectiveIDs []*uint

	for _, existingStrategicObjective := range existingRisk.StrategicObjectives {
		if _, exists := newStrategicObjectiveIDs[existingStrategicObjective.StrategicObjectiveID]; !exists {
			// Append strategic objective to delete list
			deleteStrategicObjectiveIDs = append(deleteStrategicObjectiveIDs, existingStrategicObjective.ID)
		}
	}

	// Delete strategic objectives not in the new list
	if len(deleteStrategicObjectiveIDs) > 0 {
		if err := rrs.repo.DeleteRiskStrategicObjectives(deleteStrategicObjectiveIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range existingRisk.AssociatedDocuments {
		existingDocumentIDs[doc.DocumentID] = true
	}

	// Collect new document IDs
	newDocumentIDs := make(map[uint]bool)
	for _, doc := range riskRegisterFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	// Add new documents
	var documentsToAdd []*models.RiskAssociatedDocument
	for _, newDocumentID := range riskRegisterFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			// Create a new risk document record
			document := &models.RiskAssociatedDocument{
				RiskRegisterID: id,
				DocumentID:     *newDocumentID,
				CreatedByID:    currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	// Save the new documents
	if len(documentsToAdd) > 0 {
		if err := rrs.repo.CreateRiskDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove documents not in the new list
	var deleteDocumentIDs []*uint

	for _, existingDocument := range existingRisk.AssociatedDocuments {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			// Append document to delete list
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	// Delete documents not in the new list
	if len(deleteDocumentIDs) > 0 {
		if err := rrs.repo.DeleteRiskDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingAssetIDs := make(map[uint]bool)
	for _, asset := range existingRisk.Assets {
		existingAssetIDs[asset.AssetRegisterID] = true
	}

	// Collect new asset IDs
	newAssetIDs := make(map[uint]bool)
	for _, assetID := range riskRegisterFormDTO.AssetIds {
		newAssetIDs[*assetID] = true
	}

	// Add new assets
	var assetsToAdd []*models.RiskAsset
	for _, newAssetID := range riskRegisterFormDTO.AssetIds {
		if _, exists := existingAssetIDs[*newAssetID]; !exists {
			// Create a new risk asset record
			asset := &models.RiskAsset{
				RiskRegisterID:  id,
				AssetRegisterID: *newAssetID,
				CreatedByID:     currentUserID,
			}
			assetsToAdd = append(assetsToAdd, asset)
		}
	}

	// Save the new assets
	if len(assetsToAdd) > 0 {
		if err := rrs.repo.CreateRiskAssets(assetsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove assets not in the new list
	var deleteAssetIDs []*uint
	for _, existingAsset := range existingRisk.Assets {
		if _, exists := newAssetIDs[existingAsset.AssetRegisterID]; !exists {
			// Append asset ID to delete list
			deleteAssetIDs = append(deleteAssetIDs, existingAsset.ID)
		}
	}

	// Delete assets not in the new list
	if len(deleteAssetIDs) > 0 {
		if err := rrs.repo.DeleteRiskAssets(deleteAssetIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

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

			url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new RiskDocument instance
			document := &models.RiskDocument{
				RiskRegisterID: id,
				URL:            url,
				Name:           doc.FileName,
				Size:           doc.Size,
				Token:          encryptedFilename,
				OriginalName:   doc.OriginalName,
				Type:           doc.Type,
				CreatedByID:    currentUserID,
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

	go func() {
		// Log the update activity
		err := rrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			id,
			riskRegisterFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingRiskRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "risk-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Risk Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل المخاطر",
				TitleEn:              existingRiskRegister.Title,
				ReferenceNumber:      existingRiskRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               riskRegisterFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Risk Register Responsibility Assigned",
				HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
				TitleEn:              existingRiskRegister.Title,
				ReferenceNumber:      existingRiskRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := rrs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := rrs.userSvc.GetByID(riskRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := rrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Removed", existingRiskRegister.Title)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":          oldResponsibleUser.FirstNameEn,
				"LastNameEn":           oldResponsibleUser.LastNameEn,
				"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
				"RiskTitle":            existingRiskRegister.Title,
				"RiskDescription":      existingRiskRegister.Description,
				"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
				"RedirectURL":          fmt.Sprintf("%srisk/risk-register/%d/overview", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "risk/templates/email/risk_register_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", existingRiskRegister.Title)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"RiskReferenceNumber":   existingRiskRegister.ReferenceNumber,
				"RiskTitle":             existingRiskRegister.Title,
				"RiskDescription":       existingRiskRegister.Description,
				"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":  updatedByUser.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/risk-register/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "risk/templates/email/risk_register_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = rrs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = rrs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

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
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
	}
	documentsToDelete := make(map[uint][]*riskDTOs.RiskDocumentDTO)

	for _, id := range ids {
		risk, err := rrs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
		}{
			OwnerID:         &risk.ResponsibleUserID,
			TitleEn:         risk.Title,
			Description:     *risk.Description,
			ReferenceNumber: risk.ReferenceNumber,
		})
		// Check if controls exist
		if len(risk.Controls) > 0 {
			var deleteControlIDs []*uint
			for _, ctrl := range risk.Controls {
				deleteControlIDs = append(deleteControlIDs, ctrl.ID)
			}
			if err := rrs.repo.DeleteRiskControls(deleteControlIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete each compliance requirement if requirements exist
		if len(risk.ComplianceRequirements) > 0 {
			var deleteRequirementIds []*uint
			for _, req := range risk.ComplianceRequirements {
				deleteRequirementIds = append(deleteRequirementIds, req.ID)
			}
			if err := rrs.repo.DeleteRiskComplianceRequirements(deleteRequirementIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Check if controls exist
		if len(risk.Processes) > 0 {
			var deleteProcessIDs []*uint
			for _, process := range risk.Processes {
				deleteProcessIDs = append(deleteProcessIDs, process.ID)
			}
			if err := rrs.repo.DeleteRiskProcesses(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if len(risk.SLARegisters) > 0 {
			var deleteSLARegisterIDs []*uint
			for _, sla := range risk.SLARegisters {
				deleteSLARegisterIDs = append(deleteSLARegisterIDs, sla.ID)
			}
			if err := rrs.repo.DeleteRiskSLAs(deleteSLARegisterIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if len(risk.StrategicObjectives) > 0 {
			var deleteStrategicObjectiveIds []*uint
			for _, obj := range risk.StrategicObjectives {
				deleteStrategicObjectiveIds = append(deleteStrategicObjectiveIds, obj.ID)
			}
			if err := rrs.repo.DeleteRiskStrategicObjectives(deleteStrategicObjectiveIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(risk.AssociatedDocuments) > 0 {
			var deleteDocumentIds []*uint
			for _, document := range risk.AssociatedDocuments {
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}
			if err := rrs.repo.DeleteRiskDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(risk.Assets) > 0 {
			var deleteAssetIds []*uint
			for _, ast := range risk.Assets {
				deleteAssetIds = append(deleteAssetIds, ast.ID)
			}
			if err := rrs.repo.DeleteRiskAssets(deleteAssetIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete timeline
		// todo
		if err := rrs.riskTimelineSvc.Delete(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
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

			// Send DB and email notifications for each successfully deleted risk register
			for _, notification := range notifications {
				urlParam := ""
				entity := "risk-register"
				action := "delete"

				// Prepare DB notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Risk Register Deleted",
					HeadingAr:            "تم حذف سجل المخاطر",
					TitleEn:              notification.TitleEn,
					ReferenceNumber:      notification.ReferenceNumber,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Fetch owner details (responsible user)
				ownerUser, err := rrs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Risk Register/%s: Deleted", notification.TitleEn)
				templateRelativePath := "risk/templates/email/risk_register_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":         ownerUser.FirstNameEn,
					"LastNameEn":          ownerUser.LastNameEn,
					"RiskReferenceNumber": notification.ReferenceNumber,
					"RiskTitle":           notification.TitleEn,
					"RiskDescription":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = rrs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "risk")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}()
	}

	return failedIds, nil
}

func (rrs *RiskRegisterService) UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO) error {
	tx := rrs.repo.BeginTransaction()

	existingRisk, err := rrs.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRisk.RiskResponseStrategyID = &riskResponseStrategyFormDTO.RiskResponseStrategyID

	if err := rrs.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}
	// update timeline
	// todo
	timelineType, err := rrs.riskTimelineSvc.GetByType("risk_response")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := rrs.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := rrs.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rrs *RiskRegisterService) DeleteResponseStrategy(riskID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	existingRisk, err := rrs.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRisk.RiskResponseStrategyID = nil

	if err := rrs.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}
	// update timeline
	// todo
	timelineType, err := rrs.riskTimelineSvc.GetByType("risk_response")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := rrs.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    nil,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      nil,
	}
	if err := rrs.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rrs *RiskRegisterService) GetControlPerformance(riskID uint) (map[string]float64, error) {
	_, err := rrs.GetByID(riskID)
	if err != nil {
		return nil, err
	}

	inherentScore, err := rrs.getControlScoreByType(riskID, "inherent")
	if err != nil {
		return nil, err
	}

	residualScore, err := rrs.getControlScoreByType(riskID, "residual")
	if err != nil {
		return nil, err
	}

	scores := map[string]float64{
		"inherent_control_score": inherentScore,
		"residual_control_score": residualScore,
	}

	return scores, nil
}

func (rrs *RiskRegisterService) getControlScoreByType(riskID uint, controlType string) (float64, error) {
	riskControlAnalyses, err := rrs.repo.GetControlSummary(riskID, controlType)
	if err != nil {
		return 0, err
	}

	if len(riskControlAnalyses) == 0 {
		return 0, nil
	}

	implementationStatusCount := make(map[string]int)
	totalControls := len(riskControlAnalyses)

	for _, rca := range riskControlAnalyses {
		implStatus := rca.ControlImplementationStatus.Type
		implementationStatusCount[implStatus]++
	}

	controlScore := (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100

	return controlScore, nil
}

func (rrs *RiskRegisterService) UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction, evaluationType string) error {
	existingRisk, err := rrs.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if evaluationType == "inherent" {
		existingRisk.InherentImpactID = &riskEvaluationDetails.RiskMatrixImpactID
		existingRisk.InherentLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
		existingRisk.InherentRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
		existingRisk.InherentRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID
	} else if evaluationType == "residual" {
		existingRisk.ResidualImpactID = &riskEvaluationDetails.RiskMatrixImpactID
		existingRisk.ResidualLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
		existingRisk.ResidualRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
		existingRisk.ResidualRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID
	}

	existingRisk.CurrentImpactID = &riskEvaluationDetails.RiskMatrixImpactID
	existingRisk.CurrentLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
	existingRisk.CurrentRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
	existingRisk.CurrentRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID

	if err := rrs.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateProcessRiskRating(existingRisk.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateControlRiskRating(existingRisk.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rrs *RiskRegisterService) UpdateProcessRiskRating(riskID uint, tx utility.DatabaseTransaction) error {
	// Get the risk register details for the given riskID
	risk, err := rrs.repo.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Iterate through each associated process and update their risk rating
	for _, riskProcess := range risk.Processes {
		// Get the full process details (including associated risks)
		process, err := rrs.processSvc.GetFullDetailsByID(riskProcess.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Find the maximum risk rating value from all associated risks
		var maxRiskRatingValue *models.RiskMatrixRiskRatingValue
		for _, processRisk := range process.Risks {

			riskDetail, err := rrs.repo.GetFullDetailsByID(processRisk.RiskRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if maxRiskRatingValue == nil || (riskDetail.CurrentRiskRatingValue != nil && riskDetail.CurrentRiskRatingValue.Value > maxRiskRatingValue.Value) {
				maxRiskRatingValue = riskDetail.CurrentRiskRatingValue
			}
		}

		existingProcess, err := rrs.processSvc.GetByID(riskProcess.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// If we found a max risk rating value, update the process with the corresponding rating
		if maxRiskRatingValue != nil {
			// Update the process risk rating values
			existingProcess.RiskMatrixRiskRatingID = &maxRiskRatingValue.RiskMatrixRiskRatingID
			existingProcess.RiskMatrixRiskRatingValueID = &maxRiskRatingValue.ID

			// Update the process risk rating in the database
			err = rrs.processSvc.UpdateProcessRiskRating(existingProcess, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (rrs *RiskRegisterService) UpdateControlRiskRating(riskID uint, tx utility.DatabaseTransaction) error {
	// Get the risk register details for the given riskID
	risk, err := rrs.repo.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Iterate through each associated control and update their risk rating
	for _, riskControl := range risk.Controls {
		// Get the full control details (including associated risks)
		control, err := rrs.controlSvc.GetFullDetailsByID(riskControl.ControlID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Find the maximum risk rating value from all associated risks
		var maxRiskRatingValue *models.RiskMatrixRiskRatingValue
		for _, controlRisk := range control.Risks {

			riskDetail, err := rrs.repo.GetFullDetailsByID(controlRisk.RiskRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if maxRiskRatingValue == nil || (riskDetail.CurrentRiskRatingValue != nil && riskDetail.CurrentRiskRatingValue.Value > maxRiskRatingValue.Value) {
				maxRiskRatingValue = riskDetail.CurrentRiskRatingValue
			}
		}

		existingControl, err := rrs.controlSvc.GetByID(riskControl.ControlID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// If we found a max risk rating value, update the control with the corresponding rating
		if maxRiskRatingValue != nil {
			// Update the control risk rating values
			existingControl.RiskMatrixRiskRatingID = &maxRiskRatingValue.RiskMatrixRiskRatingID
			existingControl.RiskMatrixRiskRatingValueID = &maxRiskRatingValue.ID

			// Update the control risk rating in the database
			err = rrs.controlSvc.UpdateControlRiskRating(existingControl, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (rrs *RiskRegisterService) ClearRiskRating(riskID uint, updateDTO *riskDTOs.RiskRegisterRiskRatingUpdateDTO, tx utility.DatabaseTransaction) error {
	riskRegister, err := rrs.GetByID(riskID)
	if err != nil {
		return err
	}

	riskRegister.InherentImpactID = updateDTO.InherentImpactID
	riskRegister.InherentLikelihoodID = updateDTO.InherentLikelihoodID
	riskRegister.InherentRiskRatingID = updateDTO.InherentRiskRatingID
	riskRegister.InherentRiskRatingValueID = updateDTO.InherentRiskRatingValueID

	riskRegister.CurrentImpactID = updateDTO.CurrentImpactID
	riskRegister.CurrentLikelihoodID = updateDTO.CurrentLikelihoodID
	riskRegister.CurrentRiskRatingID = updateDTO.CurrentRiskRatingID
	riskRegister.CurrentRiskRatingValueID = updateDTO.CurrentRiskRatingValueID

	riskRegister.ResidualImpactID = updateDTO.ResidualImpactID
	riskRegister.ResidualLikelihoodID = updateDTO.ResidualLikelihoodID
	riskRegister.ResidualRiskRatingID = updateDTO.ResidualRiskRatingID
	riskRegister.ResidualRiskRatingValueID = updateDTO.ResidualRiskRatingValueID

	if err := rrs.repo.Update(riskRegister, tx); err != nil {
		return err
	}

	if err := rrs.UpdateProcessRiskRating(riskRegister.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateControlRiskRating(riskRegister.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rrs *RiskRegisterService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	if _, err := rrs.riskCategorySvc.GetByID(riskRegisterFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}
	if riskRegisterFormDTO.RiskSourceID != nil {
		if _, err := rrs.riskSourceSvc.GetByID(*riskRegisterFormDTO.RiskSourceID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk source does not exist"), "not_exist", "risk source")
		}
	}

	if _, err := rrs.userSvc.GetByID(riskRegisterFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if riskRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := rrs.reviewFrequencySvc.GetByID(*riskRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	controlIds := riskRegisterFormDTO.ControlIDs

	controlIDs := make(map[uint]bool)

	for _, controlId := range controlIds {
		if controlIDs[*controlId] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", controlId), "duplicate", "control")
		}

		if _, err := rrs.controlSvc.GetByID(*controlId); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlId), "not_exist", "control")
		}

		controlIDs[*controlId] = true
	}

	requirementIds := riskRegisterFormDTO.ComplianceRequirementIDs

	requirementIDs := make(map[uint]bool)

	for _, requirementId := range requirementIds {
		if requirementIDs[*requirementId] {
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", requirementId), "duplicate", "compliance requirement")
		}

		if _, err := rrs.complianceRequirementSvc.GetByComplianceRequirementID(*requirementId); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", requirementId), "not_exist", "compliance requirement")
		}

		requirementIDs[*requirementId] = true
	}

	processIds := riskRegisterFormDTO.ProcessIDs

	processIDs := make(map[uint]bool)

	for _, processId := range processIds {
		if processIDs[*processId] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", processId), "duplicate", "process")
		}

		if _, err := rrs.processSvc.GetByID(*processId); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processId), "not_exist", "process")
		}

		processIDs[*processId] = true
	}

	slaRegisterIds := riskRegisterFormDTO.SLARegisterIDs

	slaRegisterIDs := make(map[uint]bool)

	for _, slaRegisterId := range slaRegisterIds {
		if slaRegisterIDs[*slaRegisterId] {
			return utility.NewValidationError(fmt.Errorf("duplicate slaRegister ID: %d", slaRegisterId), "duplicate", "slaRegister")
		}

		if _, err := rrs.slaRegisterSvc.GetByID(*slaRegisterId); err != nil {
			return utility.NewValidationError(fmt.Errorf("slaRegister does not exist for ID: %d", slaRegisterId), "not_exist", "slaRegister")
		}

		slaRegisterIDs[*slaRegisterId] = true
	}

	strategicObjectiveIds := riskRegisterFormDTO.StrategicObjectiveIDs

	strategicObjectiveIDs := make(map[uint]bool)

	for _, strategicObjectiveId := range strategicObjectiveIds {
		if strategicObjectiveIDs[*strategicObjectiveId] {
			return utility.NewValidationError(fmt.Errorf("duplicate strategic objective ID: %d", strategicObjectiveId), "duplicate", "strategic objective")
		}

		if _, err := rrs.strategicObjectiveSvc.GetByID(*strategicObjectiveId); err != nil {
			return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveId), "not_exist", "strategic objective")
		}

		strategicObjectiveIDs[*strategicObjectiveId] = true
	}

	assets := riskRegisterFormDTO.AssetIds
	assetIDs := make(map[uint]bool)

	for _, asset := range assets {
		if assetIDs[*asset] {
			return utility.NewValidationError(fmt.Errorf("duplicate asset ID: %d", asset), "duplicate", "asset")
		}

		if _, err := rrs.assetRegisterSvc.GetByID(*asset); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", asset), "not_exist", "asset")
		}

		// Add the asset ID to the map
		assetIDs[*asset] = true
	}

	return nil
}

func (rrs *RiskRegisterService) parseRiskRegisterDates(dto *riskDTOs.RiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterParsedDates{}

	var err error

	if dto.ReportedAt != nil && *dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(*dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported at"), "invalid_date", "reported at")
		}
	}

	if dto.NextReviewDate != nil && *dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}

func (rrs *RiskRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	// Get full details of the risk register by ID
	existingRiskRegister, err := rrs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	// Check if there are any associated documents
	if len(existingRiskRegister.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk register")
	}

	// Decrypt the token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Define the storage path
	storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(id), 10), "documents")

	// Get the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the specified path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (cs *RiskRegisterService) ExportRiskRegister(data []*riskDTOs.RiskRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := cs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Ref Number",
			"Risk Title",
			"Description",
			"Risk Category",
			"Responsible Unit (TBD)",
			"Responsible User",
			"Risk Source",
			"Risk Status",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Risk Response",
			"Residual Likelihood",
			"Residual Impact",
			"Residual Risk Rating",
			"Reported At",
			"Review Frequency",
			"Next Review Date",
		}
	} else {
		headers = []string{
			"SL No.",
			"Ref Number",
			"Risk Title",
			"Description",
			"Risk Category",
			"Responsible Unit (TBD)",
			"Responsible User",
			"Risk Source",
			"Risk Status",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Risk Response",
			"Residual Likelihood",
			"Residual Impact",
			"Residual Risk Rating",
			"Reported At",
			"Review Frequency",
			"Next Review Date",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, risk := range data {
		// Initialize data variables with default values
		var riskTitle, riskCategory, responsibleUnit, responsibleUser, riskSource, riskStatus string
		var riskResponse, residualLikelihood, residualImpact, residualRiskRating string
		var description, likelihood, impact, riskRating string
		var reviewFrequency, nextReviewDate, reportedAt string

		// Assign values based on language preference and check for nil pointers
		if userSetting.Language == "ar" {
			riskCategory = risk.RiskCategoryTitleAr
			responsibleUnit = risk.ResponsibleUnitTitleAr
			responsibleUser = risk.ResponsibleUserFirstNameAr + " " + risk.ResponsibleUserLastNameAr
			riskSource = utility.GetStringValue(risk.RiskSourceTitleAr)
			riskStatus = risk.RiskStatusTitleAr
			riskResponse = risk.RiskResponseStrategyTitleAr
			likelihood = utility.GetStringValue(risk.InherentLikelihoodTitleAr)
			impact = utility.GetStringValue(risk.InherentImpactTitleAr)
			riskRating = utility.GetStringValue(risk.InherentRiskRatingTitleAr)
			residualLikelihood = utility.GetStringValue(risk.ResidualLikelihoodTitleAr)
			residualImpact = utility.GetStringValue(risk.ResidualImpactTitleAr)
			residualRiskRating = utility.GetStringValue(risk.ResidualRiskRatingTitleAr)
			reviewFrequency = utility.GetStringValue(risk.ReviewFrequencyTitleAr)
		} else {
			riskCategory = risk.RiskCategoryTitleEn
			responsibleUnit = risk.ResponsibleUnitTitleEn
			responsibleUser = risk.ResponsibleUserFirstNameEn + " " + risk.ResponsibleUserLastNameEn
			riskSource = utility.GetStringValue(risk.RiskSourceTitleEn)
			riskStatus = risk.RiskStatusTitleEn
			riskResponse = risk.RiskResponseStrategyTitleEn
			likelihood = utility.GetStringValue(risk.InherentLikelihoodTitleEn)
			impact = utility.GetStringValue(risk.InherentImpactTitleEn)
			riskRating = utility.GetStringValue(risk.InherentRiskRatingTitleEn)
			residualLikelihood = utility.GetStringValue(risk.ResidualLikelihoodTitleEn)
			residualImpact = utility.GetStringValue(risk.ResidualImpactTitleEn)
			residualRiskRating = utility.GetStringValue(risk.ResidualRiskRatingTitleEn)
			reviewFrequency = utility.GetStringValue(risk.ReviewFrequencyTitleEn)
		}

		riskTitle = risk.Title
		description = utility.GetStringValue(risk.Description)
		nextReviewDate = utility.GetStringValue(risk.NextReviewDate)
		reportedAt = utility.GetStringValue(risk.ReportedAt)

		rows[i] = []interface{}{
			i + 1,
			risk.ReferenceNumber,
			riskTitle,
			description,
			riskCategory,
			responsibleUnit,
			responsibleUser,
			riskSource,
			riskStatus,
			likelihood,
			impact,
			riskRating,
			riskResponse,
			residualLikelihood,
			residualImpact,
			residualRiskRating,
			reportedAt,
			reviewFrequency,
			nextReviewDate,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

// func getStringValue(pointer *string) string {
// 	if pointer != nil {
// 		return *pointer
// 	}
// 	return ""
// }
