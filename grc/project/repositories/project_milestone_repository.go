package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/project/dtos"

	"gorm.io/gorm"
)

type ProjectMilestoneRepositoryInterface interface {
	GetAll(projectID uint, params *dtos.ProjectMilestoneListQueryParamsDTO) (*dtos.ProjectMilestoneListResponseDTO, error)
	GetByID(id uint) (*models.ProjectMilestone, error)
	GetFullDetailsByID(id uint) (*models.ProjectMilestone, error)
	Create(milestone *models.ProjectMilestone) (uint, error)
	Update(milestone *models.ProjectMilestone) error
	Delete(id uint) error
}

type ProjectMilestoneRepository struct {
	db *gorm.DB
}

func NewProjectMilestoneRepository(db *gorm.DB) *ProjectMilestoneRepository {
	return &ProjectMilestoneRepository{db: db}
}

func (repo *ProjectMilestoneRepository) GetAll(projectID uint, params *dtos.ProjectMilestoneListQueryParamsDTO) (*dtos.ProjectMilestoneListResponseDTO, error) {
	var milestones []*dtos.ProjectMilestoneListDTO
	var totalCount int64

	query := repo.db.Table("project_milestones").
		Select(`
			project_milestones.id,
			project_milestones.project_id,
			project_milestones.title_en,
			project_milestones.title_ar,
			project_milestones.description_en,
			project_milestones.description_ar,
			CONVERT(VARCHAR, project_milestones.planned_date, 23) AS planned_date,
			CONVERT(VARCHAR, project_milestones.actual_date, 23) AS actual_date,
			project_milestones.project_milestone_status_id,
			status.title_en AS project_milestone_status_en,
			status.title_ar AS project_milestone_status_ar,
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
		Joins(`LEFT JOIN projects ON project_milestones.project_id = projects.id`).
		Joins(`LEFT JOIN project_milestone_statuses AS status ON project_milestones.project_milestone_status_id = status.id`).
		Joins(`LEFT JOIN users AS created_user ON project_milestones.created_by_id = created_user.id`).
		Joins(`LEFT JOIN users AS updated_user ON project_milestones.updated_by_id = updated_user.id`)

	// Filter by project_id if provided
	if projectID != 0 {
		query = query.Where("project_milestones.project_id = ?", projectID)
	}

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_milestones.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_milestones.title_en LIKE ? OR
			project_milestones.title_ar LIKE ? OR
			projects.title_en LIKE ? OR
			projects.title_ar LIKE ?)
		`, keyword, keyword, keyword, keyword)
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&milestones).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.ProjectMilestoneListResponseDTO{
		Items: milestones,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ProjectMilestoneRepository) GetByID(id uint) (*models.ProjectMilestone, error) {
	var milestone models.ProjectMilestone
	result := repo.db.First(&milestone, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &milestone, nil
}

func (repo *ProjectMilestoneRepository) GetFullDetailsByID(id uint) (*models.ProjectMilestone, error) {
	var milestone models.ProjectMilestone
	result := repo.db.
		Preload("Project").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ProjectMilestoneStatus").
		First(&milestone, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &milestone, nil
}

func (repo *ProjectMilestoneRepository) Create(milestone *models.ProjectMilestone) (uint, error) {
	if err := repo.db.Create(milestone).Error; err != nil {
		return 0, err
	}
	return milestone.ID, nil
}

func (repo *ProjectMilestoneRepository) Update(milestone *models.ProjectMilestone) error {
	if err := repo.db.Save(milestone).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProjectMilestoneRepository) Delete(id uint) error {
	var milestone *models.ProjectMilestone
	if err := repo.db.First(&milestone, id).Error; err != nil {
		return err
	}

	if err := repo.db.Delete(&models.ProjectMilestone{}, id).Error; err != nil {
		return err
	}

	return nil
}
