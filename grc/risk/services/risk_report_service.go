package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"strings"

	"gorm.io/gorm"
)

type RiskReportServiceInterface interface {
	GetAll(params *dtos.RiskReportListQueryParamsDTO) (*dtos.RiskReportListResponseDTO, error)
	GetRiskReportDetails(riskReportID uint) (*dtos.RiskReportDetailDTO, error)
	// GetRiskReportDetails(riskReportID uint, tx ...utility.DatabaseTransaction) (*dtos.RiskReportDetailDTO, error)
	Create(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error)
	GetSectionByID(riskReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportSection, error)
	UpdateSection(riskID uint, id uint, currentUserID uint, riskReporSectionFormDTO *dtos.RiskReportSectionFormDTO) error
	DeleteRiskReportSection(riskReportID uint, sectionID uint, currentUserID uint) error
	DeleteRiskReport(riskReportID uint, currentUserID uint) error
}

type RiskReportService struct {
	repo           repositories.RiskReportRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewRiskReportService(
	repo repositories.RiskReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskReportServiceInterface {
	return &RiskReportService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rrs *RiskReportService) GetAll(params *dtos.RiskReportListQueryParamsDTO) (*dtos.RiskReportListResponseDTO, error) {
	return rrs.repo.GetAll(params)
}

// func (rrs *RiskReportService) GetRiskReportDetails(riskReportID uint) (*dtos.RiskReportDetailDTO, error) {
// 	// Fetch the full risk report with all related data using Preload
// 	riskReport, err := rrs.repo.GetFullDetailsByID(riskReportID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Initialize DTO
// 	var riskReportDTO dtos.RiskReportDetailDTO

// 	// Use mapstructure to map the retrieved model to the DTO
// 	err = mapstructure.Decode(riskReport, &riskReportDTO)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ensure date formatting for the DTO
// 	riskReportDTO.StartDate = riskReport.StartDate
// 	riskReportDTO.EndDate = riskReport.EndDate
// 	riskReportDTO.CreatedAt = riskReport.CreatedAt

// 	return &riskReportDTO, nil
// }

func (rrs *RiskReportService) GetRiskReportDetails(riskReportID uint) (*dtos.RiskReportDetailDTO, error) {
	riskReport, err := rrs.repo.GetRiskReportByID(riskReportID)
	if err != nil {
		return nil, err
	}

	sections, err := rrs.repo.GetRiskReportSections(riskReportID)
	if err != nil {
		return nil, err
	}

	risks, err := rrs.repo.GetRiskIdentifications(riskReportID)
	if err != nil {
		return nil, err
	}

	actionPlans, err := rrs.repo.GetRiskActionPlans(riskReportID)
	if err != nil {
		return nil, err
	}

	kris, err := rrs.repo.GetRiskKRIs(riskReportID)
	if err != nil {
		return nil, err
	}

	// Construct Response DTO
	response := &dtos.RiskReportDetailDTO{
		ID:                       riskReport.ID,
		Title:                    riskReport.Title,
		StartDate:                riskReport.StartDate,
		EndDate:                  riskReport.EndDate,
		CreatedAt:                riskReport.CreatedAt,
		CreatedUserFirstNameEn:   riskReport.CreatedUserFirstNameEn,
		CreatedUserFirstNameAr:   riskReport.CreatedUserFirstNameAr,
		CreatedUserLastNameEn:    riskReport.CreatedUserLastNameEn,
		CreatedUserLastNameAr:    riskReport.CreatedUserLastNameAr,
		OrganizationUnitTitlesEn: riskReport.OrganizationUnitTitlesEn,
		OrganizationUnitTitlesAr: riskReport.OrganizationUnitTitlesAr,
		Sections:                 sections,
		Risks:                    risks,
		ActionPlans:              actionPlans,
		KRIs:                     kris,
	}

	return response, nil
}

func (rrs *RiskReportService) GetSectionByID(riskReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportSection, error) {
	section, err := rrs.repo.GetSectionByID(sectionID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}
	return section, nil
}

func (rrs *RiskReportService) Create(currentUserID uint, riskReportFormDTO *dtos.RiskReportFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

	dates, err := rrs.parseRiskReportDates(riskReportFormDTO)
	if err != nil {
		return 0, err
	}

	// Step 1: Create Risk Report Entry
	riskReport := &models.RiskReport{
		Title:       riskReportFormDTO.Title,
		StartDate:   *dates.StartDate,
		EndDate:     *dates.EndDate,
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

	risks, err := rrs.repo.GetRisksByDateRange(*dates.StartDate, *dates.EndDate, riskReportFormDTO.OrganizationUnitIDs, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, risk := range risks {

		controlTitles, err := rrs.repo.GetControlTitlesByRiskID(risk.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to fetch control titles for risk %d: %w", risk.ID, err)
		}

		existingControls := strings.Join(controlTitles, ", ")

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
			ExistingControls:            &existingControls,
			// ProposedControls:            risk.ProposedControls,
			RiskStatusTitleEn:   risk.RiskStatusTitleEn,
			RiskStatusTitleAr:   risk.RiskStatusTitleAr,
			RiskStatusColorCode: risk.RiskStatusColorCode,
		}

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

			proposedControls, err := rrs.repo.GetActionPlanControlsByActionPlanID(actionPlan.ID, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			riskActionPlan := &models.RiskReportActionPlan{
				RiskReportID:               riskReportID,
				RiskReportIdentificationID: riskIdentification.ID,
				ActionPlanTitle:            actionPlan.ActionPlanTitle,
				ActionPlanDescription:      actionPlan.ActionPlanDescription,
				SeverityLevelTitleEn:       &actionPlan.SeverityLevelTitleEn,
				SeverityLevelTitleAr:       &actionPlan.SeverityLevelTitleAr,
				ResponsibleUserNameEn:      actionPlan.ResponsibleUserNameEn,
				ResponsibleUserNameAr:      actionPlan.ResponsibleUserNameAr,
				TargetDate:                 &actionPlan.TargetDate,
				ActionPlanStatusTitleEn:    actionPlan.ActionPlanStatusTitleEn,
				ActionPlanStatusTitleAr:    actionPlan.ActionPlanStatusTitleAr,
				ActionPlanStatusColorCode:  actionPlan.ActionPlanStatusColorCode,
				ProposedControls:           &proposedControls,
			}

			if err := rrs.repo.CreateRiskActionPlan(riskActionPlan, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

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

		documents, err := rrs.repo.GetDocumentsByRiskID(risk.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, document := range documents {
			riskReportDocument := &models.RiskReportDocument{
				RiskReportID:               riskReportID,
				RiskReportIdentificationID: riskIdentification.ID,
				URL:                        document.URL,
				Name:                       *document.Name,
				Size:                       document.Size,
				Token:                      document.Token,
				OriginalName:               document.OriginalName,
				Type:                       document.Type,
				CreatedByID:                currentUserID,
			}

			if err := rrs.repo.CreateRiskDocument(riskReportDocument, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	tx.Commit()

	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskReport{}.TableName(),
			riskReportID,
			riskReport,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return riskReportID, nil
}

func (rrs *RiskReportService) UpdateSection(riskReportID uint, sectionID uint, currentUserID uint, riskReporSectionFormDTO *dtos.RiskReportSectionFormDTO) error {
	tx := rrs.repo.BeginTransaction()

	existingReportSection, err := rrs.GetSectionByID(riskReportID, sectionID, tx)
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
	existingReportSection.TitleEn = riskReporSectionFormDTO.TitleEn
	existingReportSection.TitleAr = riskReporSectionFormDTO.TitleAr
	existingReportSection.DescriptionEn = *riskReporSectionFormDTO.DescriptionEn
	existingReportSection.DescriptionAr = *riskReporSectionFormDTO.DescriptionAr

	err = rrs.repo.UpdateSection(existingReportSection, tx)
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
			err := rrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskReportSection{}.TableName(),
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

func (rrs *RiskReportService) DeleteRiskReportSection(riskReportID uint, sectionID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	section, err := rrs.GetSectionByID(riskReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the section
	if err := rrs.repo.DeleteRiskReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reorder remaining sections
	if err := rrs.repo.ReorderRiskReportSections(riskReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskReportSection{}.TableName(),
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

func (rrs *RiskReportService) DeleteRiskReport(riskReportID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	existingReport, err := rrs.repo.GetRiskReportByID(riskReportID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskReportOrganizationUnits(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskKRIs(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskActionPlans(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskIdentifications(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskReportSections(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.repo.DeleteRiskReport(riskReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskReport{}.TableName(),
			riskReportID,
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
	case "context":
		return "Risk Governance & Context", "حوكمة المخاطر والسياق"
	case "identification":
		return "Risk Identification", "تحديد المخاطر"
	case "assessment":
		return "Risk Assessment Details", "تفاصيل تقييم المخاطر"
	case "mitigation":
		return "Risk Response & Mitigation", "الاستجابة للمخاطر والتخفيف منها"
	case "kri":
		return "Key Risk Indicators (KRIs) & Performance", "مؤشرات المخاطر الرئيسية والأداء"
	case "communication":
		return "Risk Monitoring, Reporting & Communication", "مراقبة المخاطر والإبلاغ والتواصل"
	case "recommendation":
		return "Recommendations & Action Plan", "التوصيات وخطة العمل"
	case "documents":
		return "Appendices & Supporting Documents", "الملاحق والمستندات الداعمة"
	default:
		return sectionType, sectionType // Default to sectionType if not found
	}
}

func (rrs *RiskReportService) parseRiskReportDates(dto *dtos.RiskReportFormDTO) (*dtos.RiskReportParsedDates, error) {
	parsed := &dtos.RiskReportParsedDates{}

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
