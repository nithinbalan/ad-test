package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ControlRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error)
	GetAllControls() ([]*models.Control, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error)
	Create(control *models.Control, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceRegister *models.Control, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)

	//lifecycle status
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)
	CheckControlAutomationExists(id uint) (bool, error)
	CheckControlConsistencyExists(id uint) (bool, error)
}

type ControlRepository struct {
	db *gorm.DB
}

func NewControlRepository(db *gorm.DB) *ControlRepository {
	return &ControlRepository{
		db: db,
	}
}

func (repo *ControlRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlRepository) GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error) {
	var controls []*organizationDTOs.ControlListDTO
	var totalCount int64

	query := repo.db.Table("controls").
		Select(`
			controls.*,
			control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			control_automations.title_en AS control_automation_title_en,
			control_automations.title_ar AS control_automation_title_ar,
			control_automations.type AS control_automation_type,
			control_automations.score AS control_automation_score,
			control_automations.status AS control_automation_status,
			control_consistencies.title_en AS control_consistency_title_en,
			control_consistencies.title_ar AS control_consistency_title_ar,
			control_consistencies.type AS control_consistency_type,
			control_consistencies.score AS control_consistency_score,
			control_consistencies.status AS control_consistency_status,
			control_consistencies.description_en AS control_consistency_description_en,
			control_consistencies.description_ar AS control_consistency_description_ar,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			test_responsible_user.id AS test_responsible_user_id,
			test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
			test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
			test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
			test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
			test_responsible_user.initial AS test_responsible_user_initial,
			test_responsible_user.color AS test_responsible_user_color,
			test_responsible_user.email AS test_responsible_user_email,
			test_responsible_user.profile_image_path AS test_responsible_user_profile_image_path,
			test_frequencies.title_en AS test_frequency_title_en,
			test_frequencies.title_ar AS test_frequency_title_ar,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_rating_values.value AS risk_rating_value,
			control_effectiveness_ratings.id AS control_effectiveness_rating_id,
			control_effectiveness_ratings.title_en AS control_effectiveness_rating_title_en,
			control_effectiveness_ratings.title_ar AS control_effectiveness_rating_title_ar,
			control_effectiveness_ratings.min_score AS control_effectiveness_rating_min_score,
			control_effectiveness_ratings.max_score AS control_effectiveness_rating_max_score,
			control_effectiveness_ratings.type AS control_effectiveness_rating_type
		`).
		Joins("LEFT JOIN control_types ON control_types.id = controls.control_type_id").
		Joins("LEFT JOIN control_categories ON control_categories.id = controls.control_category_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = controls.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = controls.responsible_user_id").
		Joins("LEFT JOIN control_automations ON control_automations.id = controls.control_automation_id").
		Joins("LEFT JOIN control_consistencies ON control_consistencies.id = controls.control_consistency_id").
		Joins("LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = controls.review_frequency_id").
		Joins("LEFT JOIN users AS test_responsible_user ON test_responsible_user.id = controls.test_responsible_user_id").
		Joins("LEFT JOIN review_frequencies AS test_frequencies ON test_frequencies.id = controls.test_frequency_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = controls.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = controls.updated_by_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = controls.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = controls.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN control_effectiveness_ratings ON control_effectiveness_ratings.id = controls.control_effectiveness_rating_id")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("controls."+columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if len(params.ComplianceRegisterIDs) > 0 {
		query = query.Where("compliance_register.id IN (?)", params.ComplianceRegisterIDs)
	}

	// Apply last_review_date filters
	if params.LastReviewDateFrom != "" {
		query = query.Where("last_review_date >= ?", params.LastReviewDateFrom)
	}
	if params.LastReviewDateTo != "" {
		query = query.Where("last_review_date <= ?", params.LastReviewDateTo)
	}

	// Apply next_review_date filters
	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

	// Apply last_test_date filters
	if params.LastTestDateFrom != "" {
		query = query.Where("last_test_date >= ?", params.LastTestDateFrom)
	}
	if params.LastTestDateTo != "" {
		query = query.Where("last_test_date <= ?", params.LastTestDateTo)
	}

	// Apply next_test_date filters
	if params.NextTestDateFrom != "" {
		query = query.Where("next_test_date >= ?", params.NextTestDateFrom)
	}
	if params.NextTestDateTo != "" {
		query = query.Where("next_test_date <= ?", params.NextTestDateTo)
	}

	// Apply implemented_at filters
	if params.ImplementedAtFrom != "" {
		query = query.Where("implemented_at >= ?", params.ImplementedAtFrom)
	}
	if params.ImplementedAtTo != "" {
		query = query.Where("implemented_at <= ?", params.ImplementedAtTo)
	}

	// Apply standard code filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "applicable":
			query = query.Where("controls.is_applicable = ?", true)
		case "not_applicable":
			query = query.Where("controls.is_applicable = ?", false)
		case "implemented":
			query = query.Where("control_implementation_statuses.type = ?", "implemented")
		case "not_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "not-implemented")
		case "partially_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "partially-implemented")
		case "efficient":
			query = query.Where("control_efficiency_statuses.type = ?", "efficient")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")
		case "SAMA", "NCA", "ISMS":
			// Add the compliance_controls join only when filtering by these codes
			query = query.Joins(`
            INNER JOIN compliance_controls ON compliance_controls.control_id = controls.id
            INNER JOIN compliance_register ON compliance_register.id = compliance_controls.compliance_register_id
        `).Where("compliance_register.standard_code = ?", filter)
		case "high_risk":
			query = query.Where("risk_matrix_risk_ratings.type IN (?)", []string{"high", "very_high"})
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			controls.description LIKE ? OR
			controls.objective LIKE ? OR
			controls.code LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			test_responsible_user.first_name_en LIKE ? OR
			test_responsible_user.first_name_ar LIKE ? OR
			test_responsible_user.last_name_en LIKE ? OR
			test_responsible_user.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&controls).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &organizationDTOs.ControlListResponseDTO{
		Items: controls,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ControlRepository) GetAllControls() ([]*models.Control, error) {
	var controls []*models.Control
	result := repo.db.Select("id, title_en").Find(&controls)
	return controls, result.Error
}

func (repo *ControlRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var control *models.Control

	result := db.First(&control, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return control, nil

}

func (repo *ControlRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var control *models.Control

	latestControlTestSubQuery := db.
		Table("control_tests").
		Select("id").
		Where("control_id = ?", id).
		Order("tested_at desc").
		Limit(1)

	latestControlReviewSubQuery := db.
		Table("control_reviews").
		Select("id").
		Where("control_id = ?", id).
		Order("reviewed_at desc").
		Limit(1)

	result := db.
		Preload("ControlType").
		Preload("ControlCategory").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ControlImplementationStatus").
		Preload("ControlEfficiencyStatus").
		Preload("TestResponsibleUser").
		Preload("TestResponsibleUser.OrganizationUnit1").
		Preload("ReviewFrequency").
		Preload("ControlReviewStatus").
		Preload("TestFrequency").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Risks").
		Preload("RiskMatrixRating").
		Preload("RiskMatrixRatingValue").
		Preload("ControlAutomation").
		Preload("ControlAutomation.CreatedByUser").
		Preload("ControlConsistency").
		Preload("ControlConsistency.CreatedByUser").
		Preload("AssociatedDocuments").
		Preload("AssociatedDocuments.Document").
		Preload("ControlEffectivenessRating").
		Preload("LatestControlTest", func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)", latestControlTestSubQuery).
				Preload("TestedByUser")
		}).
		Preload("LatestControlReview", func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)", latestControlReviewSubQuery).
				Preload("ReviewedByUser").
				Preload("UpdatedByUser").
				Preload("Documents")
		}).
		First(&control, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return control, nil
}

