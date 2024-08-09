package repositories

import (
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type DocumentClauseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *documentDTOs.DocumentClauseListQueryParamsDTO, documentVersionId uint) ([]*models.DocumentVersionClause, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	ExistsByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (bool, error)
	Create(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) (uint, error)
	Update(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateDocumentControls(documentControls []*models.DocumentControl, tx ...utility.DatabaseTransaction) error
	DeleteDocumentControls(controlIDs []*uint, tx ...utility.DatabaseTransaction) error
}

// Repository implementation
type DocumentClauseRepository struct {
	db *gorm.DB
}

// Constructor for DocumentClauseRepository
func NewDocumentClauseRepository(db *gorm.DB) *DocumentClauseRepository {
	return &DocumentClauseRepository{
		db: db,
	}
}

func (repo *DocumentClauseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentClauseRepository) GetAll(params *documentDTOs.DocumentClauseListQueryParamsDTO, documentVersionId uint) ([]*models.DocumentVersionClause, error) {
	var clauses []*models.DocumentVersionClause
	query := repo.db.Model(&models.DocumentVersionClause{}).
		Preload("DocumentControls").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("DocumentControls.Control").
		Where("document_version_id = ?", documentVersionId)

	// Dynamic filtering based on the keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(document_version_clauses.clause_number LIKE ? OR
			document_version_clauses.clause_title_en LIKE ? OR
			document_version_clauses.clause_title_ar LIKE ?)`,
			keyword, keyword, keyword)
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "clause_number" // Default sort by clause_number if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&clauses).Error; err != nil {
		return nil, err
	}

	return clauses, nil
}

func (repo *DocumentClauseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	var clause models.DocumentVersionClause
	db := repo.db

	// Check if transaction is provided and use it
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Where("id = ?", id).First(&clause).Error
	if err != nil {
		return nil, err
	}
	return &clause, nil
}

func (repo *DocumentClauseRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentVersionClause *models.DocumentVersionClause

	result := db.
		Preload("ParentClause").
		Preload("ChildClauses").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("DocumentControls").
		First(&documentVersionClause, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentVersionClause, nil
}

func (repo *DocumentClauseRepository) GetByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	var clause models.DocumentVersionClause
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Where("clause_number = ?", clauseNumber).First(&clause).Error
	if err != nil {
		return nil, err
	}
	return &clause, nil
}

func (repo *DocumentClauseRepository) ExistsByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (bool, error) {
	var count int64
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Model(&models.DocumentVersionClause{}).Where("clause_number = ?", clauseNumber).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *DocumentClauseRepository) Create(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentVersionClause).Error; err != nil {
		return 0, err
	}

	return documentVersionClause.ID, nil
}

func (repo *DocumentClauseRepository) Update(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentVersionClause).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentClauseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.DocumentVersionClause{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentClauseRepository) CreateDocumentControls(documentControls []*models.DocumentControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documentControls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentClauseRepository) DeleteDocumentControls(controlIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", controlIDs).Delete(&models.DocumentControl{}).Error; err != nil {
		return err
	}

	return nil
}
