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

type AuditPlanChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) (*auditDTOs.AuditPlanChecklistListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanChecklist, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanChecklist, error)
	Create(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateAuditPlanChecklistMembers(members []*models.AuditPlanChecklistAuditMember, tx ...utility.DatabaseTransaction) error
	DeleteAuditPlanChecklistMembers(memberIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateAuditPlanChecklistAuditableItems(items []*models.AuditPlanChecklistAuditableItem, tx ...utility.DatabaseTransaction) error
	DeleteAuditPlanChecklistAuditableItems(itemIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditPlanChecklistRepository struct {
	db *gorm.DB
}

func NewAuditPlanChecklistRepository(db *gorm.DB) *AuditPlanChecklistRepository {
	return &AuditPlanChecklistRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanChecklistRepository) GetAll(params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) (*auditDTOs.AuditPlanChecklistListResponseDTO, error) {
	var auditPlanChecklists []*auditDTOs.AuditPlanChecklistListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_checklists").
		Select(`
			audit_plan_checklists.*,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			CONVERT(VARCHAR, audit_plan_checklists.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audit_plan_checklists.target_date, 23) AS target_date,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.profile_image_path AS team_lead_profile_image_path,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams_team_lead.first_name_en AS audit_team_team_lead_first_name_en,
			audit_teams_team_lead.first_name_ar AS audit_team_team_lead_first_name_ar,
			audit_teams_team_lead.last_name_en AS audit_team_team_lead_last_name_en,
			audit_teams_team_lead.last_name_ar AS audit_team_team_lead_last_name_ar,
			audit_teams_team_lead.email AS audit_team_team_lead_email,
			audit_teams_team_lead.profile_image_path AS audit_team_team_lead_profile_image_path,
			created_user.first_name_en AS created_user_first_name_en,
        	created_user.first_name_ar AS created_user_first_name_ar,
        	created_user.last_name_en AS created_user_last_name_en,
        	created_user.last_name_ar AS created_user_last_name_ar,   
        	created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,   
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
			LEFT JOIN audit_categories ON audit_categories.id = audit_plan_checklists.audit_category_id
			LEFT JOIN audit_teams ON audit_teams.id = audit_plan_checklists.audit_team_id
        	LEFT JOIN users AS audit_teams_team_lead ON audit_teams_team_lead.id = audit_teams.team_lead_id
			LEFT JOIN audit_statuses ON audit_statuses.id = audit_plan_checklists.audit_status_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plan_checklists.team_lead_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_checklists.created_by_id
        	LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_checklists.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("start_date >= ?", params.StartDate)
	}
	if params.TargetDate != "" {
		query = query.Where("target_date <= ?", params.TargetDate)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_plan_checklists.title LIKE ? OR
			audit_plan_checklists.description LIKE ? OR
			audit_plan_checklists.reference_number LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			audit_teams.title_en LIKE ? OR
			audit_teams.title_ar LIKE ? OR
			audit_statuses.title_en LIKE ? OR
			audit_statuses.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlanChecklists).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanChecklistListResponseDTO{
		Items: auditPlanChecklists,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanChecklistRepository) GetByID(id uint) (*models.AuditPlanChecklist, error) {
	var auditPlanChecklist *models.AuditPlanChecklist

	result := repo.db.First(&auditPlanChecklist, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanChecklist, nil
}

func (repo *AuditPlanChecklistRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanChecklist, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlanChecklist *models.AuditPlanChecklist

	result := db.
		Preload("AuditCategory").
		Preload("AuditTeam").
		Preload("AuditStatus").
		Preload("TeamLead").
		Preload("Members").
		Preload("AuditableItems").
		First(&auditPlanChecklist, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanChecklist, nil
}

func (repo *AuditPlanChecklistRepository) Create(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanChecklist).Error; err != nil {
		return 0, err
	}

	return auditPlanChecklist.ID, nil
}

func (repo *AuditPlanChecklistRepository) Update(auditPlanChecklist *models.AuditPlanChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditPlanChecklist).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanChecklistRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Plan Checklist
	if err := db.Delete(&models.AuditPlanChecklist{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanChecklistRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AuditPlanChecklist{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "APC", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

// Audit Plan Checklist Members
func (repo *AuditPlanChecklistRepository) CreateAuditPlanChecklistMembers(members []*models.AuditPlanChecklistAuditMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&members).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanChecklistRepository) DeleteAuditPlanChecklistMembers(memberIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", memberIds).Delete(&models.AuditPlanChecklistAuditMember{}).Error; err != nil {
		return err
	}

	return nil
}

// Audit Plan Checklist Auditable Items
func (repo *AuditPlanChecklistRepository) CreateAuditPlanChecklistAuditableItems(items []*models.AuditPlanChecklistAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&items).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanChecklistRepository) DeleteAuditPlanChecklistAuditableItems(itemIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", itemIds).Delete(&models.AuditPlanChecklistAuditableItem{}).Error; err != nil {
		return err
	}

	return nil
}
