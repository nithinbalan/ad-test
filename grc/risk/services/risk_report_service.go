package services

import (
	"fmt"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type RiskReportServiceInterface interface {
	GetRiskReportDetails(riskReportID uint, tx ...utility.DatabaseTransaction) (*dtos.RiskReportDetailDTO, error)
	GenerateRiskReport(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error)
}

type RiskReportService struct {
	repo repositories.RiskReportRepositoryInterface
}

func NewRiskReportService(
	repo repositories.RiskReportRepositoryInterface,
) RiskReportServiceInterface {
	return &RiskReportService{
		repo: repo,
	}
}

func (rrs *RiskReportService) GetRiskReportDetails(riskReportID uint, tx ...utility.DatabaseTransaction) (*dtos.RiskReportDetailDTO, error) {
	// Fetch Risk Report Details
	riskReport, err := rrs.repo.GetRiskReportByID(riskReportID, tx...)
	if err != nil {
		return nil, err
	}

	// Initialize DTO
	riskReportDTO := &dtos.RiskReportDetailDTO{
		ID:        riskReport.ID,
		Title:     riskReport.Title,
		StartDate: riskReport.StartDate,
		EndDate:   riskReport.EndDate,
		// CreatedBy: riskReport.CreatedByUser.FirstNameEn(),
		CreatedAt: riskReport.CreatedAt,
	}

	// Fetch Sections
	sections, err := rrs.repo.GetSectionsByRiskReportID(riskReportID, tx...)
	if err != nil {
		return nil, err
	}
	for _, section := range sections {
		riskReportDTO.Sections = append(riskReportDTO.Sections, dtos.RiskReportSectionDTO{
			Type:          section.Type,
			TitleEn:       section.TitleEn,
			TitleAr:       section.TitleAr,
			DescriptionEn: section.DescriptionEn,
			DescriptionAr: section.DescriptionAr,
			OrderNum:      section.OrderNum,
		})
	}

	// Fetch Identified Risks
	identifications, err := rrs.repo.GetIdentificationsByRiskReportID(riskReportID, tx...)
	if err != nil {
		return nil, err
	}
	for _, identification := range identifications {
		riskReportDTO.Identifications = append(riskReportDTO.Identifications, dtos.RiskIdentificationDTO{
			ReferenceNumber: identification.ReferenceNumber,
			Title:           identification.Title,
			Description:     identification.Description,
			RiskCategory:    identification.RiskCategoryTitleEn,
			// RiskStatus:      identification.RiskStatusTitleEn,
		})
	}

	// Fetch Action Plans
	actionPlans, err := rrs.repo.GetActionPlansByRiskReportID(riskReportID, tx...)
	if err != nil {
		return nil, err
	}
	for _, actionPlan := range actionPlans {
		riskReportDTO.ActionPlans = append(riskReportDTO.ActionPlans, dtos.RiskActionPlanDTO{
			// ReferenceNumber: actionPlan.RiskReportIdentification.ReferenceNumber,
			Title:           actionPlan.ActionPlanTitle,
			Description:     actionPlan.ActionPlanDescription,
			SeverityLevel:   actionPlan.SeverityLevelTitleEn,
			ResponsibleUser: actionPlan.ResponsibleUserNameEn,
			// OrganizationUnit: actionPlan.OrganizationUnitTitleEn,
			TargetDate: actionPlan.TargetDate.Format("2006-01-02"),
			Status:     actionPlan.ActionPlanStatusTitleEn,
		})
	}

	// Fetch KRIs
	kris, err := rrs.repo.GetKRIsByRiskReportID(riskReportID, tx...)
	if err != nil {
		return nil, err
	}
	for _, kri := range kris {
		riskReportDTO.KRIs = append(riskReportDTO.KRIs, dtos.RiskReportKriDTO{
			// ReferenceNumber: kri.RiskReportIdentification.ReferenceNumber,
			Title:       kri.Title,
			Description: kri.Description,
			// Threshold:       fmt.Sprintf("Low: %v, Medium: %v, High: %v", kri.LowThreshold, kri.MediumThreshold, kri.HighThreshold),
			CurrentValue: fmt.Sprintf("%v", kri.CurrentValue),
			Status:       "Normal", // You can add logic to determine the status based on thresholds
		})
	}

	return riskReportDTO, nil
}

func (rrs *RiskReportService) GenerateRiskReport(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

	// Step 1: Create Risk Report Entry
	riskReport := &models.RiskReport{
		Title:       riskReportFormDTO.Title,
		StartDate:   riskReportFormDTO.StartDate,
		EndDate:     riskReportFormDTO.EndDate,
		CreatedBYID: currentUserID,
	}

	riskReportID, err := rrs.repo.Create(riskReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, orgUnitID := range riskReportFormDTO.OrganizationUnitIDs {
		orgUnit := &models.RiskReportOrganizationUnit{
			RiskReportID:       riskReportID,
			OrganizationUnitID: orgUnitID,
			CreatedBYID:        currentUserID,
		}

		if err := rrs.repo.CreateRiskReportOrganizationUnit(orgUnit, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	sectionTypes := []string{
		"executive-summary", "context", "identification", "assessment",
		"mitigation", "kri", "communication", "recommendation", "documents",
	}

	for i, sectionType := range sectionTypes {
		titleEn, titleAr := formatTitle(sectionType)
		section := models.RiskReportSection{
			RiskReportID:  riskReportID,
			Type:          sectionType,
			TitleEn:       titleEn,
			TitleAr:       titleAr,
			DescriptionEn: "",
			DescriptionAr: "",
			OrderNum:      i + 1,
		}
		if err := rrs.repo.CreateSection(&section, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Step 2: Fetch Risks in Date Range
	risks, err := rrs.repo.GetRisksByDateRange(riskReportFormDTO.StartDate, riskReportFormDTO.EndDate, riskReportFormDTO.OrganizationUnitIDs, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Step 3: Populate Risk Report Identifications
	for _, risk := range risks {
		riskIdentification := &models.RiskReportIdentification{
			RiskReportID:                riskReportID,
			ReferenceNumber:             risk.ReferenceNumber,
			Title:                       risk.Title,
			Description:                 risk.Description,
			RiskReportedAt:              risk.ReportedAt,
			RiskCategoryTitleEn:         risk.RiskCategoryTitleEn,
			RiskCategoryTitleAr:         risk.RiskCategoryTitleAr,
			RiskSourceTitleEn:           risk.RiskSourceTitleEn,
			RiskSourceTitleAr:           risk.RiskSourceTitleAr,
			ResponsibleUnitTitleEn:      risk.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:      risk.ResponsibleUnitTitleAr,
			InherentRiskRatingTitleEn:   risk.InherentRiskRatingTitleEn,
			InherentRiskRatingTitleAr:   risk.InherentRiskRatingTitleAr,
			InherentRiskRatingColorCode: risk.InherentRiskRatingColorCode,
			InherentLikelihoodTitleEn:   risk.InherentLikelihoodTitleEn,
			InherentLikelihoodTitleAr:   risk.InherentLikelihoodTitleAr,
			InherentImpactTitleEn:       risk.InherentImpactTitleEn,
			InherentImpactTitleAr:       risk.InherentImpactTitleAr,
			ResidualRiskRatingTitleEn:   risk.ResidualRiskRatingTitleEn,
			ResidualRiskRatingTitleAr:   risk.ResidualRiskRatingTitleAr,
			ResidualRiskRatingColorCode: risk.ResidualRiskRatingColorCode,
			ResidualLikelihoodTitleEn:   risk.ResidualLikelihoodTitleEn,
			ResidualLikelihoodTitleAr:   risk.ResidualLikelihoodTitleAr,
			ResidualImpactTitleEn:       risk.ResidualImpactTitleEn,
			ResidualImpactTitleAr:       risk.ResidualImpactTitleAr,
			// ExistingControls:            risk.ExistingControls,
			// ProposedControls:            risk.ProposedControls,
			RiskStatusTitleEn:   risk.RiskStatusTitleEn,
			RiskStatusTitleAr:   risk.RiskStatusTitleAr,
			RiskStatusColorCode: risk.RiskStatusColorCode,
		}

		// Step 4: Insert Risk Identification Entry
		if err := rrs.repo.CreateRiskIdentification(riskIdentification, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		actionPlans, err := rrs.repo.GetActionPlansByRiskID(risk.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, actionPlan := range actionPlans {
			riskActionPlan := &models.RiskReportActionPlan{
				RiskReportID:               riskReportID,
				RiskReportIdentificationID: riskIdentification.ID,
				ActionPlanTitle:            actionPlan.ActionPlanTitle,
				ActionPlanDescription:      actionPlan.ActionPlanDescription,
				SeverityLevelTitleEn:       actionPlan.SeverityLevelTitleEn,
				SeverityLevelTitleAr:       actionPlan.SeverityLevelTitleAr,
				ResponsibleUserNameEn:      actionPlan.ResponsibleUserNameEn,
				ResponsibleUserNameAr:      actionPlan.ResponsibleUserNameAr,
				TargetDate:                 actionPlan.TargetDate,
				ActionPlanStatusTitleEn:    actionPlan.ActionPlanStatusTitleEn,
				ActionPlanStatusTitleAr:    actionPlan.ActionPlanStatusTitleAr,
				ActionPlanStatusColorCode:  actionPlan.ActionPlanStatusColorCode,
			}

			if err := rrs.repo.CreateRiskActionPlan(riskActionPlan, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

		//Update KRI table
		kris, err := rrs.repo.GetKRIsByRiskID(risk.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, kri := range kris {
			riskKRI := &models.RiskReportKRI{
				RiskReportID:               riskReportID,
				RiskReportIdentificationID: riskIdentification.ID,
				Title:                      kri.Title,
				Description:                kri.Description,
				LowThreshold:               &kri.LowThreshold,
				MediumThreshold:            &kri.MediumThreshold,
				UnitTitleEn:                &kri.UnitTitleEn,
				UnitTitleAr:                &kri.UnitTitleAr,
				CurrentValue:               &kri.CurrentValue,
			}

			if err := rrs.repo.CreateRiskKRI(riskKRI, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	// Step 5: Commit Transaction
	tx.Commit()
	return riskReportID, nil
}

func formatTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "context":
		return "Risk Governance & Context", "حوكمة المخاطر والسياق"
	case "identification":
		return "Risk Identification & Assessment", "تحديد وتقييم المخاطر"
	case "assessment":
		return "Risk Evaluation & Impact", "تقييم المخاطر وتأثيرها"
	case "mitigation":
		return "Risk Mitigation Strategies", "استراتيجيات التخفيف من المخاطر"
	case "kri":
		return "Key Risk Indicators (KRIs)", "مؤشرات المخاطر الرئيسية"
	case "communication":
		return "Risk Communication & Reporting", "الاتصال والإبلاغ عن المخاطر"
	case "recommendation":
		return "Recommendations & Action Plan", "التوصيات وخطة العمل"
	case "documents":
		return "Supporting Documents", "المستندات الداعمة"
	default:
		return sectionType, sectionType // Default to sectionType if not found
	}
}
