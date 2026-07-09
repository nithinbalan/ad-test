package repositories

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type SLAEscalationMatrixRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*dtos.SLAEscalationMatrixListDTO, error)
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

func (repo *SLAEscalationMatrixRepository) GetAll(sort, sortBy, keyword, status string) ([]*dtos.SLAEscalationMatrixListDTO, error) {
	var matrices []*dtos.SLAEscalationMatrixListDTO

	query := repo.db.
		Table("sla_escalation_matrix").
		Select(`
			sla_escalation_matrix.id,
			sla_escalation_matrix.severity_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.type AS severity_level_type,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.weight AS severity_level_weight,
			sla_escalation_matrix.remediation_days,
			sla_escalation_matrix.is_email,
			sla_escalation_matrix.is_alert,
			sla_escalation_matrix.role_id,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar,
			sla_escalation_matrix.status
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = sla_escalation_matrix.severity_id").
		Joins("LEFT JOIN roles ON roles.id = sla_escalation_matrix.role_id")

	if keyword != "" {
		keyword = "%" + keyword + "%"

		query = query.Where(`
			sla_escalation_matrix.status LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			roles.title_en LIKE ? OR
			roles.title_ar LIKE ?
		`,
			keyword,
			keyword,
			keyword,
			keyword,
			keyword,
		)
	}

	if status == "active" || status == "inactive" {
		query = query.Where("sla_escalation_matrix.status = ?", status)
	}

	if sortBy == "" {
		sortBy = "sla_escalation_matrix.id"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)

	if err := query.Find(&matrices).Error; err != nil {
		return nil, err
	}

	return matrices, nil
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

	if err := repo.db.Model(&matrix).Updates(map[string]interface{}{
		"severity_id":      updatedMatrix.SeverityID,
		"remediation_days": updatedMatrix.RemediationDays,
		"is_email":         updatedMatrix.IsEmail,
		"is_alert":         updatedMatrix.IsAlert,
		"role_id":          updatedMatrix.RoleID,
		"status":           updatedMatrix.Status,
	}).Error; err != nil {
		return nil, err
	}

	if err := repo.db.First(&matrix, id).Error; err != nil {
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