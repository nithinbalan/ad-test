package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	Create(riskComment *models.RiskComment, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
}

type RiskCommentRepository struct {
	db *gorm.DB
}

func NewRiskCommentRepository(db *gorm.DB) RiskCommentRepositoryInterface {
	return &RiskCommentRepository{
		db: db,
	}
}

func (repo *RiskCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskCommentRepository) Create(riskComment *models.RiskComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskComment).Error; err != nil {
		return 0, err
	}
	return riskComment.ID, nil
}

func (repo *RiskCommentRepository) UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(risk).Update("risk_status_id", risk.RiskStatusID).Error; err != nil {
		return err
	}

	return nil
}
