package repositories

import (
	"gorm.io/gorm"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/masters/dtos"
	"jethur/grc/models"
)

type ActionPlanMasterRepositoryInterface interface {
	GetAll(params *dtos.ActionPlanMasterListQueryParamsDTO) (*dtos.ActionPlanMasterListResponseDTO, error)
	GetByID(id uint) (*models.ActionPlanMaster, error)
	GetFullDetailsByID(id uint) (*models.ActionPlanMaster, error)
	GetIncidentCategoryByID(id uint) (*models.IncidentCategory, error)
	GetModuleByID(id uint) (*models.Module, error)
	ExistsByTitle(title string) (bool, error)
	ExistsByTitleExcludingID(id uint, title string, moduleID uint) (bool, error)
	Create(actionPlanMaster *models.ActionPlanMaster) (uint, error)
	Update(actionPlanMaster *models.ActionPlanMaster) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

type ActionPlanMasterRepository struct {
	db *gorm.DB
}

func NewActionPlanMasterRepository(db *gorm.DB) *ActionPlanMasterRepository {
	return &ActionPlanMasterRepository{
		db: db,
	}
}

func (repo *ActionPlanMasterRepository) GetAll(params *dtos.ActionPlanMasterListQueryParamsDTO) (*dtos.ActionPlanMasterListResponseDTO, error) {

	var actionPlanMasterListDTO []*dtos.ActionPlanMasterListDTO
	var totalCount int64

	query := repo.db.Table("action_plan_master").
		Select(`
			action_plan_master.id,
			action_plan_master.title,
			action_plan_master.description,
			action_plan_master.status,

			module.id AS module_id,
			module.identifier AS module_identifier,
			module.title_en AS module_title_en,
			module.title_ar AS module_title_ar,

			incident_category.id AS incident_category_id,
			incident_category.title_en AS incident_category_title_en,
			incident_category.title_ar AS incident_category_title_ar,

			severity_level.id AS severity_level_id,
			severity_level.title_en AS severity_title_en,
			severity_level.title_ar AS severity_title_ar,
			severity_level.color_code AS severity_color_code,
			severity_level.type AS severity_type
		`).
		Joins("INNER JOIN modules module ON action_plan_master.module_id = module.id").
		Joins("LEFT JOIN incident_categories incident_category ON action_plan_master.incident_category_id = incident_category.id").
		Joins("LEFT JOIN severity_levels severity_level ON action_plan_master.severity_level_id = severity_level.id")
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"(action_plan_master.title LIKE ? OR action_plan_master.description LIKE ?)",
			keyword,
			keyword,
		)
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("action_plan_master.status = ?", params.Status)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "action_plan_master.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Find(&actionPlanMasterListDTO).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.ActionPlanMasterListResponseDTO{
		Items: actionPlanMasterListDTO,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ActionPlanMasterRepository) GetByID(id uint) (*models.ActionPlanMaster, error) {
	var actionPlanMaster models.ActionPlanMaster

	if err := repo.db.First(&actionPlanMaster, id).Error; err != nil {
		return nil, err
	}

	return &actionPlanMaster, nil
}

func (repo *ActionPlanMasterRepository) GetFullDetailsByID(id uint) (*models.ActionPlanMaster, error) {
	var actionPlanMaster models.ActionPlanMaster

	result := repo.db.
		Preload("Module").
		Preload("IncidentCategory").
		Preload("SeverityLevel").
		First(&actionPlanMaster, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &actionPlanMaster, nil
}

func (repo *ActionPlanMasterRepository) GetIncidentCategoryByID(id uint) (*models.IncidentCategory, error) {
	var incidentCategory *models.IncidentCategory
	result := repo.db.First(&incidentCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentCategory, nil
}

func (repo *ActionPlanMasterRepository) GetModuleByID(id uint) (*models.Module, error) {
	var module *models.Module
	result := repo.db.First(&module, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

func (repo *ActionPlanMasterRepository) Create(actionPlanMaster *models.ActionPlanMaster) (uint, error) {

	if err := repo.db.Create(actionPlanMaster).Error; err != nil {
		return 0, err
	}

	return actionPlanMaster.ID, nil
}

func (repo *ActionPlanMasterRepository) Update(actionPlan *models.ActionPlanMaster) error {
	if err := repo.db.
		Model(&models.ActionPlanMaster{}).
		Where("id = ?", actionPlan.ID).
		Updates(actionPlan).
		Error; err != nil {
		return err
	}

	return nil
}

func (repo *ActionPlanMasterRepository) Delete(id uint) error {

	var actionPlan *models.ActionPlanMaster
	if err := repo.db.First(&actionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ActionPlanMaster{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ActionPlanMasterRepository) ExistsByTitle(title string) (bool, error) {

	var count int64

	err := repo.db.Model(&models.ActionPlanMaster{}).
		Where("title = ?", title).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ActionPlanMasterRepository) ExistsByTitleExcludingID(id uint, title string, moduleID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ActionPlanMaster{}).
		Where("id <> ? AND title = ? AND module_id = ?", id, title, moduleID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ActionPlanMasterRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ActionPlanMaster{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ActionPlanMaster{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ActionPlanMasterRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.ActionPlanMaster{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ActionPlanMaster{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}
