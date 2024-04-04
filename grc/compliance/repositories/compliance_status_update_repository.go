package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"

	"gorm.io/gorm"
)

type ComplianceStatusUpdateRepositoryInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error)
	Create(complianceUpdate *models.ComplianceStatusUpdate) (uint, error)
}

type ComplianceStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceStatusUpdateRepository(db *gorm.DB) *ComplianceStatusUpdateRepository {
	return &ComplianceStatusUpdateRepository{
		db: db,
	}
}

func (repo *ComplianceStatusUpdateRepository) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error) {

	var complianceUpdates []*complianceDTOs.ComplianceStatusUpdateListDTO

	query := repo.db.Table("compliance_status_updates").
		Select(`
            compliance_status_updates.*,
			compliances.title AS compliance_title,
			compliances.description AS compliance_description,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar   
        `).
		Joins(`
		    LEFT JOIN compliances ON compliances.id = compliance_status_updates.compliance_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_status_updates.compliance_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_status_updates.created_by
		`).
		Where("compliance_status_updates.compliance_id = ?", complianceID)

	if keyword != "" {
		query = query.Where(`
			(compliances.title LIKE ? OR
			compliance_status_updates.comment LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceUpdates)

	return complianceUpdates, result.Error
}

func (repo *ComplianceStatusUpdateRepository) Create(complianceUpdate *models.ComplianceStatusUpdate) (uint, error) {

	if err := repo.db.Create(complianceUpdate).Error; err != nil {
		return 0, err
	}

	return complianceUpdate.ID, nil
}
