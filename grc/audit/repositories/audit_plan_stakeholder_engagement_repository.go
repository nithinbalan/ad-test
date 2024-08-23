package repositories

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AuditPlanStakeholderEngagementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanStakeholderEngagementListQueryParamsDTO) (*auditDTOs.AuditPlanStakeholderEngagementListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanStakeholderEngagement, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanStakeholderEngagement, error)
	Create(auditPlanStakeholderEngagement *models.AuditPlanStakeholderEngagement, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanStakeholderEngagement *models.AuditPlanStakeholderEngagement, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type AuditPlanStakeholderEngagementRepository struct {
	db *gorm.DB
}

func NewAuditPlanStakeholderEngagementRepository(db *gorm.DB) *AuditPlanStakeholderEngagementRepository {
	return &AuditPlanStakeholderEngagementRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanStakeholderEngagementRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanStakeholderEngagementRepository) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanStakeholderEngagementListQueryParamsDTO) (*auditDTOs.AuditPlanStakeholderEngagementListResponseDTO, error) {
	var engagements []*auditDTOs.AuditPlanStakeholderEngagementListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_stakeholder_engagements").
		Select(`
			audit_plan_stakeholder_engagements.id,
			audit_plan_stakeholder_engagements.reference_number,
			audit_plan_stakeholder_engagements.audit_engagement_type_id,
			audit_engagement_types.title_en AS audit_engagement_type_title_en,
			audit_engagement_types.title_ar AS audit_engagement_type_title_ar,
			audit_plan_stakeholder_engagements.title,
			audit_plan_stakeholder_engagements.objectives,
			audit_plan_stakeholder_engagements.responsible_user_id,
			audit_plan_stakeholder_engagements.audit_plan_id AS audit_plan_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			audit_plan_stakeholder_engagements.audit_schedule_status_id,
			audit_schedule_statuses.title_en AS audit_schedule_status_title_en,
			audit_schedule_statuses.title_ar AS audit_schedule_status_title_ar,
			audit_schedule_statuses.color_code AS audit_schedule_status_color_code,
			audit_plan_stakeholder_engagements.created_by_id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			audit_plan_stakeholder_engagements.updated_by_id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, audit_plan_stakeholder_engagements.created_at, 23) AS created_at,
			CONVERT(VARCHAR, audit_plan_stakeholder_engagements.updated_at, 23) AS updated_at
		`).
		Joins(`
			LEFT JOIN audit_engagement_types ON audit_engagement_types.id = audit_plan_stakeholder_engagements.audit_engagement_type_id
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_stakeholder_engagements.audit_plan_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = audit_plan_stakeholder_engagements.responsible_user_id
			LEFT JOIN audit_schedule_statuses ON audit_schedule_statuses.id = audit_plan_stakeholder_engagements.audit_schedule_status_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_stakeholder_engagements.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_stakeholder_engagements.updated_by_id
		`).Where(`audit_plan_stakeholder_engagements.audit_plan_id = ?`, auditPlanID)
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply filters and pagination as needed
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			audit_plan_stakeholder_engagements.title LIKE ? OR
			audit_plan_stakeholder_engagements.objectives LIKE ? OR
			audit_engagement_types.title_en LIKE ? OR
			audit_engagement_types.title_ar LIKE ? OR
			audit_schedule_statuses.title_en LIKE ? OR
			audit_schedule_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default sort is ascending
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&engagements).Error
	if err != nil {
		return nil, err
	}
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanStakeholderEngagementListResponseDTO{
		Items: engagements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanStakeholderEngagementRepository) GetByID(id uint) (*models.AuditPlanStakeholderEngagement, error) {
	var engagement *models.AuditPlanStakeholderEngagement

	result := repo.db.First(&engagement, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return engagement, nil
}

func (repo *AuditPlanStakeholderEngagementRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanStakeholderEngagement, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var engagement *models.AuditPlanStakeholderEngagement

	result := db.
		Preload("AuditPlan").
		Preload("AuditEngagementType").
		Preload("ResponsibleUser").
		Preload("AuditScheduleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&engagement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return engagement, nil
}

func (repo *AuditPlanStakeholderEngagementRepository) Create(engagement *models.AuditPlanStakeholderEngagement, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(engagement).Error; err != nil {
		return 0, err
	}

	return engagement.ID, nil
}

func (repo *AuditPlanStakeholderEngagementRepository) Update(engagement *models.AuditPlanStakeholderEngagement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(engagement).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanStakeholderEngagementRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.AuditPlanStakeholderEngagement{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanStakeholderEngagementRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AuditPlanStakeholderEngagement{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "APSE", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
