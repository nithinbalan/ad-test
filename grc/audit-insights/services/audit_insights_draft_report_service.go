package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
)

type AuditInsightsDraftReportServiceInterface interface {
	GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error)
	GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error)
	ExportAuditInsightsDraftReports(data []*auditDTOs.AuditDraftReportListDTO, currentUserID uint) (*os.File, error)
}

type AuditInsightsDraftReportService struct {
	repo                repositories.AuditInsightsDraftReportRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	auditSvc            auditServices.AuditServiceInterface
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewAuditInsightsDraftReportService(
	repo repositories.AuditInsightsDraftReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditSvc auditServices.AuditServiceInterface,
	auditDraftReportSvc auditServices.AuditDraftReportServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) AuditInsightsDraftReportServiceInterface {
	return &AuditInsightsDraftReportService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		auditSvc:            auditSvc,
		auditDraftReportSvc: auditDraftReportSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (aidrs *AuditInsightsDraftReportService) GetAll(params *auditDTOs.AuditDraftReportListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditDraftReportListResponseDTO, error) {
	orgUnitID, err := aidrs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aidrs.repo.GetAll(params, orgUnitID)
}

func (aidrs *AuditInsightsDraftReportService) GetFullDetailsByID(id uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftReportDetailDTO, error) {
	auditDraftReport, err := aidrs.auditDraftReportSvc.GetByID(id)
	if err != nil {
		return nil, err
	}
	auditID := auditDraftReport.AuditID
	audit, err := aidrs.auditDraftReportSvc.GetFullDetailsByID(auditID, currentUserID, "audit-insight", tx...)
	if err != nil {
		return nil, err
	}
	return audit, nil
}

func (aidrs *AuditInsightsDraftReportService) ExportAuditInsightsDraftReports(data []*auditDTOs.AuditDraftReportListDTO, currentUserID uint) (*os.File, error) {
	// Get user settings (for language preference)
	userSetting, err := aidrs.userSettingSvc.GetUserSetting(currentUserID)
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
			"Report Date",
			"Draft Report Status (AR)",
			"Internal Review Submitted By (AR)",
			"Auditee Review Submitted By (AR)",
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
			"Report Date",
			"Draft Report Status (EN)",
			"Internal Review Submitted By (EN)",
			"Auditee Review Submitted By (EN)",
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
				item.ReportDate,
				item.AuditDraftReportStatusTitleAr,
				item.InternalReviewSubmittedUserFirstNameAr + " " + item.InternalReviewSubmittedUserLastNameAr,
				item.AuditeeReviewSubmittedUserFirstNameAr + " " + item.AuditeeReviewSubmittedUserLastNameAr,
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
				item.ReportDate,
				item.AuditDraftReportStatusTitleEn,
				item.InternalReviewSubmittedUserFirstNameEn + " " + item.InternalReviewSubmittedUserLastNameEn,
				item.AuditeeReviewSubmittedUserFirstNameEn + " " + item.AuditeeReviewSubmittedUserLastNameEn,
				item.CreatedUserFirstNameEn + " " + item.CreatedUserLastNameEn,
			}
		}
	}

	// Call utility to generate Excel file
	return utility.GenerateExcelFile(headers, rows)
}
