package repositories

import (
	"fmt"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type OrganizationUnitRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *settingsDTOs.OrganizationUnitListQueryParamsDTO) ([]*settingsDTOs.OrganizationUnitListDTO, error)
	GetByID(id uint) (*models.OrganizationUnit, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.OrganizationUnit, error)
	Create(complianceRegister *models.OrganizationUnit, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceRegister *models.OrganizationUnit, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)

	GetAllDescendantUnitIDs(parentID uint) []uint
}

type OrganizationUnitRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitRepository(db *gorm.DB) *OrganizationUnitRepository {
	return &OrganizationUnitRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *OrganizationUnitRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *OrganizationUnitRepository) GetAll(params *settingsDTOs.OrganizationUnitListQueryParamsDTO) ([]*settingsDTOs.OrganizationUnitListDTO, error) {
	var organizationUnits []*settingsDTOs.OrganizationUnitListDTO

	query := repo.db.Table("organization_units").
		Select(`
            organization_units.id,
			organization_units.organization_level_id,
			organization_levels.level AS organization_level,
            organization_levels.title_en AS organization_level_title_en,
            organization_levels.title_ar AS organization_level_title_ar,
			parent.id AS parent_organization_unit_id,
            parent.title_en AS parent_organization_unit_title_en,
            parent.title_ar AS parent_organization_unit_title_ar,
            organization_units.title_en,
            organization_units.title_ar,
			organization_units.status
        `).
		Joins("LEFT JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id").
		Joins("LEFT JOIN organization_units AS parent ON organization_units.organization_unit_id = parent.id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("organization_units."+columnName+" IN (?)", values)
		}
	}

	// Apply filters
	if params.Status != "" {
		query = query.Where("organization_units.status = ?", params.Status)
	}

	if params.Primary {
		query = query.Where("organization_levels.is_primary = ?", true)
	}

	if params.Level != 0 {
		query = query.Where("organization_units.organization_level_id = ?", params.Level)
	}
	if len(params.ParentUnitIDs) > 0 {
		query = query.Where("organization_units.organization_unit_id IN (?)", params.ParentUnitIDs)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("organization_units.title_en LIKE ? OR organization_units.title_ar LIKE ?", keyword, keyword)
	}

	// Apply sorting
	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&organizationUnits).Error
	if err != nil {
		return nil, err
	}

	return organizationUnits, nil
}

func (repo *OrganizationUnitRepository) GetByID(id uint) (*models.OrganizationUnit, error) {
	var organizationUnit *models.OrganizationUnit

	result := repo.db.First(&organizationUnit, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return organizationUnit, nil

}

func (repo *OrganizationUnitRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.OrganizationUnit, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var organizationUnit *models.OrganizationUnit

	result := db.
		Preload("ComplianceCategory").
		Preload("ComplianceStatus").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ReviewFrequency").
		Preload("ReviewResponsibleUser").
		Preload("SeverityLevel").
		Preload("LifecycleStatus").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&organizationUnit, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return organizationUnit, nil
}

func (repo *OrganizationUnitRepository) Create(complianceRegister *models.OrganizationUnit, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceRegister).Error; err != nil {
		return 0, err
	}

	return complianceRegister.ID, nil
}

func (repo *OrganizationUnitRepository) Update(complianceRegister *models.OrganizationUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *OrganizationUnitRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.OrganizationUnit{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *OrganizationUnitRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.OrganizationUnit{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.OrganizationUnit{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *OrganizationUnitRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.OrganizationUnit{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.OrganizationUnit{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *OrganizationUnitRepository) GetAllDescendantUnitIDs(parentID uint) []uint {
	var allDescendantIDs []uint
	var queue []uint

	queue = append(queue, parentID)

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		var childUnits []models.OrganizationUnit
		err := repo.db.Where("organization_unit_id = ?", currentID).Find(&childUnits).Error
		if err != nil {
			fmt.Printf("Error fetching child units for ID %d: %v\n", currentID, err)
			continue
		}

		for _, child := range childUnits {
			allDescendantIDs = append(allDescendantIDs, child.ID)
			queue = append(queue, child.ID)
		}
	}

	return allDescendantIDs
}
