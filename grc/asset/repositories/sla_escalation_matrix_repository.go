package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type SLAEscalationMatrixRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.SLAEscalationMatrix, error)
	GetByID(id uint) (*models.SLAEscalationMatrix, error)
	Create(matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error)
	Update(id uint, matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error)
	Delete(id uint) error

	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)

	ExistsBySeverityID(severityID uint) (bool, error)
	ExistsBySeverityIDExcludingID(id uint, severityID uint) (bool, error)
}

type SLAEscalationMatrixRepository struct {
	db *gorm.DB
}

func NewSLAEscalationMatrixRepository(db *gorm.DB) *SLAEscalationMatrixRepository {
	return &SLAEscalationMatrixRepository{db: db}
}

func (repo *SLAEscalationMatrixRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.SLAEscalationMatrix, error) {
	var matrices []*models.SLAEscalationMatrix
	query := repo.db.Model(&models.SLAEscalationMatrix{})

	if keyword != "" {
		query = query.Where("status LIKE ?", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&matrices)

	return matrices, result.Error
}

func (repo *SLAEscalationMatrixRepository) GetByID(id uint) (*models.SLAEscalationMatrix, error) {
	var matrix models.SLAEscalationMatrix
	result := repo.db.First(&matrix, id)
	return &matrix, result.Error
}

func (repo *SLAEscalationMatrixRepository) Create(matrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error) {
	result := repo.db.Create(matrix)
	return matrix, result.Error
}

func (repo *SLAEscalationMatrixRepository) Update(id uint, updatedMatrix *models.SLAEscalationMatrix) (*models.SLAEscalationMatrix, error) {
	var matrix models.SLAEscalationMatrix
	if err := repo.db.First(&matrix, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&matrix).Updates(updatedMatrix).Error; err != nil {
		return nil, err
	}

	return &matrix, nil
}

func (repo *SLAEscalationMatrixRepository) Delete(id uint) error {
	result := repo.db.Delete(&models.SLAEscalationMatrix{}, id)
	return result.Error
}

func (repo *SLAEscalationMatrixRepository) ExistsBySeverityID(severityID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SLAEscalationMatrix{}).
		Where("severity_id = ?", severityID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *SLAEscalationMatrixRepository) ExistsBySeverityIDExcludingID(id uint, severityID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.SLAEscalationMatrix{}).
		Where("id <> ? AND severity_id = ?", id, severityID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *SLAEscalationMatrixRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SLAEscalationMatrix{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SLAEscalationMatrix{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SLAEscalationMatrixRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SLAEscalationMatrix{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SLAEscalationMatrix{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}