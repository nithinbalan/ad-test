package services

import (
	"errors"
	"os"
	"strings"

	auditListDTOs "jethur/grc/audit/dtos"
	auditDTOs "jethur/grc/kpi/dtos"
	auditRepositories "jethur/grc/kpi/repositories"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditDashboardServiceInterface interface {
	GetAuditStatusBreakdown(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditStatusBreakdownDTO, error)
	GetAuditSummary(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditSummaryDTO, error)
	GetAuditsByStatus(statusID uint, params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditListDTOs.AuditListResponseDTO, error)
	ExportAuditsByStatus(data []*auditListDTOs.AuditListDTO, currentUserID uint) (*os.File, error)
	GetAuditsByStatusExportFilename(statusID uint) string
	GetAuditFindingsSummary(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditFindingsSummaryDTO, error)
	GetAuditFindingsSeverityBreakdown() (*auditDTOs.AuditFindingsSeverityBreakdownDTO, error)
	GetAuditFindingsBySeverityLevel(severityLevelID uint, params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditFindingsBySeverityLevelListResponseDTO, error)
	ExportAuditFindingsBySeverityLevel(data []*auditDTOs.AuditFindingsBySeverityLevelListDTO, currentUserID uint) (*os.File, error)
	GetAuditFindingsBySeverityLevelExportFilename(severityLevelID uint) string
	GetAuditFindingsTrend() ([]auditDTOs.AuditFindingsTrendMonthDTO, error)
	GetOverdueAuditActionPlans(params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditOverdueActionPlanListResponseDTO, error)
	ExportOverdueAuditActionPlans(data []*auditDTOs.AuditOverdueActionPlanListDTO, currentUserID uint) (*os.File, error)
}

type AuditDashboardService struct {
	repo           auditRepositories.AuditDashboardRepositoryInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewAuditDashboardService(repo auditRepositories.AuditDashboardRepositoryInterface, userSettingSvc organizationServices.UserSettingServiceInterface) AuditDashboardServiceInterface {
	return &AuditDashboardService{repo: repo, userSettingSvc: userSettingSvc}
}

func (ads *AuditDashboardService) GetAuditStatusBreakdown(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditStatusBreakdownDTO, error) {
	return ads.repo.GetAuditStatusBreakdown(params)
}

func (ads *AuditDashboardService) GetAuditSummary(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditSummaryDTO, error) {
	return ads.repo.GetAuditSummary(params)
}

func (ads *AuditDashboardService) GetAuditsByStatus(statusID uint, params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditListDTOs.AuditListResponseDTO, error) {
	return ads.repo.GetAuditsByStatus(statusID, params)
}

func (ads *AuditDashboardService) ExportAuditsByStatus(data []*auditListDTOs.AuditListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ads.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"SL No.",
		"Audit Plan Title",
		"Audit Plan Description",
		"Audit Plan Purpose",
		"Audit Plan Objective",
		"Start Date",
		"End Date",
		"Audit Category",
		"Audit Status",
		"Team Lead",
		"Created By",
	}

	rows := make([][]interface{}, 0, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows = append(rows, []interface{}{
				i + 1,
				item.AuditPlanTitle,
				item.AuditPlanDescription,
				item.AuditPlanPurpose,
				item.AuditPlanObjective,
				item.StartDate,
				item.EndDate,
				item.AuditCategoryTitleAr,
				item.AuditStatusTitleAr,
				item.TeamLeadFirstNameAr + " " + item.TeamLeadLastNameAr,
				item.CreatedUserFirstNameAr + " " + item.CreatedUserLastNameAr,
			})
		} else {
			rows = append(rows, []interface{}{
				i + 1,
				item.AuditPlanTitle,
				item.AuditPlanDescription,
				item.AuditPlanPurpose,
				item.AuditPlanObjective,
				item.StartDate,
				item.EndDate,
				item.AuditCategoryTitleEn,
				item.AuditStatusTitleEn,
				item.TeamLeadFirstNameEn + " " + item.TeamLeadLastNameEn,
				item.CreatedUserFirstNameEn + " " + item.CreatedUserLastNameEn,
			})
		}
	}

	file, err := utility.GenerateExcelFile(headers, rows)
	if err != nil {
		return nil, errors.New("failed to generate export file: " + err.Error())
	}
	return file, nil
}

func (ads *AuditDashboardService) GetAuditsByStatusExportFilename(statusID uint) string {
	title, err := ads.repo.GetAuditStatusTitle(statusID)
	if err != nil || title == "" {
		title = "Audits"
	} else {
		title += " Audits"
	}
	return sanitizeFilename(title + ".xlsx")
}

func (ads *AuditDashboardService) GetAuditFindingsSummary(params *auditDTOs.AuditDashboardQueryParamsDTO) (*auditDTOs.AuditFindingsSummaryDTO, error) {
	return ads.repo.GetAuditFindingsSummary(params)
}

