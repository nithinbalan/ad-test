// Repository - repositories/policy_exception_repository.go
package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type PolicyExceptionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.PolicyExceptionListQueryParamsDTO) (*dtos.PolicyExceptionListResponseDTO, error)
	GetByID(id uint) (*models.PolicyException, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.PolicyException, error)
	Create(pe *models.PolicyException, tx ...utility.DatabaseTransaction) (uint, error)
	Update(pe *models.PolicyException, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	GetByRiskID(riskID uint) (*models.RiskRegister, error)

	CreatePolicyExceptionRisks(mappings []models.PolicyExceptionRisk, tx utility.DatabaseTransaction) error
	DeleteSinglePolicyExceptionRisk(policyExceptionID uint, riskID uint, tx utility.DatabaseTransaction) error
}

type PolicyExceptionRepository struct {
	db *gorm.DB
}

func NewPolicyExceptionRepository(db *gorm.DB) *PolicyExceptionRepository {
	return &PolicyExceptionRepository{db: db}
}

func (repo *PolicyExceptionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PolicyExceptionRepository) GetAll(params *dtos.PolicyExceptionListQueryParamsDTO) (*dtos.PolicyExceptionListResponseDTO, error) {
	var policyExceptions []*dtos.PolicyExceptionListDTO
	var totalCount int64

	query := repo.db.Table("policy_exceptions").
		Select(`
		policy_exceptions.*,
		documents.title_en AS document_title_en,
		documents.title_ar AS document_title_ar,
		document_requirement_groups.title_en AS document_requirement_group_title_en,
		document_requirement_groups.title_ar AS document_requirement_group_title_ar,
		document_requirements.title_en AS document_requirement_title_en,
		document_requirements.title_ar AS document_requirement_title_ar,
		workflow_statuses.title_en AS workflow_status_title_en,
		workflow_statuses.title_ar AS workflow_status_title_ar,
		workflow_statuses.color_code AS workflow_status_color,
		requested_user.first_name_en AS requested_by_first_name_en,
		requested_user.first_name_ar AS requested_by_first_name_ar,
		requested_user.last_name_en AS requested_by_last_name_en,
		requested_user.last_name_ar AS requested_by_last_name_ar,
		requested_user.initial AS requested_by_initial,
		requested_user.email AS requested_by_email,
		requested_user.profile_image_path AS requested_by_profile_image_path,
		created_user.first_name_en AS created_by_first_name_en,
		created_user.first_name_ar AS created_by_first_name_ar,
		created_user.last_name_en AS created_by_last_name_en,
		created_user.last_name_ar AS created_by_last_name_ar,
		created_user.initial AS created_by_initial,
		created_user.color AS created_by_color,
		created_user.email AS created_by_email,
		created_user.profile_image_path AS created_by_profile_image_path,
		updated_user.first_name_en AS updated_by_first_name_en,
		updated_user.first_name_ar AS updated_by_first_name_ar,
		updated_user.last_name_en AS updated_by_last_name_en,
		updated_user.last_name_ar AS updated_by_last_name_ar,
		updated_user.initial AS updated_by_initial,
		updated_user.color AS updated_by_color,
		updated_user.email AS updated_by_email,
		updated_user.profile_image_path AS updated_by_profile_image_path
		`).
		Joins(`LEFT JOIN documents ON documents.id = policy_exceptions.document_id`).
		Joins(`LEFT JOIN document_versions ON document_versions.id = policy_exceptions.document_version_id`).
		Joins(`LEFT JOIN document_requirement_groups ON document_requirement_groups.id = policy_exceptions.document_requirement_group_id`).
		Joins(`LEFT JOIN document_requirements ON document_requirements.id = policy_exceptions.document_requirement_id`).
		Joins(`LEFT JOIN workflow_statuses ON workflow_statuses.id = policy_exceptions.workflow_status_id`).
		Joins(`LEFT JOIN users AS requested_user ON requested_user.id = policy_exceptions.requested_by_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = policy_exceptions.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = policy_exceptions.updated_by_id`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("policy_exceptions.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			policy_exceptions.reference_number LIKE ? OR
			documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			requested_user.first_name_en LIKE ? OR
			requested_user.first_name_ar LIKE ? OR
			requested_user.last_name_en LIKE ? OR
			requested_user.last_name_ar LIKE ? OR
			document_versions.title_en LIKE ? OR
			document_versions.title_ar LIKE ? OR
			document_requirements.title_en LIKE ? OR
			document_requirements.title_ar LIKE ? OR
			document_requirement_groups.title_en LIKE ? OR
			document_requirement_groups.title_ar LIKE ?`,
			kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw,
		)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&policyExceptions).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.PolicyExceptionListResponseDTO{
		Items: policyExceptions,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
func (repo *PolicyExceptionRepository) GetByID(id uint) (*models.PolicyException, error) {
	var pe models.PolicyException
	err := repo.db.First(&pe, id).Error
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

func (repo *PolicyExceptionRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.PolicyException, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var pe *models.PolicyException
	err := db.
		Preload("Document").
		Preload("Document.LatestDocumentVersion").
		Preload("Document.LatestDocumentVersion.DocumentClauses").
		Preload("DocumentRequirementGroup").
		Preload("DocumentRequirementGroup.DocumentRequirements").
		Preload("DocumentRequirement").
		Preload("RequestedBy").
		Preload("Risks").
		Preload("Risks.RiskRegister").
		Preload("OrganizationUnit").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&pe, id).Error
	if err != nil {
		return nil, err
	}
	return pe, nil
}

func (repo *PolicyExceptionRepository) Create(pe *models.PolicyException, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(pe).Error; err != nil {
		return 0, err
	}
	return pe.ID, nil
}

func (repo *PolicyExceptionRepository) Update(pe *models.PolicyException, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(pe).Error
}

func (repo *PolicyExceptionRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.PolicyException{}, id).Error
}

func (repo *PolicyExceptionRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.PolicyException{}).
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

	newReferenceNumber := fmt.Sprintf("PE-%d-%04d", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *PolicyExceptionRepository) GetByRiskID(riskID uint) (*models.RiskRegister, error) {
	var riskRegister *models.RiskRegister
	err := repo.db.Where("id = ?", riskID).Find(&riskRegister).Error
	return riskRegister, err
}

func (repo *PolicyExceptionRepository) CreatePolicyExceptionRisks(mappings []models.PolicyExceptionRisk, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.Create(&mappings).Error
}

func (repo *PolicyExceptionRepository) DeleteSinglePolicyExceptionRisk(policyExceptionID uint, riskID uint, tx utility.DatabaseTransaction) error {
	db := tx.DB()
	return db.
		Where("policy_exception_id = ? AND risk_register_id = ?", policyExceptionID, riskID).
		Delete(&models.PolicyExceptionRisk{}).Error
}
