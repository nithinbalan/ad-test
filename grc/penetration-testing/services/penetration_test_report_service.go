package services

import (
	"errors"
	"fmt"
	"io"
	"jethur/config"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	penetrationTestRepo "jethur/grc/penetration-testing/repositories"
	settingSvc "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestReportServiceInterface interface {
	GetFullReportDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestReportPDFDetailDTO, error)
	ExportPenetrationTestReport(id uint, currentUserID uint) (string, error)
	GetPenetrationTestReportExportZip(penetrationTestReportID uint, currentUserID uint, language string) (string, error)
	GetAll(params *penetrationTestDTOs.PenetrationTestReportListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestReportListResponseDTO, error)
	Create(currentUserID uint, dto penetrationTestDTOs.PenetrationTestReportRequestDTO) (uint, error)
	GetPenetrationTestReportDetails(id uint) (*penetrationTestDTOs.PenetrationTestReportDetailsDTO, error)
	DeletePenetrationTestReport(penetrationTestReportID uint, currentUserID uint) error
	GetSectionByID(penetrationTestReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReportSection, error)
	UpdateSection(penetrationTestReportID uint, sectionID uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestReportSectionFormDTO) error
	DeletePenetrationTestReportSection(penetrationTestReportID uint, sectionID uint, currentUserID uint) error
}

