package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TPRMComplianceResponseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TPRMComplianceResponse, error)
	GetByID(id uint) (*models.TPRMComplianceResponse, error)
	GetFullDetailsByID(id uint) (*models.TPRMComplianceResponse, error)
	Create(entity *models.TPRMComplianceResponse, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entity *models.TPRMComplianceResponse, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	ExistsByScore(score int) (bool, error)
	ExistsByScoreExcludingID(id uint, score int) (bool, error)
}

type TPRMComplianceResponseRepository struct {
	db *gorm.DB
}

func NewTPRMComplianceResponseRepository(db *gorm.DB) *TPRMComplianceResponseRepository {
	return &TPRMComplianceResponseRepository{db: db}
}

func (repo *TPRMComplianceResponseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TPRMComplianceResponseRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TPRMComplianceResponse, error) {
	var complianceResponses []*models.TPRMComplianceResponse
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "score"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceResponses)

	return complianceResponses, result.Error
}

func (repo *TPRMComplianceResponseRepository) GetByID(id uint) (*models.TPRMComplianceResponse, error) {
	var entity *models.TPRMComplianceResponse
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *TPRMComplianceResponseRepository) GetFullDetailsByID(id uint) (*models.TPRMComplianceResponse, error) {
	var entity *models.TPRMComplianceResponse
	if err := repo.db.Preload("CreatedByUser").
		First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *TPRMComplianceResponseRepository) Create(entity *models.TPRMComplianceResponse, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *TPRMComplianceResponseRepository) Update(entity *models.TPRMComplianceResponse, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(entity).Error
}

func (repo *TPRMComplianceResponseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.TPRMComplianceResponse{}, id).Error
}

func (repo *TPRMComplianceResponseRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.TPRMComplianceResponse{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *TPRMComplianceResponseRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.TPRMComplianceResponse{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *TPRMComplianceResponseRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *TPRMComplianceResponseRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *TPRMComplianceResponseRepository) ExistsByScore(score int) (bool, error) {
	var count int64

	err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("score = ?", score).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *TPRMComplianceResponseRepository) ExistsByScoreExcludingID(id uint, score int) (bool, error) {
	var count int64

	err := repo.db.Model(&models.TPRMComplianceResponse{}).
		Where("id <> ? AND (score = ?)", id, score).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
