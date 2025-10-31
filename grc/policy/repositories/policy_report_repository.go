package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"math"

	"jethur/utility"

	"gorm.io/gorm"
)

type PolicyReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *policyDTOs.PolicyReportListQueryParamsDTO) (*policyDTOs.PolicyReportLisResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PolicyReport, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyReportListDTO, error)
	Create(policyReport *models.PolicyReport, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(policyReportID uint) error

	GetPolicyOverview(documentID uint) (*policyDTOs.PolicyReportPolicyOverviewDTO, error)
	GetPolicyKeySummary(documentID uint) (*policyDTOs.PolicyReportKeySummaryDTO, error)
	GetPolicyAcknowledgementList(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) ([]*policyDTOs.PolicyAcknowledgementListDTO, error)
	GetPolicyAcknowldegementOrganizationUnitSummary(
		documentID uint,
		versionID uint,
		primaryOrgUnitField string,
		isPublic bool,
	) (*policyDTOs.PolicyAcknowledgementOrganizatationUnitSummaryDTO, error)
	GetPolicyAssociatedRisks(documentID uint) ([]*commonDTOs.AssociatedRiskListDTO, error)
	GetPolicyAssociatedtrategicObjectives(documentID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error)
	GetPolicyAssociatedComplianceRequirements(documentID uint) ([]*documentDTOs.DocumentComplianceRequirementListDTO, error)
	GetPolicyVersions(documentID uint) ([]*policyDTOs.PolicyVersionDTO, error)
	GePolicyAccessUsers(documentID uint, primaryOrgUnitField string) ([]*policyDTOs.PolicyAccessUserDTO, error)
	GetPolicyApprovalHistory(VersionID uint) ([]*policyDTOs.PolicyApprovalHistoryDTO, error)
	GetPolicyAttachments(documentID uint) ([]*policyDTOs.PolicyAttachmentDTO, error)
}

type PolicyReportRepository struct {
	db *gorm.DB
}

func NewPolicyReportRepository(db *gorm.DB) PolicyReportRepositoryInterface {
	return &PolicyReportRepository{
		db: db,
	}
}

