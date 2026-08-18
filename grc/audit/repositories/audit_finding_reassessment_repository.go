package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingReassessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(findingID uint, params *auditDTOs.AuditFindingReassessmentListQueryParamsDTO) (*auditDTOs.AuditFindingReassessmentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFindingReassessment, error)
	GetDetailsByID(id uint) (*models.AuditFindingReassessment, error)
	GetLatestID(findingID uint, tx ...utility.DatabaseTransaction) (uint, error)
	Create(auditFindingReassessment *models.AuditFindingReassessment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFindingReassessment *models.AuditFindingReassessment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	// Document methods
	CreateDocuments(docs []*models.AuditFindingReassessmentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditFindingReassessmentDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetAuditFindingReassessmentActionPlanSummary(reassessmentID uint) ([]auditDTOs.AuditFindingActionPlanStatusCountResponseDTO, error)
}

type AuditFindingReassessmentRepository struct {
	db *gorm.DB
}

func NewAuditFindingReassessmentRepository(db *gorm.DB) AuditFindingReassessmentRepositoryInterface {
	return &AuditFindingReassessmentRepository{
		db: db,
	}
}

func (repo *AuditFindingReassessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingReassessmentRepository) GetAll(findingID uint, params *auditDTOs.AuditFindingReassessmentListQueryParamsDTO) (*auditDTOs.AuditFindingReassessmentListResponseDTO, error) {
	var auditFindingReassessments []*auditDTOs.AuditFindingReassessmentListDTO
	var totalCount int64

	query := repo.db.Table("audit_finding_reassessments").
		Select(`
			audit_finding_reassessments.id,
			audit_finding_reassessments.audit_finding_id,
			audit_finding_reassessments.outcome,
			audit_finding_reassessments.comment,
			audit_finding_reassessments.created_at,
			audit_finding_reassessments.updated_at,
			audit_finding_reassessments.created_by_id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			audit_finding_reassessments.updated_by_id AS updated_user_id,
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
			LEFT JOIN users AS created_user ON created_user.id = audit_finding_reassessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_reassessments.updated_by_id
		`).
		Where("audit_finding_reassessments.audit_finding_id = ?", findingID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_finding_reassessments.comment LIKE ? OR
			audit_finding_reassessments.outcome LIKE ?)`,
			keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_finding_reassessments.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	if err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindingReassessments).Error; err != nil {
		return nil, err
	}

	latestID, err := repo.GetLatestID(findingID)
	if err != nil {
		return nil, err
	}

	for _, reassessment := range auditFindingReassessments {
		var documents []*auditDTOs.AuditFindingReassessmentDocumentDTO
		repo.db.Model(&models.AuditFindingReassessmentDocument{}).
			Where("audit_finding_reassessment_id = ?", reassessment.ID).
			Find(&documents)
		reassessment.Documents = documents
		reassessment.IsLatest = latestID != 0 && reassessment.ID == latestID

		counts, err := repo.GetAuditFindingReassessmentActionPlanSummary(reassessment.ID)
		if err != nil {
			return nil, err
		}
		totalActionPlans := 0
		for _, count := range counts {
			totalActionPlans += count.Count
		}
		reassessment.TotalActionPlans = totalActionPlans
		reassessment.ActionPlanSummary = counts
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFindingReassessmentListResponseDTO{
		Items: auditFindingReassessments,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AuditFindingReassessmentRepository) GetLatestID(findingID uint, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var latestID uint
	err := db.Model(&models.AuditFindingReassessment{}).
		Where("audit_finding_id = ?", findingID).
		Order("created_at DESC, id DESC").
		Limit(1).
		Pluck("id", &latestID).Error
	if err != nil {
		return 0, err
	}

	return latestID, nil
}

func (repo *AuditFindingReassessmentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFindingReassessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditFindingReassessment *models.AuditFindingReassessment
	result := db.Preload("Documents").First(&auditFindingReassessment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingReassessment, nil
}

func (repo *AuditFindingReassessmentRepository) GetDetailsByID(id uint) (*models.AuditFindingReassessment, error) {
	var auditFindingReassessment *models.AuditFindingReassessment

	result := repo.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		First(&auditFindingReassessment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingReassessment, nil
}

func (repo *AuditFindingReassessmentRepository) Create(auditFindingReassessment *models.AuditFindingReassessment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFindingReassessment).Error; err != nil {
		return 0, err
	}

	return auditFindingReassessment.ID, nil
}

func (repo *AuditFindingReassessmentRepository) Update(auditFindingReassessment *models.AuditFindingReassessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFindingReassessment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingReassessmentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditFindingReassessment{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingReassessmentRepository) CreateDocuments(docs []*models.AuditFindingReassessmentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditFindingReassessmentRepository) GetDocumentByID(id uint) (*models.AuditFindingReassessmentDocument, error) {
	var document models.AuditFindingReassessmentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditFindingReassessmentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditFindingReassessmentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingReassessmentRepository) GetAuditFindingReassessmentActionPlanSummary(reassessmentID uint) ([]auditDTOs.AuditFindingActionPlanStatusCountResponseDTO, error) {
	actionPlanStatusCounts := []auditDTOs.AuditFindingActionPlanStatusCountResponseDTO{}

	query := `
        SELECT 
            action_plan_statuses.id, 
            action_plan_statuses.type, 
            action_plan_statuses.title_en, 
            action_plan_statuses.title_ar, 
            action_plan_statuses.color_code,
            COUNT(audit_finding_action_plans.id) AS count
        FROM 
            action_plan_statuses
        LEFT JOIN 
            audit_finding_action_plans 
                ON audit_finding_action_plans.action_plan_status_id = action_plan_statuses.id
                AND audit_finding_action_plans.audit_finding_reassessment_id = ?
        GROUP BY 
            action_plan_statuses.id, 
            action_plan_statuses.type, 
            action_plan_statuses.title_en, 
            action_plan_statuses.title_ar, 
            action_plan_statuses.color_code
    `

	err := repo.db.Raw(query, reassessmentID).Scan(&actionPlanStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return actionPlanStatusCounts, nil
}
