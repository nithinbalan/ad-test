package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	complianceDTOs "jethur/grc/compliance/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	logServices "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"
	"sort"
	"time"

	"gorm.io/gorm"
)

type AuditFinalReportServiceInterface interface {
	SetAuditReportWorkflowService(svc AuditReportWorkflowServiceInterface)
	GetByID(id uint) (*models.AuditFinalReport, error)
	GetFullDetailsByID(auditID uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error)
	Create(auditID uint, currentUserID uint) (uint, error)
	Update(auditFinalReport *models.AuditFinalReport) error
	GetByAuditID(auditID uint) (*models.AuditFinalReport, error)
	Delete(auditID uint, currentUserID uint) error
}

type AuditFinalReportService struct {
	repo                      repositories.AuditFinalReportRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditSvc                  AuditServiceInterface
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface
	auditDraftReportSvc       AuditDraftReportServiceInterface
	auditDraftFindingSvc      AuditDraftFindingServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface
	auditFindingSvc           AuditFindingServiceInterface
	auditReportWorkflowSvc    AuditReportWorkflowServiceInterface
}

func NewAuditFinalReportService(
	repo repositories.AuditFinalReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftReportStatusSvc AuditDraftReportStatusServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
) AuditFinalReportServiceInterface {
	return &AuditFinalReportService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditSvc:                  auditSvc,
		auditDraftReportStatusSvc: auditDraftReportStatusSvc,
		auditDraftReportSvc:       auditDraftReportSvc,
		auditDraftFindingSvc:      auditDraftFindingSvc,
		userSvc:                   userSvc,
		auditFinalReportStatusSvc: auditFinalReportStatusSvc,
		auditFindingSvc:           auditFindingSvc,
	}
}
func (afrs *AuditFinalReportService) SetAuditReportWorkflowService(svc AuditReportWorkflowServiceInterface) {
	afrs.auditReportWorkflowSvc = svc
}

func (afrs *AuditFinalReportService) GetByID(id uint) (*models.AuditFinalReport, error) {
	auditFinalReport, err := afrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinalReport, nil
}

