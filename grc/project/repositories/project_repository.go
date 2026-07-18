package repositories

import (
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ProjectListQueryParamsDTO) (*dtos.ProjectListResponseDTO, error)
	GetByID(id uint) (*models.Project, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	Create(project *models.Project, tx ...utility.DatabaseTransaction) (uint, error)
	Update(project *models.Project, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
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
			CONVERT(VARCHAR, projects.start_date, 23) AS start_date,
			CONVERT(VARCHAR, projects.end_date, 23) AS end_date,
			project_types.id AS project_type_id,
			project_types.title_en AS project_type_title_en,
			project_types.title_ar AS project_type_title_ar,
			project_statuses.id AS project_status_id,
			project_statuses.title_en AS project_status_title_en,
			project_statuses.title_ar AS project_status_title_ar,
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
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN project_types ON project_types.id = projects.project_type_id").
		Joins("LEFT JOIN project_statuses ON project_statuses.id = projects.project_status_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = projects.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = projects.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = projects.updated_by_id")

	// Apply dynamic filters
	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	// Quick filters can be implemented here if needed

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
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ?`,
			like, like, like, like,
			like, like, like, like,
			like, like, like, like,
			like)
	}

	// Count before pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
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

func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
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
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
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
	if err := db.Save(project).Error; err != nil {
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
