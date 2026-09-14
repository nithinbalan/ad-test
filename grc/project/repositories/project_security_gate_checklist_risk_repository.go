package repositories

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistRiskRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	CreateChecklistRisk(checklistRisk *models.ProjectSecurityGateChecklistRisk, tx ...utility.DatabaseTransaction) error
	GetChecklistRiskByRiskID(checklistID uint, riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistRisk, error)

	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
	CreateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error
	GetRiskRegisterByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)

	GetChecklistResponseStatusByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error)

	// Masters
	GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error)
	GetRiskCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskCategory, error)
	GetRiskMatrixByType(matrixType string, tx ...utility.DatabaseTransaction) (*models.RiskMatrix, error)
}

type ProjectSecurityGateChecklistRiskRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistRiskRepository(db *gorm.DB) ProjectSecurityGateChecklistRiskRepositoryInterface {
	return &ProjectSecurityGateChecklistRiskRepository{db: db}
}

func (repo *ProjectSecurityGateChecklistRiskRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateChecklistRiskRepository) CreateChecklistRisk(checklistRisk *models.ProjectSecurityGateChecklistRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(checklistRisk).Error
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetChecklistRiskByRiskID(checklistID uint, riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistRisk, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var checklistRisk *models.ProjectSecurityGateChecklistRisk
	result := db.
		Where("checklist_id = ? AND risk_register_id = ?", checklistID, riskRegisterID).
		First(&checklistRisk)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, riskRegisterID)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return checklistRisk, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetChecklistResponseStatusByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistResponseStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var responseStatus *models.ProjectSecurityGateChecklistResponseStatus
	result := db.First(&responseStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return responseStatus, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var lastReferenceNumber string
	err := db.Model(&models.RiskRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Where("reference_number LIKE ?", "IR-%").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "IR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) CreateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(riskRegister).Error; err != nil {
		return 0, err
	}

	return riskRegister.ID, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) UpdateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(riskRegister).Error
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetRiskRegisterByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRegister *models.RiskRegister
	result := db.First(&riskRegister, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskRegister, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskStatus *models.RiskStatus
	result := db.Where("type = ?", status).First(&riskStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskStatus, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetRiskCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskCategory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskCategory *models.RiskCategory
	result := db.First(&riskCategory, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskCategory, nil
}

func (repo *ProjectSecurityGateChecklistRiskRepository) GetRiskMatrixByType(matrixType string, tx ...utility.DatabaseTransaction) (*models.RiskMatrix, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskMatrix *models.RiskMatrix
	result := db.Where("type = ?", matrixType).First(&riskMatrix)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, fmt.Sprintf("Risk Matrix of type %s", matrixType))
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskMatrix, nil
}
