package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentClosureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.IncidentClosure, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosure, error)
	Create(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateLessons(lessons []*models.IncidentClosureLesson, tx ...utility.DatabaseTransaction) error
	DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error

}

type IncidentClosureRepository struct {
	db *gorm.DB
}

func NewIncidentClosureRepository(db *gorm.DB) *IncidentClosureRepository {
	return &IncidentClosureRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentClosureRepository) GetByID(id uint) (*models.IncidentClosure, error) {
	var incidentClosure *models.IncidentClosure
	result := repo.db.First(&incidentClosure, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosure, nil
}

func (repo *IncidentClosureRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosure, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentClosure *models.IncidentClosure

	result := db.
		Preload("IncidentRegister").
		Preload("Lessons").
		First(&incidentClosure, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosure, nil
}

func (repo *IncidentClosureRepository) Create(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentClosure).Error; err != nil {
		return 0, err
	}
	return incidentClosure.ID, nil
}

func (repo *IncidentClosureRepository) Update(incidentClosure *models.IncidentClosure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentClosure).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentClosure *models.IncidentClosure
	if err := db.First(&incidentClosure, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentClosure{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentClosureRepository) CreateLessons(lessons []*models.IncidentClosureLesson, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&lessons).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentClosureRepository) DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", lessonsIDs).Delete(&models.IncidentClosureLesson{}).Error; err != nil {
		return err
	}

	return nil
}