package repositories

import (
	"fmt"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	thirdPartyDTOs "jethur/grc/tprm/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThirdPartyRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *thirdPartyDTOs.ThirdPartyRegisterListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyRegister, error)
	Create(register *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(register *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type ThirdPartyRegisterRepository struct {
	db *gorm.DB
}

func NewThirdPartyRegisterRepository(db *gorm.DB) *ThirdPartyRegisterRepository {
	return &ThirdPartyRegisterRepository{db: db}
}

func (repo *ThirdPartyRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyRegisterRepository) GetAll(params *thirdPartyDTOs.ThirdPartyRegisterListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyRegisterListResponseDTO, error) {
	var results []*thirdPartyDTOs.ThirdPartyRegisterListDTO
	var totalCount int64

	query := repo.db.Table("third_party_register").
		Select(`
			third_party_register.*,
			CONVERT(VARCHAR, onboarding_date, 23) AS onboarding_date,
			CONVERT(VARCHAR, service_start_date, 23) AS service_start_date,
			CONVERT(VARCHAR, service_end_date, 23) AS service_end_date,
			countries.id AS country_id,
			countries.title_en AS country_title_en,
			countries.title_ar AS country_title_ar,
			third_party_types.title_en AS third_party_type_title_en,
			third_party_types.title_ar AS third_party_type_title_ar,
			risk_levels.title_en AS risk_level_title_en,
        	risk_levels.title_ar AS risk_level_title_ar,
        	risk_levels.color_code AS risk_level_color_code,
        	risk_levels.type AS risk_level_type,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.last_name_en AS created_user_last_name_en,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.last_name_en AS updated_user_last_name_en
		`).
		Joins("LEFT JOIN third_party_types ON third_party_types.id = third_party_register.third_party_type_id").
		Joins("LEFT JOIN countries ON countries.id = third_party_register.country_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = third_party_register.organization_unit_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = third_party_register.severity_level_id").
		Joins("LEFT JOIN risk_levels ON risk_levels.id = third_party_register.risk_level_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = third_party_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = third_party_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = third_party_register.updated_by_id").
		Joins(`
        LEFT JOIN (
        SELECT 
            tpa.third_party_register_id,
            tpa.third_party_assessment_status_id,
            tpa.risk_level_id,
			tpa.score as score,
            tpa.created_at,
            ROW_NUMBER() OVER (
                PARTITION BY tpa.third_party_register_id
                ORDER BY tpa.created_at DESC
            ) AS rn
        FROM third_party_assessments tpa
        ) AS latest
        ON latest.third_party_register_id = third_party_register.id
        AND latest.rn = 1
        `)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "latest_status_id" {
				query = query.Where("latest.third_party_assessment_status_id IN (?)", values)
				continue
			}

			if columnName == "latest_risk_level_id" {
				query = query.Where("latest.risk_level_id IN (?)", values)
				continue
			}

			if columnName == "organization_unit_id" {
				query = query.Where("third_party_register.organization_unit_id IN (?)", values)
			} else {
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	now := time.Now()
	exp30 := now.AddDate(0, 0, 30)
	exp60 := now.AddDate(0, 0, 60)

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "low":
			query = query.Where("severity_levels.type = ?", "low")
		case "medium":
			query = query.Where("severity_levels.type = ?", "medium")
		case "high":
			query = query.Where("severity_levels.type = ?", "high")
		case "expired":
			query = query.Where("third_party_register.service_end_date < ?", now)
		case "expiring_30":
			query = query.Where("third_party_register.service_end_date BETWEEN ? AND ?", now, exp30)
		case "expiring_60":
			query = query.Where("third_party_register.service_end_date BETWEEN ? AND ?", exp30, exp60)
		}

	}

	if params.Keyword != "" {
		likePattern := "%" + params.Keyword + "%"
		query = query.Where(`
			third_party_register.title LIKE ? OR
			third_party_register.reference_number LIKE ? OR
			third_party_register.representative_name LIKE ? OR
			third_party_register.representative_role LIKE ? OR
			third_party_register.email LIKE ? OR
			third_party_register.service_start_date LIKE ? OR
			third_party_register.service_end_date LIKE ? OR
	
			countries.title_en LIKE ? OR
			countries.title_ar LIKE ? OR
	
			third_party_types.title_en LIKE ? OR
			third_party_types.title_ar LIKE ? OR
	
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
	
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
	
			responsible_user.first_name_en LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.email LIKE ?
		`,
			likePattern, likePattern,
			likePattern, likePattern,
			likePattern, likePattern,
			likePattern, likePattern,
			likePattern, likePattern, likePattern,
			likePattern, likePattern, likePattern,
			likePattern, likePattern, likePattern, likePattern,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &thirdPartyDTOs.ThirdPartyRegisterListResponseDTO{
		Items: results,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ThirdPartyRegisterRepository) GetByID(id uint) (*models.ThirdPartyRegister, error) {
	var entity *models.ThirdPartyRegister
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ThirdPartyRegister
	if err := db.Preload("ThirdPartyType").
		Preload("Country").
		Preload("OrgUnit").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("RiskLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ThirdPartyRegisterRepository) Create(entity *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ThirdPartyRegisterRepository) Update(entity *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(entity).Error
}

func (repo *ThirdPartyRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ThirdPartyRegister{}, id).Error
}

func (repo *ThirdPartyRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReference string
	if err := repo.db.Model(&models.ThirdPartyRegister{}).Select("reference_number").Order("id desc").Limit(1).Pluck("reference_number", &lastReference).Error; err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReference != "" && len(lastReference) >= 4 {
		lastDigits = lastReference[len(lastReference)-4:]
	}
	lastNum, err := strconv.Atoi(lastDigits)
	if err != nil {
		return "", err
	}

	ref := fmt.Sprintf("TPR-%d-%04d", time.Now().Year(), lastNum+1)
	return ref, nil
}

func (repo *ThirdPartyRegisterRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	/* ================= ACTION PLAN TREE ================= */

	queries := []string{

		`DELETE tprm_assessment_section_question_action_plan_status_update_documents
		 FROM tprm_assessment_section_question_action_plan_status_update_documents 
		 JOIN tprm_assessment_section_question_action_plan_status_updates ON tprm_assessment_section_question_action_plan_status_updates.id = tprm_assessment_section_question_action_plan_status_update_documents.tprm_assessment_section_question_action_plan_status_update_id
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_updates
		 FROM tprm_assessment_section_question_action_plan_status_updates 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_comments
		 FROM tprm_assessment_section_question_action_plan_comments 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_comments.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tasks
		 FROM tasks 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tasks.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_update_documents
		 FROM tprm_assessment_section_question_action_plan_status_update_documents 
		 JOIN tprm_assessment_section_question_action_plan_status_updates ON tprm_assessment_section_question_action_plan_status_updates.id = tprm_assessment_section_question_action_plan_status_update_documents.tprm_assessment_section_question_action_plan_status_update_id
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_question_answer_options ON tprm_assessment_question_answer_options.id = tprm_assessment_answer_options.answer_option_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_options.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_comment_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_comments
		 FROM tprm_assessment_section_question_comments 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_comments.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plans
		 FROM tprm_assessment_section_question_action_plans 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_question_answer_options
		 FROM tprm_assessment_question_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answers
		 FROM tprm_assessment_answers 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_documents.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_comment_documents.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_updates
		 FROM tprm_assessment_section_question_action_plan_status_updates 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_comments
		 FROM tprm_assessment_section_question_action_plan_comments 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_comments.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tasks
		 FROM tasks 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tasks.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_question_answer_options ON tprm_assessment_question_answer_options.id = tprm_assessment_answer_options.answer_option_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_options.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_comment_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE stakeholder_action_plan_status_update_documents
		 FROM stakeholder_action_plan_status_update_documents 
		 JOIN stakeholder_action_plan_status_updates ON stakeholder_action_plan_status_updates.id = stakeholder_action_plan_status_update_documents.stakeholder_action_plan_status_update_id
		 JOIN stakeholder_action_plans ON stakeholder_action_plans.id = stakeholder_action_plan_status_updates.stakeholder_action_plan_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_action_plans.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_action_plan_comment_documents
		 FROM stakeholder_action_plan_comment_documents 
		 JOIN stakeholder_action_plan_comments ON stakeholder_action_plan_comments.id = stakeholder_action_plan_comment_documents.stakeholder_action_plan_comment_id
		 JOIN stakeholder_action_plans ON stakeholder_action_plans.id = stakeholder_action_plan_comments.stakeholder_action_plan_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_action_plans.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE tprm_assessment_workflow_reviews
		 FROM tprm_assessment_workflow_reviews 
		 JOIN tprm_assessment_workflows ON tprm_assessment_workflows.id = tprm_assessment_workflow_reviews.tprm_assessment_workflow_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_workflows.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_workflow_participants
		 FROM tprm_assessment_workflow_participants 
		 JOIN tprm_assessment_workflows ON tprm_assessment_workflows.id = tprm_assessment_workflow_participants.tprm_assessment_workflow_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_workflows.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_questions
		 FROM tprm_assessment_section_questions 
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_comments
		 FROM tprm_assessment_section_question_comments 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_comments.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_action_plans
		 FROM tprm_assessment_section_question_action_plans 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_question_answer_options
		 FROM tprm_assessment_question_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answers
		 FROM tprm_assessment_answers 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_options.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_documents.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_comment_documents.tprm_assessment_question_id
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE stakeholder_needs_and_expectation_evaluations
		 FROM stakeholder_needs_and_expectation_evaluations 
		 JOIN stakeholder_needs_and_expectations ON stakeholder_needs_and_expectations.id = stakeholder_needs_and_expectation_evaluations.stakeholder_needs_and_expectation_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_needs_and_expectations.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_engagement_evaluations
		 FROM stakeholder_engagement_evaluations 
		 JOIN stakeholder_engagements ON stakeholder_engagements.id = stakeholder_engagement_evaluations.stakeholder_engagement_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_engagements.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_action_plan_status_updates
		 FROM stakeholder_action_plan_status_updates 
		 JOIN stakeholder_action_plans ON stakeholder_action_plans.id = stakeholder_action_plan_status_updates.stakeholder_action_plan_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_action_plans.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_action_plan_comments
		 FROM stakeholder_action_plan_comments 
		 JOIN stakeholder_action_plans ON stakeholder_action_plans.id = stakeholder_action_plan_comments.stakeholder_action_plan_id
		 JOIN stakeholders ON stakeholders.id = stakeholder_action_plans.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE tprm_assessment_workflows
		 FROM tprm_assessment_workflows 
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_workflows.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_sections
		 FROM tprm_assessment_sections 
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_sections.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_questions
		 FROM tprm_assessment_section_questions 
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_questions.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE tprm_assessment_section_question_comments
		 FROM tprm_assessment_section_question_comments 
		 JOIN third_party_assessments ON third_party_assessments.id = tprm_assessment_section_question_comments.third_party_assessment_id
		 WHERE third_party_assessments.third_party_register_id = ?`,

		`DELETE stakeholder_needs_and_expectations
		 FROM stakeholder_needs_and_expectations 
		 JOIN stakeholders ON stakeholders.id = stakeholder_needs_and_expectations.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_engagements
		 FROM stakeholder_engagements 
		 JOIN stakeholders ON stakeholders.id = stakeholder_engagements.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE stakeholder_action_plans
		 FROM stakeholder_action_plans 
		 JOIN stakeholders ON stakeholders.id = stakeholder_action_plans.stakeholder_id
		 WHERE stakeholders.third_party_register_id = ?`,

		`DELETE FROM third_party_strategic_objectives WHERE third_party_strategic_objectives.third_party_register_id = ?`,
		`DELETE FROM third_party_risks WHERE third_party_risks.third_party_register_id = ?`,
		`DELETE FROM third_party_processes WHERE third_party_processes.third_party_register_id = ?`,
		`DELETE FROM third_party_incidents WHERE third_party_incidents.third_party_register_id = ?`,
		`DELETE FROM third_party_controls WHERE third_party_controls.third_party_register_id = ?`,
		`DELETE FROM third_party_compliance_requirements WHERE third_party_compliance_requirements.third_party_register_id = ?`,
		`DELETE FROM third_party_associated_documents WHERE third_party_associated_documents.third_party_register_id = ?`,
		`DELETE FROM third_party_assets WHERE third_party_assets.third_party_register_id = ?`,
		`DELETE FROM third_party_assessments WHERE third_party_assessments.third_party_register_id = ?`,
		`DELETE FROM stakeholders WHERE stakeholders.third_party_register_id = ?`,
		`DELETE FROM process_dependencies WHERE process_dependencies.third_party_register_id = ?`,
		`DELETE FROM kpi_third_parties WHERE kpi_third_parties.third_party_register_id = ?`,
	}

	for _, q := range queries {
		if err := db.Exec(q, id).Error; err != nil {
			return err
		}
	}

	/* ================= ROOT ================= */
	return db.Delete(&models.ThirdPartyRegister{}, id).Error
}
