package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateClassification, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateClassification, error)
	GetByType(classificationType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateClassification, error)
}

type ProjectSecurityGateClassificationRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateClassificationRepository(db *gorm.DB) *ProjectSecurityGateClassificationRepository {
	return &ProjectSecurityGateClassificationRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateClassification, error) {
	var classifications []*models.ProjectSecurityGateClassification
	query := repo.db

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&classifications)

	return classifications, result.Error
}

func (repo *ProjectSecurityGateClassificationRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateClassification, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var classification *models.ProjectSecurityGateClassification
	result := db.First(&classification, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, result.Error
}

func (repo *ProjectSecurityGateClassificationRepository) GetByType(classificationType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateClassification, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var classification *models.ProjectSecurityGateClassification

	result := db.Where("type = ?", classificationType).First(&classification)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Security Gate Classification with type "+classificationType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, nil
}
