package repositories

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMExecutionRecordRepositoryInterface interface {
	GetAll(params *bcmDTOs.BCMExecutionRecordListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordDTO, error)
	GetByID(id uint) (*models.BCMExecutionRecord, error)
	GetFullDetailsByID(id uint) (*models.BCMExecutionRecord, error)
	Create(record *models.BCMExecutionRecord) (uint, error)
	Update(record *models.BCMExecutionRecord) error
	Delete(id uint) error
}

// Implementation
type BCMExecutionRecordRepository struct {
	db *gorm.DB
}

func NewBCMExecutionRecordRepository(db *gorm.DB) *BCMExecutionRecordRepository {
	return &BCMExecutionRecordRepository{
		db: db,
	}
}

func (repo *BCMExecutionRecordRepository) GetAll(params *bcmDTOs.BCMExecutionRecordListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordDTO, error) {
	var records []*bcmDTOs.BCMExecutionRecordDTO

	query := repo.db.Table("bcm_execution_records").
		Select(`
			bcm_execution_records.id,
			bcm_execution_records.title_en,
			bcm_execution_records.title_ar,
			bcm_execution_records.owner_entity_id,
			owner_unit.title_en AS owner_entity_title_en,
			owner_unit.title_ar AS owner_entity_title_ar,
			bcm_execution_records.record_type_id,
			record_type.title_en AS record_type_title_en,
			record_type.title_ar AS record_type_title_ar,
			bcm_execution_records.is_back_up_available,
			bcm_execution_records.location_id,
			location.title_en AS location_title_en,
			location.title_ar AS location_title_ar,
			bcm_execution_records.recovery_plan,
			bcm_execution_records.created_at,
			bcm_execution_records.updated_at
		`).
		Joins("LEFT JOIN organization_units AS owner_unit ON owner_unit.id = bcm_execution_records.owner_entity_id").
		Joins("LEFT JOIN bcm_execution_record_types AS record_type ON record_type.id = bcm_execution_records.record_type_id").
		Joins("LEFT JOIN locations AS location ON location.id = bcm_execution_records.location_id")

	if params.Keyword != "" {
		query = query.Where(`
				(bcm_execution_records.title_en LIKE ? OR
				bcm_execution_records.title_ar LIKE ? OR
				owner_unit.title_en LIKE ? OR
				owner_unit.title_ar LIKE ? OR
				record_type.title_en LIKE ? OR
				record_type.title_ar LIKE ? OR
				location.title_en LIKE ? OR
				location.title_ar LIKE ? OR
				bcm_execution_records.recovery_plan LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order("bcm_execution_records." + params.SortBy + " " + params.Sort)

	result := query.Find(&records)

	return records, result.Error
}

func (repo *BCMExecutionRecordRepository) GetByID(id uint) (*models.BCMExecutionRecord, error) {
	var record *models.BCMExecutionRecord
	result := repo.db.Preload("OwnerEntity").Preload("RecordType").Preload("Location").Preload("CreatedByUser").Preload("UpdatedByUser").First(&record, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}

func (repo *BCMExecutionRecordRepository) GetFullDetailsByID(id uint) (*models.BCMExecutionRecord, error) {
	var record *models.BCMExecutionRecord
	result := repo.db.
		Preload("OwnerEntity").
		Preload("RecordType").
		Preload("Location").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&record, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return record, nil
}

func (repo *BCMExecutionRecordRepository) Create(record *models.BCMExecutionRecord) (uint, error) {
	if err := repo.db.Create(record).Error; err != nil {
		return 0, err
	}

	return record.ID, nil
}

func (repo *BCMExecutionRecordRepository) Update(record *models.BCMExecutionRecord) error {
	if err := repo.db.Save(record).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BCMExecutionRecordRepository) Delete(id uint) error {
	var record *models.BCMExecutionRecord
	if err := repo.db.First(&record, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.BCMExecutionRecord{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
