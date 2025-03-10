package services

import (
	"fmt"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
)

type RiskReportServiceInterface interface {
	GenerateRiskReport(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error)
}

type RiskReportService struct {
	repo     repositories.RiskReportRepositoryInterface
	riskRepo repositories.RiskRegisterRepositoryInterface
}

func NewRiskReportService(
	repo repositories.RiskReportRepositoryInterface,
	riskRepo repositories.RiskRegisterRepositoryInterface,
) RiskReportServiceInterface {
	return &RiskReportService{
		repo:     repo,
		riskRepo: riskRepo,
	}
}

func (rrs *RiskReportService) GenerateRiskReport(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

	// Step 1: Create Risk Report Entry
	riskReport := &models.RiskReport{
		Title:       riskReportFormDTO.Title,
		StartDate:   riskReportFormDTO.StartDate,
		EndDate:     riskReportFormDTO.EndDate,
		CreatedByID: currentUserID,
	}

	riskReportID, err := rrs.repo.Create(riskReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to create risk report: %w", err)
	}

	// Step 2: Fetch Risks in Date Range
	risks, err := rrs.repo.GetRisksByDateRange(riskReportFormDTO.StartDate, riskReportFormDTO.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to fetch risks: %w", err)
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
			return 0, fmt.Errorf("failed to create risk identification: %w", err)
		}
	}

	// Step 5: Commit Transaction
	tx.Commit()
	return riskReportID, nil
}
