package services

import (
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskReportServiceInterface interface {
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
