package repositories

import (
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskClosureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint) ([]*dtos.RiskClosureListDTO, error)
	GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*models.RiskClosure, error)
	GetByID(id uint) (*models.RiskClosure, error)
	Create(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.RiskClosureDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskClosureDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
}

type RiskClosureRepository struct {
	db *gorm.DB
}

func NewRiskClosureRepository(db *gorm.DB) RiskClosureRepositoryInterface {
	return &RiskClosureRepository{
		db: db,
	}
}

func (repo *RiskClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskClosureRepository) GetAll(riskID uint) ([]*dtos.RiskClosureListDTO, error) {
	var closures []*dtos.RiskClosureListDTO
	query := repo.db.Table("risk_closures").
		Select(`
			risk_closures.*,
			risk_register.current_risk_rating_id,
			risk_matrix_risk_ratings.title_en AS current_risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS current_risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS current_risk_rating_color_code,
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
		Joins(`
			LEFT JOIN users AS created_user ON created_user.id = risk_closures.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_closures.updated_by_id
			LEFT JOIN risk_register ON risk_register.id = risk_closures.risk_register_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
		`).
		Where("risk_closures.risk_register_id = ?", riskID)

	if err := query.Find(&closures).Error; err != nil {
		return nil, err
	}

	return closures, nil
}

func (repo *RiskClosureRepository) GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*models.RiskClosure, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskClosure models.RiskClosure
	result := db.
		Preload("RiskRegister.CurrentRiskRating").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("risk_register_id = ? AND id = ?", riskID, id).
		First(&riskClosure)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskClosure, nil
}

func (repo *RiskClosureRepository) GetByID(id uint) (*models.RiskClosure, error) {
	var riskClosure models.RiskClosure
	result := repo.db.First(&riskClosure, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &riskClosure, nil
}

func (repo *RiskClosureRepository) Create(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskClosure).Error; err != nil {
		return 0, err
	}
	return riskClosure.ID, nil
}

func (repo *RiskClosureRepository) Update(riskClosure *models.RiskClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(riskClosure).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskClosure{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) CreateDocuments(docs []*models.RiskClosureDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskClosureRepository) GetDocumentByID(id uint) (*models.RiskClosureDocument, error) {
	var document models.RiskClosureDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskClosureRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskClosureDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskClosureRepository) UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskRegister{}).Where("id = ?", risk.ID).Updates(risk).Error; err != nil {
		return err
	}

	return nil
}
