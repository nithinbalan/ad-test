package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ViolationReportStatus, error)
	GetByID(id uint) (*models.ViolationReportStatus, error)
	GetByType(typeValue string) (*models.ViolationReportStatus, error)
}

type ViolationReportStatusRepository struct {
	db *gorm.DB
}

func NewViolationReportStatusRepository(db *gorm.DB) *ViolationReportStatusRepository {
	return &ViolationReportStatusRepository{
		db: db,
	}
}

func (repo *ViolationReportStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ViolationReportStatus, error) {

	var violationReportStatuses []*models.ViolationReportStatus

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
	result := query.Find(&violationReportStatuses)

	return violationReportStatuses, result.Error
}

func (repo *ViolationReportStatusRepository) GetByID(id uint) (*models.ViolationReportStatus, error) {
	var status *models.ViolationReportStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *ViolationReportStatusRepository) GetByType(typeValue string) (*models.ViolationReportStatus, error) {
	var violationReportStatus *models.ViolationReportStatus

	result := repo.db.Where("type = ?", typeValue).First(&violationReportStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Violation Report Status with type "+typeValue)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return violationReportStatus, nil
}
