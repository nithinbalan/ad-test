package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"

	"gorm.io/gorm"
)

type ComplianceActionPlanStatusUpdateRepositoryInterface interface {
	GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error)
	Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate) (uint, error)
	Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate) error
}

type ComplianceActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceActionPlanStatusUpdateRepository(db *gorm.DB) *ComplianceActionPlanStatusUpdateRepository {
	return &ComplianceActionPlanStatusUpdateRepository{
		db: db,
	}
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {

	var complianceActionPlanUpdates []*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO

	query := repo.db.Table("compliance_action_plan_status_updates").
		Select(`
            compliance_action_plan_status_updates.*,
            compliance_action_plans.title AS compliance_action_plan_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
            compliance_statuses.color_code AS compliance_status_color_code,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar 
        `).
		Joins(`
            LEFT JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
            LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_action_plan_status_updates.compliance_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_action_plan_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_action_plan_status_updates.updated_by
        `).
		Where("compliance_action_plan_status_updates.compliance_action_plan_id = ?", complianceActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (compliance_action_plans.title LIKE ? OR
            compliance_action_plan_status_updates.comment LIKE ? OR
            compliance_statuses.title_en LIKE ? OR
            compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceActionPlanUpdates)

	return complianceActionPlanUpdates, result.Error
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	var complianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate
	result := repo.db.First(&complianceActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlanStatusUpdate, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate) (uint, error) {

	if err := repo.db.Create(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceActionPlanStatusUpdate.ID, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate) error {
	if err := repo.db.Save(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}
