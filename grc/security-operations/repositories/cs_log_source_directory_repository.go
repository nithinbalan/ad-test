package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CSLogSourceDirectoryRepositoryInterface interface {
	GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error)
	Create(entry *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entry *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	BeginTransaction() utility.DatabaseTransaction
	ExistsDataEntityType(id uint) (bool, error)
	ExistsOrganizationUnit(id uint) (bool, error)
	ExistsReceptionMechanism(id uint) (bool, error)
	ExistsCommunicationChannel(id uint) (bool, error)

	CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx ...utility.DatabaseTransaction) error
	DeleteReceptionMechanisms(sourceDirectoryID uint, tx ...utility.DatabaseTransaction) error
	CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx ...utility.DatabaseTransaction) error
	DeleteCommunicationChannels(sourceDirectoryID uint, tx ...utility.DatabaseTransaction) error
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

func (repo *CSLogSourceDirectoryRepository) exists(model any, id uint) (bool, error) {
	var count int64
	err := repo.db.Model(model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
func (repo *CSLogSourceDirectoryRepository) ExistsDataEntityType(id uint) (bool, error) {
	return repo.exists(&models.CSLogDataEntityType{}, id)
}
func (repo *CSLogSourceDirectoryRepository) ExistsOrganizationUnit(id uint) (bool, error) {
	return repo.exists(&models.OrganizationUnit{}, id)
}
func (repo *CSLogSourceDirectoryRepository) ExistsReceptionMechanism(id uint) (bool, error) {
	return repo.exists(&models.CSLogDataReceptionMechanism{}, id)
}
func (repo *CSLogSourceDirectoryRepository) ExistsCommunicationChannel(id uint) (bool, error) {
	return repo.exists(&models.CSLogCommunicationChannel{}, id)
}

func (repo *CSLogSourceDirectoryRepository) GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error) {
	var entries []*socDtos.CSLogSourceDirectoryListDTO
	query := repo.db.Table("cs_log_source_directory AS sd").
		Select(`
			sd.id,
			sd.reference_number,
			sd.entity_name,
			sd.data_entity_type_id,
			data_entity_types.title_en AS data_entity_type_title_en,
			data_entity_types.title_ar AS data_entity_type_title_ar,
			sd.responsible_entity_id,
			responsible_entity.title_en AS responsible_entity_title_en,
			responsible_entity.title_ar AS responsible_entity_title_ar,
			sd.beneficiary_entity_id,
			beneficiary_entity.title_en AS beneficiary_entity_title_en,
			beneficiary_entity.title_ar AS beneficiary_entity_title_ar,
			sd.created_at,
			sd.updated_at,
			created_user.id AS created_by_id,
			created_user.first_name_en AS created_by_first_name_en,
			created_user.first_name_ar AS created_by_first_name_ar,
			created_user.last_name_en AS created_by_last_name_en,
			created_user.last_name_ar AS created_by_last_name_ar,
			created_user.email AS created_by_email,
			updated_user.id AS updated_by_id,
			updated_user.first_name_en AS updated_by_first_name_en,
			updated_user.first_name_ar AS updated_by_first_name_ar,
			updated_user.last_name_en AS updated_by_last_name_en,
			updated_user.last_name_ar AS updated_by_last_name_ar,
			updated_user.email AS updated_by_email
		`).
		Joins("LEFT JOIN cs_log_data_entity_types AS data_entity_types ON data_entity_types.id = sd.data_entity_type_id").
		Joins("LEFT JOIN organization_units AS responsible_entity ON responsible_entity.id = sd.responsible_entity_id").
		Joins("LEFT JOIN organization_units AS beneficiary_entity ON beneficiary_entity.id = sd.beneficiary_entity_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = sd.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = sd.updated_by_id")

	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where("(sd.entity_name LIKE ? OR sd.reference_number LIKE ?)", like, like)
	}
	if params.DataEntityTypeID != nil {
		query = query.Where("sd.data_entity_type_id = ?", *params.DataEntityTypeID)
	}
	if params.ResponsibleEntityID != nil {
		query = query.Where("sd.responsible_entity_id = ?", *params.ResponsibleEntityID)
	}
	if params.BeneficiaryEntityID != nil {
		query = query.Where("sd.beneficiary_entity_id = ?", *params.BeneficiaryEntityID)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Scan(&entries).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	return &socDtos.CSLogSourceDirectoryListResponseDTO{
		Items: entries,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *CSLogSourceDirectoryRepository) GetByID(id uint) (*models.CSLogSourceDirectory, error) {
	var entry models.CSLogSourceDirectory
	if err := repo.db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (repo *CSLogSourceDirectoryRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var entry *models.CSLogSourceDirectory
	result := db.
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

func (repo *CSLogSourceDirectoryRepository) Create(entry *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&entry).Error; err != nil {
		return 0, err
	}
	return entry.ID, nil
}

func (repo *CSLogSourceDirectoryRepository) Update(updated *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(updated).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSourceDirectoryRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.CSLogSourceDirectory{}, id).Error
}

func (repo *CSLogSourceDirectoryRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.CSLogSourceDirectory{}).
		Select("reference_number").
		Order("reference_number desc").
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
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SRC-%d-%04d", time.Now().Year(), lastReferenceID+1), nil
}

func (repo *CSLogSourceDirectoryRepository) CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(links) == 0 {
		return nil
	}
	return db.Create(&links).Error
}

func (repo *CSLogSourceDirectoryRepository) DeleteReceptionMechanisms(sourceDirectoryID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("source_directory_id = ?", sourceDirectoryID).
		Delete(&models.CSLogSourceDirectoryReceptionMechanism{}).Error
}

func (repo *CSLogSourceDirectoryRepository) CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(links) == 0 {
		return nil
	}
	return db.Create(&links).Error
}

func (repo *CSLogSourceDirectoryRepository) DeleteCommunicationChannels(sourceDirectoryID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("source_directory_id = ?", sourceDirectoryID).
		Delete(&models.CSLogSourceDirectoryCommunicationChannel{}).Error
}
