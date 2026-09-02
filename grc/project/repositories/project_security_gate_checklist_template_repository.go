package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistTemplateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll() ([]*models.ProjectSecurityGateChecklistTemplate, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplate, error)
	GetDetailsByID(id uint) (*models.ProjectSecurityGateChecklistTemplate, error)
	GetVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error)
	GetMaxVersionByTemplateID(templateID uint, tx ...utility.DatabaseTransaction) (int, error)
	ExistsVersionByStatus(templateID uint, status string, tx ...utility.DatabaseTransaction) (bool, error)
	CreateVersion(version *models.ProjectSecurityGateChecklistTemplateVersion, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error)
	UpdateVersion(version *models.ProjectSecurityGateChecklistTemplateVersion, tx ...utility.DatabaseTransaction) error
	GetRequirementByID(id uint) (*models.ProjectSecurityGateChecklistTemplateRequirement, error)
	CreateRequirement(requirement *models.ProjectSecurityGateChecklistTemplateRequirement) (*models.ProjectSecurityGateChecklistTemplateRequirement, error)
	UpdateRequirement(id uint, requirement *models.ProjectSecurityGateChecklistTemplateRequirement) (*models.ProjectSecurityGateChecklistTemplateRequirement, error)
	DeleteRequirement(id uint) error
}

type ProjectSecurityGateChecklistTemplateRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistTemplateRepository(db *gorm.DB) *ProjectSecurityGateChecklistTemplateRepository {
	return &ProjectSecurityGateChecklistTemplateRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetAll() ([]*models.ProjectSecurityGateChecklistTemplate, error) {
	var templates []*models.ProjectSecurityGateChecklistTemplate

	result := repo.db.Preload("Stage").Order("id asc").Find(&templates)

	return templates, result.Error
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var template *models.ProjectSecurityGateChecklistTemplate
	result := db.Preload("Stage").First(&template, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return template, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetDetailsByID(id uint) (*models.ProjectSecurityGateChecklistTemplate, error) {
	var template *models.ProjectSecurityGateChecklistTemplate
	result := repo.db.
		Preload("Stage").
		Preload("Versions", func(db *gorm.DB) *gorm.DB {
			return db.Order("version desc")
		}).
		Preload("Versions.CreatedByUser").
		Preload("Versions.PublishedByUser").
		Preload("Versions.Requirements", func(db *gorm.DB) *gorm.DB {
			return db.Order("id asc")
		}).
		First(&template, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return template, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var version *models.ProjectSecurityGateChecklistTemplateVersion

	result := db.
		Preload("CreatedByUser").
		Preload("PublishedByUser").
		Preload("Requirements", func(db *gorm.DB) *gorm.DB {
			return db.Order("id asc")
		}).
		First(&version, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return version, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetMaxVersionByTemplateID(templateID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maxVersion *int

	result := db.Model(&models.ProjectSecurityGateChecklistTemplateVersion{}).
		Where("template_id = ?", templateID).
		Select("MAX(version)").
		Scan(&maxVersion)

	if result.Error != nil {
		return 0, result.Error
	}

	if maxVersion == nil {
		return 0, nil
	}

	return *maxVersion, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) ExistsVersionByStatus(templateID uint, status string, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64

	err := db.Model(&models.ProjectSecurityGateChecklistTemplateVersion{}).
		Where("template_id = ? AND status = ?", templateID, status).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) CreateVersion(version *models.ProjectSecurityGateChecklistTemplateVersion, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(version)
	return version, result.Error
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) UpdateVersion(version *models.ProjectSecurityGateChecklistTemplateVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(version).Error
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) GetRequirementByID(id uint) (*models.ProjectSecurityGateChecklistTemplateRequirement, error) {
	var requirement *models.ProjectSecurityGateChecklistTemplateRequirement
	result := repo.db.First(&requirement, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return requirement, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) CreateRequirement(requirement *models.ProjectSecurityGateChecklistTemplateRequirement) (*models.ProjectSecurityGateChecklistTemplateRequirement, error) {
	result := repo.db.Create(requirement)
	return requirement, result.Error
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) UpdateRequirement(id uint, updatedRequirement *models.ProjectSecurityGateChecklistTemplateRequirement) (*models.ProjectSecurityGateChecklistTemplateRequirement, error) {
	var requirement *models.ProjectSecurityGateChecklistTemplateRequirement
	if err := repo.db.First(&requirement, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&requirement).Updates(map[string]interface{}{
		"title_en":    updatedRequirement.TitleEn,
		"title_ar":    updatedRequirement.TitleAr,
		"is_critical": updatedRequirement.IsCritical,
	}).Error; err != nil {
		return nil, err
	}

	return requirement, nil
}

func (repo *ProjectSecurityGateChecklistTemplateRepository) DeleteRequirement(id uint) error {
	var requirement *models.ProjectSecurityGateChecklistTemplateRequirement
	if err := repo.db.First(&requirement, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.ProjectSecurityGateChecklistTemplateRequirement{}, id).Error
}
