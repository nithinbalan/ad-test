package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"math"

	"jethur/utility"

	"gorm.io/gorm"
)

type PolicyReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *policyDTOs.PolicyReportListQueryParamsDTO) (*policyDTOs.PolicyReportLisResponseDTO, error)
	GetByID(id uint) (*models.PolicyReport, error)
	GetFullDetailsByID(id uint) (*policyDTOs.PolicyReportListDTO, error)
	Create(policyReport *models.PolicyReport, tx ...utility.DatabaseTransaction) (uint, error)
	Delete(policyReportID uint, tx ...utility.DatabaseTransaction) error

	GetPolicyOverview(documentID uint) (*policyDTOs.PolicyReportPolicyOverviewDTO, error)
	GetPolicyKeyMetrics(documentID uint, primaryOrgUnitField string, isPublic bool) (*policyDTOs.PolicyReportKeySummaryDTO, error)
	GetPolicyAcknowledgementList(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) ([]*policyDTOs.PolicyAcknowledgementListDTO, error)
	GetPolicyAcknowldegementOrganizationUnit(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) ([]*policyDTOs.PolicyAcknowledgementOrganizatationUnitDTO, error)
	GetPolicyAssociatedRisks(documentID uint) ([]*policyDTOs.AssociatedRiskListDTO, error)
	GetPolicyAssociatedtrategicObjectives(documentID uint) ([]*policyDTOs.AssociatedStrategicObjectiveListDTO, error)
	GetPolicyVersions(documentID uint) ([]*policyDTOs.PolicyVersionDTO, error)
	GePolicyAccessUsers(documentID uint, primaryOrgUnitField string) ([]*policyDTOs.PolicyAccessUserDTO, error)
	GetPolicyApprovalHistory(VersionID uint) ([]*policyDTOs.PolicyApprovalHistoryDTO, error)
	GetPolicyAttachments(documentID uint) ([]*policyDTOs.PolicyAttachmentDTO, error)
	GetPolicyAssocatedCompliance(documentID uint) ([]*commonDTOs.AssociatedClauseListDTO, error)

	CreatePolicyReportOverview(overview *models.PolicyReportOverview, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportKeyMetrics(keyMetrics *models.PolicyReportKeyMetrics, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportClause(clause *models.PolicyReportClause, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportAcknowledgmentUser(user *models.PolicyReportAcknowledgementUser, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportAcknowledgmentOrganizationUnit(organizationUnit *models.PolicyReportAcknowledgementOrganizationUnit, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportRisk(risk *models.PolicyReportRisk, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportStrategicObjective(strategicObjective *models.PolicyReportStrategicObjective, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportCompliance(compliance *models.PolicyReportCompliance, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportDocumentVersion(documentVersion *models.PolicyReportDocumentVersion, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportDocumentAccess(documentAccess *models.PolicyReportDocumentAccess, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportDocumentApprovalHistory(documentApprovalHistory *models.PolicyReportDocumentApprovalHistory, tx ...utility.DatabaseTransaction) error
	CreatePolicyReportDocumentAttachment(documentAttachment *models.PolicyReportDocumentAttachment, tx ...utility.DatabaseTransaction) error

	DeletePolicyReportOverview(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportKeyMetrics(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportClause(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportAcknowledgmentUser(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportAcknowledgmentOrganizationUnit(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportRisk(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportStrategicObjective(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportCompliance(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportDocumentVersion(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportDocumentAccess(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportDocumentApprovalHistory(policyReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePolicyReportDocumentAttachment(policyReportID uint, tx ...utility.DatabaseTransaction) error

	GetPolicyReportOverview(policyReportID uint) (*policyDTOs.PolicyReportOverviewResponseDTO, error)
	GetPolicyReportKeyMetrics(policyReportID uint) (*policyDTOs.PolicyReportKeyMetricsResponseDTO, error)
	GetPolicyReportClause(policyReportID uint) ([]*policyDTOs.PolicyReportClauseResponseDTO, error)
	GetPolicyReportAcknowledgmentUser(policyReportID uint) ([]*policyDTOs.PolicyReportAcknowledgementUserResponseDTO, error)
	GetPolicyReportAcknowledgmentOrganizationUnit(policyReportID uint) ([]*policyDTOs.PolicyReportAcknowledgementOrganizationUnitResponseDTO, error)
	GetPolicyReportRisk(policyReportID uint) ([]*policyDTOs.PolicyReportRiskResponseDTO, error)
	GetPolicyReportStrategicObjective(policyReportID uint) ([]*policyDTOs.PolicyReportStrategicObjectiveResponseDTO, error)
	GetPolicyReportCompliance(policyReportID uint) ([]*policyDTOs.PolicyReportComplianceResponseDTO, error)
	GetPolicyReportDocumentVersion(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentVersionResponseDTO, error)
	GetPolicyReportDocumentAccess(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentAccessResponseDTO, error)
	GetPolicyReportDocumentApprovalHistory(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentApprovalHistoryResponseDTO, error)
	GetPolicyReportDocumentAttachment(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentAttachmentResponseDTO, error)
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

func (repo *PolicyReportRepository) GetByID(id uint) (*models.PolicyReport, error) {

	var policyReport *models.PolicyReport

	result := repo.db.First(&policyReport, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return policyReport, nil
}

func (repo *PolicyReportRepository) GetFullDetailsByID(id uint) (*policyDTOs.PolicyReportListDTO, error) {

	var policyReport policyDTOs.PolicyReportListDTO

	result := repo.db.Table("policy_reports").
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

func (repo *PolicyReportRepository) Delete(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", policyReportID).
		Delete(&models.PolicyReport{}).Error
}

func (repo *PolicyReportRepository) GetPolicyOverview(documentID uint) (*policyDTOs.PolicyReportPolicyOverviewDTO, error) {
	var result policyDTOs.PolicyReportPolicyOverviewDTO

	err := repo.db.Table("documents").
		Select(`
			documents.id,
			documents.document_type_id,
			documents.reference_number,
			documents.title_en,
			documents.title_ar,
			documents.description_en,
			documents.description_ar,
			documents.created_at,
			documents.updated_at,
			FORMAT(documents.next_review_date, 'dd-MMM-yyyy') AS next_review_date,
			dc.title_en AS document_category_title_en,
			dc.title_ar AS document_category_title_ar,
			dt.title_en AS document_type_title_en,
			dt.title_ar AS document_type_title_ar,
			da.title_en AS document_access_type_title_en,
			da.title_ar AS document_access_type_title_ar,
			da.color_code AS document_access_type_color_code,
			ou.title_en AS responsible_unit_title_en,
			ou.title_ar AS responsible_unit_title_ar,
			u.first_name_en AS responsible_user_first_name_en,
			u.first_name_ar AS responsible_user_first_name_ar,
			u.last_name_en AS responsible_user_last_name_en,
			u.last_name_ar AS responsible_user_last_name_ar,
			ds.title_en AS document_status_title_en,
			ds.title_ar AS document_status_title_ar,
			ds.color_code AS document_status_color_code,
			ds.type AS document_status_type,
			dv.version AS document_version,
			FORMAT(dv.published_at, 'dd-MMM-yyyy') AS published_at
		`).
		Joins("LEFT JOIN document_categories dc ON dc.id = documents.document_category_id").
		Joins("LEFT JOIN document_types dt ON dt.id = documents.document_type_id").
		Joins("LEFT JOIN document_access_types da ON da.id = documents.document_access_type_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = documents.responsible_unit_id").
		Joins("LEFT JOIN users u ON u.id = documents.responsible_user_id").
		Joins("LEFT JOIN document_statuses ds ON ds.id = documents.document_status_id").
		Joins("LEFT JOIN document_versions dv ON dv.document_id = documents.id AND dv.is_latest_version = 1").
		Where("documents.id = ?", documentID).
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching policy overview: %w", err)
	}

	return &result, nil
}

func (repo *PolicyReportRepository) GetPolicyKeyMetrics(documentID uint, primaryOrgUnitField string, isPublic bool) (*policyDTOs.PolicyReportKeySummaryDTO, error) {
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

	// Get latest document version
	if err := repo.db.Table("document_versions dv").
		Select("dv.id").
		Where("dv.document_id = ?", documentID).
		Where("dv.is_latest_version = 1").
		Take(&latestVersionID).Error; err != nil {
		return nil, err
	}

	// Clauses defined
	if err := repo.db.Table("document_version_clauses dvc").
		Where("dvc.document_version_id = ?", latestVersionID).
		Count(&clausesDefined).Error; err != nil {
		return nil, err
	}

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

	// Count total recipients
	if err := usersScope.Count(&totalRecipients).Error; err != nil {
		return nil, err
	}

	// Count acknowledged users
	ackScope := repo.db.Table("document_acknowledge_users dau").
		Joins("JOIN users u ON u.id = dau.user_id").
		Joins("JOIN user_statuses us ON us.id = u.user_status_id AND us.is_active = 1").
		Where("dau.document_id = ?", documentID)

	if !isPublic {
		ackScope = ackScope.Where(`
			EXISTS (
				SELECT 1
				FROM document_shared_users dsu
				WHERE dsu.document_id = ? AND dsu.user_id = dau.user_id
			) OR EXISTS (
				SELECT 1
				FROM document_shared_units dsu2
				JOIN organization_units ou2 ON ou2.id = dsu2.organization_unit_id
				WHERE dsu2.document_id = ?
				  AND u.`+primaryOrgUnitField+` = ou2.id
			)
		`, documentID, documentID)
	}

	if err := ackScope.Count(&ackUsers).Error; err != nil {
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
		TotalAcknowledgments:      totalRecipients,
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
			u.id AS user_id,
			u.first_name_en AS first_name_en,
			u.first_name_ar AS first_name_ar,
			u.last_name_en AS last_name_en,
			u.last_name_ar AS last_name_ar,
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar,
			CAST(CASE WHEN dau.id IS NULL THEN 0 ELSE 1 END AS BIT) AS acknowledged,
			FORMAT(dau.created_at, 'dd-MMM-yyyy') AS acknowledged_on,
			dau.comment AS comment
		`).
		Order("acknowledged DESC, u.first_name_en ASC, u.last_name_en ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (repo *PolicyReportRepository) GetPolicyAcknowldegementOrganizationUnit(documentID uint, versionID uint, primaryOrgUnitField string, isPublic bool) ([]*policyDTOs.PolicyAcknowledgementOrganizatationUnitDTO, error) {

	results := []*policyDTOs.PolicyAcknowledgementOrganizatationUnitDTO{}

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
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar,
			COUNT(DISTINCT u.id) AS total_users,
			COUNT(DISTINCT dau.user_id) AS acknowledged_users,
			COUNT(DISTINCT u.id) - COUNT(DISTINCT dau.user_id) AS not_acknowledged_users,
			CASE WHEN COUNT(DISTINCT u.id) = 0 THEN 0
			     ELSE CAST(COUNT(DISTINCT dau.user_id) AS FLOAT) * 100.0 / COUNT(DISTINCT u.id)
			END AS acknowledged_percent
		`).
		Group("ou.id, ou.title_en, ou.title_ar").
		Order("organization_unit_title_en ASC").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PolicyReportRepository) GetPolicyAssociatedRisks(documentID uint) ([]*policyDTOs.AssociatedRiskListDTO, error) {
	var documentRisk []*policyDTOs.AssociatedRiskListDTO

	query := repo.db.Table("risk_associated_documents").
		Select(`
            risk_register.title,
            risk_register.description,
            risk_register.created_at,
            risk_register.updated_at,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating.color_code AS current_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = risk_associated_documents.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
		`).
		Where("risk_associated_documents.document_id = ?", documentID)

	err := query.Find(&documentRisk).Error
	if err != nil {
		return nil, err
	}

	return documentRisk, nil
}

func (repo *PolicyReportRepository) GetPolicyAssociatedtrategicObjectives(documentID uint) ([]*policyDTOs.AssociatedStrategicObjectiveListDTO, error) {
	var objectiveList []*policyDTOs.AssociatedStrategicObjectiveListDTO

	query := repo.db.Table("strategic_objective_associated_documents").
		Select(`
			strategic_objectives.title_en,
			strategic_objectives.title_ar,
			FORMAT(strategic_objectives.start_date, 'dd-MMM-yyyy') AS start_date,
			FORMAT(strategic_objectives.end_date, 'dd-MMM-yyyy') AS end_date,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar
        `).
		Joins(`
			LEFT JOIN strategic_objectives ON strategic_objectives.id = strategic_objective_associated_documents.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
        `).Where("strategic_objective_associated_documents.document_id = ?", documentID)

	err := query.Find(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	return objectiveList, nil
}

func (repo *PolicyReportRepository) GetPolicyAssocatedCompliance(documentID uint) ([]*commonDTOs.AssociatedClauseListDTO, error) {

	var clauseList []commonDTOs.AssociatedClauseRowDTO
	query := repo.db.Table("document_compliance_requirements").
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
			LEFT JOIN compliance_clauses ON document_compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON document_compliance_requirements.compliance_register_id = compliance_register.id
		`).
		Where("document_compliance_requirements.document_id = ?", documentID)

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

	return groupedClauses, nil
}

func (repo *PolicyReportRepository) GetPolicyVersions(documentID uint) ([]*policyDTOs.PolicyVersionDTO, error) {
	var versions []*policyDTOs.PolicyVersionDTO

	err := repo.db.Table("document_versions").
		Select(`
			document_versions.id,
			document_versions.document_id,
			document_versions.version,
			FORMAT(document_versions.published_at, 'dd-MMM-yyyy') AS published_at,
			published_by_users.first_name_en AS published_by_first_name_en,
			published_by_users.first_name_ar AS published_by_first_name_ar,
			published_by_users.last_name_en AS published_by_last_name_en,
			published_by_users.last_name_ar AS published_by_last_name_ar
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
			u.first_name_en AS first_name_en,
			u.last_name_en AS last_name_en,
			u.first_name_ar AS first_name_ar,
			u.last_name_ar AS last_name_ar,
			des.title_en AS designation_title_en,
			des.title_ar AS designation_title_ar,
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar
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
			dw.level AS workflow_level,
			u.id AS approved_by_id,
			u.first_name_en AS approved_by_first_name_en,
			u.first_name_ar AS approved_by_first_name_ar,
			u.last_name_en AS approved_by_last_name_en,
			u.last_name_ar AS approved_by_last_name_ar,
			des.title_en AS designation_title_en,
			des.title_ar AS designation_title_ar,
			dws.title_en AS action_title_en,
			dws.title_ar AS action_title_ar,
			dws.type AS action_type,
			dws.color_code AS action_color_code,
			FORMAT(dwr.created_at, 'dd-MMM-yyyy') AS action_date,
			dwr.comment AS comment
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
			d.title_en AS document_title_en,
			d.title_ar AS document_title_ar,
			dt.title_en AS document_type_title_en,
			dt.title_ar AS document_type_title_ar,
			FORMAT(dad.created_at, 'dd-MMM-yyyy') AS created_at
		`).
		Order("dad.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (repo *PolicyReportRepository) CreatePolicyReportOverview(overview *models.PolicyReportOverview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(overview).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportKeyMetrics(keyMetrics *models.PolicyReportKeyMetrics, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(keyMetrics).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportClause(clause *models.PolicyReportClause, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(clause).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportAcknowledgmentUser(user *models.PolicyReportAcknowledgementUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(user).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportAcknowledgmentOrganizationUnit(organizationUnit *models.PolicyReportAcknowledgementOrganizationUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(organizationUnit).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportRisk(risk *models.PolicyReportRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(risk).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportStrategicObjective(strategicObjective *models.PolicyReportStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(strategicObjective).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportCompliance(compliance *models.PolicyReportCompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(compliance).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportDocumentVersion(documentVersion *models.PolicyReportDocumentVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(documentVersion).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportDocumentAccess(documentAccess *models.PolicyReportDocumentAccess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(documentAccess).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportDocumentApprovalHistory(documentApprovalHistory *models.PolicyReportDocumentApprovalHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(documentApprovalHistory).Error
}

func (repo *PolicyReportRepository) CreatePolicyReportDocumentAttachment(documentAttachment *models.PolicyReportDocumentAttachment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(documentAttachment).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportOverview(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportOverview{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportKeyMetrics(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportKeyMetrics{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportClause(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportClause{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportAcknowledgmentUser(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportAcknowledgementUser{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportAcknowledgmentOrganizationUnit(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportAcknowledgementOrganizationUnit{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportRisk(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportRisk{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportStrategicObjective(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportStrategicObjective{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportCompliance(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportCompliance{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportDocumentVersion(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportDocumentVersion{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportDocumentAccess(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportDocumentAccess{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportDocumentApprovalHistory(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportDocumentApprovalHistory{}).Error
}

func (repo *PolicyReportRepository) DeletePolicyReportDocumentAttachment(policyReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("policy_report_id = ?", policyReportID).Delete(&models.PolicyReportDocumentAttachment{}).Error
}

func (repo *PolicyReportRepository) GetPolicyReportOverview(policyReportID uint) (*policyDTOs.PolicyReportOverviewResponseDTO, error) {
	var overview policyDTOs.PolicyReportOverviewResponseDTO
	err := repo.db.Model(&models.PolicyReportOverview{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&overview).Error

	if err != nil {
		return nil, err
	}
	return &overview, nil
}

func (repo *PolicyReportRepository) GetPolicyReportKeyMetrics(policyReportID uint) (*policyDTOs.PolicyReportKeyMetricsResponseDTO, error) {
	var keyMetrics policyDTOs.PolicyReportKeyMetricsResponseDTO
	err := repo.db.Model(&models.PolicyReportKeyMetrics{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&keyMetrics).Error

	if err != nil {
		return nil, err
	}
	return &keyMetrics, nil
}

func (repo *PolicyReportRepository) GetPolicyReportClause(policyReportID uint) ([]*policyDTOs.PolicyReportClauseResponseDTO, error) {
	clauses := []*policyDTOs.PolicyReportClauseResponseDTO{}
	err := repo.db.Model(&models.PolicyReportClause{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&clauses).Error

	if err != nil {
		return nil, err
	}
	return clauses, nil
}

func (repo *PolicyReportRepository) GetPolicyReportAcknowledgmentUser(policyReportID uint) ([]*policyDTOs.PolicyReportAcknowledgementUserResponseDTO, error) {
	users := []*policyDTOs.PolicyReportAcknowledgementUserResponseDTO{}
	err := repo.db.Model(&models.PolicyReportAcknowledgementUser{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *PolicyReportRepository) GetPolicyReportAcknowledgmentOrganizationUnit(policyReportID uint) ([]*policyDTOs.PolicyReportAcknowledgementOrganizationUnitResponseDTO, error) {
	organizationUnits := []*policyDTOs.PolicyReportAcknowledgementOrganizationUnitResponseDTO{}
	err := repo.db.Model(&models.PolicyReportAcknowledgementOrganizationUnit{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&organizationUnits).Error

	if err != nil {
		return nil, err
	}
	return organizationUnits, nil
}

func (repo *PolicyReportRepository) GetPolicyReportRisk(policyReportID uint) ([]*policyDTOs.PolicyReportRiskResponseDTO, error) {
	risks := []*policyDTOs.PolicyReportRiskResponseDTO{}
	err := repo.db.Model(&models.PolicyReportRisk{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&risks).Error

	if err != nil {
		return nil, err
	}
	return risks, nil
}

func (repo *PolicyReportRepository) GetPolicyReportStrategicObjective(policyReportID uint) ([]*policyDTOs.PolicyReportStrategicObjectiveResponseDTO, error) {
	strategicObjectives := []*policyDTOs.PolicyReportStrategicObjectiveResponseDTO{}
	err := repo.db.Model(&models.PolicyReportStrategicObjective{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&strategicObjectives).Error

	if err != nil {
		return nil, err
	}
	return strategicObjectives, nil
}

func (repo *PolicyReportRepository) GetPolicyReportCompliance(policyReportID uint) ([]*policyDTOs.PolicyReportComplianceResponseDTO, error) {
	compliances := []*policyDTOs.PolicyReportComplianceResponseDTO{}
	err := repo.db.Model(&models.PolicyReportCompliance{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&compliances).Error

	if err != nil {
		return nil, err
	}
	return compliances, nil
}

func (repo *PolicyReportRepository) GetPolicyReportDocumentVersion(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentVersionResponseDTO, error) {
	documentVersions := []*policyDTOs.PolicyReportDocumentVersionResponseDTO{}
	err := repo.db.Model(&models.PolicyReportDocumentVersion{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&documentVersions).Error

	if err != nil {
		return nil, err
	}
	return documentVersions, nil
}

func (repo *PolicyReportRepository) GetPolicyReportDocumentAccess(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentAccessResponseDTO, error) {
	documentAccesses := []*policyDTOs.PolicyReportDocumentAccessResponseDTO{}
	err := repo.db.Model(&models.PolicyReportDocumentAccess{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&documentAccesses).Error

	if err != nil {
		return nil, err
	}
	return documentAccesses, nil
}

func (repo *PolicyReportRepository) GetPolicyReportDocumentApprovalHistory(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentApprovalHistoryResponseDTO, error) {
	approvalHistories := []*policyDTOs.PolicyReportDocumentApprovalHistoryResponseDTO{}
	err := repo.db.Model(&models.PolicyReportDocumentApprovalHistory{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&approvalHistories).Error

	if err != nil {
		return nil, err
	}
	return approvalHistories, nil
}

func (repo *PolicyReportRepository) GetPolicyReportDocumentAttachment(policyReportID uint) ([]*policyDTOs.PolicyReportDocumentAttachmentResponseDTO, error) {
	attachments := []*policyDTOs.PolicyReportDocumentAttachmentResponseDTO{}
	err := repo.db.Model(&models.PolicyReportDocumentAttachment{}).
		Where("policy_report_id = ?", policyReportID).
		Scan(&attachments).Error

	if err != nil {
		return nil, err
	}
	return attachments, nil
}
