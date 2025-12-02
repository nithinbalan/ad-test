package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

type TenableSettingsRepository interface {
	GetFirst() (*models.TenableSettings, error)
	GetByID(id uint) (*models.TenableSettings, error)
	Create(settings *models.TenableSettings) error
	Update(settings *models.TenableSettings) error
	Activate(id uint) error
	GetSettingsDTO() (*dtos.TenableSettingsResponseDTO, error)
}

type tenableSettingsRepo struct {
	db *gorm.DB
}

func NewTenableSettingsRepository(db *gorm.DB) TenableSettingsRepository {
	return &tenableSettingsRepo{db: db}
}

func (repo *tenableSettingsRepo) GetFirst() (*models.TenableSettings, error) {
	var settings models.TenableSettings
	err := repo.db.First(&settings).Error
	return &settings, err
}

func (repo *tenableSettingsRepo) GetByID(id uint) (*models.TenableSettings, error) {
	var st models.TenableSettings
	err := repo.db.First(&st, id).Error
	return &st, err
}

func (repo *tenableSettingsRepo) Create(settings *models.TenableSettings) error {
	return repo.db.Create(settings).Error
}

func (repo *tenableSettingsRepo) Update(settings *models.TenableSettings) error {
	return repo.db.Save(settings).Error
}

func (repo *tenableSettingsRepo) GetSettingsDTO() (*dtos.TenableSettingsResponseDTO, error) {
	var st models.TenableSettings
	if err := repo.db.First(&st).Error; err != nil {
		return nil, err
	}

	dto := &dtos.TenableSettingsResponseDTO{
		ID:                 st.ID,
		BaseURL:            st.BaseURL,
		AccessKey:          st.AccessKey,
		SecretKey:          st.SecretKey,
		SyncFrequencyHours: st.SyncFrequencyHours,
		LookbackDays:       st.LookbackDays,
		Status:             st.Status,
	}

	return dto, nil
}

func (repo *tenableSettingsRepo) Activate(id uint) error {
	return repo.db.
		Model(&models.TenableSettings{}).
		Where("id = ?", id).
		Update("status", "active").
		Error
}
