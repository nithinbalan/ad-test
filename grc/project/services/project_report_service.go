package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectReportServiceInterface interface {
	ExportProjectReport(id uint, currentUserID uint) (string, error)
	GetFullReportDetailsByID(id uint) (*dtos.ProjectReportDetailDTO, error)

	Create(currentUserID uint, dto dtos.ProjectReportRequestDTO) (uint, error)
	GetAll(params *dtos.ProjectReportListQueryParamsDTO) (*dtos.ProjectReportListResponseDTO, error)
	GetProjectReportDetails(id uint) (*dtos.ProjectReportDetailsDTO, error)
	DeleteProjectReport(id uint, currentUserID uint) error
	GetProjectReportSection(projectReportID uint, sectionID uint) (*dtos.ProjectReportSectionDTO, error)
	UpdateProjectReportSection(projectReportID uint, sectionID uint, form dtos.ProjectReportSectionFormDTO, currentUserID uint) error
	DeleteProjectReportSection(projectReportID uint, sectionID uint, currentUserID uint) error
	GetProjectReportExportZip(projectReportID uint, currentUserID uint) (string, error)
}

type ProjectReportService struct {
	repo      repositories.ProjectReportRepositoryInterface
	settigSvc settingSvc.ReportTemplateSettingServiceInterface
}

func NewProjectReportService(
	repo repositories.ProjectReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,
) ProjectReportServiceInterface {
	return &ProjectReportService{
		repo:      repo,
		settigSvc: settigSvc,
	}
}

