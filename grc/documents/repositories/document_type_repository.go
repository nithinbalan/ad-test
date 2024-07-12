package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, docmentType string, status string) ([]*models.DocumentType, error)
	GetByID(id uint) (*models.DocumentType, error)
	Create(documentType *models.DocumentType) (*models.DocumentType, error)
	Update(id uint, documentType *models.DocumentType) (*models.DocumentType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)
}

// Implementation
type DocumentTypeRepository struct {
	db *gorm.DB
}

func NewDocumentTypeRepository(db *gorm.DB) *DocumentTypeRepository {
	return &DocumentTypeRepository{
		db: db,
	}
}

func (repo *DocumentTypeRepository) GetAll(sort string, sortBy string, keyword string, docmentType string, status string) ([]*models.DocumentType, error) {
	var documentTypes []*models.DocumentType
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

	if docmentType == "policy" || docmentType == "procedure" || docmentType == "others" {
		query = query.Where("type = ?", docmentType)
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&documentTypes)

	return documentTypes, result.Error
}

func (repo *DocumentTypeRepository) GetByID(id uint) (*models.DocumentType, error) {
	var documentType *models.DocumentType
	result := repo.db.First(&documentType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentType, nil
}

func (repo *DocumentTypeRepository) Create(documentType *models.DocumentType) (*models.DocumentType, error) {
	result := repo.db.Create(&documentType)

	return documentType, result.Error
}

func (repo *DocumentTypeRepository) Update(id uint, updatedDocumentType *models.DocumentType) (*models.DocumentType, error) {
	var documentType *models.DocumentType
	if err := repo.db.First(&documentType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&documentType).Updates(updatedDocumentType).Error; err != nil {
		return nil, err
	}

	return documentType, nil
}

func (repo *DocumentTypeRepository) Delete(id uint) error {
	var documentType *models.DocumentType
	if err := repo.db.First(&documentType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.DocumentType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DocumentTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *DocumentTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
