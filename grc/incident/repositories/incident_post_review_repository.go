package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentPostReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(incidentID uint) (*models.IncidentPostReview, error)
	GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*models.IncidentPostReview, error)
	Create(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateLessons(lessons []*models.IncidentPostReviewLesson, tx ...utility.DatabaseTransaction) error
	DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error
	CreateParticipants(users []*models.IncidentPostReviewUser, tx ...utility.DatabaseTransaction) error
	DeleteParticipants(userIDs []*uint, tx ...utility.DatabaseTransaction) error
	CreateRecommendations(recommendations []*models.IncidentPostReviewRecommendation, tx ...utility.DatabaseTransaction) error
	DeleteRecommendations(recommendationIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentPostReviewRepository struct {
	db *gorm.DB
}

func NewIncidentPostReviewRepository(db *gorm.DB) *IncidentPostReviewRepository {
	return &IncidentPostReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentPostReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentPostReviewRepository) GetByID(incidentID uint) (*models.IncidentPostReview, error) {
	var incidentPostReview *models.IncidentPostReview
	result := repo.db.Where("incident_register_id = ?", incidentID).First(&incidentPostReview)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentPostReview, nil
}

func (repo *IncidentPostReviewRepository) GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*models.IncidentPostReview, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentPostReview *models.IncidentPostReview

	result := db.Where("incident_register_id = ?", incidentID).
		Preload("IncidentRegister").
		Preload("IncidentPostReviewEffectiveness").
		Preload("Recommendations").
		Preload("Users.User").
		Preload("Lessons").
		First(&incidentPostReview)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentPostReview, nil
}

func (repo *IncidentPostReviewRepository) Create(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentPostReview).Error; err != nil {
		return 0, err
	}
	return incidentPostReview.ID, nil
}

func (repo *IncidentPostReviewRepository) Update(incidentPostReview *models.IncidentPostReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentPostReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentPostReview *models.IncidentPostReview
	if err := db.First(&incidentPostReview, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentPostReview{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateLessons(lessons []*models.IncidentPostReviewLesson, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&lessons).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteLessons(lessonsIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete compliance requirements by IDs
	if err := db.Where("id IN (?)", lessonsIDs).Delete(&models.IncidentPostReviewLesson{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateParticipants(users []*models.IncidentPostReviewUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteParticipants(userIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete users by IDs
	if err := db.Where("id IN (?)", userIDs).Delete(&models.IncidentPostReviewUser{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentPostReviewRepository) CreateRecommendations(recommendations []*models.IncidentPostReviewRecommendation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&recommendations).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentPostReviewRepository) DeleteRecommendations(recommendationIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete Incident PostReview Recommendation by IDs
	if err := db.Where("id IN (?)", recommendationIDs).Delete(&models.IncidentPostReviewRecommendation{}).Error; err != nil {
		return err
	}

	return nil
}
