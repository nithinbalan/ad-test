package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentRequirementRepositoryInterface interface {
	GetAll(documentRequirementGroupID uint, sort string, sortBy string, keyword string) ([]*models.DocumentRequirement, error)
	GetByID(id uint) (*models.DocumentRequirement, error)
	Create(requirementGroup *models.DocumentRequirement) (uint, error)
	Update(id uint, requirementGroup *models.DocumentRequirement) (*models.DocumentRequirement, error)
	Delete(id uint) error
}

// Implementation
type DocumentRequirementRepository struct {
	db *gorm.DB
}

func NewDocumentRequirementRepository(db *gorm.DB) *DocumentRequirementRepository {
	return &DocumentRequirementRepository{
		db: db,
	}
}

func (repo *DocumentRequirementRepository) GetAll(documentRequirementGroupID uint, sort string, sortBy string, keyword string) ([]*models.DocumentRequirement, error) {
	documentRequirements := []*models.DocumentRequirement{}
	query := repo.db.Preload("DocumentRequirementGroup")

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Where("document_requirement_group_id = ?", documentRequirementGroupID).Order(sortBy + " " + sort)

	result := query.Find(&documentRequirements)

	return documentRequirements, result.Error
}

func (repo *DocumentRequirementRepository) GetByID(id uint) (*models.DocumentRequirement, error) {
	var requirementGroup *models.DocumentRequirement

	result := repo.db.
		Preload("DocumentRequirementGroup").
		Preload("DocumentRequirementGroup.Document").
		First(&requirementGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return requirementGroup, nil
}

func (repo *DocumentRequirementRepository) Create(requirementGroup *models.DocumentRequirement) (uint, error) {
	result := repo.db.Create(&requirementGroup)

	return requirementGroup.ID, result.Error
}

func (repo *DocumentRequirementRepository) Update(id uint, updatedCategory *models.DocumentRequirement) (*models.DocumentRequirement, error) {
	var requirementGroup *models.DocumentRequirement
	if err := repo.db.First(&requirementGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&requirementGroup).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return requirementGroup, nil
}

func (repo *DocumentRequirementRepository) Delete(id uint) error {
	var requirementGroup *models.DocumentRequirement
	if err := repo.db.First(&requirementGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.DocumentRequirement{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
