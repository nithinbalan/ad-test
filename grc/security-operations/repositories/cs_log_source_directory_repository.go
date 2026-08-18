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
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.CSLogSourceDirectory, error)
	Create(entry *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entry *models.CSLogSourceDirectory, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx ...utility.DatabaseTransaction) error
	DeleteReceptionMechanisms(ids []*uint, tx ...utility.DatabaseTransaction) error
	CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx ...utility.DatabaseTransaction) error
	DeleteCommunicationChannels(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type CSLogSourceDirectoryRepository struct {
	db *gorm.DB
}

func NewCSLogSourceDirectoryRepository(db *gorm.DB) *CSLogSourceDirectoryRepository {
	return &CSLogSourceDirectoryRepository{
		db: db,
	}
}

func (repo *CSLogSourceDirectoryRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *CSLogSourceDirectoryRepository) GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error) {
	var entries []*socDtos.CSLogSourceDirectoryListDTO
	var totalCount int64

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
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN cs_log_data_entity_types AS data_entity_types ON data_entity_types.id = sd.data_entity_type_id").
		Joins("LEFT JOIN organization_units AS responsible_entity ON responsible_entity.id = sd.responsible_entity_id").
		Joins("LEFT JOIN organization_units AS beneficiary_entity ON beneficiary_entity.id = sd.beneficiary_entity_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = sd.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = sd.updated_by_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("(sd.entity_name LIKE ? OR sd.reference_number LIKE ?)", keyword, keyword)
	}

	filterColumns := map[string]string{
		"data_entity_type_id":   "sd.data_entity_type_id",
		"responsible_entity_id": "sd.responsible_entity_id",
		"beneficiary_entity_id": "sd.beneficiary_entity_id",
	}
	for filterName, values := range params.Filters {
		if columnName, ok := filterColumns[filterName]; ok && len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&entries).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &socDtos.CSLogSourceDirectoryListResponseDTO{
		Items: entries,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *CSLogSourceDirectoryRepository) GetByID(id uint) (*models.CSLogSourceDirectory, error) {
	var entry *models.CSLogSourceDirectory

	result := repo.db.First(&entry, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return entry, nil
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
	if err := db.Create(entry).Error; err != nil {
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
	if err := db.Delete(&models.CSLogSourceDirectory{}, id).Error; err != nil {
		return err
	}

	return nil
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
	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("SRC-%d-%04d", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *CSLogSourceDirectoryRepository) CreateReceptionMechanisms(links []*models.CSLogSourceDirectoryReceptionMechanism, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&links).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSourceDirectoryRepository) DeleteReceptionMechanisms(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", ids).Delete(&models.CSLogSourceDirectoryReceptionMechanism{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSourceDirectoryRepository) CreateCommunicationChannels(links []*models.CSLogSourceDirectoryCommunicationChannel, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&links).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSourceDirectoryRepository) DeleteCommunicationChannels(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", ids).Delete(&models.CSLogSourceDirectoryCommunicationChannel{}).Error; err != nil {
		return err
	}
	return nil
}
