package repositories

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditTeamRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) (*auditDTOs.AuditTeamListResponseDTO, error)
	GetByID(id uint) (*models.AuditTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditTeam, error)
	Create(auditTeam *models.AuditTeam, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditTeam *models.AuditTeam, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateMembers(users []*models.AuditTeamMember, tx ...utility.DatabaseTransaction) error
	DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditTeamRepository struct {
	db *gorm.DB
}

func NewAuditTeamRepository(db *gorm.DB) *AuditTeamRepository {
	return &AuditTeamRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditTeamRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditTeamRepository) GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) (*auditDTOs.AuditTeamListResponseDTO, error) {

	var auditTeams []*auditDTOs.AuditTeamListDTO
	var totalCount int64

	query := repo.db.Table("audit_teams").
		Select(`
        audit_teams.*,
        users.first_name_en AS team_lead_first_name_en,
        users.first_name_ar AS team_lead_first_name_ar,
        users.last_name_en AS team_lead_last_name_en,
        users.last_name_ar AS team_lead_last_name_ar,
        users.email AS team_lead_email,
        users.profile_image_path AS team_lead_profile_image_path,
        created_user.first_name_en AS created_user_first_name_en,
        created_user.first_name_ar AS created_user_first_name_ar,
        created_user.last_name_en AS created_user_last_name_en,
        created_user.last_name_ar AS created_user_last_name_ar,   
        created_user.email AS created_user_email,
        created_user.profile_image_path AS created_user_profile_image_path,   
        updated_user.first_name_en AS updated_user_first_name_en,
        updated_user.first_name_ar AS updated_user_first_name_ar,
        updated_user.last_name_en AS updated_user_last_name_en,
        updated_user.last_name_ar AS updated_user_last_name_ar,
        updated_user.email AS updated_user_email,
        updated_user.profile_image_path AS updated_user_profile_image_path
    `).
		Joins(`
        LEFT JOIN users ON users.id = audit_teams.team_lead_id
        LEFT JOIN users AS created_user ON created_user.id = audit_teams.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = audit_teams.updated_by_id
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", columnName), values)
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
            (audit_teams.title_en LIKE ? OR
            audit_teams.title_ar LIKE ? OR
            users.first_name_en LIKE ? OR
            users.first_name_ar LIKE ? OR
            users.last_name_en LIKE ? OR
            users.last_name_ar LIKE ? OR
            users.email LIKE ? OR
            users.profile_image_path LIKE ? OR
            created_user.first_name_en LIKE ? OR
            created_user.first_name_ar LIKE ? OR
            created_user.last_name_en LIKE ? OR
            created_user.last_name_ar LIKE ? OR
            created_user.email LIKE ? OR
            created_user.profile_image_path LIKE ? OR
            updated_user.first_name_en LIKE ? OR
            updated_user.first_name_ar LIKE ? OR
            updated_user.last_name_en LIKE ? OR
            updated_user.last_name_ar LIKE ? OR
            updated_user.email LIKE ? OR
            updated_user.profile_image_path LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditTeams).Error
	if err != nil {
		return nil, err
	}
	for _, auditTeam := range auditTeams {
		var teamMembers []*auditDTOs.AuditTeamMembersDTO

		repo.db.Model(&models.AuditTeamMember{}).Where("audit_team_id = ?", auditTeam.ID).Find(&teamMembers)
		auditTeam.Users = teamMembers
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditTeamListResponseDTO{
		Items: auditTeams,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditTeamRepository) GetByID(id uint) (*models.AuditTeam, error) {
	var auditTeam *models.AuditTeam
	result := repo.db.First(&auditTeam, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditTeam, nil
}

func (repo *AuditTeamRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditTeam, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditTeam models.AuditTeam

	// Load the auditable item with all related data
	result := db.
		Preload("TeamLead").
		Preload("Users").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditTeam, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditTeam, nil
}

func (repo *AuditTeamRepository) Create(auditTeam *models.AuditTeam, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditTeam).Error; err != nil {
		return 0, err
	}
	return auditTeam.ID, nil
}

func (repo *AuditTeamRepository) Update(auditTeam *models.AuditTeam, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(auditTeam).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditTeamRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditTeam *models.AuditTeam
	if err := db.First(&auditTeam, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.AuditTeam{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AuditTeamRepository) CreateMembers(users []*models.AuditTeamMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditTeamRepository) DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete users by IDs
	if err := db.Where("id IN (?)", teamMembersIDs).Delete(&models.AuditTeamMember{}).Error; err != nil {
		return err
	}

	return nil
}