func (repo *ControlRepository) Create(control *models.Control, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(control).Error; err != nil {
		return 0, err
	}

	return control.ID, nil
}

func (repo *ControlRepository) Update(control *models.Control, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(control).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlRepository) Delete(id uint) error {
	if err := repo.db.Delete(&models.Control{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.Control{}).
		Joins("JOIN lifecycle_statuses ON controls.lifecycle_status_id = lifecycle_statuses.id").
		Where("controls.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("controls.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *ControlRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {
	if err := repo.db.Model(&models.Control{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *ControlRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var lastReferenceNumber string
	err := db.Model(&models.Control{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "CTRL", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ControlRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	queries := []string{
		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans
		ON compliance_action_plans.id = tasks.compliance_action_plan_id
		JOIN compliance_assessment_controls
		ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls
		ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_action_plan_status_update_documents
		 FROM compliance_action_plan_status_update_documents
		 JOIN compliance_action_plan_status_updates ON compliance_action_plan_status_updates.id = compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		 JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_action_plan_status_updates
		 FROM compliance_action_plan_status_updates
		 JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_action_plan_comment_documents
		 FROM compliance_action_plan_comment_documents
		 JOIN compliance_action_plan_comments ON compliance_action_plan_comments.id = compliance_action_plan_comment_documents.compliance_action_plan_comment_id
		 JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_action_plan_comments
		 FROM compliance_action_plan_comments
		 JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_action_plans
		 FROM compliance_action_plans
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_assessment_control_comments
		 FROM compliance_assessment_control_comments
		 JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_assessment_control_comments.compliance_assessment_control_id
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE compliance_assessment_controls
		 FROM compliance_assessment_controls
		 JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		 WHERE compliance_controls.control_id = ?`,

		`DELETE audit_finding_action_plan_status_update_documents
		 FROM audit_finding_action_plan_status_update_documents
		 JOIN audit_finding_action_plan_status_updates ON audit_finding_action_plan_status_updates.id = audit_finding_action_plan_status_update_documents.audit_finding_action_plan_status_update_id
		 JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_action_plan_status_updates
		 FROM audit_finding_action_plan_status_updates
		 JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_action_plan_comment_documents
		 FROM audit_finding_action_plan_comment_documents
		 JOIN audit_finding_action_plan_comments ON audit_finding_action_plan_comments.id = audit_finding_action_plan_comment_documents.audit_finding_action_plan_comment_id
		 JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_action_plan_comments
		 FROM audit_finding_action_plan_comments
		 JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE tasks
		 FROM tasks
		 JOIN audit_finding_action_plans ON audit_finding_action_plans.id = tasks.audit_finding_action_plan_id
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_action_plans
		 FROM audit_finding_action_plans
		 JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_draft_finding_conversation_documents
		 FROM audit_draft_finding_conversation_documents
		 JOIN audit_draft_finding_conversations ON audit_draft_finding_conversations.id = audit_draft_finding_conversation_documents.audit_draft_finding_conversation_id
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_draft_finding_conversations
		 FROM audit_draft_finding_conversations
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_draft_finding_documents
		 FROM audit_draft_finding_documents
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_documents.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_draft_finding_action_plans
		 FROM audit_draft_finding_action_plans
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_draft_findings
		 FROM audit_draft_findings
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_strategic_objectives
		 FROM audit_finding_strategic_objectives
		 JOIN audit_findings ON audit_findings.id = audit_finding_strategic_objectives.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_root_cause_analyses
		 FROM audit_finding_root_cause_analyses
		 JOIN audit_findings ON audit_findings.id = audit_finding_root_cause_analyses.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_risks
		 FROM audit_finding_risks
		 JOIN audit_findings ON audit_findings.id = audit_finding_risks.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_processes
		 FROM audit_finding_processes
		 JOIN audit_findings ON audit_findings.id = audit_finding_processes.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_kpis
		 FROM audit_finding_kpis
		 JOIN audit_findings ON audit_findings.id = audit_finding_kpis.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_incidents
		 FROM audit_finding_incidents
		 JOIN audit_findings ON audit_findings.id = audit_finding_incidents.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_documents
		 FROM audit_finding_documents
		 JOIN audit_findings ON audit_findings.id = audit_finding_documents.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE FROM audit_finding_controls WHERE audit_finding_controls.control_id = ?`,

		`DELETE audit_finding_controls
		 FROM audit_finding_controls
		 JOIN audit_findings ON audit_findings.id = audit_finding_controls.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_compliance_requirements
		 FROM audit_finding_compliance_requirements
		 JOIN audit_findings ON audit_findings.id = audit_finding_compliance_requirements.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_associated_documents
		 FROM audit_finding_associated_documents
		 JOIN audit_findings ON audit_findings.id = audit_finding_associated_documents.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_finding_assets
		 FROM audit_finding_assets
		 JOIN audit_findings ON audit_findings.id = audit_finding_assets.audit_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_findings
		 FROM audit_findings
		 JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_checklist_documents
		 FROM audit_checklist_documents
		 JOIN audit_checklists ON audit_checklists.id = audit_checklist_documents.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_checklists
		 FROM audit_checklists
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.control_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		 FROM audit_plan_auditable_item_risk_controls
		 JOIN audit_plan_auditable_item_risks ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		 FROM audit_plan_auditable_item_risks
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_plan_auditable_items
		 FROM audit_plan_auditable_items
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_schedule_auditable_items
		 FROM audit_schedule_auditable_items
		 JOIN auditable_items ON auditable_items.id = audit_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_plan_schedule_auditable_items
		 FROM audit_plan_schedule_auditable_items
		 JOIN auditable_items ON auditable_items.id = audit_plan_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_plan_checklists
		 FROM audit_plan_checklists
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE audit_document_request_auditable_items
		 FROM audit_document_request_auditable_items
		 JOIN auditable_items ON auditable_items.id = audit_document_request_auditable_items.auditable_item_id
		WHERE auditable_items.control_id = ?`,

		`DELETE control_action_plan_status_update_documents
		 FROM control_action_plan_status_update_documents
		 JOIN control_action_plan_status_updates ON control_action_plan_status_updates.id = control_action_plan_status_update_documents.control_action_plan_status_update_id
		 JOIN control_action_plans ON control_action_plans.id = control_action_plan_status_updates.control_action_plan_id
		 WHERE control_action_plans.control_id = ?`,

		`DELETE control_action_plan_status_updates
		 FROM control_action_plan_status_updates
		 JOIN control_action_plans ON control_action_plans.id = control_action_plan_status_updates.control_action_plan_id
		 WHERE control_action_plans.control_id = ?`,

		`DELETE control_action_plan_comment_documents
		 FROM control_action_plan_comment_documents
		 JOIN control_action_plan_comments ON control_action_plan_comments.id = control_action_plan_comment_documents.control_action_plan_comment_id
		 JOIN control_action_plans ON control_action_plans.id = control_action_plan_comments.control_action_plan_id
		 WHERE control_action_plans.control_id = ?`,

		`DELETE control_action_plan_comments
		 FROM control_action_plan_comments
		 JOIN control_action_plans ON control_action_plans.id = control_action_plan_comments.control_action_plan_id
		 WHERE control_action_plans.control_id = ?`,

		`DELETE asset_control_documents
		 FROM asset_control_documents
		 JOIN asset_controls ON asset_controls.id = asset_control_documents.asset_control_id
		 WHERE asset_controls.control_id = ?`,

		`DELETE control_review_documents
		 FROM control_review_documents
		 JOIN control_reviews ON control_reviews.id = control_review_documents.control_review_id
		 WHERE control_reviews.control_id = ?`,

		`DELETE control_test_documents
		 FROM control_test_documents
		 JOIN control_tests ON control_tests.id = control_test_documents.control_test_id
		 WHERE control_tests.control_id = ?`,

		`DELETE risk_control_analyses
		 FROM risk_control_analyses
		 JOIN risk_controls ON risk_controls.id = risk_control_analyses.risk_control_id
		 WHERE risk_controls.control_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN control_action_plans ON control_action_plans.id = tasks.control_action_plan_id
		WHERE control_action_plans.control_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		FROM audit_plan_auditable_item_risk_controls
		WHERE audit_plan_auditable_item_risk_controls.control_id = ?`,

		`DELETE FROM risk_controls WHERE control_id = ?`,
		`DELETE FROM control_tests WHERE control_id = ?`,
		`DELETE FROM control_reviews WHERE control_id = ?`,
		`DELETE FROM asset_controls WHERE control_id = ?`,
		`DELETE FROM auditable_items WHERE control_id = ?`,
		`DELETE FROM control_action_plans WHERE control_id = ?`,
		`DELETE FROM compliance_controls WHERE control_id = ?`,
		`DELETE FROM bcp_controls WHERE control_id = ?`,
		`DELETE FROM control_associated_documents WHERE control_id = ?`,
		`DELETE FROM control_clauses WHERE control_id = ?`,
		`DELETE FROM document_controls WHERE control_id = ?`,
		`DELETE FROM incident_controls WHERE control_id = ?`,
		`DELETE FROM kpi_controls WHERE control_id = ?`,
		`DELETE FROM kri_controls WHERE control_id = ?`,
		`DELETE FROM meeting_controls WHERE control_id = ?`,
		`DELETE FROM process_controls WHERE control_id = ?`,
		`DELETE FROM project_controls WHERE control_id = ?`,
		`DELETE FROM risk_action_plan_controls WHERE control_id = ?`,
		`DELETE FROM strategic_objective_controls WHERE control_id = ?`,
		`DELETE FROM third_party_controls WHERE control_id = ?`,
	}

	for _, q := range queries {
		res := db.Exec(q, id)
		if res.Error != nil {
			return res.Error
		}
	}

	// Force hard delete (ignore soft delete)
	return db.Delete(&models.Control{}, id).Error
}

func (repo *ControlRepository) CheckControlAutomationExists(id uint) (bool, error) {
	var count int64
	if err := repo.db.Model(&models.ControlAutomation{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *ControlRepository) CheckControlConsistencyExists(id uint) (bool, error) {
	var count int64
	if err := repo.db.Model(&models.ControlConsistency{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
