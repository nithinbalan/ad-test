package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateReadinessStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateReadinessStatus, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error)
	GetByType(readinessStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error)
}

type ProjectSecurityGateReadinessStatusRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateReadinessStatusRepository(db *gorm.DB) *ProjectSecurityGateReadinessStatusRepository {
	return &ProjectSecurityGateReadinessStatusRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateReadinessStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateReadinessStatus, error) {
	var readinessStatuses []*models.ProjectSecurityGateReadinessStatus
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
	result := query.Find(&readinessStatuses)

	return readinessStatuses, result.Error
}

func (repo *ProjectSecurityGateReadinessStatusRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var readinessStatus *models.ProjectSecurityGateReadinessStatus
	result := db.First(&readinessStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return readinessStatus, result.Error
}

func (repo *ProjectSecurityGateReadinessStatusRepository) GetByType(readinessStatusType string, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReadinessStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var readinessStatus *models.ProjectSecurityGateReadinessStatus

	result := db.Where("type = ?", readinessStatusType).First(&readinessStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Security Gate Readiness Status with type "+readinessStatusType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return readinessStatus, nil
}
