package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BusinessImpactAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error)
	GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error)
	GetProcessDetails(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
	CreateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error
	UpdateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error
	CreateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error
	UpdateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error

	DeleteImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	DeleteRecoveryObjectivesByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	GetImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error)
	GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error)
}

type BusinessImpactAssessmentRepository struct {
	db *gorm.DB
}

func NewBusinessImpactAssessmentRepository(db *gorm.DB) *BusinessImpactAssessmentRepository {
	return &BusinessImpactAssessmentRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BusinessImpactAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BusinessImpactAssessmentRepository) GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error) {
	var businessImpactAssessments []*businessContinuityDTOs.BusinessImpactAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("business_impact_assessments").
		Select(`
            business_impact_assessments.*,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			asset_register.title_en AS asset_register_title_en,
			asset_register.title_ar AS asset_register_title_ar,
			impact_categories.title_en AS impact_category_title_en,
			impact_categories.title_ar AS impact_category_title_ar,
			impact_categories.type AS impact_category_type,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			process_criticality_levels.id AS process_criticality_level_id,
			process_criticality_levels.title_en AS process_criticality_level_title_en,
			process_criticality_levels.title_ar AS process_criticality_level_title_ar,
			risk_levels.type AS risk_level_type,
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
		    LEFT JOIN processes ON processes.id = business_impact_assessments.process_id
		    LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id
		    LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
		    LEFT JOIN asset_register ON asset_register.id = business_impact_assessments.asset_register_id
			LEFT JOIN impact_categories ON impact_categories.id = business_impact_assessments.impact_category_id
		    LEFT JOIN risk_levels ON risk_levels.id = business_impact_assessments.risk_level_id
			LEFT JOIN users AS created_user ON created_user.id = business_impact_assessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_impact_assessments.updated_by_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "high":
			query = query.Where("process_criticality_levels.type = ?", "high")
		case "medium":
			query = query.Where("process_criticality_levels.type = ?", "medium")
		case "low":
			query = query.Where("process_criticality_levels.type = ?", "low")
		case "all":
			continue
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
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			impact_categories.title_en LIKE ? OR
			impact_categories.title_ar LIKE ? OR
			risk_levels.title_en LIKE ? OR
			risk_levels.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order("business_impact_assessments.id")

	err := query.Find(&businessImpactAssessments).Error
	if err != nil {
		return nil, err
	}

	return businessImpactAssessments, nil
}

func (repo *BusinessImpactAssessmentRepository) GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error) {
	var recoveryObjectives []*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("process_dependencies").
		Select(`
            process_dependencies.id,
			process_dependencies.process_id,
			parent_process.title_en AS process_title_en,
			parent_process.title_ar AS process_title_ar,

			process_dependencies.process_dependency_sub_type_id,
			process_dependency_sub_types.type AS process_dependency_sub_type_type,
			process_dependency_sub_types.title_en AS process_dependency_sub_type_title_en,
			process_dependency_sub_types.title_ar AS process_dependency_sub_type_title_ar,

			process_dependency_types.title_en AS process_dependency_type_title_en,
			process_dependency_types.title_ar AS process_dependency_type_title_ar,
			process_dependency_types.type AS process_dependency_type_type,

			process_dependencies.dependent_process_id,
			dependent_process.title_en AS dependent_process_title_en,
			dependent_process.title_ar AS dependent_process_title_ar,

			process_dependencies.organization_unit_id,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,

			process_dependencies.asset_id,
			asset_register.title_en AS asset_title_en,
			asset_register.title_ar AS asset_title_ar,

			process_dependencies.title,
			process_dependencies.description,
			process_dependencies.rto,
			process_dependencies.rto_unit_id,
			process_dependencies.mtpd,
			process_dependencies.mtpd_unit_id,
			process_dependencies.rpo,
			process_dependencies.rpo_unit_id,

			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
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
		    LEFT JOIN processes AS parent_process ON parent_process.id = process_dependencies.process_id
			LEFT JOIN process_dependency_sub_types ON process_dependency_sub_types.id = process_dependencies.process_dependency_sub_type_id
			LEFT JOIN process_dependency_types ON process_dependency_types.id = process_dependency_sub_types.process_dependency_type_id
			LEFT JOIN processes AS dependent_process ON dependent_process.id = process_dependencies.dependent_process_id
			LEFT JOIN organization_units ON organization_units.id = process_dependencies.organization_unit_id
			LEFT JOIN asset_register ON asset_register.id = process_dependencies.asset_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = process_dependencies.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = process_dependencies.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = process_dependencies.rpo_unit_id
			LEFT JOIN users AS created_user ON created_user.id = process_dependencies.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = process_dependencies.updated_by_id
        `)
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
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order("process_dependencies.id")

	err := query.Find(&recoveryObjectives).Error
	if err != nil {
		return nil, err
	}

	return recoveryObjectives, nil
}

func (repo *BusinessImpactAssessmentRepository) GetProcessDetails(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var process *models.Process

	result := db.
		Preload("ProcessGroup").
		Preload("OwnerUnit").
		Preload("Owner").
		Preload("Owner.OrganizationUnit1").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("CreatedByUser.OrganizationUnit1").
		Preload("UpdatedByUser").
		Preload("UpdatedByUser.OrganizationUnit1").
		Preload("RiskMatrixRating").
		Preload("RiskMatrixRatingValue").
		Preload("Risks").
		Preload("Risks.Risk").
		Preload("Risks.Risk.RiskCategory").
		Preload("Risks.Risk.RiskSource").
		Preload("Risks.Risk.CurrentLikelihood").
		Preload("Risks.Risk.CurrentImpact").
		Preload("Risks.Risk.CurrentRiskRating").
		Preload("Risks.Risk.CurrentRiskRatingValue").
		Preload("Incidents").
		Preload("Incidents.IncidentRegister").
		Preload("Incidents.IncidentRegister.IncidentCategory").
		Preload("Incidents.IncidentRegister.IncidentType").
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("RiskLevel").
		Preload("ProcessCriticalityLevel").
		Preload("LastAssessedByUser").
		Preload("BusinessImpactAssessments").
		Preload("ProcessDependencies").
		Preload("ProcessDependencies.RtoUnit").
		Preload("ProcessDependencies.MtpdUnit").
		Preload("ProcessDependencies.RpoUnit").
		Preload("ProcessDependencies.ProcessDependencySubType").
		Preload("ProcessDependencies.ProcessDependencySubType.ProcessDependencyType").
		First(&process, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return process, nil
}

func (repo *BusinessImpactAssessmentRepository) CreateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessImpactAssessment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessImpactAssessmentRepository) UpdateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(businessImpactAssessment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessImpactAssessmentRepository) CreateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessImpactAssessmentRecoveryObjective).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessImpactAssessmentRepository) UpdateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(businessImpactAssessmentRecoveryObjective).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessImpactAssessmentRepository) DeleteImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("process_id = ?", processID).Delete(&models.BusinessImpactAssessment{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessImpactAssessmentRepository) DeleteRecoveryObjectivesByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("process_id = ?", processID).Delete(&models.BusinessImpactAssessmentRecoveryObjective{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessImpactAssessmentRepository) GetImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error) {
	var businessImpactAssessments []*models.BusinessImpactAssessment
	err := repo.db.Table("business_impact_assessments").
		Select("business_impact_assessments.*").
		Where("process_id = ?", processID).
		Order("id").
		Find(&businessImpactAssessments).Error
	if err != nil {
		return nil, err
	}
	return businessImpactAssessments, nil
}

func (repo *BusinessImpactAssessmentRepository) GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error) {
	var results = []businessContinuityDTOs.BIAMatrixResponseDTO{}

	query := repo.db.Table("business_function_processes").
		Select(`
			business_functions.id AS business_function_id,
			business_functions.title AS business_function_title,
			business_functions.highest_process_criticality_level_id,
			pcl_bf.title_en AS business_function_criticality_title_en,
			pcl_bf.title_ar AS business_function_criticality_title_ar,
			pcl_bf.color_code AS business_function_criticality_color,

			processes.id AS process_id,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			processes.process_criticality_level_id,
			pcl_proc.title_en AS process_criticality_title_en,
			pcl_proc.title_ar AS process_criticality_title_ar,
			pcl_proc.color_code AS process_criticality_color,

			processes.rto AS process_rto,
			processes.mtpd AS process_mtpd,
			processes.rpo AS process_rpo,
			time_units_rto.title_en AS process_rto_unit_title_en,
			time_units_rto.title_ar AS process_rto_unit_title_ar,
			processes.mtpd_unit_id AS process_mtpd_unit_id,
			time_units_mtpd.title_en AS process_mtpd_unit_title_en,
			time_units_mtpd.title_ar AS process_mtpd_unit_title_ar,
			processes.rpo_unit_id AS process_rpo_unit_id,
			time_units_rpo.title_en AS process_rpo_unit_title_en,
			time_units_rpo.title_ar AS process_rpo_unit_title_ar,
			processes.risk_level_id AS risk_level_id,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code
		`).
		Joins(`
			INNER JOIN business_functions ON business_functions.id = business_function_processes.business_function_id
			INNER JOIN processes ON processes.id = business_function_processes.process_id
		    LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id
			LEFT JOIN process_criticality_levels AS pcl_bf ON pcl_bf.id = business_functions.highest_process_criticality_level_id
			LEFT JOIN process_criticality_levels AS pcl_proc ON pcl_proc.id = processes.process_criticality_level_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN time_units AS time_units_rto ON time_units_rto.id = processes.rto_unit_id
			LEFT JOIN time_units AS time_units_mtpd ON time_units_mtpd.id = processes.mtpd_unit_id
			LEFT JOIN time_units AS time_units_rpo ON time_units_rpo.id = processes.rpo_unit_id
		`).
		Where("processes.is_bia_performed = ?", 1)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "high":
			query = query.Where("pcl_proc.type = ?", "high")
		case "medium":
			query = query.Where("pcl_proc.type = ?", "medium")
		case "low":
			query = query.Where("pcl_proc.type = ?", "low")
		case "all":
			continue
		}

	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			business_functions.title LIKE ? OR 
			processes.title_en LIKE ? OR 
			processes.title_ar LIKE ? OR 
			risk_levels.title_en LIKE ? OR 
			risk_levels.title_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
