package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRequirementServiceInterface interface {
	GetAll(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error)
	GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
	ExportComplianceRequirements(data []*complianceDTOs.ComplianceRequirementListDTO, currentUserID uint) (*os.File, error)

	GenerateRequirementTemplate(complianceID uint) (*os.File, error)
	ValidateRequirementTemplate(complianceID uint, file io.Reader) ([]map[string]interface{}, error)
	ValidateImportRequirement(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportRequirements(complianceID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error)

	// Inherent Risk Analysis methods
	CreateInherentRiskAnalysis(complianceRequirementID uint, currentUserID uint, riskAnalysisFormDTO []*complianceDTOs.ComplianceRequirementInherentRiskFormDTO) error
	DeleteInherentRiskAnalysis(complianceRequirementID uint, currentUserID uint) error
}

// Service and implementations
type ComplianceRequirementService struct {
	repo                             complianceRepositories.ComplianceRequirementRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	complianceRegisterSvc            ComplianceRegisterServiceInterface
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
	severityLevelSvc                 masterServices.SeverityLevelServiceInterface
	controlSvc                       organizationServices.ControlServiceInterface
	strategicObjectiveSvc            organizationServices.StrategicObjectiveServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
	userSettingSvc                   organizationServices.UserSettingServiceInterface
	complianceClauseSvc              ComplianceClauseServiceInterface
	requiredDocumentTitleSvc         masterServices.RequiredDocumentTitleServiceInterface
	organizationUnitSvc              settingsServices.OrganizationUnitServiceInterface
}

func NewComplianceRequirementService(
	repo complianceRepositories.ComplianceRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	strategicObjectiveSvc organizationServices.StrategicObjectiveServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	complianceClauseSvc ComplianceClauseServiceInterface,
	requiredDocumentTitleSvc masterServices.RequiredDocumentTitleServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
) ComplianceRequirementServiceInterface {
	return &ComplianceRequirementService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		complianceRegisterSvc:            complianceRegisterSvc,
		complianceRequirementCategorySvc: complianceRequirementCategorySvc,
		complianceStatusSvc:              complianceStatusSvc,
		severityLevelSvc:                 severityLevelSvc,
		controlSvc:                       controlSvc,
		strategicObjectiveSvc:            strategicObjectiveSvc,
		userSvc:                          userSvc,
		notificationSvc:                  notificationSvc,
		userSettingSvc:                   userSettingSvc,
		complianceClauseSvc:              complianceClauseSvc,
		requiredDocumentTitleSvc:         requiredDocumentTitleSvc,
		organizationUnitSvc:              organizationUnitSvc,
	}
}

func (crs *ComplianceRequirementService) GetAll(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	if complianceID != 0 {
		if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
			return nil, err
		}
	}
	return crs.repo.GetAll(complianceID, params)
}

