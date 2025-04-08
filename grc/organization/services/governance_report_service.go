package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"
	"strings"

	"gorm.io/gorm"
)

type GovernanceReportServiceInterface interface {
	GetAll(params *dtos.GovernanceReportListQueryParamsDTO) (*dtos.GovernanceReportListResponseDTO, error)
	GetGovernanceReportDetails(governanceReportID uint) (*dtos.GovernanceReportDetailDTO, error)
	Create(currentUserID uint, governanceReportFormDTO *dtos.GovernanceReportFormDTO) (uint, error)
	GetSectionByID(governanceReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.GovernanceReportSection, error)
	UpdateSection(governanceID uint, id uint, currentUserID uint, governanceReporSectionFormDTO *dtos.GovernanceReportSectionFormDTO) error
	DeleteGovernanceReportSection(governanceReportID uint, sectionID uint, currentUserID uint) error
	DeleteGovernanceReport(governanceReportID uint, currentUserID uint) error
}

type GovernanceReportService struct {
	repo                 repositories.GovernanceReportRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
}

func NewGovernanceReportService(
	repo repositories.GovernanceReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
) GovernanceReportServiceInterface {
	return &GovernanceReportService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (grs *GovernanceReportService) GetAll(params *dtos.GovernanceReportListQueryParamsDTO) (*dtos.GovernanceReportListResponseDTO, error) {
	return grs.repo.GetAll(params)
}

func (grs *GovernanceReportService) GetGovernanceReportDetails(governanceReportID uint) (*dtos.GovernanceReportDetailDTO, error) {
	governanceReport, err := grs.repo.GetGovernanceReportByID(governanceReportID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, governanceReportID)
		}
		return nil, err
	}

	sections, err := grs.repo.GetGovernanceReportSections(governanceReportID)
	if err != nil {
		return nil, err
	}

	documents, err := grs.repo.GetGovernanceReportDocuments(governanceReportID)
	if err != nil {
		return nil, err
	}

	meetings, err := grs.repo.GetGovernanceReportMeetings(governanceReportID)
	if err != nil {
		return nil, err
	}

	tasks, err := grs.repo.GetGovernanceReportTasks(governanceReportID)
	if err != nil {
		return nil, err
	}

	// Construct Response DTO
	response := &dtos.GovernanceReportDetailDTO{
		ID:                     governanceReport.ID,
		Title:                  governanceReport.Title,
		StartDate:              governanceReport.StartDate,
		EndDate:                governanceReport.EndDate,
		CreatedAt:              governanceReport.CreatedAt,
		CreatedUserFirstNameEn: governanceReport.CreatedUserFirstNameEn,
		CreatedUserFirstNameAr: governanceReport.CreatedUserFirstNameAr,
		CreatedUserLastNameEn:  governanceReport.CreatedUserLastNameEn,
		CreatedUserLastNameAr:  governanceReport.CreatedUserLastNameAr,
		Sections:               sections,
		Documents:              documents,
		Meetings:               meetings,
		Tasks:                  tasks,
	}

	return response, nil
}

func (grs *GovernanceReportService) GetSectionByID(governanceReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.GovernanceReportSection, error) {
	section, err := grs.repo.GetSectionByID(sectionID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}
	return section, nil
}

