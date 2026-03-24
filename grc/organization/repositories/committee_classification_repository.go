package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type CommitteeClassificationRepositoryInterface interface {
	GetAll(params *dtos.CommitteeClassificationListQueryParamsDTO) ([]*dtos.CommitteeClassificationListDTO, error)
	GetByID(id uint) (*models.CommitteeClassification, error)
	Create(cc *models.CommitteeClassification) (uint, error)
	Update(cc *models.CommitteeClassification) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type CommitteeClassificationRepository struct {
	db *gorm.DB
}

func NewCommitteeClassificationRepository(db *gorm.DB) *CommitteeClassificationRepository {
	return &CommitteeClassificationRepository{db: db}
}

func (repo *CommitteeClassificationRepository) GetAll(params *dtos.CommitteeClassificationListQueryParamsDTO) ([]*dtos.CommitteeClassificationListDTO, error) {
	var ccs []*dtos.CommitteeClassificationListDTO

	query := repo.db.Table("committee_classifications").
		Select("committee_classifications.*")

	if params.Keyword != "" {
		query = query.Where(`
				(committee_classifications.title_en LIKE ? OR
				committee_classifications.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("committee_classifications.status = ?", params.Status)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&ccs)

	return ccs, result.Error
}

func (repo *CommitteeClassificationRepository) GetByID(id uint) (*models.CommitteeClassification, error) {
	var cc *models.CommitteeClassification
	result := repo.db.First(&cc, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return cc, result.Error
}

func (repo *CommitteeClassificationRepository) Create(cc *models.CommitteeClassification) (uint, error) {
	if err := repo.db.Create(cc).Error; err != nil {
		return 0, err
	}
	return cc.ID, nil
}

func (repo *CommitteeClassificationRepository) Update(updatedCC *models.CommitteeClassification) error {
	if err := repo.db.Save(updatedCC).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CommitteeClassificationRepository) Delete(id uint) error {
	if err := repo.db.Delete(&models.CommitteeClassification{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CommitteeClassificationRepository) Activate(ids []uint) ([]uint, error) {
	var successfulIds []uint

	for _, id := range ids {
		result := repo.db.Model(&models.CommitteeClassification{}).
			Where("id = ?", id).
			Update("status", "active")

		if result.Error == nil && result.RowsAffected > 0 {
			successfulIds = append(successfulIds, id)
		}
	}

	return successfulIds, nil
}

func (repo *CommitteeClassificationRepository) Deactivate(ids []uint) ([]uint, error) {
	var successfulIds []uint

	for _, id := range ids {
		result := repo.db.Model(&models.CommitteeClassification{}).
			Where("id = ?", id).
			Update("status", "inactive")

		if result.Error == nil && result.RowsAffected > 0 {
			successfulIds = append(successfulIds, id)
		}
	}

	return successfulIds, nil
}

func (repo *CommitteeClassificationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CommitteeClassification{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	return count > 0, err
}

func (repo *CommitteeClassificationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CommitteeClassification{}).
		Where("id != ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	return count > 0, err
}
