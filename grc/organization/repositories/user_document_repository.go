package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type UserDocumentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(userID uint, params *orgDTOs.UserDocumentListQueryParamsDTO) ([]*orgDTOs.UserDocumentListDTO, error)
	GetByID(id uint) (*models.UserDocument, error)
	GetFullDetailsByID(id uint) (*models.UserDocument, error)
	Create(item *models.UserDocument) (uint, error)
	Update(updatedUserDocument *models.UserDocument) error
	Delete(id uint) error
}

// Implementation
type UserDocumentRepository struct {
	db *gorm.DB
}

func NewUserDocumentRepository(db *gorm.DB) *UserDocumentRepository {
	return &UserDocumentRepository{
		db: db,
	}
}

func (repo *UserDocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *UserDocumentRepository) GetAll(userID uint, params *orgDTOs.UserDocumentListQueryParamsDTO) ([]*orgDTOs.UserDocumentListDTO, error) {
	var result []*orgDTOs.UserDocumentListDTO

	query := repo.db.Table("user_documents").
		Select(`
			user_documents.*,
			CONVERT(VARCHAR, user_documents.expiry_date, 23) AS expiry_date,
			user_document_types.title_en AS user_document_type_title_en,
			user_document_types.title_ar AS user_document_type_title_ar,

			user_document_sub_types.title_en AS user_document_sub_type_title_en,
			user_document_sub_types.title_ar AS user_document_sub_type_title_ar,

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
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`LEFT JOIN user_document_types ON user_document_types.id = user_documents.user_document_type_id`).
		Joins(`LEFT JOIN user_document_sub_types ON user_document_sub_types.id = user_documents.user_document_sub_type_id`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = user_documents.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = user_documents.updated_by_id").
		Where("user_documents.user_id = ?", userID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("user_documents."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			user_document_sub_types.title_en LIKE ? OR
			user_document_sub_types.title_ar LIKE ? OR
			user_document_types.title_en LIKE ? OR
			user_document_types.title_ar LIKE ?
		`, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "user_document_types.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Debug().Find(&result).Error
	return result, err
}

func (repo *UserDocumentRepository) GetByID(id uint) (*models.UserDocument, error) {
	var userDocument *models.UserDocument
	err := repo.db.First(&userDocument, id).Error
	return userDocument, err
}

func (repo *UserDocumentRepository) GetFullDetailsByID(id uint) (*models.UserDocument, error) {

	var userDocument *models.UserDocument

	result := repo.db.
		Preload("UserDocumentType").
		Preload("UserDocumentSubType").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&userDocument, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return userDocument, nil
}

func (repo *UserDocumentRepository) Create(userDocument *models.UserDocument) (uint, error) {
	err := repo.db.Create(userDocument).Error
	return userDocument.ID, err
}

func (repo *UserDocumentRepository) Update(updatedUserDocument *models.UserDocument) error {

	if err := repo.db.Save(updatedUserDocument).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserDocumentRepository) Delete(id uint) error {
	return repo.db.Delete(&models.UserDocument{}, id).Error
}
