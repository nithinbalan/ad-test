package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type TprmAssessmentSectionQuestionCommentListDTO struct {
	ID                                 uint      `json:"id"`
	TprmAssessmentSectionQuestionID    uint      `json:"tprm_assessment_section_question_id"`
	TprmAssessmentSectionQuestionOrder int       `json:"tprm_assessment_section_question_order"`
	TprmAssessmentSectionTitleEn       string    `json:"tprm_assessment_section_title_en"`
	TprmAssessmentSectionTitleAr       string    `json:"tprm_assessment_section_title_ar"`
	ParentCommentID                    *uint     `json:"parent_comment_id"`
	Comment                            string    `json:"comment"`
	Level                              int       `json:"level"`
	CreatedByID                        uint      `json:"created_by_id"`
	CreatedAt                          time.Time `json:"created_at"`

	// User Information
	CreatedUserFirstNameEn  string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string  `json:"created_user_last_name_ar"`
	CreatedUserInitial      string  `json:"created_user_initial"`
	CreatedUserColor        string  `json:"created_user_color"`
	CreatedUserEmail        string  `json:"created_user_email"`
	CreatedUserProfileImage *string `json:"created_user_profile_image_path"`

	// Replies
	ChildComments []*TprmAssessmentSectionQuestionCommentListDTO `json:"child_comments" gorm:"-"`
}

type TprmAssessmentSectionQuestionCommentFormDTO struct {
	TprmAssessmentSectionQuestionCommentID *uint  `json:"tprm_assessment_section_question_comment_id"` // parent comment id (nullable)
	Comment                                string `json:"comment" binding:"required"`
}

type TprmAssessmentSectionQuestionCommentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}

type TprmAssessmentSectionQuestionCommentDetailDTO struct {
	ID              uint                            `json:"id"`
	Comment         string                          `json:"comment"`
	Level           int                             `json:"level"`
	ParentCommentID *uint                           `json:"parent_comment_id"`
	CreatedAt       time.Time                       `json:"created_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO `json:"created_by"`
}
