package repositories

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditTeamRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) ([]*models.AuditTeam, int64, error)
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

func (repo *AuditTeamRepository) GetAll(params *auditDTOs.AuditTeamListQueryParamsDTO) ([]*models.AuditTeam, int64, error) {
	var auditTeams []*models.AuditTeam
	var totalCount int64

	// Build the base query with preloading relationships
	query := repo.db.Preload("TeamLead").
		Preload("Users").
		Preload("Users.User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Model(&models.AuditTeam{})

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", columnName), values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (audit_teams.title_en LIKE ? OR
            audit_teams.title_ar LIKE ?)`,
			keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	if err := query.Find(&auditTeams).Error; err != nil {
		return nil, 0, err
	}

	return auditTeams, totalCount, nil
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
		Preload("Users.User").
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
