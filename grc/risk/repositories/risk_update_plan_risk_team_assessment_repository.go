package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskTeamAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetRiskTeamAssessmentUsers(timelineID uint) ([]*models.RiskUpdatePlanRiskTeamAssessment, error)
	CreateRiskTeamAssessments(assessments []*models.RiskUpdatePlanRiskTeamAssessment, tx ...utility.DatabaseTransaction) error
	DeleteRiskTeamAssessmentsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error
	GetRiskTeamDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskTeamDocument, error)
	CreateRiskTeamDocument(document *models.RiskUpdatePlanRiskTeamDocument, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanRiskTeamDocument, error)
	DeleteRiskTeamDocument(id uint, tx ...utility.DatabaseTransaction) error
	GetRiskTeamDocumentByID(id uint) (*models.RiskUpdatePlanRiskTeamDocument, error)
}

type RiskUpdatePlanRiskTeamAssessmentRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanRiskTeamAssessmentRepository(db *gorm.DB) *RiskUpdatePlanRiskTeamAssessmentRepository {
	return &RiskUpdatePlanRiskTeamAssessmentRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) GetRiskTeamAssessmentUsers(timelineID uint) ([]*models.RiskUpdatePlanRiskTeamAssessment, error) {
	var assessments []*models.RiskUpdatePlanRiskTeamAssessment

	err := repo.db.
		Preload("RiskUpdatePlanTimeline").
		Preload("RiskChampion").
		Preload("RiskChampion.OrganizationUnit1").
		Preload("RiskChampion.OrganizationUnit2").
		Preload("RiskChampion.OrganizationUnit3").
		Preload("RiskChampion.OrganizationUnit4").
		Preload("RiskChampion.OrganizationUnit5").
		Preload("RiskChampion.UserStatus").
		Preload("RiskChampion.Designation").
		Preload("RiskChampion.UserRoles").
		Preload("RiskChampion.UserSetting").
		Preload("RiskChampion.UserSetting.Timezone").
		Preload("RiskChampion.UserRoles.Role").
		Preload("RiskChampion.UserRoles.Role.LicenseType").
		Preload("RiskChampion.ReportingUser").
		Preload("RiskChampion.LicenseType").
		Preload("RiskChampion.DelegationUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("risk_update_plan_timeline_id = ?", timelineID).
		Find(&assessments).Error

	if err != nil {
		return nil, err
	}

	return assessments, nil
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) CreateRiskTeamAssessments(assessments []*models.RiskUpdatePlanRiskTeamAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(assessments) == 0 {
		return nil
	}

	return db.Create(&assessments).Error
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) DeleteRiskTeamAssessmentsByTimelineID(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_update_plan_timeline_id = ?", timelineID).Delete(&models.RiskUpdatePlanRiskTeamAssessment{}).Error
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) GetRiskTeamDocumentsByTimelineID(timelineID uint) ([]*models.RiskUpdatePlanRiskTeamDocument, error) {
	var documents []*models.RiskUpdatePlanRiskTeamDocument

	err := repo.db.
		Where("risk_update_plan_timeline_id = ?", timelineID).
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) CreateRiskTeamDocument(document *models.RiskUpdatePlanRiskTeamDocument, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanRiskTeamDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return nil, err
	}

	return document, nil
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) DeleteRiskTeamDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.RiskUpdatePlanRiskTeamDocument{}, id).Error
}

func (repo *RiskUpdatePlanRiskTeamAssessmentRepository) GetRiskTeamDocumentByID(id uint) (*models.RiskUpdatePlanRiskTeamDocument, error) {
	var document models.RiskUpdatePlanRiskTeamDocument

	err := repo.db.First(&document, id).Error
	if err != nil {
		return nil, err
	}

	return &document, nil
}
