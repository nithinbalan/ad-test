package repositories

import (
	"errors"
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type PenetrationTestExecutorRepositoryInterface interface {
	GetAll(string, string, string, string) ([]*models.PenetrationTestExecutor, error)
	GetByID(uint) (*models.PenetrationTestExecutor, error)
	Create(*models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error)
	Update(uint, *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error)
	Delete(uint) error
	Activate([]uint) ([]uint, error)
	Deactivate([]uint) ([]uint, error)
	ExistsByTitle(string, string) (bool, error)
	ExistsByTitleExcludingID(uint, string, string) (bool, error)
}

type PenetrationTestExecutorRepository struct {
	db *gorm.DB
}

func NewPenetrationTestExecutorRepository(db *gorm.DB) *PenetrationTestExecutorRepository {
	return &PenetrationTestExecutorRepository{
		db: db,
	}
}

func (repo *PenetrationTestExecutorRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutor, error) {
	var v []*models.PenetrationTestExecutor
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

func (repo *PenetrationTestExecutorRepository) GetByID(id uint) (*models.PenetrationTestExecutor, error) {
	var v *models.PenetrationTestExecutor
	x := repo.db.First(&v, id)
	if errors.Is(x.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(x.Error, id)
	}
	if x.Error != nil {
		return nil, x.Error
	}
	return v, nil
}

func (repo *PenetrationTestExecutorRepository) Create(v *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error) {
	x := repo.db.Create(&v)
	return v, x.Error
}

func (repo *PenetrationTestExecutorRepository) Update(id uint, u *models.PenetrationTestExecutor) (*models.PenetrationTestExecutor, error) {
	var v *models.PenetrationTestExecutor
	if x := repo.db.First(&v, id).Error; x != nil {
		return nil, x
	}
	if x := repo.db.Model(&v).Updates(u).Error; x != nil {
		return nil, x
	}
	return v, nil
}

func (repo *PenetrationTestExecutorRepository) Delete(id uint) error {
	var v *models.PenetrationTestExecutor
	if x := repo.db.First(&v, id).Error; x != nil {
		return x
	}
	return repo.db.Delete(&models.PenetrationTestExecutor{}, id).Error
}

func (repo *PenetrationTestExecutorRepository) Activate(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return repo.updateStatus(ids, "inactive", "active")
}

func (repo *PenetrationTestExecutorRepository) Deactivate(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return repo.updateStatus(ids, "active", "inactive")
}

func (repo *PenetrationTestExecutorRepository) updateStatus(ids []uint, a, b string) ([]uint, error) {
	var u []uint
	if x := repo.db.Model(&models.PenetrationTestExecutor{}).Where("id IN ? AND status = ?", ids, a).Pluck("id", &u).Error; x != nil {
		return nil, x
	}
	if len(u) > 0 {
		if x := repo.db.Model(&models.PenetrationTestExecutor{}).Where("id IN ?", u).Update("status", b).Error; x != nil {
			return nil, x
		}
	}
	return u, nil
}

func (repo *PenetrationTestExecutorRepository) ExistsByTitle(a, b string) (bool, error) {
	var c int64
	x := repo.db.Model(&models.PenetrationTestExecutor{}).Where("title_en = ? OR title_ar = ?", a, b).Count(&c).Error
	if x != nil {
		return false, utility.NewInternalServerError(x)
	}
	return c > 0, nil
}

func (repo *PenetrationTestExecutorRepository) ExistsByTitleExcludingID(id uint, a, b string) (bool, error) {
	var c int64
	x := repo.db.Model(&models.PenetrationTestExecutor{}).Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, a, b).Count(&c).Error
	if x != nil {
		return false, utility.NewInternalServerError(x)
	}
	return c > 0, nil
}
