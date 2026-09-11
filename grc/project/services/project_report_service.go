package services

import (
	"errors"
	"fmt"
	"time"

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
