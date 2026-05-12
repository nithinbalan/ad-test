package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskWorkshopAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskWorkshopAssessmentDocument(document *models.RiskUpdatePlanRiskWorkshopAssessmentDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetRiskWorkshopAssessmentDocumentByID(id uint) (*models.RiskUpdatePlanRiskWorkshopAssessmentDocument, error)
	GetRiskWorkshopAssessmentDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskWorkshopAssessmentDocument, error)
	DeleteRiskWorkshopAssessmentDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanRiskWorkshopAssessmentRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskWorkshopAssessmentRepository(db *gorm.DB) *RiskUpdatePlanRiskWorkshopAssessmentRepository {
	return &RiskUpdatePlanRiskWorkshopAssessmentRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskWorkshopAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskWorkshopAssessmentRepository) CreateRiskWorkshopAssessmentDocument(document *models.RiskUpdatePlanRiskWorkshopAssessmentDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskUpdatePlanRiskWorkshopAssessmentRepository) GetRiskWorkshopAssessmentDocumentByID(id uint) (*models.RiskUpdatePlanRiskWorkshopAssessmentDocument, error) {
	var document models.RiskUpdatePlanRiskWorkshopAssessmentDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskUpdatePlanRiskWorkshopAssessmentRepository) GetRiskWorkshopAssessmentDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskWorkshopAssessmentDocument, error) {
	var documents []*models.RiskUpdatePlanRiskWorkshopAssessmentDocument
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *RiskUpdatePlanRiskWorkshopAssessmentRepository) DeleteRiskWorkshopAssessmentDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanRiskWorkshopAssessmentDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}
