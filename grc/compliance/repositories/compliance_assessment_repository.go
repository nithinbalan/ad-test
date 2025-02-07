package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceAssessment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error)
	Create(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	IsAssessmentExists(organizationUnitID uint, compliancRegisterID uint) (bool, error)
	IsAssessmentExistExcludingID(id uint, organizationUnitID uint, compliancRegisterID uint) (bool, error)

	CreateComplianceAssessmentRequirement(requirements []*models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentRequirements(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRepository(db *gorm.DB) *ComplianceAssessmentRepository {
	return &ComplianceAssessmentRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentRepository) GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error) {

	var complianceAssessments []*complianceDTOs.ComplianceAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("compliance_assessments").
		Select(`
            compliance_assessments.*,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
            review_user.first_name_ar AS review_user_first_name_ar,
            review_user.last_name_en AS review_user_last_name_en,
            review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
            review_user.color AS review_user_color,
            review_user.email AS review_user_email,
            review_user.profile_image_path AS review_user_profile_image_path,
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
		Joins(`
		    LEFT JOIN compliance_register ON compliance_register.id = compliance_assessments.compliance_register_id
		    LEFT JOIN organization_units ON organization_units.id = compliance_assessments.organization_unit_id
		    LEFT JOIN users AS review_user ON review_user.id = compliance_assessments.review_user_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_assessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_assessments.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			review_user.first_name_en LIKE ? OR
			review_user.first_name_ar LIKE ? OR
			review_user.last_name_en LIKE ? OR
			review_user.last_name_ar LIKE ? OR
			review_user.email LIKE ? OR
            review_user.initial LIKE ? OR
            review_user.color LIKE ? OR
            review_user.profile_image_path LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Calculate total count for pagination
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&complianceAssessments).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &complianceDTOs.ComplianceAssessmentListResponseDTO{
		Items: complianceAssessments,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentRepository) GetByID(id uint) (*models.ComplianceAssessment, error) {
	var complianceAssessment *models.ComplianceAssessment

	result := repo.db.First(&complianceAssessment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessment, nil

}

func (repo *ComplianceAssessmentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var complianceAssessment *models.ComplianceAssessment

	result := db.
		Preload("ComplianceRegister").
		Preload("OrganizationUnit").
		Preload("AssessmentRequirements").
		Preload("AssessmentRequirements.ComplianceRequirement").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceStatus").
		Preload("ReviewUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&complianceAssessment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessment, nil
}

func (repo *ComplianceAssessmentRepository) Create(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceAssessment).Error; err != nil {
		return 0, err
	}

	return complianceAssessment.ID, nil
}

func (repo *ComplianceAssessmentRepository) Update(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Compliance  Assessment
	if err := db.Delete(&models.ComplianceAssessment{}, id).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) IsAssessmentExists(organizationUnitID uint, complianceRegisterID uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("compliance_register_id = ?", complianceRegisterID)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}

func (repo *ComplianceAssessmentRepository) IsAssessmentExistExcludingID(id uint, organizationUnitID uint, complianceRegisterID uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("id <> ? AND compliance_register_id = ?", id, complianceRegisterID)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirement(requirements []*models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&requirements).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirements(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirement{}).Error; err != nil {
		return err
	}

	return nil
}
