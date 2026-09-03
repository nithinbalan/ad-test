package repositories

import (
	"errors"
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type PenetrationTestEnvironmentReadinessRepositoryInterface interface {
	GetAll(string, string, string, string) ([]*models.PenetrationTestEnvironmentReadiness, error)
	GetByID(uint) (*models.PenetrationTestEnvironmentReadiness, error)
	Create(*models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error)
	Update(uint, *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error)
	Delete(uint) error
	Activate([]uint) ([]uint, error)
	Deactivate([]uint) ([]uint, error)
	ExistsByTitle(string, string) (bool, error)
	ExistsByTitleExcludingID(uint, string, string) (bool, error)
}

type PenetrationTestEnvironmentReadinessRepository struct {
	db *gorm.DB
}

func NewPenetrationTestEnvironmentReadinessRepository(db *gorm.DB) *PenetrationTestEnvironmentReadinessRepository {
	return &PenetrationTestEnvironmentReadinessRepository{
		db: db,
	}
}

func (repo *PenetrationTestEnvironmentReadinessRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestEnvironmentReadiness, error) {
	var v []*models.PenetrationTestEnvironmentReadiness
	q := repo.db
	if keyword != "" {
		q = q.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if status == "active" || status == "inactive" {
		q = q.Where("status = ?", status)
	}
	return v, q.Order(sortBy + " " + sort).Find(&v).Error
}

func (repo *PenetrationTestEnvironmentReadinessRepository) GetByID(id uint) (*models.PenetrationTestEnvironmentReadiness, error) {
	var v *models.PenetrationTestEnvironmentReadiness
	x := repo.db.First(&v, id)
	if errors.Is(x.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(x.Error, id)
	}
	if x.Error != nil {
		return nil, x.Error
	}
	return v, nil
}

func (repo *PenetrationTestEnvironmentReadinessRepository) Create(v *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error) {
	x := repo.db.Create(&v)
	return v, x.Error
}

func (repo *PenetrationTestEnvironmentReadinessRepository) Update(id uint, u *models.PenetrationTestEnvironmentReadiness) (*models.PenetrationTestEnvironmentReadiness, error) {
	var v *models.PenetrationTestEnvironmentReadiness
	if x := repo.db.First(&v, id).Error; x != nil {
		return nil, x
	}
	if x := repo.db.Model(&v).Updates(u).Error; x != nil {
		return nil, x
	}
	return v, nil
}

func (repo *PenetrationTestEnvironmentReadinessRepository) Delete(id uint) error {
	var v *models.PenetrationTestEnvironmentReadiness
	if x := repo.db.First(&v, id).Error; x != nil {
		return x
	}
	return repo.db.Delete(&models.PenetrationTestEnvironmentReadiness{}, id).Error
}

func (repo *PenetrationTestEnvironmentReadinessRepository) Activate(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return repo.updateStatus(ids, "inactive", "active")
}

func (repo *PenetrationTestEnvironmentReadinessRepository) Deactivate(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return repo.updateStatus(ids, "active", "inactive")
}

func (repo *PenetrationTestEnvironmentReadinessRepository) updateStatus(ids []uint, a, b string) ([]uint, error) {
	var u []uint
	if x := repo.db.Model(&models.PenetrationTestEnvironmentReadiness{}).Where("id IN ? AND status = ?", ids, a).Pluck("id", &u).Error; x != nil {
		return nil, x
	}
	if len(u) > 0 {
		if x := repo.db.Model(&models.PenetrationTestEnvironmentReadiness{}).Where("id IN ?", u).Update("status", b).Error; x != nil {
			return nil, x
		}
	}
	return u, nil
}

func (repo *PenetrationTestEnvironmentReadinessRepository) ExistsByTitle(a, b string) (bool, error) {
	var c int64
	x := repo.db.Model(&models.PenetrationTestEnvironmentReadiness{}).Where("title_en = ? OR title_ar = ?", a, b).Count(&c).Error
	if x != nil {
		return false, utility.NewInternalServerError(x)
	}
	return c > 0, nil
}

func (repo *PenetrationTestEnvironmentReadinessRepository) ExistsByTitleExcludingID(id uint, a, b string) (bool, error) {
	var c int64
	x := repo.db.Model(&models.PenetrationTestEnvironmentReadiness{}).Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, a, b).Count(&c).Error
	if x != nil {
		return false, utility.NewInternalServerError(x)
	}
	return c > 0, nil
}