func (prs *ProjectReportService) GetFullReportDetailsByID(id uint) (*dtos.ProjectReportDetailDTO, error) {
	project, err := prs.repo.GetFullReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if _, err := prs.settigSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := prs.settigSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover_image_missing")
	}

	var projectDTO dtos.ProjectReportDetailDTO
	err = mapstructure.Decode(project, &projectDTO)
	if err != nil {
		return nil, err
	}

	projectDTO.StartDateString = project.StartDate.Format("2006-01-02")
	projectDTO.EndDateString = project.EndDate.Format("2006-01-02")
	projectDTO.StartDateFormatted = project.StartDate.Format("01/02/2006")
	projectDTO.EndDateFormatted = project.EndDate.Format("01/02/2006")
	projectDTO.CreatedAt = project.CreatedAt
	projectDTO.UpdatedAt = project.UpdatedAt

	if project.IncidentActionPlanID != nil {
		projectDTO.IncidentActionPlanID = project.IncidentActionPlanID
	}
	if project.IncidentActionPlan != nil {
		actionPlanDTO := &dtos.ProjectIncidentActionPlanDTO{
			ID:                 project.IncidentActionPlan.ID,
			IncidentRegisterID: project.IncidentActionPlan.IncidentRegisterID,
			Title:              project.IncidentActionPlan.Title,
			Description:        project.IncidentActionPlan.Description,
			Percentage:         project.IncidentActionPlan.Percentage,
			ActionPlanStatusID: project.IncidentActionPlan.ActionPlanStatusID,
			ResponsibleUserID:  project.IncidentActionPlan.ResponsibleUserID,
		}
		if project.IncidentActionPlan.IncidentRegister != nil {
			actionPlanDTO.IncidentReferenceNumber = &project.IncidentActionPlan.IncidentRegister.ReferenceNumber
			actionPlanDTO.IncidentTitle = &project.IncidentActionPlan.IncidentRegister.Title
		}
		if project.IncidentActionPlan.TargetDate != nil {
			targetDateStr := project.IncidentActionPlan.TargetDate.Format("2006-01-02")
			actionPlanDTO.TargetDate = &targetDateStr
		}
		if project.IncidentActionPlan.ActionPlanStatus != nil {
			actionPlanDTO.ActionPlanStatusTitleEn = &project.IncidentActionPlan.ActionPlanStatus.TitleEn
			actionPlanDTO.ActionPlanStatusTitleAr = &project.IncidentActionPlan.ActionPlanStatus.TitleAr
			actionPlanDTO.ActionPlanStatusColor = &project.IncidentActionPlan.ActionPlanStatus.ColorCode
		}
		if project.IncidentActionPlan.ResponsibleUser != nil {
			actionPlanDTO.ResponsibleUser = &organizationDTOs.UserMinimalDTO{
				ID:               int(project.IncidentActionPlan.ResponsibleUser.ID),
				FirstNameEN:      project.IncidentActionPlan.ResponsibleUser.FirstNameEn,
				LastNameEN:       project.IncidentActionPlan.ResponsibleUser.LastNameEn,
				FirstNameAR:      project.IncidentActionPlan.ResponsibleUser.FirstNameAr,
				LastNameAR:       project.IncidentActionPlan.ResponsibleUser.LastNameAr,
				ProfileImagePath: project.IncidentActionPlan.ResponsibleUser.ProfileImagePath,
				Email:            project.IncidentActionPlan.ResponsibleUser.Email,
				Initial:          project.IncidentActionPlan.ResponsibleUser.Initial,
				Color:            project.IncidentActionPlan.ResponsibleUser.Color,
			}
		}
		projectDTO.IncidentActionPlan = actionPlanDTO
	}

	milestones, err := prs.repo.GetAllMilestones(id)
	if err == nil {
		projectDTO.Milestones = milestones
	}

	controls, err := prs.repo.GetAllControls(id)
	if err == nil {
		projectDTO.Controls = controls
	}

	risks, err := prs.repo.GetAllRisks(id)
	if err == nil {
		projectDTO.Risks = risks
	}

	processes, err := prs.repo.GetAllProcesses(id)
	if err == nil {
		projectDTO.Processes = processes
	}

	objectives, err := prs.repo.GetAllStrategicObjectives(id)
	if err == nil {
		projectDTO.StrategicObjectives = objectives
	}

	kpis, err := prs.repo.GetAllKPIs(id)
	if err == nil {
		projectDTO.KPIs = kpis
	}

	docs, err := prs.repo.GetAllAssociatedDocuments(id)
	if err == nil {
		projectDTO.AssociatedDocuments = docs
	}

	assets, err := prs.repo.GetAllAssets(id)
	if err == nil {
		projectDTO.Assets = assets
	}

	incidents, err := prs.repo.GetAllIncidents(id)
	if err == nil {
		projectDTO.Incidents = incidents
	}

	complianceRequirements, err := prs.repo.GetAllComplianceRequirements(id)
	if err == nil {
		projectDTO.ComplianceRequirements = complianceRequirements
	}

	actionPlans, err := prs.repo.GetAllActionPlans(id)
	if err == nil {
		projectDTO.ActionPlans = actionPlans
	}

	changeRequests, err := prs.repo.GetAllChangeRequests(id)
	if err == nil {
		projectDTO.ChangeRequests = changeRequests
	}

	decisionLogs, err := prs.repo.GetAllDecisionLogs(id)
	if err == nil {
		projectDTO.DecisionLogs = decisionLogs
	}

	supportingDocs, err := prs.repo.GetAllSupportingDocuments(id)
	if err == nil {
		projectDTO.SupportingDocuments = supportingDocs
	}

	closure, err := prs.repo.GetProjectClosure(id)
	if err != nil {
		return nil, err
	}
	if closure != nil {
		projectDTO.Closure = closure
	}

	return &projectDTO, nil
}

