package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskAssessmentWorkshopRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateRiskAssessmentWorkshopDocument(document *models.RiskUpdatePlanRiskAssessmentWorkshopDocument, tx ...utility.DatabaseTransaction) (uint, error)
	GetRiskAssessmentWorkshopDocumentByID(id uint) (*models.RiskUpdatePlanRiskAssessmentWorkshopDocument, error)
	GetRiskAssessmentWorkshopDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskAssessmentWorkshopDocument, error)
	DeleteRiskAssessmentWorkshopDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanRiskAssessmentWorkshopRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskAssessmentWorkshopRepository(db *gorm.DB) *RiskUpdatePlanRiskAssessmentWorkshopRepository {
	return &RiskUpdatePlanRiskAssessmentWorkshopRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskAssessmentWorkshopRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskAssessmentWorkshopRepository) CreateRiskAssessmentWorkshopDocument(document *models.RiskUpdatePlanRiskAssessmentWorkshopDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *RiskUpdatePlanRiskAssessmentWorkshopRepository) GetRiskAssessmentWorkshopDocumentByID(id uint) (*models.RiskUpdatePlanRiskAssessmentWorkshopDocument, error) {
	var document models.RiskUpdatePlanRiskAssessmentWorkshopDocument
	if err := repo.db.First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *RiskUpdatePlanRiskAssessmentWorkshopRepository) GetRiskAssessmentWorkshopDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskAssessmentWorkshopDocument, error) {
	var documents []*models.RiskUpdatePlanRiskAssessmentWorkshopDocument
	if err := repo.db.Where("risk_update_plan_timeline_id = ?", timelineID).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *RiskUpdatePlanRiskAssessmentWorkshopRepository) DeleteRiskAssessmentWorkshopDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskUpdatePlanRiskAssessmentWorkshopDocument{}, id).Error; err != nil {
		return err
	}

	return nil
}
