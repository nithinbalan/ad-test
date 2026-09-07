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
	CreateBusinessImpactAssessmentRating(rating *models.BusinessImpactAssessmentRating, tx ...utility.DatabaseTransaction) error
	GetRatingByID(id uint) (*models.BusinessImpactAssessmentRating, error)
	UpdateBusinessImpactAssessmentRating(rating *models.BusinessImpactAssessmentRating, tx ...utility.DatabaseTransaction) error
	DeleteBusinessImpactAssessmentRatingByProcessID(processID uint, tx ...utility.DatabaseTransaction) error

	DeleteImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	DeleteRecoveryObjectivesByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	GetImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error)
	GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error)

	// New methods for BusinessImpactAssessmentAnalysis
	CreateBusinessImpactAssessmentAnalysis(analysis []*models.BusinessImpactAssessmentAnalysis, tx ...utility.DatabaseTransaction) error
	DeleteBusinessImpactAssessmentAnalysisByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error
	DeleteBusinessImpactAssessmentAnalysisByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	GetImpactAssessmentsWithAnalysisByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error)
	GetCriticalSystems() ([]*businessContinuityDTOs.CriticalSystemListDTO, error)
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
            updated_user.profile_image_path AS updated_user_profile_image_path,
			bia_processes.workflow_status_id AS workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.type AS workflow_status_type,
			workflow_statuses.color_code AS workflow_status_color_code,
			bia_processes.approved_by_id AS approved_by_id,
			approved_by.first_name_en AS approved_by_first_name_en,
			approved_by.first_name_ar AS approved_by_first_name_ar,
			approved_by.last_name_en AS approved_by_last_name_en,
			approved_by.last_name_ar AS approved_by_last_name_ar,
			approved_by.initial AS approved_by_initial,
			approved_by.color AS approved_by_color,
			approved_by.email AS approved_by_email,
			approved_by.profile_image_path AS approved_by_profile_image_path,
			CONVERT(VARCHAR, bia_processes.approved_at, 23) AS approved_at
        `).
		Joins(`
		    LEFT JOIN processes ON processes.id = business_impact_assessments.process_id
			LEFT JOIN bia_processes ON bia_processes.process_id = processes.id
			LEFT JOIN workflow_statuses ON workflow_statuses.id = bia_processes.workflow_status_id
			LEFT JOIN users AS approved_by ON approved_by.id = bia_processes.approved_by_id
		    LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
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
			switch columnName {
			case "workflow_status_ids":
				query = query.Where("bia_processes.workflow_status_id IN (?)", values)
			case "process_id", "process_ids":
				query = query.Where("business_impact_assessments.process_id IN (?)", values)
			default:
				query = query.Where(columnName+" IN (?)", values)
			}
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
			workflow_statuses.title_en LIKE ? OR
			workflow_statuses.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			approved_by.first_name_en LIKE ? OR
			approved_by.first_name_ar LIKE ? OR
			approved_by.last_name_en LIKE ? OR
			approved_by.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	if params.IsCyberUser {
		query = query.Where("business_impact_assessments.is_cyber = ?", true)
	} else {
		query = query.Where("business_impact_assessments.is_cyber = ?", false)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "business_impact_assessments.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

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

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "process_dependencies.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

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
		Preload("ResponsibleUnit").
		Preload("Section").
		Preload("SubSection").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
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
		Preload("BusinessImpactAssessments.ImpactCategory").
		Preload("BusinessImpactAssessments.BusinessImpactAssessmentAnalysis").
		Preload("BusinessImpactAssessments.BusinessImpactAssessmentAnalysis.OutageBucket").
		Preload("BusinessImpactAssessments.BusinessImpactAssessmentAnalysis.OutageBucket.DurationUnit").
		Preload("BusinessImpactAssessments.BusinessImpactAssessmentAnalysis.ImpactAnalysisLevel").
		Preload("BusinessImpactAssessmentRatings").
		Preload("BusinessImpactAssessmentRatings.OverallImpactLevel").
		Preload("BusinessImpactAssessmentRatings.OverallOutageBucket").
		Preload("BusinessImpactAssessmentRatings.BCMBIAImpactLevel").
		Preload("BusinessImpactAssessmentRatings.BCMBIAImpactLevel.RtoUnit").
		Preload("BusinessImpactAssessmentRatings.BCMBIAImpactLevel.MtpdUnit").
		Preload("BusinessImpactAssessmentRatings.BCMBIAImpactLevel.RpoUnit").
		Preload("BusinessImpactAssessmentRatings.BCMBIAImpactLevel.RspUnit").
		Preload("BusinessImpactAssessmentRatings.RtoUnit").
		Preload("BusinessImpactAssessmentRatings.MtpdUnit").
		Preload("BusinessImpactAssessmentRatings.RpoUnit").
		Preload("BusinessImpactAssessmentRatings.RspUnit").
		Preload("ProcessDependencies").
		Preload("ProcessDependencies.RtoUnit").
		Preload("ProcessDependencies.MtpdUnit").
		Preload("ProcessDependencies.RpoUnit").
		Preload("ProcessDependencies.ProcessDependencySubType").
		Preload("ProcessDependencies.ProcessDependencySubType.ProcessDependencyType").
		Preload("ProcessEmployeeAnalysis").
		Preload("ProcessEmployeeAnalysis.Designation").
		Preload("ProcessEmployeeAnalysis.Employees").
		Preload("ProcessEmployeeAnalysis.Employees.User").
		Preload("BCMStaffAllocations").
		Preload("BCMStaffAllocations.ProcessEmployeeAnalysis").
		Preload("BCMStaffAllocations.ProcessEmployeeAnalysis.Designation").
		Preload("BCMStaffAllocations.ProcessEmployeeAnalysis.Employees").
		Preload("BCMStaffAllocations.ProcessEmployeeAnalysis.Employees.User").
		Preload("BCMStaffAllocations.OutageBuckets").
		Preload("BCMStaffAllocations.OutageBuckets.OutageBucket").
		Preload("BCMStaffAllocations.OutageBuckets.Users").
		Preload("BCMStaffAllocations.OutageBuckets.Users.User").
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

			processes.impact_analysis_level_id AS impact_analysis_level_id,
			impact_levels.title_en AS impact_analysis_level_title_en,
			impact_levels.title_ar AS impact_analysis_level_title_ar,
			impact_levels.color_code AS impact_analysis_level_color_code,
			impact_levels.type AS impact_analysis_level_type,

			processes.rto_from_value AS process_rto_from_value,
			processes.rto_to_value AS process_rto_to_value,
			processes.rto_condition AS process_rto_condition,
			processes.rto_unit_id AS process_rto_unit_id,
			time_units_rto.title_en AS process_rto_unit_title_en,
			time_units_rto.title_ar AS process_rto_unit_title_ar,
			processes.mtpd_from_value AS process_mtpd_from_value,
			processes.mtpd_to_value AS process_mtpd_to_value,
			processes.mtpd_condition AS process_mtpd_condition,
			processes.mtpd_unit_id AS process_mtpd_unit_id,
			time_units_mtpd.title_en AS process_mtpd_unit_title_en,
			time_units_mtpd.title_ar AS process_mtpd_unit_title_ar,
			processes.rpo_from_value AS process_rpo_from_value,
			processes.rpo_to_value AS process_rpo_to_value,
			processes.rpo_condition AS process_rpo_condition,
			processes.rpo_unit_id AS process_rpo_unit_id,
			time_units_rpo.title_en AS process_rpo_unit_title_en,
			time_units_rpo.title_ar AS process_rpo_unit_title_ar,
			processes.rsp_from_value AS process_rsp_from_value,
			processes.rsp_to_value AS process_rsp_to_value,
			processes.rsp_condition AS process_rsp_condition,
			processes.rsp_unit_id AS process_rsp_unit_id,
			time_units_rsp.title_en AS process_rsp_unit_title_en,
			time_units_rsp.title_ar AS process_rsp_unit_title_ar,

			processes.risk_level_id AS risk_level_id,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code
		`).
		Joins(`
			INNER JOIN business_functions ON business_functions.id = business_function_processes.business_function_id
			INNER JOIN processes ON processes.id = business_function_processes.process_id
		    LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN process_criticality_levels AS pcl_bf ON pcl_bf.id = business_functions.highest_process_criticality_level_id
			LEFT JOIN process_criticality_levels AS pcl_proc ON pcl_proc.id = processes.process_criticality_level_id
			LEFT JOIN bcm_impact_analysis_levels AS impact_levels ON impact_levels.id = processes.impact_analysis_level_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN business_continuity_duration_units AS time_units_rto ON time_units_rto.id = processes.rto_unit_id
			LEFT JOIN business_continuity_duration_units AS time_units_mtpd ON time_units_mtpd.id = processes.mtpd_unit_id
			LEFT JOIN business_continuity_duration_units AS time_units_rpo ON time_units_rpo.id = processes.rpo_unit_id
			LEFT JOIN business_continuity_duration_units AS time_units_rsp ON time_units_rsp.id = processes.rsp_unit_id
		`).
		Where("processes.is_bia_performed = ?", 1).
		Where("business_functions.is_cyber = ?", params.IsCyberUser)

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

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "processes.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *BusinessImpactAssessmentRepository) CreateBusinessImpactAssessmentRating(rating *models.BusinessImpactAssessmentRating, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(rating).Error; err != nil {
		return err
	}

	return nil
}

// New methods for BusinessImpactAssessmentAnalysis
func (repo *BusinessImpactAssessmentRepository) CreateBusinessImpactAssessmentAnalysis(analysis []*models.BusinessImpactAssessmentAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(analysis) == 0 {
		return nil
	}

	return db.Create(&analysis).Error
}

func (repo *BusinessImpactAssessmentRepository) DeleteBusinessImpactAssessmentAnalysisByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("business_impact_assessment_id = ?", assessmentID).Delete(&models.BusinessImpactAssessmentAnalysis{}).Error
}

func (repo *BusinessImpactAssessmentRepository) DeleteBusinessImpactAssessmentAnalysisByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("business_impact_assessment_id IN (SELECT id FROM business_impact_assessments WHERE process_id = ?)", processID).Delete(&models.BusinessImpactAssessmentAnalysis{}).Error
}

func (repo *BusinessImpactAssessmentRepository) GetImpactAssessmentsWithAnalysisByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error) {
	var businessImpactAssessments []*models.BusinessImpactAssessment
	err := repo.db.Table("business_impact_assessments").
		Preload("BusinessImpactAssessmentAnalysis.OutageBucket").
		Preload("BusinessImpactAssessmentAnalysis.ImpactAnalysisLevel").
		Preload("ImpactCategory").
		Where("process_id = ?", processID).
		Order("id").
		Find(&businessImpactAssessments).Error
	if err != nil {
		return nil, err
	}
	return businessImpactAssessments, nil
}

func (repo *BusinessImpactAssessmentRepository) GetRatingByID(id uint) (*models.BusinessImpactAssessmentRating, error) {
	var rating *models.BusinessImpactAssessmentRating
	err := repo.db.First(&rating, id).Error
	return rating, err
}

func (repo *BusinessImpactAssessmentRepository) UpdateBusinessImpactAssessmentRating(rating *models.BusinessImpactAssessmentRating, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.BusinessImpactAssessmentRating{}).Where("id = ?", rating.ID).
		Select("RtoFromValue", "IsEdited", "OverrideReason", "UpdatedBYID").
		Updates(rating).Error
}

func (repo *BusinessImpactAssessmentRepository) DeleteBusinessImpactAssessmentRatingByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("process_id = ?", processID).Delete(&models.BusinessImpactAssessmentRating{}).Error
}

func (repo *BusinessImpactAssessmentRepository) GetCriticalSystems() ([]*businessContinuityDTOs.CriticalSystemListDTO, error) {
	type systemRow struct {
		SystemRegisterID      uint    `gorm:"column:system_register_id"`
		SystemRegisterTitleEn string  `gorm:"column:system_register_title_en"`
		SystemRegisterTitleAr string  `gorm:"column:system_register_title_ar"`
		ServiceDescription    *string `gorm:"column:service_description"`
		RTO                   *int    `gorm:"column:rto"`
		RTOUnitTitleEn        *string `gorm:"column:rto_unit_title_en"`
		RTOUnitTitleAr        *string `gorm:"column:rto_unit_title_ar"`
		MTPD                  *int    `gorm:"column:mtpd"`
		MTPDUnitTitleEn       *string `gorm:"column:mtpd_unit_title_en"`
		MTPDUnitTitleAr       *string `gorm:"column:mtpd_unit_title_ar"`
		RPO                   *int    `gorm:"column:rpo"`
		RPOUnitTitleEn        *string `gorm:"column:rpo_unit_title_en"`
		RPOUnitTitleAr        *string `gorm:"column:rpo_unit_title_ar"`
		RecoveryManagerPhone  string  `gorm:"column:recovery_manager_phone"`
		RecoveryManagerEmail  string  `gorm:"column:recovery_manager_email"`
		RecoveryLocation      *string `gorm:"column:recovery_location"`
		VendorRTO             *int    `gorm:"column:vendor_rto"`
		VendorRTOUnitTitleEn  *string `gorm:"column:vendor_rto_unit_title_en"`
		VendorRTOUnitTitleAr  *string `gorm:"column:vendor_rto_unit_title_ar"`
	}

	type processRow struct {
		SystemRegisterID             uint    `gorm:"column:system_register_id"`
		ProcessID                    uint    `gorm:"column:process_id"`
		ProcessTitleEn               string  `gorm:"column:process_title_en"`
		ProcessTitleAr               string  `gorm:"column:process_title_ar"`
		ImpactAnalysisLevelID        *uint   `gorm:"column:impact_analysis_level_id"`
		ImpactAnalysisLevelTitleEn   *string `gorm:"column:impact_analysis_level_title_en"`
		ImpactAnalysisLevelTitleAr   *string `gorm:"column:impact_analysis_level_title_ar"`
		ImpactAnalysisLevelColorCode *string `gorm:"column:impact_analysis_level_color_code"`
		RtoFromValue                 *int    `gorm:"column:rto_from_value"`
		RtoToValue                   *int    `gorm:"column:rto_to_value"`
		RtoCondition                 *string `gorm:"column:rto_condition"`
		RtoUnitTitleEn               *string `gorm:"column:rto_unit_title_en"`
		RtoUnitTitleAr               *string `gorm:"column:rto_unit_title_ar"`
		MtpdFromValue                *int    `gorm:"column:mtpd_from_value"`
		MtpdToValue                  *int    `gorm:"column:mtpd_to_value"`
		MtpdCondition                *string `gorm:"column:mtpd_condition"`
		MtpdUnitTitleEn              *string `gorm:"column:mtpd_unit_title_en"`
		MtpdUnitTitleAr              *string `gorm:"column:mtpd_unit_title_ar"`
		RpoFromValue                 *int    `gorm:"column:rpo_from_value"`
		RpoToValue                   *int    `gorm:"column:rpo_to_value"`
		RpoCondition                 *string `gorm:"column:rpo_condition"`
		RpoUnitTitleEn               *string `gorm:"column:rpo_unit_title_en"`
		RpoUnitTitleAr               *string `gorm:"column:rpo_unit_title_ar"`
	}

	var systemRows []systemRow
	err := repo.db.Table("system_registers sr").
		Select(`
			sr.id                       AS system_register_id,
			sr.title_en                 AS system_register_title_en,
			sr.title_ar                 AS system_register_title_ar,
			sr.service_description      AS service_description,
			sr.rto                      AS rto,
			rto_unit.title_en           AS rto_unit_title_en,
			rto_unit.title_ar           AS rto_unit_title_ar,
			sr.mtpd                     AS mtpd,
			mtpd_unit.title_en          AS mtpd_unit_title_en,
			mtpd_unit.title_ar          AS mtpd_unit_title_ar,
			sr.rpo                      AS rpo,
			rpo_unit.title_en           AS rpo_unit_title_en,
			rpo_unit.title_ar           AS rpo_unit_title_ar,
			sr.recovery_manager_phone   AS recovery_manager_phone,
			sr.recovery_manager_email   AS recovery_manager_email,
			sr.recovery_location        AS recovery_location,
			sr.vendor_rto               AS vendor_rto,
			vrto_unit.title_en          AS vendor_rto_unit_title_en,
			vrto_unit.title_ar          AS vendor_rto_unit_title_ar
		`).
		Joins(`
			INNER JOIN process_system_registers psr ON psr.system_register_id = sr.id
			INNER JOIN processes p ON p.id = psr.process_id
			INNER JOIN bcm_impact_analysis_levels ial ON ial.id = p.impact_analysis_level_id
			LEFT  JOIN business_continuity_duration_units rto_unit  ON rto_unit.id  = sr.rto_unit_id
			LEFT  JOIN business_continuity_duration_units mtpd_unit ON mtpd_unit.id = sr.mtpd_unit_id
			LEFT  JOIN business_continuity_duration_units rpo_unit  ON rpo_unit.id  = sr.rpo_unit_id
			LEFT  JOIN business_continuity_duration_units vrto_unit ON vrto_unit.id = sr.vendor_rto_unit_id
		`).
		Where("p.impact_analysis_level_id IS NOT NULL AND ial.type != ?", "non_applicable").
		Group("sr.id, sr.title_en, sr.title_ar, sr.service_description, sr.rto, rto_unit.title_en, rto_unit.title_ar, sr.mtpd, mtpd_unit.title_en, mtpd_unit.title_ar, sr.rpo, rpo_unit.title_en, rpo_unit.title_ar, sr.recovery_manager_phone, sr.recovery_manager_email, sr.recovery_location, sr.vendor_rto, vrto_unit.title_en, vrto_unit.title_ar").
		Order("sr.id").
		Find(&systemRows).Error
	if err != nil {
		return nil, err
	}

	if len(systemRows) == 0 {
		return []*businessContinuityDTOs.CriticalSystemListDTO{}, nil
	}

	// Collect system IDs
	systemIDs := make([]uint, 0, len(systemRows))
	for _, s := range systemRows {
		systemIDs = append(systemIDs, s.SystemRegisterID)
	}

	// Fetch critical processes grouped by system register
	var processRows []processRow
	err = repo.db.Table("processes p").
		Select(`
			psr.system_register_id              AS system_register_id,
			p.id                                AS process_id,
			p.title_en                          AS process_title_en,
			p.title_ar                          AS process_title_ar,
			p.impact_analysis_level_id          AS impact_analysis_level_id,
			ial.title_en                        AS impact_analysis_level_title_en,
			ial.title_ar                        AS impact_analysis_level_title_ar,
			ial.color_code                      AS impact_analysis_level_color_code,
			p.rto_from_value                    AS rto_from_value,
			p.rto_to_value                      AS rto_to_value,
			p.rto_condition                     AS rto_condition,
			rto_unit.title_en                   AS rto_unit_title_en,
			rto_unit.title_ar                   AS rto_unit_title_ar,
			p.mtpd_from_value                   AS mtpd_from_value,
			p.mtpd_to_value                     AS mtpd_to_value,
			p.mtpd_condition                    AS mtpd_condition,
			mtpd_unit.title_en                  AS mtpd_unit_title_en,
			mtpd_unit.title_ar                  AS mtpd_unit_title_ar,
			p.rpo_from_value                    AS rpo_from_value,
			p.rpo_to_value                      AS rpo_to_value,
			p.rpo_condition                     AS rpo_condition,
			rpo_unit.title_en                   AS rpo_unit_title_en,
			rpo_unit.title_ar                   AS rpo_unit_title_ar
		`).
		Joins(`
			INNER JOIN process_system_registers psr ON psr.process_id = p.id
			INNER JOIN bcm_impact_analysis_levels ial ON ial.id = p.impact_analysis_level_id
			LEFT  JOIN business_continuity_duration_units rto_unit  ON rto_unit.id  = p.rto_unit_id
			LEFT  JOIN business_continuity_duration_units mtpd_unit ON mtpd_unit.id = p.mtpd_unit_id
			LEFT  JOIN business_continuity_duration_units rpo_unit  ON rpo_unit.id  = p.rpo_unit_id
		`).
		Where("p.impact_analysis_level_id IS NOT NULL AND ial.type != ? AND psr.system_register_id IN (?)", "non_applicable", systemIDs).
		Order("psr.system_register_id, p.id").
		Find(&processRows).Error
	if err != nil {
		return nil, err
	}

	// Build process map keyed by system_register_id
	processMap := make(map[uint][]businessContinuityDTOs.CriticalSystemProcessDTO)
	for _, pr := range processRows {
		processMap[pr.SystemRegisterID] = append(processMap[pr.SystemRegisterID], businessContinuityDTOs.CriticalSystemProcessDTO{
			ProcessID:                    pr.ProcessID,
			ProcessTitleEn:               pr.ProcessTitleEn,
			ProcessTitleAr:               pr.ProcessTitleAr,
			ImpactAnalysisLevelID:        pr.ImpactAnalysisLevelID,
			ImpactAnalysisLevelTitleEn:   pr.ImpactAnalysisLevelTitleEn,
			ImpactAnalysisLevelTitleAr:   pr.ImpactAnalysisLevelTitleAr,
			ImpactAnalysisLevelColorCode: pr.ImpactAnalysisLevelColorCode,
			RtoFromValue:                 pr.RtoFromValue,
			RtoToValue:                   pr.RtoToValue,
			RtoCondition:                 pr.RtoCondition,
			RtoUnitTitleEn:               pr.RtoUnitTitleEn,
			RtoUnitTitleAr:               pr.RtoUnitTitleAr,
			MtpdFromValue:                pr.MtpdFromValue,
			MtpdToValue:                  pr.MtpdToValue,
			MtpdCondition:                pr.MtpdCondition,
			MtpdUnitTitleEn:              pr.MtpdUnitTitleEn,
			MtpdUnitTitleAr:              pr.MtpdUnitTitleAr,
			RpoFromValue:                 pr.RpoFromValue,
			RpoToValue:                   pr.RpoToValue,
			RpoCondition:                 pr.RpoCondition,
			RpoUnitTitleEn:               pr.RpoUnitTitleEn,
			RpoUnitTitleAr:               pr.RpoUnitTitleAr,
		})
	}

	result := make([]*businessContinuityDTOs.CriticalSystemListDTO, 0, len(systemRows))
	for _, s := range systemRows {
		processes := processMap[s.SystemRegisterID]
		if processes == nil {
			processes = []businessContinuityDTOs.CriticalSystemProcessDTO{}
		}
		result = append(result, &businessContinuityDTOs.CriticalSystemListDTO{
			SystemRegisterID:      s.SystemRegisterID,
			SystemRegisterTitleEn: s.SystemRegisterTitleEn,
			SystemRegisterTitleAr: s.SystemRegisterTitleAr,
			ServiceDescription:    s.ServiceDescription,
			RTO:                   s.RTO,
			RTOUnitTitleEn:        s.RTOUnitTitleEn,
			RTOUnitTitleAr:        s.RTOUnitTitleAr,
			MTPD:                  s.MTPD,
			MTPDUnitTitleEn:       s.MTPDUnitTitleEn,
			MTPDUnitTitleAr:       s.MTPDUnitTitleAr,
			RPO:                   s.RPO,
			RPOUnitTitleEn:        s.RPOUnitTitleEn,
			RPOUnitTitleAr:        s.RPOUnitTitleAr,
			RecoveryManagerPhone:  s.RecoveryManagerPhone,
			RecoveryManagerEmail:  s.RecoveryManagerEmail,
			RecoveryLocation:      s.RecoveryLocation,
			VendorRTO:             s.VendorRTO,
			VendorRTOUnitTitleEn:  s.VendorRTOUnitTitleEn,
			VendorRTOUnitTitleAr:  s.VendorRTOUnitTitleAr,
			CriticalProcessCount:  len(processes),
			CriticalProcesses:     processes,
		})
	}

	return result, nil
}
