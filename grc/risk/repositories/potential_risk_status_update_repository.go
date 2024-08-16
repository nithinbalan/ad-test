package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"

	"gorm.io/gorm"
)

type PotentialRiskStatusUpdateRepositoryInterface interface {
	GetAll(potentialRiskID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.PotentialRiskStatusUpdateListDTO, error)
	Create(potentialRiskUpdate *models.PotentialRiskStatusUpdate) (uint, error)
}

type PotentialRiskStatusUpdateRepository struct {
	db *gorm.DB
}

func NewPotentialRiskStatusUpdateRepository(db *gorm.DB) *PotentialRiskStatusUpdateRepository {
	return &PotentialRiskStatusUpdateRepository{
		db: db,
	}
}

func (repo *PotentialRiskStatusUpdateRepository) GetAll(potentialRiskID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.PotentialRiskStatusUpdateListDTO, error) {

	var potentialRiskUpdates []*riskDTOs.PotentialRiskStatusUpdateListDTO

	query := repo.db.Table("potential_risk_status_updates").
		Select(`
            potential_risk_status_updates.*,
			potential_risks.title AS potential_risk_title,
			potential_risks.description AS potential_risk_description,
			potential_risk_statuses.title_en AS potential_risk_status_title_en,
			potential_risk_statuses.title_ar AS potential_risk_status_title_ar,
			potential_risk_statuses.color_code AS potential_risk_status_color_code,
			created_user.id AS created_user_id,
       		created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar   
        `).
		Joins(`
		    LEFT JOIN potential_risks ON potential_risks.id = potential_risk_status_updates.potential_risk_id
		    LEFT JOIN potential_risk_statuses ON potential_risk_statuses.id = potential_risk_status_updates.potential_risk_status_id
			LEFT JOIN users AS created_user ON created_user.id = potential_risk_status_updates.created_by_id
		`).
		Where("potential_risk_status_updates.potential_risk_id = ?", potentialRiskID)

	if keyword != "" {
		query = query.Where(`
			(potential_risks.title LIKE ? OR
			potential_risk_status_updates.comment LIKE ? OR
			potential_risk_statuses.title_en LIKE ? OR
			potential_risk_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "potential_risk_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&potentialRiskUpdates)

	return potentialRiskUpdates, result.Error
}

func (repo *PotentialRiskStatusUpdateRepository) Create(potentialRiskUpdate *models.PotentialRiskStatusUpdate) (uint, error) {

	if err := repo.db.Create(potentialRiskUpdate).Error; err != nil {
		return 0, err
	}

	return potentialRiskUpdate.ID, nil
}
