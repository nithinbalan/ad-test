package services

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepo "jethur/grc/meeting/repositories"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"time"
)

type MeetingReportServiceInterface interface {
	ExportMeetingReport(id uint, currentUserID uint) (string, error)
	GetFullReportDetailsByID(id uint) (*meetingDTOs.MeetingReportDetailDTO, error)
}

// Service and implementations
type MeetingReportService struct {
	repo      meetingRepo.MeetingReportRepositoryInterface
	settigSvc settingSvc.ReportTemplateSettingServiceInterface
}

func NewIncidentReportService(
	repo meetingRepo.MeetingReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,

) MeetingReportServiceInterface {
	return &MeetingReportService{
		repo: repo,
		settigSvc :settigSvc,
	}
}

func (mrs MeetingReportService) GetFullReportDetailsByID(id uint) (*meetingDTOs.MeetingReportDetailDTO, error) {
	// Initialize your DTO
	var meetingDTO meetingDTOs.MeetingReportDetailDTO
	meeting, err := mrs.repo.GetFullReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	if _, err := mrs.settigSvc.GetImagePath("logo"); err != nil {
		return nil, utility.NewForbiddenError("logo image missing")
	}
	if _, err := mrs.settigSvc.GetImagePath("cover-page-bg"); err != nil {
		return nil, utility.NewForbiddenError("cover image missing")
	}

	err = mapstructure.Decode(meeting, &meetingDTO)
	if err != nil {
		return nil, err
	}

	kpis, err := mrs.repo.GetAllMeetingKPI(id)
	if err != nil {
		return nil, err
	}

	meetingDTO.KPIs = kpis

	strategicObjectives, err := mrs.repo.GetAllMeetingStrategicObjectives(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.StrategicObjectives = strategicObjectives

	actionPlans, err := mrs.repo.GetAllMeetingActionPlans(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.ActionPlans = actionPlans

	assets, err := mrs.repo.GetAllMeetingAssets(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.Assets = assets

	associatedDocuments, err := mrs.repo.GetAllMeetingDocuments(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.AssociatedDocuments = associatedDocuments

	controls, err := mrs.repo.GetAllMeetingControls(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.Controls = controls

	risks, err := mrs.repo.GetAllMeetingRisks(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.Risks = risks

	process, err := mrs.repo.GetAllMeetingProcesses(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.Processes = process

	complianceRequirements, err := mrs.repo.GetAllMeetingComplianceRequirements(id)
	if err != nil {
		return nil, err
	}
	meetingDTO.ComplianceRequirements = complianceRequirements

	meetingDTO.AppUrl = os.Getenv("APP_URL")
	meetingDTO.CreatedAt = meeting.CreatedAt
	meetingDTO.UpdatedAt = meeting.UpdatedAt

	return &meetingDTO, nil
}

func (irs MeetingReportService) ExportMeetingReport(id uint, currentUserID uint) (string, error) {
	meetingDTO, err := irs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}

	meetingDTO.GeneratedAt = time.Now().Format("2006-01-02 03:04 PM")

	// Main content template
	templatePaths := []string{"meeting/templates/pdf/meeting_report.html"}
	htmlContent, err := utility.GeneratePDFContent(meetingDTO, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Header template - process with data to resolve template variables
	headerTemplatePaths := []string{"meeting/templates/pdf/meeting_report_header.html"}
	headerContent, err := utility.GeneratePDFContent(meetingDTO, headerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer template - process with data to resolve template variables
	footerTemplatePaths := []string{"meeting/templates/pdf/meeting_report_footer.html"}
	footerContent, err := utility.GeneratePDFContent(meetingDTO, footerTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	// cover template - process with data to resolve template variables
	coverTemplatePaths := []string{"meeting/templates/pdf/meeting_report_cover.html"}
	coverContent, err := utility.GeneratePDFContent(meetingDTO, coverTemplatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	//Convert HTML → PDF with header and footer
	pdfFileName := fmt.Sprintf("MeetingReport_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}
