package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type CommentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommentStatus, error)
	GetByID(id uint) (*models.CommentStatus, error)
	GetByType(status string) (*models.CommentStatus, error)
}

// Implementation
type CommentStatusRepository struct {
	db *gorm.DB
}

func NewCommentStatusRepository(db *gorm.DB) *CommentStatusRepository {
	return &CommentStatusRepository{
		db: db,
	}
}

func (repo *CommentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommentStatus, error) {
	var commentStatuses []*models.CommentStatus

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
	result := query.Find(&commentStatuses)

	return commentStatuses, result.Error
}

func (repo *CommentStatusRepository) GetByID(id uint) (*models.CommentStatus, error) {
	var commentStatus *models.CommentStatus
	result := repo.db.First(&commentStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return commentStatus, result.Error
}

func (repo *CommentStatusRepository) GetByType(status string) (*models.CommentStatus, error) {
	var commentStatus *models.CommentStatus

	result := repo.db.Where("type = ?", status).First(&commentStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Action Plan Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return commentStatus, nil
}
