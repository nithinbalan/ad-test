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
	var statusCounts []riskDTOs.ControlCountByStatusDTO
	err = repo.db.Table("controls").
		Select(`
			control_implementation_statuses.type AS implementation_status_type,
			control_implementation_statuses.title_en,
			control_implementation_statuses.title_ar,
			control_implementation_statuses.color_code,
			COUNT(controls.id) AS count
		`).
		Joins("INNER JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Where("controls.responsible_unit_id = ?", businessUnitID).
		Where("controls.is_applicable = ?", true).
		Group("control_implementation_statuses.type, control_implementation_statuses.title_en, control_implementation_statuses.title_ar, control_implementation_statuses.color_code").
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

	// Build response
	response := &riskDTOs.ControlRegisterUpdatedDTO{
		NotApplicable: notApplicableCount,
	}

	for _, statusCount := range statusCounts {
		switch statusCount.ImplementationStatusType {
		case "not-implemented":
			response.NotImplemented = &riskDTOs.ControlCountByStatusDTO{
				ImplementationStatusType: statusCount.ImplementationStatusType,
				TitleEn:                  statusCount.TitleEn,
				TitleAr:                  statusCount.TitleAr,
				ColorCode:                statusCount.ColorCode,
				Count:                    statusCount.Count,
			}
		case "partially-implemented":
			response.PartiallyImplemented = &riskDTOs.ControlCountByStatusDTO{
				ImplementationStatusType: statusCount.ImplementationStatusType,
				TitleEn:                  statusCount.TitleEn,
				TitleAr:                  statusCount.TitleAr,
				ColorCode:                statusCount.ColorCode,
				Count:                    statusCount.Count,
			}
		case "implemented":
			response.Implemented = &riskDTOs.ControlCountByStatusDTO{
				ImplementationStatusType: statusCount.ImplementationStatusType,
				TitleEn:                  statusCount.TitleEn,
				TitleAr:                  statusCount.TitleAr,
				ColorCode:                statusCount.ColorCode,
				Count:                    statusCount.Count,
			}
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
