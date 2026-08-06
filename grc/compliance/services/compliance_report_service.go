package services

import (
	"fmt"
	"html/template"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
)

type ComplianceReportService struct {
	repo                repositories.ComplianceReportRepositoryInterface
	activityLogService  logServices.ActivityLogServiceInterface
	complianceStatusSvc ComplianceStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
}

type ComplianceReportServiceInterface interface {
	GetAll(currentUserID uint, queryParams *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error)
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
	complianceStatusSvc ComplianceStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) *ComplianceReportService {
	return &ComplianceReportService{
		repo:                repo,
		activityLogService:  activityLogService,
		complianceStatusSvc: complianceStatusSvc,
		userSvc:             userSvc,
	}
}

func (crs *ComplianceReportService) GetAll(currentUserID uint, queryParams *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error) {
	if currentUserID != 0 {
		isCyberUser, err := crs.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			return nil, err
		}
		queryParams.IsCyber = isCyberUser
	}
	reports, err := crs.repo.GetAll(queryParams)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (crs *ComplianceReportService) Create(currentUserID uint, request *dtos.ComplianceReportFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	if err := crs.validateComplianceReportFormDTO(request); err != nil {
		return 0, err
	}

	parsedDates, err := crs.parseComplianceReportDates(request)
	if err != nil {
		return 0, err
	}

	complianceReport := &models.ComplianceReport{
		Title:          request.Title,
		StartDate:      *parsedDates.StartDate,
		EndDate:        *parsedDates.EndDate,
		CreatedBYID:    currentUserID,
		AssessmentType: request.AssessmentType,
	}

	var isCyberUser bool
	if request.AssessmentType == "department" {
		isCyberUser, err = crs.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if request.AssessmentType == "framework" {
		complianceReport.ComplianceRegisterID = request.ComplianceRegisterID
		register, err := crs.repo.GetComplianceRegisterByID(*request.ComplianceRegisterID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		complianceReport.IsCyber = register.IsCyber
	} else {
		complianceReport.ResponsibleUnitID = request.ResponsibleUnitID
		complianceReport.IsCyber = isCyberUser
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

	var complianceAssessments []models.ComplianceAssessment
	if request.AssessmentType == "framework" {
		complianceAssessments, err = crs.repo.GetComplianceAssessmentByComplianceRegisterID(*request.ComplianceRegisterID, *parsedDates.StartDate, *parsedDates.EndDate, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		complianceAssessments, err = crs.repo.GetComplianceAssessmentByResponsibleUnitID(*request.ResponsibleUnitID, isCyberUser, *parsedDates.StartDate, *parsedDates.EndDate, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	orgUnitAssessments := make(map[uint][]models.ComplianceAssessment)
	for _, assessment := range complianceAssessments {
		for _, requirement := range assessment.AssessmentRequirements {
			for _, responsibleUnit := range requirement.ResponsibleUnits {
				orgUnitAssessments[responsibleUnit.ResponsibleUnitID] = append(orgUnitAssessments[responsibleUnit.ResponsibleUnitID], assessment)
			}
		}
	}

	frameworkType, maxValueInFramework := inferFrameworkMetaFromAssessments(complianceAssessments)

	latestAsmtIDByUnit := make(map[uint]uint)
	latestAsmtAtByUnit := make(map[uint]time.Time)
	for _, asmt := range complianceAssessments {
		asmtAt := asmt.CreatedAt
		seenRU := make(map[uint]struct{})
		for _, req := range asmt.AssessmentRequirements {
			for _, ru := range req.ResponsibleUnits {
				if ru == nil {
					continue
				}
				uid := ru.ResponsibleUnitID
				if _, done := seenRU[uid]; done {
					continue
				}
				seenRU[uid] = struct{}{}
				if cur, ok := latestAsmtAtByUnit[uid]; !ok || asmtAt.After(cur) {
					latestAsmtAtByUnit[uid] = asmtAt
					latestAsmtIDByUnit[uid] = asmt.ID
				}
			}
		}
	}

	if request.AssessmentType == "department" {
		unitID := *request.ResponsibleUnitID

		type reqAgg struct {
			Requirement *models.ComplianceAssessmentRequirement
			RUs         []*models.ComplianceAssessmentRequirementResponsibleUnit
		}

		type regAgg struct {
			Reqs map[uint]*reqAgg
		}

		byRegister := make(map[uint]*regAgg)

		for _, asmt := range complianceAssessments {
			for _, req := range asmt.AssessmentRequirements {
				if req == nil || req.ComplianceRequirement == nil {
					continue
				}

				for _, ru := range req.ResponsibleUnits {
					if ru == nil || ru.ResponsibleUnitID != unitID {
						continue
					}

					if req.ComplianceRequirement.ComplianceClause == nil || req.ComplianceRequirement.ComplianceClause.ComplianceRegister == nil {
						continue
					}

					if isCyberUser && !req.ComplianceRequirement.ComplianceClause.ComplianceRegister.IsCyber {
						continue
					}

					regID := req.ComplianceRequirement.ComplianceClause.ComplianceRegisterID

					if byRegister[regID] == nil {
						byRegister[regID] = &regAgg{Reqs: make(map[uint]*reqAgg)}
					}

					rid := req.ComplianceRequirementID
					if byRegister[regID].Reqs[rid] == nil {
						byRegister[regID].Reqs[rid] = &reqAgg{
							Requirement: req,
							RUs:         []*models.ComplianceAssessmentRequirementResponsibleUnit{},
						}
					}

					byRegister[regID].Reqs[rid].RUs = append(byRegister[regID].Reqs[rid].RUs, ru)
				}
			}
		}

		for regID, reg := range byRegister {
			ruBuckets := make([][]*models.ComplianceAssessmentRequirementResponsibleUnit, 0)
			allRUs := make([]*models.ComplianceAssessmentRequirementResponsibleUnit, 0)
			for _, ra := range reg.Reqs {
				if len(ra.RUs) > 0 {
					ruBuckets = append(ruBuckets, ra.RUs)
					allRUs = append(allRUs, ra.RUs...)
				}
			}

			deptAvg, status := crs.deptAverageForUnit(
				ruBuckets,
				frameworkType,
				maxValueInFramework,
			)

			isRiskBased, isRegNotification := assessmentTypesFromRUs(allRUs)

			dept := &models.ComplianceReportDepartment{
				ComplianceReportID:                 complianceReportID,
				PrimaryUnitID:                      unitID,
				ComplianceRegisterID:               &regID,
				ComplianceStatusID:                 &status,
				CompliancePercentage:               deptAvg,
				IsRiskBasedAssessment:              isRiskBased,
				IsRegulationNotificationAssessment: isRegNotification,
				CreatedAt:                          time.Now(),
				CreatedBYID:                        currentUserID,
			}

			deptID, err := crs.repo.CreateComplianceReportDepartment(dept, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			clauseMap := make(map[uint]uint)

			for _, ra := range reg.Reqs {
				clause := ra.Requirement.ComplianceRequirement.ComplianceClause
				if clause == nil {
					continue
				}
				if _, ok := clauseMap[clause.ID]; ok {
					continue
				}

				rc := &models.ComplianceReportClause{
					ComplianceReportID:           complianceReportID,
					ComplianceReportDepartmentID: deptID,
					ClauseNo:                     clause.ClauseNo,
					ClauseTitleEn:                clause.TitleEn,
					ClauseTitleAr:                clause.TitleAr,
					DescriptionEn:                clause.DescriptionEn,
					DescriptionAr:                clause.DescriptionAr,
					CreatedBYID:                  currentUserID,
				}

				cid, err := crs.repo.CreateComplianceReportClause(rc, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
				clauseMap[clause.ID] = cid
			}

			for _, ra := range reg.Reqs {
				var currStatusID *uint
				if ra.Requirement.ComplianceStatus != nil {
					currStatusID = &ra.Requirement.ComplianceStatus.ID
				} else if ra.Requirement.ComplianceStatusID != nil {
					currStatusID = ra.Requirement.ComplianceStatusID
				}

				avgScore, avgStatusID := avgScoreAndStatusFromRUsFrameworkAware(ra.RUs, "percentage", 100)
				if avgStatusID == nil {
					avgStatusID = currStatusID
				}

				isAck, ackUserID, ackComment := acknowledgementFromRUs(ra.RUs)

				inherentScore := ra.Requirement.ComplianceRequirement.InherentRiskScore
				inherentRiskRatingID := ra.Requirement.ComplianceRequirement.InherentRiskRatingID
				residualScore := ra.Requirement.AverageResidualRiskScore
				if residualScore == nil {
					residualScore = residualScoreFromRUs(ra.RUs)
				}
				residualRiskRatingID := ra.Requirement.ComplianceRequirement.ResidualRiskRatingID

				isRiskBased, isRegNotification := assessmentTypesFromRUs(ra.RUs)

				cr := &models.ComplianceReportRequirement{
					ComplianceReportID:                 complianceReportID,
					ComplianceRegisterID:               dept.ComplianceRegisterID,
					ComplianceReportClauseID:           clauseMap[ra.Requirement.ComplianceRequirement.ComplianceClause.ID],
					TitleEn:                            ra.Requirement.ComplianceRequirement.TitleEn,
					TitleAr:                            ra.Requirement.ComplianceRequirement.TitleAr,
					AverageScore:                       &avgScore,
					InherentScore:                      inherentScore,
					InherentRiskRatingID:               inherentRiskRatingID,
					ResidualScore:                      residualScore,
					ResidualRiskRatingID:               residualRiskRatingID,
					ComplianceStatusID:                 currStatusID,
					AverageComplianceStatusID:          avgStatusID,
					IsAcknowledged:                     isAck,
					AcknowledgeUserID:                  ackUserID,
					AcknowledgeComment:                 ackComment,
					IsRiskBasedAssessment:              isRiskBased,
					IsRegulationNotificationAssessment: isRegNotification,
					CreatedBYID:                        currentUserID,
				}

				crID, err := crs.repo.CreateComplianceReportRequirement(cr, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				for _, ru := range ra.RUs {
					for _, ap := range ru.ActionPlans {
						rap := &models.ComplianceReportActionPlan{
							ComplianceReportID:            complianceReportID,
							ComplianceReportRequirementID: crID,
							Title:                         ap.Title,
							Description:                   ap.Description,
							ActionPlanStatusID:            ap.ActionPlanStatusID,
							TargetDate:                    ap.TargetDate,
							SeverityLevelID:               ap.SeverityLevelID,
							CreatedAt:                     time.Now(),
							CreatedBYID:                   currentUserID,
						}
						if _, err := crs.repo.CreateComplianceReportActionPlan(rap, tx); err != nil {
							tx.Rollback()
							return 0, err
						}
					}
				}
			}
		}
	} else {
		type reqAgg struct {
			Requirement *models.ComplianceAssessmentRequirement
			RUs         []*models.ComplianceAssessmentRequirementResponsibleUnit
			RegisterID  uint
		}
		type unitAgg struct {
			Reqs map[uint]*reqAgg
		}
		orgUnitAgg := make(map[uint]*unitAgg)

		for _, asmt := range complianceAssessments {
			for _, req := range asmt.AssessmentRequirements {
				if req == nil || req.ComplianceRequirement == nil || req.ComplianceRequirement.ComplianceClause == nil {
					continue
				}
				for _, ru := range req.ResponsibleUnits {
					if ru == nil {
						continue
					}
					unitID := ru.ResponsibleUnitID

					latestID, ok := latestAsmtIDByUnit[unitID]
					if !ok || asmt.ID != latestID {
						continue
					}
					if _, ok := orgUnitAgg[unitID]; !ok {
						orgUnitAgg[unitID] = &unitAgg{Reqs: make(map[uint]*reqAgg)}
					}
					rid := req.ComplianceRequirementID
					if _, ok := orgUnitAgg[unitID].Reqs[rid]; !ok {
						orgUnitAgg[unitID].Reqs[rid] = &reqAgg{
							Requirement: req,
							RUs:         make([]*models.ComplianceAssessmentRequirementResponsibleUnit, 0, 1),
							RegisterID:  req.ComplianceRequirement.ComplianceClause.ComplianceRegisterID,
						}
					}
					orgUnitAgg[unitID].Reqs[rid].RUs = append(orgUnitAgg[unitID].Reqs[rid].RUs, ru)
				}
			}
		}

		for unitID, ua := range orgUnitAgg {
			ruBuckets := make([][]*models.ComplianceAssessmentRequirementResponsibleUnit, 0, len(ua.Reqs))
			allRUs := make([]*models.ComplianceAssessmentRequirementResponsibleUnit, 0)
			var complianceRegisterId uint
			for _, rag := range ua.Reqs {
				if rag != nil && len(rag.RUs) > 0 {
					ruBuckets = append(ruBuckets, rag.RUs)
					allRUs = append(allRUs, rag.RUs...)
					complianceRegisterId = rag.RegisterID
				}
			}

			deptAvg, status := crs.deptAverageForUnit(ruBuckets, frameworkType, maxValueInFramework)
			isRiskBased, isRegNotification := assessmentTypesFromRUs(allRUs)

			dept := &models.ComplianceReportDepartment{
				ComplianceReportID:                 complianceReportID,
				PrimaryUnitID:                      unitID,
				CompliancePercentage:               deptAvg,
				ComplianceStatusID:                 &status,
				ComplianceRegisterID:               &complianceRegisterId,
				IsRiskBasedAssessment:              isRiskBased,
				IsRegulationNotificationAssessment: isRegNotification,
				CreatedAt:                          time.Now(),
				CreatedBYID:                        currentUserID,
			}
			deptID, err := crs.repo.CreateComplianceReportDepartment(dept, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			clauseIDByClausePK := make(map[uint]uint)

			for _, rag := range ua.Reqs {
				req := rag.Requirement
				clause := req.ComplianceRequirement.ComplianceClause
				if clause == nil {
					continue
				}
				if _, exists := clauseIDByClausePK[clause.ID]; exists {
					continue
				}
				rc := &models.ComplianceReportClause{
					ComplianceReportID:           complianceReportID,
					ComplianceReportDepartmentID: deptID,
					ClauseNo:                     clause.ClauseNo,
					ClauseTitleEn:                clause.TitleEn,
					ClauseTitleAr:                clause.TitleAr,
					DescriptionEn:                clause.DescriptionEn,
					DescriptionAr:                clause.DescriptionAr,
					CreatedBYID:                  currentUserID,
				}
				newID, err := crs.repo.CreateComplianceReportClause(rc, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
				clauseIDByClausePK[clause.ID] = newID
			}

			for _, rag := range ua.Reqs {
				req := rag.Requirement
				if req == nil || req.ComplianceRequirement == nil || req.ComplianceRequirement.ComplianceClause == nil {
					continue
				}

				avgScore, statusIDPtr := avgScoreAndStatusFromRUsFrameworkAware(rag.RUs, frameworkType, maxValueInFramework)

				var currStatusID *uint
				if req.ComplianceStatus != nil {
					currStatusID = req.ComplianceRequirement.ComplianceStatusID
				}

				isAck, ackUserID, ackComment := acknowledgementFromRUs(rag.RUs)

				inherentScore := req.ComplianceRequirement.InherentRiskScore
				inherentRiskRatingID := req.ComplianceRequirement.InherentRiskRatingID
				residualScore := req.AverageResidualRiskScore
				if residualScore == nil {
					residualScore = residualScoreFromRUs(rag.RUs)
				}
				residualRiskRatingID := req.ComplianceRequirement.ResidualRiskRatingID

				isRiskBased, isRegNotification := assessmentTypesFromRUs(rag.RUs)

				crr := &models.ComplianceReportRequirement{
					ComplianceReportID:                 complianceReportID,
					ComplianceRegisterID:               &complianceRegisterId,
					ComplianceReportClauseID:           clauseIDByClausePK[req.ComplianceRequirement.ComplianceClause.ID],
					TitleEn:                            req.ComplianceRequirement.TitleEn,
					TitleAr:                            req.ComplianceRequirement.TitleAr,
					AverageScore:                       &avgScore,
					InherentScore:                      inherentScore,
					InherentRiskRatingID:               inherentRiskRatingID,
					ResidualScore:                      residualScore,
					ResidualRiskRatingID:               residualRiskRatingID,
					ComplianceStatusID:                 currStatusID,
					AverageComplianceStatusID:          statusIDPtr,
					IsAcknowledged:                     isAck,
					AcknowledgeUserID:                  ackUserID,
					AcknowledgeComment:                 ackComment,
					IsRiskBasedAssessment:              isRiskBased,
					IsRegulationNotificationAssessment: isRegNotification,
					CreatedBYID:                        currentUserID,
				}

				reqReportID, err := crs.repo.CreateComplianceReportRequirement(crr, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				for _, ru := range rag.RUs {
					for _, ap := range ru.ActionPlans {
						rAp := &models.ComplianceReportActionPlan{
							ComplianceReportID:            complianceReportID,
							ComplianceReportRequirementID: reqReportID,
							Title:                         ap.Title,
							Description:                   ap.Description,
							ActionPlanStatusID:            ap.ActionPlanStatusID,
							TargetDate:                    ap.TargetDate,
							SeverityLevelID:               ap.SeverityLevelID,
							CreatedAt:                     time.Now(),
							CreatedBYID:                   currentUserID,
						}
						if _, err := crs.repo.CreateComplianceReportActionPlan(rAp, tx); err != nil {
							tx.Rollback()
							return 0, err
						}
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
	complianceReportDetails.AssessmentType = complianceReport.AssessmentType

	// Map compliance register for framework type
	if complianceReport.AssessmentType == "framework" {
		complianceRegister := &dtos.ComplianceRegister{}
		if err := mapstructure.Decode(
			complianceReport.ComplianceRegister,
			complianceRegister,
		); err != nil {
			return nil, err
		}
		complianceReportDetails.ComplianceRegister = complianceRegister
	}

	// Map responsible unit for department type
	if complianceReport.AssessmentType == "department" {
		responsibleUnit := &organizationDTOs.OrganizationUnitDTO{}
		if err := mapstructure.Decode(
			complianceReport.ResponsibleUnit,
			responsibleUnit,
		); err != nil {
			return nil, err
		}
		complianceReportDetails.ResponsibleUnit = responsibleUnit
		complianceReportDetails.ResponsibleUnitID = complianceReport.ResponsibleUnitID
	}

	// Map created by user
	complianceReportDetails.CreatedByUser.ID = int(complianceReport.CreatedByUser.ID)
	complianceReportDetails.CreatedByUser.FirstNameEN = complianceReport.CreatedByUser.FirstNameEn
	complianceReportDetails.CreatedByUser.FirstNameAR = complianceReport.CreatedByUser.FirstNameAr
	complianceReportDetails.CreatedByUser.LastNameEN = complianceReport.CreatedByUser.LastNameEn
	complianceReportDetails.CreatedByUser.LastNameAR = complianceReport.CreatedByUser.LastNameAr
	complianceReportDetails.CreatedByUser.Email = complianceReport.CreatedByUser.Email

	// Map sections
	for _, section := range complianceReport.Sections {
		var descriptionEn template.HTML
		if section.DescriptionEn != nil {
			descriptionEn = template.HTML(*section.DescriptionEn)
		}
		var descriptionAr template.HTML
		if section.DescriptionAr != nil {
			descriptionAr = template.HTML(*section.DescriptionAr)
		}
		sectionDTO := &dtos.ComplianceReportSection{
			ID:            section.ID,
			TitleEn:       section.TitleEn,
			TitleAr:       section.TitleAr,
			DescriptionEn: descriptionEn,
			DescriptionAr: descriptionAr,
			Type:          section.Type,
			Order:         section.OrderNum,
		}
		complianceReportDetails.ComplianceReportSections = append(complianceReportDetails.ComplianceReportSections, sectionDTO)

		if section.Type == "compliance-status-by-department" {
			complianceReportDeptSection := make(map[uint]*dtos.ComplianceReportDepartmentSection)

			for _, department := range departments {
				if department == nil || department.ComplianceRegisterID == nil {
					continue
				}

				var complianceStatus *dtos.ComplianceStatusDTO
				if department.ComplianceStatus != nil {
					complianceStatus = &dtos.ComplianceStatusDTO{
						ID:        int(department.ComplianceStatus.ID),
						TitleEn:   department.ComplianceStatus.TitleEn,
						TitleAr:   department.ComplianceStatus.TitleAr,
						ColorCode: department.ComplianceStatus.ColorCode,
						Type:      department.ComplianceStatus.Type,
						ScoreFrom: department.ComplianceStatus.ScoreFrom,
						ScoreTo:   department.ComplianceStatus.ScoreTo,
						Status:    department.ComplianceStatus.Status,
					}
				}

				var complianceRegister *dtos.ComplianceRegister
				if department.ComplianceRegister != nil {
					complianceRegister = &dtos.ComplianceRegister{
						ID:      department.ComplianceRegister.ID,
						TitleEn: department.ComplianceRegister.TitleEn,
						TitleAr: department.ComplianceRegister.TitleAr,
					}
				}

				departmentDTO := &dtos.ComplianceReportDepartment{
					OrganizationUnit: settingsDTOs.OrganizationUnitDTO{
						ID:      int(department.OrganizationUnit.ID),
						TitleEn: department.OrganizationUnit.TitleEn,
						TitleAr: department.OrganizationUnit.TitleAr,
					},
					ComplianceRegister:                 complianceRegister,
					ComplianceStatus:                   complianceStatus,
					CompliancePercentage:               &department.CompliancePercentage,
					IsRiskBasedAssessment:              department.IsRiskBasedAssessment,
					IsRegulationNotificationAssessment: department.IsRegulationNotificationAssessment,
					NonCompliantAreas:                  []*dtos.ComplianceReportRequirement{},
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
							ID:                                 nonCompliantRequirement.ID,
							TitleEn:                            nonCompliantRequirement.TitleEn,
							TitleAr:                            nonCompliantRequirement.TitleAr,
							Score:                              nonCompliantRequirement.Score,
							AverageScore:                       nonCompliantRequirement.AverageScore,
							InherentScore:                      nonCompliantRequirement.InherentScore,
							InherentRiskRatingID:               nonCompliantRequirement.InherentRiskRatingID,
							ResidualScore:                      nonCompliantRequirement.ResidualScore,
							ResidualRiskRatingID:               nonCompliantRequirement.ResidualRiskRatingID,
							IsAcknowledged:                     nonCompliantRequirement.IsAcknowledged,
							AcknowledgeComment:                 nonCompliantRequirement.AcknowledgeComment,
							IsRiskBasedAssessment:              nonCompliantRequirement.IsRiskBasedAssessment,
							IsRegulationNotificationAssessment: nonCompliantRequirement.IsRegulationNotificationAssessment,
							ActionPlans:                        []*dtos.ComplianceReportActionPlan{},
						}

						if nonCompliantRequirement.ComplianceStatus != nil {
							nonCompliantArea.ComplianceStatus = &dtos.ComplianceStatusDTO{
								ID:        int(nonCompliantRequirement.ComplianceStatus.ID),
								Type:      nonCompliantRequirement.ComplianceStatus.Type,
								TitleEn:   nonCompliantRequirement.ComplianceStatus.TitleEn,
								TitleAr:   nonCompliantRequirement.ComplianceStatus.TitleAr,
								ColorCode: nonCompliantRequirement.ComplianceStatus.ColorCode,
							}
						}

						if nonCompliantRequirement.AverageComplianceStatus != nil {
							nonCompliantArea.AverageComplianceStatus = &dtos.ComplianceStatusDTO{
								ID:        int(nonCompliantRequirement.AverageComplianceStatus.ID),
								Type:      nonCompliantRequirement.AverageComplianceStatus.Type,
								TitleEn:   nonCompliantRequirement.AverageComplianceStatus.TitleEn,
								TitleAr:   nonCompliantRequirement.AverageComplianceStatus.TitleAr,
								ColorCode: nonCompliantRequirement.AverageComplianceStatus.ColorCode,
							}
						}

						if nonCompliantRequirement.AcknowledgeUser != nil {
							nonCompliantArea.AcknowledgeUser = &organizationDTOs.UserMinimalDTO{
								ID:          int(nonCompliantRequirement.AcknowledgeUser.ID),
								FirstNameEN: nonCompliantRequirement.AcknowledgeUser.FirstNameEn,
								FirstNameAR: nonCompliantRequirement.AcknowledgeUser.FirstNameAr,
								LastNameEN:  nonCompliantRequirement.AcknowledgeUser.LastNameEn,
								LastNameAR:  nonCompliantRequirement.AcknowledgeUser.LastNameAr,
							}
						}

						if nonCompliantRequirement.InherentRiskRating != nil {
							nonCompliantArea.InherentRiskRating = &dtos.ComplianceAssessmentRatingDTO{
								ID:        nonCompliantRequirement.InherentRiskRating.ID,
								TitleEn:   nonCompliantRequirement.InherentRiskRating.TitleEn,
								TitleAr:   nonCompliantRequirement.InherentRiskRating.TitleAr,
								ColorCode: nonCompliantRequirement.InherentRiskRating.ColorCode,
								Type:      nonCompliantRequirement.InherentRiskRating.Type,
							}
						}

						if nonCompliantRequirement.ResidualRiskRating != nil {
							nonCompliantArea.ResidualRiskRating = &dtos.ComplianceAssessmentRatingDTO{
								ID:        nonCompliantRequirement.ResidualRiskRating.ID,
								TitleEn:   nonCompliantRequirement.ResidualRiskRating.TitleEn,
								TitleAr:   nonCompliantRequirement.ResidualRiskRating.TitleAr,
								ColorCode: nonCompliantRequirement.ResidualRiskRating.ColorCode,
								Type:      nonCompliantRequirement.ResidualRiskRating.Type,
							}
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
							nonCompliantArea.ActionPlans = append(nonCompliantArea.ActionPlans, actionPlanDTO)
						}

						departmentDTO.NonCompliantAreas = append(departmentDTO.NonCompliantAreas, nonCompliantArea)
					}
				}

				regID := *department.ComplianceRegisterID
				if complianceReportDeptSection[regID] == nil {
					complianceReportDeptSection[regID] = &dtos.ComplianceReportDepartmentSection{
						ComplianceRegister:         complianceRegister,
						ComplianceReportDepartment: []*dtos.ComplianceReportDepartment{},
					}
				}
				complianceReportDeptSection[regID].ComplianceReportDepartment = append(
					complianceReportDeptSection[regID].ComplianceReportDepartment,
					departmentDTO,
				)
			}

			for _, deptSection := range complianceReportDeptSection {
				sectionDTO.ComplianceReportDepartments = append(sectionDTO.ComplianceReportDepartments, deptSection)
			}
		}

		if section.Type == "compliance-status-by-standard" {
			complianceReportClauseSections := make(map[uint]*dtos.ComplianceReportClauseSection)

			for _, department := range departments {
				if department == nil || department.ComplianceRegisterID == nil {
					continue
				}

				regID := *department.ComplianceRegisterID

				if complianceReportClauseSections[regID] == nil {
					var complianceRegister *dtos.ComplianceRegister
					if department.ComplianceRegister != nil {
						complianceRegister = &dtos.ComplianceRegister{
							ID:      department.ComplianceRegister.ID,
							TitleEn: department.ComplianceRegister.TitleEn,
							TitleAr: department.ComplianceRegister.TitleAr,
						}
					}
					complianceReportClauseSections[regID] = &dtos.ComplianceReportClauseSection{
						ComplianceRegister:     complianceRegister,
						ComplianceReportClause: []*dtos.ComplianceReportClause{},
					}
				}

				uniqueClauses := make(map[string]*dtos.ComplianceReportClause)
				uniqueRequirements := make(map[string]*dtos.ComplianceReportRequirement)

				for _, clause := range department.ComplianceReportClauses {
					if clause == nil {
						continue
					}

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

					for _, requirement := range clause.ComplianceReportRequirements {
						if requirement == nil {
							continue
						}

						requirementKey := fmt.Sprintf("%s_%s_%s", clause.ClauseNo, requirement.TitleEn, requirement.TitleAr)

						var requirementDTO *dtos.ComplianceReportRequirement
						if existingReq, exists := uniqueRequirements[requirementKey]; exists {
							requirementDTO = existingReq
						} else {
							requirementDTO = &dtos.ComplianceReportRequirement{
								ID:                                 requirement.ID,
								TitleEn:                            requirement.TitleEn,
								TitleAr:                            requirement.TitleAr,
								Score:                              requirement.Score,
								AverageScore:                       requirement.AverageScore,
								InherentScore:                      requirement.InherentScore,
								InherentRiskRatingID:               requirement.InherentRiskRatingID,
								ResidualScore:                      requirement.ResidualScore,
								ResidualRiskRatingID:               requirement.ResidualRiskRatingID,
								IsAcknowledged:                     requirement.IsAcknowledged,
								AcknowledgeComment:                 requirement.AcknowledgeComment,
								IsRiskBasedAssessment:              requirement.IsRiskBasedAssessment,
								IsRegulationNotificationAssessment: requirement.IsRegulationNotificationAssessment,
								ActionPlans:                        []*dtos.ComplianceReportActionPlan{},
							}

							if requirement.ComplianceStatus != nil {
								requirementDTO.ComplianceStatus = &dtos.ComplianceStatusDTO{
									ID:        int(requirement.ComplianceStatus.ID),
									Type:      requirement.ComplianceStatus.Type,
									TitleEn:   requirement.ComplianceStatus.TitleEn,
									TitleAr:   requirement.ComplianceStatus.TitleAr,
									ColorCode: requirement.ComplianceStatus.ColorCode,
								}
							}

							if requirement.AverageComplianceStatus != nil {
								requirementDTO.AverageComplianceStatus = &dtos.ComplianceStatusDTO{
									ID:        int(requirement.AverageComplianceStatus.ID),
									Type:      requirement.AverageComplianceStatus.Type,
									TitleEn:   requirement.AverageComplianceStatus.TitleEn,
									TitleAr:   requirement.AverageComplianceStatus.TitleAr,
									ColorCode: requirement.AverageComplianceStatus.ColorCode,
								}
							}

							if requirement.AcknowledgeUser != nil {
								requirementDTO.AcknowledgeUser = &organizationDTOs.UserMinimalDTO{
									ID:          int(requirement.AcknowledgeUser.ID),
									FirstNameEN: requirement.AcknowledgeUser.FirstNameEn,
									FirstNameAR: requirement.AcknowledgeUser.FirstNameAr,
									LastNameEN:  requirement.AcknowledgeUser.LastNameEn,
									LastNameAR:  requirement.AcknowledgeUser.LastNameAr,
								}
							}

							if requirement.InherentRiskRating != nil {
								requirementDTO.InherentRiskRating = &dtos.ComplianceAssessmentRatingDTO{
									ID:        requirement.InherentRiskRating.ID,
									TitleEn:   requirement.InherentRiskRating.TitleEn,
									TitleAr:   requirement.InherentRiskRating.TitleAr,
									ColorCode: requirement.InherentRiskRating.ColorCode,
									Type:      requirement.InherentRiskRating.Type,
								}
							}

							if requirement.ResidualRiskRating != nil {
								requirementDTO.ResidualRiskRating = &dtos.ComplianceAssessmentRatingDTO{
									ID:        requirement.ResidualRiskRating.ID,
									TitleEn:   requirement.ResidualRiskRating.TitleEn,
									TitleAr:   requirement.ResidualRiskRating.TitleAr,
									ColorCode: requirement.ResidualRiskRating.ColorCode,
									Type:      requirement.ResidualRiskRating.Type,
								}
							}

							uniqueRequirements[requirementKey] = requirementDTO
							uniqueClauses[clause.ClauseNo].ComplianceRequirements = append(uniqueClauses[clause.ClauseNo].ComplianceRequirements, requirementDTO)
						}

						for _, actionPlan := range requirement.ComplianceReportActionPlans {
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

				var clauseSlice []*dtos.ComplianceReportClause
				for _, clauseDTO := range uniqueClauses {
					clauseSlice = append(clauseSlice, clauseDTO)
				}

				sort.Slice(clauseSlice, func(i, j int) bool {
					iParts := strings.Split(clauseSlice[i].ClauseNo, ".")
					jParts := strings.Split(clauseSlice[j].ClauseNo, ".")

					for k := 0; k < len(iParts) && k < len(jParts); k++ {
						iNum, errI := strconv.Atoi(iParts[k])
						jNum, errJ := strconv.Atoi(jParts[k])

						if errI != nil || errJ != nil {
							return iParts[k] < jParts[k]
						}

						if iNum != jNum {
							return iNum < jNum
						}
					}

					return len(iParts) < len(jParts)
				})

				complianceReportClauseSections[regID].ComplianceReportClause = append(
					complianceReportClauseSections[regID].ComplianceReportClause,
					clauseSlice...,
				)
			}

			for _, clauseSection := range complianceReportClauseSections {
				sectionDTO.ComplianceReportStandards = append(sectionDTO.ComplianceReportStandards, clauseSection)
			}
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
				requirement.ComplianceStatus.Type != "not_applicable" {
				// requirement.Score != nil {
				// totalScore += float64(*requirement.Score)
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

	var descriptionEn template.HTML
	if section.DescriptionEn != nil {
		descriptionEn = template.HTML(*section.DescriptionEn)
	}
	var descriptionAr template.HTML
	if section.DescriptionAr != nil {
		descriptionAr = template.HTML(*section.DescriptionAr)
	}
	sectionDTO := &dtos.ComplianceReportSection{
		ID:            section.ID,
		TitleEn:       section.TitleEn,
		TitleAr:       section.TitleAr,
		DescriptionEn: descriptionEn,
		DescriptionAr: descriptionAr,
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
			"compliance", "compliance_report", "", nil,
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

func (crs *ComplianceReportService) validateComplianceReportFormDTO(dto *dtos.ComplianceReportFormDTO) error {
	if dto.AssessmentType != "framework" && dto.AssessmentType != "department" {
		return utility.NewValidationError(fmt.Errorf("invalid assessment type"), "invalid_assessment_type", "assessment_type")
	}

	if dto.AssessmentType == "framework" && dto.ComplianceRegisterID == nil {
		return utility.NewValidationError(fmt.Errorf("compliance register ID is required for framework assessment"), "required", "compliance_register_id")
	}

	if dto.AssessmentType == "department" && dto.ResponsibleUnitID == nil {
		return utility.NewValidationError(fmt.Errorf("responsible unit ID is required for department assessment"), "required", "responsible_unit_id")
	}

	return nil
}

func inferFrameworkMetaFromAssessments(as []models.ComplianceAssessment) (frameworkType string, maxValue int) {
	frameworkType = "percentage"
	maxValue = 100

	for _, a := range as {
		for _, r := range a.AssessmentRequirements {
			for _, ru := range r.ResponsibleUnits {
				if ru == nil {
					continue
				}
				if ru.Score != nil {
					return "percentage", 100
				}
			}
		}
	}

	return frameworkType, maxValue
}

func avgScoreAndStatusFromRUsFrameworkAware(rus []*models.ComplianceAssessmentRequirementResponsibleUnit, frameworkType string, maxValueInFramework int) (avgPct float64, modalStatusID *uint) {

	if maxValueInFramework <= 0 {
		switch frameworkType {
		case "percentage":
			maxValueInFramework = 100
		case "binary":
			maxValueInFramework = 1
		default:
			maxValueInFramework = 100
		}
	}

	perReqMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			return float64(maxValueInFramework)
		default:
			return 100.0
		}
	}()

	var totalScore, maxPossible float64
	statusFreq := map[uint]int{}

	for _, ru := range rus {
		if ru == nil {
			continue
		}

		if ru.IsNotApplicable {
			continue
		}

		switch frameworkType {
		case "point", "maturity":
			if ru.Score != nil {
				v := float64(*ru.Score)
				maxPossible += perReqMax
				totalScore += v
			}
		case "percentage":
			if ru.Score != nil {
				maxPossible += perReqMax
				totalScore += float64(*ru.Score)
			}
		case "binary":
			maxPossible += perReqMax
			if ru.Score != nil && *ru.Score > 0 {
				totalScore += 1.0
			}
		default:
			if ru.Score != nil {
				maxPossible += 100.0
				totalScore += float64(*ru.Score)
			}
		}

		if ru.ComplianceStatusID != nil {
			statusFreq[*ru.ComplianceStatusID]++
		}
	}

	if maxPossible <= 0 {
		return 0, nil
	}

	pct := (totalScore / maxPossible) * 100.0
	if pct == 100 {
		pct = 99.99
	}
	avgPct = math.Round(pct*100) / 100

	var modeID uint
	var best int
	for k, c := range statusFreq {
		if c > best {
			best = c
			modeID = k
		}
	}
	if best == 0 {
		return avgPct, nil
	}
	return avgPct, &modeID
}

func acknowledgementFromRUs(rus []*models.ComplianceAssessmentRequirementResponsibleUnit) (isAck bool, ackUserID *uint, ackComment *string) {
	if len(rus) == 0 {
		return false, nil, nil
	}

	acknowledged := false
	allAcknowledged := true
	var firstAckRU *models.ComplianceAssessmentRequirementResponsibleUnit

	for _, ru := range rus {
		if ru == nil {
			continue
		}
		if ru.IsAcknowledged {
			acknowledged = true
			if firstAckRU == nil {
				firstAckRU = ru
			}
		} else {
			allAcknowledged = false
		}
	}

	if !acknowledged {
		return false, nil, nil
	}

	if firstAckRU != nil {
		return allAcknowledged, firstAckRU.AcknowledgeUserID, firstAckRU.AcknowledgeComment
	}
	return false, nil, nil
}

func residualScoreFromRUs(rus []*models.ComplianceAssessmentRequirementResponsibleUnit) *float64 {
	if len(rus) == 0 {
		return nil
	}

	var total float64
	var count int
	for _, ru := range rus {
		if ru == nil {
			continue
		}
		if ru.ResidualRiskScore != nil {
			total += *ru.ResidualRiskScore
			count++
		}
	}

	if count == 0 {
		return nil
	}

	avg := math.Round((total/float64(count))*100) / 100
	return &avg
}

func assessmentTypesFromRUs(rus []*models.ComplianceAssessmentRequirementResponsibleUnit) (isRiskBased *bool, isRegNotification *bool) {
	if len(rus) == 0 {
		return nil, nil
	}

	hasRiskBased := false
	hasRegNotification := false

	for _, ru := range rus {
		if ru == nil {
			continue
		}
		if ru.IsRiskBasedAssessment != nil && *ru.IsRiskBasedAssessment {
			hasRiskBased = true
		}
		if ru.IsRegulationNotificationAssessment != nil && *ru.IsRegulationNotificationAssessment {
			hasRegNotification = true
		}
	}

	var riskBasedPtr *bool
	var regNotificationPtr *bool

	if hasRiskBased {
		riskBasedPtr = &hasRiskBased
	}
	if hasRegNotification {
		regNotificationPtr = &hasRegNotification
	}

	return riskBasedPtr, regNotificationPtr
}

func (crs *ComplianceReportService) deptAverageForUnit(ruBuckets [][]*models.ComplianceAssessmentRequirementResponsibleUnit, frameworkType string, maxValueInFramework int) (float64, uint) {

	var sum float64
	var count int

	statusNA, err := crs.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		return 0, 0
	}

	statusC, err := crs.complianceStatusSvc.GetByType("compliant")
	if err != nil {
		return 0, 0
	}

	statusP, err := crs.complianceStatusSvc.GetByType("partially_compliant")
	if err != nil {
		return 0, 0
	}

	statusNI, err := crs.complianceStatusSvc.GetByType("not_implemented")
	if err != nil {
		return 0, 0
	}

	hasC, hasP, hasNI := false, false, false

	for _, rus := range ruBuckets {
		if len(rus) == 0 {
			continue
		}

		avg, _ := avgScoreAndStatusFromRUsFrameworkAware(
			rus,
			frameworkType,
			maxValueInFramework,
		)

		applicable := false
		for _, ru := range rus {
			if ru != nil && !ru.IsNotApplicable {
				applicable = true

				if ru.ComplianceAssessmentRequirement != nil && ru.ComplianceAssessmentRequirement.ComplianceStatusID != nil {
					switch *ru.ComplianceAssessmentRequirement.ComplianceStatusID {
					case statusC.ID:
						hasC = true
					case statusP.ID:
						hasP = true
					case statusNI.ID:
						hasNI = true
					}
				}

				break
			}
		}

		if !applicable {
			continue
		}

		sum += avg
		count++
	}

	statusID := statusNA.ID

	switch {
	case hasC && !hasP && !hasNI:
		statusID = statusC.ID
	case hasC && (hasP || hasNI):
		statusID = statusP.ID
	case hasP:
		statusID = statusP.ID
	case hasNI:
		statusID = statusNI.ID
	default:
		statusID = statusNA.ID
	}

	if count == 0 {
		return 0, statusID
	}

	deptAvg := math.Round((sum/float64(count))*100) / 100

	return deptAvg, statusID
}
