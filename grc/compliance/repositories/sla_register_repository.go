package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SLARegisterRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegisterWithAllField, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(article *models.SLARegister) (*models.SLARegister, error)
	Update(id uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Delete(id uint) error
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SLARegisterRepository struct {
	db *gorm.DB
}

func NewSLARegisterRepository(db *gorm.DB) *SLARegisterRepository {
	return &SLARegisterRepository{
		db: db,
	}
}

func (repo *SLARegisterRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegisterWithAllField, error) {
    var slaRegisters []*models.SLARegisterWithAllField

    query := repo.db.Table("sla_registers").
        Select(
            "sla_registers.id",
            "sla_registers.reference_id",
            "sla_registers.sla_category_id",
            "sla_registers.severity_level_id",
            "sla_registers.responsible_unit_id",
            "sla_registers.responsible_user_id",
            "sla_registers.lifecycle_status_id",
            "sla_registers.title_en",
            "sla_registers.title_ar",
            "sla_registers.description",
            "sla_registers.client",
            "sla_registers.created_by",
            "sla_registers.updated_by",
            "sla_registers.created_at",
            "sla_registers.updated_at",
            "sla_categories.title_en as sla_category_title_en",
            "sla_categories.title_ar as sla_category_title_ar",
            "severity_levels.title_en as severity_level_title_en",
            "severity_levels.title_ar as severity_level_title_ar",
            "organization_units.title_en as responsible_unit_title_en",
            "organization_units.title_ar as responsible_unit_title_ar",
            "responsible_users.first_name_en as responsible_user_first_name_en",
            "responsible_users.first_name_ar as responsible_user_first_name_ar",
            "responsible_users.last_name_en as responsible_user_last_name_en",
            "responsible_users.last_name_ar as responsible_user_last_name_ar",
            "lifecycle_statuses.title_en as lifecycle_status_title_en",
            "lifecycle_statuses.title_ar as lifecycle_status_title_ar",
            "created_by_users.first_name_en as created_by_first_name_en",
            "created_by_users.first_name_ar as created_by_first_name_ar",
            "created_by_users.last_name_en as created_by_last_name_en",
            "created_by_users.last_name_ar as created_by_last_name_ar",
            "updated_by_users.first_name_en as updated_by_first_name_en",
            "updated_by_users.first_name_ar as updated_by_first_name_ar",
            "updated_by_users.last_name_en as updated_by_last_name_en",
            "updated_by_users.last_name_ar as updated_by_last_name_ar",
        ).
		Joins("LEFT JOIN sla_categories ON sla_registers.sla_category_id = sla_categories.id").
        Joins("LEFT JOIN severity_levels ON sla_registers.severity_level_id = severity_levels.id").
        Joins("LEFT JOIN organization_units ON sla_registers.responsible_unit_id = organization_units.id").
        Joins("LEFT JOIN users AS responsible_users ON sla_registers.responsible_user_id = responsible_users.id").
        Joins("LEFT JOIN lifecycle_statuses ON sla_registers.lifecycle_status_id = lifecycle_statuses.id").
        Joins("LEFT JOIN users AS created_by_users ON sla_registers.created_by = created_by_users.id").
        Joins("LEFT JOIN users AS updated_by_users ON sla_registers.updated_by = updated_by_users.id")
    if keyword != "" {
        query = query.Where("sla_registers.title_en LIKE ? OR sla_registers.title_ar LIKE ? OR sla_registers.description LIKE ? OR sla_registers.client LIKE ?",
            "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
    }
    if sortBy == "" {
        sortBy = "title_en"
    }
    if sort != "asc" && sort != "desc" {
        sort = "asc"
    }

    // Order the results
    query = query.Order("sla_registers." + sortBy + " " + sort)

    // Execute the query and return result/error
    if err := query.Find(&slaRegisters).Error; err != nil {
        return nil, err
    }

    return slaRegisters, nil
}


func (repo *SLARegisterRepository) GetByID(id uint) (*models.SLARegister, error) {
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

func (repo *SLARegisterRepository) Create(slaRegister *models.SLARegister) (*models.SLARegister, error) {
	result := repo.db.Create(&slaRegister)
	
	return slaRegister, result.Error
}

func (repo *SLARegisterRepository) Update(id uint, updatedRegister *models.SLARegister) (*models.SLARegister, error) {
	var slaRegister *models.SLARegister

	if err := repo.db.First(&slaRegister, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&slaRegister).Updates(updatedRegister).Error; err != nil {
		return nil, err
	}

	return slaRegister, nil
}

func (repo *SLARegisterRepository) Delete(id uint) error {

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


func (repo *SLARegisterRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.SLARegister{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}


func (repo *SLARegisterRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SLARegister{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}