func (prs *ProjectReportService) ExportProjectReport(id uint, currentUserID uint) (string, error) {
	projectDTO, err := prs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}

	projectDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	// Main content template
	templatePaths := []string{"project/templates/pdf/project_report.html"}
	htmlContent, err := utility.GeneratePDFContent(projectDTO, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template
	headerTemplatePaths := []string{"project/templates/pdf/project_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(projectDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template
	footerTemplatePaths := []string{"project/templates/pdf/project_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(projectDTO, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// Cover template
	coverTemplatePaths := []string{"project/templates/pdf/project_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(projectDTO, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// Convert HTML -> PDF with cover, header, footer
	pdfFileName := fmt.Sprintf("ProjectReport_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

type ProjectReportParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}

func (prs *ProjectReportService) parseProjectReportDates(dto dtos.ProjectReportRequestDTO) (*ProjectReportParsedDates, error) {
	var startDate, endDate *time.Time

	if dto.StartDate != nil && *dto.StartDate != "" {
		t, err := time.Parse("2006-01-02", *dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid start_date format, expected YYYY-MM-DD"), "invalid", "start_date")
		}
		startDate = &t
	}

	if dto.EndDate != nil && *dto.EndDate != "" {
		t, err := time.Parse("2006-01-02", *dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid end_date format, expected YYYY-MM-DD"), "invalid", "end_date")
		}
		endDate = &t
	}

	return &ProjectReportParsedDates{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

func formatProjectSectionTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "scope":
		return "Scope", "النطاق"
	case "objectives":
		return "Objectives", "الأهداف"
	case "overview":
		return "Overview", "نظرة عامة"
	case "details":
		return "Details", "التفاصيل"
	case "closure":
		return "Project Closure", "إغلاق المشروع"
	default:
		return sectionType, sectionType
	}
}

func (prs *ProjectReportService) Create(currentUserID uint, dto dtos.ProjectReportRequestDTO) (uint, error) {
	tx := prs.repo.BeginTransaction()

	template, err := prs.settigSvc.GetFirst()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if template.CoverPageBGImagePath == nil || template.LogoImagePath == nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Logo and Background Image not added.")
	}

	dates, err := prs.parseProjectReportDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	report := &models.ProjectReport{
		Title:       dto.Title,
		StartDate:   dates.StartDate,
		EndDate:     dates.EndDate,
		CreatedBYID: currentUserID,
	}

	reportID, err := prs.repo.CreateReport(report, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"executive-summary", "scope", "objectives", "overview", "details", "closure",
	}

	for i, s := range sectionTypes {
		en, ar := formatProjectSectionTitle(s)
		sec := models.ProjectReportSection{
			ProjectReportID: reportID,
			Type:            s,
			TitleEn:         en,
			TitleAr:         ar,
			OrderNum:        i + 1,
		}
		if _, err := prs.repo.CreateSection(&sec, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var projects []*models.Project

	switch dto.Type {
	case "date-range":
		if dates.StartDate == nil || dates.EndDate == nil {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("start_date and end_date are required for date-range reports"),
				"required",
				"start_date",
			)
		}
		if dates.StartDate.After(*dates.EndDate) {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("start_date must be less than or equal to end_date"),
				"invalid",
				"start_date",
			)
		}
		projects, err = prs.repo.GetProjectsByDateRange(*dates.StartDate, *dates.EndDate, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if len(projects) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("no projects found for the specified criteria"),
				"empty_result",
				"date_range",
			)
		}

	case "individual":
		if dto.ProjectIDs == nil || len(*dto.ProjectIDs) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("project_ids are required for individual reports"),
				"required",
				"project_ids",
			)
		}
		projects, err = prs.repo.GetProjectsByIDs(*dto.ProjectIDs, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if len(projects) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(
				errors.New("no projects found for the specified IDs"),
				"empty_result",
				"project_ids",
			)
		}
	}

	for _, proj := range projects {
		var respUserEn, respUserAr string
		if proj.ResponsibleUser != nil {
			respUserEn = fmt.Sprintf("%s %s", proj.ResponsibleUser.FirstNameEn, proj.ResponsibleUser.LastNameEn)
			respUserAr = fmt.Sprintf("%s %s", proj.ResponsibleUser.FirstNameAr, proj.ResponsibleUser.LastNameAr)
		}

		var projTypeEn, projTypeAr *string
		if proj.ProjectType != nil {
			projTypeEn = &proj.ProjectType.TitleEn
			projTypeAr = &proj.ProjectType.TitleAr
		}

		var projStatusEn, projStatusAr *string
		if proj.ProjectStatus != nil {
			projStatusEn = &proj.ProjectStatus.TitleEn
			projStatusAr = &proj.ProjectStatus.TitleAr
		}

		var projPhaseEn, projPhaseAr *string
		if proj.ProjectPhase != nil {
			projPhaseEn = &proj.ProjectPhase.TitleEn
			projPhaseAr = &proj.ProjectPhase.TitleAr
		}

		var projClassEn, projClassAr *string
		if proj.ProjectClassification != nil {
			projClassEn = &proj.ProjectClassification.TitleEn
			projClassAr = &proj.ProjectClassification.TitleAr
		}

		var incidentActionPlanID *uint
		var incidentActionPlanTitle *string
		if proj.IncidentActionPlan != nil {
			incidentActionPlanID = proj.IncidentActionPlanID
			incidentActionPlanTitle = &proj.IncidentActionPlan.Title
		}

		reportProj := &models.ProjectReportProject{
			ProjectReportID:              reportID,
			ProjectID:                    proj.ID,
			TitleEn:                      proj.TitleEn,
			TitleAr:                      proj.TitleAr,
			DescriptionEn:                proj.DescriptionEn,
			DescriptionAr:                proj.DescriptionAr,
			Budget:                       proj.Budget,
			Scope:                        proj.Scope,
			StartDate:                    proj.StartDate,
			EndDate:                      proj.EndDate,
			ProjectTypeTitleEn:           projTypeEn,
			ProjectTypeTitleAr:           projTypeAr,
			ProjectStatusTitleEn:         projStatusEn,
			ProjectStatusTitleAr:         projStatusAr,
			ProjectPhaseTitleEn:          projPhaseEn,
			ProjectPhaseTitleAr:          projPhaseAr,
			ProjectClassificationTitleEn: projClassEn,
			ProjectClassificationTitleAr: projClassAr,
			ResponsibleUserNameEn:        respUserEn,
			ResponsibleUserNameAr:        respUserAr,
			IncidentActionPlanID:         incidentActionPlanID,
			IncidentActionPlanTitle:      incidentActionPlanTitle,
			CreatedAt:                    proj.CreatedAt,
		}

		if err := prs.repo.CreateReportProject(reportProj, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		supportingDocs, _ := prs.repo.GetAllSupportingDocuments(proj.ID)
		if len(supportingDocs) > 0 {
			var reportDocs []*models.ProjectReportDocument
			for _, sd := range supportingDocs {
				for _, doc := range sd.Documents {
					name := doc.Name
					size := doc.Size
					token := doc.Token
					docType := doc.Type
					reportDocs = append(reportDocs, &models.ProjectReportDocument{
						ProjectReportProjectID: reportProj.ID,
						Name:                   &name,
						Size:                   &size,
						Token:                  &token,
						OriginalName:           doc.OriginalName,
						Type:                   &docType,
						CreatedAt:              time.Now(),
					})
				}
			}
			if len(reportDocs) > 0 {
				if err := prs.repo.CreateProjectDocuments(reportDocs, tx); err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	tx.Commit()
	return reportID, nil
}

func (prs *ProjectReportService) GetAll(params *dtos.ProjectReportListQueryParamsDTO) (*dtos.ProjectReportListResponseDTO, error) {
	return prs.repo.GetAllReports(params)
}

func (prs *ProjectReportService) GetProjectReportDetails(id uint) (*dtos.ProjectReportDetailsDTO, error) {
	report, err := prs.repo.GetReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var sectionsDTO []*dtos.ProjectReportSectionDTO
	for _, sec := range report.Sections {
		sectionsDTO = append(sectionsDTO, &dtos.ProjectReportSectionDTO{
			ID:            sec.ID,
			Type:          sec.Type,
			TitleEn:       sec.TitleEn,
			TitleAr:       sec.TitleAr,
			DescriptionEn: sec.DescriptionEn,
			DescriptionAr: sec.DescriptionAr,
			OrderNum:      sec.OrderNum,
		})
	}

	var projectsDTO []*dtos.ProjectReportProjectDetailDTO
	for _, rp := range report.Projects {
		var docsDTO []*dtos.ProjectReportDocumentDTO
		for _, d := range rp.Documents {
			docsDTO = append(docsDTO, &dtos.ProjectReportDocumentDTO{
				ID:           d.ID,
				Name:         d.Name,
				Size:         d.Size,
				Token:        d.Token,
				OriginalName: d.OriginalName,
				Type:         d.Type,
			})
		}

		milestones, _ := prs.repo.GetAllMilestones(rp.ProjectID)
		controls, _ := prs.repo.GetAllControls(rp.ProjectID)
		risks, _ := prs.repo.GetAllRisks(rp.ProjectID)
		processes, _ := prs.repo.GetAllProcesses(rp.ProjectID)
		objectives, _ := prs.repo.GetAllStrategicObjectives(rp.ProjectID)
		kpis, _ := prs.repo.GetAllKPIs(rp.ProjectID)
		docs, _ := prs.repo.GetAllAssociatedDocuments(rp.ProjectID)
		assets, _ := prs.repo.GetAllAssets(rp.ProjectID)
		incidents, _ := prs.repo.GetAllIncidents(rp.ProjectID)
		complianceRequirements, _ := prs.repo.GetAllComplianceRequirements(rp.ProjectID)
		actionPlans, _ := prs.repo.GetAllActionPlans(rp.ProjectID)
		changeRequests, _ := prs.repo.GetAllChangeRequests(rp.ProjectID)
		decisionLogs, _ := prs.repo.GetAllDecisionLogs(rp.ProjectID)
		supportingDocs, _ := prs.repo.GetAllSupportingDocuments(rp.ProjectID)
		closure, _ := prs.repo.GetProjectClosure(rp.ProjectID)

		var incidentActionPlanDTO *dtos.ProjectIncidentActionPlanDTO
		if rp.Project != nil && rp.Project.IncidentActionPlan != nil {
			iap := rp.Project.IncidentActionPlan
			incidentActionPlanDTO = &dtos.ProjectIncidentActionPlanDTO{
				ID:                 iap.ID,
				IncidentRegisterID: iap.IncidentRegisterID,
				Title:              iap.Title,
				Description:        iap.Description,
				Percentage:         iap.Percentage,
				ActionPlanStatusID: iap.ActionPlanStatusID,
				ResponsibleUserID:  iap.ResponsibleUserID,
			}
			if iap.IncidentRegister != nil {
				incidentActionPlanDTO.IncidentReferenceNumber = &iap.IncidentRegister.ReferenceNumber
				incidentActionPlanDTO.IncidentTitle = &iap.IncidentRegister.Title
			}
			if iap.TargetDate != nil {
				targetDateStr := iap.TargetDate.Format("2006-01-02")
				incidentActionPlanDTO.TargetDate = &targetDateStr
			}
			if iap.ActionPlanStatus != nil {
				incidentActionPlanDTO.ActionPlanStatusTitleEn = &iap.ActionPlanStatus.TitleEn
				incidentActionPlanDTO.ActionPlanStatusTitleAr = &iap.ActionPlanStatus.TitleAr
				incidentActionPlanDTO.ActionPlanStatusColor = &iap.ActionPlanStatus.ColorCode
			}
			if iap.ResponsibleUser != nil {
				incidentActionPlanDTO.ResponsibleUser = &organizationDTOs.UserMinimalDTO{
					ID:               int(iap.ResponsibleUser.ID),
					FirstNameEN:      iap.ResponsibleUser.FirstNameEn,
					LastNameEN:       iap.ResponsibleUser.LastNameEn,
					FirstNameAR:      iap.ResponsibleUser.FirstNameAr,
					LastNameAR:       iap.ResponsibleUser.LastNameAr,
					ProfileImagePath: iap.ResponsibleUser.ProfileImagePath,
					Email:            iap.ResponsibleUser.Email,
					Initial:          iap.ResponsibleUser.Initial,
					Color:            iap.ResponsibleUser.Color,
				}
			}
		}

		projectsDTO = append(projectsDTO, &dtos.ProjectReportProjectDetailDTO{
			ID:                           rp.ID,
			ProjectID:                    rp.ProjectID,
			TitleEn:                      rp.TitleEn,
			TitleAr:                      rp.TitleAr,
			DescriptionEn:                rp.DescriptionEn,
			DescriptionAr:                rp.DescriptionAr,
			Budget:                       rp.Budget,
			Scope:                        rp.Scope,
			StartDate:                    rp.StartDate.Format("2006-01-02"),
			EndDate:                      rp.EndDate.Format("2006-01-02"),
			ProjectTypeTitleEn:           rp.ProjectTypeTitleEn,
			ProjectTypeTitleAr:           rp.ProjectTypeTitleAr,
			ProjectStatusTitleEn:         rp.ProjectStatusTitleEn,
			ProjectStatusTitleAr:         rp.ProjectStatusTitleAr,
			ProjectPhaseTitleEn:          rp.ProjectPhaseTitleEn,
			ProjectPhaseTitleAr:          rp.ProjectPhaseTitleAr,
			ProjectClassificationTitleEn: rp.ProjectClassificationTitleEn,
			ProjectClassificationTitleAr: rp.ProjectClassificationTitleAr,
			ResponsibleUserNameEn:        rp.ResponsibleUserNameEn,
			ResponsibleUserNameAr:        rp.ResponsibleUserNameAr,
			IncidentActionPlanID:         rp.IncidentActionPlanID,
			IncidentActionPlanTitle:      rp.IncidentActionPlanTitle,
			IncidentActionPlan:           incidentActionPlanDTO,
			CreatedAt:                    rp.CreatedAt.Format("2006-01-02"),
			Documents:                    docsDTO,
			Milestones:                   milestones,
			Controls:                     controls,
			Risks:                        risks,
			Processes:                    processes,
			StrategicObjectives:          objectives,
			KPIs:                         kpis,
			AssociatedDocuments:          docs,
			Assets:                       assets,
			Incidents:                    incidents,
			ComplianceRequirements:       complianceRequirements,
			ActionPlans:                  actionPlans,
			ChangeRequests:               changeRequests,
			DecisionLogs:                 decisionLogs,
			SupportingDocuments:          supportingDocs,
			Closure:                      closure,
		})
	}

	var createdUserDTO *organizationDTOs.UserMinimalDTO
	if report.CreatedByUser != nil {
		createdUserDTO = &organizationDTOs.UserMinimalDTO{
			ID:               int(report.CreatedByUser.ID),
			FirstNameEN:      report.CreatedByUser.FirstNameEn,
			LastNameEN:       report.CreatedByUser.LastNameEn,
			FirstNameAR:      report.CreatedByUser.FirstNameAr,
			LastNameAR:       report.CreatedByUser.LastNameAr,
			Email:            report.CreatedByUser.Email,
			Initial:          report.CreatedByUser.Initial,
			Color:            report.CreatedByUser.Color,
			ProfileImagePath: report.CreatedByUser.ProfileImagePath,
		}
	}

	return &dtos.ProjectReportDetailsDTO{
		ID:            report.ID,
		Title:         report.Title,
		StartDate:     report.StartDate,
		EndDate:       report.EndDate,
		CreatedAt:     report.CreatedAt,
		CreatedByUser: createdUserDTO,
		Sections:      sectionsDTO,
		Projects:      projectsDTO,
	}, nil
}

func (prs *ProjectReportService) DeleteProjectReport(id uint, currentUserID uint) error {
	tx := prs.repo.BeginTransaction()

	_, err := prs.repo.GetReportByID(id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	if err := prs.repo.DeleteReportDocuments(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := prs.repo.DeleteReportProjects(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := prs.repo.DeleteReportSections(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := prs.repo.DeleteReport(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (prs *ProjectReportService) GetProjectReportSection(projectReportID uint, sectionID uint) (*dtos.ProjectReportSectionDTO, error) {
	section, err := prs.repo.GetSectionByID(sectionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}
	if section.ProjectReportID != projectReportID {
		return nil, utility.NewResourceNotFoundError(fmt.Errorf("section does not belong to report"), sectionID)
	}

	return &dtos.ProjectReportSectionDTO{
		ID:            section.ID,
		Type:          section.Type,
		TitleEn:       section.TitleEn,
		TitleAr:       section.TitleAr,
		DescriptionEn: section.DescriptionEn,
		DescriptionAr: section.DescriptionAr,
		OrderNum:      section.OrderNum,
	}, nil
}

func (prs *ProjectReportService) UpdateProjectReportSection(projectReportID uint, sectionID uint, form dtos.ProjectReportSectionFormDTO, currentUserID uint) error {
	tx := prs.repo.BeginTransaction()

	section, err := prs.repo.GetSectionByID(sectionID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, sectionID)
		}
		return err
	}
	if section.ProjectReportID != projectReportID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("section does not belong to report"), sectionID)
	}

	section.TitleEn = form.TitleEn
	section.TitleAr = form.TitleAr
	if form.DescriptionEn != nil {
		section.DescriptionEn = *form.DescriptionEn
	}
	if form.DescriptionAr != nil {
		section.DescriptionAr = *form.DescriptionAr
	}

	if err := prs.repo.UpdateSection(section, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (prs *ProjectReportService) DeleteProjectReportSection(projectReportID uint, sectionID uint, currentUserID uint) error {
	tx := prs.repo.BeginTransaction()

	section, err := prs.repo.GetSectionByID(sectionID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, sectionID)
		}
		return err
	}
	if section.ProjectReportID != projectReportID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("section does not belong to report"), sectionID)
	}

	if err := prs.repo.DeleteReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := prs.repo.ReorderReportSections(projectReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (prs *ProjectReportService) GetProjectReportExportZip(projectReportID uint, currentUserID uint) (string, error) {
	report, err := prs.GetProjectReportDetails(projectReportID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project report details: %w", err)
	}

	var executiveSummary, scope, objective string
	for _, section := range report.Sections {
		switch section.Type {
		case "executive-summary":
			executiveSummary = section.DescriptionEn
		case "scope":
			scope = section.DescriptionEn
		case "objectives":
			objective = section.DescriptionEn
		}
	}

	var startDateStr, endDateStr string
	if report.StartDate != nil {
		startDateStr = report.StartDate.Format("Jan _2, 2006")
	}
	if report.EndDate != nil {
		endDateStr = report.EndDate.Format("Jan _2, 2006")
	}

	var createdUserFirstNameEn, createdUserFirstNameAr, createdUserLastNameEn, createdUserLastNameAr string
	if report.CreatedByUser != nil {
		createdUserFirstNameEn = report.CreatedByUser.FirstNameEN
		createdUserFirstNameAr = report.CreatedByUser.FirstNameAR
		createdUserLastNameEn = report.CreatedByUser.LastNameEN
		createdUserLastNameAr = report.CreatedByUser.LastNameAR
	}

	pdfDetail := dtos.ProjectReportPDFDetailDTO{
		ID:                     report.ID,
		Title:                  report.Title,
		StartDate:              startDateStr,
		EndDate:                endDateStr,
		ReportedAt:             report.CreatedAt.Format("02, Jan 2006"),
		CreatedUserFirstNameEn: createdUserFirstNameEn,
		CreatedUserFirstNameAr: createdUserFirstNameAr,
		CreatedUserLastNameEn:  createdUserLastNameEn,
		CreatedUserLastNameAr:  createdUserLastNameAr,
		CreatedAt:              report.CreatedAt.Format("2006-01-02"),
		GeneratedAt:            time.Now().Format("Jan _2, 2006, 03:04 PM"),
		ExecutiveSummary:       executiveSummary,
		Scope:                  scope,
		Objective:              objective,
		Projects:               report.Projects,
	}

	// Render main body
	templatePaths := []string{"project/templates/pdf/project_multiple_report_body.html"}
	htmlContent, err := utility.GeneratePDFContent(pdfDetail, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	headerTemplatePaths := []string{"project/templates/pdf/project_multiple_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(pdfDetail, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	footerTemplatePaths := []string{"project/templates/pdf/project_multiple_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(pdfDetail, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	coverTemplatePaths := []string{"project/templates/pdf/project_multiple_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(pdfDetail, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	pdfFileName := fmt.Sprintf("ProjectReport_%d.pdf", projectReportID)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	filesToZip := []string{pdfFileName}

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	for _, p := range report.Projects {
		for _, doc := range p.Documents {
			if doc.Name != nil && *doc.Name != "" {
				storagePath := filepath.Join("project", strconv.FormatUint(uint64(p.ProjectID), 10), "supporting-documents", "documents", *doc.Name)
				fullPath := filepath.Join(baseStoragePath, storagePath)
				if _, err := os.Stat(fullPath); err == nil {
					filesToZip = append(filesToZip, fullPath)
				}
			}
		}
	}

	zipDir := filepath.Join("storage", "zip")
	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("project_report_%d.zip", projectReportID))

	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
	}

	err = utility.CreateZipFile(zipFilePath, filesToZip)
	if err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		_ = os.Remove(pdfFileName)
	}()

	return zipFilePath, nil
}
