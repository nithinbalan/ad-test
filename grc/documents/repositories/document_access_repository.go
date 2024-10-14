package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentAccessRepositoryInterface interface {
	GetSharedUsers(documentId uint, page uint, pageSize uint) ([]*models.User, commonDTOs.PaginationInfoDTO, error)
	GetReviewUsers(documentId uint) ([]*models.User, error)
	GetSharedUnits(documentId uint) ([]*models.OrganizationUnit, error)

	GetExistingSharedUserIDs(documentID uint, userIDs []uint) ([]uint, error)
	GetExistingSharedUnitIDs(documentID uint, unitIDs []uint) ([]uint, error)

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error
	DeleteSharedUsers(documentID uint, userIDs []*uint) error
	DeleteSharedUnits(documentID uint, unitIDs []*uint) error
}

type DocumentAccessRepository struct {
	db *gorm.DB
}

func NewDocumentAccessRepository(db *gorm.DB) *DocumentAccessRepository {
	return &DocumentAccessRepository{db: db}
}

func (repo *DocumentAccessRepository) GetSharedUsers(documentId uint, page uint, pageSize uint) ([]*models.User, commonDTOs.PaginationInfoDTO, error) {
	var users []*models.User
	var totalCount int64

	// First, count the total eligible records
	countQuery := repo.db.Model(&models.DocumentSharedUser{}).
		Where("document_id = ? AND type = 'shared-user'", documentId)
	countQuery.Count(&totalCount)

	// Then fetch the paginated results
	offset := (page - 1) * pageSize
	err := countQuery.
		Joins("JOIN users ON users.id = document_shared_users.user_id").
		Select("users.*").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&users).Error

	if err != nil {
		return nil, commonDTOs.PaginationInfoDTO{}, err
	}

	pagination := commonDTOs.PaginationInfoDTO{
		CurrentPage: int(page),
		PageSize:    int(pageSize),
		TotalCount:  totalCount,
		TotalPages:  int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	return users, pagination, nil
}

func (repo *DocumentAccessRepository) GetReviewUsers(documentId uint) ([]*models.User, error) {
	var users []*models.User
	err := repo.db.
		Model(&models.DocumentSharedUser{}).
		Where("document_id = ? AND type = 'review-user'", documentId).
		Joins("JOIN users ON users.id = document_shared_users.user_id").
		Select("users.*").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *DocumentAccessRepository) GetSharedUnits(documentId uint) ([]*models.OrganizationUnit, error) {
	var units []*models.OrganizationUnit
	err := repo.db.
		Model(&models.DocumentSharedUnit{}).
		Where("document_id = ?", documentId).
		Joins("JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id").
		Select("organization_units.*").
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (repo *DocumentAccessRepository) GetExistingSharedUserIDs(documentID uint, userIDs []uint) ([]uint, error) {
	var existingUserIDs []uint
	err := repo.db.Model(&models.DocumentSharedUser{}).
		Where("document_id = ? AND type = 'shared-user' AND user_id IN (?)", documentID, userIDs).
		Pluck("user_id", &existingUserIDs).Error
	if err != nil {
		return nil, err
	}
	return existingUserIDs, nil
}

func (repo *DocumentAccessRepository) GetExistingSharedUnitIDs(documentID uint, unitIDs []uint) ([]uint, error) {
	var existingUnitIDs []uint
	err := repo.db.Model(&models.DocumentSharedUnit{}).
		Where("document_id = ? AND organization_unit_id IN (?)", documentID, unitIDs).
		Pluck("organization_unit_id", &existingUnitIDs).Error
	if err != nil {
		return nil, err
	}
	return existingUnitIDs, nil
}

func (repo *DocumentAccessRepository) CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(sharedUsers).Error
}

func (repo *DocumentAccessRepository) CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(sharedUnits).Error
}

func (repo *DocumentAccessRepository) DeleteSharedUsers(documentID uint, userIDs []*uint) error {
	return repo.db.Where("document_id = ? AND user_id IN (?)", documentID, userIDs).
		Delete(&models.DocumentSharedUser{}).Error
}

func (repo *DocumentAccessRepository) DeleteSharedUnits(documentID uint, unitIDs []*uint) error {
	return repo.db.Where("document_id = ? AND organization_unit_id IN (?)", documentID, unitIDs).
		Delete(&models.DocumentSharedUnit{}).Error
}
