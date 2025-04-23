package services

import (
	"fmt"
	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ComplianceReportService struct {
	repo               repositories.ComplianceReportRepositoryInterface
	activityLogService logServices.ActivityLogServiceInterface
}

type ComplianceReportServiceInterface interface {
	GetAll(queryParams *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error)
	Create(currentUserID uint, request *dtos.ComplianceReportFormDTO) (uint, error)
	GetComplianceReportDetails(complianceReportID uint) (*dtos.ComplianceReportDetailsDTO, error)
	GetComplianceReportSectionDetails(complianceReportID uint, sectionID uint) (*dtos.ComplianceReportSection, error)
	UpdateSection(complianceReportID uint, sectionID uint, currentUserID uint, formDTO *dtos.ComplianceReportSectionFormDTO) error
	DeleteComplianceReport(complianceReportID uint, currentUserID uint) error
	ExportComplianceReportToPDF(complianceReportID uint, currentUserID uint) (string, error)
}

func NewComplianceReportService(
	repo repositories.ComplianceReportRepositoryInterface,
	activityLogService logServices.ActivityLogServiceInterface,
) *ComplianceReportService {
	return &ComplianceReportService{
		repo:               repo,
		activityLogService: activityLogService,
	}
}

func (crs *ComplianceReportService) GetAll(queryParams *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error) {
	reports, err := crs.repo.GetAll(queryParams)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (crs *ComplianceReportService) calculateAverageScoreAndStatus(requirement *models.ComplianceRequirement, assessments []models.ComplianceAssessment, orgUnitAssessments map[uint][]models.ComplianceAssessment) (*float64, *uint, error) {
	var totalScore float64
	var scoreCount int

	// Collect all scores for this requirement across all department assessments
	for _, deptAssessments := range orgUnitAssessments {
		for _, assessment := range deptAssessments {
			for _, otherRequirement := range assessment.AssessmentRequirements {
				if otherRequirement.ComplianceRequirement.ID == requirement.ID &&
					otherRequirement.Score != nil &&
					otherRequirement.ComplianceStatus != nil &&
					otherRequirement.ComplianceStatus.Type != "not_applicable" {
					totalScore += float64(*otherRequirement.Score)
					scoreCount++
				}
			}
		}
	}

	var averageScore *float64
	var averageComplianceStatusID *uint

	if scoreCount > 0 {
		avg := totalScore / float64(scoreCount)
		averageScore = &avg

		// Get compliance status based on average score
		statuses, err := crs.repo.GetAllComplianceStatuses()
		if err != nil {
			return nil, nil, err
		}

		for _, status := range statuses {
			if status.ScoreFrom != nil && status.ScoreTo != nil {
				if avg >= float64(*status.ScoreFrom) && avg <= float64(*status.ScoreTo) {
					averageComplianceStatusID = &status.ID
					break
				}
			}
		}
	}

	return averageScore, averageComplianceStatusID, nil
}

func (crs *ComplianceReportService) Create(currentUserID uint, request *dtos.ComplianceReportFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	parsedDates, err := crs.parseComplianceReportDates(request)
	if err != nil {
		return 0, err
	}

	complianceReport := &models.ComplianceReport{
		Title:                request.Title,
		StartDate:            *parsedDates.StartDate,
		EndDate:              *parsedDates.EndDate,
		ComplianceRegisterID: request.ComplianceRegisterID,
		CreatedBYID:          currentUserID,
	}

	complianceReportID, err := crs.repo.Create(complianceReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sectionTypes := []string{
		"executive-summary", "compliance-status-by-department", "compliance-status-by-standard",
		"action-plan-implementation", "conclusion-recommendations",
	}

	for i, sectionType := range sectionTypes {
		titleEn, titleAr := formatTitle(sectionType)
		section := &models.ComplianceReportSection{
			ComplianceReportID: complianceReportID,
			Type:               sectionType,
			TitleEn:            titleEn,
			TitleAr:            titleAr,
			DescriptionEn:      nil,
			DescriptionAr:      nil,
			OrderNum:           i + 1,
		}

		_, err := crs.repo.CreateSection(section, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	complianceAssessments, err := crs.repo.GetComplianceAssessmentByComplianceRegisterID(request.ComplianceRegisterID, *parsedDates.StartDate, *parsedDates.EndDate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Group assessments by organization unit
	orgUnitAssessments := make(map[uint][]models.ComplianceAssessment)
	for _, assessment := range complianceAssessments {
		orgUnitAssessments[assessment.PrimaryUnitID] = append(orgUnitAssessments[assessment.PrimaryUnitID], assessment)
	}

	// Create report department for each organization unit
	for orgUnitID, assessments := range orgUnitAssessments {
		// Calculate average score for this organization unit's assessments
		averageScore := calculateAverageScore(assessments)

		reportDepartment := &models.ComplianceReportDepartment{
			ComplianceReportID:   complianceReportID,
			PrimaryUnitID:        orgUnitID,
			CompliancePercentage: averageScore,
			CreatedAt:            time.Now(),
			CreatedBYID:          currentUserID,
		}

		reportDepartmentID, err := crs.repo.CreateComplianceReportDepartment(reportDepartment, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Create compliance report clauses first
		clauseMap := make(map[uint]*models.ComplianceReportClause)
		for _, assessment := range assessments {
			for _, requirementAssessment := range assessment.AssessmentRequirements {
				if _, exists := clauseMap[requirementAssessment.ComplianceRequirement.ComplianceClause.ID]; !exists {
					complianceReportClause := &models.ComplianceReportClause{
						ComplianceReportID:           complianceReportID,
						ComplianceReportDepartmentID: reportDepartmentID,
						ClauseNo:                     requirementAssessment.ComplianceRequirement.ComplianceClause.ClauseNo,
						ClauseTitleEn:                requirementAssessment.ComplianceRequirement.ComplianceClause.TitleEn,
						ClauseTitleAr:                requirementAssessment.ComplianceRequirement.ComplianceClause.TitleAr,
						DescriptionEn:                requirementAssessment.ComplianceRequirement.ComplianceClause.DescriptionEn,
						DescriptionAr:                requirementAssessment.ComplianceRequirement.ComplianceClause.DescriptionAr,
						CreatedBYID:                  currentUserID,
					}

					clauseID, err := crs.repo.CreateComplianceReportClause(complianceReportClause, tx)
					if err != nil {
						tx.Rollback()
						return 0, err
					}
					clauseMap[requirementAssessment.ComplianceRequirement.ComplianceClause.ID] = &models.ComplianceReportClause{ID: clauseID}
				}

			}
		}

		// Now create requirements
		for _, assessment := range assessments {
			for _, requirement := range assessment.AssessmentRequirements {
				averageScore, averageComplianceStatusID, err := crs.calculateAverageScoreAndStatus(requirement.ComplianceRequirement, assessments, orgUnitAssessments)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				complianceReportRequirement := &models.ComplianceReportRequirement{
					ComplianceReportID:        complianceReportID,
					ComplianceReportClauseID:  clauseMap[requirement.ComplianceRequirement.ComplianceClause.ID].ID,
					TitleEn:                   requirement.ComplianceRequirement.TitleEn,
					TitleAr:                   requirement.ComplianceRequirement.TitleAr,
					Score:                     requirement.Score,
					AverageScore:              averageScore,
					ComplianceStatusID:        nil,
					AverageComplianceStatusID: averageComplianceStatusID,
					CreatedBYID:               currentUserID,
				}

				if requirement.ComplianceStatus != nil {
					complianceReportRequirement.ComplianceStatusID = &requirement.ComplianceStatus.ID
				}

				complianceReportRequirementID, err := crs.repo.CreateComplianceReportRequirement(complianceReportRequirement, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				for _, actionPlan := range requirement.ActionPlans {

					complianceReportActionPlan := &models.ComplianceReportActionPlan{
						ComplianceReportID:            complianceReportID,
						ComplianceReportRequirementID: complianceReportRequirementID,
						Title:                         actionPlan.Title,
						Description:                   actionPlan.Description,
						ActionPlanStatusID:            actionPlan.ActionPlanStatusID,
						TargetDate:                    actionPlan.TargetDate,
						SeverityLevelID:               actionPlan.SeverityLevelID,
						CreatedAt:                     time.Now(),
						CreatedBYID:                   currentUserID,
					}

					_, err := crs.repo.CreateComplianceReportActionPlan(complianceReportActionPlan, tx)
					if err != nil {
						tx.Rollback()
						return 0, err
					}
				}
			}
		}

	}

	tx.Commit()

	return complianceReportID, nil
}

func (crs *ComplianceReportService) GetComplianceReportDetails(complianceReportID uint) (*dtos.ComplianceReportDetailsDTO, error) {
	var complianceReportDetails dtos.ComplianceReportDetailsDTO

	complianceReport, err := crs.repo.GetComplianceReportByID(complianceReportID)
	if err != nil {
		return nil, err
	}

	departments, err := crs.repo.GetComplianceReportDepartments(complianceReportID)
	if err != nil {
		return nil, err
	}

	// Map basic report details
	complianceReportDetails.ID = complianceReport.ID
	complianceReportDetails.Title = complianceReport.Title
	complianceReportDetails.StartDate = complianceReport.StartDate.Format("2006-01-02")
	complianceReportDetails.EndDate = complianceReport.EndDate.Format("2006-01-02")
	complianceReportDetails.ComplianceRegisterID = complianceReport.ComplianceRegisterID
	complianceReportDetails.CreatedAt = complianceReport.CreatedAt.Format("2006-01-02")

	// Map compliance register
	complianceReportDetails.ComplianceRegister.ID = complianceReport.ComplianceRegister.ID
	complianceReportDetails.ComplianceRegister.TitleEn = complianceReport.ComplianceRegister.TitleEn
	complianceReportDetails.ComplianceRegister.TitleAr = complianceReport.ComplianceRegister.TitleAr

	// Map created by user
	complianceReportDetails.CreatedByUser.ID = int(complianceReport.CreatedByUser.ID)
	complianceReportDetails.CreatedByUser.FirstNameEN = complianceReport.CreatedByUser.FirstNameEn
	complianceReportDetails.CreatedByUser.FirstNameAR = complianceReport.CreatedByUser.FirstNameAr
	complianceReportDetails.CreatedByUser.LastNameEN = complianceReport.CreatedByUser.LastNameEn
	complianceReportDetails.CreatedByUser.LastNameAR = complianceReport.CreatedByUser.LastNameAr
	complianceReportDetails.CreatedByUser.Email = complianceReport.CreatedByUser.Email

	// Map sections
	for _, section := range complianceReport.Sections {
		sectionDTO := &dtos.ComplianceReportSection{
			ID:            section.ID,
			TitleEn:       section.TitleEn,
			TitleAr:       section.TitleAr,
			DescriptionEn: section.DescriptionEn,
			DescriptionAr: section.DescriptionAr,
			Type:          section.Type,
			Order:         section.OrderNum,
		}
		complianceReportDetails.ComplianceReportSections = append(complianceReportDetails.ComplianceReportSections, sectionDTO)

		if section.Type == "compliance-status-by-department" {
			for _, department := range departments {
				if department == nil {
					continue
				}
				fmt.Println("department", department.OrganizationUnit.ID)
				departmentDTO := &dtos.ComplianceReportDepartment{
					OrganizationUnit: settingsDTOs.OrganizationUnitDTO{
						ID:      int(department.OrganizationUnit.ID),
						TitleEn: department.OrganizationUnit.TitleEn,
						TitleAr: department.OrganizationUnit.TitleAr,
					},
					CompliancePercentage: &department.CompliancePercentage,
					NonCompliantAreas:    []*dtos.ComplianceReportRequirement{},
				}

				for _, clause := range department.ComplianceReportClauses {
					if clause == nil {
						continue
					}
					for _, nonCompliantRequirement := range clause.ComplianceReportNonCompliantRequirements {
						if nonCompliantRequirement == nil {
							continue
						}

						nonCompliantArea := &dtos.ComplianceReportRequirement{
							ID:          nonCompliantRequirement.ID,
							TitleEn:     nonCompliantRequirement.TitleEn,
							TitleAr:     nonCompliantRequirement.TitleAr,
							Score:       nonCompliantRequirement.Score,
							ActionPlans: []*dtos.ComplianceReportActionPlan{},
						}

						if nonCompliantRequirement.ComplianceStatus.ID != 0 {
							nonCompliantArea.ComplianceStatus = &dtos.ComplianceStatusDTO{
								ID:        int(nonCompliantRequirement.ComplianceStatus.ID),
								Type:      nonCompliantRequirement.ComplianceStatus.Type,
								TitleEn:   nonCompliantRequirement.ComplianceStatus.TitleEn,
								TitleAr:   nonCompliantRequirement.ComplianceStatus.TitleAr,
								ColorCode: nonCompliantRequirement.ComplianceStatus.ColorCode,
							}
						} else {
							nonCompliantArea.ComplianceStatus = nil
						}

						for _, actionPlan := range nonCompliantRequirement.ComplianceReportActionPlans {
							actionPlanDTO := &dtos.ComplianceReportActionPlan{
								Title:       actionPlan.Title,
								Description: actionPlan.Description,
								ActionPlanStatus: masterDTOs.ActionPlanStatusDTO{
									ID:        actionPlan.ActionPlanStatus.ID,
									TitleEn:   actionPlan.ActionPlanStatus.TitleEn,
									TitleAr:   actionPlan.ActionPlanStatus.TitleAr,
									ColorCode: actionPlan.ActionPlanStatus.ColorCode,
								},
								TargetDate: func() string {
									if actionPlan.TargetDate == nil {
										return ""
									}
									return actionPlan.TargetDate.Format("2006-01-02")
								}(),
							}
							if actionPlan.SeverityLevel != nil {
								actionPlanDTO.SeverityLevel = &masterDTOs.SeverityLevelDTO{
									ID:        actionPlan.SeverityLevel.ID,
									TitleEn:   actionPlan.SeverityLevel.TitleEn,
									TitleAr:   actionPlan.SeverityLevel.TitleAr,
									ColorCode: actionPlan.SeverityLevel.ColorCode,
								}
							}
							nonCompliantArea.ActionPlans = append(nonCompliantArea.ActionPlans, actionPlanDTO)
						}

						departmentDTO.NonCompliantAreas = append(departmentDTO.NonCompliantAreas, nonCompliantArea)
					}
				}

				sectionDTO.ComplianceReportDepartments = append(sectionDTO.ComplianceReportDepartments, departmentDTO)
			}

		}

		if section.Type == "compliance-status-by-standard" {
			// Create maps to store unique clauses and requirements
			uniqueClauses := make(map[string]*dtos.ComplianceReportClause)           // Use clause_no as key
			uniqueRequirements := make(map[string]*dtos.ComplianceReportRequirement) // Track unique requirements by clause and title

			for _, department := range departments {
				if department == nil {
					continue
				}

				for _, clause := range department.ComplianceReportClauses {
					if clause == nil {
						continue
					}

					// If clause is not in uniqueClauses, add it
					if _, exists := uniqueClauses[clause.ClauseNo]; !exists {
						clauseDTO := &dtos.ComplianceReportClause{
							ClauseNo:               clause.ClauseNo,
							ClauseTitleEn:          clause.ClauseTitleEn,
							ClauseTitleAr:          clause.ClauseTitleAr,
							DescriptionEn:          clause.DescriptionEn,
							DescriptionAr:          clause.DescriptionAr,
							ComplianceRequirements: []*dtos.ComplianceReportRequirement{},
						}
						uniqueClauses[clause.ClauseNo] = clauseDTO
					}

					// Add unique requirements for this clause
					for _, requirement := range clause.ComplianceReportRequirements {
						if requirement == nil {
							continue
						}

						// Create a unique key for the requirement using clause number and title
						requirementKey := fmt.Sprintf("%s_%s_%s", clause.ClauseNo, requirement.TitleEn, requirement.TitleAr)

						// Get or create requirement DTO
						var requirementDTO *dtos.ComplianceReportRequirement
						if existingReq, exists := uniqueRequirements[requirementKey]; exists {
							requirementDTO = existingReq
						} else {
							requirementDTO = &dtos.ComplianceReportRequirement{
								ID:          requirement.ID,
								TitleEn:     requirement.TitleEn,
								TitleAr:     requirement.TitleAr,
								Score:       requirement.Score,
								ActionPlans: []*dtos.ComplianceReportActionPlan{},
							}

							if requirement.AverageComplianceStatus.ID != 0 {
								requirementDTO.ComplianceStatus = &dtos.ComplianceStatusDTO{
									ID:        int(requirement.AverageComplianceStatus.ID),
									Type:      requirement.AverageComplianceStatus.Type,
									TitleEn:   requirement.AverageComplianceStatus.TitleEn,
									TitleAr:   requirement.AverageComplianceStatus.TitleAr,
									ColorCode: requirement.AverageComplianceStatus.ColorCode,
								}
							}

							uniqueRequirements[requirementKey] = requirementDTO
							uniqueClauses[clause.ClauseNo].ComplianceRequirements = append(uniqueClauses[clause.ClauseNo].ComplianceRequirements, requirementDTO)
						}

						// Add action plans for this requirement
						for _, actionPlan := range requirement.ComplianceReportActionPlans {
							if actionPlan == nil {
								continue
							}

							actionPlanDTO := &dtos.ComplianceReportActionPlan{
								Title:       actionPlan.Title,
								Description: actionPlan.Description,
								ActionPlanStatus: masterDTOs.ActionPlanStatusDTO{
									ID:        actionPlan.ActionPlanStatus.ID,
									TitleEn:   actionPlan.ActionPlanStatus.TitleEn,
									TitleAr:   actionPlan.ActionPlanStatus.TitleAr,
									ColorCode: actionPlan.ActionPlanStatus.ColorCode,
								},
								TargetDate: func() string {
									if actionPlan.TargetDate == nil {
										return ""
									}
									return actionPlan.TargetDate.Format("2006-01-02")
								}(),
							}

							if actionPlan.SeverityLevel != nil {
								actionPlanDTO.SeverityLevel = &masterDTOs.SeverityLevelDTO{
									ID:        actionPlan.SeverityLevel.ID,
									TitleEn:   actionPlan.SeverityLevel.TitleEn,
									TitleAr:   actionPlan.SeverityLevel.TitleAr,
									ColorCode: actionPlan.SeverityLevel.ColorCode,
								}
							}

							requirementDTO.ActionPlans = append(requirementDTO.ActionPlans, actionPlanDTO)
						}
					}
				}
			}

			// Add unique clauses to the section
			// Convert map to slice for sorting
			var clauseSlice []*dtos.ComplianceReportClause
			for _, clauseDTO := range uniqueClauses {
				clauseSlice = append(clauseSlice, clauseDTO)
			}

			// Sort clauses by clause number
			sort.Slice(clauseSlice, func(i, j int) bool {
				// Split clause numbers into parts (e.g., "1.2" -> ["1", "2"])
				iParts := strings.Split(clauseSlice[i].ClauseNo, ".")
				jParts := strings.Split(clauseSlice[j].ClauseNo, ".")

				// Compare each part numerically
				for k := 0; k < len(iParts) && k < len(jParts); k++ {
					// Convert parts to integers for comparison
					iNum, errI := strconv.Atoi(iParts[k])
					jNum, errJ := strconv.Atoi(jParts[k])

					// If conversion fails, fall back to string comparison
					if errI != nil || errJ != nil {
						return iParts[k] < jParts[k]
					}

					if iNum != jNum {
						return iNum < jNum
					}
				}

				// If all parts are equal up to the shortest length,
				// the shorter number comes first
				return len(iParts) < len(jParts)
			})

			// Add sorted clauses to the section
			sectionDTO.ComplianceReportStandards = clauseSlice
		}

		if section.Type == "action-plan-implementation" {
			for _, department := range departments {
				if department == nil {
					continue
				}

				for _, clause := range department.ComplianceReportClauses {
					if clause == nil {
						continue
					}

					for _, requirement := range clause.ComplianceReportRequirements {
						if requirement == nil {
							continue
						}

						for _, actionPlan := range requirement.ComplianceReportActionPlans {
							if actionPlan == nil {
								continue
							}

							actionPlanDTO := &dtos.ComplianceReportActionPlan{
								Title:       actionPlan.Title,
								Description: actionPlan.Description,
								ActionPlanStatus: masterDTOs.ActionPlanStatusDTO{
									ID:        actionPlan.ActionPlanStatus.ID,
									TitleEn:   actionPlan.ActionPlanStatus.TitleEn,
									TitleAr:   actionPlan.ActionPlanStatus.TitleAr,
									ColorCode: actionPlan.ActionPlanStatus.ColorCode,
								},
								TargetDate: func() string {
									if actionPlan.TargetDate == nil {
										return ""
									}
									return actionPlan.TargetDate.Format("2006-01-02")
								}(),
								OrganizationUnit: settingsDTOs.OrganizationUnitDTO{
									ID:      int(department.OrganizationUnit.ID),
									TitleEn: department.OrganizationUnit.TitleEn,
									TitleAr: department.OrganizationUnit.TitleAr,
								},
							}
							if actionPlan.SeverityLevel != nil {
								actionPlanDTO.SeverityLevel = &masterDTOs.SeverityLevelDTO{
									ID:        actionPlan.SeverityLevel.ID,
									TitleEn:   actionPlan.SeverityLevel.TitleEn,
									TitleAr:   actionPlan.SeverityLevel.TitleAr,
									ColorCode: actionPlan.SeverityLevel.ColorCode,
								}
							} else {
								actionPlanDTO.SeverityLevel = nil
							}
							sectionDTO.ComplianceReportActionPlans = append(sectionDTO.ComplianceReportActionPlans, actionPlanDTO)
						}
					}
				}
			}
		}
	}

	return &complianceReportDetails, nil
}

func (crs *ComplianceReportService) parseComplianceReportDates(dto *dtos.ComplianceReportFormDTO) (*dtos.ComplianceReportParsedDates, error) {
	parsed := &dtos.ComplianceReportParsedDates{}

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

func formatTitle(sectionType string) (string, string) {
	switch sectionType {
	case "executive-summary":
		return "Executive Summary", "الملخص التنفيذي"
	case "compliance-status-by-department":
		return "Compliance Status by Department", "الحالة المرافقة للمواد المنظمة"
	case "compliance-status-by-standard":
		return "Compliance Status by Standard", "الحالة المرافقة للمعايير"
	case "action-plan-implementation":
		return "Action Plan Implementation", "تنفيذ الخطة"
	case "conclusion-recommendations":
		return "Conclusion and Recommendations", "الاستنتاجات والتوصيات"
	}
	return "", ""
}

func calculateAverageScore(assessments []models.ComplianceAssessment) float64 {
	if len(assessments) == 0 {
		return 0
	}

	var totalScore float64
	var totalCount int

	for _, assessment := range assessments {
		for _, requirement := range assessment.AssessmentRequirements {
			if requirement.ComplianceStatus != nil &&
				requirement.ComplianceStatus.Type != "not_applicable" &&
				requirement.Score != nil {
				totalScore += float64(*requirement.Score)
				totalCount++
			}
		}
	}

	if totalCount == 0 {
		return 0
	}

	return totalScore / float64(totalCount)
}

func (crs *ComplianceReportService) GetComplianceReportSectionDetails(complianceReportID uint, sectionID uint) (*dtos.ComplianceReportSection, error) {
	section, err := crs.repo.GetComplianceReportSectionByID(sectionID)
	if err != nil {
		return nil, err
	}

	sectionDTO := &dtos.ComplianceReportSection{
		ID:            section.ID,
		TitleEn:       section.TitleEn,
		TitleAr:       section.TitleAr,
		DescriptionEn: section.DescriptionEn,
		DescriptionAr: section.DescriptionAr,
		Type:          section.Type,
		Order:         section.OrderNum,
	}

	return sectionDTO, nil
}

func (crs *ComplianceReportService) UpdateSection(complianceReportID uint, sectionID uint, currentUserID uint, formDTO *dtos.ComplianceReportSectionFormDTO) error {
	tx := crs.repo.BeginTransaction()

	_, err := crs.repo.GetComplianceReportByID(complianceReportID)
	if err != nil {
		return err
	}

	section, err := crs.repo.GetComplianceReportSectionByID(sectionID)
	if err != nil {
		return err
	}

	section.TitleEn = formDTO.TitleEn
	section.TitleAr = formDTO.TitleAr
	section.DescriptionEn = formDTO.DescriptionEn
	section.DescriptionAr = formDTO.DescriptionAr

	err = crs.repo.UpdateComplianceReportSection(section, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (crs *ComplianceReportService) DeleteComplianceReport(complianceReportID uint, currentUserID uint) error {
	tx := crs.repo.BeginTransaction()

	existingReport, err := crs.repo.GetComplianceReportByID(complianceReportID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReportActionPlans(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReportRequirements(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReportClauses(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReportDepartments(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReportSections(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.repo.DeleteComplianceReport(complianceReportID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go func() {
		err := crs.activityLogService.LogActivity(
			models.DeleteAction,
			models.ComplianceReport{}.TableName(),
			complianceReportID,
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

func (crs *ComplianceReportService) ExportComplianceReportToPDF(complianceReportID uint, currentUserID uint) (string, error) {
	complianceReport, err := crs.GetComplianceReportDetails(complianceReportID)
	if err != nil {
		return "", err
	}

	// Render template and generate PDF
	templatePaths := []string{"compliance/templates/pdf/compliance_report.html"}
	htmlContent, err := utility.GeneratePDFContent(complianceReport, templatePaths, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	pdfFileName := fmt.Sprintf("ComplianceReport_%d.pdf", complianceReportID)
	if err := utility.ConvertHTMLContentToPDF(htmlContent, pdfFileName); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (crs *ComplianceReportService) getComplianceStatusByScore(score float64) (*uint, error) {
	// Get all compliance statuses
	statuses, err := crs.repo.GetAllComplianceStatuses()
	if err != nil {
		return nil, err
	}

	// Find the status that matches the score range
	for _, status := range statuses {
		if status.ScoreFrom != nil && status.ScoreTo != nil {
			if score >= float64(*status.ScoreFrom) && score <= float64(*status.ScoreTo) {
				return &status.ID, nil
			}
		}
	}

	return nil, nil
}
