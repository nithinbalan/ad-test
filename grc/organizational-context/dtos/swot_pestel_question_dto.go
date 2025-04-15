package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type SwotPestelQuestionListDTO struct {
	ID                             int        `json:"id"`
	TitleEn                        string     `json:"title_en"`
	TitleAr                        string     `json:"title_ar"`
	SwotPestelQuestionGroupID      uint       `json:"swot_pestel_question_group_id"`
	SwotPestelQuestionGroupTitleEn string     `json:"swot_pestel_question_group_title_en"`
	SwotPestelQuestionGroupTitleAr string     `json:"swot_pestel_question_group_title_ar"`
	Order                          int        `json:"order"`
	Status                         string     `json:"status"`
	IsSwot                         bool       `json:"is_swot"`
	IsPestel                       bool       `json:"is_pestel"`
	Answer                         *string    `json:"answer"`
	AnsweredByID                   *uint      `json:"answered_by_id"`
	AnsweredAt                     *time.Time `json:"answered_at"`
	AnsweredUserFirstNameEn        string     `json:"answered_user_first_name_en"`
	AnsweredUserFirstNameAr        string     `json:"answered_user_first_name_ar"`
	AnsweredUserLastNameEn         string     `json:"answered_user_last_name_en"`
	AnsweredUserLastNameAr         string     `json:"answered_user_last_name_ar"`
	AnsweredUserInitial            string     `json:"answered_user_initial"`
	AnsweredUserColor              string     `json:"answered_user_color"`
	AnsweredUserEmail              string     `json:"answered_user_email"`
	AnsweredUserProfileImagePath   string     `json:"answered_user_profile_image_path"`
}

type SwotPestelQuestionDTO struct {
	ID                        int        `json:"id"`
	SwotPestelQuestionGroupID uint       `json:"swot_pestel_question_group_id"`
	TitleEn                   string     `json:"title_en"`
	TitleAr                   string     `json:"title_ar"`
	Order                     int        `json:"order"`
	Status                    string     `json:"status"`
	IsSwot                    bool       `json:"is_swot"`
	IsPestel                  bool       `json:"is_pestel"`
	Answer                    *string    `json:"answer"`
	AnsweredByID              *uint      `json:"answered_by_id"`
	AnsweredAt                *time.Time `json:"answered_at"`

	AnsweredByUser             *organizationDTOs.UserMinimalDTO     `json:"answered_by"`
	SwotPestelQuestionGroup    *SwotPestelQuestionGroupDTO          `json:"swot_pestel_question_group"`
	SwotPestelQuestionInvitees []*SwotPestelQuestionInviteeDTO      `json:"swot_pestel_question_invitees" gorm:"foreignKey:SwotPestelQuestionID;references:ID"`
	SwotPestelQuestionFactors  []*SwotPestelQuestionFactorDetailDTO `json:"swot_pestel_question_factors" gorm:"foreignKey:SwotPestelQuestionID;references:ID"`
}

type AnswerFormDTO struct {
	Answer *string `json:"answer" binding:"required"`
}

type InviteeFormDTO struct {
	UserIds []*uint `json:"user_ids" binding:"required"`
}

type SwotPestelQuestionInviteeDTO struct {
	ID                   uint                            `json:"id"`
	SwotPestelQuestionID uint                            `json:"swot_pestel_question_id"`
	UserID               uint                            `json:"user_id"`
	User                 organizationDTOs.UserMinimalDTO `json:"user"`
}
