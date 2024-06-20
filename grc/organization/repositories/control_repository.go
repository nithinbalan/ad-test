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

type ControlRepositoryInterface interface {
	GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error)
	GetByID(id uint) (*models.Control, error)
	GetFullDetailsByID(id uint) (*models.Control, error)
	Create(control *models.Control) (uint, error)
	Update(complianceRegister *models.Control, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	GenerateReferenceNumber() (string, error)

	//lifecycle status
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)
}

type ControlRepository struct {
	db *gorm.DB
}

func NewControlRepository(db *gorm.DB) *ControlRepository {
	return &ControlRepository{
		db: db,
	}
}

func (repo *ControlRepository) GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error) {
	var controls []*organizationDTOs.ControlListDTO
	var totalCount int64

	query := repo.db.Table("controls").
		Select(`
				controls.*,
				control_types.title_en AS control_type_title_en,
				control_types.title_ar AS control_type_title_ar,
				control_categories.title_en AS control_category_title_en,
				control_categories.title_ar AS control_category_title_ar,
				organization_units.title_en AS owner_unit_title_en,
				organization_units.title_ar AS owner_unit_title_ar,
				owner.first_name_en AS owner_first_name_en,
				owner.first_name_ar AS owner_first_name_ar,
				owner.last_name_en AS owner_last_name_en,
				owner.last_name_ar AS owner_last_name_ar,
				review_frequencies.title_en AS review_frequency_title_en,
				review_frequencies.title_ar AS review_frequency_title_ar,
				test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
				test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
				test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
				test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
				test_frequencies.title_en AS test_frequency_title_en,
				test_frequencies.title_ar AS test_frequency_title_ar,
				lifecycle_statuses.title_en AS lifecycle_status_title_en,
				lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
				lifecycle_statuses.color_code AS lifecycle_status_color_code,
				control_implementation_statuses.title_en AS control_implementation_status_title_en,
				control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
				control_implementation_statuses.color_code AS control_implementation_status_color_code,
				control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
				control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
				control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
				created_user.first_name_en AS created_user_first_name_en,
				created_user.first_name_ar AS created_user_first_name_ar,
				created_user.last_name_en AS created_user_last_name_en,
				created_user.last_name_ar AS created_user_last_name_ar,   
				updated_user.first_name_en AS updated_user_first_name_en,
				updated_user.first_name_ar AS updated_user_first_name_ar,
				updated_user.last_name_en AS updated_user_last_name_en,
				updated_user.last_name_ar AS updated_user_last_name_ar
        `).
		Joins("LEFT JOIN control_types ON control_types.id = controls.control_type_id").
		Joins("LEFT JOIN control_categories ON control_categories.id = controls.control_category_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = controls.owner_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = controls.owner_id").
		Joins("LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = controls.review_frequency_id").
		Joins("LEFT JOIN users AS test_responsible_user ON test_responsible_user.id = controls.test_responsible_user_id").
		Joins("LEFT JOIN review_frequencies AS test_frequencies ON test_frequencies.id = controls.test_frequency_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = controls.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = controls.updated_by_id")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply last_review_date filters
	if params.LastReviewDateFrom != "" {
		query = query.Where("last_review_date >= ?", params.LastReviewDateFrom)
	}
	if params.LastReviewDateTo != "" {
		query = query.Where("last_review_date <= ?", params.LastReviewDateTo)
	}

	// Apply next_review_date filters
	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

	// Apply last_test_date filters
	if params.LastTestDateFrom != "" {
		query = query.Where("last_test_date >= ?", params.LastTestDateFrom)
	}
	if params.LastTestDateTo != "" {
		query = query.Where("last_test_date <= ?", params.LastTestDateTo)
	}

	// Apply next_test_date filters
	if params.NextTestDateFrom != "" {
		query = query.Where("next_test_date >= ?", params.NextTestDateFrom)
	}
	if params.NextTestDateTo != "" {
		query = query.Where("next_test_date <= ?", params.NextTestDateTo)
	}

	// Apply implemented_at filters
	if params.ImplementedAtFrom != "" {
		query = query.Where("implemented_at >= ?", params.ImplementedAtFrom)
	}
	if params.ImplementedAtTo != "" {
		query = query.Where("implemented_at <= ?", params.ImplementedAtTo)
	}

	//Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "implemented":
			query = query.Where("control_implementation_statuses.type = ?", "implemented")
		case "not_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "not_implemented")
		case "partially_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "partially_implemented")
		case "efficient":
			query = query.Where("control_efficiency_statuses.type = ?", "efficient")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")

		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			controls.description LIKE ? OR
			controls.objective LIKE ? OR
			controls.code LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			test_responsible_user.first_name_en LIKE ? OR
			test_responsible_user.first_name_ar LIKE ? OR
			test_responsible_user.last_name_en LIKE ? OR
			test_responsible_user.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
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

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&controls).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ControlListResponseDTO{
		Items: controls,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ControlRepository) GetByID(id uint) (*models.Control, error) {
	var control *models.Control

	result := repo.db.First(&control, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return control, nil

}

func (repo *ControlRepository) GetFullDetailsByID(id uint) (*models.Control, error) {
	var control *models.Control

	latestControlTestSubQuery := repo.db.
		Table("control_tests").
		Select("id").
		Where("control_id = ?", id).
		Order("tested_at desc").
		Limit(1)

	latestControlReviewSubQuery := repo.db.
		Table("control_reviews").
		Select("id").
		Where("control_id = ?", id).
		Order("reviewed_at desc").
		Limit(1)

	result := repo.db.
		Preload("ControlType").
		Preload("ControlCategory").
		Preload("OwnerUnit").
		Preload("Owner").
		Preload("Owner.OrganizationUnit1").
		Preload("ControlImplementationStatus").
		Preload("ControlEfficiencyStatus").
		Preload("TestResponsibleUser").
		Preload("TestResponsibleUser.OrganizationUnit1").
		Preload("ReviewFrequency").
		Preload("ControlReviewStatus").
		Preload("TestFrequency").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("LatestControlTest", func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)", latestControlTestSubQuery).
				Preload("TestedByUser")
		}).
		Preload("LatestControlReview", func(db *gorm.DB) *gorm.DB {
			return db.Where("id IN (?)", latestControlReviewSubQuery).
				Preload("ReviewedByUser").
				Preload("UpdatedByUser").
				Preload("Documents")
		}).
		First(&control, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return control, nil
}

func (repo *ControlRepository) Create(control *models.Control) (uint, error) {
	if err := repo.db.Create(control).Error; err != nil {
		return 0, err
	}

	return control.ID, nil
}

func (repo *ControlRepository) Update(control *models.Control, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(control).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlRepository) Delete(id uint) error {
	if err := repo.db.Delete(&models.Control{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.Control{}).
		Joins("JOIN lifecycle_statuses ON controls.lifecycle_status_id = lifecycle_statuses.id").
		Where("controls.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("controls.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *ControlRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {
	if err := repo.db.Model(&models.Control{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *ControlRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Control{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "CTRL", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
