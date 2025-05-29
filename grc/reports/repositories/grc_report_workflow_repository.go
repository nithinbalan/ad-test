package repositories

import (
	"jethur/grc/models"

	"jethur/utility"

	"gorm.io/gorm"
)

type GRCReportWorkflowRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.GRCReportWorkflow, error)
	GetByType(id uint, section string) (*models.GRCReportWorkflow, error)
	Update(grcReportWorkflow *models.GRCReportWorkflow, tx ...utility.DatabaseTransaction) error
}

type GRCReportWorkflowRepository struct {
	db *gorm.DB
}

func NewGRCReportWorkflowRepository(db *gorm.DB) GRCReportWorkflowRepositoryInterface {
	return &GRCReportWorkflowRepository{db: db}
}

func (repo *GRCReportWorkflowRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *GRCReportWorkflowRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.GRCReportWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow *models.GRCReportWorkflow

	result := db.
		Preload("GRCReport").
		Preload("GRCReportWorkflowStatus").
		Preload("ReviewUser").
		Preload("ApproveUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&workflow, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflow, nil
}

func (repo *GRCReportWorkflowRepository) GetByType(id uint, section string) (*models.GRCReportWorkflow, error) {
	var grcReportWorkflow *models.GRCReportWorkflow

	result := repo.db.Where("type = ?", section).Where("grc_report_id = ?", id).First(&grcReportWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return grcReportWorkflow, nil
}

func (repo *GRCReportWorkflowRepository) Update(grcReportWorkflow *models.GRCReportWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(grcReportWorkflow).Error; err != nil {
		return err
	}

	return nil
}