func (repo *PolicyReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PolicyReportRepository) GetAll(params *policyDTOs.PolicyReportListQueryParamsDTO) (*policyDTOs.PolicyReportLisResponseDTO, error) {
	var reports []*policyDTOs.PolicyReportListDTO
	var totalCount int64

	query := repo.db.Table("policy_reports").
		Select(`
			policy_reports.id, 
			policy_reports.title, 
			CONVERT(VARCHAR, policy_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, policy_reports.end_date, 23) AS end_date,
			documents.id AS document_id,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			policy_reports.created_at,
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
		Joins("LEFT JOIN documents ON documents.id = policy_reports.document_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = policy_reports.created_by_id")

	if params.Keyword != "" {
		query = query.Where(`(policy_reports.title LIKE ? )`, "%"+params.Keyword+"%")
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

	response := &policyDTOs.PolicyReportLisResponseDTO{
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

func (repo *PolicyReportRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PolicyReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var policyReport *models.PolicyReport

	result := db.First(&policyReport, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return policyReport, nil
}

func (repo *PolicyReportRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*policyDTOs.PolicyReportListDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var policyReport policyDTOs.PolicyReportListDTO

	result := db.Table("policy_reports").
		Select(`
			policy_reports.id, 
			policy_reports.title, 
			CONVERT(VARCHAR, policy_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, policy_reports.end_date, 23) AS end_date,
			documents.id AS document_id,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			policy_reports.created_at,
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
		Joins("LEFT JOIN documents ON documents.id = policy_reports.document_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = policy_reports.created_by_id").
		Where("policy_reports.id = ?", id).
		Scan(&policyReport)

	if result.Error != nil {
		return nil, result.Error
	}

	return &policyReport, nil
}

func (repo *PolicyReportRepository) Create(policyReport *models.PolicyReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(policyReport).Error; err != nil {
		return 0, err
	}
	return policyReport.ID, nil
}

func (repo *PolicyReportRepository) Delete(policyReportID uint) error {

	return repo.db.Where("id = ?", policyReportID).
		Delete(&models.PolicyReport{}).Error
}

func (repo *PolicyReportRepository) GetPolicyOverview(documentID uint) (*policyDTOs.PolicyReportPolicyOverviewDTO, error) {
	var result policyDTOs.PolicyReportPolicyOverviewDTO

	err := repo.db.Table("documents").
		Select(`
			documents.id,
			documents.document_type_id,
			documents.title_en,
			documents.title_ar,
			documents.created_at,
			documents.updated_at,
			documents.next_review_date,
			dc.title_en AS document_category_title_en,
			dc.title_ar AS document_category_title_ar,
			dt.title_en AS document_type_title_en,
			dt.title_ar AS document_type_title_ar,
			da.title_en AS document_access_type_title_en,
			da.title_ar AS document_access_type_title_ar,
			ou.title_en AS responsible_unit_title_en,
			ou.title_ar AS responsible_unit_title_ar,
			u.id AS responsible_user_id,
			u.first_name_en AS responsible_user_first_name_en,
			u.first_name_ar AS responsible_user_first_name_ar,
			u.last_name_en AS responsible_user_last_name_en,
			u.last_name_ar AS responsible_user_last_name_ar,
			u.initial AS responsible_user_initial,
			u.color AS responsible_user_color,
			u.email AS responsible_user_email,
			u.profile_image_path AS responsible_user_profile_image_path,
			ds.title_en AS document_status_title_en,
			ds.title_ar AS document_status_title_ar,
			ds.color_code AS document_status_color_code,
			ds.type AS document_status_type,
			dv.version AS document_version,
			dv.published_at AS published_at,
			pub.first_name_en AS published_user_first_name_en,
			pub.first_name_ar AS published_user_first_name_ar,
			pub.last_name_en AS published_user_last_name_en,
			pub.last_name_ar AS published_user_last_name_ar,
			pub.initial AS published_user_initial,
			pub.color AS published_user_color,
			pub.email AS published_user_email,
			pub.profile_image_path AS published_user_profile_image_path
		`).
		Joins("LEFT JOIN document_categories dc ON dc.id = documents.document_category_id").
		Joins("LEFT JOIN document_types dt ON dt.id = documents.document_type_id").
		Joins("LEFT JOIN document_access_types da ON da.id = documents.document_access_type_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = documents.responsible_unit_id").
		Joins("LEFT JOIN users u ON u.id = documents.responsible_user_id").
		Joins("LEFT JOIN document_statuses ds ON ds.id = documents.document_status_id").
		Joins("LEFT JOIN document_versions dv ON dv.document_id = documents.id AND dv.is_latest_version = 1").
		Joins("LEFT JOIN users pub ON pub.id = dv.published_by_id").
		Where("documents.id = ?", documentID).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching policy overview: %w", err)
	}

	return &result, nil
}

