package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceEvidenceSubmissionRepositoryInterface interface {
	GetAll(params *complianceDTOs.ComplianceEvidenceSubmissionListQueryParamsDTO, currentUserID uint) (*complianceDTOs.ComplianceEvidenceSubmissionListResponseDTO, error)
	GetRequirementAssessmentFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementResponsibleUnit, error)
}

type ComplianceEvidenceSubmissionRepository struct {
	db *gorm.DB
}

func NewComplianceEvidenceSubmissionRepository(db *gorm.DB) ComplianceEvidenceSubmissionRepositoryInterface {
	return &ComplianceEvidenceSubmissionRepository{db: db}
}

func (repo *ComplianceEvidenceSubmissionRepository) GetAll(params *complianceDTOs.ComplianceEvidenceSubmissionListQueryParamsDTO, currentUserID uint) (*complianceDTOs.ComplianceEvidenceSubmissionListResponseDTO, error) {
	var items []*complianceDTOs.ComplianceEvidenceSubmissionListDTO
	var total int64

	baseQuery := `
	(
		SELECT
			carru.id,
			carru.is_completed,
			carru.is_acknowledged,
			ca.id AS compliance_assessment_id,
			clauses.id AS clause_id,
			clauses.clause_no AS compliance_clause_no,
			clauses.title_en AS compliance_clause_title_en,
			clauses.title_ar AS compliance_clause_title_ar,
			requirements.id AS compliance_requirement_id,
			requirements.title_en AS compliance_requirement_title_en,
			requirements.title_ar AS compliance_requirement_title_ar,
			registers.title_en AS compliance_register_title_en,
			registers.title_ar AS compliance_register_title_ar,
			req_statuses.id AS compliance_assessment_requirement_status_id,
			req_statuses.title_en AS compliance_assessment_requirement_status_title_en,
			req_statuses.title_ar AS compliance_assessment_requirement_status_title_ar,
			req_statuses.color_code AS compliance_assessment_requirement_status_color_code,
			req_statuses.type AS compliance_assessment_requirement_status_type,
			org_units.id AS responsible_unit_id,
			org_units.title_en AS responsible_unit_title_en,
			org_units.title_ar AS responsible_unit_title_ar,
			responsible_users.id AS responsible_user_id,
			responsible_users.first_name_en AS responsible_user_first_name_en,
			responsible_users.first_name_ar AS responsible_user_first_name_ar,
			responsible_users.last_name_en AS responsible_user_last_name_en,
			responsible_users.last_name_ar AS responsible_user_last_name_ar,
			responsible_users.initial AS responsible_user_initial,
			responsible_users.color AS responsible_user_color,
			responsible_users.email AS responsible_user_email,
			responsible_users.profile_image_path AS responsible_user_profile_image_path,
			completed_users.id AS completed_user_id,
			completed_users.first_name_en AS completed_user_first_name_en,
			completed_users.first_name_ar AS completed_user_first_name_ar,
			completed_users.last_name_en AS completed_user_last_name_en,
			completed_users.last_name_ar AS completed_user_last_name_ar,
			completed_users.initial AS completed_user_initial,
			completed_users.color AS completed_user_color,
			completed_users.email AS completed_user_email,
			completed_users.profile_image_path AS completed_user_profile_image_path,
			CONVERT(VARCHAR, carru.created_at, 23) AS created_at,
			CONVERT(VARCHAR, carru.updated_at, 23) AS updated_at,
			CONVERT(VARCHAR, carru.completed_at, 23) AS completed_at,
			carru.is_risk_based_assessment,
			carru.is_regulation_notification_assessment,
			ca.review_user_id AS review_user_id,
			registers.responsible_user_id AS register_responsible_user_id
		FROM compliance_assessment_requirement_responsible_units carru
		LEFT JOIN compliance_assessment_requirements car
			ON car.id = carru.compliance_assessment_requirement_id
		LEFT JOIN compliance_assessments ca
			ON ca.id = car.compliance_assessment_id
		LEFT JOIN compliance_requirements requirements
			ON requirements.id = car.compliance_requirement_id
		LEFT JOIN compliance_clauses clauses
			ON clauses.id = requirements.compliance_clause_id
		LEFT JOIN compliance_register registers
			ON registers.id = clauses.compliance_register_id
		LEFT JOIN compliance_assessment_requirement_statuses req_statuses
			ON req_statuses.id = carru.compliance_assessment_requirement_status_id
		LEFT JOIN organization_units org_units
			ON org_units.id = carru.responsible_unit_id
		LEFT JOIN users responsible_users
			ON responsible_users.id = carru.responsible_user_id
		LEFT JOIN users completed_users
			ON completed_users.id = carru.completed_by_user_id
		WHERE (carru.is_risk_based_assessment = 1 OR carru.is_regulation_notification_assessment = 1)
	) AS merged_data`

	query := repo.db.Table(baseQuery)

	if params.IsPassiveUser {
		query = query.Where("responsible_user_id = ?", currentUserID)
	} else if !params.IsSuperAdmin {
		query = query.Where(
			"responsible_user_id = ? OR review_user_id = ? OR register_responsible_user_id = ?",
			currentUserID, currentUserID, currentUserID,
		)
	}

	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(column+" IN (?)", values)
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "completed":
			query = query.Where("is_completed = ?", true)
		case "pending":
			query = query.Where("is_completed = ?", false)
		case "acknowledged":
			query = query.Where("is_acknowledged = ?", true)
		case "not_acknowledged":
			query = query.Where("is_acknowledged = ?", false)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(
			LOWER(compliance_clause_title_en) LIKE LOWER(?) OR
			LOWER(compliance_clause_title_ar) LIKE LOWER(?) OR
			LOWER(compliance_register_title_en) LIKE LOWER(?) OR
			LOWER(compliance_register_title_ar) LIKE LOWER(?) OR
			LOWER(compliance_requirement_title_en) LIKE LOWER(?) OR
			LOWER(compliance_requirement_title_ar) LIKE LOWER(?)
		)`, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Order(fmt.Sprintf("%s %s", params.SortBy, params.Sort)).
			Offset(offset).
			Limit(params.PageSize)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &complianceDTOs.ComplianceEvidenceSubmissionListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  total,
		},
	}, nil
}

func (repo *ComplianceEvidenceSubmissionRepository) GetRequirementAssessmentFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementResponsibleUnit, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var responsibleUnit *models.ComplianceAssessmentRequirementResponsibleUnit

	result := db.
		Preload("ComplianceAssessmentRequirement").
		Preload("ComplianceAssessmentRequirement.ComplianceAssessment").
		Preload("ComplianceAssessmentRequirement.ComplianceAssessment.ComplianceAssessmentStatus").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ResponsibleUser").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ResponsibleUnit").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ComplianceCategory").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.SeverityLevel").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.LegislativeBody").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.LegislationType").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ImportanceClassification").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.SecrecyClassification").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.IssuingAuthorities").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.IssuingAuthorities.IssuingAuthority").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.IssuingAuthorities.Entity").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.IssuingAuthorities.Document").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.CreatedByUser").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ComplianceRegisterSharedUnit").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.AverageResidualRiskRating").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ComplianceRegisterSharedUnit.Department").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceClause.ComplianceRegister.ComplianceAssessmentFramework.Values.ComplianceStatus").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.SeverityLevel").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceRequirementCategory").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement.ComplianceStatus").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("CompletedByUser").
		Preload("ComplianceAssessmentRequirementStatus").
		Preload("ComplianceAssessmentFrameworkValue").
		Preload("ComplianceAssessmentFrameworkValue.ComplianceAssessmentFramework").
		Preload("ComplianceAssessmentFrameworkValue.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceAssessmentFrameworkValue.ComplianceStatus").
		Preload("RequiredDocumentTitles").
		Preload("RequiredDocumentTitles.RequiredDocumentTitle").
		Preload("RequiredDocumentTitles.Documents").
		Preload("RequiredDocumentTitles.Documents.ApprovedByUser").
		Preload("RequiredDocumentTitles.Documents.RejectedByUser").
		Preload("ActionPlans").
		Preload("ActionPlans.ActionPlanStatus").
		Preload("ActionPlans.ResponsibleUser").
		Preload("ActionPlans.ReviewUser").
		Preload("RequirementComments").
		Preload("RequirementComments.ComplianceAssessmentRequirementStatus").
		Preload("RequirementComments.CreatedByUser").
		First(&responsibleUnit, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return responsibleUnit, nil
}
