package repositories

import (
	"fmt"
	"sort"
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentControlRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAllClauses(assessmentID uint, currentUserID uint) ([]*models.ComplianceClause, error)
	GetClauseControlCounts(assessmentID uint) (map[uint]complianceDTOs.ClauseControlCountDTO, error)
	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentControlListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentControl, error)
	GetFullDetailsByID(controlId uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentControl, error)
	Update(control *models.ComplianceAssessmentControl, tx ...utility.DatabaseTransaction) error

	CreateComment(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) error
	Delete(controlID uint, tx ...utility.DatabaseTransaction) error
	ComplianceAssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentControlRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentControlRepository(db *gorm.DB) *ComplianceAssessmentControlRepository {
	return &ComplianceAssessmentControlRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceAssessmentControlRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentControlRepository) GetAllClauses(assessmentID uint, currentUserID uint) ([]*models.ComplianceClause, error) {
	var clauses []*models.ComplianceClause

	err := repo.db.
		Preload("ChildClauses", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ChildClauses")
		}).
		Joins("JOIN compliance_register ON compliance_register.id = compliance_clauses.compliance_register_id").
		Joins("JOIN compliance_assessments ON compliance_assessments.compliance_register_id = compliance_register.id").
		Where("compliance_assessments.id = ?", assessmentID).
		Find(&clauses).Error

	if err != nil {
		return nil, err
	}

	// sort numerically
	sort.SliceStable(clauses, func(i, j int) bool {
		return compareClauseNumbers(clauses[i].ClauseNo, clauses[j].ClauseNo)
	})
	for _, clause := range clauses {
		sortChildClauses(clause)
	}
	return clauses, nil
}

func (repo *ComplianceAssessmentControlRepository) GetClauseControlCounts(assessmentID uint) (map[uint]complianceDTOs.ClauseControlCountDTO, error) {
	var results []complianceDTOs.ClauseControlCountDTO

	err := repo.db.Table("compliance_clauses cc").
		Select(`
			cc.id as id,
			COUNT(cac.id) as total_controls,
			SUM(CASE WHEN cac.is_completed = 1 THEN 1 ELSE 0 END) as completed_controls
		`).
		Joins("JOIN compliance_controls con ON con.compliance_clause_id = cc.id").
		Joins("JOIN compliance_assessment_controls cac ON cac.compliance_control_id = con.id").
		Joins("JOIN compliance_assessments ca ON ca.id = cac.compliance_assessment_id").
		Where("ca.id = ?", assessmentID).
		Group("cc.id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]complianceDTOs.ClauseControlCountDTO)
	for _, r := range results {
		countMap[r.ID] = r
	}
	return countMap, nil
}

func (repo *ComplianceAssessmentControlRepository) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentControlListResponseDTO, error) {
	var controls []*complianceDTOs.ComplianceAssessmentControlListDTO
	var total int64

	query := repo.db.Table("compliance_assessment_controls").
		Select(`
			compliance_assessment_controls.*,
			compliance_clauses.id AS clause_id,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title_en AS compliance_clause_title_en,
			compliance_clauses.title_ar AS compliance_clause_title_ar,
			controls.id AS control_id,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			compliance_assessment_requirement_statuses.title_en AS status_title_en,
			compliance_assessment_requirement_statuses.title_ar AS status_title_ar,
			compliance_assessment_requirement_statuses.color_code AS status_color_code,
			compliance_assessment_requirement_statuses.type AS status_type
		`).
		Joins(`
			LEFT JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
			LEFT JOIN controls ON controls.id = compliance_controls.control_id
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_controls.compliance_clause_id
			LEFT JOIN compliance_assessment_requirement_statuses ON compliance_assessment_requirement_statuses.id = compliance_assessment_controls.compliance_assessment_requirement_status_id
		`).
		Where("compliance_assessment_controls.compliance_assessment_id = ?", assessmentID)

	if params.ClauseID != 0 {
		query = query.Where("compliance_controls.compliance_clause_id = ?", params.ClauseID)
	}
	// dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "completed":
			query = query.Where("compliance_assessment_controls.is_completed = ?", true)
		case "not_completed":
			query = query.Where("compliance_assessment_controls.is_completed = ?", false)
		}
	}

	// keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_clauses.clause_no LIKE ? OR
			 compliance_clauses.title_en LIKE ? OR
			 compliance_clauses.title_ar LIKE ? OR
			 controls.title_en LIKE ? OR
			 controls.title_ar LIKE ? )
		`, keyword, keyword, keyword, keyword, keyword)
	}

	sortBy := params.SortBy
	sort := strings.ToUpper(strings.TrimSpace(params.Sort))
	if sort != "ASC" && sort != "DESC" {
		sort = "ASC"
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	query = query.Order(fmt.Sprintf("TRIM(%s) %s", sortBy, sort))

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&controls).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &complianceDTOs.ComplianceAssessmentControlListResponseDTO{
		Items: controls,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  total,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentControlRepository) GetByID(id uint) (*models.ComplianceAssessmentControl, error) {
	var control *models.ComplianceAssessmentControl
	result := repo.db.First(&control, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return control, nil
}

func (repo *ComplianceAssessmentControlRepository) GetFullDetailsByID(controlId uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentControl, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var control models.ComplianceAssessmentControl
	query := db.
		Preload("ComplianceAssessment").
		Preload("ComplianceAssessment.ComplianceRegister").
		Preload("ComplianceControl").
		Preload("ComplianceControl.Control").
		Preload("ComplianceControl.ComplianceClause").
		Preload("CompletedByUser").
		Preload("ComplianceAssessmentFrameworkValue").
		Preload("ComplianceAssessmentFrameworkValue.ComplianceStatus").
		Preload("ComplianceAssessmentRequirementStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser")

	if err := query.First(&control, controlId).Error; err != nil {
		return nil, err
	}
	return &control, nil
}

func (repo *ComplianceAssessmentControlRepository) Update(control *models.ComplianceAssessmentControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(control).Error
}

func (repo *ComplianceAssessmentControlRepository) CreateComment(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(comment).Error
}

func (repo *ComplianceAssessmentControlRepository) Delete(controlID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ComplianceAssessmentControl{}, controlID).Error
}

func (repo *ComplianceAssessmentControlRepository) ComplianceAssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessment).Error; err != nil {
		return err
	}
	return nil
}
