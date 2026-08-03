package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PrivilegedAccountRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PrivilegedAccount, error)
	GetByID(id uint) (*models.PrivilegedAccount, error)
	Create(privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error)
	Update(id uint, privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)
}

type PrivilegedAccountRepository struct {
	db *gorm.DB
}

func NewPrivilegedAccountRepository(db *gorm.DB) *PrivilegedAccountRepository {
	return &PrivilegedAccountRepository{db: db}
}

func (repo *PrivilegedAccountRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PrivilegedAccount, error) {
	var privilegedAccounts []*models.PrivilegedAccount
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if sortBy != "id" && sortBy != "title_en" && sortBy != "title_ar" && sortBy != "status" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	result := query.Order(sortBy + " " + sort).Find(&privilegedAccounts)
	return privilegedAccounts, result.Error
}

func (repo *PrivilegedAccountRepository) GetByID(id uint) (*models.PrivilegedAccount, error) {
	var privilegedAccount models.PrivilegedAccount
	if err := repo.db.First(&privilegedAccount, id).Error; err != nil {
		return nil, err
	}
	return &privilegedAccount, nil
}

func (repo *PrivilegedAccountRepository) Create(privilegedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error) {
	if err := repo.db.Create(privilegedAccount).Error; err != nil {
		return nil, err
	}
	return privilegedAccount, nil
}

func (repo *PrivilegedAccountRepository) Update(id uint, updatedAccount *models.PrivilegedAccount) (*models.PrivilegedAccount, error) {
	var privilegedAccount models.PrivilegedAccount
	if err := repo.db.First(&privilegedAccount, id).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&privilegedAccount).Updates(updatedAccount).Error; err != nil {
		return nil, err
	}
	return &privilegedAccount, nil
}

func (repo *PrivilegedAccountRepository) Delete(id uint) error {
	var privilegedAccount models.PrivilegedAccount
	if err := repo.db.First(&privilegedAccount, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&privilegedAccount).Error
}

func (repo *PrivilegedAccountRepository) Activate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "inactive", "active")
}

func (repo *PrivilegedAccountRepository) Deactivate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "active", "inactive")
}

func (repo *PrivilegedAccountRepository) updateStatus(ids []uint, currentStatus string, newStatus string) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.PrivilegedAccount{}).
		Where("id IN ? AND status = ?", ids, currentStatus).
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PrivilegedAccount{}).
			Where("id IN ?", idsToUpdate).
			Update("status", newStatus).Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *PrivilegedAccountRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PrivilegedAccount{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *PrivilegedAccountRepository) ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PrivilegedAccount{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
