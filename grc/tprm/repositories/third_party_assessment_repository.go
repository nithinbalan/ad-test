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

type ThirdPartyAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *thirdPartyDTOs.ThirdPartyAssessmentListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyAssessment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error)
	Create(assessments *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(assessments *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateAssessmentSection(section *models.TPRMAssessmentSection, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAssessmentQuestion(question *models.TPRMAssessmentSectionQuestion, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAssessmentQuestionAnswerOption(option *models.TPRMAssessmentQuestionAnswerOption, tx ...utility.DatabaseTransaction) (uint, error)

	DeleteSectionsByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error
	DeleteQuestionsByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error
	DeleteAnswerOptionsByQuestionID(questionID uint, tx ...utility.DatabaseTransaction) error

	CreateAssessmentAnswer(assessments *models.TPRMAssessmentAnswer, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAssessmentAnswerOption(assessments *models.TPRMAssessmentAnswerOption, tx ...utility.DatabaseTransaction) (uint, error)

	GetAnswerSummary(assessmentID uint) ([]thirdPartyDTOs.AnswerSummaryItemDTO, error)
	GetQuestionByID(id uint) (*models.TPRMAssessmentSectionQuestion, error)
	UpdateAnswer(answer *models.TPRMAssessmentAnswer, tx ...utility.DatabaseTransaction) error
	DeleteAnswerOptionsByAnswerID(answer uint, tx ...utility.DatabaseTransaction) error
	GetAnswerOptionByQuestionID(questionID uint, complianceResponseID uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentQuestionAnswerOption, error)
	GetAnswerOptionByID(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentQuestionAnswerOption, error)

	GetAssessmentSectionDetails(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentSection, error)
	GetAssessmentSectionQuestionDetails(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentSectionQuestion, error)

	GetAllAnswerDocuments(answerID uint, tx ...utility.DatabaseTransaction) ([]thirdPartyDTOs.AnswerDocumentListDTO, error)
	CreateAnswerDocument(item *models.TPRMAssessmentAnswerDocument, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteAnswerDocument(id uint, tx ...utility.DatabaseTransaction) error
}

type ThirdPartyAssessmentRepository struct {
	db *gorm.DB
}

func NewThirdPartyAssessmentRepository(db *gorm.DB) *ThirdPartyAssessmentRepository {
	return &ThirdPartyAssessmentRepository{db: db}
}

func (repo *ThirdPartyAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyAssessmentRepository) GetAll(params *thirdPartyDTOs.ThirdPartyAssessmentListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO, error) {
	var results []*thirdPartyDTOs.ThirdPartyAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("third_party_assessments").
		Select(`
			third_party_assessments.*,
			CONVERT(VARCHAR, due_date, 23) AS due_date,
			question_sets.id AS question_set_id,
			question_sets.title_en AS question_set_title_en,
			question_sets.title_ar AS question_set_title_ar,
			third_party_register.title AS third_party_register_title,
			third_party_assessment_statuses.title_en AS third_party_assessment_status_title_en,
			third_party_assessment_statuses.title_ar AS third_party_assessment_status_title_ar,
			third_party_assessment_statuses.type AS third_party_assessment_status_type,
			third_party_assessment_statuses.color_code AS third_party_assessment_status_color_code,
			risk_levels.title_en AS risk_level_title_en,
			risk_levels.title_ar AS risk_level_title_ar,
			risk_levels.color_code AS risk_level_color_code,
			risk_levels.type AS risk_level_type,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.last_name_en AS created_user_last_name_en,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.last_name_en AS updated_user_last_name_en
		`).
		Joins("LEFT JOIN third_party_register ON third_party_register.id = third_party_assessments.third_party_register_id").
		Joins("LEFT JOIN question_sets ON question_sets.id = third_party_assessments.question_set_id").
		Joins("LEFT JOIN risk_levels ON risk_levels.id = third_party_assessments.risk_level_id").
		Joins("LEFT JOIN third_party_assessment_statuses ON third_party_assessment_statuses.id = third_party_assessments.third_party_assessment_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = third_party_assessments.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = third_party_assessments.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "third_party_type_id" {
				query = query.Where("third_party_register.third_party_type_id IN (?)", values)
				continue
			}
			if columnName == "risk_level_id" {
				query = query.Where("third_party_assessments.risk_level_id IN (?)", values)
				continue
			}
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	today := time.Now().Truncate(24 * time.Hour)
	next7 := today.AddDate(0, 0, 7)
	next14 := today.AddDate(0, 0, 14)

	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("CAST(third_party_assessments.due_date AS date) < CAST(? AS date)", today)
		case "upcoming_7":
			query = query.Where("CAST(third_party_assessments.due_date AS date) BETWEEN CAST(? AS date) AND CAST(? AS date)", today, next7)
		case "upcoming_14":
			query = query.Where("CAST(third_party_assessments.due_date AS date) BETWEEN CAST(? AS date) AND CAST(? AS date)", next7, next14)
		}
	}

	if params.Keyword != "" {
		likePattern := "%" + params.Keyword + "%"
		query = query.Where(`
			third_party_assessments.title LIKE ? OR
			third_party_assessments.reference_number LIKE ? OR
			third_party_assessments.due_date LIKE ? OR
			third_party_register.title LIKE ? OR
			question_sets.title_en LIKE ? OR
			question_sets.title_ar LIKE ?
		`,
			likePattern, likePattern,
			likePattern, likePattern,
			likePattern, likePattern,
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

	if params.Export {
		return &thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO{
			Items: results,
		}, nil
	}

	response := &thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO{
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

func (repo *ThirdPartyAssessmentRepository) GetByID(id uint) (*models.ThirdPartyAssessment, error) {
	var entity *models.ThirdPartyAssessment
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyAssessmentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ThirdPartyAssessment
	if err := db.Preload("ThirdPartyRegister").
		Preload("ThirdPartyRegister.ThirdPartyType").
		Preload("ThirdPartyRegister.Country").
		Preload("ThirdPartyRegister.OrgUnit").
		Preload("ThirdPartyRegister.ResponsibleUser").
		Preload("ThirdPartyRegister.ResponsibleUnit").
		Preload("ThirdPartyRegister.SeverityLevel").
		Preload("ThirdPartyAssessmentStatus").
		Preload("ThirdPartyAssessmentWorkflows").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowParticipants").
		Preload("ThirdPartyAssessmentWorkflows.ThirdPartyAssessmentWorkflowReviews").
		Preload("Sections").
		Preload("Sections.Questions").
		// Preload("Sections.Questions.QuestionLibrary.QuestionType").
		// Preload("Sections.Questions.AnswerOptions").
		// Preload("Sections.Questions.AnswerOptions.RiskLevel").
		// Preload("Sections.Questions.AnswerOptions.TPRMComplianceResponse").
		// Preload("Sections.Questions.Answer").
		// Preload("Sections.Questions.Answer.AnswerOptions.AnswerOption").
		// Preload("Sections.Questions.Answer.AnswerOptions.AnswerOption.RiskLevel").
		// Preload("Sections.Questions.Answer.AnswerOptions.AnswerOption.TPRMComplianceResponse").
		Preload("QuestionSet").
		Preload("RiskLevel").
		//Preload("Questions").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("SubmittedByUser").
		Preload("SentByUser").
		Preload("ImportedByUser").
		First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ThirdPartyAssessmentRepository) Create(entity *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) Update(entity *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(entity).Error
}

func (repo *ThirdPartyAssessmentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ThirdPartyAssessment{}, id).Error
}

func (repo *ThirdPartyAssessmentRepository) GenerateReferenceNumber() (string, error) {
	var lastReference string
	if err := repo.db.Model(&models.ThirdPartyAssessment{}).Select("reference_number").Order("id desc").Limit(1).Pluck("reference_number", &lastReference).Error; err != nil {
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

	ref := fmt.Sprintf("TPA-%d-%04d", time.Now().Year(), lastNum+1)
	return ref, nil
}

func (repo *ThirdPartyAssessmentRepository) CreateAssessmentSection(section *models.TPRMAssessmentSection, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(section).Error; err != nil {
		return 0, err
	}
	return section.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) CreateAssessmentQuestion(question *models.TPRMAssessmentSectionQuestion, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(question).Error; err != nil {
		return 0, err
	}
	return question.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) CreateAssessmentQuestionAnswerOption(option *models.TPRMAssessmentQuestionAnswerOption, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(option).Error; err != nil {
		return 0, err
	}
	return option.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) DeleteAnswerOptionsByQuestionID(questionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("tprm_assessment_question_id = ?", questionID).Delete(&models.TPRMAssessmentQuestionAnswerOption{}).Error
}

func (repo *ThirdPartyAssessmentRepository) DeleteQuestionsByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("third_party_assessment_id = ?", assessmentID).Delete(&models.TPRMAssessmentSectionQuestion{}).Error
}

func (repo *ThirdPartyAssessmentRepository) DeleteSectionsByAssessmentID(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("third_party_assessment_id = ?", assessmentID).Delete(&models.TPRMAssessmentSection{}).Error
}

func (repo *ThirdPartyAssessmentRepository) CreateAssessmentAnswer(answer *models.TPRMAssessmentAnswer, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(answer).Error; err != nil {
		return 0, err
	}
	return answer.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) CreateAssessmentAnswerOption(option *models.TPRMAssessmentAnswerOption, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(option).Error; err != nil {
		return 0, err
	}
	return option.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) GetAnswerSummary(assessmentID uint) ([]thirdPartyDTOs.AnswerSummaryItemDTO, error) {
	var results []thirdPartyDTOs.AnswerSummaryItemDTO

	query := `
       SELECT
			cr.title_en AS title_en,
			cr.title_ar AS title_ar,
			cr.color_code AS color_code,
			COUNT(ao.id) AS count
		FROM tprm_compliance_responses cr
		LEFT JOIN tprm_assessment_question_answer_options qo
			ON qo.tprm_compliance_response_id = cr.id
		LEFT JOIN tprm_assessment_answer_options ao
			ON ao.answer_option_id = qo.id
		AND ao.tprm_assessment_question_id IN (
				SELECT id
				FROM tprm_assessment_section_questions
				WHERE third_party_assessment_id = ?
		)
		GROUP BY cr.id, cr.title_en, cr.title_ar,cr.color_code
		ORDER BY cr.title_en;
    `

	err := repo.db.Raw(query, assessmentID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ThirdPartyAssessmentRepository) GetQuestionByID(id uint) (*models.TPRMAssessmentSectionQuestion, error) {
	var question *models.TPRMAssessmentSectionQuestion
	if err := repo.db.Preload("Answer").Preload("Answer.Documents").First(&question, id).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (repo *ThirdPartyAssessmentRepository) UpdateAnswer(answer *models.TPRMAssessmentAnswer, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(answer).Error
}

func (repo *ThirdPartyAssessmentRepository) DeleteAnswerOptionsByAnswerID(answerID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("tprm_assessment_answer_id = ?", answerID).Delete(&models.TPRMAssessmentAnswerOption{}).Error
}

func (repo *ThirdPartyAssessmentRepository) GetAnswerOptionByID(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentQuestionAnswerOption, error) {
	var question *models.TPRMAssessmentQuestionAnswerOption
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.First(&question, id).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (repo *ThirdPartyAssessmentRepository) GetAnswerOptionByQuestionID(questionID uint, complianceResponseID uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentQuestionAnswerOption, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var answerOption models.TPRMAssessmentQuestionAnswerOption

	result := db.Where("tprm_assessment_question_id = ? AND tprm_compliance_response_id = ?", questionID, complianceResponseID).First(&answerOption)
	if result.Error != nil {
		return nil, result.Error
	}

	return &answerOption, nil
}

func (repo *ThirdPartyAssessmentRepository) GetAssessmentSectionDetails(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section *models.TPRMAssessmentSection
	if err := db.Preload("Questions").
		Preload("Questions.QuestionLibrary.QuestionType").
		First(&section, id).Error; err != nil {
		return nil, err
	}

	return section, nil
}

func (repo *ThirdPartyAssessmentRepository) GetAssessmentSectionQuestionDetails(id uint, tx ...utility.DatabaseTransaction) (*models.TPRMAssessmentSectionQuestion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section *models.TPRMAssessmentSectionQuestion
	if err := db.Preload("AnswerOptions").
		Preload("AnswerOptions.RiskLevel").
		Preload("AnswerOptions.TPRMComplianceResponse").
		Preload("Answer").
		Preload("Answer.AnswerOptions.AnswerOption").
		Preload("Answer.AnswerOptions.AnswerOption.RiskLevel").
		Preload("Answer.AnswerOptions.AnswerOption.TPRMComplianceResponse").
		First(&section, id).Error; err != nil {
		return nil, err
	}

	return section, nil
}

func (repo *ThirdPartyAssessmentRepository) GetAllAnswerDocuments(answerID uint, tx ...utility.DatabaseTransaction) ([]thirdPartyDTOs.AnswerDocumentListDTO, error) {
	var result []thirdPartyDTOs.AnswerDocumentListDTO

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	query := db.Table("tprm_assessment_answer_documents").
		Select(`
			tprm_assessment_answer_documents.*`).
		Where("tprm_assessment_answer_documents.tprm_assessment_answer_id = ?", answerID)

	err := query.Find(&result).Error
	return result, err
}

func (repo *ThirdPartyAssessmentRepository) CreateAnswerDocument(answerDocument *models.TPRMAssessmentAnswerDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(answerDocument).Error
	return answerDocument.ID, err
}

func (repo *ThirdPartyAssessmentRepository) DeleteAnswerDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.TPRMAssessmentAnswerDocument{}, id).Error
}

func (repo *ThirdPartyAssessmentRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

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
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_updates
		 FROM tprm_assessment_section_question_action_plan_status_updates 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_comments
		 FROM tprm_assessment_section_question_action_plan_comments 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_comments.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tasks
		 FROM tasks 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tasks.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_update_documents
		 FROM tprm_assessment_section_question_action_plan_status_update_documents 
		 JOIN tprm_assessment_section_question_action_plan_status_updates ON tprm_assessment_section_question_action_plan_status_updates.id = tprm_assessment_section_question_action_plan_status_update_documents.tprm_assessment_section_question_action_plan_status_update_id
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_question_answer_options ON tprm_assessment_question_answer_options.id = tprm_assessment_answer_options.answer_option_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_options.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_comment_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_comments
		 FROM tprm_assessment_section_question_comments 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_comments.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plans
		 FROM tprm_assessment_section_question_action_plans 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_question_answer_options
		 FROM tprm_assessment_question_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answers
		 FROM tprm_assessment_answers 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_options.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_documents.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_comment_documents.tprm_assessment_question_id
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_status_updates
		 FROM tprm_assessment_section_question_action_plan_status_updates 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_status_updates.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plan_comments
		 FROM tprm_assessment_section_question_action_plan_comments 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tprm_assessment_section_question_action_plan_comments.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tasks
		 FROM tasks 
		 JOIN tprm_assessment_section_question_action_plans ON tprm_assessment_section_question_action_plans.id = tasks.tprm_assessment_section_question_action_plan_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_question_answer_options ON tprm_assessment_question_answer_options.id = tprm_assessment_answer_options.answer_option_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_options.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_answers ON tprm_assessment_answers.id = tprm_assessment_answer_comment_documents.tprm_assessment_answer_id
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_workflow_reviews
		 FROM tprm_assessment_workflow_reviews 
		 JOIN tprm_assessment_workflows ON tprm_assessment_workflows.id = tprm_assessment_workflow_reviews.tprm_assessment_workflow_id
		 WHERE tprm_assessment_workflows.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_workflow_participants
		 FROM tprm_assessment_workflow_participants 
		 JOIN tprm_assessment_workflows ON tprm_assessment_workflows.id = tprm_assessment_workflow_participants.tprm_assessment_workflow_id
		 WHERE tprm_assessment_workflows.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_questions
		 FROM tprm_assessment_section_questions 
		 JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id
		 WHERE tprm_assessment_sections.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_comments
		 FROM tprm_assessment_section_question_comments 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_comments.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_section_question_action_plans
		 FROM tprm_assessment_section_question_action_plans 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_action_plans.tprm_assessment_section_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_question_answer_options
		 FROM tprm_assessment_question_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_question_answer_options.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answers
		 FROM tprm_assessment_answers 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answers.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_options
		 FROM tprm_assessment_answer_options 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_options.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_documents
		 FROM tprm_assessment_answer_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_documents.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE tprm_assessment_answer_comment_documents
		 FROM tprm_assessment_answer_comment_documents 
		 JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_answer_comment_documents.tprm_assessment_question_id
		 WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE FROM tprm_assessment_workflows WHERE tprm_assessment_workflows.third_party_assessment_id = ?`,
		`DELETE FROM tprm_assessment_sections WHERE tprm_assessment_sections.third_party_assessment_id = ?`,
		`DELETE FROM tprm_assessment_section_questions WHERE tprm_assessment_section_questions.third_party_assessment_id = ?`,

		`DELETE FROM tprm_assessment_section_question_comments WHERE tprm_assessment_section_question_comments.third_party_assessment_id = ?`,
	}

	for _, q := range queries {
		if err := db.Exec(q, id).Error; err != nil {
			return err
		}
	}

	/* ================= ROOT ================= */
	return db.Delete(&models.ThirdPartyAssessment{}, id).Error
}
