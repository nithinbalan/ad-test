package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategicObjectiveAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	CreateStrategicObjectiveClauses(clauses []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeleteStrategicObjectiveClauses(strategicObjectiveID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error
	GetClausesByStrategicObjectiveID(strategicObjectiveID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
}

type StrategicObectiveAssociatedItemRepository struct {
	db *gorm.DB
}

func NewStrategicObectiveAssociatedItemRepository(db *gorm.DB) *StrategicObectiveAssociatedItemRepository {
	return &StrategicObectiveAssociatedItemRepository{
		db: db,
	}
}

func (repo *StrategicObectiveAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategicObectiveAssociatedItemRepository) GetAllClauses(strategicObjectiveID uint, params *commonDTOs.AssociatedItemClauseListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {

	var clauseList []commonDTOs.AssociatedClauseRowDTO
	var totalCount int64
	query := repo.db.Table("compliance_requirement_strategic_objectives").
		Select(`
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_clauses.clause_no,
			compliance_clauses.compliance_register_id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins(`
			LEFT JOIN compliance_clauses ON compliance_requirement_strategic_objectives.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON compliance_requirement_strategic_objectives.compliance_register_id = compliance_register.id
		`).
		Where("compliance_requirement_strategic_objectives.strategic_objective_id = ?", strategicObjectiveID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_clauses.clause_no LIKE ? OR
			compliance_clauses.title_en  LIKE ? OR
			compliance_clauses.title_ar  LIKE ?`,
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_register.title_en"
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&clauseList).Error; err != nil {
		return nil, err
	}

	clausesByComplianceRegister := make(map[uint]*commonDTOs.AssociatedClauseListDTO)
	for _, row := range clauseList {
		if _, exists := clausesByComplianceRegister[row.ComplianceRegisterID]; !exists {
			clausesByComplianceRegister[row.ComplianceRegisterID] = &commonDTOs.AssociatedClauseListDTO{
				ComplianceRegisterID:      row.ComplianceRegisterID,
				ComplianceRegisterTitleEn: row.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: row.ComplianceRegisterTitleAr,
				Clauses:                   []commonDTOs.AssociatedClauseDTO{},
			}
		}
		clausesByComplianceRegister[row.ComplianceRegisterID].Clauses = append(clausesByComplianceRegister[row.ComplianceRegisterID].Clauses, commonDTOs.AssociatedClauseDTO{
			ComplianceClauseID: row.ComplianceClauseID,
			ClauseNo:           row.ClauseNo,
			ClauseTitleEn:      row.ClauseTitleEn,
			ClauseTitleAr:      row.ClauseTitleAr,
		})
	}

	var groupedClauses []*commonDTOs.AssociatedClauseListDTO
	for _, v := range clausesByComplianceRegister {
		groupedClauses = append(groupedClauses, v)
	}

	if groupedClauses == nil {
		groupedClauses = []*commonDTOs.AssociatedClauseListDTO{}
	}

	totalCount = int64(len(groupedClauses))
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	start := (params.Page - 1) * params.PageSize
	if start > int(totalCount) {
		start = int(totalCount)
	}

	end := start + params.PageSize
	if end > int(totalCount) {
		end = int(totalCount)
	}

	pagedClauses := groupedClauses[start:end]
	response := &commonDTOs.AssociatedItemClauseListResponseDTO{
		Items: pagedClauses,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *StrategicObectiveAssociatedItemRepository) GetClausesByStrategicObjectiveID(strategicObjectiveID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauseIDs []uint
	if err := db.
		Model(&models.ComplianceRequirementStrategicObjective{}).
		Where("strategic_objective_id = ?", strategicObjectiveID).
		Pluck("compliance_clause_id", &clauseIDs).Error; err != nil {
		return nil, err
	}

	return clauseIDs, nil
}

func (repo *StrategicObectiveAssociatedItemRepository) CreateStrategicObjectiveClauses(clauses []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&clauses).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategicObectiveAssociatedItemRepository) DeleteStrategicObjectiveClauses(strategicObjectiveID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.
		Where("strategic_objective_id = ? AND compliance_clause_id IN ?", strategicObjectiveID, clauseIDs).
		Delete(&models.ComplianceRequirementStrategicObjective{}).Error; err != nil {
		return err
	}

	return nil
}