func (repo *PolicyReportRepository) GetPolicyKeySummary(documentID uint) (*policyDTOs.PolicyReportKeySummaryDTO, error) {

	var (
		clausesDefined            int64
		totalRecipients           int64
		ackUsers                  int64
		risksReported             int64
		associatedStandardsCount  int64
		associatedClausesCount    int64
		strategicObjectivesLinked int64
		approvalCycles            int64
		latestVersionID           uint
	)

	if err := repo.db.Table("document_versions dv").
		Select("dv.id").
		Where("dv.document_id = ?", documentID).
		Where("dv.is_latest_version = 1").
		Take(&latestVersionID).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("document_version_clauses dvc").
		Where("dvc.document_version_id = ?", latestVersionID).
		Count(&clausesDefined).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("users a").
		Joins("JOIN user_statuses aus ON aus.id = a.user_status_id AND aus.is_active = 1").
		Where(`
			EXISTS (
				SELECT 1 FROM document_shared_users dsu
				WHERE dsu.document_id = ? AND dsu.user_id = a.id
			) OR EXISTS (
				SELECT 1
				FROM document_shared_units dsu2
				JOIN organization_units ou ON ou.id = dsu2.organization_unit_id
				WHERE dsu2.document_id = ?
				  AND (
					  a.organization_unit_1_id = ou.id OR
					  a.organization_unit_2_id = ou.id OR
					  a.organization_unit_3_id = ou.id OR
					  a.organization_unit_4_id = ou.id OR
					  a.organization_unit_5_id = ou.id
				  )
			)
		`, documentID, documentID).
		Count(&totalRecipients).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("document_acknowledge_users dau").
		Joins("JOIN users u ON u.id = dau.user_id").
		Joins("JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1").
		Where("dau.document_id = ?", documentID).
		Where(`
			EXISTS (
				SELECT 1 FROM document_shared_users dsu
				WHERE dsu.document_id = ? AND dsu.user_id = dau.user_id
			) OR EXISTS (
				SELECT 1
				FROM document_shared_units dsu2
				JOIN organization_units ou ON ou.id = dsu2.organization_unit_id
				WHERE dsu2.document_id = ?
				  AND (
					  u.organization_unit_1_id = ou.id OR
					  u.organization_unit_2_id = ou.id OR
					  u.organization_unit_3_id = ou.id OR
					  u.organization_unit_4_id = ou.id OR
					  u.organization_unit_5_id = ou.id
				  )
			)
		`, documentID, documentID).
		Count(&ackUsers).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("risk_associated_documents rad").
		Where("rad.document_id = ?", documentID).
		Count(&risksReported).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("document_compliance_requirements dcr").
		Where("dcr.document_id = ?", documentID).
		Select("COUNT(DISTINCT dcr.compliance_register_id)").
		Scan(&associatedStandardsCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("document_compliance_requirements dcr2").
		Where("dcr2.document_id = ?", documentID).
		Where("dcr2.compliance_clause_id IS NOT NULL").
		Select("COUNT(DISTINCT dcr2.compliance_clause_id)").
		Scan(&associatedClausesCount).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("strategic_objective_associated_documents soad").
		Where("soad.document_id = ?", documentID).
		Count(&strategicObjectivesLinked).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("document_workflow_reviews dwr").
		Joins("JOIN document_workflows dw ON dw.id = dwr.document_workflow_id").
		Where("dw.document_version_id = ?", latestVersionID).
		Count(&approvalCycles).Error; err != nil {
		return nil, err
	}

	pending := totalRecipients - ackUsers
	if pending < 0 {
		pending = 0
	}
	var rate float64
	if totalRecipients > 0 {
		rate = math.Round((float64(ackUsers)/float64(totalRecipients)*100.0)*100) / 100
	}

	out := &policyDTOs.PolicyReportKeySummaryDTO{
		ClausesDefined:            clausesDefined,
		AcknowledgmentRate:        rate,
		AcknowledgmentsCompleted:  ackUsers,
		AcknowledgmentsPending:    pending,
		RiskReported:              risksReported,
		AssociatedStandardsCount:  associatedStandardsCount,
		AssociatedClausesCount:    associatedClausesCount,
		StrategicObjectivesLinked: strategicObjectivesLinked,
		ApprovalCycle:             approvalCycles,
	}
	return out, nil
}

func (repo *PolicyReportRepository) GetPolicyAcknowledgementList(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) ([]*policyDTOs.PolicyAcknowledgementListDTO, error) {

	usersScope := repo.db.Table("users u").
		Joins("JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1")

	if !isPublic {
		usersScope = usersScope.Where(`
			EXISTS (
				SELECT 1
				FROM document_shared_users dsu
				WHERE dsu.document_id = ? AND dsu.user_id = u.id
			) OR EXISTS (
				SELECT 1
				FROM document_shared_units dsu2
				JOIN organization_units ou2 ON ou2.id = dsu2.organization_unit_id
				WHERE dsu2.document_id = ?
				  AND u.`+primaryOrgUnitField+` = ou2.id
			)
		`, documentID, documentID)
	}

	usersScope = usersScope.Joins("LEFT JOIN organization_units ou ON ou.id = u." + primaryOrgUnitField)

	usersScope = usersScope.Joins(`
		LEFT JOIN document_acknowledge_users dau
		  ON dau.user_id = u.id
		 AND dau.document_id = ?
		 AND dau.document_version_id = ?
	`, documentID, versionID)

	rows := []*policyDTOs.PolicyAcknowledgementListDTO{}
	if err := usersScope.
		Select(`
			u.id                                 AS user_id,
			ISNULL(u.first_name_en, '')          AS first_name_en,
			ISNULL(u.first_name_ar, '')          AS first_name_ar,
			ISNULL(u.last_name_en, '')           AS last_name_en,
			ISNULL(u.last_name_ar, '')           AS last_name_ar,
			ISNULL(u.email, '')                  AS email,
			ISNULL(u.initial, '')                AS initial,
			ISNULL(u.color, '')                  AS color,
			ISNULL(u.profile_image_path, '')     AS profile_image_path,
			ISNULL(ou.title_en, '')              AS department_title_en,
			ISNULL(ou.title_ar, '')              AS department_title_ar,
			CASE WHEN dau.id IS NULL THEN 0 ELSE 1 END AS acknowledged,
			CASE WHEN dau.id IS NULL THEN NULL ELSE CONVERT(VARCHAR, dau.created_at, 23) END AS acknowledged_on,
			dau.comment AS comment
		`).
		Order("acknowledged DESC, u.first_name_en ASC, u.last_name_en ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (repo *PolicyReportRepository) GetPolicyAcknowldegementOrganizationUnitSummary(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) (*policyDTOs.PolicyAcknowledgementOrganizatationUnitSummaryDTO, error) {

	items := []policyDTOs.PolicyAcknowledgementOrganizatationUnitSummaryItemDTO{}

	base := repo.db.Table("users u").
		Joins("JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1")

	if !isPublic {
		base = base.Where(`
			EXISTS (
				SELECT 1
				FROM document_shared_users dsu
				WHERE dsu.document_id = ? AND dsu.user_id = u.id
			) OR EXISTS (
				SELECT 1
				FROM document_shared_units dsu2
				JOIN organization_units ou2 ON ou2.id = dsu2.organization_unit_id
				WHERE dsu2.document_id = ?
				  AND u.`+primaryOrgUnitField+` = ou2.id
			)
		`, documentID, documentID)
	}

	base = base.
		Joins("LEFT JOIN organization_units ou ON ou.id = u."+primaryOrgUnitField).
		Joins(`
			LEFT JOIN document_acknowledge_users dau
			  ON dau.user_id = u.id
			 AND dau.document_id = ?
			 AND dau.document_version_id = ?
		`, documentID, versionID)

	if err := base.
		Select(`
			ou.id AS organization_unit_id,
			ISNULL(ou.title_en, '') AS organization_unit_title_en,
			ISNULL(ou.title_ar, '') AS organization_unit_title_ar,
			COUNT(DISTINCT u.id) AS total_users,
			COUNT(DISTINCT dau.user_id) AS acknowledged_users,
			COUNT(DISTINCT u.id) - COUNT(DISTINCT dau.user_id) AS not_acknowledged_users,
			CASE WHEN COUNT(DISTINCT u.id) = 0 THEN 0
			     ELSE CAST(COUNT(DISTINCT dau.user_id) AS FLOAT) * 100.0 / COUNT(DISTINCT u.id)
			END AS acknowledged_percent
		`).
		Group("ou.id, ou.title_en, ou.title_ar").
		Order("organization_unit_title_en ASC").
		Scan(&items).Error; err != nil {
		return nil, err
	}

	var totalUsers int64
	var totalAck int64
	for _, it := range items {
		totalUsers += it.TotalUsers
		totalAck += it.AcknowledgedUsers
	}
	totalPending := totalUsers - totalAck

	var totalPercent float64
	if totalUsers > 0 {
		totalPercent = float64(totalAck) * 100.0 / float64(totalUsers)
	}

	resp := &policyDTOs.PolicyAcknowledgementOrganizatationUnitSummaryDTO{
		Items:                     items,
		TotalUsers:                totalUsers,
		TotalAcknowledgedUsers:    totalAck,
		TotalNotAcknowledgedUsers: totalPending,
		TotalAcknowledgedPercent:  totalPercent,
	}
	return resp, nil
}

func (repo *PolicyReportRepository) GetPolicyAssociatedRisks(documentID uint) ([]*commonDTOs.AssociatedRiskListDTO, error) {
	var documentRisk []*commonDTOs.AssociatedRiskListDTO

	query := repo.db.Table("risk_associated_documents").
		Select(`
		    risk_associated_documents.id,
		    risk_associated_documents.risk_register_id,
		    risk_associated_documents.document_id,
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
			LEFT JOIN risk_register ON risk_register.id = risk_associated_documents.risk_register_id
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
		Where("risk_associated_documents.document_id = ?", documentID)

	err := query.Find(&documentRisk).Error
	if err != nil {
		return nil, err
	}

	return documentRisk, nil
}

func (repo *PolicyReportRepository) GetPolicyAssociatedtrategicObjectives(documentID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error) {
	var objectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO

	query := repo.db.Table("strategic_objective_associated_documents").
		Select(`
			strategic_objectives.*,
			strategic_objective_associated_documents.strategic_objective_id,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,

			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

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
			LEFT JOIN strategic_objectives ON strategic_objectives.id = strategic_objective_associated_documents.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
        `).Where("strategic_objective_associated_documents.document_id = ?", documentID)

	err := query.Find(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	return objectiveList, nil
}

func (repo *PolicyReportRepository) GetPolicyAssociatedComplianceRequirements(documentID uint) ([]*documentDTOs.DocumentComplianceRequirementListDTO, error) {
	var complianceRequirementList []*documentDTOs.DocumentComplianceRequirementListDTO

	query := repo.db.Table("compliance_requirement_associated_documents").
		Select(`
			compliance_requirement_associated_documents.compliance_register_id AS compliance_register_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			compliance_requirement_associated_documents.compliance_requirement_id AS compliance_requirement_id,
			compliance_clauses.title_en AS compliance_clause_title_en,
			compliance_clauses.title_ar AS compliance_clause_title_ar,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			compliance_requirements.title_en AS compliance_requirement_title_en,
			compliance_requirements.title_ar AS compliance_requirement_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code
		`).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_associated_documents.compliance_requirement_id
			LEFT JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id
			LEFT JOIN compliance_requirement_categories ON compliance_requirements.compliance_requirement_category_id = compliance_requirement_categories.id
			LEFT JOIN compliance_register ON compliance_register.id = compliance_requirement_associated_documents.compliance_register_id
		`).
		Where("compliance_requirement_associated_documents.document_id = ?", documentID)

	err := query.Find(&complianceRequirementList).Error
	if err != nil {
		return nil, err
	}

	return complianceRequirementList, nil
}

func (repo *PolicyReportRepository) GetPolicyVersions(documentID uint) ([]*policyDTOs.PolicyVersionDTO, error) {
	var versions []*policyDTOs.PolicyVersionDTO

	err := repo.db.Table("document_versions").
		Select(`
			document_versions.id,
			document_versions.document_id,
			document_versions.version,
			document_versions.is_latest_version,
			document_versions.enable_acknowledgement,
			document_versions.created_at,
			document_versions.updated_at,
			published_by_users.first_name_en AS published_by_first_name_en,
			published_by_users.last_name_en AS published_by_last_name_en,
			published_by_users.email AS published_by_email
		`).
		Joins("LEFT JOIN users published_by_users ON published_by_users.id = document_versions.published_by_id").
		Where("document_versions.document_id = ?", documentID).
		Order("document_versions.version DESC").
		Scan(&versions).Error

	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (repo *PolicyReportRepository) GePolicyAccessUsers(documentID uint, primaryOrgUnitField string) ([]*policyDTOs.PolicyAccessUserDTO, error) {
	var rows []*policyDTOs.PolicyAccessUserDTO

	err := repo.db.Table("document_shared_users dsu").
		Joins("JOIN users u ON u.id = dsu.user_id").
		Joins("LEFT JOIN designations des ON des.id = u.designation_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = u."+primaryOrgUnitField).
		Where("dsu.document_id = ?", documentID).
		Select(`
			DISTINCT
			u.id AS user_id,
			ISNULL(u.first_name_en, '') AS first_name_en,
			ISNULL(u.last_name_en, '')  AS last_name_en,
			ISNULL(u.first_name_ar, '') AS first_name_ar,
			ISNULL(u.last_name_ar, '')  AS last_name_ar,
			ISNULL(u.email, '')         AS email,
			u.mobile_number             AS phone,
			u.profile_image_path        AS profile_image_path,
			ISNULL(des.title_en, '')    AS designation_title_en,
			ISNULL(des.title_ar, '')    AS designation_title_ar,
			ISNULL(ou.title_en, '')     AS organization_unit_title_en,
			ISNULL(ou.title_ar, '')     AS organization_unit_title_ar
		`).
		Order("first_name_en ASC, last_name_en ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (repo *PolicyReportRepository) GetPolicyApprovalHistory(versionID uint) ([]*policyDTOs.PolicyApprovalHistoryDTO, error) {
	var history []*policyDTOs.PolicyApprovalHistoryDTO

	err := repo.db.Table("document_workflow_reviews dwr").
		Joins("JOIN document_workflows dw ON dw.id = dwr.document_workflow_id").
		Joins("JOIN users u ON u.id = dwr.created_by_id").
		Joins("LEFT JOIN designations des ON des.id = u.designation_id").
		Joins("JOIN document_workflow_statuses dws ON dws.id = dwr.document_workflow_status_id").
		Where("dw.document_version_id = ?", versionID).
		Select(`
			dwr.id AS id,
			dwr.document_workflow_id AS document_workflow_id,
			ISNULL(dw.level, 0) AS workflow_level,
			u.id AS approved_by_id,
			ISNULL(u.first_name_en, '') AS approved_by_first_name_en,
			ISNULL(u.first_name_ar, '') AS approved_by_first_name_ar,
			ISNULL(u.last_name_en, '') AS approved_by_last_name_en,
			ISNULL(u.last_name_ar, '') AS approved_by_last_name_ar,
			ISNULL(des.title_en, '') AS designation_title_en,
			ISNULL(des.title_ar, '') AS designation_title_ar,
			ISNULL(dws.title_en, '') AS action_title_en,
			ISNULL(dws.title_ar, '') AS action_title_ar,
			ISNULL(dws.type, '') AS action_type,
			ISNULL(dws.color_code, '') AS action_color_code,
			CASE WHEN dwr.created_at IS NULL THEN NULL ELSE CONVERT(VARCHAR, dwr.created_at, 106) END AS date,
			ISNULL(dwr.comment, '') AS comment
		`).
		Order("dwr.created_at ASC, dwr.id ASC").
		Scan(&history).Error
	if err != nil {
		return nil, err
	}

	return history, nil
}

func (repo *PolicyReportRepository) GetPolicyAttachments(documentID uint) ([]*policyDTOs.PolicyAttachmentDTO, error) {
	var rows []*policyDTOs.PolicyAttachmentDTO

	err := repo.db.Table("document_associated_documents dad").
		Joins("JOIN documents d ON d.id = dad.associated_document_id").
		Joins("LEFT JOIN document_types dt ON dt.id = d.document_type_id").
		Where("dad.document_id = ?", documentID).
		Select(`
			d.id AS document_id,
			ISNULL(d.title_en, '') AS document_title_en,
			ISNULL(d.title_ar, '') AS document_title_ar,
			ISNULL(dt.title_en, '') AS document_type_title_en,
			ISNULL(dt.title_ar, '') AS document_type_title_ar,
			CONVERT(VARCHAR, dad.created_at, 106) AS created_at
		`).
		Order("dad.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
