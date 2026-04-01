package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type OrganizationUnitsRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]models.OrganizationUnit, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error)
	GetOpenRiskCountByUnit(unitID uint) (int64, error)
	GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error)
	SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error
	GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllProcessesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetParentOrganizationUnitWithHead(unitID uint, primaryField string) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error)
	GetHeadUsersByOrganizationUnitID(organizationUnitID uint) ([]*models.User, error)

	// Organization Unit Report
	GetAll(params *organizationDTOs.OrganizationUnitReportListQueryParamsDTO) (*organizationDTOs.OrganizationUnitReportListReponseDTO, error)
	GetByID(id uint) (*organizationDTOs.OrganizationUnitReportListDTO, error)
	CreateOrganizationUnitReport(organizationUnitReport *models.OrganizationUnitReport, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteOrganizationUnitReport(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	GetOrganizationUnitOverview(organizationUnitID uint, level uint) (*organizationDTOs.OrganizationUnitOverviewDTO, error)
	GetOrganizationUnitKeyMetrics(organizationUnitID uint) (*organizationDTOs.OrganizationUnitKeyMatrixDTO, error)
	GetUserRoleDistribution(organizationUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserRoleDistributionDTO, error)
	GetUsersByOrganizationUnitReport(organizationUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserDetailsDTO, error)
	GetRiskByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitRiskListDTO, error)
	GetPolicyByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitPolicyListDTO, error)
	GetProcessByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitProcessListDTO, error)
	GetComplianceByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitComplianceListDTO, error)
	GetComplianceRequirementsByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitComplianceRequirementListDTO, error)

	// get summaries
	GetRiskSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitRiskSummaryDTO, error)
	GetRiskRatingDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
	GetPolicySummaryByOrganizationUnit(organizationUnitID uint, level uint) (*organizationDTOs.OrganizationUnitPolicySummaryDTO, error)
	GetComplianceSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitComplianceSummaryDTO, error)
	GetProcessHealthSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitProcessHealthSummaryDTO, error)
	GetProcessCriticalityDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
	GetRiskMitigationStatusDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)

	// create list
	CreateOrganizationUnitReportOverview(overview *models.OrganizationUnitReportOverview, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportKeyMetrics(keyMetrics *models.OrganizationUnitReportKeyMetrics, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportUserRole(userRole *models.OrganizationUnitReportUserRole, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportUser(user *models.OrganizationUnitReportUser, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportRisk(risk *models.OrganizationUnitReportRisk, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportPolicy(policy *models.OrganizationUnitReportPolicy, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportProcess(process *models.OrganizationUnitReportProcess, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportComplianceRequirement(complianceRequirement *models.OrganizationUnitReportComplianceRequirement, tx ...utility.DatabaseTransaction) error

	// create summaires
	CreateOrganizationUnitReportRiskSummary(riskSummary *models.OrganizationUnitReportRiskSummary, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportRiskRating(riskRating *models.OrganizationUnitReportRiskRating, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportRiskMitigation(riskMitigation *models.OrganizationUnitReportRiskMitigation, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportPolicySummary(policySummary *models.OrganizationUnitReportPolicySummary, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportProcessSummary(processSummary *models.OrganizationUnitReportProcessSummary, tx ...utility.DatabaseTransaction) error
	CreateOrganizationUnitReportComplianceSummary(complianceSummary *models.OrganizationUnitReportComplianceSummary, tx ...utility.DatabaseTransaction) error

	// delete lists
	DeleteOrganizationUnitReportOverview(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportKeyMetrics(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportUserRole(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportUser(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportRisk(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportPolicy(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportProcess(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportComplianceRequirement(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error

	// delete summaries
	DeleteOrganizationUnitReportRiskSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportRiskRating(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportRiskMitigation(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportPolicySummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportProcessSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteOrganizationUnitReportComplianceSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error

	// get tables
	GetOrganizationUnitReportOverview(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportOverviewResponseDTO, error)
	GetOrganizationUnitReportKeyMetrics(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportKeyMetricsResponseDTO, error)
	GetOrganizationUnitReportUserRole(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportUserRoleResponseDTO, error)
	GetOrganizationUnitReportUser(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportUserResponseDTO, error)
	GetOrganizationUnitReportRisk(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskResponseDTO, error)
	GetOrganizationUnitReportPolicy(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportPolicyResponseDTO, error)
	GetOrganizationUnitReportProcess(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportProcessResponseDTO, error)
	GetOrganizationUnitReportComplianceRequirement(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportComplianceRequirementResponseDTO, error)

	GetOrganizationUnitReportRiskSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportRiskSummaryResponseDTO, error)
	GetOrganizationUnitReportRiskRating(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskRatingResponseDTO, error)
	GetOrganizationUnitReportRiskMitigation(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskMitigationResponsDTO, error)
	GetOrganizationUnitReportPolicySummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportPolicySummaryResponseDTO, error)
	GetOrganizationUnitReportProcessSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportProcessSummaryResponseDTO, error)
	GetOrganizationUnitReportComplianceSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportComplianceSummaryResponseDTO, error)
}

type OrganizationUnitsRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitsRepository(db *gorm.DB) *OrganizationUnitsRepository {
	return &OrganizationUnitsRepository{
		db: db,
	}
}

func (repo *OrganizationUnitsRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]models.OrganizationUnit, error) {
	var units []struct {
		models.OrganizationUnit
		Level        uint   `gorm:"column:level"`
		LevelTitleEn string `gorm:"column:level_title_en"`
		LevelTitleAr string `gorm:"column:level_title_ar"`
		IsPrimary    bool   `gorm:"column:is_primary"`
	}

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Build the query with explicit joins
	query := repo.db.Table("organization_units").
		Select(`
            organization_units.*,
			organization_levels.level,
            organization_levels.title_en as level_title_en,
            organization_levels.title_ar as level_title_ar,
            organization_levels.is_primary as is_primary
        `).
		Joins("LEFT JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id").
		Order(fmt.Sprintf("%s %s", sortBy, sort))

	err := query.Scan(&units).Error
	if err != nil {
		return nil, err
	}

	// Map the results to the expected return type
	result := make([]models.OrganizationUnit, len(units))
	for i, u := range units {
		result[i] = u.OrganizationUnit
		if u.OrganizationLevel == nil {
			result[i].OrganizationLevel = &models.OrganizationLevel{}
		}
		result[i].OrganizationLevel.TitleEn = u.LevelTitleEn
		result[i].OrganizationLevel.TitleAr = u.LevelTitleAr
		result[i].OrganizationLevel.IsPrimary = u.IsPrimary
		result[i].OrganizationLevel.Level = uint8(u.Level)
	}

	return result, nil
}

func (repo *OrganizationUnitsRepository) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var users []organizationDTOs.OrganizationUnitUserPrimaryDetailsDTO
	var totalCount int64

	// Build the base query
	query := repo.db.Table("users").
		Select(`
		 users.id,
		 users.first_name_en,
		 users.last_name_en,
		 users.first_name_ar,
		 users.last_name_ar,
		 users.is_head,
		 users.initial,
		 users.email,
		 users.profile_image_path,
		 designations.title_en AS designation_en,
		 designations.title_ar AS designation_ar`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id")

	// Build the WHERE clause based on level
	if level > 0 {
		// Construct the field name using the level parameter (e.g., organization_unit_1_id)
		field := fmt.Sprintf("organization_unit_%d_id", level)
		query = query.Where(fmt.Sprintf("users.%s = ?", field), orgUnitID)
	}

	// Rest of the function remains the same...
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(int(offset)).Limit(int(params.PageSize))

	// Apply sorting
	sortOrder := "CASE WHEN users.is_head = 1 THEN 0 ELSE 1 END, " + params.SortBy + " " + params.Sort
	query = query.Order(sortOrder)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Populate job descriptions for each user
	if len(users) > 0 {
		userIDs := make([]uint, len(users))
		for i, user := range users {
			userIDs[i] = user.ID
		}

		// Get job descriptions for all users
		var jobDescriptions []struct {
			ID                       uint   `json:"id"`
			UserID                   uint   `json:"user_id"`
			JobDescriptionID         uint   `json:"job_description_id"`
			JobTitleEn               string `json:"job_title_en"`
			JobTitleAr               string `json:"job_title_ar"`
			JobDescriptionEn         string `json:"job_description_en"`
			JobDescriptionAr         string `json:"job_description_ar"`
			RolesAndResponsibilities string `json:"roles_and_responsibilities"`
		}
		jobDescQuery := repo.db.Table("user_job_descriptions").
			Select(`
				user_job_descriptions.id,
				user_job_descriptions.user_id,
				user_job_descriptions.job_description_id,
				job_descriptions.title_en AS job_title_en,
				job_descriptions.title_ar AS job_title_ar,
				job_descriptions.description_en AS job_description_en,
				job_descriptions.description_ar AS job_description_ar,
				job_descriptions.roles_and_responsibilities AS roles_and_responsibilities

			`).
			Joins("LEFT JOIN job_descriptions ON job_descriptions.id = user_job_descriptions.job_description_id").
			Where("user_job_descriptions.user_id IN (?)", userIDs)

		if err := jobDescQuery.Find(&jobDescriptions).Error; err != nil {
			return nil, err
		}

		// Group job descriptions by user ID
		jobDescMap := make(map[uint][]*organizationDTOs.JobDescriptionListDTO)
		for i := range jobDescriptions {
			userID := jobDescriptions[i].UserID
			if jobDescMap[userID] == nil {
				jobDescMap[userID] = []*organizationDTOs.JobDescriptionListDTO{}
			}
			jobDescMap[userID] = append(jobDescMap[userID], &organizationDTOs.JobDescriptionListDTO{
				ID:                       jobDescriptions[i].ID,
				JobDescriptionID:         jobDescriptions[i].JobDescriptionID,
				JobTitleEn:               jobDescriptions[i].JobTitleEn,
				JobTitleAr:               jobDescriptions[i].JobTitleAr,
				JobDescriptionEn:         jobDescriptions[i].JobDescriptionEn,
				JobDescriptionAr:         jobDescriptions[i].JobDescriptionAr,
				RolesAndResponsibilities: jobDescriptions[i].RolesAndResponsibilities,
			})
		}

		// Assign job descriptions to users
		for i := range users {
			users[i].JobDescriptions = jobDescMap[users[i].ID]
			if users[i].JobDescriptions == nil {
				users[i].JobDescriptions = []*organizationDTOs.JobDescriptionListDTO{}
			}
		}
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error) {
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

func (repo *OrganizationUnitsRepository) GetOpenRiskCountByUnit(unitID uint) (int64, error) {
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

func (repo *OrganizationUnitsRepository) GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error) {
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

func (repo *OrganizationUnitsRepository) SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.User{}).
		Where(fmt.Sprintf("%s = ? AND is_head = ?", primaryField), orgUnitID, true).
		Update("is_head", false).Error
}

func (repo *OrganizationUnitsRepository) GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("documents").
		Select(`
			documents.id AS id,
			documents.id AS document_id,
			documents.reference_number AS reference_number,
			documents.document_type_id AS document_type_id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.last_checkout_at AS last_checkout_at,
			documents.last_checkin_at AS last_checkin_at,
			documents.is_locked AS is_locked,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			documents.responsible_user_id,
			documents.created_at,
			documents.updated_at,
			documents.published_at,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_type_access_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
        	responsible_user.color AS responsible_user_color,
        	responsible_user.email AS responsible_user_email,
        	responsible_user.profile_image_path AS responsible_user_profile_image_path,
            document_versions.url AS document_url,
            document_versions.name AS document_name,
            document_versions.size AS document_size,
            document_versions.token AS document_token,
            document_versions.original_name AS document_original_name,
            document_versions.type AS document_type,
            document_versions.version AS document_version,
            document_statuses.title_en AS document_status_title_en,
            document_statuses.title_ar AS document_status_title_ar,
            document_statuses.color_code AS document_status_color_code,
            document_statuses.type AS document_status_type,
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
        	updated_user.profile_image_path AS updated_user_profile_image_path,
            published_user.first_name_en AS published_user_first_name_en,
            published_user.first_name_ar AS published_user_first_name_ar,
            published_user.last_name_en AS published_user_last_name_en,
            published_user.last_name_ar AS published_user_last_name_ar,
            published_user.initial AS published_user_initial,
            published_user.color AS published_user_color,
            published_user.email AS published_user_email,
            published_user.profile_image_path AS published_user_profile_image_path
		`).
		Joins(`
            LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
            LEFT JOIN document_types ON document_types.id = documents.document_type_id
            LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
            LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
            LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
            LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
            LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
            LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
            LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
        `).
		Where("documents.responsible_unit_id = ?", unitID).
		Where("document_types.type = 'policy'").
		Where("document_statuses.type = 'published'")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			documents.reference_number LIKE ? OR
			document_categories.title_en LIKE ? OR
			document_categories.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			document_statuses.title_en LIKE ? OR
			document_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "documents.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedDocumentListResponseDTO{
		Items: documentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationUnitsRepository) GetAllProcessesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
	var processes []*organizationDTOs.ProcessListDTO
	var totalCount int64

	query := repo.db.Table("processes").
		Select(`
			processes.id,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			processes.marked_as_auditable_item,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			processes.rto as rto,
			processes.rto_unit_id as rto_unit_id,
			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			processes.mtpd as mtpd,
			processes.mtpd_unit_id as mtpd_unit_id,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			processes.rpo as rpo,
			processes.rpo_unit_id as rpo_unit_id,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
			processes.is_bia_performed,
			processes.last_assessed_at,
			last_assessed_user.id AS last_assessed_user_id,
			last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
			last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
			last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
			last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
			last_assessed_user.initial AS last_assessed_user_initial,
			last_assessed_user.color AS last_assessed_user_color,
			last_assessed_user.email AS last_assessed_user_email,
			last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			process_criticality_levels.id AS process_criticality_level_id,
			process_criticality_levels.title_en AS process_criticality_level_title_en,
			process_criticality_levels.title_ar AS process_criticality_level_title_ar,
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
			updated_user.profile_image_path AS updated_user_profile_image_path,
			processes.created_at,
			processes.updated_at,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_rating_values.value AS risk_rating_value,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			risk_levels.type AS risk_level_type
		`).
		Joins("LEFT JOIN process_groups ON process_groups.id = processes.process_group_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id").
		Joins("LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id").
		Joins("LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = processes.rto_unit_id").
		Joins("LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = processes.mtpd_unit_id").
		Joins("LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = processes.rpo_unit_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id").
		Joins("LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id").
		Where("processes.responsible_unit_id = ?", unitID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.reference_number LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "processes.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processes).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ProcessListResponseDTO{
		Items: processes,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationUnitsRepository) GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*commonDTOs.AssociatedRiskListDTO

	query := repo.db.Table("risk_register").
		Select(`
		    risk_register.id as id,
			risk_register.id as risk_register_id,
			 risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
        	users.color AS responsible_user_color,
        	users.email AS responsible_user_email,
        	users.profile_image_path AS responsible_user_profile_image_path,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
        	updated_user.profile_image_path AS updated_user_profile_image_path,
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating_value.value AS current_risk_rating_value,
            current_risk_rating.color_code AS current_risk_rating_color_code,
            inherent_likelihood.title_en AS inherent_likelihood_title_en,
            inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
            inherent_impact.title_en AS inherent_impact_title_en,
            inherent_impact.title_ar AS inherent_impact_title_ar,
            inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
            inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
            inherent_risk_rating_value.value AS inherent_risk_rating_value,
            inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
            residual_likelihood.title_en AS residual_likelihood_title_en,
            residual_likelihood.title_ar AS residual_likelihood_title_ar,
            residual_impact.title_en AS residual_impact_title_en,
            residual_impact.title_ar AS residual_impact_title_ar,
            residual_risk_rating.title_en AS residual_risk_rating_title_en,
            residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
            residual_risk_rating_value.value AS residual_risk_rating_value,
            residual_risk_rating.color_code AS residual_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
            LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
            LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
		`).
		Where("risk_register.responsible_unit_id = ?", unitID).
		Where("risk_statuses.type = ?", "open")

	// Count total records

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	// Filter by risk rating value
	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	// Sorting defaults
	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &commonDTOs.AssociatedRiskListResponseDTO{
		Items: riskList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

// GetParentOrganizationUnitWithHead retrieves the main organization unit and its head user
func (repo *OrganizationUnitsRepository) GetParentOrganizationUnitWithHead(unitID uint, primaryField string) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error) {
	var result organizationDTOs.ParentOrganizationUnitWithHeadDTO

	// First, get the main organization unit (where organization_unit_id is null)
	parentQuery := `
        SELECT 
            COALESCE(parent_ou.id, child_ou.id) as id,
            COALESCE(parent_ou.title_en, child_ou.title_en) as title_en,
            COALESCE(parent_ou.title_ar, child_ou.title_ar) as title_ar
        FROM organization_units child_ou
        LEFT JOIN organization_units parent_ou ON parent_ou.id = child_ou.organization_unit_id
        WHERE child_ou.id = ?`

	var parentUnit organizationDTOs.OrganizationUnitDTO
	err := repo.db.Raw(parentQuery, unitID).Scan(&parentUnit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No parent unit found, return empty result
			return &result, nil
		}
		return nil, fmt.Errorf("failed to fetch parent organization unit: %w", err)
	}

	result.OrganizationUnit = parentUnit

	// If we found a parent unit, try to find its head user from primary organization unit
	if parentUnit.ID != 0 {
		headUserQuery := fmt.Sprintf(`
            SELECT
                u.id,
                u.first_name_en,
                u.last_name_en,
                u.first_name_ar,
                u.last_name_ar,
                u.profile_image_path,
                u.is_head,
				u.initial
            FROM users u
            WHERE u.%s = ? AND u.is_head = 1`, primaryField)

		var headUser organizationDTOs.OrganizationUnitUserDTO
		err = repo.db.Raw(headUserQuery, unitID).Scan(&headUser).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to fetch head user: %w", err)
		}

		// Only set head user if found
		if headUser.ID != 0 {
			result.HeadUser = &headUser
		}
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetAll(params *organizationDTOs.OrganizationUnitReportListQueryParamsDTO) (*organizationDTOs.OrganizationUnitReportListReponseDTO, error) {
	var reports []*organizationDTOs.OrganizationUnitReportListDTO
	var totalCount int64

	query := repo.db.Table("organization_unit_reports").
		Select(`
			organization_unit_reports.id, 
			organization_unit_reports.title, 
			CONVERT(VARCHAR, organization_unit_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, organization_unit_reports.end_date, 23) AS end_date,
			organization_units.id AS organization_unit_id,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			organization_units.color_code AS organization_unit_color_code,
			organization_unit_reports.created_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units ON organization_units.id = organization_unit_reports.organization_unit_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = organization_unit_reports.created_by_id")

	if params.Keyword != "" {
		query = query.Where(`(organization_unit_reports.title LIKE ? )`, "%"+params.Keyword+"%")
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.OrganizationUnitReportListReponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationUnitsRepository) GetByID(id uint) (*organizationDTOs.OrganizationUnitReportListDTO, error) {
	var organizationUnitReport organizationDTOs.OrganizationUnitReportListDTO

	result := repo.db.Table("organization_unit_reports").
		Select(`
			organization_unit_reports.id, 
			organization_unit_reports.title, 
			CONVERT(VARCHAR, organization_unit_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, organization_unit_reports.end_date, 23) AS end_date,
			organization_units.id AS organization_unit_id,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			organization_units.color_code AS organization_unit_color_code,
			organization_unit_reports.created_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units ON organization_units.id = organization_unit_reports.organization_unit_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = organization_unit_reports.created_by_id").
		Where("organization_unit_reports.id = ?", id).
		Scan(&organizationUnitReport)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &organizationUnitReport, nil
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReport(organizationUnitReport *models.OrganizationUnitReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(organizationUnitReport).Error; err != nil {
		return 0, err
	}
	return organizationUnitReport.ID, nil
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReport(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", organizationUnitReportID).
		Delete(&models.OrganizationUnitReport{}).Error
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitOverview(organizationUnitID uint, level uint) (*organizationDTOs.OrganizationUnitOverviewDTO, error) {
	var result organizationDTOs.OrganizationUnitOverviewDTO

	field := fmt.Sprintf("organization_unit_%d_id", level)

	query := repo.db.Table("organization_units AS ou").
		Select(`
			ou.id AS organization_unit_id, 
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar,
			CASE
    			WHEN head.first_name_en IS NULL AND head.last_name_en IS NULL THEN NULL
    			ELSE LTRIM(RTRIM(CONCAT(head.first_name_en, ' ', head.last_name_en)))
			END AS organization_unit_head_name_en,
			CASE
				WHEN head.first_name_ar IS NULL AND head.last_name_ar IS NULL THEN NULL
				ELSE LTRIM(RTRIM(CONCAT(head.first_name_ar, ' ', head.last_name_ar)))
			END AS organization_unit_head_name_ar,
			head.email AS organization_unit_head_email,
			COUNT(DISTINCT users.id) AS total_employees
		`).
		Joins(fmt.Sprintf(`
			LEFT JOIN users AS head
				ON head.is_head = 1
				AND head.%s = ou.id 
		`, field)).
		Joins(fmt.Sprintf(`
			LEFT JOIN users
				ON users.%s = ou.id
		`, field)).
		Where("ou.status = ?", "active").
		Where("ou.id = ?", organizationUnitID).
		Group(`
		ou.id, 
		ou.title_en,
		ou.title_ar, 
		head.first_name_en, 
		head.last_name_en, 
		head.first_name_ar,
		head.last_name_ar,
		head.email
		`)

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitKeyMetrics(organizationUnitID uint) (*organizationDTOs.OrganizationUnitKeyMatrixDTO, error) {
	result := &organizationDTOs.OrganizationUnitKeyMatrixDTO{}

	riskQuery := repo.db.Table("risk_register").Where("responsible_unit_id = ?", organizationUnitID)
	if err := riskQuery.Count(&result.TotalRisks).Error; err != nil {
		return nil, err
	}

	policyQuery := repo.db.Table("documents").
		Joins("JOIN document_types ON document_types.id = documents.document_type_id").
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("documents.responsible_unit_id = ?", organizationUnitID).
		Where("document_types.type = ?", "policy").
		Where("document_statuses.type = ?", "published")

	if err := policyQuery.Count(&result.TotalPolicies).Error; err != nil {
		return nil, err
	}

	processQuery := repo.db.Table("processes").Where("responsible_unit_id = ?", organizationUnitID)
	if err := processQuery.Count(&result.TotalProcesses).Error; err != nil {
		return nil, err
	}

	complianceQuery := repo.db.Table("compliance_register").Where("responsible_unit_id = ?", organizationUnitID)
	if err := complianceQuery.Count(&result.TotalCompliances).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *OrganizationUnitsRepository) GetUserRoleDistribution(organizationUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserRoleDistributionDTO, error) {

	results := []organizationDTOs.OrganizationUnitUserRoleDistributionDTO{}

	query := repo.db.Table("user_roles").
		Select(`
			roles.id,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar,
			COUNT(DISTINCT users.id) AS count
		`).
		Joins("JOIN users ON users.id = user_roles.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id")

	if level > 0 {
		field := fmt.Sprintf("organization_unit_%d_id", level)
		query = query.Where(fmt.Sprintf("users.%s = ?", field), organizationUnitID)
	}

	query = query.Group("roles.id, roles.title_en, roles.title_ar")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetUsersByOrganizationUnitReport(organizationUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserDetailsDTO, error) {

	var users []organizationDTOs.OrganizationUnitUserDetailsDTO
	var totalCount int64

	query := repo.db.Table("users").
		Select(`
			users.id,
			users.first_name_en,
			users.last_name_en,
			users.first_name_ar,
			users.last_name_ar,
			users.is_head,
			users.initial,
			users.email,
			users.profile_image_path,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			roles.title_en AS role_title_en,
			roles.title_ar AS role_title_ar
		`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN user_roles ON user_roles.user_id = users.id").
		Joins("LEFT JOIN roles ON roles.id = user_roles.role_id")

	if level > 0 {
		field := fmt.Sprintf("organization_unit_%d_id", level)
		query = query.Where(fmt.Sprintf("users.%s = ?", field), organizationUnitID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *OrganizationUnitsRepository) GetRiskByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitRiskListDTO, error) {

	var risks []organizationDTOs.OrganizationUnitRiskListDTO

	query := repo.db.Table("risk_register").
		Select(`
			risk_register.id AS id,
			risk_register.id AS risk_register_id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
			risk_register.created_at,
			risk_register.updated_at,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code,
			current_likelihood.title_en AS current_likelihood_title_en,
			current_likelihood.title_ar AS current_likelihood_title_ar,
			current_impact.title_en AS current_impact_title_en,
			current_impact.title_ar AS current_impact_title_ar,
			current_risk_rating.title_en AS current_risk_rating_title_en,
			current_risk_rating.title_ar AS current_risk_rating_title_ar,
			current_risk_rating.color_code AS current_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
			LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
			LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
		`).
		Where("risk_register.responsible_unit_id = ?", organizationUnitID).
		Where("risk_statuses.type = ?", "open")

	query = query.Order("risk_register.reported_at DESC")

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	return risks, nil
}

func (repo *OrganizationUnitsRepository) GetPolicyByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitPolicyListDTO, error) {
	var policyList []organizationDTOs.OrganizationUnitPolicyListDTO

	query := `
		SELECT
			documents.id AS id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
			document_versions.version AS document_version,
			document_versions.url AS document_url, 
			document_versions.name AS document_name,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			documents.created_at AS created_at,
			documents.updated_at AS updated_at,
			users.first_name_en AS responsible_user_first_name_en,
			users.last_name_en AS responsible_user_last_name_en,
			users.first_name_ar AS responsible_user_frist_name_ar,
			users.last_name_ar AS responsible_user_last_name_ar,

			-- Check if all required users acknowledged
			CASE
				WHEN dat.type = 'public' THEN
					CASE
						WHEN NOT EXISTS (
							SELECT 1 FROM users u
							JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1
							WHERE NOT EXISTS (
								SELECT 1 FROM document_acknowledge_users dau
								WHERE dau.user_id = u.id AND dau.document_id = documents.id
							)
						) THEN 1 ELSE 0
					END
				ELSE
					CASE 
						WHEN NOT EXISTS (
							SELECT 1 FROM (
								SELECT u.id
								FROM users u
								JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1
								-- Join primary organization level
								LEFT JOIN organization_units oup ON (
									u.organization_unit_1_id = oup.id OR
									u.organization_unit_2_id = oup.id OR
									u.organization_unit_3_id = oup.id OR
									u.organization_unit_4_id = oup.id OR
									u.organization_unit_5_id = oup.id
								)
								LEFT JOIN organization_levels olp ON olp.id = oup.organization_level_id AND olp.is_primary = 1
								WHERE
									EXISTS (
										SELECT 1 FROM document_shared_users dsu
										WHERE dsu.document_id = documents.id AND dsu.user_id = u.id
									)
									OR EXISTS (
										SELECT 1 FROM document_shared_units dsu2
										JOIN organization_units ou2 ON ou2.id = dsu2.organization_unit_id
										JOIN organization_levels ol2 ON ol2.id = ou2.organization_level_id AND ol2.is_primary = 1
										WHERE dsu2.document_id = documents.id
										AND oup.id = ou2.id -- Only compare if both are primary-level units
									)
								AND NOT EXISTS (
									SELECT 1 FROM document_acknowledge_users dau
									WHERE dau.user_id = u.id AND dau.document_id = documents.id
								)
							) AS pending_users
						) THEN 1 ELSE 0
					END
			END AS acknowledged
		FROM documents
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types dat ON dat.id = documents.document_access_type_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users ON users.id = documents.responsible_user_id
		WHERE documents.responsible_unit_id = ? 
			AND document_types.type = 'policy' 
			AND document_statuses.type = 'published'
	`

	if err := repo.db.Raw(query, organizationUnitID).Scan(&policyList).Error; err != nil {

		return nil, err
	}

	return policyList, nil
}

func (repo *OrganizationUnitsRepository) GetProcessByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitProcessListDTO, error) {
	var processes []organizationDTOs.OrganizationUnitProcessListDTO

	query := repo.db.Table("processes").
		Select(`
			processes.id,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			organization_units.title_en As responsible_unit_title_en,
			organization_units.title_ar As responsible_unit_title_ar,
			organization_units.color_code As responsible_unit_color_code,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			processes.marked_as_auditable_item,
			processes.risk_level_id,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			risk_levels.type AS risk_level_type,
			processes.process_criticality_level_id,
			process_criticality_levels.title_en AS process_criticality_level_title_en,
			process_criticality_levels.title_ar AS process_criticality_level_title_ar,
			processes.created_at,
			processes.updated_at,
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_rating_values.value AS risk_rating_value
		`).
		Joins(`
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id
		`).
		Where("processes.responsible_unit_id = ?", organizationUnitID)

	if err := query.Find(&processes).Error; err != nil {
		return nil, err
	}

	return processes, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitComplianceListDTO, error) {
	var complianceList []organizationDTOs.OrganizationUnitComplianceListDTO

	query := repo.db.Table("compliance_register").
		Select(`
			compliance_register.id AS id,
			compliance_categories.title_en AS compliance_category_title_en,
			compliance_categories.title_ar AS compliance_category_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			compliance_register.reference_number AS reference_number,
			compliance_register.standard_code AS standard_code,
			compliance_register.title_en AS title_en,
			compliance_register.title_ar AS title_ar,
			compliance_register.description_en AS description_en,
			compliance_register.description_ar AS description_ar,
			compliance_register.version AS version,
			CONVERT(VARCHAR, compliance_register.deadline, 23) AS deadline,
			CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
			CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
			CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
			CONVERT(VARCHAR, compliance_register.created_at, 23) AS created_at,
			CONVERT(VARCHAR, compliance_register.updated_at, 23) AS updated_at
		`).
		Joins(`
			LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
		`).
		Where("compliance_register.responsible_unit_id = ?", organizationUnitID)

	if err := query.Find(&complianceList).Error; err != nil {
		return nil, err
	}

	return complianceList, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceRequirementsByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitComplianceRequirementListDTO, error) {
	var result []organizationDTOs.OrganizationUnitComplianceRequirementListDTO

	query := repo.db.Table("compliance_requirements cr").
		Select(`
			cr.id AS id,
			cr.title_en AS title_en,
			cr.title_ar AS title_ar,
			cr.score AS score,
			cr.deadline AS deadline,
			cr.created_at AS created_at,
			cr.updated_at AS updated_at,
			crc.title_en AS compliance_requirement_category_title_en,
			crc.title_ar AS compliance_requirement_category_title_ar,
			cs.title_en AS compliance_status_title_en,
			cs.title_ar AS compliance_status_title_ar,
			cs.color_code AS compliance_status_color_code,
			sl.title_en AS severity_level_title_en,
			sl.title_ar AS severity_level_title_ar,
			sl.color_code AS severity_level_color_code
		`).
		Joins(`
			JOIN compliance_requirement_responsible_units crru
				ON crru.compliance_requirement_id = cr.id
		`).
		Joins(`
			LEFT JOIN compliance_requirement_categories crc
				ON crc.id = cr.compliance_requirement_category_id
		`).
		Joins(`
			LEFT JOIN compliance_statuses cs
				ON cs.id = cr.compliance_status_id
		`).
		Joins(`
			LEFT JOIN severity_levels sl
				ON sl.id = cr.severity_level_id
		`).
		Where("crru.responsible_unit_id = ?", organizationUnitID)

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *OrganizationUnitsRepository) GetRiskSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitRiskSummaryDTO, error) {
	var result organizationDTOs.OrganizationUnitRiskSummaryDTO

	query := repo.db.Table("risk_register").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_register.responsible_unit_id = ?", organizationUnitID).
		Where("risk_statuses.type IN ?", []string{"open", "closed"}).
		Select(`
			SUM(CASE WHEN risk_statuses.type = 'open' THEN 1 ELSE 0 END) AS open_risks,
			SUM(CASE WHEN risk_statuses.type = 'closed' THEN 1 ELSE 0 END) AS closed_risks,
			COUNT(*) AS total_risks
		`)

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	if result.TotalRisks > 0 {
		result.ClosurePercentage = (float64(result.ClosedRisks) / float64(result.TotalRisks)) * 100
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetRiskRatingDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	query := repo.db.Table("risk_register").
		Select(`
			ISNULL(current_risk_rating.id, 0) AS id,
			ISNULL(current_risk_rating.title_en, '') AS title_en,
			ISNULL(current_risk_rating.title_ar, '') AS title_ar,
			ISNULL(current_risk_rating.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN risk_statuses
				ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating
				ON current_risk_rating.id = risk_register.current_risk_rating_id
		`).
		Where("risk_register.responsible_unit_id = ?", organizationUnitID).
		Where("risk_statuses.type = ?", "open")

	query = query.Group(`
			current_risk_rating.id,
			current_risk_rating.title_en,
			current_risk_rating.title_ar,
			current_risk_rating.color_code
		`).
		Order("current_risk_rating.id ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetPolicySummaryByOrganizationUnit(organizationUnitID uint, level uint) (*organizationDTOs.OrganizationUnitPolicySummaryDTO, error) {
	var result organizationDTOs.OrganizationUnitPolicySummaryDTO
	var totalUsers, acknowledgedUsers int64

	orgLevelField := fmt.Sprintf("organization_unit_%d_id", level)

	baseQuery := repo.db.Table("documents").
		Joins("LEFT JOIN document_types dt ON dt.id = documents.document_type_id").
		Where("dt.type = 'policy'").
		Where("documents.responsible_unit_id = ?", organizationUnitID)

	if err := baseQuery.Session(&gorm.Session{}).
		Select("COUNT(*)").
		Scan(&result.TotalPolicies).Error; err != nil {
		return nil, err
	}

	if err := baseQuery.Session(&gorm.Session{}).
		Joins("LEFT JOIN document_statuses ds ON ds.id = documents.document_status_id").
		Where("(ds.type IN ('approved', 'published') OR documents.document_status_id IS NULL)").
		Select("COUNT(*)").
		Scan(&result.ActivePolicies).Error; err != nil {
		return nil, err
	}

	if err := baseQuery.Session(&gorm.Session{}).
		Joins("LEFT JOIN document_statuses ds ON ds.id = documents.document_status_id").
		Where("ds.type = 'in-review'").
		Select("COUNT(*)").
		Scan(&result.PoliciesUnderReview).Error; err != nil {
		return nil, err
	}

	totalUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id)
		FROM (
			-- Explicitly shared users
			SELECT u.id AS user_id
			FROM document_shared_users dsu
			JOIN documents d ON d.id = dsu.document_id
			JOIN document_types dt ON dt.id = d.document_type_id
			JOIN users u ON u.id = dsu.user_id
			JOIN user_statuses us ON us.id = u.user_status_id
			WHERE dt.type = 'policy'
				AND d.responsible_unit_id = ?
				AND us.is_active = 1

			UNION

			-- Unit-based shared users
			SELECT u.id AS user_id
			FROM document_shared_units dsu
			JOIN documents d ON d.id = dsu.document_id
			JOIN document_types dt ON dt.id = d.document_type_id
			JOIN organization_units ou ON ou.id = dsu.organization_unit_id
			JOIN users u ON u.%s = ou.id
			JOIN user_statuses us ON us.id = u.user_status_id
			WHERE dt.type = 'policy'
				AND d.responsible_unit_id = ?
				AND us.is_active = 1
		) AS combined_users
	`, orgLevelField)

	if err := repo.db.Raw(totalUsersQuery, organizationUnitID, organizationUnitID).Scan(&totalUsers).Error; err != nil {
		return nil, fmt.Errorf("error fetching total users: %v", err)
	}

	acknowledgedUsersQuery := `
		SELECT COUNT(DISTINCT dau.user_id)
		FROM document_acknowledge_users dau
		JOIN documents d ON d.id = dau.document_id
		JOIN document_types dt ON dt.id = d.document_type_id
		JOIN users u ON u.id = dau.user_id
		JOIN user_statuses us ON us.id = u.user_status_id
		WHERE dt.type = 'policy'
			AND d.responsible_unit_id = ?
			AND us.is_active = 1
	`

	if err := repo.db.Raw(acknowledgedUsersQuery, organizationUnitID).Scan(&acknowledgedUsers).Error; err != nil {
		return nil, fmt.Errorf("error fetching acknowledged users: %v", err)
	}

	if result.TotalPolicies > 0 {
		result.PolicyHealthPercentage = float64(result.ActivePolicies) / float64(result.TotalPolicies) * 100
	}
	if totalUsers > 0 {
		result.PolicyAwarenessPercentage = float64(acknowledgedUsers) / float64(totalUsers) * 100
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetProcessHealthSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitProcessHealthSummaryDTO, error) {
	var result organizationDTOs.OrganizationUnitProcessHealthSummaryDTO

	query := repo.db.Table("processes").
		Select(`
		SUM(CASE WHEN lifecycle_statuses.type = 'active' THEN 1 ELSE 0 END) AS active_processes,
		COUNT(*) AS total_processes
	`).
		Joins("JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id").
		Where("processes.responsible_unit_id = ?", organizationUnitID)

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	if result.TotalProcesses > 0 {
		result.ActivePercentage = (float64(result.ActiveProcesses) / float64(result.TotalProcesses)) * 100
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetProcessCriticalityDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	query := repo.db.Table("processes").
		Select(`
			ISNULL(process_criticality_levels.id, 0) AS id,
			ISNULL(process_criticality_levels.title_en, '') AS title_en,
			ISNULL(process_criticality_levels.title_ar, '') AS title_ar,
			ISNULL(process_criticality_levels.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN process_criticality_levels
				ON process_criticality_levels.id = processes.process_criticality_level_id
		`).
		Where("processes.responsible_unit_id = ?", organizationUnitID)

	query = query.Group(`
			process_criticality_levels.id,
			process_criticality_levels.title_en,
			process_criticality_levels.title_ar,
			process_criticality_levels.color_code
		`).
		Order("process_criticality_levels.title_en ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceSummaryByOrganizationUnit(organizationUnitID uint) (*organizationDTOs.OrganizationUnitComplianceSummaryDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	query := repo.db.Table("compliance_requirements AS cr").
		Select(`
			cs.id AS id,
			ISNULL(cs.title_en, '') AS title_en,
			ISNULL(cs.title_ar, '') AS title_ar,
			ISNULL(cs.type, '') AS type,
			ISNULL(cs.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			INNER JOIN compliance_requirement_responsible_units AS crru 
				ON crru.compliance_requirement_id = cr.id
			INNER JOIN compliance_statuses AS cs 
				ON cs.id = cr.compliance_status_id
		`).
		Where("crru.responsible_unit_id = ? AND cr.compliance_status_id IS NOT NULL", organizationUnitID).
		Group(`
			cs.id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code
		`).
		Order("cs.title_en ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	var total int64
	for _, item := range results {
		total += item.Count
	}

	summary := &organizationDTOs.OrganizationUnitComplianceSummaryDTO{
		TotalRequirements:  total,
		ComplianceStatuses: results,
	}

	return summary, nil
}

func (repo *OrganizationUnitsRepository) GetRiskMitigationStatusDistributionByOrganizationUnit(organizationUnitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {

	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	query := repo.db.Table("risk_action_plans rap").
		Select(`
			ISNULL(aps.id, 0) AS id,
			ISNULL(aps.title_en, '') AS title_en,
			ISNULL(aps.title_ar, '') AS title_ar,
			ISNULL(aps.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN action_plan_statuses aps
				ON aps.id = rap.action_plan_status_id
			JOIN risk_register rr
				ON rr.id = rap.risk_register_id
		`).
		Where("rr.responsible_unit_id = ?", organizationUnitID)

	query = query.Group(`
			aps.id,
			aps.title_en,
			aps.title_ar,
			aps.color_code
		`).
		Order("aps.id ASC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportOverview(overview *models.OrganizationUnitReportOverview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(overview).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportKeyMetrics(keyMetrics *models.OrganizationUnitReportKeyMetrics, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(keyMetrics).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportUserRole(userRole *models.OrganizationUnitReportUserRole, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(userRole).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportUser(user *models.OrganizationUnitReportUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(user).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportRisk(risk *models.OrganizationUnitReportRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(risk).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportPolicy(policy *models.OrganizationUnitReportPolicy, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(policy).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportProcess(process *models.OrganizationUnitReportProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(process).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportComplianceRequirement(complianceRequirement *models.OrganizationUnitReportComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(complianceRequirement).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportRiskSummary(riskSummary *models.OrganizationUnitReportRiskSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskSummary).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportRiskRating(riskRating *models.OrganizationUnitReportRiskRating, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskRating).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportRiskMitigation(riskMitigation *models.OrganizationUnitReportRiskMitigation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskMitigation).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportPolicySummary(policySummary *models.OrganizationUnitReportPolicySummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(policySummary).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportProcessSummary(processSummary *models.OrganizationUnitReportProcessSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(processSummary).Error
}

func (repo *OrganizationUnitsRepository) CreateOrganizationUnitReportComplianceSummary(complianceSummary *models.OrganizationUnitReportComplianceSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(complianceSummary).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportOverview(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportOverview{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportKeyMetrics(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportKeyMetrics{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportUserRole(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportUserRole{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportUser(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportUser{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportRisk(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportRisk{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportPolicy(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportPolicy{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportProcess(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportProcess{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportComplianceRequirement(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportComplianceRequirement{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportRiskSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportRiskSummary{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportRiskRating(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportRiskRating{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportRiskMitigation(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportRiskMitigation{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportPolicySummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportPolicySummary{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportProcessSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportProcessSummary{}).Error
}

func (repo *OrganizationUnitsRepository) DeleteOrganizationUnitReportComplianceSummary(organizationUnitReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("organization_unit_report_id = ?", organizationUnitReportID).Delete(&models.OrganizationUnitReportComplianceSummary{}).Error
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportOverview(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportOverviewResponseDTO, error) {
	var overview organizationDTOs.OrganizationUnitReportOverviewResponseDTO
	err := repo.db.Table("organization_unit_report_overviews").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&overview).Error

	if err != nil {
		return nil, err
	}
	return &overview, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportKeyMetrics(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportKeyMetricsResponseDTO, error) {
	var keyMetrics organizationDTOs.OrganizationUnitReportKeyMetricsResponseDTO
	err := repo.db.Table("organization_unit_report_key_metrics").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&keyMetrics).Error

	if err != nil {
		return nil, err
	}
	return &keyMetrics, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportUserRole(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportUserRoleResponseDTO, error) {
	userRoles := []*organizationDTOs.OrganizationUnitReportUserRoleResponseDTO{}
	err := repo.db.Table("organization_unit_report_user_roles").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&userRoles).Error

	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportUser(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportUserResponseDTO, error) {
	users := []*organizationDTOs.OrganizationUnitReportUserResponseDTO{}
	err := repo.db.Table("organization_unit_report_users").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportRisk(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskResponseDTO, error) {
	var risks []*organizationDTOs.OrganizationUnitReportRiskResponseDTO
	err := repo.db.Table("organization_unit_report_risks").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&risks).Error

	if err != nil {
		return nil, err
	}
	return risks, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportPolicy(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportPolicyResponseDTO, error) {
	var policies []*organizationDTOs.OrganizationUnitReportPolicyResponseDTO
	err := repo.db.Table("organization_unit_report_policies").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&policies).Error

	if err != nil {
		return nil, err
	}
	return policies, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportProcess(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportProcessResponseDTO, error) {
	var processes []*organizationDTOs.OrganizationUnitReportProcessResponseDTO
	err := repo.db.Table("organization_unit_report_processes").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&processes).Error

	if err != nil {
		return nil, err
	}
	return processes, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportComplianceRequirement(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportComplianceRequirementResponseDTO, error) {
	var complianceRequirements []*organizationDTOs.OrganizationUnitReportComplianceRequirementResponseDTO
	err := repo.db.Table("organization_unit_report_compliance_requirements").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&complianceRequirements).Error

	if err != nil {
		return nil, err
	}
	return complianceRequirements, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportRiskSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportRiskSummaryResponseDTO, error) {
	var riskSummary organizationDTOs.OrganizationUnitReportRiskSummaryResponseDTO
	err := repo.db.Table("organization_unit_report_risk_summary").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&riskSummary).Error

	if err != nil {
		return nil, err
	}
	return &riskSummary, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportRiskRating(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskRatingResponseDTO, error) {
	riskRatings := []*organizationDTOs.OrganizationUnitReportRiskRatingResponseDTO{}
	err := repo.db.Table("organization_unit_report_risk_ratings").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&riskRatings).Error

	if err != nil {
		return nil, err
	}
	return riskRatings, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportRiskMitigation(organizationUnitReportID uint) ([]*organizationDTOs.OrganizationUnitReportRiskMitigationResponsDTO, error) {
	riskMitigations := []*organizationDTOs.OrganizationUnitReportRiskMitigationResponsDTO{}
	err := repo.db.Table("organization_unit_report_risk_mitigations").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Find(&riskMitigations).Error

	if err != nil {
		return nil, err
	}
	return riskMitigations, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportPolicySummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportPolicySummaryResponseDTO, error) {
	var policySummary organizationDTOs.OrganizationUnitReportPolicySummaryResponseDTO
	err := repo.db.Table("organization_unit_report_policy_summary").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&policySummary).Error

	if err != nil {
		return nil, err
	}
	return &policySummary, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportProcessSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportProcessSummaryResponseDTO, error) {
	var processSummary organizationDTOs.OrganizationUnitReportProcessSummaryResponseDTO
	err := repo.db.Table("organization_unit_report_process_summary").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&processSummary).Error

	if err != nil {
		return nil, err
	}
	return &processSummary, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitReportComplianceSummary(organizationUnitReportID uint) (*organizationDTOs.OrganizationUnitReportComplianceSummaryResponseDTO, error) {
	var complianceSummary organizationDTOs.OrganizationUnitReportComplianceSummaryResponseDTO
	err := repo.db.Table("organization_unit_report_compliance_summary").
		Where("organization_unit_report_id = ?", organizationUnitReportID).
		Scan(&complianceSummary).Error

	if err != nil {
		return nil, err
	}
	return &complianceSummary, nil
}

// GetHeadUsersByOrganizationUnitID gets all head users for a specific organization unit
func (repo *OrganizationUnitsRepository) GetHeadUsersByOrganizationUnitID(organizationUnitID uint) ([]*models.User, error) {
	var users []*models.User

	err := repo.db.
		Where("organization_unit_2_id = ? AND is_head = ?", organizationUnitID, true).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}