func (grs *GovernanceReportService) Create(currentUserID uint, governanceReportFormDTO *dtos.GovernanceReportFormDTO) (uint, error) {
	tx := grs.repo.BeginTransaction()

	dates, err := grs.parseGovernanceReportDates(governanceReportFormDTO)
	if err != nil {
		return 0, err
	}

	governanceReport := &models.GovernanceReport{
		Title:       governanceReportFormDTO.Title,
		StartDate:   *dates.StartDate,
		EndDate:     *dates.EndDate,
		CreatedBYID: currentUserID,
	}

	governanceReportID, err := grs.repo.Create(governanceReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"executive-summary", "policy-acknowledgement", "meetings", "actions", "conclusion",
	}

	for i, sectionType := range sectionTypes {
		titleEn, titleAr := formatTitle(sectionType)
		section := models.GovernanceReportSection{
			GovernanceReportID: governanceReportID,
			Type:               sectionType,
			TitleEn:            titleEn,
			TitleAr:            titleAr,
			DescriptionEn:      "",
			DescriptionAr:      "",
			OrderNum:           i + 1,
		}
		if err := grs.repo.CreateSection(&section, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	documents, err := grs.repo.GetPolicyDocumentByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, document := range documents {
		summary, err := grs.GetDocumentAcknowledgmentSummary(document)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get acknowledgment summary for document ID %d: %w", document.ID, err)
		}

		governanceDocument := &models.GovernanceReportDocument{
			GovernanceReportID:     governanceReportID,
			ReferenceNumber:        document.ReferenceNumber,
			DocumentTitleEn:        document.TitleEn,
			DocumentTitleAr:        document.TitleAr,
			ResponsibleUserNameEn:  fmt.Sprintf("%s %s", document.ResponsibleUser.FirstNameEn, document.ResponsibleUser.LastNameEn),
			ResponsibleUserNameAr:  fmt.Sprintf("%s %s", document.ResponsibleUser.FirstNameAr, document.ResponsibleUser.LastNameAr),
			ResponsibleUnitTitleEn: document.ResponsibleUnit.TitleEn,
			ResponsibleUnitTitleAr: document.ResponsibleUnit.TitleAr,
			TotalEmployees:         int(summary.TotalUsers),
			AcknowledgedCount:      int(summary.AcknowledgedUsers),
			PendingCount:           int(summary.NotAcknowledgedUsers),
			CreatedByID:            currentUserID,
		}

		if err := grs.repo.CreateGovernanceIDocument(governanceDocument, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Step: Fetch Meetings in the date range
	meetings, err := grs.repo.GetMeetingsByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch meetings: %w", err)
	}

	// Step: For each meeting, gather required data and insert
	for _, meeting := range meetings {
		participantCount := len(meeting.Participants)

		var keyDecisionsList []string
		for _, minute := range meeting.MeetingMinutes {
			if minute.Minutes != nil && strings.TrimSpace(*minute.Minutes) != "" {
				keyDecisionsList = append(keyDecisionsList, strings.TrimSpace(*minute.Minutes))
			}
		}
		keyDecisions := strings.Join(keyDecisionsList, ", ")

		var assignedActions []string
		for _, actionplans := range meeting.MeetingActionPlans {
			assignedActions = append(assignedActions, actionplans.Title)
		}

		meetingRecord := &models.GovernanceReportMeeting{
			GovernanceReportID: governanceReportID,
			MeetingTitle:       meeting.Title,
			MeetingDate:        meeting.Date,
			ParticipantsCount:  participantCount,
			KeyDecisions:       keyDecisions,
			ActionsAssigned:    strings.Join(assignedActions, ", "),
			CreatedByID:        currentUserID,
		}

		if err := grs.repo.CreateGovernanceReportMeeting(meetingRecord, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert governance_report_meeting: %w", err)
		}
	}

	tasks, err := grs.repo.GetTasksByDateRange(*dates.StartDate, *dates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch governance tasks: %w", err)
	}

	for _, task := range tasks {
		taskRecord := &models.GovernanceReportTask{
			GovernanceReportID:    governanceReportID,
			TaskTitle:             task.Title,
			TaskDescription:       task.Description,
			ResponsibleUserNameEn: fmt.Sprintf("%s %s", task.ResponsibleUser.FirstNameEn, task.ResponsibleUser.LastNameEn),
			ResponsibleUserNameAr: fmt.Sprintf("%s %s", task.ResponsibleUser.FirstNameAr, task.ResponsibleUser.LastNameAr),
			TaskStatusTitleEn:     task.TaskStatus.TitleEn,
			TaskStatusTitleAr:     task.TaskStatus.TitleAr,
			TaskStatusColorCode:   task.TaskStatus.ColorCode,
			TargetDate:            *task.TargetDate,
			CreatedByID:           currentUserID,
		}

		if err := grs.repo.CreateGovernanceReportTask(taskRecord, tx); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert governance_report_task: %w", err)
		}
	}

	tx.Commit()

	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.GovernanceReport{}.TableName(),
			governanceReportID,
			governanceReport,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return governanceReportID, nil
}

func (grs *GovernanceReportService) UpdateSection(governanceReportID uint, sectionID uint, currentUserID uint, governanceReporSectionFormDTO *dtos.GovernanceReportSectionFormDTO) error {
	tx := grs.repo.BeginTransaction()

	existingReportSection, err := grs.GetSectionByID(governanceReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Capture old values before update
	oldValues := map[string]interface{}{
		"title_en":       existingReportSection.TitleEn,
		"title_ar":       existingReportSection.TitleAr,
		"description_en": existingReportSection.DescriptionEn,
		"description_ar": existingReportSection.DescriptionAr,
	}

	// Update section details
	existingReportSection.TitleEn = governanceReporSectionFormDTO.TitleEn
	existingReportSection.TitleAr = governanceReporSectionFormDTO.TitleAr
	existingReportSection.DescriptionEn = *governanceReporSectionFormDTO.DescriptionEn
	existingReportSection.DescriptionAr = *governanceReporSectionFormDTO.DescriptionAr

	err = grs.repo.UpdateSection(existingReportSection, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Capture new values after update
	newValues := map[string]interface{}{
		"title_en":       existingReportSection.TitleEn,
		"title_ar":       existingReportSection.TitleAr,
		"description_en": existingReportSection.DescriptionEn,
		"description_ar": existingReportSection.DescriptionAr,
	}

	// Identify changes
	changedData := utility.CompareFieldChanges(oldValues, newValues)

	// Log activity only if there are changes
	if len(changedData) > 0 {
		go func() {
			err := grs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.GovernanceReportSection{}.TableName(),
				sectionID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return nil
}

func (grs *GovernanceReportService) DeleteGovernanceReportSection(governanceReportID uint, sectionID uint, currentUserID uint) error {
	tx := grs.repo.BeginTransaction()

	section, err := grs.GetSectionByID(governanceReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the section
	if err := grs.repo.DeleteGovernanceReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reorder remaining sections
	if err := grs.repo.ReorderGovernanceReportSections(governanceReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.GovernanceReportSection{}.TableName(),
			sectionID,
			section,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (grs *GovernanceReportService) DeleteGovernanceReport(governanceReportID uint, currentUserID uint) error {
	tx := grs.repo.BeginTransaction()

	existingReport, err := grs.repo.GetGovernanceReportByID(governanceReportID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGovernanceReportTasks(governanceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGovernanceReportMeetings(governanceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGovernanceReportDocuments(governanceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGovernanceReportSections(governanceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := grs.repo.DeleteGovernanceReport(governanceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		err := grs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.GovernanceReport{}.TableName(),
			governanceReportID,
			existingReport,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	tx.Commit()
	return nil
}

func formatTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "policy-acknowledgement":
		return "Policy Acknowledgment Status", "حالة إقرار السياسات"
	case "meetings":
		return "Governance Meetings & Decisions", "اجتماعات الحوكمة والقرارات"
	case "actions":
		return "Governance Tasks & Action Items", "مهام الحوكمة وعناصر العمل"
	case "conclusion":
		return "Conclusion", "الخاتمة"
	default:
		return sectionType, sectionType // fallback to sectionType if no match
	}
}

func (grs *GovernanceReportService) parseGovernanceReportDates(dto *dtos.GovernanceReportFormDTO) (*dtos.GovernanceReportParsedDates, error) {
	parsed := &dtos.GovernanceReportParsedDates{}

	var err error

	parsed.StartDate, err = utility.ParseDate(dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
	}

	parsed.EndDate, err = utility.ParseDate(dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
	}

	return parsed, nil
}

func (grs *GovernanceReportService) GetDocumentAcknowledgmentSummary(document *models.Document) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {

	if document.DocumentAccessType.Type == "public" {
		return grs.repo.GetPublicDocumentAcknowledgmentSummary(document.ID)
	}

	primaryLevel, err := grs.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	var primaryOrgUnitField string
	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return grs.repo.GetDocumentAcknowledgmentSummary(document.ID, primaryOrgUnitField)
}
