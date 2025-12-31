package repositories

import (
	// "encoding/json"
	"fmt"
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"

	"gorm.io/gorm"
)

type TprmAssessmentSectionQuestionCommentRepositoryInterface interface {
	GetAll(questionID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error)
	GetByID(id uint) (*models.TPRMAssessmentSectionQuestionComment, error)
	GetFullDetailsByID(id uint) (*models.TPRMAssessmentSectionQuestionComment, error)
	GetAllByAssessmentID(assessmentID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error)
	Create(comment *models.TPRMAssessmentSectionQuestionComment) error
	Update(comment *models.TPRMAssessmentSectionQuestionComment) error
	Delete(id uint) error
}

type TprmAssessmentSectionQuestionCommentRepository struct {
	db *gorm.DB
}

func NewTprmAssessmentSectionQuestionCommentRepository(db *gorm.DB) TprmAssessmentSectionQuestionCommentRepositoryInterface {
	return &TprmAssessmentSectionQuestionCommentRepository{db: db}
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) GetAll(questionID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error) {

	var comments []*dtos.TprmAssessmentSectionQuestionCommentListDTO

	query := repo.db.
		Table("tprm_assessment_section_question_comments").
		Joins("LEFT JOIN users ON users.id = tprm_assessment_section_question_comments.created_by_id").
		Joins("LEFT JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = tprm_assessment_section_question_comments.tprm_assessment_section_question_id").
		Joins("LEFT JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id").
		Where("tprm_assessment_section_question_comments.tprm_assessment_section_question_id = ?", questionID).
		Select(`
        tprm_assessment_section_question_comments.id,
        tprm_assessment_section_question_comments.tprm_assessment_section_question_id,
        tprm_assessment_section_question_comments.tprm_assessment_section_question_comment_id AS parent_comment_id,
        tprm_assessment_section_question_comments.comment,
        tprm_assessment_section_question_comments.level,
        tprm_assessment_section_question_comments.created_by_id,
        tprm_assessment_section_question_comments.created_at,
		tprm_assessment_section_questions.[order] AS tprm_assessment_section_question_order,
		tprm_assessment_sections.title_en AS tprm_assessment_section_title_en,
		tprm_assessment_sections.title_ar AS tprm_assessment_section_title_ar,

        users.first_name_en AS created_user_first_name_en,
        users.first_name_ar AS created_user_first_name_ar,
        users.last_name_en AS created_user_last_name_en,
        users.last_name_ar AS created_user_last_name_ar,
        users.initial AS created_user_initial,
        users.color AS created_user_color,
        users.email AS created_user_email,
        users.profile_image_path AS created_user_profile_image_path
       `).Order("tprm_assessment_section_question_comments.created_at DESC")

	if params.Keyword != "" {
		query = query.Where("tprm_assessment_section_question_comments.comment LIKE ?", "%"+params.Keyword+"%")
	}

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("tprm_assessment_section_question_comments."+columnName+" IN (?)", values)
		}
	}

	// Sort validation

	if err := query.Debug().Scan(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) GetAllByAssessmentID(
	assessmentID uint,
	params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO,
) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error) {

	var results []*dtos.TprmAssessmentSectionQuestionCommentListDTO

	query := repo.db.
		Table("tprm_assessment_section_question_comments c").
		Joins("LEFT JOIN tprm_assessment_section_questions q ON q.id = c.tprm_assessment_section_question_id").
		Joins("LEFT JOIN tprm_assessment_section_questions ON tprm_assessment_section_questions.id = c.tprm_assessment_section_question_id").
		Joins("LEFT JOIN tprm_assessment_sections ON tprm_assessment_sections.id = tprm_assessment_section_questions.tprm_assessment_section_id").
		Joins("LEFT JOIN users u ON u.id = c.created_by_id").
		Where("c.third_party_assessment_id = ? OR q.third_party_assessment_id = ?", assessmentID, assessmentID).
		Select(`
            c.id,
            c.tprm_assessment_section_question_id,
            c.tprm_assessment_section_question_comment_id AS parent_comment_id,
            c.comment,
            c.level,
            c.created_by_id,
            c.created_at,
			tprm_assessment_section_questions.[order] AS tprm_assessment_section_question_order,
			tprm_assessment_sections.title_en AS tprm_assessment_section_title_en,
			tprm_assessment_sections.title_ar AS tprm_assessment_section_title_ar,


            u.first_name_en AS created_user_first_name_en,
            u.first_name_ar AS created_user_first_name_ar,
            u.last_name_en AS created_user_last_name_en,
            u.last_name_ar AS created_user_last_name_ar,
            u.initial AS created_user_initial,
            u.color AS created_user_color,
            u.email AS created_user_email,
            u.profile_image_path AS created_user_profile_image_path
        `)

	for col, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		if col == "section_ids" || col == "section_id" {
			sectionIDs := values
			query = query.Where("q.tprm_assessment_section_id IN (?)", sectionIDs)
		} else {
			query = query.Where(col+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where("c.comment LIKE ?", "%"+params.Keyword+"%")
	}

	query = query.Order("c.created_at DESC")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) GetByID(id uint) (*models.TPRMAssessmentSectionQuestionComment, error) {
	var comment models.TPRMAssessmentSectionQuestionComment
	if err := repo.db.First(&comment, id).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) GetFullDetailsByID(id uint) (*models.TPRMAssessmentSectionQuestionComment, error) {

	var comment models.TPRMAssessmentSectionQuestionComment

	err := repo.db.
		Preload("CreatedByUser").
		First(&comment, id).
		Error

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) Create(comment *models.TPRMAssessmentSectionQuestionComment) error {

	err := repo.db.Create(comment).Error
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) Update(comment *models.TPRMAssessmentSectionQuestionComment) error {

	err := repo.db.Model(&models.TPRMAssessmentSectionQuestionComment{}).
		Where("id = ?", comment.ID).
		Updates(comment).Error

	if err != nil {
		return fmt.Errorf("failed to update comment with id %d: %w", comment.ID, err)
	}

	return nil
}

func (repo *TprmAssessmentSectionQuestionCommentRepository) Delete(id uint) error {
	return repo.db.Delete(&models.TPRMAssessmentSectionQuestionComment{}, id).Error
}
