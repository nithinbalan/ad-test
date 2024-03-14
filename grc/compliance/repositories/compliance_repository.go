package repositories

import (
	"errors"
	complianceModels "jethur/grc/compliance/models"
	masterModels "jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceModels.ComplianceRegister, error)
	GetByID(id uint) (*complianceModels.ComplianceRegister, error)
	Create(currentUserID uint, compliance *complianceModels.ComplianceRegister) (*complianceModels.ComplianceRegister, error)
	Update(id uint, compliance *complianceModels.ComplianceRegister) (*complianceModels.ComplianceRegister, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
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

func (repo *ComplianceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*complianceModels.ComplianceRegister, error) {
	var compliances []*complianceModels.ComplianceRegister

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&compliances)

	return compliances, result.Error
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

func (repo *ComplianceRepository) Create(currentUserID uint, compliance *complianceModels.ComplianceRegister) (*complianceModels.ComplianceRegister, error) {

	var lifecycleStatus masterModels.LifecycleStatus
	err := repo.db.Where("type = ?", "active").First(&lifecycleStatus).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(err, lifecycleStatus.ID)
	} else if err != nil {
		return nil, err
	}

	// Retrieve the last compliance entry
	var lastID uint
	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).Select("id").Order("id desc").Limit(1).Pluck("id", &lastID).Error; err != nil {
		return nil, err
	}

	// Generate reference number using the last ID
	referenceNumber := utility.GenerateReferenceNumber(compliance.StandardCode, lastID)

	compliance.ReferenceNumber = referenceNumber
	compliance.CreatedBy = currentUserID
	compliance.LifecycleStatusID = lifecycleStatus.ID

	result := repo.db.Create(&compliance)
	return compliance, result.Error
}

func (repo *ComplianceRepository) Update(id uint, updatedCompliance *complianceModels.ComplianceRegister) (*complianceModels.ComplianceRegister, error) {
	var compliance *complianceModels.ComplianceRegister

	if err := repo.db.First(&compliance, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&compliance).Updates(updatedCompliance).Error; err != nil {
		return nil, err
	}

	return compliance, nil
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

	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&complianceModels.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
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