func (afrs *AuditFinalReportService) GetFullDetailsByID(auditID uint, currentUserID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFinalReportDetailDTO, error) {
	// Check if the audit exists
	if _, err := afrs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}

	// Fetch the full details of the final report
	auditFinalReport, err := afrs.repo.GetFullDetailsByID(auditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if no data found
			return nil, nil
		}
		return nil, err
	}

	createdByUser := mapUserToUserMinimalDTO(auditFinalReport.CreatedByUser)
	var updatedByUser *organizationDTOs.UserMinimalDTO
	if auditFinalReport.UpdatedByUser != nil {
		updatedByUserDTO := mapUserToUserMinimalDTO(auditFinalReport.UpdatedByUser)
		updatedByUser = &updatedByUserDTO
	}
	// Map the AuditDraftReportStatus
	auditFinalReportStatusDTO := auditDTOs.AuditFinalReportStatusDTO{
		ID:        auditFinalReport.AuditFinalReportStatus.ID,
		TitleEn:   auditFinalReport.AuditFinalReportStatus.TitleEn,
		TitleAr:   auditFinalReport.AuditFinalReportStatus.TitleAr,
		Type:      auditFinalReport.AuditFinalReportStatus.Type,
		ColorCode: auditFinalReport.AuditFinalReportStatus.ColorCode,
	}

	// Map audit final report sections
	mainSections := []*auditDTOs.AuditFinalReportSectionDTO{}
	// Sort the sections by the Order field (handling nil values as well)
	sort.Slice(auditFinalReport.AuditFinalReportSections, func(i, j int) bool {
		// Handle nil values by treating them as greater (so they appear later in the sorted list)
		if auditFinalReport.AuditFinalReportSections[i].Order == nil {
			return false
		}
		if auditFinalReport.AuditFinalReportSections[j].Order == nil {
			return true
		}
		return *auditFinalReport.AuditFinalReportSections[i].Order < *auditFinalReport.AuditFinalReportSections[j].Order
	})
	for _, section := range auditFinalReport.AuditFinalReportSections {
		if section.AuditFinalReportSectionID == nil {
			// Main section
			mainSectionDTO := mapSectionToDTO(section)

			// If this section is for findings, map the findings directly to it
			if section.Type == "findings" {
				findingsDTOs := []*auditDTOs.AuditFinalReportFindingDTO{}
				for _, finding := range auditFinalReport.Audit.AuditFindings {
					// Map action plans for the finding
					actionPlansDTOs := []*auditDTOs.AuditFindingActionPlanDTO{}
					for _, actionPlan := range finding.AuditFindingActionPlans {
						var description *string
						if actionPlan.Description != nil {
							description = actionPlan.Description
						}
						// Map the ResponsibleUser
						responsibleUserDTO := mapUserToUserMinimalDTO(actionPlan.ResponsibleUser)

						// Map the ActionPlanStatus
						actionPlanStatusDTO := auditDTOs.AuditFindingActionPlanStatusDTO{
							ID:        actionPlan.ActionPlanStatus.ID,
							TitleEn:   actionPlan.ActionPlanStatus.TitleEn,
							TitleAr:   actionPlan.ActionPlanStatus.TitleAr,
							ColorCode: actionPlan.ActionPlanStatus.ColorCode,
							Type:      actionPlan.ActionPlanStatus.Type,
						}
						actionPlanDTO := auditDTOs.AuditFindingActionPlanDTO{
							ID:                 actionPlan.ID,
							Title:              actionPlan.Title,
							Description:        description,
							ResponsibleUserID:  actionPlan.ResponsibleUserID,
							ResponsibleUser:    responsibleUserDTO,
							TargetDate:         actionPlan.TargetDate.Format("2006-01-02"),
							ActionPlanStatusID: actionPlan.ActionPlanStatusID,
							ActionPlanStatus:   actionPlanStatusDTO,
							CreatedByUser:      mapUserToUserMinimalDTO(actionPlan.CreatedByUser),
							CreatedAt:          actionPlan.CreatedAt.Format("2006-01-02"),
						}
						actionPlansDTOs = append(actionPlansDTOs, &actionPlanDTO)
					}
					var findingDescription *string
					var recommendation *string
					if finding.Description != nil {
						findingDescription = finding.Description
					}
					if finding.Recommendations != nil {
						recommendation = finding.Recommendations
					}
					var departmentResponsibleUserDTO *organizationDTOs.UserMinimalDTO
					if finding.DepartmentResponseByUser != nil {
						mappedUser := mapUserToUserMinimalDTO(finding.DepartmentResponseByUser)
						departmentResponsibleUserDTO = &mappedUser
					}
					var departmentResponseAt *string
					if finding.DepartmentResponseAt != nil {
						formattedDepartmentResponseAt := finding.DepartmentResponseAt.Format("2006-01-02")
						departmentResponseAt = &formattedDepartmentResponseAt
					}
					// Map the finding DTO
					findingDTO := auditDTOs.AuditFinalReportFindingDTO{
						ID:                       finding.ID,
						Title:                    finding.Title,
						Description:              findingDescription,
						SeverityLevelID:          finding.SeverityLevelID,
						DepartmentResponse:       finding.DepartmentResponse,
						DepartmentResponseByID:   finding.DepartmentResponseBYID,
						DepartmentResponseByUser: departmentResponsibleUserDTO,
						DepartmentResponseAt:     departmentResponseAt,
						AuditFindingCategoryID:   finding.AuditFindingCategoryID,
						Recommendations:          recommendation,
						ActionPlans:              actionPlansDTOs,
						CreatedByUser:            mapUserToUserMinimalDTO(finding.CreatedByUser),
						CreatedAt:                finding.CreatedAt.Format("2006-01-02"),
					}
					// Map AuditableItem to DTO
					findingDTO.AuditableItem = auditDTOs.AuditableItemDetailDTO{
						ID:                      finding.AuditableItem.ID,
						Type:                    finding.AuditableItem.Type,
						TitleEn:                 finding.AuditableItem.TitleEn,
						TitleAr:                 finding.AuditableItem.TitleAr,
						Comment:                 &finding.AuditableItem.Comment,
						ProcessID:               finding.AuditableItem.ProcessID,
						ControlID:               finding.AuditableItem.ControlID,
						DocumentID:              finding.AuditableItem.DocumentID,
						IncidentRegisterID:      finding.AuditableItem.IncidentRegisterID,
						ComplianceRequirementID: finding.AuditableItem.ComplianceRequirementID,
						RiskRegisterID:          finding.AuditableItem.RiskRegisterID,
					}
					// Map AuditFindingCategory to DTO
					if finding.AuditFindingCategory != nil {
						findingDTO.AuditFindingCategory = auditDTOs.AuditFindingCategoryDTO{
							ID:            int(finding.AuditFindingCategory.ID),
							TitleEn:       finding.AuditFindingCategory.TitleEn,
							TitleAr:       finding.AuditFindingCategory.TitleAr,
							ColorCode:     finding.AuditFindingCategory.ColorCode,
							DescriptionEn: *finding.AuditFindingCategory.DescriptionEn,
							DescriptionAr: *finding.AuditFindingCategory.DescriptionAr,
						}
					}

					// Map SeverityLevel to DTO
					if finding.SeverityLevel != nil {
						findingDTO.SeverityLevel = masterDTOs.SeverityLevelDTO{
							ID:        finding.SeverityLevel.ID,
							TitleEn:   finding.SeverityLevel.TitleEn,
							TitleAr:   finding.SeverityLevel.TitleAr,
							ColorCode: finding.SeverityLevel.ColorCode,
							Status:    finding.SeverityLevel.Status,
						}
					}
					// Map additional related entities if they exist
					if finding.AuditableItem.Process != nil {
						findingDTO.AuditableItem.Process = &organizationDTOs.ProcessDetailDTO{
							TitleEn: safeStringPointer(&finding.AuditableItem.Process.TitleEn),
							TitleAr: safeStringPointer(&finding.AuditableItem.Process.TitleAr),
						}
					}
					if finding.AuditableItem.Control != nil {
						findingDTO.AuditableItem.Control = &organizationDTOs.ControlDetailDTO{
							TitleEn: safeStringPointer(&finding.AuditableItem.Control.TitleEn),
							TitleAr: safeStringPointer(&finding.AuditableItem.Control.TitleAr),
						}
					}
					if finding.AuditableItem.Document != nil {
						findingDTO.AuditableItem.Document = &documentDTOs.DocumentDetailsDTO{
							TitleEn: safeStringPointer(&finding.AuditableItem.Document.TitleEn),
							TitleAr: safeStringPointer(&finding.AuditableItem.Document.TitleAr),
						}
					}
					if finding.AuditableItem.ComplianceRequirement != nil {
						findingDTO.AuditableItem.ComplianceRequirement = &complianceDTOs.ComplianceRequirementDetailDTO{
							Title: safeStringPointer(&finding.AuditableItem.ComplianceRequirement.Title),
						}
					}
					if finding.AuditableItem.IncidentRegister != nil {
						findingDTO.AuditableItem.IncidentRegister = &incidentDTOs.IncidentRegisterDetailDTO{
							Title: safeStringPointer(&finding.AuditableItem.IncidentRegister.Title),
						}
					}
					if finding.AuditableItem.RiskRegister != nil {
						findingDTO.AuditableItem.RiskRegister = &riskDTOs.RiskRegisterDetailDTO{
							Title: safeStringPointer(&finding.AuditableItem.RiskRegister.Title),
						}
					}
					if finding.AuditableItem.User != nil {
						userDTO := mapUserToUserMinimalDTO(finding.AuditableItem.User)
						findingDTO.AuditableItem.User = &userDTO
					}
					if finding.AuditableItem.ResponsibleUnit != nil {
						findingDTO.AuditableItem.ResponsibleUnit = settingsDTOs.OrganizationUnitDTO{
							ID:        int(finding.AuditableItem.ResponsibleUnit.ID),
							TitleEn:   finding.AuditableItem.ResponsibleUnit.TitleEn,
							TitleAr:   finding.AuditableItem.ResponsibleUnit.TitleAr,
							ColorCode: finding.AuditableItem.ResponsibleUnit.ColorCode,
							Status:    finding.AuditableItem.ResponsibleUnit.Status,
						}
					}
					findingsDTOs = append(findingsDTOs, &findingDTO)
				}
				// Assign findings to the main section DTO
				mainSectionDTO.Findings = findingsDTOs
			}

			// Add subsections (if any)
			subSections := []*auditDTOs.AuditFinalReportSectionDTO{}
			for _, subSection := range auditFinalReport.AuditFinalReportSections {
				if subSection.AuditFinalReportSectionID != nil && *subSection.AuditFinalReportSectionID == section.ID {
					subSections = append(subSections, mapSectionToDTO(subSection))
				}
			}
			mainSectionDTO.SubSections = subSections

			// Add the main section DTO to the main sections list
			mainSections = append(mainSections, mainSectionDTO)
		}
	}

	// Map other fields and create the DTO
	auditReportDTO := &auditDTOs.AuditFinalReportDetailDTO{
		ID:                     auditFinalReport.ID,
		AuditID:                auditFinalReport.AuditID,
		NextReviewLevel:        auditFinalReport.NextReviewLevel,
		Title:                  auditFinalReport.Audit.AuditPlan.Title,
		Description:            auditFinalReport.Audit.AuditPlan.Description,
		ReferenceNumber:        auditFinalReport.Audit.AuditPlan.ReferenceNumber,
		MainSections:           mainSections,
		AuditFinalReportStatus: &auditFinalReportStatusDTO,
		PreparedBy:             createdByUser,
		UpdatedByUser:          updatedByUser,
		PreparedAt:             auditFinalReport.CreatedAt.Format("2006-01-02"),
	}

	// Return the DTO
	return auditReportDTO, nil
}