func (ads *AuditDashboardService) GetAuditFindingsSeverityBreakdown() (*auditDTOs.AuditFindingsSeverityBreakdownDTO, error) {
	return ads.repo.GetAuditFindingsSeverityBreakdown()
}

func (ads *AuditDashboardService) GetAuditFindingsBySeverityLevel(severityLevelID uint, params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditFindingsBySeverityLevelListResponseDTO, error) {
	return ads.repo.GetAuditFindingsBySeverityLevel(severityLevelID, params)
}

func (ads *AuditDashboardService) ExportAuditFindingsBySeverityLevel(data []*auditDTOs.AuditFindingsBySeverityLevelListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ads.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"SL No.",
		"Reference Number",
		"Title",
		"Audit Title",
		"Auditable Item",
		"Finding Category",
		"Severity Level",
		"Status",
		"Created By",
		"Created At",
		"Updated By",
	}

	rows := make([][]interface{}, 0, len(data))
	for i, item := range data {
		updatedByEn := ""
		if item.UpdatedUserFirstNameEn != nil && item.UpdatedUserLastNameEn != nil {
			updatedByEn = *item.UpdatedUserFirstNameEn + " " + *item.UpdatedUserLastNameEn
		}
		updatedByAr := ""
		if item.UpdatedUserFirstNameAr != nil && item.UpdatedUserLastNameAr != nil {
			updatedByAr = *item.UpdatedUserFirstNameAr + " " + *item.UpdatedUserLastNameAr
		}
		if userSetting.Language == "ar" {
			rows = append(rows, []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.Title,
				item.AuditTitle,
				item.AuditableItemTitleAr,
				item.AuditFindingCategoryTitleAr,
				item.SeverityLevelTitleAr,
				item.AuditFindingStatusTitleAr,
				item.CreatedUserFirstNameAr + " " + item.CreatedUserLastNameAr,
				item.CreatedAt,
				updatedByAr,
			})
		} else {
			rows = append(rows, []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.Title,
				item.AuditTitle,
				item.AuditableItemTitleEn,
				item.AuditFindingCategoryTitleEn,
				item.SeverityLevelTitleEn,
				item.AuditFindingStatusTitleEn,
				item.CreatedUserFirstNameEn + " " + item.CreatedUserLastNameEn,
				item.CreatedAt,
				updatedByEn,
			})
		}
	}

	file, err := utility.GenerateExcelFile(headers, rows)
	if err != nil {
		return nil, errors.New("failed to generate export file: " + err.Error())
	}
	return file, nil
}

func (ads *AuditDashboardService) GetAuditFindingsBySeverityLevelExportFilename(severityLevelID uint) string {
	title, err := ads.repo.GetSeverityLevelTitle(severityLevelID)
	if err != nil || title == "" {
		title = "Audit Findings"
	}
	return sanitizeFilename(title + " Severity Findings.xlsx")
}

func (ads *AuditDashboardService) GetAuditFindingsTrend() ([]auditDTOs.AuditFindingsTrendMonthDTO, error) {
	return ads.repo.GetAuditFindingsTrend()
}

func (ads *AuditDashboardService) GetOverdueAuditActionPlans(params *auditDTOs.AuditDashboardListQueryParamsDTO) (*auditDTOs.AuditOverdueActionPlanListResponseDTO, error) {
	return ads.repo.GetOverdueAuditActionPlans(params)
}

func (ads *AuditDashboardService) ExportOverdueAuditActionPlans(data []*auditDTOs.AuditOverdueActionPlanListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ads.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"SL No.",
		"Reference Number",
		"Finding Title",
		"Action Plan Title",
		"Severity Level",
		"Responsible User",
		"Department",
		"Target Date",
		"Overdue Days",
		"Status",
		"Created At",
	}

	rows := make([][]interface{}, 0, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows = append(rows, []interface{}{
				i + 1,
				item.AuditFindingReferenceNumber,
				item.AuditFindingTitle,
				item.Title,
				item.SeverityLevelTitleAr,
				item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr,
				item.ResponsibleUserDepartmentTitleAr,
				item.TargetDate,
				item.OverdueDays,
				item.ActionPlanStatusTitleAr,
				item.CreatedAt,
			})
		} else {
			rows = append(rows, []interface{}{
				i + 1,
				item.AuditFindingReferenceNumber,
				item.AuditFindingTitle,
				item.Title,
				item.SeverityLevelTitleEn,
				item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn,
				item.ResponsibleUserDepartmentTitleEn,
				item.TargetDate,
				item.OverdueDays,
				item.ActionPlanStatusTitleEn,
				item.CreatedAt,
			})
		}
	}

	file, err := utility.GenerateExcelFile(headers, rows)
	if err != nil {
		return nil, errors.New("failed to generate export file: " + err.Error())
	}
	return file, nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", "*", "-", "?", "", "\"", "", "<", "", ">", "", "|", "")
	return strings.TrimSpace(replacer.Replace(name))
}
