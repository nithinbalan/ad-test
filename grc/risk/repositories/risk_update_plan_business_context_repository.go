package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanBusinessContextRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error
	GetBusinessContextByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanBusinessContext, error)
	GetBusinessContextByID(id uint) (*models.RiskUpdatePlanBusinessContext, error)
	UpdateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error
	UpdateBusinessContextStatus(id uint, status bool, tx ...utility.DatabaseTransaction) error
	GetBusinessUnitIDByPlanID(planID uint) (uint, error)
	GetOrganizationUnitBusinessContextByBusinessUnit(businessUnitID uint) (*models.OrganizationUnitBusinessContext, error)
	GetBusinessContextIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetProcessIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	GetControlIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
}

type RiskUpdatePlanBusinessContextRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanBusinessContextRepository(db *gorm.DB) *RiskUpdatePlanBusinessContextRepository {
	return &RiskUpdatePlanBusinessContextRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanBusinessContextRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanBusinessContextRepository) CreateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(businessContext).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetBusinessContextByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanBusinessContext, error) {
	var businessContexts []*models.RiskUpdatePlanBusinessContext
	if err := repo.db.
		Preload("RiskUpdatePlanTimeline").
		Preload("CreatedByUser").
		Where("risk_update_plan_timeline_id = ?", timelineID).
		Find(&businessContexts).Error; err != nil {
		return nil, err
	}

	return businessContexts, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetBusinessContextByID(id uint) (*models.RiskUpdatePlanBusinessContext, error) {
	var businessContext models.RiskUpdatePlanBusinessContext
	if err := repo.db.
		Preload("RiskUpdatePlanTimeline").
		Preload("CreatedByUser").
		First(&businessContext, id).Error; err != nil {
		return nil, err
	}

	return &businessContext, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) UpdateBusinessContext(businessContext *models.RiskUpdatePlanBusinessContext, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(businessContext).Error
}

func (repo *RiskUpdatePlanBusinessContextRepository) UpdateBusinessContextStatus(id uint, status bool, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("id = ?", id).
		Update("is_business_context_updated", status).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetBusinessUnitIDByPlanID(planID uint) (uint, error) {
	var businessUnitID uint
	if err := repo.db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Pluck("business_unit_id", &businessUnitID).Error; err != nil {
		return 0, err
	}

	return businessUnitID, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetOrganizationUnitBusinessContextByBusinessUnit(businessUnitID uint) (*models.OrganizationUnitBusinessContext, error) {
	var orgBusinessContext models.OrganizationUnitBusinessContext
	if err := repo.db.Where("department_id = ?", businessUnitID).First(&orgBusinessContext).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &orgBusinessContext, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetBusinessContextIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.OrganizationUnitBusinessContext{}).
		Where("department_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetProcessIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.Process{}).
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *RiskUpdatePlanBusinessContextRepository) GetControlIDsByBusinessUnit(businessUnitID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ids []uint
	if err := db.Model(&models.Control{}).
		Where("responsible_unit_id = ?", businessUnitID).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}
