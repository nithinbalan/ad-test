package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentRequirementGroupRepositoryInterface interface {
	GetAll(documentId uint, sort string, sortBy string, keyword string, status string) ([]*models.DocumentRequirementGroup, error)
	GetByID(id uint) (*models.DocumentRequirementGroup, error)
	Create(requirementGroup *models.DocumentRequirementGroup) (uint, error)
	Update(id uint, requirementGroup *models.DocumentRequirementGroup) (*models.DocumentRequirementGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(documentId uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(documentId uint, id uint, titleEn string, titleAr string) (bool, error)
}

// Implementation
type DocumentRequirementGroupRepository struct {
	db *gorm.DB
}

func NewDocumentRequirementGroupRepository(db *gorm.DB) *DocumentRequirementGroupRepository {
	return &DocumentRequirementGroupRepository{
		db: db,
	}
}

func (repo *DocumentRequirementGroupRepository) GetAll(documentId uint, sort string, sortBy string, keyword string, status string) ([]*models.DocumentRequirementGroup, error) {
	var documentRequirementGroups []*models.DocumentRequirementGroup
	query := repo.db.Preload("Document")

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Where("document_id = ?", documentId).Order(sortBy + " " + sort)

	result := query.Find(&documentRequirementGroups)

	return documentRequirementGroups, result.Error
}

func (repo *DocumentRequirementGroupRepository) GetByID(id uint) (*models.DocumentRequirementGroup, error) {
	var requirementGroup *models.DocumentRequirementGroup

	result := repo.db.Preload("Document").First(&requirementGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return requirementGroup, nil
}

func (repo *DocumentRequirementGroupRepository) Create(requirementGroup *models.DocumentRequirementGroup) (uint, error) {
	result := repo.db.Create(&requirementGroup)

	return requirementGroup.ID, result.Error
}

func (repo *DocumentRequirementGroupRepository) Update(id uint, updatedCategory *models.DocumentRequirementGroup) (*models.DocumentRequirementGroup, error) {
	var requirementGroup *models.DocumentRequirementGroup
	if err := repo.db.First(&requirementGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&requirementGroup).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return requirementGroup, nil
}

func (repo *DocumentRequirementGroupRepository) Delete(id uint) error {
	var requirementGroup *models.DocumentRequirementGroup
	if err := repo.db.First(&requirementGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.DocumentRequirementGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DocumentRequirementGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentRequirementGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentRequirementGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentRequirementGroupRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentRequirementGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentRequirementGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *DocumentRequirementGroupRepository) ExistsByTitle(documentId uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentRequirementGroup{}).
		Where("document_id = ? AND (title_en = ? OR title_ar = ?)", documentId, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *DocumentRequirementGroupRepository) ExistsByTitleExcludingID(documentId, id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentRequirementGroup{}).
		Where("document_id = ? AND id <> ? AND (title_en = ? OR title_ar = ?)", documentId, id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
