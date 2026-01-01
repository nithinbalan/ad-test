package services

import (
	"errors"
	// "fmt"
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TprmAssessmentSectionQuestionCommentServiceInterface interface {
	GetAll(questionID uint, assessmentID uint, currentUserID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error)
	GetFullDetailsByID(questionID uint, assessmentID uint, id uint) (*dtos.TprmAssessmentSectionQuestionCommentDetailDTO, error)
	GetAllByAssessmentID(assessmentID uint, currentUserID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error)
	Create(questionID uint, assessmentID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO, currentUserID uint) (uint, error)
	CreateAssessmentComment(assessmentID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO, currentUserID uint) (uint, error)
	Update(assessmentID uint, id uint, currentUserID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO) error
	Delete(assessmentID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type TprmAssessmentSectionQuestionCommentService struct {
	repo              repositories.TprmAssessmentSectionQuestionCommentRepositoryInterface
	assessmentService ThirdPartyAssessmentServiceInterface
}

func NewTprmAssessmentSectionQuestionCommentService(
	repo repositories.TprmAssessmentSectionQuestionCommentRepositoryInterface,
	assessmentService ThirdPartyAssessmentServiceInterface,
) TprmAssessmentSectionQuestionCommentServiceInterface {
	return &TprmAssessmentSectionQuestionCommentService{
		repo:              repo,
		assessmentService: assessmentService,
	}
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) GetAll(questionID uint, assessmentID uint, currentUserID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error) {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return nil, err
	}

	comments, err := tasqcs.repo.GetAll(questionID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*dtos.TprmAssessmentSectionQuestionCommentListDTO)
	var roots []*dtos.TprmAssessmentSectionQuestionCommentListDTO

	for _, c := range comments {
		dto := &dtos.TprmAssessmentSectionQuestionCommentListDTO{
			ID:                                 c.ID,
			TprmAssessmentSectionQuestionID:    c.TprmAssessmentSectionQuestionID,
			TprmAssessmentSectionQuestionOrder: c.TprmAssessmentSectionQuestionOrder,
			TprmAssessmentSectionTitleEn:       c.TprmAssessmentSectionTitleEn,
			TprmAssessmentSectionTitleAr:       c.TprmAssessmentSectionTitleAr,
			ParentCommentID:                    c.ParentCommentID,
			Comment:                            c.Comment,
			Level:                              c.Level,
			CreatedByID:                        c.CreatedByID,
			CreatedAt:                          c.CreatedAt,

			CreatedUserFirstNameEn:  c.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:  c.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:   c.CreatedUserLastNameEn,
			CreatedUserLastNameAr:   c.CreatedUserLastNameAr,
			CreatedUserInitial:      c.CreatedUserInitial,
			CreatedUserColor:        c.CreatedUserColor,
			CreatedUserEmail:        c.CreatedUserEmail,
			CreatedUserProfileImage: c.CreatedUserProfileImage,

			ChildComments: []*dtos.TprmAssessmentSectionQuestionCommentListDTO{},
		}

		dtoMap[c.ID] = dto
	}

	// Build threaded format
	for _, c := range comments {
		dto := dtoMap[c.ID]

		if c.ParentCommentID != nil {
			parent := dtoMap[*c.ParentCommentID]
			parent.ChildComments = append(parent.ChildComments, dto)
		} else {
			roots = append(roots, dto)
		}
	}

	return roots, nil
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) GetAllByAssessmentID(assessmentID uint, currentUserID uint, params *dtos.TprmAssessmentSectionQuestionCommentListQueryParamsDTO) ([]*dtos.TprmAssessmentSectionQuestionCommentListDTO, error) {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return nil, err
	}

	comments, err := tasqcs.repo.GetAllByAssessmentID(assessmentID, params)
	if err != nil {
		return nil, err
	}

	dtoMap := make(map[uint]*dtos.TprmAssessmentSectionQuestionCommentListDTO)
	var roots []*dtos.TprmAssessmentSectionQuestionCommentListDTO

	// Build DTO map
	for _, c := range comments {
		dto := &dtos.TprmAssessmentSectionQuestionCommentListDTO{
			ID:                                 c.ID,
			TprmAssessmentSectionQuestionID:    c.TprmAssessmentSectionQuestionID,
			TprmAssessmentSectionQuestionOrder: c.TprmAssessmentSectionQuestionOrder,
			TprmAssessmentSectionTitleEn:       c.TprmAssessmentSectionTitleEn,
			TprmAssessmentSectionTitleAr:       c.TprmAssessmentSectionTitleAr,
			ParentCommentID:                    c.ParentCommentID,
			Comment:                            c.Comment,
			Level:                              c.Level,
			CreatedByID:                        c.CreatedByID,
			CreatedAt:                          c.CreatedAt,

			CreatedUserFirstNameEn:  c.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:  c.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:   c.CreatedUserLastNameEn,
			CreatedUserLastNameAr:   c.CreatedUserLastNameAr,
			CreatedUserInitial:      c.CreatedUserInitial,
			CreatedUserColor:        c.CreatedUserColor,
			CreatedUserEmail:        c.CreatedUserEmail,
			CreatedUserProfileImage: c.CreatedUserProfileImage,

			ChildComments: []*dtos.TprmAssessmentSectionQuestionCommentListDTO{},
		}
		dtoMap[c.ID] = dto
	}

	for _, c := range comments {
		dto := dtoMap[c.ID]

		if c.ParentCommentID != nil {
			parent := dtoMap[*c.ParentCommentID]
			parent.ChildComments = append(parent.ChildComments, dto)
		} else {
			roots = append(roots, dto)
		}
	}

	return roots, nil
}
func (tasqcs *TprmAssessmentSectionQuestionCommentService) Create(questionID uint, assessmentID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO, currentUserID uint) (uint, error) {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return 0, err
	}

	comment := &models.TPRMAssessmentSectionQuestionComment{
		ThirdPartyAssessmentID:                 assessmentID,
		TprmAssessmentSectionQuestionID:        &questionID,
		TprmAssessmentSectionQuestionCommentID: dto.TprmAssessmentSectionQuestionCommentID,
		Comment:                                dto.Comment,
		Level:                                  *assessment.NextReviewLevel,
		CreatedBYID:                            currentUserID,
	}

	if err := tasqcs.repo.Create(comment); err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func SafeString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) GetFullDetailsByID(questionID uint, assessmentID uint, id uint) (*dtos.TprmAssessmentSectionQuestionCommentDetailDTO, error) {
	var commentDetail dtos.TprmAssessmentSectionQuestionCommentDetailDTO
	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return nil, err
	}
	comment, err := tasqcs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Convert to DTO
	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil

	// return comment, nil
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) Update(assessmentID uint, id uint, currentUserID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO) error {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return err
	}

	existingComment, err := tasqcs.repo.GetByID(id)
	if err != nil {
		return err
	}

	oldComment, err := tasqcs.repo.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existingComment.Comment = dto.Comment

	err = tasqcs.repo.Update(existingComment)
	if err != nil {
		return err
	}

	go func() {
		newComment, err := tasqcs.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changes := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": newComment.Comment},
		)

		if len(changes) > 0 {
			// Add logging later if needed
		}
	}()

	return nil
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) Delete(assessmentID uint, ids []uint, currentUserID uint) ([]uint, error) {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {

		if _, err := tasqcs.repo.GetByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		err := tasqcs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		return failedIDs, nil
	}

	return failedIDs, nil
}

func (tasqcs *TprmAssessmentSectionQuestionCommentService) CreateAssessmentComment(assessmentID uint, dto *dtos.TprmAssessmentSectionQuestionCommentFormDTO, currentUserID uint) (uint, error) {

	assessment, err := tasqcs.assessmentService.GetByID(assessmentID)
	if assessment == nil {
		return 0, err
	}

	comment := &models.TPRMAssessmentSectionQuestionComment{
		ThirdPartyAssessmentID:                 assessmentID,
		TprmAssessmentSectionQuestionID:        nil,
		TprmAssessmentSectionQuestionCommentID: dto.TprmAssessmentSectionQuestionCommentID,
		Comment:                                dto.Comment,
		Level:                                  *assessment.NextReviewLevel,
		CreatedBYID:                            currentUserID,
	}

	if err := tasqcs.repo.Create(comment); err != nil {
		return 0, err
	}

	return comment.ID, nil
}
