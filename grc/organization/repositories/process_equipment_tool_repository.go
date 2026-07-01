package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessEquipmentToolRepositoryInterface interface {
	GetAll(params *dtos.ProcessEquipmentToolListQueryParamsDTO) ([]*models.ProcessEquipmentTool, error)
	GetByID(id uint) (*models.ProcessEquipmentTool, error)
	Create(tool *models.ProcessEquipmentTool) (uint, error)
	Update(tool *models.ProcessEquipmentTool) error
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ProcessEquipmentToolRepository struct {
	db *gorm.DB
}

func NewProcessEquipmentToolRepository(db *gorm.DB) *ProcessEquipmentToolRepository {
	return &ProcessEquipmentToolRepository{db: db}
}

func (repo *ProcessEquipmentToolRepository) GetAll(params *dtos.ProcessEquipmentToolListQueryParamsDTO) ([]*models.ProcessEquipmentTool, error) {
	var tools []*models.ProcessEquipmentTool

	query := repo.db.Model(&models.ProcessEquipmentTool{})

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	return tools, query.Find(&tools).Error
}

func (repo *ProcessEquipmentToolRepository) GetByID(id uint) (*models.ProcessEquipmentTool, error) {
	var tool models.ProcessEquipmentTool
	if err := repo.db.First(&tool, id).Error; err != nil {
		return nil, err
	}
	return &tool, nil
}

func (repo *ProcessEquipmentToolRepository) Create(tool *models.ProcessEquipmentTool) (uint, error) {
	if err := repo.db.Create(tool).Error; err != nil {
		return 0, err
	}
	return tool.ID, nil
}

func (repo *ProcessEquipmentToolRepository) Update(tool *models.ProcessEquipmentTool) error {
	return repo.db.Save(tool).Error
}

func (repo *ProcessEquipmentToolRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessEquipmentTool{}, id).Error
}

func (repo *ProcessEquipmentToolRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessEquipmentTool{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	return count > 0, err
}

func (repo *ProcessEquipmentToolRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessEquipmentTool{}).
		Where("id != ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	return count > 0, err
}
