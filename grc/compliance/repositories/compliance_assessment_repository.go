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

	IsAssessmentExists(organizationUnitID uint, compliancRegisterID uint, ids []uint) (bool, error)
	IsAssessmentExistExcludingID(id uint, organizationUnitID uint, compliancRegisterID uint, ids []uint) (bool, error)
	HasRequirementForRegister(registerID uint) (bool, error)

	CreateComplianceAssessmentRequirement(requirements *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteComplianceAssessmentRequirements(ids []*uint, tx ...utility.DatabaseTransaction) error
	CreateComplianceAssessmentRequirementDocumentTitle(requirementDocuments []*models.ComplianceAssessmentRequirementDocumentTitle, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentRequirementDocumentTitles(ids []*uint, tx ...utility.DatabaseTransaction) error

	CreateComplianceAssessmentRequirementResponsibleUnit(resUnits *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error
	CreateComplianceAssessmentRequirementResponsibleUnits(resUnits []*models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentRequirementResponsibleUnits(ids []*uint, tx ...utility.DatabaseTransaction) error

	GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error)

	AssessmentCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error)
	GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error)

	UpdateComplianceRequirementResponsibleUintValue(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error
	UpdateComplianceRequirementStatusUpdate(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error
	UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error
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
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
            review_user.first_name_ar AS review_user_first_name_ar,
            review_user.last_name_en AS review_user_last_name_en,
            review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
            review_user.color AS review_user_color,
            review_user.email AS review_user_email,
            review_user.profile_image_path AS review_user_profile_image_path,
			compliance_assessment_statuses.title_en AS compliance_assessment_status_title_en,
			compliance_assessment_statuses.title_ar AS compliance_assessment_status_title_ar,
			compliance_assessment_statuses.color_code AS compliance_assessment_status_color_code,
			compliance_assessment_statuses.type AS compliance_assessment_status_type,
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
		    LEFT JOIN users AS review_user ON review_user.id = compliance_assessments.review_user_id
			LEFT JOIN compliance_assessment_statuses ON compliance_assessment_statuses.id = compliance_assessments.compliance_assessment_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_assessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_assessments.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "draft":
			query = query.Where("compliance_assessment_statuses.type = ?", "draft")
		case "in_review":
			query = query.Where("compliance_assessment_statuses.type = ?", "in_review")
		case "approved":
			query = query.Where("compliance_assessment_statuses.type = ?", "approved")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
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
			keyword, keyword, keyword)
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
		Preload("ComplianceRegister.ComplianceAssessmentFramework").
		Preload("ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceRegister.ComplianceAssessmentFramework.Values").
		Preload("ComplianceRegister.ResponsibleUser").
		Preload("AssessmentRequirements").
		Preload("AssessmentRequirements.ComplianceStatus").
		Preload("AssessmentRequirements.ResponsibleUnits").
		Preload("AssessmentRequirements.ResponsibleUnits.ResponsibleUnit").
		Preload("AssessmentRequirements.ResponsibleUnits.RequiredDocumentTitles").
		Preload("AssessmentRequirements.ResponsibleUnits.ComplianceAssessmentRequirementStatus").
		Preload("AssessmentRequirements.ResponsibleUnits.ComplianceAssessmentFrameworkValue").
		Preload("AssessmentRequirements.ComplianceRequirement").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceStatus").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceClause").
		Preload("ReviewUser").
		Preload("ComplianceAssessmentStatus").
		Preload("SubmittedByUser").
		Preload("ReviewedByUser").
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
func (repo *ComplianceAssessmentRepository) IsAssessmentExists(organizationUnitID uint, complianceRegisterID uint, ids []uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("compliance_register_id = ?", complianceRegisterID).
		Where("compliance_assessment_status_id IN ?", ids)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}
func (repo *ComplianceAssessmentRepository) HasRequirementForRegister(registerID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", registerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (repo *ComplianceAssessmentRepository) IsAssessmentExistExcludingID(id uint, organizationUnitID uint, complianceRegisterID uint, ids []uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("id <> ? AND compliance_register_id = ?", id, complianceRegisterID).
		Where("compliance_assessment_status_id IN ?", ids)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirement(requirements *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&requirements).Error; err != nil {
		return 0, err
	}
	return requirements.ID, nil
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

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementDocumentTitle(requirementDocuments []*models.ComplianceAssessmentRequirementDocumentTitle, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&requirementDocuments).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirementDocumentTitles(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirementDocumentTitle{}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) AssessmentCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error) {
	var out complianceDTOs.AssessmentCompletionRollupDTO

	const q = `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) AS completed,
			SUM(CASE WHEN carru.is_completed = 0 THEN 1 ELSE 0 END) AS not_completed,
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE CAST(SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) * 100.0 
				          / COUNT(*) AS DECIMAL(5,2))
			END AS completion_percentage
		FROM compliance_assessment_requirement_responsible_units carru
		JOIN compliance_assessment_requirements car 
			ON car.id = carru.compliance_assessment_requirement_id
		WHERE car.compliance_assessment_id = ?
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error) {
	var out []complianceDTOs.ComplianceStatusCountDTO

	const q = `
		SELECT
			cs.id   AS compliance_status_id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code,
			COALESCE(agg.cnt, 0) AS count
		FROM compliance_statuses cs
		LEFT JOIN (
			SELECT 
				cafv.compliance_status_id AS compliance_status_id,
				COUNT(*) AS cnt
			FROM compliance_assessment_requirement_responsible_units carru
			JOIN compliance_assessment_requirements car
				ON car.id = carru.compliance_assessment_requirement_id
			LEFT JOIN compliance_assessment_framework_values cafv
				ON cafv.id = carru.compliance_assessment_framework_value_id
			WHERE car.compliance_assessment_id = ?
			GROUP BY cafv.compliance_status_id
		) agg
			ON agg.compliance_status_id = cs.id
		-- optional: only active statuses
		-- WHERE cs.status = 'active'
		ORDER BY cs.id
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
	var result complianceDTOs.ComplianceScoreSummaryDTO

	const query = `
		SELECT 
			COUNT(*) AS total_requirements, 
			COALESCE(AVG(CAST(score AS FLOAT)), 0) AS average_score
		FROM compliance_assessment_requirements 
		WHERE compliance_assessment_id = ? AND score IS NOT NULL
	`

	if err := repo.db.Raw(query, assessmentID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementResponsibleUnits(resUnits []*models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&resUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirementResponsibleUnits(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirementResponsibleUnit{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementResponsibleUnit(resUnits *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(resUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) UpdateComplianceRequirementResponsibleUintValue(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ComplianceRequirementResponsibleUnit{}).Where("responsible_unit_id = ? AND compliance_requirement_id = ?", complianceRequirementValueUpdateDTO.ResponsibleUnitID, complianceRequirementValueUpdateDTO.ComplianceRequirementID).
		Updates(map[string]interface{}{
			"compliance_assessment_framework_value_id": complianceRequirementValueUpdateDTO.ComplianceAssessmentFrameworkValueID,
			"comments":             complianceRequirementValueUpdateDTO.Comments,
			"completed_by_user_id": complianceRequirementValueUpdateDTO.CompletedByUserID,
			"completed_at":         complianceRequirementValueUpdateDTO.CompletedAt,
		}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceRequirementStatusUpdate(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.ComplianceRequirement{}).Where("id = ?", complianceRequirementValueUpdateDTO.ComplianceRequirementID).
		Update("compliance_status_id", complianceRequirementValueUpdateDTO.ComplianceStatusID).Update("score", complianceRequirementValueUpdateDTO.Score).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceRegister{}).
		Where("id = ?", complianceRegisterID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceAssessmentRequirement{}).
		Where("id = ?", assessmentReqID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRepository) UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ComplianceClause{}).
		Where("id = ?", clauseID).
		Update("compliance_status_id", statusID).
		Update("score", score).Error
}
