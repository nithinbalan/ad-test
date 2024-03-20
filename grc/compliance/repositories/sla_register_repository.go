package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SLARepositoryInterface interface {
    GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) error
	Delete(id uint) error
}

// implementation
type SLARepository struct {
	db *gorm.DB
}

func NewSLARepository(db *gorm.DB) *SLARepository {
	return &SLARepository{
		db: db,
	}
}

func (repo *SLARepository) GetAll(sort string, sortBy string, keyword string, page int, pageSize int) (*slaDTOs.SLARegisterListResponseDTO, error) {

	var slaRegisters []*slaDTOs.SLARegisterListDTO
	var totalCount int64

    query := repo.db.Table("sla_registers").
        Select(`
            sla_registers.*,
            sla_categories.title_en AS sla_category_title_en,
            sla_categories.title_ar AS sla_category_title_ar,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            responsible_users.first_name_en AS responsible_user_first_name_en,
            responsible_users.first_name_ar AS responsible_user_first_name_ar,
            responsible_users.last_name_en AS responsible_user_last_name_en,
            responsible_users.last_name_ar AS responsible_user_last_name_ar,
            lifecycle_statuses.title_en AS lifecycle_status_title_en,
            lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
            created_by_users.first_name_en AS created_by_first_name_en,
            created_by_users.first_name_ar AS created_by_first_name_ar,
            created_by_users.last_name_en AS created_by_last_name_en,
            created_by_users.last_name_ar AS created_by_last_name_ar,
            updated_by_users.first_name_en AS updated_by_first_name_en,
            updated_by_users.first_name_ar AS updated_by_first_name_ar,
            updated_by_users.last_name_en AS updated_by_last_name_en,
            updated_by_users.last_name_ar AS updated_by_last_name_ar
        `).
        Joins(`
            LEFT JOIN sla_categories ON sla_registers.sla_category_id = sla_categories.id
            LEFT JOIN severity_levels ON sla_registers.severity_level_id = severity_levels.id
            LEFT JOIN organization_units ON sla_registers.responsible_unit_id = organization_units.id
            LEFT JOIN users AS responsible_users ON sla_registers.responsible_user_id = responsible_users.id
            LEFT JOIN lifecycle_statuses ON sla_registers.lifecycle_status_id = lifecycle_statuses.id
            LEFT JOIN users AS created_by_users ON sla_registers.created_by = created_by_users.id
            LEFT JOIN users AS updated_by_users ON sla_registers.updated_by = updated_by_users.id
        `)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
    // Apply search keyword if exists
	if keyword != "" {
		query = query.Where(`
			(sla_registers.title_en LIKE ? OR
			sla_registers.title_ar LIKE ? OR
			sla_registers.reference_number LIKE ? OR
			sla_registers.description LIKE ? OR
			sla_registers.client LIKE ? OR
			sla_categories.title_en LIKE ? OR
			sla_categories.title_ar LIKE ? OR
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
			"%"+keyword+"%")
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
	err := query.Offset(offset).Limit(pageSize).Find(&slaRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(pageSize) - 1) / int64(pageSize))

	response := &slaDTOs.SLARegisterListResponseDTO{
		Items: slaRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}


func (repo *SLARepository) GetByID(id uint) (*models.SLARegister, error) {
	var slaRegister *models.SLARegister
	// result := repo.db.First(&slaRegister, id)
	result := repo.db.Preload("SLACategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&slaRegister, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return slaRegister, result.Error
}

func (repo *SLARepository) Create(currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) (uint, error) {

	// Retrieve the last reference number from the database
	var lastReferenceNumber string
	err := repo.db.Model(&models.SLARegister{}).
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
	code := "SLA-JETHUR"
	// Update lastDigits based on the last reference number if available
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	// Generate the new reference number using a utility function
	referenceNumber := utility.GenerateReferenceNumber(code, lastDigits)

	// Map DTO to SLARegister model
	slaRegister := &models.SLARegister{
		SLACategoryID:           slaRegFormDTO.SLACategoryID,
		SeverityLevelID:         slaRegFormDTO.SeverityLevelID,
		ResponsibleUnitID:       slaRegFormDTO.ResponsibleUnitID,
		ResponsibleUserID:       slaRegFormDTO.ResponsibleUserID,
		LifecycleStatusID:       slaRegFormDTO.LifecycleStatusID,
		ReferenceNumber:         referenceNumber,
		TitleEn:                 slaRegFormDTO.TitleEn,
		TitleAr:                 slaRegFormDTO.TitleAr,
		Description:             slaRegFormDTO.Description,
		Client:                  slaRegFormDTO.Client,
		CreatedBy:               currentUserID,
	}
    
	if err := repo.db.Create(slaRegister).Error; err != nil {
		return 0, err
	}

	return slaRegister.ID, nil
}

func (repo *SLARepository) Update(id uint, currentUserID uint, slaRegFormDTO *slaDTOs.SLARegisterFormDTO) error {
	// Prepare the update map
	updateFields := map[string]interface{}{
		"SLACategoryID":    slaRegFormDTO.SLACategoryID,
		"SeverityLevelID":         slaRegFormDTO.SeverityLevelID,
		"ResponsibleUnitID":       slaRegFormDTO.ResponsibleUnitID,
		"ResponsibleUserID":       slaRegFormDTO.ResponsibleUserID,
		"LifecycleStatusID":       slaRegFormDTO.LifecycleStatusID,
		"TitleEn":                 slaRegFormDTO.TitleEn,
		"TitleAr":                 slaRegFormDTO.TitleAr,
		"Description":             slaRegFormDTO.Description,
		"Client":                  slaRegFormDTO.Client,
		"UpdatedBy":               currentUserID,
	}
	// Perform the update
	if err := repo.db.Model(&models.SLARegister{}).
		Where("id = ?", id).
		Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARepository) Delete(id uint) error {

	var slaRegister *models.SLARegister
	if err := repo.db.First(&slaRegister, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLARegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}




