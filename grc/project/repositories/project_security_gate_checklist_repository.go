package repositories

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateChecklists(checklists []*models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error
	GetSummaryBySecurityGateID(securityGateID uint) (*projectDTOs.ProjectSecurityGateChecklistSummaryDTO, error)
	GetAll(securityGateID uint, params *projectDTOs.ProjectSecurityGateChecklistListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateChecklistListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error)
	GetFullDetailsByID(id uint) (*models.ProjectSecurityGateChecklist, error)
	DeleteBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error
	UpdateResponse(checklist *models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error
	CreateResponseDocuments(documents []*models.ProjectSecurityGateChecklistResponseDocument, tx ...utility.DatabaseTransaction) error
	GetResponseDocumentByID(id uint) (*models.ProjectSecurityGateChecklistResponseDocument, error)
	DeleteResponseDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error

	DeleteChecklistControls(checklistID uint, tx ...utility.DatabaseTransaction) error
	CreateChecklistControls(controls []*models.ProjectSecurityGateChecklistControl, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateChecklistRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistRepository(db *gorm.DB) *ProjectSecurityGateChecklistRepository {
	return &ProjectSecurityGateChecklistRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateChecklistRepository) CreateChecklists(checklists []*models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(checklists) == 0 {
		return nil
	}

	return db.Create(&checklists).Error
}

func (repo *ProjectSecurityGateChecklistRepository) GetAll(securityGateID uint, params *projectDTOs.ProjectSecurityGateChecklistListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateChecklistListResponseDTO, error) {
	var checklists []*projectDTOs.ProjectSecurityGateChecklistListDTO
	var totalCount int64

	query := repo.db.Table("project_security_gate_checklists").
		Select(`
			project_security_gate_checklists.id,
			project_security_gate_checklists.security_gate_id,
			project_security_gate_checklists.requirement_id,
			requirements.title_en AS title_en,
			requirements.title_ar AS title_ar,
			requirements.is_critical AS is_critical,
			project_security_gate_checklists.response_status_id AS response_status_id,
			response_status.title_en AS response_status_title_en,
			response_status.title_ar AS response_status_title_ar,
			response_status.type AS response_status_type,
			response_status.color_code AS response_status_color_code,
			project_security_gate_checklists.comment AS response_comment,
			project_security_gate_checklists.is_exception AS response_is_exception,
			project_security_gate_checklists.exception_justification AS response_exception_justification,
			project_security_gate_checklists.created_by_id AS response_created_by_id,
			response_created_user.first_name_en AS response_created_user_first_name_en,
			response_created_user.first_name_ar AS response_created_user_first_name_ar,
			response_created_user.last_name_en AS response_created_user_last_name_en,
			response_created_user.last_name_ar AS response_created_user_last_name_ar,
			response_created_user.initial AS response_created_user_initial,
			response_created_user.color AS response_created_user_color,
			response_created_user.email AS response_created_user_email,
			response_created_user.profile_image_path AS response_created_user_profile_image_path,
			project_security_gate_checklists.created_at AS response_created_at,
			project_security_gate_checklists.updated_by_id AS response_updated_by_id,
			response_updated_user.first_name_en AS response_updated_user_first_name_en,
			response_updated_user.first_name_ar AS response_updated_user_first_name_ar,
			response_updated_user.last_name_en AS response_updated_user_last_name_en,
			response_updated_user.last_name_ar AS response_updated_user_last_name_ar,
			response_updated_user.initial AS response_updated_user_initial,
			response_updated_user.color AS response_updated_user_color,
			response_updated_user.email AS response_updated_user_email,
			response_updated_user.profile_image_path AS response_updated_user_profile_image_path,
			project_security_gate_checklists.updated_at AS response_updated_at
		`).
		Joins(`
			JOIN project_security_gate_checklist_template_requirements AS requirements
				ON requirements.id = project_security_gate_checklists.requirement_id
			LEFT JOIN project_security_gate_checklist_response_statuses AS response_status
				ON response_status.id = project_security_gate_checklists.response_status_id
			LEFT JOIN users AS response_created_user
				ON response_created_user.id = project_security_gate_checklists.created_by_id
			LEFT JOIN users AS response_updated_user
				ON response_updated_user.id = project_security_gate_checklists.updated_by_id
		`).
		Where("project_security_gate_checklists.security_gate_id = ?", securityGateID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("requirements.title_en LIKE ? OR requirements.title_ar LIKE ?", keyword, keyword)
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

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&checklists).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &projectDTOs.ProjectSecurityGateChecklistListResponseDTO{
		Items: checklists,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProjectSecurityGateChecklistRepository) GetSummaryBySecurityGateID(securityGateID uint) (*projectDTOs.ProjectSecurityGateChecklistSummaryDTO, error) {
	type summaryRow struct {
		Total     int64
		Critical  int64
		Completed int64
	}

	var row summaryRow
	err := repo.db.Table("project_security_gate_checklists").
		Select(`
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN requirements.is_critical = 1 THEN 1 ELSE 0 END), 0) AS critical,
			COALESCE(SUM(CASE WHEN project_security_gate_checklists.created_by_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS completed
		`).
		Joins(`
			JOIN project_security_gate_checklist_template_requirements AS requirements
				ON requirements.id = project_security_gate_checklists.requirement_id
		`).
		Where("project_security_gate_checklists.security_gate_id = ?", securityGateID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	outstanding := row.Total - row.Completed

	var completionPercentage float64
	if row.Total > 0 {
		completionPercentage = float64(row.Completed) / float64(row.Total) * 100
	}

	return &projectDTOs.ProjectSecurityGateChecklistSummaryDTO{
		Total:                row.Total,
		Critical:             row.Critical,
		Completed:            row.Completed,
		Outstanding:          outstanding,
		CompletionPercentage: completionPercentage,
	}, nil
}

func (repo *ProjectSecurityGateChecklistRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var checklist *models.ProjectSecurityGateChecklist
	result := db.
		Preload("Documents").
		First(&checklist, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklist, nil
}

func (repo *ProjectSecurityGateChecklistRepository) GetFullDetailsByID(id uint) (*models.ProjectSecurityGateChecklist, error) {
	var checklist *models.ProjectSecurityGateChecklist
	result := repo.db.
		Preload("Requirement").
		Preload("ResponseStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_security_gate_checklist_comments.created_at DESC")
		}).
		Preload("Comments.CreatedByUser").
		Preload("Comments.Documents").
		Preload("Controls").
		Preload("Controls.Control").
		First(&checklist, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklist, nil
}

func (repo *ProjectSecurityGateChecklistRepository) DeleteBySecurityGateID(securityGateID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("security_gate_id = ?", securityGateID).Delete(&models.ProjectSecurityGateChecklist{}).Error
}

func (repo *ProjectSecurityGateChecklistRepository) UpdateResponse(checklist *models.ProjectSecurityGateChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(checklist).Error
}

func (repo *ProjectSecurityGateChecklistRepository) CreateResponseDocuments(documents []*models.ProjectSecurityGateChecklistResponseDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(documents) == 0 {
		return nil
	}

	return db.Create(&documents).Error
}

func (repo *ProjectSecurityGateChecklistRepository) GetResponseDocumentByID(id uint) (*models.ProjectSecurityGateChecklistResponseDocument, error) {
	var document *models.ProjectSecurityGateChecklistResponseDocument
	result := repo.db.First(&document, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return document, nil
}

func (repo *ProjectSecurityGateChecklistRepository) DeleteResponseDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(documentIDs) == 0 {
		return nil
	}

	return db.Where("id IN (?)", documentIDs).Delete(&models.ProjectSecurityGateChecklistResponseDocument{}).Error
}

func (repo *ProjectSecurityGateChecklistRepository) DeleteChecklistControls(checklistID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("checklist_id = ?", checklistID).Delete(&models.ProjectSecurityGateChecklistControl{}).Error
}

func (repo *ProjectSecurityGateChecklistRepository) CreateChecklistControls(controls []*models.ProjectSecurityGateChecklistControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(controls) == 0 {
		return nil
	}

	return db.Create(&controls).Error
}
