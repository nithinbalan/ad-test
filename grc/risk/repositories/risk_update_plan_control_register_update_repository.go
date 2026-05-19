package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanControlRegisterUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error)
	UpdateControlTestPlanID(controlTestID uint, planID uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanControlRegisterUpdateRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanControlRegisterUpdateRepository(db *gorm.DB) *RiskUpdatePlanControlRegisterUpdateRepository {
	return &RiskUpdatePlanControlRegisterUpdateRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanControlRegisterUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanControlRegisterUpdateRepository) GetControlCountsByImplementationStatus(planID uint) (*riskDTOs.ControlRegisterUpdatedDTO, error) {
	// Get the business unit ID from the plan
	var businessUnitID uint
	err := repo.db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Pluck("business_unit_id", &businessUnitID).Error
	if err != nil {
		return nil, err
	}

	// Get counts grouped by implementation status type for applicable controls
	var statusCounts []struct {
		ImplementationStatusType string
		Count                    int64
	}
	err = repo.db.Table("controls").
		Select(`
			control_implementation_statuses.type AS implementation_status_type,
			COUNT(controls.id) AS count
		`).
		Joins("INNER JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Where("controls.responsible_unit_id = ?", businessUnitID).
		Where("controls.is_applicable = ?", true).
		Group("control_implementation_statuses.type").
		Find(&statusCounts).Error
	if err != nil {
		return nil, err
	}

	// Count not applicable controls
	var notApplicableCount int64
	err = repo.db.Model(&models.Control{}).
		Where("responsible_unit_id = ?", businessUnitID).
		Where("is_applicable = ?", false).
		Count(&notApplicableCount).Error
	if err != nil {
		return nil, err
	}

	// Build response with all statuses initialized to 0
	response := &riskDTOs.ControlRegisterUpdatedDTO{
		NotImplemented:       0,
		PartiallyImplemented: 0,
		Implemented:          0,
		NotApplicable:        notApplicableCount,
	}

	// Fill in actual counts
	for _, statusCount := range statusCounts {
		switch statusCount.ImplementationStatusType {
		case "not-implemented":
			response.NotImplemented = statusCount.Count
		case "partially-implemented":
			response.PartiallyImplemented = statusCount.Count
		case "implemented":
			response.Implemented = statusCount.Count
		}
	}

	return response, nil
}

func (repo *RiskUpdatePlanControlRegisterUpdateRepository) UpdateControlTestPlanID(controlTestID uint, planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Model(&models.ControlTest{}).
		Where("id = ?", controlTestID).
		Update("risk_register_update_plan_id", planID).Error
	if err != nil {
		return err
	}

	return nil
}
