package repositories

import (
	"fmt"
	bcDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type JobTeamRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *bcDTOs.JobTeamListQueryParamsDTO) ([]*models.JobTeam, int64, error)
	GetByID(id uint) (*models.JobTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.JobTeam, error)
	Create(jobTeam *models.JobTeam, tx ...utility.DatabaseTransaction) (uint, error)
	Update(jobTeam *models.JobTeam, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateMembers(members []*models.JobTeamMember, tx ...utility.DatabaseTransaction) error
	DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type JobTeamRepository struct {
	db *gorm.DB
}

func NewJobTeamRepository(db *gorm.DB) *JobTeamRepository {
	return &JobTeamRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *JobTeamRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *JobTeamRepository) GetAll(params *bcDTOs.JobTeamListQueryParamsDTO) ([]*models.JobTeam, int64, error) {
	var jobTeams []*models.JobTeam
	var totalCount int64

	// Build the base query with preloading relationships and joins
	query := repo.db.Preload("TeamLeadDesignation").
		Preload("JobTeamMembers").
		Preload("JobTeamMembers.Designation").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Joins("LEFT JOIN designations ON designations.id = job_teams.team_lead_designation_id").
		Model(&models.JobTeam{})

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", columnName), values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (job_teams.title_en LIKE ? OR
            job_teams.title_ar LIKE ? OR
            job_teams.scope LIKE ? OR
            job_teams.objective LIKE ? OR
            designations.title_en LIKE ? OR
            designations.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc" // Default order is descending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	if err := query.Find(&jobTeams).Error; err != nil {
		return nil, 0, err
	}

	return jobTeams, totalCount, nil
}

func (repo *JobTeamRepository) GetByID(id uint) (*models.JobTeam, error) {
	var jobTeam *models.JobTeam
	result := repo.db.First(&jobTeam, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return jobTeam, nil
}

func (repo *JobTeamRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.JobTeam, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var jobTeam models.JobTeam

	// Load the job team with all related data
	result := db.
		Preload("TeamLeadDesignation").
		Preload("JobTeamMembers").
		Preload("JobTeamMembers.Designation").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&jobTeam, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &jobTeam, nil
}

func (repo *JobTeamRepository) Create(jobTeam *models.JobTeam, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(jobTeam).Error; err != nil {
		return 0, err
	}
	return jobTeam.ID, nil
}

func (repo *JobTeamRepository) Update(jobTeam *models.JobTeam, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(jobTeam).Error; err != nil {
		return err
	}
	return nil
}

func (repo *JobTeamRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var jobTeam *models.JobTeam
	if err := db.First(&jobTeam, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.JobTeam{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *JobTeamRepository) CreateMembers(members []*models.JobTeamMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&members).Error; err != nil {
		return err
	}
	return nil
}

func (repo *JobTeamRepository) DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete members by IDs
	if err := db.Where("id IN (?)", teamMembersIDs).Delete(&models.JobTeamMember{}).Error; err != nil {
		return err
	}

	return nil
}
