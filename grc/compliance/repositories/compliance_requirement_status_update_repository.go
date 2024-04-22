package repositories

import (
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type ComplianceRequirementStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error)
	GetByID(id uint) (*models.ComplianceRequirementStatusUpdate, error)
	Create(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) error
	//Document
	SaveComplianceRequirementStatusUpdateDocument(doc *models.ComplianceRequirementStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetComplianceRequirementStatusUpdateDocumentByID(id uint) (*models.ComplianceRequirementStatusUpdateDocument, error)
	StoreComplianceRequirementStatusUpdateDocumentFile(complianceId uint, requirementId uint, updatesId uint, fileName string) error
	RemoveComplianceRequirementStatusUpdateDocumentFile(complianceId uint, requirementId uint, updatesId uint, fileNames string) error
	DeleteComplianceRequirementStatusUpdateDocumentID(docId uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceRequirementStatusUpdateRepository struct {
	db *gorm.DB
}

func NewComplianceRequirementStatusUpdateRepository(db *gorm.DB) *ComplianceRequirementStatusUpdateRepository {
	return &ComplianceRequirementStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceRequirementStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetAll(complianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error) {

	var complianceRequirementStatusUpdates []*complianceDTOs.ComplianceRequirementStatusUpdateListDTO

	query := repo.db.Table("compliance_requirement_status_updates").
		Select(`
            compliance_requirement_status_updates.*,
            compliance_requirements.title AS compliance_requirement_title,
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
            LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_status_updates.compliance_requirement_id
            LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_requirement_status_updates.compliance_status_id
            LEFT JOIN users AS created_user ON created_user.id = compliance_requirement_status_updates.created_by
            LEFT JOIN users AS updated_user ON updated_user.id = compliance_requirement_status_updates.updated_by
        `).
		Where("compliance_requirement_status_updates.compliance_requirement_id = ?", complianceRequirementID)

	if keyword != "" {
		query = query.Where(`
            (compliance_requirements.title LIKE ? OR
            compliance_requirement_status_updates.comment LIKE ? OR
            compliance_statuses.title_en LIKE ? OR
            compliance_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_requirement_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceRequirementStatusUpdates)

	return complianceRequirementStatusUpdates, result.Error
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetByID(id uint) (*models.ComplianceRequirementStatusUpdate, error) {
	var complianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate
	result := repo.db.First(&complianceRequirementStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRequirementStatusUpdate, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) Create(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(ComplianceRequirementStatusUpdate).Error; err != nil {
		return 0, err
	}

	return ComplianceRequirementStatusUpdate.ID, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) Update(ComplianceRequirementStatusUpdate *models.ComplianceRequirementStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ComplianceRequirementStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) GetComplianceRequirementStatusUpdateDocumentByID(id uint) (*models.ComplianceRequirementStatusUpdateDocument, error) {
	var document models.ComplianceRequirementStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) StoreComplianceRequirementStatusUpdateDocumentFile(id uint, requirementId uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	temporaryStoragePath := filepath.Join(storagePath, "tmp", fileName)
	directoryPath := filepath.Join(storagePath, "compliance", "compliance", strconv.FormatUint(uint64(id), 10), "requirements", strconv.FormatUint(uint64(requirementId), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents")
	destinationPath := filepath.Join(directoryPath, fileName)
	if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
		return err
	}
	if err := os.Rename(temporaryStoragePath, destinationPath); err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) SaveComplianceRequirementStatusUpdateDocument(doc *models.ComplianceRequirementStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) RemoveComplianceRequirementStatusUpdateDocumentFile(complianceId uint, requirementId uint, updatesId uint, fileName string) error {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(storagePath, "compliance", "compliance", strconv.FormatUint(uint64(complianceId), 10), "requirements", strconv.FormatUint(uint64(requirementId), 10), "updates", strconv.FormatUint(uint64(updatesId), 10), "documents", fileName)
	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRequirementStatusUpdateRepository) DeleteComplianceRequirementStatusUpdateDocumentID(docID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.ComplianceRequirementStatusUpdateDocument{}, docID).Error; err != nil {
		return err
	}

	return nil
}
