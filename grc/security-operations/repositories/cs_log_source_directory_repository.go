package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogSourceDirectoryRepositoryInterface interface {
	GetAll(sort, sortBy, keyword string, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID *uint) ([]*models.CSLogSourceDirectory, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	Create(entry *models.CSLogSourceDirectory, tx utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error)
	Update(id uint, entry *models.CSLogSourceDirectory, tx utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error)
	Delete(id uint) error
	ExistsByEntityName(entityName string) (bool, error)
	ExistsByEntityNameExcludingID(id uint, entityName string) (bool, error)
	SetReferenceNumber(id uint, referenceNumber string, tx utility.DatabaseTransaction) error
	BeginTransaction() utility.DatabaseTransaction

	CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx utility.DatabaseTransaction) error
	DeleteReceptionMechanisms(sourceDirectoryID uint, tx utility.DatabaseTransaction) error
	CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx utility.DatabaseTransaction) error
	DeleteCommunicationChannels(sourceDirectoryID uint, tx utility.DatabaseTransaction) error
}

type CSLogSourceDirectoryRepository struct {
	db *gorm.DB
}

func NewCSLogSourceDirectoryRepository(db *gorm.DB) *CSLogSourceDirectoryRepository {
	return &CSLogSourceDirectoryRepository{db: db}
}

func (repo *CSLogSourceDirectoryRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *CSLogSourceDirectoryRepository) GetAll(sort, sortBy, keyword string, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID *uint) ([]*models.CSLogSourceDirectory, error) {
	var entries []*models.CSLogSourceDirectory
	query := repo.db.
		Preload("DataEntityType").
		Preload("ResponsibleEntity").
		Preload("BeneficiaryEntity").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ReceptionMechanisms.ReceptionMechanism").
		Preload("CommunicationChannels.CommunicationChannel")

	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("(entity_name LIKE ? OR reference_number LIKE ?)", like, like)
	}
	if dataEntityTypeID != nil {
		query = query.Where("data_entity_type_id = ?", *dataEntityTypeID)
	}
	if responsibleEntityID != nil {
		query = query.Where("responsible_entity_id = ?", *responsibleEntityID)
	}
	if beneficiaryEntityID != nil {
		query = query.Where("beneficiary_entity_id = ?", *beneficiaryEntityID)
	}

	if sortBy == "" {
		sortBy = "id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	query = query.Order(sortBy + " " + sort)

	result := query.Find(&entries)
	return entries, result.Error
}

func (repo *CSLogSourceDirectoryRepository) GetByID(id uint) (*models.CSLogSourceDirectory, error) {
	var entry *models.CSLogSourceDirectory
	result := repo.db.
		Preload("DataEntityType").
		Preload("ResponsibleEntity").
		Preload("BeneficiaryEntity").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ReceptionMechanisms.ReceptionMechanism").
		Preload("CommunicationChannels.CommunicationChannel").
		First(&entry, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return entry, nil
}

func (repo *CSLogSourceDirectoryRepository) Create(entry *models.CSLogSourceDirectory, tx utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.Create(&entry).Error; err != nil {
		return nil, err
	}
	return entry, nil
}

func (repo *CSLogSourceDirectoryRepository) Update(id uint, updated *models.CSLogSourceDirectory, tx utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	var existing *models.CSLogSourceDirectory
	if err := db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	if err := db.Model(&existing).Updates(updated).Error; err != nil {
		return nil, err
	}
	return existing, nil
}

func (repo *CSLogSourceDirectoryRepository) Delete(id uint) error {
	var entry *models.CSLogSourceDirectory
	if err := repo.db.First(&entry, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.CSLogSourceDirectory{}, id).Error
}

func (repo *CSLogSourceDirectoryRepository) ExistsByEntityName(entityName string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogSourceDirectory{}).
		Where("entity_name = ?", entityName).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CSLogSourceDirectoryRepository) ExistsByEntityNameExcludingID(id uint, entityName string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogSourceDirectory{}).
		Where("id <> ? AND entity_name = ?", id, entityName).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CSLogSourceDirectoryRepository) SetReferenceNumber(id uint, referenceNumber string, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Model(&models.CSLogSourceDirectory{}).
		Where("id = ?", id).
		Update("reference_number", referenceNumber).Error
}

func (repo *CSLogSourceDirectoryRepository) CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if len(links) == 0 {
		return nil
	}
	return db.Create(&links).Error
}

func (repo *CSLogSourceDirectoryRepository) DeleteReceptionMechanisms(sourceDirectoryID uint, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("source_directory_id = ?", sourceDirectoryID).
		Delete(&models.CSLogSourceDirectoryReceptionMechanism{}).Error
}

func (repo *CSLogSourceDirectoryRepository) CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if len(links) == 0 {
		return nil
	}
	return db.Create(&links).Error
}

func (repo *CSLogSourceDirectoryRepository) DeleteCommunicationChannels(sourceDirectoryID uint, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("source_directory_id = ?", sourceDirectoryID).
		Delete(&models.CSLogSourceDirectoryCommunicationChannel{}).Error
}
