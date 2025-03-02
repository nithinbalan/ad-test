package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpRecoveryPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpRecoveryPlanListQueryParamsDTO) ([]*businessContinuityDTOs.BcpRecoveryPlanListDTO, error)
	GetByID(id uint) (*models.BcpRecoveryPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpRecoveryPlan, error)
	Create(BcpRecoveryPlan *models.BcpRecoveryPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(BcpRecoveryPlan *models.BcpRecoveryPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type BcpRecoveryPlanRepository struct {
	db *gorm.DB
}

func NewBcpRecoveryPlanRepository(db *gorm.DB) *BcpRecoveryPlanRepository {
	return &BcpRecoveryPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BcpRecoveryPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpRecoveryPlanRepository) GetAll(bcpId uint, params *businessContinuityDTOs.BcpRecoveryPlanListQueryParamsDTO) ([]*businessContinuityDTOs.BcpRecoveryPlanListDTO, error) {
	var BcpRecoveryPlans []*businessContinuityDTOs.BcpRecoveryPlanListDTO
	var totalCount int64

	query := repo.db.Table("bcp_recovery_plans").
		Select(`
            bcp_recovery_plans.*,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			bcp_phases.title_en AS bcp_phase_title_en,
			bcp_phases.title_ar AS bcp_phase_title_ar,
			processes.title_en AS bcp_process_title_en,
			processes.title_ar AS bcp_process_title_ar,
			target_completion_unit.title_en AS target_completion_unit_title_en,
			target_completion_unit.title_ar AS target_completion_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path  
        `).
		Joins(`
		    LEFT JOIN severity_levels ON severity_levels.id = bcp_recovery_plans.severity_level_id
			LEFT JOIN bcp_phases ON bcp_phases.id = bcp_recovery_plans.bcp_phase_id
			LEFT JOIN bcp_processes ON bcp_processes.id = bcp_recovery_plans.bcp_process_id
			LEFT JOIN processes ON processes.id = bcp_processes.process_id
			LEFT JOIN time_units AS target_completion_unit ON target_completion_unit.id = bcp_recovery_plans.target_completion_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = bcp_recovery_plans.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = bcp_recovery_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = bcp_recovery_plans.updated_by_id
        `).Where("bcp_recovery_plans.business_continuity_plan_id = ?", bcpId)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(bcp_recovery_plans.title LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			bcp_phases.title_en LIKE ? OR
			bcp_phases.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&BcpRecoveryPlans).Error
	if err != nil {
		return nil, err
	}

	return BcpRecoveryPlans, nil
}

func (repo *BcpRecoveryPlanRepository) GetByID(id uint) (*models.BcpRecoveryPlan, error) {
	var BcpRecoveryPlan *models.BcpRecoveryPlan

	result := repo.db.First(&BcpRecoveryPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return BcpRecoveryPlan, nil

}

func (repo *BcpRecoveryPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpRecoveryPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var BcpRecoveryPlan *models.BcpRecoveryPlan

	result := db.
		Preload("BusinessContinuityPlan").
		Preload("BusinessContinuityPlan.BusinessFunction").
		Preload("SeverityLevel").
		Preload("BcpProcess").
		Preload("BcpProcess.Process").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("BcpPhase").
		Preload("TargetCompletionUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&BcpRecoveryPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return BcpRecoveryPlan, nil
}

func (repo *BcpRecoveryPlanRepository) Create(BcpRecoveryPlan *models.BcpRecoveryPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(BcpRecoveryPlan).Error; err != nil {
		return 0, err
	}

	return BcpRecoveryPlan.ID, nil
}

func (repo *BcpRecoveryPlanRepository) Update(BcpRecoveryPlan *models.BcpRecoveryPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(BcpRecoveryPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpRecoveryPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Recovery Plan
	if err := db.Delete(&models.BcpRecoveryPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}