func (crs *ComplianceRequirementService) GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	// Step 1: Fetch paged requirements + total count
	requirements, totalCount, err := crs.repo.GetRequirementsByComplianceID(complianceID, params)
	if err != nil {
		return nil, err
	}
	if len(requirements) == 0 {
		return &complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO{
			Items:      []*complianceDTOs.ComplianceRequirementByComplianceIDListDTO{},
			Pagination: commonDTOs.PaginationInfoDTO{CurrentPage: params.Page, PageSize: params.PageSize},
		}, nil
	}

	// Step 2: Collect requirement IDs
	reqIDs := make([]uint, 0, len(requirements))
	for _, r := range requirements {
		reqIDs = append(reqIDs, r.ID)
	}

	// Step 3: Fetch responsible units
	var unitRows []complianceDTOs.ComplianceRequirementResponsibleUnitRowDTO
	err = crs.repo.GetResponsibleUnitRowsByRequirementIDs(reqIDs, &unitRows)
	if err != nil {
		return nil, err
	}

	// Step 4: Map responsible units
	unitMap := make(map[uint][]*complianceDTOs.ComplianceRequirementResponsibleUnitListDTO)
	for _, row := range unitRows {
		unit := &settingsDTOs.OrganizationUnitDTO{
			ID:        int(row.ResponsibleUnitID),
			TitleEn:   row.TitleEn,
			TitleAr:   row.TitleAr,
			Status:    row.Status,
			ColorCode: row.ColorCode,
		}
		if row.ParentID != nil {
			unit.ParentUnit = &settingsDTOs.OrganizationUnitDTO{ID: int(*row.ParentID)}
		}
		unitDTO := &complianceDTOs.ComplianceRequirementResponsibleUnitListDTO{
			ID:                                   row.ID,
			ComplianceRequirementID:              row.ComplianceRequirementID,
			ResponsibleUnitID:                    row.ResponsibleUnitID,
			IsRegulationNotificationAssessment:   row.IsRegulationNotificationAssessment,
			IsRiskBasedAssessment:                row.IsRiskBasedAssessment,
			ComplianceAssessmentFrameworkValueID: row.ComplianceAssessmentFrameworkValueID,
			ResponsibleUnit:                      unit,
		}
		unitMap[row.ComplianceRequirementID] = append(unitMap[row.ComplianceRequirementID], unitDTO)
	}

	// Step 5: Assign responsible units
	for _, r := range requirements {
		if units, ok := unitMap[r.ID]; ok {
			r.ResponsibleUnits = units
		} else {
			r.ResponsibleUnits = []*complianceDTOs.ComplianceRequirementResponsibleUnitListDTO{}
		}
	}

	// Step 6: Group into clauses (preserve DB order)
	clauseIndex := make(map[uint]int)
	clauseList := make([]*complianceDTOs.ComplianceRequirementByComplianceIDListDTO, 0, len(requirements))

	for _, r := range requirements { // already ordered by SQL
		idx, ok := clauseIndex[r.ClauseID]
		if !ok {
			idx = len(clauseList)
			clauseIndex[r.ClauseID] = idx
			clauseList = append(clauseList, &complianceDTOs.ComplianceRequirementByComplianceIDListDTO{
				ClauseID:      r.ClauseID,
				ClauseNo:      r.ComplianceClauseNo,
				ClauseTitleEn: r.ComplianceClauseTitleEn,
				ClauseTitleAr: r.ComplianceClauseTitleAr,
				Requirements:  []*complianceDTOs.ComplianceRequirementListDTO{},
			})
		}
		clauseList[idx].Requirements = append(clauseList[idx].Requirements, r)
	}

	// Step 7: Build response with correct pagination
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	return &complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO{
		Items: clauseList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (crs ComplianceRequirementService) GetByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRequirement, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceRequirement, err := crs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error) {
	complianceRequirement, err := crs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetFullDetailsByID(complianceID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRequirementDetailDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var complianceRequirementDTO complianceDTOs.ComplianceRequirementDetailDTO

	complianceRequirement, err := crs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to DTO
	err = mapstructure.Decode(complianceRequirement, &complianceRequirementDTO)
	if err != nil {
		return nil, err
	}

	if complianceRequirement.Deadline != nil {
		deadline := complianceRequirement.Deadline.Format("2006-01-02")
		complianceRequirementDTO.DeadlineString = &deadline
	}

	return &complianceRequirementDTO, nil
}

func (crs *ComplianceRequirementService) Create(complianceID uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	_, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return 0, err
	}

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	_, err = crs.complianceClauseSvc.GetByID(complianceRequirementFormDTO.ComplianceClauseID)
	if err != nil {
		return 0, err
	}

	complianceRequirement := &models.ComplianceRequirement{
		ComplianceRequirementCategoryID: complianceRequirementFormDTO.ComplianceRequirementCategoryID,
		// ComplianceStatusID:              complianceRequirementFormDTO.ComplianceStatusID,
		SeverityLevelID:    complianceRequirementFormDTO.SeverityLevelID,
		TitleEn:            complianceRequirementFormDTO.TitleEn,
		TitleAr:            complianceRequirementFormDTO.TitleAr,
		Deadline:           dates.Deadline,
		CreatedByID:        currentUserID,
		ComplianceClauseID: complianceRequirementFormDTO.ComplianceClauseID,
	}

	id, err := crs.repo.Create(complianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// complianceStatusID := complianceRequirementFormDTO.ComplianceStatusID
	// complianceRequirementUpdate := &models.ComplianceRequirementStatusUpdate{
	// 	ComplianceRequirementID: id,
	// 	ComplianceStatusID:      *complianceStatusID,
	// 	Comment:                 "Initial compliance status",
	// 	CreatedByID:             currentUserID,
	// }

	// _, err = crs.repo.CreateInitialComplianceRequirementStatusUpdate(complianceRequirementUpdate, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// averageScore, err := crs.repo.CalculateAverageScore(complianceID, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// Update the ComplianceRegister with the new total score
	// err = crs.complianceRegisterSvc.UpdateComplianceTotalScore(complianceID, averageScore, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }
	var responsibleUnits []*models.ComplianceRequirementResponsibleUnit
	for _, resUnit := range complianceRequirementFormDTO.ResponsibleUnits {
		unit := &models.ComplianceRequirementResponsibleUnit{
			ComplianceRequirementID:            id,
			ResponsibleUnitID:                  resUnit.ResponsibleUnitID,
			IsRegulationNotificationAssessment: resUnit.IsRegulationNotificationAssessment,
			IsRiskBasedAssessment:              resUnit.IsRiskBasedAssessment,
			CreatedByID:                        currentUserID,
		}
		responsibleUnits = append(responsibleUnits, unit)
	}

	if len(responsibleUnits) > 0 {
		if err := crs.repo.CreateComplianceRequirementResponsibleUnits(responsibleUnits, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add Controls
	var controls []*models.ComplianceControl
	for _, controlID := range complianceRequirementFormDTO.ControlIDs {
		control := &models.ComplianceControl{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: &id,
			ControlID:               *controlID,
			ComplianceClauseID:      &complianceRequirementFormDTO.ComplianceClauseID,
			CreatedByID:             currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add Strategic Objectives
	var strategicObjectives []*models.ComplianceRequirementStrategicObjective
	for _, strategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		strategicObjective := &models.ComplianceRequirementStrategicObjective{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: &id,
			StrategicObjectiveID:    *strategicObjectiveID,
			CreatedByID:             currentUserID,
		}
		strategicObjectives = append(strategicObjectives, strategicObjective)
	}

	if len(strategicObjectives) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectives, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.ComplianceRequirementAssociatedDocument
	for _, documentID := range complianceRequirementFormDTO.DocumentIds {
		document := &models.ComplianceRequirementAssociatedDocument{
			ComplianceRegisterID:    complianceID,
			ComplianceRequirementID: id,
			DocumentID:              *documentID,
			CreatedByID:             currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := crs.repo.CreateComplianceRequirementDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var complianceRequirementRequiredDocumentTitle []*models.ComplianceRequirementRequiredDocumentTitle
	for _, requiredDocumentTitleID := range complianceRequirementFormDTO.RequiredDocumentTitleIDs {
		crrDocumentTitle := &models.ComplianceRequirementRequiredDocumentTitle{
			ComplianceRequirementID: id,
			RequiredDocumentTitleID: *requiredDocumentTitleID,
			CreatedBYID:             currentUserID,
		}
		complianceRequirementRequiredDocumentTitle = append(complianceRequirementRequiredDocumentTitle, crrDocumentTitle)
	}

	if len(complianceRequirementRequiredDocumentTitle) > 0 {
		if err := crs.repo.CreateComplianceRequirmentRequiredDocumentTitles(complianceRequirementRequiredDocumentTitle, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingComplianceRequirement, err := crs.GetFullDetailsByID(complianceID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			existingComplianceRequirement,
			currentUserID,
			"compliance", "compliance_register", "requirement", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRequirement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// prepare notification
		// urlParam := fmt.Sprintf(`{"compliance_register_id":%d,"requirement_id":%d}`, complianceID, id)
		// queryParam := ""
		// entity := "compliance-requirement"
		// action := "responsible-user-assign"

		// dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 	UserIDs:              []uint{currentUserID},
		// 	ModuleIdentifier:     "compliance",
		// 	NotificationType:     string(commonDTOs.AssignmentType),
		// 	HeadingEn:            "Compliance Requirement Responsiblity Assigned",
		// 	HeadingAr:            "تم تعيين مسؤولية متطلبات الامتثال",
		// 	TitleEn:              complianceRequirementFormDTO.TitleEn,
		// 	TitleAr:              complianceRequirementFormDTO.TitleAr,
		// 	ReferenceNumber:      compliance.ReferenceNumber,
		// 	Entity:               &entity,
		// 	Action:               &action,
		// 	URLParam:             urlParam,
		// 	QueryParam:           queryParam,
		// 	IsSystemNotification: false,
		// 	CreatedByID:          &currentUserID,
		// }

		// responsibleUser, err := crs.userSvc.GetByID(currentUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch responsible user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }
		// assignedBy, err := crs.userSvc.GetByID(currentUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch assigning user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// complianceRequirementCategory, err := crs.complianceRequirementCategorySvc.GetByID(complianceRequirementFormDTO.ComplianceRequirementCategoryID)
		// if err != nil {
		// 	return
		// }

		// emailTo := []string{responsibleUser.Email}
		// subject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Assigned", complianceRequirementFormDTO.TitleEn)
		// appURL := utility.GetAppURL()
		// deadline := ""
		// if dates.Deadline != nil {
		// 	deadline = dates.Deadline.Format("2006-01-02")
		// }
		// emailData := map[string]interface{}{
		// 	"FirstNameEn":                   responsibleUser.FirstNameEn,
		// 	"LastNameEn":                    responsibleUser.LastNameEn,
		// 	"ComplianceReferenceNumber":     compliance.ReferenceNumber,
		// 	"ComplianceTitleEn":             compliance.TitleEn,
		// 	"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
		// 	"RequirementTitleEn":            complianceRequirementFormDTO.TitleEn,
		// 	"ComplianceClauseTitle":         clause.TitleEn,
		// 	"Deadline":                      deadline,
		// 	"AssignedByFirstNameEn":         assignedBy.FirstNameEn,
		// 	"AssignedByLastNameEn":          assignedBy.LastNameEn,
		// 	"AssignedOn":                    time.Now().Format("2006-01-02"),
		// 	"RedirectURL":                   fmt.Sprintf("%s/compliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
		// }

		// emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 	EmailTo:              emailTo,
		// 	Subject:              subject,
		// 	TemplateRelativePath: "compliance/templates/email/compliance_requirement_create_responsible_user_mail.html",
		// 	EmailData:            emailData,
		// }

		// // send notification and mail
		// err = crs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_requirement")
		// if err != nil {
		// 	utility.LogError("Failed to send notifications", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// }
	}()

	return id, nil
}

func (crs *ComplianceRequirementService) Update(complianceID uint, id uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	tx := crs.repo.BeginTransaction()

	_, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return err
	}

	existingComplianceRequirement, err := crs.GetByID(complianceID, id, tx)
	if err != nil {
		return err
	}

	oldComplianceRequirement, err := crs.GetFullDetailsByID(complianceID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := crs.validateComplianceRequirementFormDTO(complianceRequirementFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRequirementDates(complianceRequirementFormDTO)
	if err != nil {
		return err
	}

	clause, err := crs.complianceClauseSvc.GetByID(complianceRequirementFormDTO.ComplianceClauseID)
	if err != nil {
		return err
	}

	// Update Compliance Requirement
	existingComplianceRequirement.ComplianceClauseID = clause.ID
	existingComplianceRequirement.ComplianceRequirementCategoryID = complianceRequirementFormDTO.ComplianceRequirementCategoryID
	// existingComplianceRequirement.ComplianceStatusID = complianceRequirementFormDTO.ComplianceStatusID
	existingComplianceRequirement.SeverityLevelID = complianceRequirementFormDTO.SeverityLevelID
	existingComplianceRequirement.TitleEn = complianceRequirementFormDTO.TitleEn
	existingComplianceRequirement.TitleAr = complianceRequirementFormDTO.TitleAr
	existingComplianceRequirement.Deadline = dates.Deadline
	existingComplianceRequirement.UpdatedByID = &currentUserID

	existingComplianceRequirement.ComplianceClause = nil // Reset preloaded association
	err = crs.repo.Update(existingComplianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Controls
	complianceRequirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Get existing responsible unit IDs
	existingUnitIDs := make(map[uint]uint) // map[responsibleUnitID]recordID
	existingUnitAssessmentMap := make(map[uint]*complianceDTOs.ComplianceRequirementResponsibleUnitDTO)
	for _, unit := range complianceRequirement.ResponsibleUnits {
		if unit.ID != 0 && unit.ResponsibleUnitID != 0 {
			existingUnitIDs[unit.ResponsibleUnitID] = unit.ID
			existingUnitAssessmentMap[unit.ResponsibleUnitID] = unit
		}
	}

	// Get new unit IDs from form
	newUnitMap := make(map[uint]*complianceDTOs.ComplianceRequirementResponsibleUnitFormDTO)
	for _, unit := range complianceRequirementFormDTO.ResponsibleUnits {
		newUnitMap[unit.ResponsibleUnitID] = unit
	}

	// Add new responsible units
	var unitsToAdd []*models.ComplianceRequirementResponsibleUnit
	var unitsToUpdate []*models.ComplianceRequirementResponsibleUnit
	for _, unit := range complianceRequirementFormDTO.ResponsibleUnits {
		if _, exists := existingUnitIDs[unit.ResponsibleUnitID]; !exists {
			unitsToAdd = append(unitsToAdd, &models.ComplianceRequirementResponsibleUnit{
				ComplianceRequirementID:            id,
				ResponsibleUnitID:                  unit.ResponsibleUnitID,
				IsRegulationNotificationAssessment: unit.IsRegulationNotificationAssessment,
				IsRiskBasedAssessment:              unit.IsRiskBasedAssessment,
				CreatedByID:                        currentUserID,
			})
			continue
		}

		existingUnit := existingUnitAssessmentMap[unit.ResponsibleUnitID]
		if existingUnit != nil &&
			((existingUnit.IsRegulationNotificationAssessment == nil && unit.IsRegulationNotificationAssessment != nil) ||
				(existingUnit.IsRegulationNotificationAssessment != nil && unit.IsRegulationNotificationAssessment == nil) ||
				(existingUnit.IsRegulationNotificationAssessment != nil && unit.IsRegulationNotificationAssessment != nil && *existingUnit.IsRegulationNotificationAssessment != *unit.IsRegulationNotificationAssessment) ||
				(existingUnit.IsRiskBasedAssessment == nil && unit.IsRiskBasedAssessment != nil) ||
				(existingUnit.IsRiskBasedAssessment != nil && unit.IsRiskBasedAssessment == nil) ||
				(existingUnit.IsRiskBasedAssessment != nil && unit.IsRiskBasedAssessment != nil && *existingUnit.IsRiskBasedAssessment != *unit.IsRiskBasedAssessment)) {
			recordID := existingUnitIDs[unit.ResponsibleUnitID]
			unitsToUpdate = append(unitsToUpdate, &models.ComplianceRequirementResponsibleUnit{
				ID:                                 recordID,
				ComplianceRequirementID:            id,
				ResponsibleUnitID:                  unit.ResponsibleUnitID,
				IsRegulationNotificationAssessment: unit.IsRegulationNotificationAssessment,
				IsRiskBasedAssessment:              unit.IsRiskBasedAssessment,
				UpdatedByID:                        &currentUserID,
			})
		}
	}
	if len(unitsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementResponsibleUnits(unitsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(unitsToUpdate) > 0 {
		if err := crs.repo.UpdateComplianceRequirementResponsibleUnits(unitsToUpdate, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete removed responsible units
	var unitIDsToDelete []*uint
	for unitID, recordID := range existingUnitIDs {
		if _, exists := newUnitMap[unitID]; !exists {
			idToDelete := uint(recordID)
			unitIDsToDelete = append(unitIDsToDelete, &idToDelete)
		}
	}
	if len(unitIDsToDelete) > 0 {
		if err := crs.repo.DeleteComplianceRequirementResponsibleUnits(unitIDsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirement.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRequirementFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ComplianceControl
	for _, newControlID := range complianceRequirementFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new compliance control record
			control := &models.ComplianceControl{
				ComplianceRegisterID:    complianceID,
				ComplianceRequirementID: &id,
				ControlID:               *newControlID,
				ComplianceClauseID:      &clause.ID,
				CreatedByID:             currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range complianceRequirement.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Strategic Objectives
	existingStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirement.StrategicObjectives {
		existingStrategicObjectiveIDs[obj.StrategicObjectiveID] = true
	}

	// Collect new strategic objective IDs
	newStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		newStrategicObjectiveIDs[*obj] = true
	}

	// Add new strategic objectives
	var strategicObjectivesToAdd []*models.ComplianceRequirementStrategicObjective
	for _, newStrategicObjectiveID := range complianceRequirementFormDTO.StrategicObjectiveIDs {
		if _, exists := existingStrategicObjectiveIDs[*newStrategicObjectiveID]; !exists {
			// Create a new compliance strategic objective record
			strategicObjective := &models.ComplianceRequirementStrategicObjective{
				ComplianceRequirementID: &id,
				StrategicObjectiveID:    *newStrategicObjectiveID,
				CreatedByID:             currentUserID,
			}
			strategicObjectivesToAdd = append(strategicObjectivesToAdd, strategicObjective)
		}
	}

	// Save the new strategic objectives
	if len(strategicObjectivesToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementStrategicObjectives(strategicObjectivesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove strategic objectives not in the new list
	var deleteStrategicObjectiveIDs []*uint

	for _, existingStrategicObjective := range complianceRequirement.StrategicObjectives {
		if _, exists := newStrategicObjectiveIDs[existingStrategicObjective.StrategicObjectiveID]; !exists {
			// Append strategic objective to delete list
			deleteStrategicObjectiveIDs = append(deleteStrategicObjectiveIDs, existingStrategicObjective.ID)
		}
	}

	// Delete strategic objectives not in the new list
	if len(deleteStrategicObjectiveIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range complianceRequirement.Documents {
		existingDocumentIDs[doc.DocumentID] = true
	}

	// Collect new document IDs
	newDocumentIDs := make(map[uint]bool)
	for _, doc := range complianceRequirementFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	// Add new documents
	var documentsToAdd []*models.ComplianceRequirementAssociatedDocument
	for _, newDocumentID := range complianceRequirementFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			// Create a new compliance requirement document record
			document := &models.ComplianceRequirementAssociatedDocument{
				ComplianceRegisterID:    complianceID,
				ComplianceRequirementID: id,
				DocumentID:              *newDocumentID,
				CreatedByID:             currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	// Save the new documents
	if len(documentsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirementDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove documents not in the new list
	var deleteDocumentIDs []*uint

	for _, existingDocument := range complianceRequirement.Documents {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			// Append document to delete list
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	// averageScore, err := crs.repo.CalculateAverageScore(complianceID, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// err = crs.complianceRegisterSvc.UpdateComplianceTotalScore(complianceID, averageScore, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// Delete documents not in the new list
	if len(deleteDocumentIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirementDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Required Document Titles
	existingDocumentTitleIDs := make(map[uint]bool)
	for _, doc := range oldComplianceRequirement.RequiredDocumentTitles {
		existingDocumentTitleIDs[doc.RequiredDocumentTitleID] = true
	}

	// Collect new required document title IDs
	newDocumentTitleIDs := make(map[uint]bool)
	for _, docID := range complianceRequirementFormDTO.RequiredDocumentTitleIDs {
		newDocumentTitleIDs[*docID] = true
	}

	// Add new required document titles
	var requiredDocumentsToAdd []*models.ComplianceRequirementRequiredDocumentTitle
	for _, newDocumentTitleID := range complianceRequirementFormDTO.RequiredDocumentTitleIDs {
		if _, exists := existingDocumentTitleIDs[*newDocumentTitleID]; !exists {
			// Create a new compliance requirement document record
			requiredDocument := &models.ComplianceRequirementRequiredDocumentTitle{
				ComplianceRequirementID: id,
				RequiredDocumentTitleID: *newDocumentTitleID,
				CreatedBYID:             currentUserID,
			}
			requiredDocumentsToAdd = append(requiredDocumentsToAdd, requiredDocument)
		}
	}

	// Save the new required documents
	if len(requiredDocumentsToAdd) > 0 {
		if err := crs.repo.CreateComplianceRequirmentRequiredDocumentTitles(requiredDocumentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove required document titles not in the new list
	var deleteDocumentTitleIDs []*uint
	for _, existingDocument := range oldComplianceRequirement.RequiredDocumentTitles {
		if _, exists := newDocumentTitleIDs[existingDocument.RequiredDocumentTitleID]; !exists {
			// Append document to delete list
			deleteDocumentTitleIDs = append(deleteDocumentTitleIDs, existingDocument.ID)
		}
	}

	// Delete required document titles not in the new list
	if len(deleteDocumentTitleIDs) > 0 {
		if err := crs.repo.DeleteComplianceRequirmentRequiredDocumentTitles(deleteDocumentTitleIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingRequirement, err := crs.GetFullDetailsByID(complianceID, id)
		if err != nil {
			return
		}
		// Determine control
		controlsAdded := []uint{}
		controlsRemoved := []uint{}
		for id := range newControlIDs {
			if !existingControlIDs[id] {
				controlsAdded = append(controlsAdded, id)
			}
		}
		for id := range existingControlIDs {
			if !newControlIDs[id] {
				controlsRemoved = append(controlsRemoved, id)
			}
		}

		// Determine strategic objectives
		strategicObjectivesAdded := []uint{}
		strategicObjectivesRemoved := []uint{}
		for id := range newStrategicObjectiveIDs {
			if !existingStrategicObjectiveIDs[id] {
				strategicObjectivesAdded = append(strategicObjectivesAdded, id)
			}
		}
		for id := range existingStrategicObjectiveIDs {
			if !newStrategicObjectiveIDs[id] {
				strategicObjectivesRemoved = append(strategicObjectivesRemoved, id)
			}
		}

		// Determine strategic objectives
		documentsAdded := []uint{}
		documentssRemoved := []uint{}
		for id := range newDocumentIDs {
			if !existingDocumentIDs[id] {
				documentsAdded = append(documentsAdded, id)
			}
		}
		for id := range existingDocumentIDs {
			if !newDocumentIDs[id] {
				documentssRemoved = append(documentssRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(controlsAdded) > 0 || len(controlsRemoved) > 0 {
			associatedItems["control_ids"] = map[string]interface{}{
				"added":   controlsAdded,
				"removed": controlsRemoved,
			}
		}

		if len(strategicObjectivesAdded) > 0 || len(strategicObjectivesRemoved) > 0 {
			associatedItems["strategic_objective_ids"] = map[string]interface{}{
				"added":   strategicObjectivesAdded,
				"removed": strategicObjectivesRemoved,
			}
		}

		if len(documentsAdded) > 0 || len(documentssRemoved) > 0 {
			associatedItems["document_ids"] = map[string]interface{}{
				"added":   documentsAdded,
				"removed": documentssRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		// log changes
		oldData := map[string]interface{}{
			"title":             oldComplianceRequirement.TitleEn,
			"dead_line":         oldComplianceRequirement.DeadlineString,
			"compliance_clause": oldComplianceRequirement.ComplianceClause.TitleEn,
		}
		// Add category title if exists
		if oldComplianceRequirement.ComplianceRequirementCategory != nil {
			oldData["compliance_requirement_category"] = oldComplianceRequirement.ComplianceRequirementCategory.TitleEn
		}

		// Add severity level title if exists
		if oldComplianceRequirement.SeverityLevel != nil {
			oldData["severity_level"] = oldComplianceRequirement.SeverityLevel.TitleEn
		}

		newData := map[string]interface{}{
			"title":             existingRequirement.TitleEn,
			"dead_line":         existingRequirement.DeadlineString,
			"compliance_clause": existingRequirement.ComplianceClause.TitleEn,
		}
		// Add category title if exists
		if existingRequirement.ComplianceRequirementCategory != nil {
			newData["compliance_requirement_category"] = existingRequirement.ComplianceRequirementCategory.TitleEn
		}
		// Add severity level title if exists
		if existingRequirement.SeverityLevel != nil {
			newData["severity_level"] = existingRequirement.SeverityLevel.TitleEn
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := crs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceRequirement{}.TableName(),
				id,
				changedData,
				currentUserID,
				"compliance", "compliance_register", "requirement", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceRequirement{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
		// Send notifications if the responsible user has changed
		// if oldResponsibleUserID != existingComplianceRequirement.ResponsibleUserID {
		// 	// URL param for notifications
		// 	urlParam := fmt.Sprintf(`{"compliance_register_id":%d,"requirement_id":%d}`, complianceID, id)
		// 	queryParam := ""
		// 	entity := "compliance-requirement"

		// 	// Old responsible user notification (Removal)
		// 	oldResponsibleUserAction := "responsible-user-remove"
		// 	oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 		UserIDs:              []uint{oldResponsibleUserID},
		// 		ModuleIdentifier:     "compliance",
		// 		NotificationType:     string(commonDTOs.InformationType),
		// 		HeadingEn:            "Compliance Requirement Responsibility Removed",
		// 		HeadingAr:            "تمت إزالة مسؤولية متطلبات الامتثال",
		// 		TitleEn:              existingComplianceRequirement.TitleEn,
		// 		ReferenceNumber:      compliance.ReferenceNumber,
		// 		TitleAr:              existingComplianceRequirement.TitleEn,
		// 		Entity:               &entity,
		// 		Action:               &oldResponsibleUserAction,
		// 		URLParam:             urlParam,
		// 		QueryParam:           queryParam,
		// 		IsSystemNotification: false,
		// 		CreatedByID:          &currentUserID,
		// 	}

		// 	// New responsible user notification (Assignment)
		// 	newResponsibleUserAction := "responsible-user-assign"
		// 	newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 		UserIDs:              []uint{currentUserID},
		// 		ModuleIdentifier:     "compliance",
		// 		NotificationType:     string(commonDTOs.AssignmentType),
		// 		HeadingEn:            "Compliance Requirement Responsibility Assigned",
		// 		HeadingAr:            "تم تعيين المسؤولية لمتطلبات الامتثال",
		// 		TitleEn:              existingComplianceRequirement.TitleEn,
		// 		TitleAr:              existingComplianceRequirement.TitleEn,
		// 		ReferenceNumber:      compliance.ReferenceNumber,
		// 		Entity:               &entity,
		// 		Action:               &newResponsibleUserAction,
		// 		URLParam:             urlParam,
		// 		QueryParam:           queryParam,
		// 		IsSystemNotification: false,
		// 		CreatedByID:          &currentUserID,
		// 	}

		// 	// Fetch the responsible users' data
		// 	oldResponsibleUser, err := crs.userSvc.GetByID(oldResponsibleUserID)
		// 	if err != nil {
		// 		utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
		// 			"error": err.Error(),
		// 		})
		// 		return
		// 	}

		// 	newResponsibleUser, err := crs.userSvc.GetByID(currentUserID)
		// 	if err != nil {
		// 		utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
		// 			"error": err.Error(),
		// 		})
		// 		return
		// 	}

		// 	// Fetch details of the user making the assignment
		// 	updatedByUser, err := crs.userSvc.GetByID(currentUserID)
		// 	if err != nil {
		// 		utility.LogError("Failed to fetch assigning user", map[string]interface{}{
		// 			"error": err.Error(),
		// 		})
		// 		return
		// 	}

		// 	complianceRequirementCategory, err := crs.complianceRequirementCategorySvc.GetByID(complianceRequirementFormDTO.ComplianceRequirementCategoryID)
		// 	if err != nil {
		// 		return
		// 	}

		// 	appURL := utility.GetAppURL()

		// 	// Email notification for the old responsible user
		// 	oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
		// 	oldResponsibleUserSubject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Removed", existingComplianceRequirement.TitleEn)
		// 	oldResponsibleUserEmailData := map[string]interface{}{
		// 		"FirstNameEn":                   oldResponsibleUser.FirstNameEn,
		// 		"LastNameEn":                    oldResponsibleUser.LastNameEn,
		// 		"ComplianceReferenceNumber":     compliance.ReferenceNumber,
		// 		"ComplianceTitle":               compliance.TitleEn,
		// 		"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
		// 		"RequirementTitleEn":            existingComplianceRequirement.TitleEn,
		// 		"ComplianceClauseTitle":         clause.TitleEn,
		// 		"Deadline":                      dates.Deadline.Format("2006-01-02"),
		// 		"UpdatedByFirstNameEn":          updatedByUser.FirstNameEn,
		// 		"UpdatedByLastNameEn":           updatedByUser.LastNameEn,
		// 		"AssignedOn":                    time.Now().Format("2006-01-02"),
		// 		"RedirectURL":                   fmt.Sprintf("%s/compliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
		// 	}

		// 	oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 		EmailTo:              oldResponsibleUserEmailTo,
		// 		Subject:              oldResponsibleUserSubject,
		// 		TemplateRelativePath: "compliance/templates/email/compliance_requirement_remove_responsible_user_mail.html",
		// 		EmailData:            oldResponsibleUserEmailData,
		// 	}

		// 	// Email notification for the new responsible user
		// 	newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
		// 	newResponsibleUserSubject := fmt.Sprintf("Compliance Requirement/%s: Responsibility Assigned", existingComplianceRequirement.TitleEn)
		// 	newResponsibleUserEmailData := map[string]interface{}{
		// 		"FirstNameEn":                   newResponsibleUser.FirstNameEn,
		// 		"LastNameEn":                    newResponsibleUser.LastNameEn,
		// 		"ComplianceReferenceNumber":     compliance.ReferenceNumber,
		// 		"ComplianceTitle":               compliance.TitleEn,
		// 		"ComplianceRequirementCategory": complianceRequirementCategory.TitleEn,
		// 		"RequirementTitleEn":            existingComplianceRequirement.TitleEn,
		// 		"ComplianceClauseTitle":         clause.TitleEn,
		// 		"Deadline":                      dates.Deadline.Format("2006-01-02"),
		// 		"UpdatedByFirstNameEn":          updatedByUser.FirstNameEn,
		// 		"UpdatedByLastNameEn":           updatedByUser.LastNameEn,
		// 		"AssignedOn":                    time.Now().Format("2006-01-02"),
		// 		"RedirectURL":                   fmt.Sprintf("%s/compliance/register/%d/requirements?requirements=%d", appURL, complianceID, id),
		// 	}

		// 	newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 		EmailTo:              newResponsibleUserEmailTo,
		// 		Subject:              newResponsibleUserSubject,
		// 		TemplateRelativePath: "compliance/templates/email/compliance_requirement_change_responsible_user_mail.html",
		// 		EmailData:            newResponsibleUserEmailData,
		// 	}

		// 	// Send notifications and emails
		// 	err = crs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_requirement")
		// 	if err != nil {
		// 		utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
		// 			"error": err.Error(),
		// 		})
		// 	}

		// 	err = crs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "compliance_requirement")
		// 	if err != nil {
		// 		utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
		// 			"error": err.Error(),
		// 		})
		// 	}
		// }
	}()
	return nil
}

func (crs *ComplianceRequirementService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	_, err := crs.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceRequirement []map[string]interface{}

	// var notifications []struct {
	// 	OwnerID          *uint
	// 	Title            string
	// 	RequirementTitle string
	// 	ReferenceNumber  string
	// }

	// Track clause IDs for deletion
	clauseIDsToDelete := make(map[uint]bool)

	for _, id := range ids {
		requirement, err := crs.GetFullDetailsByID(complianceID, id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Prepare notification data before deletion
		// notifications = append(notifications, struct {
		// 	OwnerID          *uint
		// 	Title            string
		// 	RequirementTitle string
		// 	ReferenceNumber  string
		// }{
		// 	OwnerID:          &requirement.ResponsibleUserID,
		// 	Title:            compliance.TitleEn,
		// 	RequirementTitle: requirement.TitleEn,
		// 	ReferenceNumber:  compliance.ReferenceNumber,
		// })
		if len(requirement.ResponsibleUnits) > 0 {
			var deleteUnitIDs []*uint
			for _, unit := range requirement.ResponsibleUnits {
				deleteUnitIDs = append(deleteUnitIDs, &unit.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementResponsibleUnits(deleteUnitIDs, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Check if controls exist
		if len(requirement.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range requirement.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementControls(deleteControlIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Check if requiredDocumentTitile exist
		if len(requirement.RequiredDocumentTitles) > 0 {
			var deleteRequiredDocumentTitleIds []*uint
			for _, rdTitle := range requirement.RequiredDocumentTitles {
				deleteRequiredDocumentTitleIds = append(deleteRequiredDocumentTitleIds, rdTitle.ID)
			}
			if err := crs.repo.DeleteComplianceRequirmentRequiredDocumentTitles(deleteRequiredDocumentTitleIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Check if strategic objectives exist
		if len(requirement.StrategicObjectives) > 0 {
			var deleteStrategicObjectiveIds []*uint
			for _, obj := range requirement.StrategicObjectives {
				deleteStrategicObjectiveIds = append(deleteStrategicObjectiveIds, obj.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementStrategicObjectives(deleteStrategicObjectiveIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if len(requirement.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, document := range requirement.Documents {
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		complianceRequirementStatusUpdates, err := crs.repo.GetStatusUpdateByRequirementID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if len(complianceRequirementStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range complianceRequirementStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := crs.repo.DeleteComplianceRequirementStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if err := crs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
			// Track the clause IDs
			clauseIDsToDelete[uint(requirement.ComplianceClause.ID)] = true
		}

		deletedComplianceRequirement = append(deletedComplianceRequirement, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": complianceID,
			"title":                  requirement.TitleEn,
		})

	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Delete clauses if all their requirements are deleted
	// for clauseID := range clauseIDsToDelete {
	// 	count, err := crs.repo.CountRequirementsByClauseID(clauseID, tx)
	// 	if err != nil {
	// 		failedIds = append(failedIds, clauseID)
	// 		fmt.Printf("Error counting requirements for clause %d: %v\n", clauseID, err)
	// 		continue
	// 	} else if count == 0 {
	// 		if err := crs.repo.DeleteClause(clauseID, tx); err != nil {
	// 			failedIds = append(failedIds, clauseID)
	// 			// Log error, but proceed with other deletions
	// 			fmt.Printf("Error deleting clause %d: %v\n", clauseID, err)
	// 			continue
	// 		}
	// 	}
	// }

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Log successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirement{}.TableName(),
				successfulIds,
				deletedComplianceRequirement,
				currentUserID,
				nil,
			)
			if err != nil {
				// Log error, but proceed as it might be recoverable later
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
			// Send DB and email notifications for each successfully deleted compliance register
			// for _, notification := range notifications {
			// 	urlParam := ""
			// 	entity := "compliance-requirement"
			// 	action := "delete"

			// 	// Prepare DB notification
			// 	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			// 		UserIDs:              []uint{*notification.OwnerID},
			// 		ModuleIdentifier:     "compliance",
			// 		NotificationType:     string(commonDTOs.DeletionType),
			// 		HeadingEn:            "Compliance Requirement Deleted",
			// 		HeadingAr:            "تم حذف متطلبات الامتثال",
			// 		TitleEn:              notification.RequirementTitle,
			// 		TitleAr:              notification.RequirementTitle,
			// 		ReferenceNumber:      compliance.ReferenceNumber,
			// 		Entity:               &entity,
			// 		Action:               &action,
			// 		URLParam:             urlParam,
			// 		IsSystemNotification: false,
			// 		CreatedByID:          &currentUserID,
			// 	}

			// 	// Fetch owner details (responsible user)
			// 	ownerUser, err := crs.userSvc.GetByID(*notification.OwnerID)
			// 	if err != nil {
			// 		utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
			// 			"error": err.Error(),
			// 		})
			// 		continue
			// 	}

			// 	// Prepare email notification
			// 	ownerEmailTo := []string{ownerUser.Email}
			// 	subject := fmt.Sprintf("Compliance Requirement/%s: Deleted", notification.RequirementTitle)
			// 	templateRelativePath := "compliance/templates/email/compliance_requirement_delete_mail.html"

			// 	emailData := map[string]interface{}{
			// 		"FirstNameEn":               ownerUser.FirstNameEn,
			// 		"LastNameEn":                ownerUser.LastNameEn,
			// 		"ComplianceReferenceNumber": notification.ReferenceNumber,
			// 		"ComplianceTitle":           notification.Title,
			// 		"RequirementTitleEn":        notification.RequirementTitle,
			// 	}

			// 	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			// 		EmailTo:              ownerEmailTo,
			// 		Subject:              subject,
			// 		TemplateRelativePath: templateRelativePath,
			// 		EmailData:            emailData,
			// 	}

			// 	// Send DB and email notifications
			// 	err = crs.notificationSvc.SendNotification(
			// 		dbNotificationDTO, emailNotificationDTO, "compliance_requirement")
			// 	if err != nil {
			// 		utility.LogError("Failed to send deletion notifications", map[string]interface{}{
			// 			"error": err.Error(),
			// 		})
			// 	}
			// }
		}()
	}

	tx.Commit()

	return failedIds, nil
}

func (crs *ComplianceRequirementService) StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingComplianceRequirement, err := crs.GetByID(complianceID, id)
	if err != nil {
		return err
	}

	existingComplianceRequirement.ComplianceStatusID = &complianceStatusID
	existingComplianceRequirement.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingComplianceRequirement, tx)
	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRequirementService) calculateComplianceRequirementScore(complianceStatusID uint, severityLevelID uint, tx ...utility.DatabaseTransaction) (int, error) {
	complianceStatus, err := crs.complianceStatusSvc.GetByID(complianceStatusID)
	if err != nil {
		return 0, err
	}

	severityLevel, err := crs.severityLevelSvc.GetByID(severityLevelID)
	if err != nil {
		return 0, err
	}

	maxSeverityWeight, err := crs.severityLevelSvc.GetMaxSeverityWeight(tx...)
	if err != nil {
		return 0, err
	}

	maxBaseScore, err := crs.complianceStatusSvc.GetMaxBaseScore()
	if err != nil {
		return 0, err
	}

	maxPossibleScore := maxSeverityWeight * maxBaseScore

	// Calculate the actual score
	rawScore := complianceStatus.BaseScore * severityLevel.Weight

	// Normalize the score to be out of 100
	normalizedScore := (float64(rawScore) / float64(maxPossibleScore)) * 100

	return int(normalizedScore), nil
}

// getInherentRiskRatingByScore determines the appropriate risk rating based on the calculated score
func (crs *ComplianceRequirementService) getInherentRiskRatingByScore(score float64) (*uint, error) {
	// Get all compliance assessment ratings
	ratings, err := crs.repo.GetComplianceAssessmentRatings()
	if err != nil {
		return nil, err
	}

	// Find the appropriate rating based on score range
	for _, rating := range ratings {
		if rating.ScoreFrom != nil && rating.ScoreTo != nil {
			if score >= *rating.ScoreFrom && score <= *rating.ScoreTo {
				return &rating.ID, nil
			}
		}
	}

	return nil, nil // No rating found
}

// CreateInherentRiskAnalysis creates inherent risk analysis for a compliance requirement
func (crs *ComplianceRequirementService) CreateInherentRiskAnalysis(complianceRequirementID uint, currentUserID uint, riskAnalysisFormDTO []*complianceDTOs.ComplianceRequirementInherentRiskFormDTO) error {
	tx := crs.repo.BeginTransaction()

	// Delete existing risk analysis for this requirement
	if err := crs.repo.DeleteInherentRiskAnalysisByRequirementID(complianceRequirementID, tx); err != nil {
		tx.Rollback()
		return err
	}

	var riskAnalysisRecords []*models.ComplianceRequirementInherentRiskAnalysis
	var totalScore float64

	// Create new risk analysis records
	for _, analysis := range riskAnalysisFormDTO {
		if analysis.Score == nil {
			return fmt.Errorf("score is required")
		}
		record := &models.ComplianceRequirementInherentRiskAnalysis{
			ComplianceRequirementID:          complianceRequirementID,
			ComplianceAssessmentImpactTypeID: analysis.ComplianceAssessmentImpactTypeID,
			Score:                            *analysis.Score,
			Comment:                          analysis.Comment,
			CreatedByID:                      currentUserID,
			CreatedAt:                        time.Now(),
		}
		riskAnalysisRecords = append(riskAnalysisRecords, record)
		totalScore += float64(*analysis.Score)
	}

	// Create inherent risk analysis records
	if err := crs.repo.CreateInherentRiskAnalysis(riskAnalysisRecords, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Calculate and update inherent risk score
	if len(riskAnalysisFormDTO) > 0 {
		averageScore := totalScore / float64(len(riskAnalysisFormDTO))
		ratingID, err := crs.getInherentRiskRatingByScore(averageScore)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := crs.repo.UpdateInherentRiskScore(complianceRequirementID, averageScore, ratingID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// DeleteInherentRiskAnalysis deletes all inherent risk analysis for a compliance requirement
func (crs *ComplianceRequirementService) DeleteInherentRiskAnalysis(complianceRequirementID uint, currentUserID uint) error {
	tx := crs.repo.BeginTransaction()

	// Delete all risk analysis for this requirement
	if err := crs.repo.DeleteInherentRiskAnalysisByRequirementID(complianceRequirementID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reset inherent risk score and rating
	if err := crs.repo.UpdateInherentRiskScore(complianceRequirementID, 0, nil, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (crs *ComplianceRequirementService) validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	if len(ComplianceRequirementFormDTO.ResponsibleUnits) == 0 {
		return utility.NewValidationError(
			fmt.Errorf("at least one responsible unit is required"),
			"required",
			"responsible_units",
		)
	}

	responsibleUnitIDs := make(map[uint]bool)
	for _, responsibleUnit := range ComplianceRequirementFormDTO.ResponsibleUnits {
		if responsibleUnit.ResponsibleUnitID == 0 {
			return utility.NewValidationError(fmt.Errorf("responsible unit id is required"), "required", "responsible_unit_id")
		}
		if responsibleUnitIDs[responsibleUnit.ResponsibleUnitID] {
			return utility.NewValidationError(fmt.Errorf("duplicate responsible unit ID: %d", responsibleUnit.ResponsibleUnitID), "duplicate", "responsible unit")
		}
		if responsibleUnit.IsRegulationNotificationAssessment == nil {
			return utility.NewValidationError(
				fmt.Errorf("is regulation notification assessment is required for responsible unit id: %d", responsibleUnit.ResponsibleUnitID),
				"required",
				"is_regulation_notification_assessment",
			)
		}
		if responsibleUnit.IsRiskBasedAssessment == nil {
			return utility.NewValidationError(
				fmt.Errorf("is risk based assessment is required for responsible unit id: %d", responsibleUnit.ResponsibleUnitID),
				"required",
				"is_risk_based_assessment",
			)
		}
		if _, err := crs.organizationUnitSvc.GetByID(responsibleUnit.ResponsibleUnitID); err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible unit does not exist for ID: %d", responsibleUnit.ResponsibleUnitID), "not_exist", "responsible unit")
		}
		responsibleUnitIDs[responsibleUnit.ResponsibleUnitID] = true
	}

	if ComplianceRequirementFormDTO.ComplianceRequirementCategoryID != nil {
		if _, err := crs.complianceRequirementCategorySvc.GetByID(*ComplianceRequirementFormDTO.ComplianceRequirementCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement category does not exist"), "not_exist", "compliance requirement category")
		}
	}
	// if _, err := crs.complianceStatusSvc.GetByID(*ComplianceRequirementFormDTO.ComplianceStatusID); err != nil {
	// 	return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	// }
	if ComplianceRequirementFormDTO.SeverityLevelID != nil {
		if _, err := crs.severityLevelSvc.GetByID(*ComplianceRequirementFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}
	if _, err := crs.complianceClauseSvc.GetByID(ComplianceRequirementFormDTO.ComplianceClauseID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance clause does not exist"), "not_exist", "compliance clause")
	}

	controlIds := ComplianceRequirementFormDTO.ControlIDs

	controlIDs := make(map[uint]bool)

	for _, controlId := range controlIds {
		if controlIDs[*controlId] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", controlId), "duplicate", "control")
		}

		if _, err := crs.controlSvc.GetByID(*controlId); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlId), "not_exist", "control")
		}

		controlIDs[*controlId] = true
	}

	strategicObjectiveIds := ComplianceRequirementFormDTO.StrategicObjectiveIDs

	strategicObjectiveIDs := make(map[uint]bool)

	for _, strategicObjectiveId := range strategicObjectiveIds {
		if strategicObjectiveIDs[*strategicObjectiveId] {
			return utility.NewValidationError(fmt.Errorf("duplicate strategic objective ID: %d", strategicObjectiveId), "duplicate", "strategic objective")
		}

		if _, err := crs.strategicObjectiveSvc.GetByID(*strategicObjectiveId); err != nil {
			return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveId), "not_exist", "strategic objective")
		}

		strategicObjectiveIDs[*strategicObjectiveId] = true
	}

	requiredDocumentTitleIds := ComplianceRequirementFormDTO.RequiredDocumentTitleIDs
	requiredDocumentTitleIDs := make(map[uint]bool)

	for _, requiredDocumentTitleId := range requiredDocumentTitleIds {
		if requiredDocumentTitleIDs[*requiredDocumentTitleId] {
			return utility.NewValidationError(fmt.Errorf("duplicate required document ID: %d", requiredDocumentTitleId), "duplicate", "required document")
		}

		if _, err := crs.requiredDocumentTitleSvc.GetByID(*requiredDocumentTitleId); err != nil {
			return utility.NewValidationError(fmt.Errorf("required document does not exist for ID: %d", requiredDocumentTitleId), "not_exist", "required document")
		}

		requiredDocumentTitleIDs[*requiredDocumentTitleId] = true
	}

	return nil
}

func (crs *ComplianceRequirementService) parseComplianceRequirementDates(dto *complianceDTOs.ComplianceRequirementFormDTO) (*complianceDTOs.ComplianceRequirementParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRequirementParsedDates{}

	var err error
	if dto.Deadline != nil {
		fmt.Println("hiiii22=", parsed.Deadline)
		parsed.Deadline, err = utility.ParseDate(*dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	return parsed, nil
}

func (srs *ComplianceRequirementService) ExportComplianceRequirements(data []*complianceDTOs.ComplianceRequirementListDTO, currentUserID uint) (*os.File, error) {
	// Retrieve user settings for language preference
	userSetting, err := srs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Clause No",
			"Clause",
			"Title (AR)",
			"Category (AR)",
			"Severity Level (AR)",
			"Deadline",
			"Compliance Status (AR)",
		}
	} else {
		headers = []string{
			"SL No.",
			"Clause No",
			"Clause",
			"Title (EN)",
			"Category (EN)",
			"Severity Level (EN)",
			"Deadline",
			"Compliance Status (EN)",
		}
	}

	// Convert data to [][]interface{} based on language preference
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				item.ComplianceClauseNo,
				item.ComplianceClauseTitleAr,
				item.TitleAr,
				item.ComplianceRequirementCategoryTitleAr,
				item.SeverityLevelTitleAr,
				item.Deadline,
				item.ComplianceStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				item.ComplianceClauseNo,
				item.ComplianceClauseTitleEn,
				item.TitleEn,
				item.ComplianceRequirementCategoryTitleEn,
				item.SeverityLevelTitleEn,
				item.Deadline,
				item.ComplianceStatusTitleEn,
			}
		}
	}
	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

// Generate Excel template
func (cs *ComplianceRequirementService) GenerateRequirementTemplate(complianceID uint) (*os.File, error) {
	templateHeaders := []string{
		"Clause", "TitleEn", "TitleAr", "Category", "SeverityLevel", "Deadline", "RequiredDocuments", "ResponsibleUnits", "IsRegulationNotificationAssessment", "IsRiskBasedAssessment",
	}
	templateData := [][]interface{}{}

	// Clauses
	clauses, _ := cs.complianceClauseSvc.GetAllForImport(complianceID)
	clauseHeaders := []string{"ClauseNo", "TitleEn"}
	clauseData := [][]interface{}{}
	for _, c := range clauses {
		clauseData = append(clauseData, []interface{}{c.ClauseNo, c.TitleEn})
	}

	// Categories
	categories, _ := cs.complianceRequirementCategorySvc.GetAll("", "", "", "active")
	categoryHeaders := []string{"CategoryEn", "CategoryAr"}
	categoryData := [][]interface{}{}
	for _, cat := range categories {
		categoryData = append(categoryData, []interface{}{cat.TitleEn, cat.TitleAr})
	}

	// Severity Levels
	severities, _ := cs.severityLevelSvc.GetAll("", "", "", "active")
	severityHeaders := []string{"SeverityEn", "SeverityAr"}
	severityData := [][]interface{}{}
	for _, sev := range severities {
		severityData = append(severityData, []interface{}{sev.TitleEn, sev.TitleAr})
	}

	// Responsible Units
	params := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort:          "",
		SortBy:        "",
		Keyword:       "",
		ParentUnitIDs: []uint{},
		Primary:       true,
	}

	responsibleUnits, _ := cs.organizationUnitSvc.GetAll(params)
	responsibleUnitHeaders := []string{"ResponsibleUnitEn", "ResponsibleUnitAr"}
	responsibleUnitData := [][]interface{}{}
	for _, ru := range responsibleUnits {
		responsibleUnitData = append(responsibleUnitData, []interface{}{ru.TitleEn, ru.TitleAr})
	}

	// Required Documents
	requiredDocumentTitles, _ := cs.requiredDocumentTitleSvc.GetAll("", "", "", "active")
	requiredDocumentTitleHeaders := []string{"RequiredDocumentTitleEn", "RequiredDocumentTitleAr"}
	requiredDocumentTitleData := [][]interface{}{}
	for _, rdt := range requiredDocumentTitles {
		requiredDocumentTitleData = append(requiredDocumentTitleData, []interface{}{rdt.TitleEn, rdt.TitleAr})
	}

	dropDownConfig := map[string]bool{
		"Clause":                             true,
		"Category":                           true,
		"SeverityLevel":                      true,
		"TitleEn":                            false,
		"TitleAr":                            false,
		"Deadline":                           false,
		"ResponsibleUnits":                   false,
		"RequiredDocuments":                  false,
		"IsRegulationNotificationAssessment": false,
		"IsRiskBasedAssessment":              false,
	}
	// Generate file
	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":         append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"Clause":           append([][]interface{}{utility.ConvertToInterfaceSlice(clauseHeaders)}, clauseData...),
		"Category":         append([][]interface{}{utility.ConvertToInterfaceSlice(categoryHeaders)}, categoryData...),
		"SeverityLevel":    append([][]interface{}{utility.ConvertToInterfaceSlice(severityHeaders)}, severityData...),
		"ResponsibleUnit":  append([][]interface{}{utility.ConvertToInterfaceSlice(responsibleUnitHeaders)}, responsibleUnitData...),
		"RequiredDocument": append([][]interface{}{utility.ConvertToInterfaceSlice(requiredDocumentTitleHeaders)}, requiredDocumentTitleData...),
	}, dropDownConfig)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Validate Excel Template
func (cs *ComplianceRequirementService) ValidateRequirementTemplate(complianceID uint, file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	inputHeaders := []string{"Clause", "TitleEn", "TitleAr", "Category", "SeverityLevel", "Deadline", "RequiredDocuments", "ResponsibleUnits", "IsRegulationNotificationAssessment", "IsRiskBasedAssessment"}
	snakeHeaders := []string{"clause", "title_en", "title_ar", "category", "severity_level", "deadline", "required_documents", "responsible_units", "is_regulation_notification_assessment", "is_risk_based_assessment"}

	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Requirement Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, h := range snakeHeaders {
			if i < len(row) {
				record[h] = row[i]
			} else {
				record[h] = ""
			}
		}
		input = append(input, record)
	}

	return cs.validateRequirements(complianceID, input)
}

// Validate JSON input
func (cs *ComplianceRequirementService) ValidateImportRequirement(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	return cs.validateRequirements(complianceID, input)
}

// Import
func (cs *ComplianceRequirementService) ImportRequirements(complianceID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validated, err := cs.validateRequirements(complianceID, input)
	if err != nil {
		return nil, err
	}

	for _, r := range validated {
		if errs, ok := r["errors"].(map[string]interface{}); ok && len(errs) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validated {

		var deadlineStr *string
		if val, ok := record["deadline"].(string); ok && val != "" {
			t, err := time.Parse("2006-01-02", strings.TrimSpace(val))
			if err == nil {
				formatted := t.Format("2006-01-02")
				deadlineStr = &formatted
			}
		}

		// Convert []uint to responsible units
		var responsibleUnits []*complianceDTOs.ComplianceRequirementResponsibleUnitFormDTO
		if ids, ok := record["responsible_unit_ids"].([]uint); ok {
			isRegulationNotificationAssessment := record["is_regulation_notification_assessment"].(bool)
			isRiskBasedAssessment := record["is_risk_based_assessment"].(bool)

			for _, id := range ids {
				regulationNotificationAssessment := isRegulationNotificationAssessment
				riskBasedAssessment := isRiskBasedAssessment
				responsibleUnits = append(responsibleUnits, &complianceDTOs.ComplianceRequirementResponsibleUnitFormDTO{
					ResponsibleUnitID:                  id,
					IsRegulationNotificationAssessment: &regulationNotificationAssessment,
					IsRiskBasedAssessment:              &riskBasedAssessment,
				})
			}
		}

		var requiredDocumentIDs []*uint
		if ids, ok := record["required_document_ids"].([]uint); ok {
			for _, id := range ids {
				idCopy := id // prevent pointer reuse
				requiredDocumentIDs = append(requiredDocumentIDs, &idCopy)
			}
		}
		dto := &complianceDTOs.ComplianceRequirementFormDTO{
			ComplianceClauseID:              record["clause_id"].(uint),
			TitleEn:                         record["title_en"].(string),
			TitleAr:                         record["title_ar"].(string),
			ComplianceRequirementCategoryID: utility.GetOptionalUintValue(record["category_id"]),
			SeverityLevelID:                 utility.GetOptionalUintValue(record["severity_level_id"]),
			Deadline:                        deadlineStr,
			RequiredDocumentTitleIDs:        requiredDocumentIDs,
			ResponsibleUnits:                responsibleUnits,
		}

		id, err := cs.Create(complianceID, currentUserID, dto)
		if err != nil {
			return nil, err
		}
		createdIDs = append(createdIDs, id)
	}
	return createdIDs, nil
}

func (cs *ComplianceRequirementService) validateRequirements(complianceID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"clause", "title_en", "title_ar", "category", "severity_level", "deadline", "required_documents", "responsible_units", "is_regulation_notification_assessment", "is_risk_based_assessment"}

	// --- Fetch lookup data ---

	// Clauses
	clauses, err := cs.complianceClauseSvc.GetAllForImport(complianceID)
	if err != nil {
		return nil, fmt.Errorf("error fetching clauses: %v", err)
	}
	clauseMap := make(map[string]uint)
	for _, cl := range clauses {
		key := strings.TrimSpace(cl.TitleEn)
		clauseMap[key] = cl.ID
	}

	// Categories
	categories, err := cs.complianceRequirementCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching categories: %v", err)
	}
	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[strings.TrimSpace(cat.TitleEn)] = cat.ID
	}

	// Severity Levels
	severities, err := cs.severityLevelSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching severity levels: %v", err)
	}
	severityMap := make(map[string]uint)
	for _, sev := range severities {
		severityMap[strings.TrimSpace(sev.TitleEn)] = sev.ID
	}

	// Required Documents
	requiredDocumentTitles, err := cs.requiredDocumentTitleSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching required documents titles: %v", err)
	}
	requiredDocumentsTitleMap := make(map[string]uint)
	for _, rdt := range requiredDocumentTitles {
		requiredDocumentsTitleMap[strings.TrimSpace(rdt.TitleEn)] = rdt.ID
	}

	// Responsible Units
	params := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort:          "",
		SortBy:        "",
		Keyword:       "",
		ParentUnitIDs: []uint{},
		Primary:       true,
	}

	responsibleUnits, err := cs.organizationUnitSvc.GetAll(params)
	if err != nil {
		return nil, fmt.Errorf("error fetching responsible units: %v", err)
	}
	responsibleUnitMap := make(map[string]uint)
	for _, ru := range responsibleUnits {
		responsibleUnitMap[strings.TrimSpace(ru.TitleEn)] = ru.ID
	}

	// --- Validate each record ---
	var validated []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processed := make(map[string]interface{})

		// Copy original input
		for k, v := range record {
			processed[k] = v
		}

		// Validate Clause
		if val, ok := record["clause"].(string); !ok || val == "" {
			rowErrors["clause"] = "Field is required"
		} else if id, exists := clauseMap[strings.TrimSpace(val)]; exists {
			processed["clause_id"] = id
		} else {
			rowErrors["clause"] = fmt.Sprintf("Invalid clause: %s", val)
		}

		// TitleEn
		if val, ok := record["title_en"].(string); !ok || val == "" {
			rowErrors["title_en"] = "Field is required"
		}

		// TitleAr
		if val, ok := record["title_ar"].(string); !ok || val == "" {
			rowErrors["title_ar"] = "Field is required"
		}

		// Category
		if val, ok := record["category"].(string); !ok || val == "" {
			rowErrors["category"] = "Field is required"
		} else if id, exists := categoryMap[strings.TrimSpace(val)]; exists {
			processed["category_id"] = id
		} else {
			rowErrors["category"] = fmt.Sprintf("Invalid category: %s", val)
		}

		// SeverityLevel
		if val, ok := record["severity_level"].(string); !ok || val == "" {
			rowErrors["severity_level"] = "Field is required"
		} else if id, exists := severityMap[strings.TrimSpace(val)]; exists {
			processed["severity_level_id"] = id
		} else {
			rowErrors["severity_level"] = fmt.Sprintf("Invalid severity level: %s", val)
		}

		// Resposible Units
		if val, ok := record["responsible_units"].(string); !ok || val == "" {
			rowErrors["responsible_units"] = "Field is required"
		} else {
			unitNames := strings.Split(val, ",")
			var unitIDs []uint
			var invalidUnits []string

			for _, name := range unitNames {
				trimmed := strings.TrimSpace(name)
				if id, exists := responsibleUnitMap[trimmed]; exists {
					unitIDs = append(unitIDs, id)
				} else {
					invalidUnits = append(invalidUnits, trimmed)
				}
			}

			if len(unitIDs) > 0 {
				processed["responsible_unit_ids"] = unitIDs
			}

			if len(invalidUnits) > 0 {
				rowErrors["responsible_units"] = fmt.Sprintf("Invalid responsible units: %s", strings.Join(invalidUnits, ", "))
			}
		}

		// IsRegulationNotificationAssessment
		if val, exists := record["is_regulation_notification_assessment"]; !exists {
			rowErrors["is_regulation_notification_assessment"] = "Field is required"
		} else if parsedValue, ok := parseBoolValue(val); ok {
			processed["is_regulation_notification_assessment"] = parsedValue
		} else {
			rowErrors["is_regulation_notification_assessment"] = "Invalid boolean value (use true/false)"
		}

		// IsRiskBasedAssessment
		if val, exists := record["is_risk_based_assessment"]; !exists {
			rowErrors["is_risk_based_assessment"] = "Field is required"
		} else if parsedValue, ok := parseBoolValue(val); ok {
			processed["is_risk_based_assessment"] = parsedValue
		} else {
			rowErrors["is_risk_based_assessment"] = "Invalid boolean value (use true/false)"
		}

		// Required Documents
		if val, ok := record["required_documents"].(string); ok && val != "" {
			docNames := strings.Split(val, ",")
			var docIDs []uint
			var invalidDocs []string

			for _, name := range docNames {
				trimmed := strings.TrimSpace(name)
				if id, exists := requiredDocumentsTitleMap[trimmed]; exists {
					docIDs = append(docIDs, id)
				} else {
					invalidDocs = append(invalidDocs, trimmed)
				}
			}

			if len(docIDs) > 0 {
				processed["required_document_ids"] = docIDs
			}

			if len(invalidDocs) > 0 {
				rowErrors["required_documents"] = fmt.Sprintf("Invalid document titles: %s", strings.Join(invalidDocs, ", "))
			}
		}

		// Deadline
		if val, ok := record["deadline"].(string); ok && val != "" {
			t, err := time.Parse("2006-01-02", strings.TrimSpace(val))
			if err != nil {
				rowErrors["deadline"] = fmt.Sprintf("Invalid date format (use YYYY-MM-DD): %s", val)
			} else {
				formatted := t.Format("2006-01-02")
				processed["deadline"] = formatted
			}
		}

		// Collect errors
		if len(rowErrors) == 0 {
			processed["errors"] = nil
		} else {
			for _, h := range headers {
				if _, exists := rowErrors[h]; !exists {
					rowErrors[h] = "" // keep consistent structure
				}
			}
			processed["errors"] = rowErrors
		}

		validated = append(validated, processed)
	}

	return validated, nil
}

func parseBoolValue(value interface{}) (bool, bool) {
	switch v := value.(type) {
	case bool:
		return v, true
	case string:
		trimmed := strings.TrimSpace(strings.ToLower(v))
		switch trimmed {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		default:
			return false, false
		}
	case float64:
		if v == 1 {
			return true, true
		}
		if v == 0 {
			return false, true
		}
	}

	return false, false
}
