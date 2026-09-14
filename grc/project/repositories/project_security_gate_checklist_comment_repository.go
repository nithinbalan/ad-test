package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistComment, error)
	Create(comment *models.ProjectSecurityGateChecklistComment, tx ...utility.DatabaseTransaction) error
	Update(comment *models.ProjectSecurityGateChecklistComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(documents []*models.ProjectSecurityGateChecklistCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsByCommentID(commentID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateChecklistCommentRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistCommentRepository(db *gorm.DB) ProjectSecurityGateChecklistCommentRepositoryInterface {
	return &ProjectSecurityGateChecklistCommentRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateChecklistCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateChecklistCommentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment *models.ProjectSecurityGateChecklistComment
	result := db.
		Preload("Documents").
		First(&comment, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *ProjectSecurityGateChecklistCommentRepository) Create(comment *models.ProjectSecurityGateChecklistComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(comment).Error
}

func (repo *ProjectSecurityGateChecklistCommentRepository) Update(comment *models.ProjectSecurityGateChecklistComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(comment).Error
}

func (repo *ProjectSecurityGateChecklistCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.ProjectSecurityGateChecklistComment{}, id).Error
}

func (repo *ProjectSecurityGateChecklistCommentRepository) CreateDocuments(documents []*models.ProjectSecurityGateChecklistCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(documents) == 0 {
		return nil
	}

	return db.Create(&documents).Error
}

func (repo *ProjectSecurityGateChecklistCommentRepository) GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistCommentDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document *models.ProjectSecurityGateChecklistCommentDocument
	result := db.First(&document, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return document, nil
}

func (repo *ProjectSecurityGateChecklistCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(documentIDs) == 0 {
		return nil
	}

	return db.Where("id IN (?)", documentIDs).Delete(&models.ProjectSecurityGateChecklistCommentDocument{}).Error
}

func (repo *ProjectSecurityGateChecklistCommentRepository) DeleteDocumentsByCommentID(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("project_security_gate_checklist_comment_id = ?", commentID).Delete(&models.ProjectSecurityGateChecklistCommentDocument{}).Error
}
