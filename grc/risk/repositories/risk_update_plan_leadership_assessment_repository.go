package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanLeadershipAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetLeadershipAssessmentUsers(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error)
	CreateLeadershipAssessments(assessments []*models.RiskUpdatePlanLeadershipAssessment, tx ...utility.DatabaseTransaction) error
	DeleteLeadershipAssessmentsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error
	GetLeadershipDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanLeadershipDocument, error)
	CreateLeadershipDocument(document *models.RiskUpdatePlanLeadershipDocument, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanLeadershipDocument, error)
	DeleteLeadershipDocument(id uint, tx ...utility.DatabaseTransaction) error
	GetLeadershipDocumentByID(id uint) (*models.RiskUpdatePlanLeadershipDocument, error)
}

type RiskUpdatePlanLeadershipAssessmentRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanLeadershipAssessmentRepository(db *gorm.DB) *RiskUpdatePlanLeadershipAssessmentRepository {
	return &RiskUpdatePlanLeadershipAssessmentRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) GetLeadershipAssessmentUsers(timelineID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error) {
	var users []*riskDTOs.LeaderShipAssessmentUserListDTO
	var totalCount int64

	query := repo.db.Table("users").
		Select(`
			users.*,
			org1.title_en AS organization_unit_1_title_en,
			org2.title_en AS organization_unit_2_title_en,
			org3.title_en AS organization_unit_3_title_en,
			org4.title_en AS organization_unit_4_title_en,
			org5.title_en AS organization_unit_5_title_en,
			org1.title_ar AS organization_unit_1_title_ar,
			org2.title_ar AS organization_unit_2_title_ar,
			org3.title_ar AS organization_unit_3_title_ar,
			org4.title_ar AS organization_unit_4_title_ar,
			org5.title_ar AS organization_unit_5_title_ar,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			roles_agg.roles_en,
			roles_agg.roles_ar,
			license_types.title_en as license_type_title_en,
			license_types.title_ar as license_type_title_ar,
			license_types.type as license_type_type,
			license_types.color_code as license_type_color_code,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.is_active AS user_status_is_active,
			(SELECT COUNT(*) FROM users u2 WHERE u2.reporting_user_id = users.id) AS reporting_users_count
		`).
		Joins("INNER JOIN risk_update_plan_leadership_assessments ON risk_update_plan_leadership_assessments.user_id = users.id").
		Joins("INNER JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN license_types ON license_types.id = users.license_type_id").
		Joins("LEFT JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Joins("LEFT JOIN (SELECT user_roles.user_id, STRING_AGG(roles.title_en, ', ') AS roles_en, STRING_AGG(roles.title_ar, ', ') AS roles_ar FROM user_roles JOIN roles ON roles.id = user_roles.role_id GROUP BY user_roles.user_id) roles_agg ON roles_agg.user_id = users.id").
		Where("risk_update_plan_leadership_assessments.risk_update_plan_timeline_id = ?", timelineID)

	// Apply search keyword
	if params.Keyword != "" {
		query = query.Where(`
			(users.first_name_en LIKE ? OR
			 users.last_name_en LIKE ? OR
			 users.first_name_ar LIKE ? OR
			 users.last_name_ar LIKE ? OR
			 users.email LIKE ? OR
			 users.mobile_number LIKE ? OR
			 designations.title_en LIKE ? OR
			 designations.title_ar LIKE ? OR
			 user_statuses.title_en LIKE ? OR
			 user_statuses.title_ar LIKE ? OR
			 roles_agg.roles_en LIKE ? OR
			 roles_agg.roles_ar LIKE ?)
		`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute query
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.LeaderShipAssessmentUserListResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) CreateLeadershipAssessments(assessments []*models.RiskUpdatePlanLeadershipAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(assessments) == 0 {
		return nil
	}

	return db.Create(&assessments).Error
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) DeleteLeadershipAssessmentsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).Delete(&models.RiskUpdatePlanLeadershipAssessment{}).Error
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) GetLeadershipDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanLeadershipDocument, error) {
	var documents []*models.RiskUpdatePlanLeadershipDocument

	err := repo.db.
		Where("risk_update_plan_timeline_id = ?", timelineID).
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) CreateLeadershipDocument(document *models.RiskUpdatePlanLeadershipDocument, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanLeadershipDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return nil, err
	}

	return document, nil
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) DeleteLeadershipDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.RiskUpdatePlanLeadershipDocument{}, id).Error
}

func (repo *RiskUpdatePlanLeadershipAssessmentRepository) GetLeadershipDocumentByID(id uint) (*models.RiskUpdatePlanLeadershipDocument, error) {
	var document models.RiskUpdatePlanLeadershipDocument

	err := repo.db.First(&document, id).Error
	if err != nil {
		return nil, err
	}

	return &document, nil
}
