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

type ComplianceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*complianceModels.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Archive(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ComplianceRepository struct {
	db *gorm.DB
}

func NewComplianceRepository(db *gorm.DB) *ComplianceRepository {
	return &ComplianceRepository{
		db: db,
	}
}

func (repo *ComplianceRepository) GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {

	var complianceRegisters []*complianceDTOs.ComplianceRegisterListDTO
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

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if keyword != "" {
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
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Sort by the specified column and order
	if sortBy == "" {
		sortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(sortBy + " " + sort)

	offset := (page - 1) * pageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(pageSize).Find(&complianceRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	response := &complianceDTOs.ComplianceRegisterListResponseDTO{
		Items: complianceRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceRepository) GetByID(id uint) (*complianceModels.ComplianceRegister, error) {
	var compliance *complianceModels.ComplianceRegister
	result := repo.db.First(&compliance, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return compliance, result.Error
}

func (repo *ComplianceRepository) Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {

	var lifecycleStatus masterModels.LifecycleStatus
	if err := repo.db.Where("type = ?", "active").First(&lifecycleStatus).Error; err != nil {
		return 0, utility.NewInternalServerError(err)
	}

	// Retrieve the last reference number from the database
	var lastReferenceNumber string
	err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	// Handle potential errors from the database query
	if err != nil {
		return 0, err
	}

	// Default suffix for the reference number if no previous record exists
	lastDigits := "0000"
	// Update lastDigits based on the last reference number if available
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	// Generate the new reference number using a utility function
	referenceNumber := utility.GenerateReferenceNumber(complianceRegFormDTO.StandardCode, lastDigits)

	// Map DTO to ComplianceRegister model
	complianceRegister := &complianceModels.ComplianceRegister{
		ComplianceCategoryID:    complianceRegFormDTO.ComplianceCategoryID,
		ComplianceStatusID:      complianceRegFormDTO.ComplianceStatusID,
		ResponsibleUnitID:       complianceRegFormDTO.ResponsibleUnitID,
		ResponsibleUserID:       complianceRegFormDTO.ResponsibleUserID,
		ReviewFrequencyID:       complianceRegFormDTO.ReviewFrequencyID,
		ReviewResponsibleUserID: complianceRegFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:         complianceRegFormDTO.SeverityLevelID,
		ReferenceNumber:         referenceNumber,
		StandardCode:            complianceRegFormDTO.StandardCode,
		TitleEn:                 complianceRegFormDTO.TitleEn,
		TitleAr:                 complianceRegFormDTO.TitleAr,
		DescriptionEn:           complianceRegFormDTO.DescriptionEn,
		DescriptionAr:           complianceRegFormDTO.DescriptionAr,
		Version:                 complianceRegFormDTO.Version,
		Deadline:                complianceRegFormDTO.Deadline,
		IssueDate:               complianceRegFormDTO.IssueDate,
		ExpiryDate:              complianceRegFormDTO.ExpiryDate,
		NextReviewDate:          complianceRegFormDTO.NextReviewDate,
		CreatedBy:               currentUserID,
		LifecycleStatusID:       lifecycleStatus.ID,
	}

	if err := repo.db.Create(complianceRegister).Error; err != nil {
		return 0, err
	}

	return complianceRegister.ID, nil
}

func (repo *ComplianceRepository) Update(id uint, currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	// Prepare the update map
	updateFields := map[string]interface{}{
		"ComplianceCategoryID":    complianceRegFormDTO.ComplianceCategoryID,
		"ComplianceStatusID":      complianceRegFormDTO.ComplianceStatusID,
		"ResponsibleUnitID":       complianceRegFormDTO.ResponsibleUnitID,
		"ResponsibleUserID":       complianceRegFormDTO.ResponsibleUserID,
		"ReviewFrequencyID":       complianceRegFormDTO.ReviewFrequencyID,
		"ReviewResponsibleUserID": complianceRegFormDTO.ReviewResponsibleUserID,
		"SeverityLevelID":         complianceRegFormDTO.SeverityLevelID,
		"StandardCode":            complianceRegFormDTO.StandardCode,
		"TitleEn":                 complianceRegFormDTO.TitleEn,
		"TitleAr":                 complianceRegFormDTO.TitleAr,
		"DescriptionEn":           complianceRegFormDTO.DescriptionEn,
		"DescriptionAr":           complianceRegFormDTO.DescriptionAr,
		"Version":                 complianceRegFormDTO.Version,
		"Deadline":                complianceRegFormDTO.Deadline,
		"IssueDate":               complianceRegFormDTO.IssueDate,
		"ExpiryDate":              complianceRegFormDTO.ExpiryDate,
		"NextReviewDate":          complianceRegFormDTO.NextReviewDate,
		"UpdatedBy":               currentUserID,
	}

	// Perform the update
	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Where("id = ?", id).
		Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRepository) Delete(id uint) error {

	var compliance *complianceModels.ComplianceRegister
	if err := repo.db.First(&compliance, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&complianceModels.ComplianceRegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Retrieve IDs of ComplianceRegisters with lifecycle status "archived"
	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
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

		// Update ComplianceRegisters with the retrieved activeLifecycleStatusID
		if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("lifecycle_status_id", activeLifecycleStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRepository) Archive(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Retrieve IDs of ComplianceRegisters with lifecycle status "active"
	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
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

		// Update ComplianceRegisters with the retrieved archivedLifecycleStatusID
		if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("lifecycle_status_id", archivedLifecycleStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ComplianceRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
