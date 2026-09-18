package repositories

import (
	"errors"
	"fmt"
	"jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentEvidenceRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentRegisterID uint, keyword string) ([]*dtos.IncidentEvidenceListDTO, error)
	GetByID(id uint) (*models.IncidentEvidence, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentEvidence, error)
	Create(evidence *models.IncidentEvidence, tx ...utility.DatabaseTransaction) (uint, error)
	Update(evidence *models.IncidentEvidence) error
	CreateIncidentEvidences(evidences []*models.IncidentEvidence, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentEvidences(evidenceIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDownloadHistories(incidentEvidenceID uint) ([]*dtos.IncidentEvidenceDownloadHistoryListDTO, error)
	CreateDownloadHistory(history *models.IncidentEvidenceDownloadHistory) error
}

type IncidentEvidenceRepository struct {
	db *gorm.DB
}

func NewIncidentEvidenceRepository(db *gorm.DB) *IncidentEvidenceRepository {
	return &IncidentEvidenceRepository{db: db}
}

func (repo *IncidentEvidenceRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentEvidenceRepository) GetAll(incidentRegisterID uint, keyword string) ([]*dtos.IncidentEvidenceListDTO, error) {
	evidences := []*dtos.IncidentEvidenceListDTO{}

	query := repo.db.Table("incident_evidences").
		Select(`
			incident_evidences.id,
			incident_evidences.incident_register_id,
			incident_evidences.incident_evidence_source_id,
			incident_evidence_sources.title_en AS incident_evidence_source_title_en,
			incident_evidence_sources.title_ar AS incident_evidence_source_title_ar,
			incident_evidences.incident_evidence_type_id,
			incident_evidence_types.title_en AS incident_evidence_type_title_en,
			incident_evidence_types.title_ar AS incident_evidence_type_title_ar,
			incident_evidences.url,
			incident_evidences.name,
			incident_evidences.size,
			incident_evidences.token,
			incident_evidences.original_name,
			incident_evidences.type,
			incident_evidences.is_read_only,
			incident_evidences.is_download,
			created_by.id AS created_by_id,
			created_by.first_name_en AS created_by_first_name_en,
			created_by.first_name_ar AS created_by_first_name_ar,
			created_by.last_name_en AS created_by_last_name_en,
			created_by.last_name_ar AS created_by_last_name_ar,
			created_by.initial AS created_by_initial,
			created_by.color AS created_by_color,
			created_by.email AS created_by_email,
			created_by.profile_image_path AS created_by_profile_image_path,
			incident_evidences.created_at,
			updated_by.id AS updated_by_id,
			updated_by.first_name_en AS updated_by_first_name_en,
			updated_by.first_name_ar AS updated_by_first_name_ar,
			updated_by.last_name_en AS updated_by_last_name_en,
			updated_by.last_name_ar AS updated_by_last_name_ar,
			updated_by.initial AS updated_by_initial,
			updated_by.color AS updated_by_color,
			updated_by.email AS updated_by_email,
			updated_by.profile_image_path AS updated_by_profile_image_path,
			incident_evidences.updated_at
		`).
		Joins("LEFT JOIN incident_evidence_sources ON incident_evidence_sources.id = incident_evidences.incident_evidence_source_id").
		Joins("LEFT JOIN incident_evidence_types ON incident_evidence_types.id = incident_evidences.incident_evidence_type_id").
		Joins("LEFT JOIN users AS created_by ON created_by.id = incident_evidences.created_by_id").
		Joins("LEFT JOIN users AS updated_by ON updated_by.id = incident_evidences.updated_by_id").
		Order("incident_evidences.created_at DESC")

	if keyword != "" {
		query = query.Where(`
			incident_evidences.incident_register_id = ? AND 
			(incident_evidence_sources.title_en LIKE ? OR incident_evidence_sources.title_ar LIKE ? OR incident_evidence_types.title_en LIKE ? OR incident_evidence_types.title_ar LIKE ?)`,
			incidentRegisterID, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	} else {
		query = query.Where("incident_evidences.incident_register_id = ?", incidentRegisterID)
	}

	err := query.Find(&evidences).Error
	return evidences, err
}

func (repo *IncidentEvidenceRepository) GetByID(id uint) (*models.IncidentEvidence, error) {
	var evidence *models.IncidentEvidence
	result := repo.db.First(&evidence, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return evidence, nil
}

func (repo *IncidentEvidenceRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentEvidence, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var evidence *models.IncidentEvidence
	err := db.Preload("IncidentRegister").
		Preload("IncidentEvidenceSource").
		Preload("IncidentEvidenceType").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ?", id).
		First(&evidence).Error

	if err != nil {
		fmt.Printf("Error fetching full details for IncidentEvidence id %d: %v\n", id, err)
		return nil, err
	}
	return evidence, nil
}

func (repo *IncidentEvidenceRepository) Create(evidence *models.IncidentEvidence, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(evidence).Error; err != nil {
		return 0, err
	}
	return evidence.ID, nil
}

func (repo *IncidentEvidenceRepository) CreateIncidentEvidences(evidences []*models.IncidentEvidence, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Insert one row at a time so each evidence's generated ID is returned correctly.
	// SQL Server's driver does not reliably populate IDs for multi-row batch inserts.
	for _, evidence := range evidences {
		if err := db.Create(evidence).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo *IncidentEvidenceRepository) DeleteIncidentEvidences(evidenceIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("incident_evidence_id IN (?)", evidenceIds).Delete(&models.IncidentEvidenceDownloadHistory{}).Error; err != nil {
		return err
	}
	if err := db.Where("id IN (?)", evidenceIds).Delete(&models.IncidentEvidence{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentEvidenceRepository) Update(evidence *models.IncidentEvidence) error {
	if err := repo.db.Save(evidence).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentEvidenceRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var evidence *models.IncidentEvidence
	if err := db.First(&evidence, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	if err := db.Where("incident_evidence_id = ?", id).Delete(&models.IncidentEvidenceDownloadHistory{}).Error; err != nil {
		return err
	}
	return db.Delete(&models.IncidentEvidence{}, id).Error
}

func (repo *IncidentEvidenceRepository) GetDownloadHistories(incidentEvidenceID uint) ([]*dtos.IncidentEvidenceDownloadHistoryListDTO, error) {
	histories := []*dtos.IncidentEvidenceDownloadHistoryListDTO{}

	query := repo.db.Table("incident_evidence_download_histories").
		Select(`
			incident_evidence_download_histories.id,
			incident_evidence_download_histories.incident_evidence_id,
			incident_evidence_download_histories.downloaded_at,
			downloaded_by.id AS downloaded_by_id,
			downloaded_by.first_name_en AS downloaded_by_first_name_en,
			downloaded_by.first_name_ar AS downloaded_by_first_name_ar,
			downloaded_by.last_name_en AS downloaded_by_last_name_en,
			downloaded_by.last_name_ar AS downloaded_by_last_name_ar,
			downloaded_by.initial AS downloaded_by_initial,
			downloaded_by.color AS downloaded_by_color,
			downloaded_by.email AS downloaded_by_email,
			downloaded_by.profile_image_path AS downloaded_by_profile_image_path
		`).
		Joins("LEFT JOIN users AS downloaded_by ON downloaded_by.id = incident_evidence_download_histories.downloaded_by").
		Where("incident_evidence_download_histories.incident_evidence_id = ?", incidentEvidenceID).
		Order("incident_evidence_download_histories.downloaded_at DESC")

	err := query.Find(&histories).Error
	return histories, err
}

func (repo *IncidentEvidenceRepository) CreateDownloadHistory(history *models.IncidentEvidenceDownloadHistory) error {
	return repo.db.Create(history).Error
}
