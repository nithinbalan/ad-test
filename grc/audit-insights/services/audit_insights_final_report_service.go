package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
)

type AuditInsightsFinalReportServiceInterface interface {
	GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error)
	GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error)
	ExportAuditInsightsFinalReports(data []*auditDTOs.AuditFinalReportListDTO, currentUserID uint) (*os.File, error)
}

type AuditInsightsFinalReportService struct {
	repo                repositories.AuditInsightsFinalReportRepositoryInterface
	userSvc             organizationServices.UserServiceInterface
	auditSvc            auditServices.AuditServiceInterface
	auditFinalReportSvc auditServices.AuditFinalReportServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewAuditInsightsFinalReportService(
	repo repositories.AuditInsightsFinalReportRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	auditSvc auditServices.AuditServiceInterface,
	auditFinalReportSvc auditServices.AuditFinalReportServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) AuditInsightsFinalReportServiceInterface {
	return &AuditInsightsFinalReportService{
		repo:                repo,
		userSvc:             userSvc,
		auditSvc:            auditSvc,
		auditFinalReportSvc: auditFinalReportSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (aidrs *AuditInsightsFinalReportService) GetAll(params *auditDTOs.AuditFinalReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFinalReportListResponseDTO, error) {
	orgUnitID, err := aidrs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidrs.repo.GetAll(params, orgUnitID)
}

func (aidrs *AuditInsightsFinalReportService) GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error) {
	auditFinalReport, err := aidrs.auditFinalReportSvc.GetByID(id)
	if err != nil {
		return nil, err
	}
	auditID := auditFinalReport.AuditID
	audit, err := aidrs.auditFinalReportSvc.GetFullDetailsByID(auditID, currentUserID, "audit-insight", tx...)
	if err != nil {
		return nil, err
	}
	return audit, nil
}

func (aifrs *AuditInsightsFinalReportService) ExportAuditInsightsFinalReports(data []*auditDTOs.AuditFinalReportListDTO, currentUserID uint) (*os.File, error) {
	// Get user settings for language
	userSetting, err := aifrs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Audit Plan Title",
			"Audit Plan Description",
			"Audit Plan Purpose",
			"Audit Plan Objective",
			"Audit Category (AR)",
			"Audit Status (AR)",
			"Team Lead (AR)",
			"Submitted Date",
			"Final Report Status (AR)",
			"Submitted By (AR)",
			"Created By (AR)",
		}
	} else {
		headers = []string{
			"SL No.",
			"Audit Plan Title",
			"Audit Plan Description",
			"Audit Plan Purpose",
			"Audit Plan Objective",
			"Audit Category (EN)",
			"Audit Status (EN)",
			"Team Lead (EN)",
			"Submitted Date",
			"Final Report Status (EN)",
			"Submitted By (EN)",
			"Created By (EN)",
		}
	}

	// Prepare rows for export
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				item.AuditPlanTitle,
				item.AuditPlanDescription,
				item.AuditPlanPurpose,
				item.AuditPlanObjective,
				item.AuditCategoryTitleAr,
				item.AuditStatusTitleAr,
				item.TeamLeadFirstNameAr + " " + item.TeamLeadLastNameAr,
				item.SubmittedAt,
				item.AuditFinalReportStatusTitleAr,
				item.SubmittedUserFirstNameAr + " " + item.SubmittedUserLastNameAr,
				item.CreatedUserFirstNameAr + " " + item.CreatedUserLastNameAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				item.AuditPlanTitle,
				item.AuditPlanDescription,
				item.AuditPlanPurpose,
				item.AuditPlanObjective,
				item.AuditCategoryTitleEn,
				item.AuditStatusTitleEn,
				item.TeamLeadFirstNameEn + " " + item.TeamLeadLastNameEn,
				item.SubmittedAt,
				item.AuditFinalReportStatusTitleEn,
				item.SubmittedUserFirstNameEn + " " + item.SubmittedUserLastNameEn,
				item.CreatedUserFirstNameEn + " " + item.CreatedUserLastNameEn,
			}
		}
	}

	// Generate Excel file
	return utility.GenerateExcelFile(headers, rows)
}