func (afrs *AuditFinalReportService) GetByAuditID(auditID uint) (*models.AuditFinalReport, error) {
	auditFinalReport, err := afrs.repo.GetByAuditID(auditID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}
	return auditFinalReport, nil
}

func (afrs *AuditFinalReportService) Create(auditID uint, currentUserID uint) (uint, error) {
	tx := afrs.repo.BeginTransaction()
	if _, err := afrs.auditSvc.GetByID(auditID); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if a final report already exists for this audit
	// existingReport, err := afrs.repo.GetByAuditID(auditID)
	// if err != nil {
	//     tx.Rollback()
	//     return 0, err
	// }

	// // If an existing report is found, rollback and return an error
	// if existingReport != nil {
	//     tx.Rollback()
	//     return 0, utility.NewForbiddenError("An audit final report already exists for this audit.")
	// }
	auditDraftReport, err := afrs.auditDraftReportSvc.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if AuditDraftReportStatus is 'completed'
	completedStatus, err := afrs.auditDraftReportStatusSvc.GetByType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if auditDraftReport.AuditDraftReportStatusID != completedStatus.ID {
		tx.Rollback()
		return 0, utility.NewForbiddenError("Action not allowed: Only audit draft reports with an 'Completed' status can be accepted.")
	}
	audit, err := afrs.auditSvc.GetFullDetailsByID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if the current user is the TeamLead
	isAuthorizedUser := false
	if uint(audit.AuditPlan.TeamLead.ID) == currentUserID {
		isAuthorizedUser = true
	}

	// Check if the current user is one of the AuditPlan Members
	if !isAuthorizedUser {
		for _, member := range audit.AuditPlan.Members {
			if member.UserID == currentUserID {
				isAuthorizedUser = true
				break
			}
		}
	}

	// If the current user is not authorized, rollback and return forbidden error
	if !isAuthorizedUser {
		tx.Rollback()
		return 0, utility.NewForbiddenError("The current user is not authorized to perform this action.")
	}
	auditFinalReportStatus, err := afrs.auditFinalReportStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := afrs.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	auditFinalReport := &models.AuditFinalReport{
		AuditID:                  auditID,
		ReferenceNumber:          referenceNumber,
		AuditFinalReportStatusID: auditFinalReportStatus.ID,
		CreatedBYID:              currentUserID,
	}

	id, err := afrs.repo.Create(auditFinalReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Mapping main section titles to types
	mainSections := []struct {
		Title string
		Type  string
	}{
		{"Audit Report Overview", "overview"},
		{"Audit Findings", "finding-summary"},
		{"Findings", "findings"},
		{"Conclusions", "conclusion"},
		{"Appendix", "custom"},
	}

	var reportSections []*models.AuditFinalReportSection
	order := 1
	// Step 1: Insert main sections first
	for _, sectionInfo := range mainSections {
		// Create a new section for each main section
		section := &models.AuditFinalReportSection{
			AuditFinalReportID: id,
			TitleEn:            sectionInfo.Title,
			TitleAr:            sectionInfo.Title,
			Type:               sectionInfo.Type,
			Order:              float64Ptr(float64(order)),
			CreatedBYID:        currentUserID,
		}

		// Add the section to the slice
		reportSections = append(reportSections, section)
		order++
	}
	// Insert the main sections into the database using your repository method
	if err := afrs.repo.CreateAuditFinalReportSections(reportSections, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Step 2: For each main section, insert sub-sections based on their type
	var reportSubSections []*models.AuditFinalReportSection

	for _, section := range reportSections {
		mainSection, err := afrs.repo.GetSectionByType(section.Type, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		subOrder := *section.Order // Use the main section order as the base
		if section.Type == "overview" && auditDraftReport.ExecutiveSummary != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Executive Summary",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.ExecutiveSummary,
				DescriptionAr:             auditDraftReport.ExecutiveSummary,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Audit Objective",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.Audit.AuditPlan.Objective,
				DescriptionAr:             auditDraftReport.Audit.AuditPlan.Objective,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.2),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Audit Procedure",
				TitleAr:                   "ملخص تنفيذي",
				DescriptionEn:             auditDraftReport.Audit.AuditPlan.Purpose,
				DescriptionAr:             auditDraftReport.Audit.AuditPlan.Purpose,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.3),
				CreatedBYID:               currentUserID,
			})
		}
		if section.Type == "overview" && auditDraftReport.Audit.AuditPlan.AuditableItems != nil {
			// Prepare an HTML list for the auditable items
			auditableItemsListEn := "<ul>"
			auditableItemsListAr := "<ul>"

			// Helper function to add items to the list
			addToList := func(titleEn, titleAr string) {
				auditableItemsListEn += fmt.Sprintf("<li>%s</li>", titleEn)
				auditableItemsListAr += fmt.Sprintf("<li>%s</li>", titleAr)
			}

			for _, auditableItem := range auditDraftReport.Audit.AuditPlan.AuditableItems {
				titleEn, titleAr := "", ""

				switch auditableItem.AuditableItem.Type {
				case "process":
					titleEn = auditableItem.AuditableItem.Process.TitleEn
					titleAr = auditableItem.AuditableItem.Process.TitleAr

				case "control":
					titleEn = auditableItem.AuditableItem.Control.TitleEn
					titleAr = auditableItem.AuditableItem.Control.TitleAr

				case "document":
					titleEn = auditableItem.AuditableItem.Document.TitleEn
					titleAr = auditableItem.AuditableItem.Document.TitleAr

				case "incident":
					titleEn = auditableItem.AuditableItem.IncidentRegister.Title
					titleAr = auditableItem.AuditableItem.IncidentRegister.Title

				case "compliance_requirement":
					titleEn = auditableItem.AuditableItem.ComplianceRequirement.Title
					titleAr = auditableItem.AuditableItem.ComplianceRequirement.Title

				case "risk":
					titleEn = auditableItem.AuditableItem.RiskRegister.Title
					titleAr = auditableItem.AuditableItem.RiskRegister.Title

				case "user":
					titleEn = auditableItem.AuditableItem.User.FirstNameEn
					titleAr = auditableItem.AuditableItem.User.FirstNameAr

				case "custom":
					if auditableItem.AuditableItem.TitleEn != nil && auditableItem.AuditableItem.TitleAr != nil {
						titleEn = *auditableItem.AuditableItem.TitleEn
						titleAr = *auditableItem.AuditableItem.TitleAr
					}
				}

				// Only add to the list if titles are not empty
				if titleEn != "" || titleAr != "" {
					addToList(titleEn, titleAr)
				}
			}

			// Close the HTML lists
			auditableItemsListEn += "</ul>"
			auditableItemsListAr += "</ul>"

			// Add a single section with the full list of auditable items in both languages
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Scope of the Audit",
				TitleAr:                   "نطاق التدقيق",
				DescriptionEn:             &auditableItemsListEn,
				DescriptionAr:             &auditableItemsListAr,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.4),
				CreatedBYID:               currentUserID,
			})

		}

		if section.Type == "conclusion" && auditDraftReport.Conclusion != nil {
			reportSubSections = append(reportSubSections, &models.AuditFinalReportSection{
				AuditFinalReportID:        id,
				AuditFinalReportSectionID: &mainSection.ID,
				TitleEn:                   "Concluding Statement",
				TitleAr:                   "خاتمة",
				DescriptionEn:             auditDraftReport.Conclusion,
				DescriptionAr:             auditDraftReport.Conclusion,
				Type:                      section.Type,
				Order:                     float64Ptr(subOrder + 0.1),
				CreatedBYID:               currentUserID,
			})
		}
	}
	// findings only copying from draft findings to AuditFinding, not saving in to AuditFinalReportSection
	auditDraftFindings, err := afrs.auditDraftFindingSvc.GetAllFindingsByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Initialize slices for findings and action plans
	var actionPlans []*models.AuditFindingActionPlan

	for _, draftFinding := range auditDraftFindings {
		if draftFinding.AuditDraftFindingStatus != nil && draftFinding.AuditDraftFindingStatus.IsFinalized {
			var latestComment *string
			var latestCreatedByID *uint
			var latestCreatedAt *time.Time

			// Retrieve the latest conversation if available
			if len(draftFinding.AuditDraftFindingConversations) > 0 {
				latestConversation := draftFinding.AuditDraftFindingConversations[0] // Assumes already ordered by created_at DESC
				latestComment = &latestConversation.Comment
				latestCreatedByID = &latestConversation.CreatedBYID
				latestCreatedAt = &latestConversation.CreatedAt
			}
			// Create AuditFinding
			findingDTO := &auditDTOs.AuditFindingFormDTO{
				AuditID:                draftFinding.AuditID,
				Title:                  draftFinding.Title,
				Description:            draftFinding.Description,
				AuditableItemID:        *draftFinding.AuditableItemID,
				AuditFindingCategoryID: draftFinding.AuditFindingCategoryID,
				SeverityLevelID:        *draftFinding.SeverityLevelID,
				Recommendations:        draftFinding.Recommendations,
				DepartmentResponse:     latestComment,
				DepartmentResponseBYID: latestCreatedByID,
				DepartmentResponseAt:   latestCreatedAt,
			}
			findingID, err := afrs.auditFindingSvc.Create(currentUserID, findingDTO)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			// Copy action plans from draft findings to final findings
			for _, draftActionPlan := range draftFinding.AuditDraftFindingActionPlans {
				actionPlan := &models.AuditFindingActionPlan{
					Title:             draftActionPlan.Title,
					AuditFindingID:    findingID,
					Description:       &draftActionPlan.Description,
					ResponsibleUserID: *draftActionPlan.ResponsibleUserID,
					TargetDate:        draftActionPlan.TargetDate,
					CreatedBYID:       currentUserID,
					// SeverityLevelID:    draftActionPlan.SeverityLevelID,
					ActionPlanStatusID: 1, //todo -- add ActionPlanStatusID
				}
				// Attach action plan to its respective AuditFinding
				actionPlans = append(actionPlans, actionPlan)
			}
		}
	}

	// Bulk insert findings
	if len(actionPlans) > 0 {
		// Bulk insert action plans
		if err := afrs.repo.CreateAuditFindingActionPlans(actionPlans, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	//Bulk insert AuditFinalReportSection
	if len(reportSubSections) > 0 {
		if err := afrs.repo.CreateAuditFinalReportSections(reportSubSections, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	//Add workflow
	err = afrs.auditReportWorkflowSvc.CreateAuditReportWorkflow(id, "final-report", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := afrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinalReport{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinalReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afrs *AuditFinalReportService) Update(auditFinalReport *models.AuditFinalReport) error {
	err := afrs.repo.Update(auditFinalReport)
	if err != nil {
		return err
	}

	return nil
}

func (afrs *AuditFinalReportService) Delete(auditID uint, currentUserID uint) error {
	tx := afrs.repo.BeginTransaction()
	auditFinalReport, err := afrs.repo.GetFullDetailsByID(auditID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect section IDs
	if len(auditFinalReport.AuditFinalReportSections) > 0 {
		var deleteSectionIDs []*uint
		for _, section := range auditFinalReport.AuditFinalReportSections {
			deleteSectionIDs = append(deleteSectionIDs, &section.ID)
		}

		// Delete sections by collected IDs
		if err := afrs.repo.DeleteAuditFinalReportSections(deleteSectionIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Delete related Action Plans for each Audit Finding
	var deleteFindingIDs []*uint
	if len(auditFinalReport.Audit.AuditFindings) > 0 {
		for _, finding := range auditFinalReport.Audit.AuditFindings {
			deleteFindingIDs = append(deleteFindingIDs, &finding.ID)
			// Collect Action Plan IDs for each finding
			if len(finding.AuditFindingActionPlans) > 0 {
				var deleteActionPlanIDs []*uint
				for _, actionPlan := range finding.AuditFindingActionPlans {
					deleteActionPlanIDs = append(deleteActionPlanIDs, &actionPlan.ID)
				}

				// Delete Action Plans by collected IDs
				if err := afrs.repo.DeleteAuditFindingActionPlans(deleteActionPlanIDs, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		// Collect Finding IDs to delete them
		if len(deleteFindingIDs) > 0 {
			// Delete Findings by collected IDs
			if err := afrs.repo.DeleteAuditFindings(deleteFindingIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	//delete workflow
	var deleteAuditReportWorkflowIds []*uint
	var deleteAuditReportWorkflowParticipantIds []*uint
	var deleteAuditReportWorkflowReviewIds []*uint
	if len(auditFinalReport.AuditReportWorkflows) > 0 {
		for _, auditReportWorkflow := range auditFinalReport.AuditReportWorkflows {
			deleteAuditReportWorkflowIds = append(deleteAuditReportWorkflowIds, &auditReportWorkflow.ID)

			if len(auditReportWorkflow.AuditReportWorkflowParticipants) > 0 {
				for _, participant := range auditReportWorkflow.AuditReportWorkflowParticipants {
					deleteAuditReportWorkflowParticipantIds = append(deleteAuditReportWorkflowParticipantIds, &participant.ID)
				}
			}

			if len(auditReportWorkflow.AuditReportReviews) > 0 {
				for _, reviews := range auditReportWorkflow.AuditReportReviews {
					deleteAuditReportWorkflowReviewIds = append(deleteAuditReportWorkflowReviewIds, &reviews.ID)
				}
			}
		}
	}
	if err := afrs.auditReportWorkflowSvc.DeleteAuditReportWorkflowParticipants(deleteAuditReportWorkflowParticipantIds, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := afrs.auditReportWorkflowSvc.DeleteAuditReportReviews(deleteAuditReportWorkflowReviewIds, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := afrs.auditReportWorkflowSvc.DeleteAuditReportWorkflows(deleteAuditReportWorkflowIds, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := afrs.repo.Delete(auditFinalReport.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		err := afrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditFinalReport{}.TableName(),
			auditFinalReport.ID,
			nil,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.AuditFinalReport{}.TableName(),
				"recordID":  auditFinalReport.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

// Helper function to return a pointer to a float64 value
func float64Ptr(val float64) *float64 {
	return &val
}
func mapSectionToDTO(section *models.AuditFinalReportSection) *auditDTOs.AuditFinalReportSectionDTO {
	return &auditDTOs.AuditFinalReportSectionDTO{
		ID:                        section.ID,
		TitleEn:                   section.TitleEn,
		TitleAr:                   section.TitleAr,
		DescriptionEn:             section.DescriptionEn,
		DescriptionAr:             section.DescriptionAr,
		Type:                      section.Type,
		Order:                     section.Order,
		AuditFinalReportSectionID: section.AuditFinalReportSectionID,
		CreatedByUser:             mapUserToUserMinimalDTO(section.CreatedByUser),
		// UpdatedByUser: mapUserToUserMinimalDTO(section.UpdatedByUser),
		CreatedAt: section.CreatedAt.Format("2006-01-02"),
		// UpdatedAt: section.UpdatedAt.Format("2006-01-02"),
	}
}
