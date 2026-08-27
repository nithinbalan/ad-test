package repositories

import (
	"errors"
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type PenetrationTestPlanRepositoryInterface interface {
	GetAll(string, string, string, string) ([]*models.PenetrationTestPlan, error)
	GetByID(uint) (*models.PenetrationTestPlan, error)
	Create(*models.PenetrationTestPlan) (*models.PenetrationTestPlan, error)
	Update(uint, *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error)
	Delete(uint) error
	Activate([]uint) ([]uint, error)
	Deactivate([]uint) ([]uint, error)
	ExistsByTitle(string, string) (bool, error)
	ExistsByTitleExcludingID(uint, string, string) (bool, error)
}

type PenetrationTestPlanRepository struct{ db *gorm.DB }

func NewPenetrationTestPlanRepository(db *gorm.DB) *PenetrationTestPlanRepository {
	return &PenetrationTestPlanRepository{db: db}
}

func (repo *PenetrationTestPlanRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.PenetrationTestPlan, error) {
	var values []*models.PenetrationTestPlan
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
	return values, q.Order(sortBy + " " + sort).Find(&values).Error
}

func (repo *PenetrationTestPlanRepository) GetByID(id uint) (*models.PenetrationTestPlan, error) {
	var value *models.PenetrationTestPlan
	r := repo.db.First(&value, id)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(r.Error, id)
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return value, nil
}

func (repo *PenetrationTestPlanRepository) Create(value *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error) {
	r := repo.db.Create(&value)
	return value, r.Error
}

func (repo *PenetrationTestPlanRepository) Update(id uint, updated *models.PenetrationTestPlan) (*models.PenetrationTestPlan, error) {
	var value *models.PenetrationTestPlan
	if err := repo.db.First(&value, id).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&value).Updates(updated).Error; err != nil {
		return nil, err
	}
	return value, nil
}

func (repo *PenetrationTestPlanRepository) Delete(id uint) error {
	var value *models.PenetrationTestPlan
	if err := repo.db.First(&value, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.PenetrationTestPlan{}, id).Error
}

func (repo *PenetrationTestPlanRepository) Activate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "inactive", "active")
}

func (repo *PenetrationTestPlanRepository) Deactivate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "active", "inactive")
}

func (repo *PenetrationTestPlanRepository) updateStatus(ids []uint, from, to string) ([]uint, error) {
	var u []uint
	if err := repo.db.Model(&models.PenetrationTestPlan{}).Where("id IN ? AND status = ?", ids, from).Pluck("id", &u).Error; err != nil {
		return nil, err
	}
	if len(u) > 0 {
		if err := repo.db.Model(&models.PenetrationTestPlan{}).Where("id IN ?", u).Update("status", to).Error; err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (repo *PenetrationTestPlanRepository) ExistsByTitle(a, b string) (bool, error) {
	var c int64
	e := repo.db.Model(&models.PenetrationTestPlan{}).Where("title_en = ? OR title_ar = ?", a, b).Count(&c).Error
	if e != nil {
		return false, utility.NewInternalServerError(e)
	}
	return c > 0, nil
}

func (repo *PenetrationTestPlanRepository) ExistsByTitleExcludingID(id uint, a, b string) (bool, error) {
	var c int64
	e := repo.db.Model(&models.PenetrationTestPlan{}).Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, a, b).Count(&c).Error
	if e != nil {
		return false, utility.NewInternalServerError(e)
	}
	return c > 0, nil
}
