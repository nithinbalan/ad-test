package repositories

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMManualProcedureFormRepositoryInterface interface {
	GetAll(params *bcmDTOs.BCMManualProcedureFormListQueryParamsDTO) ([]*bcmDTOs.BCMManualProcedureFormDTO, error)
	GetByID(id uint) (*models.BCMManualProcedureForm, error)
	GetFullDetailsByID(id uint) (*models.BCMManualProcedureForm, error)
	Create(form *models.BCMManualProcedureForm) (uint, error)
	Update(form *models.BCMManualProcedureForm) error
	Delete(id uint) error
}

// Implementation
type BCMManualProcedureFormRepository struct {
	db *gorm.DB
}

func NewBCMManualProcedureFormRepository(db *gorm.DB) *BCMManualProcedureFormRepository {
	return &BCMManualProcedureFormRepository{
		db: db,
	}
}

func (repo *BCMManualProcedureFormRepository) GetAll(params *bcmDTOs.BCMManualProcedureFormListQueryParamsDTO) ([]*bcmDTOs.BCMManualProcedureFormDTO, error) {
	var forms []*bcmDTOs.BCMManualProcedureFormDTO

	query := repo.db.Table("bcm_manual_procedure_forms").
		Select(`
			bcm_manual_procedure_forms.id,
			bcm_manual_procedure_forms.title_en,
			bcm_manual_procedure_forms.title_ar,
			bcm_manual_procedure_forms.code,
			bcm_manual_procedure_forms.owner_entity_id,
			owner_unit.title_en AS owner_entity_title_en,
			owner_unit.title_ar AS owner_entity_title_ar,
			bcm_manual_procedure_forms.location_id,
			location.title_en AS location_title_en,
			location.title_ar AS location_title_ar,
			bcm_manual_procedure_forms.internal_stakeholder,
			bcm_manual_procedure_forms.notes,
			bcm_manual_procedure_forms.created_at,
			bcm_manual_procedure_forms.updated_at
		`).
		Joins("LEFT JOIN organization_units AS owner_unit ON owner_unit.id = bcm_manual_procedure_forms.owner_entity_id").
		Joins("LEFT JOIN locations AS location ON location.id = bcm_manual_procedure_forms.location_id")

	if params.Keyword != "" {
		query = query.Where(`
				(bcm_manual_procedure_forms.title_en LIKE ? OR
				bcm_manual_procedure_forms.title_ar LIKE ? OR
				bcm_manual_procedure_forms.code LIKE ? OR
				bcm_manual_procedure_forms.notes LIKE ? OR
				bcm_manual_procedure_forms.internal_stakeholder LIKE ? OR
				owner_unit.title_en LIKE ? OR
				owner_unit.title_ar LIKE ? OR
				location.title_en LIKE ? OR
				location.title_ar LIKE ?)`,
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

	query = query.Order("bcm_manual_procedure_forms." + params.SortBy + " " + params.Sort)

	result := query.Find(&forms)

	return forms, result.Error
}

func (repo *BCMManualProcedureFormRepository) GetByID(id uint) (*models.BCMManualProcedureForm, error) {
	var form *models.BCMManualProcedureForm
	result := repo.db.Preload("OwnerEntity").Preload("Location").Preload("CreatedByUser").Preload("UpdatedByUser").First(&form, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return form, nil
}

func (repo *BCMManualProcedureFormRepository) GetFullDetailsByID(id uint) (*models.BCMManualProcedureForm, error) {
	var form *models.BCMManualProcedureForm
	result := repo.db.
		Preload("OwnerEntity").
		Preload("Location").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&form, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return form, nil
}

func (repo *BCMManualProcedureFormRepository) Create(form *models.BCMManualProcedureForm) (uint, error) {
	if err := repo.db.Create(form).Error; err != nil {
		return 0, err
	}

	return form.ID, nil
}

func (repo *BCMManualProcedureFormRepository) Update(form *models.BCMManualProcedureForm) error {
	if err := repo.db.Save(form).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BCMManualProcedureFormRepository) Delete(id uint) error {
	var form *models.BCMManualProcedureForm
	if err := repo.db.First(&form, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.BCMManualProcedureForm{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
