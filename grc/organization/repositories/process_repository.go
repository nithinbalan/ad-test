package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ProcessRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetByID(id uint) (*models.Process, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
	Create(process *models.Process, tx ...utility.DatabaseTransaction) (uint, error)
	Update(process *models.Process, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateProcessControls(controls []*models.ProcessControl, tx ...utility.DatabaseTransaction) error
	DeleteProcessControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error

	//lifecycle status
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)
}

type ProcessRepository struct {
	db *gorm.DB
}

func NewProcessRepository(db *gorm.DB) *ProcessRepository {
	return &ProcessRepository{
		db: db,
	}
}

func (repo *ProcessRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProcessRepository) GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
	var processes []*organizationDTOs.ProcessListDTO
	var totalCount int64

	query := repo.db.Table("processes").
		Select(`
		processes.id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		organization_units.title_en AS owner_unit_title_en,
		organization_units.title_ar AS owner_unit_title_ar,
		owner.first_name_en AS owner_first_name_en,
		owner.first_name_ar AS owner_first_name_ar,
		owner.last_name_en AS owner_last_name_en,
		owner.last_name_ar AS owner_last_name_ar,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		updated_user.first_name_en AS updated_user_first_name_en,
		updated_user.first_name_ar AS updated_user_first_name_ar,
		updated_user.last_name_en AS updated_user_last_name_en,
		updated_user.last_name_ar AS updated_user_last_name_ar,
		processes.created_at,
		processes.updated_at
	`).
		Joins("LEFT JOIN process_groups ON process_groups.id = processes.process_group_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = processes.owner_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = processes.created_by").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	//Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")

		}
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
		(processes.reference_number LIKE ? OR
			processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.description LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&processes).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ProcessListResponseDTO{
		Items: processes,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ProcessRepository) GetByID(id uint) (*models.Process, error) {
	var process *models.Process
	result := repo.db.First(&process, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return process, nil
}

func (repo *ProcessRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
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
		Preload("Controls").
		Preload("CreatedByUser").
		Preload("CreatedByUser.OrganizationUnit1").
		Preload("UpdatedByUser").
		Preload("UpdatedByUser.OrganizationUnit1").
		First(&process, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return process, nil
}

func (repo *ProcessRepository) Create(process *models.Process, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(process).Error; err != nil {
		return 0, err
	}
	return process.ID, nil
}

func (repo *ProcessRepository) Update(process *models.Process, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(process).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.Process{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProcessRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.Process{}).
		Joins("JOIN lifecycle_statuses ON process.lifecycle_status_id = lifecycle_statuses.id").
		Where("process.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("process.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *ProcessRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {
	if err := repo.db.Model(&models.Process{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *ProcessRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Process{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "PR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ProcessRepository) CreateProcessControls(controls []*models.ProcessControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) DeleteProcessControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.ProcessControl{}).Error; err != nil {
		return err
	}

	return nil
}
