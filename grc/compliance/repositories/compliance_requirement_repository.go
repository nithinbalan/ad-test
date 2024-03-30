package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	masterModels "jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRequirementRepositoryInterface interface {
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*complianceModels.ComplianceRequirement, error)
	Create(complianceClause *complianceModels.ComplianceClause, complianceRequirement *complianceModels.ComplianceRequirement) (uint, error)
	Update(complianceRequirement *complianceModels.ComplianceRequirement) error
	Delete(id uint) error
	IsValidID(id uint) (*complianceModels.ComplianceRequirement, error)
}

type ComplianceRequirementRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementRepository(db *gorm.DB) *ComplianceRequirementRepository {
	return &ComplianceRequirementRepository{
		db: db,
	}
}

func (repo *ComplianceRequirementRepository) GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {

	var complianceRequirements []*complianceDTOs.ComplianceRequirementListDTO
	var totalCount int64

	query := repo.db.Table("compliance_register").
		Select(`
            compliance_register.*,
			compliance_categories.title_en AS compliance_category_title_en,
			compliance_categories.title_ar AS compliance_category_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_responsible_user.first_name_en AS review_responsible_user_first_name_en,
            review_responsible_user.first_name_ar AS review_responsible_user_first_name_ar,
            review_responsible_user.last_name_en AS review_responsible_user_last_name_en,
            review_responsible_user.last_name_ar AS review_responsible_user_last_name_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
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
            updated_user.last_name_ar AS updated_user_last_name_ar          
        `).
		Joins(`
		    LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
		    LEFT JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
		    LEFT JOIN review_frequencies ON review_frequencies.id = compliance_register.review_frequency_id
			LEFT JOIN users AS review_responsible_user ON review_responsible_user.id = compliance_register.review_responsible_user_id
		    LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = compliance_register.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_register.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_register.updated_by
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply deadline filters
	if params.DeadlineFrom != "" {
		query = query.Where("deadline >= ?", params.DeadlineFrom)
	}
	if params.DeadlineTo != "" {
		query = query.Where("deadline <= ?", params.DeadlineTo)
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "compliant":
			query = query.Where("compliance_statuses.type = ?", "compliant")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not_implemented")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")

		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
			compliance_register.description_en LIKE ? OR
			compliance_register.description_ar LIKE ? OR
			compliance_register.standard_code LIKE ? OR
			compliance_categories.title_en LIKE ? OR
			compliance_categories.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceRequirements).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceRequirementListResponseDTO{
		Items: complianceRequirements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceRequirementRepository) GetByID(complianceID uint, id uint) (*complianceModels.ComplianceRequirement, error) {
	var complianceRequirement *complianceModels.ComplianceRequirement
	result := repo.db.First(&complianceRequirement, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirement, result.Error
}

func (repo *ComplianceRequirementRepository) Create(complianceClause *complianceModels.ComplianceClause, complianceRequirement *complianceModels.ComplianceRequirement) (uint, error) {

	var existingClause complianceModels.ComplianceClause
	err := repo.db.Where("clause_no = ? AND title = ?", complianceClause.ClauseNo, complianceClause.Title).First(&existingClause).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := repo.db.Create(complianceClause).Error; err != nil {
				return 0, err
			}
			complianceRequirement.ComplianceClauseID = complianceClause.ID
		} else {
			return 0, err
		}
	} else {
		// If the clause exists, use its ID for the compliance requirement
		complianceRequirement.ComplianceClauseID = existingClause.ID
	}

	if err := repo.db.Create(complianceRequirement).Error; err != nil {
		return 0, err
	}

	return complianceRequirement.ID, nil
}

func (repo *ComplianceRequirementRepository) Update(complianceRequirement *complianceModels.ComplianceRequirement) error {

	if err := repo.db.Save(complianceRequirement).Error; err != nil {

		return err
	}

	return nil
}

func (repo *ComplianceRequirementRepository) Delete(id uint) error {

	var compliance *complianceModels.ComplianceRequirement
	if err := repo.db.First(&compliance, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&complianceModels.ComplianceRequirement{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceRequirementRepository) Archive(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Retrieve IDs of ComplianceRequirements with lifecycle status "active"
	if err := repo.db.Model(&complianceModels.ComplianceRequirement{}).
		Joins("JOIN lifecycle_statuses ON compliance_register.lifecycle_status_id = lifecycle_statuses.id").
		Where("compliance_register.id IN ? AND lifecycle_statuses.type = ?", ids, "active").
		Pluck("compliance_register.id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		// Retrieve the ID of the "archived" LifecycleStatus
		var archivedLifecycleStatusID uint
		if err := repo.db.Model(&masterModels.LifecycleStatus{}).
			Where("type = ?", "archived").
			Pluck("id", &archivedLifecycleStatusID).Error; err != nil {
			return nil, err
		}

		// Update ComplianceRequirements with the retrieved archivedLifecycleStatusID
		if err := repo.db.Model(&complianceModels.ComplianceRequirement{}).
			Where("id IN ?", idsToUpdate).
			Update("lifecycle_status_id", archivedLifecycleStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRequirementRepository) Restore(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Retrieve IDs of ComplianceRequirements with lifecycle status "archived"
	if err := repo.db.Model(&complianceModels.ComplianceRequirement{}).
		Joins("JOIN lifecycle_statuses ON compliance_register.lifecycle_status_id = lifecycle_statuses.id").
		Where("compliance_register.id IN ? AND lifecycle_statuses.type = ?", ids, "archived").
		Pluck("compliance_register.id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		// Retrieve the ID of the "active" LifecycleStatus
		var activeLifecycleStatusID uint
		if err := repo.db.Model(&masterModels.LifecycleStatus{}).
			Where("type = ?", "active").
			Pluck("id", &activeLifecycleStatusID).Error; err != nil {
			return nil, err
		}

		// Update ComplianceRequirements with the retrieved activeLifecycleStatusID
		if err := repo.db.Model(&complianceModels.ComplianceRequirement{}).
			Where("id IN ?", idsToUpdate).
			Update("lifecycle_status_id", activeLifecycleStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRequirementRepository) IsValidID(id uint) (*complianceModels.ComplianceRequirement, error) {
	var complianceRequirement complianceModels.ComplianceRequirement

	if err := repo.db.First(&complianceRequirement, id).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return &complianceRequirement, nil
}
