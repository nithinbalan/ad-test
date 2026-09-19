package repositories

import (
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	Create(project *models.Project, tx ...utility.DatabaseTransaction) (uint, error)
	Update(project *models.Project, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	ExistsByIncidentActionPlanID(incidentActionPlanID uint) (bool, error)
	ExistsByIncidentActionPlanIDExcludingID(id uint, incidentActionPlanID uint) (bool, error)
}

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(r.db)
}

func (r *ProjectRepository) GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error) {
	var projects []*dtos.ProjectListDTO
	var totalCount int64

	query := r.db.Table("projects").
		Select(`
			projects.id,
			projects.title_en,
			projects.title_ar,
			projects.description_en,
			projects.description_ar,
			projects.budget,
			projects.scope,
			projects.number_of_resources,
			projects.completion_percentage,
			CONVERT(VARCHAR, projects.start_date, 23) AS start_date,
			CONVERT(VARCHAR, projects.end_date, 23) AS end_date,
			project_types.id AS project_type_id,
			project_types.title_en AS project_type_title_en,
			project_types.title_ar AS project_type_title_ar,
			project_statuses.id AS project_status_id,
			project_statuses.title_en AS project_status_title_en,
			project_statuses.title_ar AS project_status_title_ar,
			project_phases.id AS project_phase_id,
			project_phases.title_en AS project_phase_title_en,
			project_phases.title_ar AS project_phase_title_ar,
			project_phases.color_code AS project_phase_color_code,
			project_classifications.id AS project_classification_id,
			project_classifications.title_en AS project_classification_title_en,
			project_classifications.title_ar AS project_classification_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			projects.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			projects.next_review_level,
			CONVERT(VARCHAR, projects.submitted_at, 120) AS submitted_at,
			submitted_user.id AS submitted_user_id,
			submitted_user.first_name_en AS submitted_user_first_name_en,
			submitted_user.first_name_ar AS submitted_user_first_name_ar,
			submitted_user.last_name_en AS submitted_user_last_name_en,
			submitted_user.last_name_ar AS submitted_user_last_name_ar,
			CONVERT(VARCHAR, projects.approved_at, 120) AS approved_at,
			approved_user.id AS approved_user_id,
			approved_user.first_name_en AS approved_user_first_name_en,
			approved_user.first_name_ar AS approved_user_first_name_ar,
			approved_user.last_name_en AS approved_user_last_name_en,
			approved_user.last_name_ar AS approved_user_last_name_ar,
			CONVERT(VARCHAR, projects.rejected_at, 120) AS rejected_at,
			rejected_user.id AS rejected_user_id,
			rejected_user.first_name_en AS rejected_user_first_name_en,
			rejected_user.first_name_ar AS rejected_user_first_name_ar,
			rejected_user.last_name_en AS rejected_user_last_name_en,
			rejected_user.last_name_ar AS rejected_user_last_name_ar,
			projects.incident_action_plan_id,
			incident_action_plans.title AS incident_action_plan_title
		`).
		Joins("LEFT JOIN project_types ON project_types.id = projects.project_type_id").
		Joins("LEFT JOIN project_statuses ON project_statuses.id = projects.project_status_id").
		Joins("LEFT JOIN project_phases ON project_phases.id = projects.project_phase_id").
		Joins("LEFT JOIN project_classifications ON project_classifications.id = projects.project_classification_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = projects.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = projects.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = projects.updated_by_id").
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = projects.workflow_status_id").
		Joins("LEFT JOIN users AS submitted_user ON submitted_user.id = projects.submitted_by_id").
		Joins("LEFT JOIN users AS approved_user ON approved_user.id = projects.approved_by_id").
		Joins("LEFT JOIN users AS rejected_user ON rejected_user.id = projects.rejected_by_id").
		Joins("LEFT JOIN incident_action_plans ON incident_action_plans.id = projects.incident_action_plan_id")

	// Apply dynamic filters
	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Quick filters can be implemented here if needed
	// Quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "pending":
			query = query.Where("workflow_statuses.type = ?", "pending")
		case "approved":
			query = query.Where("workflow_statuses.type = ?", "approve")
		case "rejected":
			query = query.Where("workflow_statuses.type = ?", "reject")
		case "reverted":
			query = query.Where("workflow_statuses.type = ?", "revert")
		case "draft":
			query = query.Where("projects.workflow_status_id IS NULL")
		}
	}

	// Status filter
	if params.Status != "" {
		query = query.Where("project_statuses.type = ?", params.Status)
	}

	// Keyword search
	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
			projects.title_en LIKE ? OR
			projects.title_ar LIKE ? OR
			projects.description_en LIKE ? OR
			projects.description_ar LIKE ? OR
			project_types.title_en LIKE ? OR
			project_types.title_ar LIKE ? OR
			project_statuses.title_en LIKE ? OR
			project_statuses.title_ar LIKE ? OR
			project_phases.title_en LIKE ? OR
			project_phases.title_ar LIKE ? OR
			project_classifications.title_en LIKE ? OR
			project_classifications.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ?`,
			like, like, like, like, like, like, like, like,
			like, like, like, like,
			like, like, like, like,
			like)
	}

	// Count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export && params.PageSize > 0 {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &dtos.ProjectListResponseDTO{
		Items: projects,
		Pagination: dtos.ProjectPaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (r *ProjectRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var project models.Project
	if err := db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var project models.Project

	result := db.
		Preload("ResponsibleUser").
		Preload("ProjectType").
		Preload("ProjectStatus").
		Preload("ProjectPhase").
		Preload("ProjectClassification").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("IncidentActionPlan").
		Preload("IncidentActionPlan.ActionPlanStatus").
		Preload("IncidentActionPlan.ResponsibleUser").
		Preload("IncidentActionPlan.IncidentRegister").
		First(&project, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (r *ProjectRepository) Create(project *models.Project, tx ...utility.DatabaseTransaction) (uint, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(project).Error; err != nil {
		return 0, err
	}
	return project.ID, nil
}

func (r *ProjectRepository) Update(project *models.Project, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Omit(clause.Associations).Save(project).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) Delete(id uint) error {
	var project *models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&models.Project{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) ExistsByIncidentActionPlanID(incidentActionPlanID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).
		Where("incident_action_plan_id = ?", incidentActionPlanID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ProjectRepository) ExistsByIncidentActionPlanIDExcludingID(id uint, incidentActionPlanID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).
		Where("id <> ? AND incident_action_plan_id = ?", id, incidentActionPlanID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
