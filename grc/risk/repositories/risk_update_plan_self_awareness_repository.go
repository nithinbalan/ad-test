package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanSelfAwarenessRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	// Self-Awareness Training Documents
	CreateSelfAwarenessTrainingDocument(document *models.RiskUpdatePlanSelfAwarenessTraining, tx ...utility.DatabaseTransaction) (uint, error)
	GetSelfAwarenessTrainingDocumentByID(id uint) (*models.RiskUpdatePlanSelfAwarenessTraining, error)
	GetSelfAwarenessTrainingDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanSelfAwarenessTraining, error)
	DeleteSelfAwarenessTrainingDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanSelfAwarenessRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanSelfAwarenessRepository(db *gorm.DB) *RiskUpdatePlanSelfAwarenessRepository {
	return &RiskUpdatePlanSelfAwarenessRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanSelfAwarenessRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// Self-Awareness Training Document Methods

func (repo *RiskUpdatePlanSelfAwarenessRepository) CreateSelfAwarenessTrainingDocument(document *models.RiskUpdatePlanSelfAwarenessTraining, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskUpdatePlanSelfAwarenessRepository) GetSelfAwarenessTrainingDocumentByID(id uint) (*models.RiskUpdatePlanSelfAwarenessTraining, error) {
	var document models.RiskUpdatePlanSelfAwarenessTraining
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *RiskUpdatePlanSelfAwarenessRepository) GetSelfAwarenessTrainingDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanSelfAwarenessTraining, error) {
	var documents []*models.RiskUpdatePlanSelfAwarenessTraining
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskUpdatePlanSelfAwarenessRepository) DeleteSelfAwarenessTrainingDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanSelfAwarenessTraining{}, id).Error; err != nil {
		return err
	}

	return nil
}