type PenetrationTestReportService struct {
	repo           penetrationTestRepo.PenetrationTestReportRepositoryInterface
	settigSvc      settingSvc.ReportTemplateSettingServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestReportService(
	repo penetrationTestRepo.PenetrationTestReportRepositoryInterface,
	settigSvc settingSvc.ReportTemplateSettingServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) PenetrationTestReportServiceInterface {
	return &PenetrationTestReportService{
		repo:           repo,
		settigSvc:      settigSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (ptrs *PenetrationTestReportService) GetAll(params *penetrationTestDTOs.PenetrationTestReportListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestReportListResponseDTO, error) {
	return ptrs.repo.GetAll(params)
}

func (ptrs *PenetrationTestReportService) parsePenetrationTestReportDates(dto penetrationTestDTOs.PenetrationTestReportRequestDTO) (*penetrationTestDTOs.PenetrationTestReportParsedDates, error) {
	parsed := &penetrationTestDTOs.PenetrationTestReportParsedDates{}

	var err error

	if dto.StartDate != nil {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != nil {
		parsed.EndDate, err = utility.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}

func (ptrs *PenetrationTestReportService) GetSectionByID(penetrationTestReportID uint, sectionID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReportSection, error) {
	section, err := ptrs.repo.GetSectionByID(sectionID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionID)
		}
		return nil, err
	}

	if section.PenetrationTestReportID != penetrationTestReportID {
		return nil, utility.NewResourceNotFoundError(
			errors.New("section does not belong to this report"),
			sectionID,
		)
	}

	return section, nil
}

func (ptrs *PenetrationTestReportService) UpdateSection(penetrationTestReportID uint, sectionID uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestReportSectionFormDTO) error {
	tx := ptrs.repo.BeginTransaction()

	existingSection, err := ptrs.GetSectionByID(penetrationTestReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldValues := map[string]interface{}{
		"title_en":       existingSection.TitleEn,
		"title_ar":       existingSection.TitleAr,
		"description_en": existingSection.DescriptionEn,
		"description_ar": existingSection.DescriptionAr,
	}

	existingSection.TitleEn = formDTO.TitleEn
	existingSection.TitleAr = formDTO.TitleAr
	if formDTO.DescriptionEn != nil {
		existingSection.DescriptionEn = formDTO.DescriptionEn
	}
	if formDTO.DescriptionAr != nil {
		existingSection.DescriptionAr = formDTO.DescriptionAr
	}

	err = ptrs.repo.UpdateSection(existingSection, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	newValues := map[string]interface{}{
		"title_en":       existingSection.TitleEn,
		"title_ar":       existingSection.TitleAr,
		"description_en": existingSection.DescriptionEn,
		"description_ar": existingSection.DescriptionAr,
	}

	changedData := utility.CompareFieldChanges(oldValues, newValues)

	if len(changedData) > 0 {
		go func() {
			err := ptrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.PenetrationTestReportSection{}.TableName(),
				sectionID,
				changedData,
				currentUserID,
				"penetration_testing", "penetration_testing_report", "report", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return nil
}

func (ptrs *PenetrationTestReportService) DeletePenetrationTestReportSection(penetrationTestReportID uint, sectionID uint, currentUserID uint) error {
	tx := ptrs.repo.BeginTransaction()

	section, err := ptrs.GetSectionByID(penetrationTestReportID, sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportSection(sectionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.ReorderPenetrationTestReportSections(penetrationTestReportID, section.OrderNum, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	go func() {
		err := ptrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestReportSection{}.TableName(),
			sectionID,
			section,
			currentUserID,
			"penetration_testing", "penetration_testing_report", "report", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (ptrs *PenetrationTestReportService) Create(currentUserID uint, dto penetrationTestDTOs.PenetrationTestReportRequestDTO) (uint, error) {
	tx := ptrs.repo.BeginTransaction()

	template, err := ptrs.settigSvc.GetFirst()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if template.CoverPageBGImagePath == nil || template.LogoImagePath == nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Logo and Background Image not added.")
	}

	dates, err := ptrs.parsePenetrationTestReportDates(dto)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	queryParams := &penetrationTestDTOs.PenetrationTestReportQueryParamsDTO{
		Years:                            dto.Years,
		PenetrationTestExecutorIDs:       dto.PenetrationTestExecutorIDs,
		PenetrationTestTypeIDs:           dto.PenetrationTestTypeIDs,
		PenetrationTestClassificationIDs: dto.PenetrationTestClassificationIDs,
		PenetrationTestPlanIDs:           dto.PenetrationTestPlanIDs,
		CSResponsibleEmployeeIDs:         dto.CSResponsibleEmployeeIDs,
		ResponsibleUnitIDs:               dto.ResponsibleUnitIDs,
		PenetrationTestStatusIDs:         dto.PenetrationTestStatusIDs,
	}
	switch dto.Type {
	case "year":
		if len(dto.Years) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(errors.New("years are required for year reports"), "required", "years")
		}
		firstYear, lastYear := dto.Years[0], dto.Years[0]
		for _, year := range dto.Years {
			if year < firstYear {
				firstYear = year
			}
			if year > lastYear {
				lastYear = year
			}
		}
		start := time.Date(firstYear, time.January, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(lastYear, time.December, 31, 0, 0, 0, 0, time.UTC)
		dates.StartDate = &start
		dates.EndDate = &end
	case "date-range":
		if dates.StartDate == nil || dates.EndDate == nil {
			tx.Rollback()
			return 0, utility.NewValidationError(errors.New("start_date and end_date are required for date-range reports"), "required", "start_date")
		}
		if dates.StartDate.After(*dates.EndDate) {
			tx.Rollback()
			return 0, utility.NewValidationError(errors.New("start_date must be less than or equal to end_date"), "invalid", "start_date")
		}
		queryParams.StartDate = dates.StartDate
		queryParams.EndDate = dates.EndDate
	case "individual":
		if dto.PenetrationTestIDs == nil || len(*dto.PenetrationTestIDs) == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(errors.New("penetration_test_ids are required for individual reports"), "required", "penetration_test_ids")
		}
		queryParams.PenetrationTestIDs = *dto.PenetrationTestIDs
	case "filtered":
	default:
		tx.Rollback()
		return 0, utility.NewValidationError(errors.New("invalid report type"), "invalid", "type")
	}
	penetrationTests, err := ptrs.repo.GetReportPenetrationTests(queryParams, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(penetrationTests) == 0 {
		tx.Rollback()
		return 0, utility.NewValidationError(errors.New("no penetration tests found for the specified criteria"), "empty_result", "filters")
	}

	report := &models.PenetrationTestReport{
		Title:       dto.Title,
		StartDate:   dates.StartDate,
		EndDate:     dates.EndDate,
		CreatedBYID: currentUserID,
	}

	reportID, err := ptrs.repo.Create(report, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"executive-summary", "scope", "objectives", "overview",
		"details", "findings", "retest-results", "management-sharing", "closure",
	}

	for i, s := range sectionTypes {
		en, ar := formatTitle(s)

		sec := models.PenetrationTestReportSection{
			PenetrationTestReportID: reportID,
			Type:                    s,
			TitleEn:                 en,
			TitleAr:                 ar,
			OrderNum:                i + 1,
		}

		if _, err := ptrs.repo.CreateSection(&sec, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	penetrationTestDocumentsMap := make(map[uint][]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO)
	retestDocumentsMap := make(map[uint][]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO)
	penetrationTestRetestsMap := make(map[uint][]penetrationTestDTOs.PenetrationTestReportRetestSourceDTO)
	storagePath := config.GetEnv("STORAGE_PATH", "")
	for _, penetrationTest := range penetrationTests {
		reportPenetrationTest := &models.PenetrationTestReportPenetrationTest{
			PenetrationTestReportID:              reportID,
			PenetrationTestID:                    penetrationTest.ID,
			ReferenceNumber:                      penetrationTest.ReferenceNumber,
			SystemApplicationName:                penetrationTest.SystemApplicationName,
			Scope:                                penetrationTest.Scope,
			Year:                                 penetrationTest.Year,
			Period:                               penetrationTest.Period,
			PenetrationTestTypeTitleEn:           penetrationTest.PenetrationTestTypeTitleEn,
			PenetrationTestTypeTitleAr:           penetrationTest.PenetrationTestTypeTitleAr,
			PenetrationTestPlanTitleEn:           penetrationTest.PenetrationTestPlanTitleEn,
			PenetrationTestPlanTitleAr:           penetrationTest.PenetrationTestPlanTitleAr,
			PenetrationTestClassificationTitleEn: penetrationTest.PenetrationTestClassificationTitleEn,
			PenetrationTestClassificationTitleAr: penetrationTest.PenetrationTestClassificationTitleAr,
			PenetrationTestExecutorTitleEn:       penetrationTest.PenetrationTestExecutorTitleEn,
			PenetrationTestExecutorTitleAr:       penetrationTest.PenetrationTestExecutorTitleAr,
			CSResponsibleEmployeeNameEn:          penetrationTest.CSResponsibleEmployeeNameEn,
			CSResponsibleEmployeeNameAr:          penetrationTest.CSResponsibleEmployeeNameAr,
			PenetrationTestStatusTitleEn:         penetrationTest.PenetrationTestStatusTitleEn,
			PenetrationTestStatusTitleAr:         penetrationTest.PenetrationTestStatusTitleAr,
			PenetrationTestStatusColorCode:       penetrationTest.PenetrationTestStatusColorCode,
			EnvironmentReadinessTitleEn:          penetrationTest.EnvironmentReadinessTitleEn,
			EnvironmentReadinessTitleAr:          penetrationTest.EnvironmentReadinessTitleAr,
			EnvironmentNotifyDate:                penetrationTest.EnvironmentNotifyDate,
			ExpectedReadyDate:                    penetrationTest.ExpectedReadyDate,
			ExecutionStartDate:                   penetrationTest.ExecutionStartDate,
			ExecutionEndDate:                     penetrationTest.ExecutionEndDate,
			ExecutionStatusTitleEn:               penetrationTest.ExecutionStatusTitleEn,
			ExecutionStatusTitleAr:               penetrationTest.ExecutionStatusTitleAr,
			SecurityResultTitleEn:                penetrationTest.SecurityResultTitleEn,
			SecurityResultTitleAr:                penetrationTest.SecurityResultTitleAr,
			ResultsSharedWithManagement:          penetrationTest.ResultsSharedWithManagement,
			ManagementNotifyDate:                 penetrationTest.ManagementNotifyDate,
			ManagementResponseDate:               penetrationTest.ManagementResponseDate,
			RemediationNotes:                     penetrationTest.RemediationNotes,
			ClosureDate:                          penetrationTest.ClosureDate,
			ClosureRemarks:                       penetrationTest.ClosureRemarks,
			ClosureOwnerNameEn:                   penetrationTest.ClosureOwnerNameEn,
			ClosureOwnerNameAr:                   penetrationTest.ClosureOwnerNameAr,
			CreatedAt:                            penetrationTest.CreatedAt,
			UpdatedAt:                            penetrationTest.UpdatedAt,
		}
		if err := ptrs.repo.CreateReportPenetrationTest(reportPenetrationTest, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		findings, err := ptrs.repo.GetFindingsByPenetrationTestID(penetrationTest.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		findingIDs := make(map[uint]uint)
		for _, finding := range findings {
			reportFinding := &models.PenetrationTestReportFinding{
				PenetrationTestReportID:                reportID,
				PenetrationTestReportPenetrationTestID: reportPenetrationTest.ID,
				PenetrationTestFindingID:               finding.ID,
				ReferenceNumber:                        finding.ReferenceNumber,
				Title:                                  finding.Title,
				Description:                            finding.Description,
				Recommendation:                         finding.Recommendation,
				SeverityLevelTitleEn:                   finding.SeverityLevelTitleEn,
				SeverityLevelTitleAr:                   finding.SeverityLevelTitleAr,
				SeverityLevelColorCode:                 finding.SeverityLevelColorCode,
				OwnerNameEn:                            finding.OwnerNameEn,
				OwnerNameAr:                            finding.OwnerNameAr,
				ActionPlanStatusTitleEn:                finding.ActionPlanStatusTitleEn,
				ActionPlanStatusTitleAr:                finding.ActionPlanStatusTitleAr,
				ActionPlanStatusColorCode:              finding.ActionPlanStatusColorCode,
				Percentage:                             finding.Percentage,
				TargetRemediationDate:                  finding.TargetRemediationDate,
				CreatedAt:                              finding.CreatedAt,
				UpdatedAt:                              finding.UpdatedAt,
			}
			if err := ptrs.repo.CreateReportFinding(reportFinding, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
			findingIDs[finding.ID] = reportFinding.ID
		}
		documents, err := ptrs.repo.GetDocumentsByPenetrationTestID(penetrationTest.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		penetrationTestDocumentsMap[penetrationTest.ID] = documents
		var reportDocs []*models.PenetrationTestReportDocument
		for _, document := range documents {
			if document.Name == nil || *document.Name == "" || document.PenetrationTestSupportingDocumentID == nil {
				continue
			}
			oldFilePath := filepath.Join(storagePath, "penetration-testing", strconv.FormatUint(uint64(penetrationTest.ID), 10), "supporting-documents", strconv.FormatUint(uint64(*document.PenetrationTestSupportingDocumentID), 10), "documents", *document.Name)
			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				continue
			}
			url := utility.GenerateDocumentURL("penetration-testing", "penetration-test-report", strconv.FormatUint(uint64(reportID), 10), "documents")
			reportDocs = append(reportDocs, &models.PenetrationTestReportDocument{
				PenetrationTestReportPenetrationTestID:      reportPenetrationTest.ID,
				PenetrationTestSupportingDocumentDocumentID: &document.ID,
				URL:          &url,
				Name:         document.Name,
				Size:         document.Size,
				Token:        document.Token,
				OriginalName: document.OriginalName,
				Type:         document.Type,
				CreatedBYID:  currentUserID,
			})

		}
		if len(reportDocs) > 0 {
			if err := ptrs.repo.CreatePenetrationTestDocuments(reportDocs, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
		retests, err := ptrs.repo.GetRetestsByPenetrationTestID(penetrationTest.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		penetrationTestRetestsMap[penetrationTest.ID] = retests
		for _, retest := range retests {
			reportRetest := &models.PenetrationTestReportRetest{
				PenetrationTestReportID:                reportID,
				PenetrationTestReportPenetrationTestID: reportPenetrationTest.ID,
				PenetrationTestRetestID:                retest.ID,
				ReferenceNumber:                        retest.ReferenceNumber,
				StartDate:                              retest.StartDate,
				EndDate:                                retest.EndDate,
				ExecutorTitleEn:                        retest.ExecutorTitleEn,
				ExecutorTitleAr:                        retest.ExecutorTitleAr,
				StatusTitleEn:                          retest.StatusTitleEn,
				StatusTitleAr:                          retest.StatusTitleAr,
				StatusColorCode:                        retest.StatusColorCode,
				Notes:                                  retest.Notes,
				CreatedAt:                              retest.CreatedAt,
				UpdatedAt:                              retest.UpdatedAt,
			}
			if err := ptrs.repo.CreateReportRetest(reportRetest, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
			retestFindings, err := ptrs.repo.GetRetestFindingsByRetestID(retest.ID, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			for _, retestFinding := range retestFindings {
				reportFindingID, exists := findingIDs[retestFinding.PenetrationTestFindingID]
				if !exists {
					tx.Rollback()
					return 0, utility.NewResourceNotFoundError(errors.New("retest finding does not belong to this penetration test"), retestFinding.PenetrationTestFindingID)
				}
				reportRetestFinding := &models.PenetrationTestReportRetestFinding{
					PenetrationTestReportRetestID:  reportRetest.ID,
					PenetrationTestReportFindingID: reportFindingID,
					PenetrationTestRetestFindingID: retestFinding.ID,
					Comment:                        retestFinding.Comment,
					CommentAt:                      retestFinding.CommentAt,
					ResultTitleEn:                  retestFinding.ResultTitleEn,
					ResultTitleAr:                  retestFinding.ResultTitleAr,
					ResultAt:                       retestFinding.ResultAt,
				}
				if err := ptrs.repo.CreateReportRetestFinding(reportRetestFinding, tx); err != nil {
					tx.Rollback()
					return 0, err
				}
				documents, err := ptrs.repo.GetDocumentsByRetestFindingID(retestFinding.ID, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
				retestDocumentsMap[retest.ID] = append(retestDocumentsMap[retest.ID], documents...)
				var reportDocs []*models.PenetrationTestReportDocument
				for _, document := range documents {
					if document.Name == nil || *document.Name == "" {
						continue
					}
					oldFilePath := filepath.Join(storagePath, "penetration-testing", "penetration-test-retests", strconv.FormatUint(uint64(retest.ID), 10), "findings", strconv.FormatUint(uint64(retestFinding.ID), 10), "documents", *document.Name)
					if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
						continue
					}
					url := utility.GenerateDocumentURL("penetration-testing", "penetration-test-report", strconv.FormatUint(uint64(reportID), 10), "documents")
					reportDocs = append(reportDocs, &models.PenetrationTestReportDocument{
						PenetrationTestReportPenetrationTestID: reportPenetrationTest.ID,
						PenetrationTestReportRetestFindingID:   &reportRetestFinding.ID,
						PenetrationTestRetestFindingDocumentID: &document.ID,
						URL:                                    &url,
						Name:                                   document.Name,
						Size:                                   document.Size,
						Token:                                  document.Token,
						OriginalName:                           document.OriginalName,
						Type:                                   document.Type,
						CreatedBYID:                            currentUserID,
					})

				}
				if len(reportDocs) > 0 {
					if err := ptrs.repo.CreatePenetrationTestDocuments(reportDocs, tx); err != nil {
						tx.Rollback()
						return 0, err
					}
				}
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	newStoragePath := filepath.Join(storagePath, "penetration-testing", "penetration-test-report", strconv.FormatUint(uint64(reportID), 10), "documents")
	if err := os.MkdirAll(newStoragePath, os.ModePerm); err != nil {
		_ = ptrs.DeletePenetrationTestReport(reportID, currentUserID)
		return 0, fmt.Errorf("failed to create storage directory: %w", err)
	}
	for _, penetrationTest := range penetrationTests {
		documents := penetrationTestDocumentsMap[penetrationTest.ID]
		for _, document := range documents {
			if document.Name == nil || *document.Name == "" || document.PenetrationTestSupportingDocumentID == nil {
				continue
			}
			oldFilePath := filepath.Join(storagePath, "penetration-testing", strconv.FormatUint(uint64(penetrationTest.ID), 10), "supporting-documents", strconv.FormatUint(uint64(*document.PenetrationTestSupportingDocumentID), 10), "documents", *document.Name)
			if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
				continue
			}
			input, err := os.Open(oldFilePath)
			if err != nil {
				fmt.Printf("Warning: Failed to open source document, skipping: %s - %v\n", oldFilePath, err)
				continue
			}
			newFilePath := filepath.Join(newStoragePath, *document.Name)
			output, err := os.Create(newFilePath)
			if err != nil {
				input.Close()
				_ = ptrs.DeletePenetrationTestReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to create destination file %s: %w", newFilePath, err)
			}
			_, err = io.Copy(output, input)
			output.Close()
			input.Close()
			if err != nil {
				_ = ptrs.DeletePenetrationTestReport(reportID, currentUserID)
				return 0, fmt.Errorf("failed to copy file from %s to %s: %w", oldFilePath, newFilePath, err)
			}
		}
	}
	for _, penetrationTest := range penetrationTests {
		for _, retest := range penetrationTestRetestsMap[penetrationTest.ID] {
			documents := retestDocumentsMap[retest.ID]
			for _, document := range documents {
				if document.Name == nil || *document.Name == "" || document.PenetrationTestRetestFindingID == nil {
					continue
				}
				oldFilePath := filepath.Join(storagePath, "penetration-testing", "penetration-test-retests", strconv.FormatUint(uint64(retest.ID), 10), "findings", strconv.FormatUint(uint64(*document.PenetrationTestRetestFindingID), 10), "documents", *document.Name)
				if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
					continue
				}
				input, err := os.Open(oldFilePath)
				if err != nil {
					fmt.Printf("Warning: Failed to open source document, skipping: %s - %v\n", oldFilePath, err)
					continue
				}
				newFilePath := filepath.Join(newStoragePath, *document.Name)
				output, err := os.Create(newFilePath)
				if err != nil {
					input.Close()
					_ = ptrs.DeletePenetrationTestReport(reportID, currentUserID)
					return 0, fmt.Errorf("failed to create destination file %s: %w", newFilePath, err)
				}
				_, err = io.Copy(output, input)
				output.Close()
				input.Close()
				if err != nil {
					_ = ptrs.DeletePenetrationTestReport(reportID, currentUserID)
					return 0, fmt.Errorf("failed to copy file from %s to %s: %w", oldFilePath, newFilePath, err)
				}
			}
		}
	}
	go func() {
		report, err := ptrs.repo.GetPenetrationTestReportByID(reportID)
		if err != nil {
			utility.LogError("Failed to get penetration test report", map[string]interface{}{"error": err.Error()})
			return
		}
		err = ptrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestReport{}.TableName(),
			reportID,
			report,
			currentUserID,
			"penetration_testing", "penetration_testing_report", "report", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return reportID, nil
}

func formatTitle(sectionType string) (string, string) {
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
	case "findings":
		return "Findings", "النتائج"
	case "retest-results":
		return "Retest Results", "نتائج إعادة الاختبار"
	case "management-sharing":
		return "Management Sharing", "مشاركة الإدارة"
	case "closure":
		return "Closure", "الإغلاق"
	default:
		return sectionType, sectionType
	}
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestReportDetails(id uint) (*penetrationTestDTOs.PenetrationTestReportDetailsDTO, error) {
	report, err := ptrs.repo.GetPenetrationTestReportDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	var response penetrationTestDTOs.PenetrationTestReportDetailsDTO
	if err := mapstructure.Decode(report, &response); err != nil {
		return nil, err
	}
	response.CreatedAt = report.CreatedAt
	response.StartDate = report.StartDate
	response.EndDate = report.EndDate
	for i, penetrationTest := range response.PenetrationTests {
		penetrationTest.EnvironmentNotifyDate = report.PenetrationTests[i].EnvironmentNotifyDate
		penetrationTest.ExpectedReadyDate = report.PenetrationTests[i].ExpectedReadyDate
		penetrationTest.ExecutionStartDate = report.PenetrationTests[i].ExecutionStartDate
		penetrationTest.ExecutionEndDate = report.PenetrationTests[i].ExecutionEndDate
		penetrationTest.ManagementNotifyDate = report.PenetrationTests[i].ManagementNotifyDate
		penetrationTest.ManagementResponseDate = report.PenetrationTests[i].ManagementResponseDate
		penetrationTest.ClosureDate = report.PenetrationTests[i].ClosureDate
		penetrationTest.CreatedAt = &report.PenetrationTests[i].CreatedAt
		penetrationTest.UpdatedAt = report.PenetrationTests[i].UpdatedAt

		for j, finding := range penetrationTest.Findings {
			finding.TargetRemediationDate = report.PenetrationTests[i].Findings[j].TargetRemediationDate
			finding.CreatedAt = &report.PenetrationTests[i].Findings[j].CreatedAt
			finding.UpdatedAt = report.PenetrationTests[i].Findings[j].UpdatedAt
		}
		for j, retest := range penetrationTest.Retests {
			retest.StartDate = report.PenetrationTests[i].Retests[j].StartDate
			retest.EndDate = report.PenetrationTests[i].Retests[j].EndDate
			retest.CreatedAt = &report.PenetrationTests[i].Retests[j].CreatedAt
			retest.UpdatedAt = report.PenetrationTests[i].Retests[j].UpdatedAt

			for k, finding := range retest.RetestFindings {
				finding.CommentAt = report.PenetrationTests[i].Retests[j].RetestFindings[k].CommentAt
				finding.ResultAt = report.PenetrationTests[i].Retests[j].RetestFindings[k].ResultAt

			}
		}

		if penetrationTest.PenetrationTest != nil {
			penetrationTest.Closure = penetrationTest.PenetrationTest.Closure
			penetrationTest.PenetrationTest = nil
		}
		if penetrationTest.Closure != nil && report.PenetrationTests[i].PenetrationTest != nil && report.PenetrationTests[i].PenetrationTest.Closure != nil {
			closureDate := report.PenetrationTests[i].PenetrationTest.Closure.ClosureDate.Format("2006-01-02")
			penetrationTest.Closure.ClosureDateString = &closureDate
		}
	}
	return &response, nil
}

func (ptrs *PenetrationTestReportService) DeletePenetrationTestReport(penetrationTestReportID uint, currentUserID uint) error {
	tx := ptrs.repo.BeginTransaction()

	existingReport, err := ptrs.repo.GetPenetrationTestReportByID(penetrationTestReportID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, penetrationTestReportID)
		}
		return err
	}

	var documentsToDelete []*models.PenetrationTestReportDocument
	db := tx.DB()
	if err := db.Where("penetration_test_report_penetration_test_id IN (SELECT id FROM penetration_test_report_penetration_tests WHERE penetration_test_report_id = ?)", penetrationTestReportID).Find(&documentsToDelete).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportDocuments(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportRetestFindings(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportRetests(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportFindings(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportPenetrationTests(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReportSections(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ptrs.repo.DeletePenetrationTestReport(penetrationTestReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(documentsToDelete) > 0 {
		storagePath := filepath.Join("penetration-testing", "penetration-test-report", strconv.FormatUint(uint64(penetrationTestReportID), 10), "documents")
		for _, doc := range documentsToDelete {
			if doc.Name != nil && *doc.Name != "" {
				if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
					utility.LogError("Failed to delete file", map[string]interface{}{"file": *doc.Name, "error": err.Error()})
				}
			}
		}
	}

	go func() {
		err := ptrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestReport{}.TableName(),
			penetrationTestReportID,
			existingReport,
			currentUserID,
			"penetration_testing", "penetration_testing_report", "report", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (ptrs *PenetrationTestReportService) GetFullReportDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestReportPDFDetailDTO, error) {
	penetrationTests, err := ptrs.repo.GetPenetrationTestsByIDs([]uint{id}, nil)
	if err != nil {
		return nil, err
	}
	if len(penetrationTests) == 0 {
		return nil, utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, id)
	}
	var penetrationTest penetrationTestDTOs.PenetrationTestReportPenetrationTestDTO
	if err := mapstructure.Decode(penetrationTests[0], &penetrationTest); err != nil {
		return nil, err
	}
	penetrationTest.EnvironmentNotifyDate = penetrationTests[0].EnvironmentNotifyDate
	penetrationTest.ExpectedReadyDate = penetrationTests[0].ExpectedReadyDate
	penetrationTest.ExecutionStartDate = penetrationTests[0].ExecutionStartDate
	penetrationTest.ExecutionEndDate = penetrationTests[0].ExecutionEndDate
	penetrationTest.ManagementNotifyDate = penetrationTests[0].ManagementNotifyDate
	penetrationTest.ManagementResponseDate = penetrationTests[0].ManagementResponseDate
	penetrationTest.ClosureDate = penetrationTests[0].ClosureDate
	penetrationTest.CreatedAt = &penetrationTests[0].CreatedAt
	penetrationTest.UpdatedAt = penetrationTests[0].UpdatedAt
	findings, err := ptrs.repo.GetFindingsByPenetrationTestID(id, nil)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(findings, &penetrationTest.Findings); err != nil {
		return nil, err
	}
	retests, err := ptrs.repo.GetRetestsByPenetrationTestID(id, nil)
	if err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(retests, &penetrationTest.Retests); err != nil {
		return nil, err
	}
	for i, finding := range penetrationTest.Findings {
		finding.TargetRemediationDate = findings[i].TargetRemediationDate
		finding.CreatedAt = &findings[i].CreatedAt
		finding.UpdatedAt = findings[i].UpdatedAt
	}
	for i, retest := range penetrationTest.Retests {
		retest.StartDate = retests[i].StartDate
		retest.EndDate = retests[i].EndDate
		retest.CreatedAt = &retests[i].CreatedAt
		retest.UpdatedAt = retests[i].UpdatedAt

		findings, err := ptrs.repo.GetRetestFindingsByRetestID(retest.ID, nil)
		if err != nil {
			return nil, err
		}
		if err := mapstructure.Decode(findings, &retest.RetestFindings); err != nil {
			return nil, err
		}
		for j, finding := range retest.RetestFindings {
			finding.PenetrationTestReportFindingID = findings[j].PenetrationTestFindingID
			finding.CommentAt = findings[j].CommentAt
			finding.ResultAt = findings[j].ResultAt
		}
	}
	source, err := ptrs.repo.GetFullReportDetailsByID(id)
	if err != nil {
		return nil, err
	}
	report := &penetrationTestDTOs.PenetrationTestReportPDFDetailDTO{
		ID:               id,
		ReferenceNumber:  penetrationTest.ReferenceNumber,
		Title:            penetrationTest.SystemApplicationName,
		PenetrationTests: []*penetrationTestDTOs.PenetrationTestReportPenetrationTestDTO{&penetrationTest},
	}
	if source.CreatedByUser != nil {
		report.CreatedUserFirstNameEn = source.CreatedByUser.FirstNameEn
		report.CreatedUserLastNameEn = source.CreatedByUser.LastNameEn
	}
	if penetrationTest.CreatedAt != nil {
		report.ReportedAt = penetrationTest.CreatedAt.Format("Jan _2, 2006, 03:04 PM")
	}
	for i, sectionType := range []string{"executive-summary", "scope", "objectives", "overview", "details", "findings", "retest-results", "management-sharing", "closure"} {
		en, ar := formatTitle(sectionType)
		report.Sections = append(report.Sections, &penetrationTestDTOs.PenetrationTestReportSectionDTO{
			Type:     sectionType,
			TitleEn:  en,
			TitleAr:  ar,
			OrderNum: i + 1,
		})
	}
	return report, nil
}

func (ptrs *PenetrationTestReportService) ExportPenetrationTestReport(id uint, currentUserID uint) (string, error) {
	reportDetail, err := ptrs.GetFullReportDetailsByID(id)
	if err != nil {
		return "", err
	}
	reportDetail.GeneratedAt = time.Now().Format("Jan _2, 2006, 03:04 PM")
	htmlContent, err := utility.GeneratePDFContent(reportDetail, []string{"penetration-testing/templates/pdf/penetration_test_direct_report_body.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate body content: %w", err)
	}
	headerContent, err := utility.GeneratePDFContent(reportDetail, []string{"penetration-testing/templates/pdf/penetration_test_direct_report_header.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}
	footerContent, err := utility.GeneratePDFContent(reportDetail, []string{"penetration-testing/templates/pdf/penetration_test_direct_report_footer.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	coverContent, err := utility.GeneratePDFContent(reportDetail, []string{"penetration-testing/templates/pdf/penetration_test_direct_report_cover.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}
	pdfFileName := fmt.Sprintf("PenetrationTest_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}
	return pdfFileName, nil
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestReportExportZip(penetrationTestReportID uint, currentUserID uint, language string) (string, error) {
	report, err := ptrs.GetPenetrationTestReportDetails(penetrationTestReportID)
	if err != nil {
		return "", err
	}
	reportDetail := penetrationTestDTOs.PenetrationTestReportPDFDetailDTO{
		ID:               report.ID,
		Title:            report.Title,
		StartDate:        report.StartDate,
		EndDate:          report.EndDate,
		ReportedAt:       report.CreatedAt.Format("02,Jan 2006"),
		GeneratedAt:      time.Now().Format("Jan _2, 2006, 03:04 PM"),
		Sections:         report.Sections,
		PenetrationTests: report.PenetrationTests,
	}
	if report.CreatedByUser != nil {
		reportDetail.CreatedUserFirstNameEn = report.CreatedByUser.FirstNameEn
		reportDetail.CreatedUserLastNameEn = report.CreatedByUser.LastNameEn
		reportDetail.CreatedUserFirstNameAr = report.CreatedByUser.FirstNameAr
		reportDetail.CreatedUserLastNameAr = report.CreatedByUser.LastNameAr
	}
	bodyTemplates := []string{"penetration-testing/templates/pdf/penetration_test_report_body.html"}
	headerTemplates := []string{"penetration-testing/templates/pdf/penetration_test_report_header.html"}
	footerTemplates := []string{"penetration-testing/templates/pdf/penetration_test_report_footer.html"}
	coverTemplates := []string{"penetration-testing/templates/pdf/penetration_test_report_cover.html"}
	if language == "ar" {
		bodyTemplates = []string{"penetration-testing/templates/pdf/arabic/penetration_test_report_body_ar.html"}
		headerTemplates = []string{"penetration-testing/templates/pdf/arabic/penetration_test_report_header_ar.html"}
		footerTemplates = []string{"penetration-testing/templates/pdf/arabic/penetration_test_report_footer_ar.html"}
		coverTemplates = []string{"penetration-testing/templates/pdf/arabic/penetration_test_report_cover_ar.html"}
		reportDetail.ReportedAt = report.CreatedAt.Format("2006-01-02")
		reportDetail.GeneratedAt = time.Now().Format("2006-01-02 15:04")
	}
	htmlContent, err := utility.GeneratePDFContent(reportDetail, bodyTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate body content: %w", err)
	}
	headerContent, err := utility.GeneratePDFContent(reportDetail, headerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}
	footerContent, err := utility.GeneratePDFContent(reportDetail, footerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}
	coverContent, err := utility.GeneratePDFContent(reportDetail, coverTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}
	pdfFileName := fmt.Sprintf("PenetrationTestReport_%d.pdf", penetrationTestReportID)
	if language == "ar" {
		pdfFileName = fmt.Sprintf("PenetrationTestReport_%d_ar.pdf", penetrationTestReportID)
	}
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	filesToZip := []string{pdfFileName}
	storagePath := config.GetEnv("STORAGE_PATH", "")
	baseDocPath := filepath.Join(storagePath, "penetration-testing", "penetration-test-report", strconv.FormatUint(uint64(penetrationTestReportID), 10), "documents")
	for _, penetrationTest := range report.PenetrationTests {
		for _, doc := range penetrationTest.Documents {
			if doc.Token != nil {
				decryptedFilename, err := utility.DecryptData(*doc.Token)
				if err != nil {
					fmt.Printf("Error decrypting token for doc ID %d: %v\n", doc.ID, err)
					continue
				}
				fullPath := filepath.Join(baseDocPath, string(decryptedFilename))
				if _, err := os.Stat(fullPath); err == nil {
					filesToZip = append(filesToZip, fullPath)
				}
			}
		}
	}
	zipDir := filepath.Join("storage", "zip")
	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("penetration_test_report_%d.zip", penetrationTestReportID))
	if language == "ar" {
		zipFilePath = filepath.Join(zipDir, fmt.Sprintf("penetration_test_report_%d_ar.zip", penetrationTestReportID))
	}
	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", zipDir, err)
	}
	if err := utility.CreateZipFile(zipFilePath, filesToZip); err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}
	go func() {
		time.Sleep(2 * time.Second)
		if err := os.Remove(pdfFileName); err != nil {
			utility.LogError("Failed to delete temp file", map[string]interface{}{"file": pdfFileName, "error": err.Error()})
		}
	}()
	return zipFilePath, nil
}
