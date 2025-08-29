package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type OrganizationStructureRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetOrganizationStructure(sort string, sortBy string) ([]models.OrganizationUnit, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, primaryOrgUnitField string) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error)
	GetOpenRiskCountByUnit(unitID uint) (int64, error)
	GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error)
	SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error
}

type OrganizationStructureRepository struct {
	db *gorm.DB
}

func NewOrganizationStructureRepository(db *gorm.DB) *OrganizationStructureRepository {
	return &OrganizationStructureRepository{
		db: db,
	}
}

func (repo *OrganizationStructureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *OrganizationStructureRepository) GetOrganizationStructure(sort string, sortBy string) ([]models.OrganizationUnit, error) {
	var units []models.OrganizationUnit

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Build the base query
	query := repo.db.Table("organization_units").
		Select("organization_units.*").
		Joins("LEFT JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id")

	// Execute the query
	query = query.Order(sortBy + " " + sort)
	err := query.Find(&units).Error

	return units, err
}

func (repo *OrganizationStructureRepository) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, primaryOrgUnitField string) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var users []organizationDTOs.OrganizationUnitUserPrimaryDetailsDTO
	var totalCount int64

	// Base query with joins
	joinClause := fmt.Sprintf("LEFT JOIN organization_units AS primary_unit ON primary_unit.id = users.%s", primaryOrgUnitField)
	query := repo.db.Table("users").
		Select(`
            users.id,
            users.first_name_en,
            users.last_name_en,
            users.first_name_ar,
            users.last_name_ar,
			users.is_head,
            designations.title_en AS designation_en,
            designations.title_ar AS designation_ar`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins(joinClause).
		Where("users."+primaryOrgUnitField+" = ?", orgUnitID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(int(offset)).Limit(int(params.PageSize))

	// Apply sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	response := &organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationStructureRepository) GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error) {
	var compliant, total int64
	var statuses []organizationDTOs.ComplianceStatusCountDTO

	// Step 1: compliant + total
	const q1 = `
		SELECT 
			COALESCE(SUM(CASE WHEN cs.type = 'compliant' THEN 1 ELSE 0 END), 0) AS compliant,
			COALESCE(COUNT(*), 0) AS total
		FROM compliance_assessment_requirement_responsible_units carru
		JOIN compliance_assessment_framework_values cafv
			ON cafv.id = carru.compliance_assessment_framework_value_id
		JOIN compliance_statuses cs
			ON cs.id = cafv.compliance_status_id
		JOIN compliance_assessment_requirement_statuses carrs
			ON carrs.id = carru.compliance_assessment_requirement_status_id
		WHERE carru.responsible_unit_id = ?
		  AND carrs.type = 'approved'
		  AND cs.type NOT IN ('exempt','not_applicable')
	`
	if err := repo.db.Raw(q1, unitID).Row().Scan(&compliant, &total); err != nil {
		return 0, 0, nil, err
	}

	// Step 2: per-status (for thresholds)
	const q2 = `
		SELECT
			cs.id   AS compliance_status_id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code,
			cs.score_from,
			cs.score_to,
			COALESCE(agg.cnt, 0) AS count
		FROM compliance_statuses cs
		LEFT JOIN (
			SELECT 
				cafv.compliance_status_id,
				COUNT(*) AS cnt
			FROM compliance_assessment_requirement_responsible_units carru
			JOIN compliance_assessment_framework_values cafv
				ON cafv.id = carru.compliance_assessment_framework_value_id
			JOIN compliance_assessment_requirement_statuses carrs
				ON carrs.id = carru.compliance_assessment_requirement_status_id
			WHERE carru.responsible_unit_id = ?
			  AND carrs.type = 'approved'
			GROUP BY cafv.compliance_status_id
		) agg
			ON agg.compliance_status_id = cs.id
		ORDER BY cs.id
	`
	if err := repo.db.Raw(q2, unitID).Scan(&statuses).Error; err != nil {
		return 0, 0, nil, err
	}

	return compliant, total, statuses, nil
}

func (repo *OrganizationStructureRepository) GetOpenRiskCountByUnit(unitID uint) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM risk_register r
		JOIN risk_statuses rs ON r.risk_status_id = rs.id
		WHERE r.responsible_unit_id = ? AND rs.type = 'open'
	`
	if err := repo.db.Raw(query, unitID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *OrganizationStructureRepository) GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error) {
	var avgScore float64
	if err := repo.db.Model(&models.Control{}).
		Select("COALESCE(AVG(efficiency_score),0)").
		Where("responsible_unit_id = ?", unitID).
		Scan(&avgScore).Error; err != nil {
		return nil, err
	}

	var status organizationDTOs.ControlEfficiencyStatusDTO
	err := repo.db.Table("control_efficiency_statuses").
		Select("id, title_en, title_ar, status, color_code, score_from, score_to").
		Where("? BETWEEN score_from AND score_to AND status = 'active'", avgScore).
		First(&status).Error
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (repo *OrganizationStructureRepository) SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.User{}).
		Where(fmt.Sprintf("%s = ? AND is_head = ?", primaryField), orgUnitID, true).
		Update("is_head", false).Error
}
