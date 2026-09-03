package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistResponseStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateChecklistResponseStatus, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error)
	GetByType(responseStatusType string) (*models.ProjectSecurityGateChecklistResponseStatus, error)
}

type ProjectSecurityGateChecklistResponseStatusRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistResponseStatusRepository(db *gorm.DB) *ProjectSecurityGateChecklistResponseStatusRepository {
	return &ProjectSecurityGateChecklistResponseStatusRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateChecklistResponseStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateChecklistResponseStatus, error) {
	var responseStatuses []*models.ProjectSecurityGateChecklistResponseStatus
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
	result := query.Find(&responseStatuses)

	return responseStatuses, result.Error
}

func (repo *ProjectSecurityGateChecklistResponseStatusRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var responseStatus *models.ProjectSecurityGateChecklistResponseStatus
	result := db.First(&responseStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return responseStatus, result.Error
}

func (repo *ProjectSecurityGateChecklistResponseStatusRepository) GetByType(responseStatusType string) (*models.ProjectSecurityGateChecklistResponseStatus, error) {
	var responseStatus *models.ProjectSecurityGateChecklistResponseStatus

	result := repo.db.Where("type = ?", responseStatusType).First(&responseStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Security Gate Checklist Response Status with type "+responseStatusType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return responseStatus, nil
}
