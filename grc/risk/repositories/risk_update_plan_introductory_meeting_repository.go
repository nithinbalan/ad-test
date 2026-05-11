package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanIntroductoryMeetingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	// Introductory Meeting Documents
	CreateIntroductoryMeetingDocument(document *models.RiskUpdatePlanIntroductoryMeetingDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetIntroductoryMeetingDocumentByID(id uint) (*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	GetIntroductoryMeetingDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanIntroductoryMeetingDocument, error)
	DeleteIntroductoryMeetingDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanIntroductoryMeetingRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanIntroductoryMeetingRepository(db *gorm.DB) *RiskUpdatePlanIntroductoryMeetingRepository {
	return &RiskUpdatePlanIntroductoryMeetingRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanIntroductoryMeetingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// Introductory Meeting Document Methods

func (repo *RiskUpdatePlanIntroductoryMeetingRepository) CreateIntroductoryMeetingDocument(document *models.RiskUpdatePlanIntroductoryMeetingDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskUpdatePlanIntroductoryMeetingRepository) GetIntroductoryMeetingDocumentByID(id uint) (*models.RiskUpdatePlanIntroductoryMeetingDocument, error) {
	var document models.RiskUpdatePlanIntroductoryMeetingDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}

	return &document, nil
}

func (repo *RiskUpdatePlanIntroductoryMeetingRepository) GetIntroductoryMeetingDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanIntroductoryMeetingDocument, error) {
	var documents []*models.RiskUpdatePlanIntroductoryMeetingDocument
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskUpdatePlanIntroductoryMeetingRepository) DeleteIntroductoryMeetingDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanIntroductoryMeetingDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}
