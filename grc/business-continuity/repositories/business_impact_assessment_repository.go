package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BusinessImpactAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error)
	GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListResponseDTO, error)
	CreateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error
	UpdateBusinessImpactAssessment(businessImpactAssessment *models.BusinessImpactAssessment, tx ...utility.DatabaseTransaction) error
	CreateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error
	UpdateRecoveryObjective(businessImpactAssessmentRecoveryObjective *models.BusinessImpactAssessmentRecoveryObjective, tx ...utility.DatabaseTransaction) error

	DeleteImpactAssessmentsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	DeleteRecoveryObjectivesByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	GetImpactAssessmentsByProcessID(processID uint) ([]*models.BusinessImpactAssessment, error)
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

func (repo *BusinessImpactAssessmentRepository) GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error) {
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
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
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

	response := &businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO{
		Items: businessImpactAssessments,
	}
	return response, nil
}

func (repo *BusinessImpactAssessmentRepository) GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListResponseDTO, error) {
	var businessImpactAssessments []*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("business_impact_assessment_recovery_objectives").
		Select(`
            business_impact_assessment_recovery_objectives.*,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			asset_register.title_en AS asset_register_title_en,
			asset_register.title_ar AS asset_register_title_ar,
			asset_register.rto AS asset_register_rto,
			asset_register.rto_unit_id AS asset_register_rto_unit_id,
			asset_register_rto_time_unit.title_en AS asset_register_rto_unit_title_en,
			asset_register_rto_time_unit.title_ar AS asset_register_rto_unit_title_ar,
			asset_register.mtpd AS asset_register_mtpd,
			asset_register.mtpd_unit_id AS asset_register_mtpd_unit_id,
			asset_register_mtpd_time_unit.title_en AS asset_register_mtpd_unit_title_en,
			asset_register_mtpd_time_unit.title_ar AS asset_register_mtpd_unit_title_ar,
			asset_register.rpo AS asset_register_rpo,
			asset_register.rpo_unit_id AS asset_register_rpo_unit_id,
			asset_register_rpo_time_unit.title_en AS asset_register_rpo_unit_title_en,
			asset_register_rpo_time_unit.title_ar AS asset_register_rpo_unit_title_ar,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			documents.rto AS document_rto,
			documents.rto_unit_id AS document_rto_unit_id,
			document_rto_time_unit.title_en AS document_rto_unit_title_en,
			document_rto_time_unit.title_ar AS document_rto_unit_title_ar,
			documents.mtpd AS document_mtpd,
			documents.mtpd_unit_id AS document_mtpd_unit_id,
			document_mtpd_time_unit.title_en AS document_mtpd_unit_title_en,
			document_mtpd_time_unit.title_ar AS document_mtpd_unit_title_ar,
			documents.rpo AS document_rpo,
			documents.rpo_unit_id AS document_rpo_unit_id,
			document_rpo_time_unit.title_en AS document_rpo_unit_title_en,
			document_rpo_time_unit.title_ar AS document_rpo_unit_title_ar,
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
		    LEFT JOIN processes ON processes.id = business_impact_assessment_recovery_objectives.process_id
		    LEFT JOIN asset_register ON asset_register.id = business_impact_assessment_recovery_objectives.asset_register_id
			LEFT JOIN time_units AS asset_register_rto_time_unit ON asset_register_rto_time_unit.id = asset_register.rto_unit_id
			LEFT JOIN time_units AS asset_register_mtpd_time_unit ON asset_register_mtpd_time_unit.id = asset_register.mtpd_unit_id
			LEFT JOIN time_units AS asset_register_rpo_time_unit ON asset_register_rpo_time_unit.id = asset_register.rpo_unit_id
			LEFT JOIN documents ON documents.id = business_impact_assessment_recovery_objectives.document_id
			LEFT JOIN time_units AS document_rto_time_unit ON document_rto_time_unit.id = documents.rto_unit_id
			LEFT JOIN time_units AS document_mtpd_time_unit ON document_mtpd_time_unit.id = documents.mtpd_unit_id
			LEFT JOIN time_units AS document_rpo_time_unit ON document_rpo_time_unit.id = documents.rpo_unit_id
			LEFT JOIN users AS created_user ON created_user.id = business_impact_assessment_recovery_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_impact_assessment_recovery_objectives.updated_by_id
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

	query = query.Order("business_impact_assessment_recovery_objectives.id")

	err := query.Find(&businessImpactAssessments).Error
	if err != nil {
		return nil, err
	}

	response := &businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListResponseDTO{
		Items: businessImpactAssessments,
	}
	return response, nil
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
