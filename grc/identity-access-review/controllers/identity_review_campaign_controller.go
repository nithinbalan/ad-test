package controllers

import (
	"errors"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IdentityReviewCampaignController struct {
	service services.IdentityReviewCampaignServiceInterface
}

func NewIdentityReviewCampaignController(service services.IdentityReviewCampaignServiceInterface) *IdentityReviewCampaignController {
	return &IdentityReviewCampaignController{
		service: service,
	}
}

func (ircc *IdentityReviewCampaignController) GetAllIdentityReviewCampaigns(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	sortBy := c.DefaultQuery("sort_by", "campaign_name")
	keyword := c.Query("keyword")

	allowedParams := []string{
		"review_type_ids",
		"campaign_status_ids",
		"reviewer_ids",
		"approver_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil || page < 1 {
		utility.SendInvalidErrResp(c, errors.New("invalid 'page' parameter"))
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil || pageSize < 1 {
		utility.SendInvalidErrResp(c, errors.New("invalid 'page_size' parameter"))
		return
	}

	validSortFields := map[string]string{
		"reference_number":         "identity_review_campaigns.reference_number",
		"campaign_name":            "identity_review_campaigns.campaign_name",
		"start_date":               "identity_review_campaigns.start_date",
		"due_date":                 "identity_review_campaigns.due_date",
		"completion_percentage":    "identity_review_campaigns.completion_percentage",
		"created_at":               "identity_review_campaigns.created_at",
		"updated_at":               "identity_review_campaigns.updated_at",
		"review_type_title_en":     "review_types.title_en",
		"review_type_title_ar":     "review_types.title_ar",
		"campaign_status_title_en": "review_campaign_statuses.title_en",
		"campaign_status_title_ar": "review_campaign_statuses.title_ar",
		"reviewer_first_name_en":   "reviewer.first_name_en",
		"reviewer_first_name_ar":   "reviewer.first_name_ar",
		"approver_first_name_en":   "approver.first_name_en",
		"approver_first_name_ar":   "approver.first_name_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

	params := &dtos.IdentityReviewCampaignListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	}

	campaigns, err := ircc.service.GetAll(params)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (ircc *IdentityReviewCampaignController) GetIdentityReviewCampaignDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	campaign, err := ircc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, campaign)
}

func (ircc *IdentityReviewCampaignController) CreateIdentityReviewCampaign(c *gin.Context) {
	var formDTO dtos.IdentityReviewCampaignFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	id, err := ircc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ircc *IdentityReviewCampaignController) UpdateIdentityReviewCampaign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.IdentityReviewCampaignFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	if err := ircc.service.Update(uint(id), currentUserID, &formDTO); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (ircc *IdentityReviewCampaignController) DeleteIdentityReviewCampaigns(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	failedIDs, err := ircc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ircc *IdentityReviewCampaignController) GetIdentityReviewCampaignDocumentPreview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := ircc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
		}
		return
	}

	c.File(documentPath)
}

func (ircc *IdentityReviewCampaignController) DownloadIdentityReviewCampaignDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := ircc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
		}
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(documentPath))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.File(documentPath)
}
