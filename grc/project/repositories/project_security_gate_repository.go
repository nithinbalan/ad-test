package repositories

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *projectDTOs.ProjectSecurityGateListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateListResponseDTO, error)
	GetSummary() (*projectDTOs.ProjectSecurityGateSummaryDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGate, error)
	GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error)
	Create(projectSecurityGate *models.ProjectSecurityGate, tx ...utility.DatabaseTransaction) error
	Update(projectSecurityGate *models.ProjectSecurityGate, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ReplaceClassificationMappings(securityGateID uint, classificationIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteClassificationMappingsBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateRepository(db *gorm.DB) *ProjectSecurityGateRepository {
	return &ProjectSecurityGateRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateRepository) GetAll(params *projectDTOs.ProjectSecurityGateListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateListResponseDTO, error) {
	var projectSecurityGates []*projectDTOs.ProjectSecurityGateListDTO
	var totalCount int64

	query := repo.db.Table("project_security_gates").
		Select(`
			project_security_gates.*,
			stages.title_en AS stage_title_en,
			stages.title_ar AS stage_title_ar,
			stages.type AS stage_type,
			stages.color_code AS stage_color_code,
			readiness_status.title_en AS readiness_status_title_en,
			readiness_status.title_ar AS readiness_status_title_ar,
			readiness_status.type AS readiness_status_type,
			readiness_status.color_code AS readiness_status_color_code,
			business_owner.id AS business_owner_id,
			business_owner.first_name_en AS business_owner_first_name_en,
			business_owner.first_name_ar AS business_owner_first_name_ar,
			business_owner.last_name_en AS business_owner_last_name_en,
			business_owner.last_name_ar AS business_owner_last_name_ar,
			business_owner.initial AS business_owner_initial,
			business_owner.color AS business_owner_color,
			business_owner.email AS business_owner_email,
			business_owner.profile_image_path AS business_owner_profile_image_path,
			division.id AS division_id,
			division.title_en AS division_title_en,
			division.title_ar AS division_title_ar,
			department.id AS department_id,
			department.title_en AS department_title_en,
			department.title_ar AS department_title_ar,
			created_user.id AS created_by_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_by_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			(SELECT COUNT(*) FROM project_security_gate_checklists
				WHERE project_security_gate_checklists.security_gate_id = project_security_gates.id) AS checklist_count
		`).
		Joins(`
			LEFT JOIN project_security_gate_stages AS stages ON stages.id = project_security_gates.stage_id
			LEFT JOIN project_security_gate_readiness_statuses AS readiness_status ON readiness_status.id = project_security_gates.readiness_status_id
			LEFT JOIN users AS business_owner ON business_owner.id = project_security_gates.business_owner_id
			LEFT JOIN organization_units AS division ON division.id = project_security_gates.division_id
			LEFT JOIN organization_units AS department ON department.id = project_security_gates.department_id
			LEFT JOIN users AS created_user ON created_user.id = project_security_gates.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = project_security_gates.updated_by_id
		`)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_security_gates.project_title_en LIKE ? OR
			project_security_gates.project_title_ar LIKE ? OR
			business_owner.first_name_en LIKE ? OR
			business_owner.last_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "classification_id" {
				query = query.Where("project_security_gates.id IN (SELECT security_gate_id FROM project_security_gate_classification_mappings WHERE classification_id IN (?))", values)
			} else {
				query = query.Where("project_security_gates."+columnName+" IN (?)", values)
			}
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "exception_pending":
			query = query.Where(`
				EXISTS (
					SELECT 1 FROM project_security_gate_checklists AS exception_checklist
					WHERE exception_checklist.security_gate_id = project_security_gates.id
						AND exception_checklist.is_exception = 1
						AND exception_checklist.approved_by_id IS NULL
				)`)
		case "outstanding":
			query = query.Where(`
				EXISTS (
					SELECT 1 FROM project_security_gate_checklists AS pending_checklist
					WHERE pending_checklist.security_gate_id = project_security_gates.id
						AND pending_checklist.response_status_id IS NULL
						AND pending_checklist.is_exception = 0
				)`)
		case "risk_raised":
			query = query.Where(`
				EXISTS (
					SELECT 1 FROM project_security_gate_checklist_risks AS checklist_risk
					JOIN project_security_gate_checklists AS risk_checklist
						ON risk_checklist.id = checklist_risk.checklist_id
					WHERE risk_checklist.security_gate_id = project_security_gates.id
				)`)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Page < 1 {
		params.Page = 1
	}

	if params.PageSize < 1 {
		params.PageSize = 10
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// An export takes every matching row, not just the current page
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&projectSecurityGates).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &projectDTOs.ProjectSecurityGateListResponseDTO{
		Items: projectSecurityGates,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProjectSecurityGateRepository) GetSummary() (*projectDTOs.ProjectSecurityGateSummaryDTO, error) {
	type summaryRow struct {
		Total  int64
		Draft  int64
		GoLive int64
	}

	var row summaryRow
	err := repo.db.Table("project_security_gates").
		Select(`
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN readiness_status.type = 'draft' THEN 1 ELSE 0 END), 0) AS draft,
			COALESCE(SUM(CASE WHEN readiness_status.type = 'go_live' THEN 1 ELSE 0 END), 0) AS go_live
		`).
		Joins(`
			LEFT JOIN project_security_gate_readiness_statuses AS readiness_status
				ON readiness_status.id = project_security_gates.readiness_status_id
		`).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &projectDTOs.ProjectSecurityGateSummaryDTO{
		Total:  row.Total,
		Draft:  row.Draft,
		GoLive: row.GoLive,
	}, nil
}

func (repo *ProjectSecurityGateRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var projectSecurityGate *models.ProjectSecurityGate
	result := db.First(&projectSecurityGate, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectSecurityGate, nil
}

func (repo *ProjectSecurityGateRepository) GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error) {
	var projectSecurityGate *models.ProjectSecurityGate
	result := repo.db.
		Preload("Stage").
		Preload("TemplateVersion").
		Preload("BusinessOwner").
		Preload("Division").
		Preload("Department").
		Preload("ReadinessStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Classifications").
		Preload("Classifications.Classification").
		First(&projectSecurityGate, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectSecurityGate, nil
}

func (repo *ProjectSecurityGateRepository) Create(projectSecurityGate *models.ProjectSecurityGate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(projectSecurityGate).Error
}

func (repo *ProjectSecurityGateRepository) Update(projectSecurityGate *models.ProjectSecurityGate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(projectSecurityGate).Error
}

func (repo *ProjectSecurityGateRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.ProjectSecurityGate{}, id).Error
}

func (repo *ProjectSecurityGateRepository) ReplaceClassificationMappings(securityGateID uint, classificationIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("security_gate_id = ?", securityGateID).Delete(&models.ProjectSecurityGateClassificationMapping{}).Error; err != nil {
		return err
	}

	if len(classificationIDs) == 0 {
		return nil
	}

	mappings := make([]models.ProjectSecurityGateClassificationMapping, 0, len(classificationIDs))
	for _, classificationID := range classificationIDs {
		mappings = append(mappings, models.ProjectSecurityGateClassificationMapping{
			SecurityGateID:   securityGateID,
			ClassificationID: classificationID,
		})
	}

	return db.Create(&mappings).Error
}

func (repo *ProjectSecurityGateRepository) DeleteClassificationMappingsBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("security_gate_id = ?", securityGateID).Delete(&models.ProjectSecurityGateClassificationMapping{}).Error
}
