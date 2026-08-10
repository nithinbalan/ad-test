package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	identityAccessReviewDTOs "jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityReviewCampaignReviewCommentController struct {
	service services.IdentityReviewCampaignReviewCommentServiceInterface
	db      *gorm.DB
}

func NewIdentityReviewCampaignReviewCommentController(service services.IdentityReviewCampaignReviewCommentServiceInterface) *IdentityReviewCampaignReviewCommentController {
	return &IdentityReviewCampaignReviewCommentController{
		service: service,
	}
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) GetAllComments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "comment")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"comment_status_ids",
		"identity_review_campaign_ids",
		"level",
		"created_by_ids",
	}

	validSortFields := map[string]string{
		"comment":                 "identity_review_campaign_review_comments.comment",
		"level":                   "identity_review_campaign_review_comments.level",
		"comment_status_title_en": "comment_statuses.title_en",
		"comment_status_title_ar": "comment_statuses.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrcc.db, c, err.Error())
		return
	}

	comments, err := ircrcc.service.GetAll(uint(campaignID), currentUserID, &queryParams)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comments)
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) GetCommentDetails(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	comment, err := ircrcc.service.GetFullDetailsByID(uint(campaignID), uint(commentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, comment)
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) CreateComment(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrcc.db, c, err.Error())
		return
	}

	id, err := ircrcc.service.Create(uint(campaignID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) UpdateComment(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrcc.db, c, err.Error())
		return
	}

	err = ircrcc.service.Update(uint(campaignID), uint(commentID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) DeleteComment(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrcc.db, c, err.Error())
		return
	}

	failedIDs, err := ircrcc.service.Delete(uint(campaignID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ircrcc *IdentityReviewCampaignReviewCommentController) ResolveComment(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ircrcc.db, c, err.Error())
		return
	}

	err = ircrcc.service.Resolve(uint(campaignID), uint(commentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(ircrcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
