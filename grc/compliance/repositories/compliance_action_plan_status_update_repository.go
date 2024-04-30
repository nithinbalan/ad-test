package repositories

import (
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type ComplianceActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error)
	Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	SaveComplianceActionPlanStatusUpdateDocument(doc *models.ComplianceActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetComplianceActionPlanStatusUpdateDocumentByID(id uint) (*models.ComplianceActionPlanStatusUpdateDocument, error)
	StoreComplianceActionPlanStatusUpdateDocumentFile(complianceId uint, actionId uint, updatesId uint, fileName string) error
	RemoveComplianceActionPlanStatusUpdateDocumentFile(complianceId uint, actionId uint, updatesId uint, fileNames string) error
	DeleteComplianceActionPlanStatusUpdateDocumentId(docId uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceActionPlanStatusUpdateRepository(db *gorm.DB) *ComplianceActionPlanStatusUpdateRepository {
	return &ComplianceActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetAll(complianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {

	var complianceActionPlanUpdates []*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO

	query := repo.db.Table("compliance_action_plan_status_updates").
		Select(`
            compliance_action_plan_status_updates.*,
            compliance_action_plans.title AS compliance_action_plan_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
            compliance_statuses.color_code AS compliance_status_color_code,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar 
        `).
		Joins(`
            LEFT JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
            LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_action_plan_status_updates.compliance_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_action_plan_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_action_plan_status_updates.updated_by
        `).
		Where("compliance_action_plan_status_updates.compliance_action_plan_id = ?", complianceActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (compliance_action_plans.title LIKE ? OR
            compliance_action_plan_status_updates.comment LIKE ? OR
            compliance_statuses.title_en LIKE ? OR
            compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceActionPlanUpdates)

	return complianceActionPlanUpdates, result.Error
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetByID(id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	var complianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate
	result := repo.db.First(&complianceActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlanStatusUpdate, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Create(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceActionPlanStatusUpdate.ID, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) Update(ComplianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ComplianceActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) GetComplianceActionPlanStatusUpdateDocumentByID(id uint) (*models.ComplianceActionPlanStatusUpdateDocument, error) {
	var document models.ComplianceActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) StoreComplianceActionPlanStatusUpdateDocumentFile(id uint, actionId uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	temporaryStoragePath := filepath.Join(storagePath, "tmp", fileName)
	directoryPath := filepath.Join(storagePath, "compliance", "compliance", strconv.FormatUint(uint64(id), 10), "action-plans", strconv.FormatUint(uint64(actionId), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents")
	destinationPath := filepath.Join(directoryPath, fileName)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(temporaryStoragePath, destinationPath); err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) SaveComplianceActionPlanStatusUpdateDocument(doc *models.ComplianceActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) RemoveComplianceActionPlanStatusUpdateDocumentFile(complianceId uint, actionId uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(storagePath, "compliance", "compliance", strconv.FormatUint(uint64(complianceId), 10), "action-plans", strconv.FormatUint(uint64(actionId), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents", fileName)
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanStatusUpdateRepository) DeleteComplianceActionPlanStatusUpdateDocumentId(docID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.ComplianceActionPlanStatusUpdateDocument{}, docID).Error; err != nil {
		return err
	}

	return nil
}
