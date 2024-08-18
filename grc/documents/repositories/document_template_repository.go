package repositories

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentTemplateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(sort string, sortBy string, keyword string, status string) ([]*documentDTOs.DocumentTemplateListDTO, error)
	GetByID(id uint) (*models.DocumentTemplate, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentTemplate, error)
	Create(template *models.DocumentTemplate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(updatedTemplate *models.DocumentTemplate, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)

	//Document
	CreateDocument(doc *models.DocumentTemplateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.DocumentTemplateDocument, error)
	DeleteDocument(docId uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type DocumentTemplateRepository struct {
	db *gorm.DB
}

func NewDocumentTemplateRepository(db *gorm.DB) *DocumentTemplateRepository {
	return &DocumentTemplateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *DocumentTemplateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentTemplateRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*documentDTOs.DocumentTemplateListDTO, error) {
	var documentTemplates []*documentDTOs.DocumentTemplateListDTO

	query := repo.db.Table("document_templates").
		Select(`
			document_templates.id,
			document_templates.title_en,
			document_templates.title_ar,
			document_templates.status,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_template_documents.url,
			document_template_documents.name,
			document_template_documents.size,
			document_template_documents.token,
			document_template_documents.original_name,
			document_template_documents.type,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
        	created_user.first_name_ar AS created_user_first_name_ar,
        	created_user.last_name_en AS created_user_last_name_en,
        	created_user.last_name_ar AS created_user_last_name_ar,   
        	created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
        	updated_user.first_name_ar AS updated_user_first_name_ar,
        	updated_user.last_name_en AS updated_user_last_name_en,
        	updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path,
			document_templates.created_at,
			document_templates.updated_at
        `).
		Joins("LEFT JOIN document_types ON document_types.id = document_templates.document_type_id").
		Joins("LEFT JOIN document_template_documents ON document_template_documents.document_template_id = document_templates.id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = document_templates.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = document_templates.updated_by_id")

	if sortBy == "" {
		sortBy = "document_templates.title_en"
	} else if sortBy == "title_en" {
		sortBy = "document_templates.title_en"
	} else if sortBy == "document_type_title_en" {
		sortBy = "document_types.title_en"
	} else if sortBy == "original_name" {
		sortBy = "document_template_documents.original_name"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	if keyword != "" {
		query = query.Where(`
				  (document_templates.title_en LIKE ? OR
					document_templates.title_ar LIKE ? OR
					document_types.title_en LIKE ? OR
					document_types.title_ar LIKE ? OR
					document_template_documents.original_name LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
			"%"+keyword+"%")
	}

	query = query.Order(sortBy + " " + sort)

	err := query.Find(&documentTemplates).Error
	if err != nil {
		return nil, err
	}

	return documentTemplates, nil

}

func (repo *DocumentTemplateRepository) GetByID(id uint) (*models.DocumentTemplate, error) {
	var template *models.DocumentTemplate
	result := repo.db.First(&template, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return template, nil
}

func (repo *DocumentTemplateRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentTemplate, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentTemplate *models.DocumentTemplate

	result := db.
		Preload("DocumentType").
		Preload("Document").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&documentTemplate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentTemplate, nil
}

func (repo *DocumentTemplateRepository) Create(template *models.DocumentTemplate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(template).Error; err != nil {
		return 0, err
	}

	return template.ID, nil
}

func (repo *DocumentTemplateRepository) Update(updatedTemplate *models.DocumentTemplate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(updatedTemplate).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentTemplateRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the SLA Register
	if err := db.Delete(&models.DocumentTemplate{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentTemplateRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentTemplate{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentTemplate{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentTemplateRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentTemplate{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentTemplate{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentTemplateRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentTemplate{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *DocumentTemplateRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentTemplate{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *DocumentTemplateRepository) CreateDocument(doc *models.DocumentTemplateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentTemplateRepository) GetDocumentByID(id uint) (*models.DocumentTemplateDocument, error) {
	var document models.DocumentTemplateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *DocumentTemplateRepository) DeleteDocument(docId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete Compliance documents by IDs
	if err := db.Where("id = ?", docId).Delete(&models.DocumentTemplateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
