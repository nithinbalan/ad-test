package controllers

import (
	"errors"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IdentityReviewUserController struct {
	service services.IdentityReviewUserServiceInterface
}

func NewIdentityReviewUserController(service services.IdentityReviewUserServiceInterface) *IdentityReviewUserController {
	return &IdentityReviewUserController{service: service}
}

func (iruc *IdentityReviewUserController) GetAllIdentityReviewUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}

	sortBy := c.DefaultQuery("sort_by", "username")
	validSortFields := map[string]string{
		"identity_review_campaign_name": "identity_review_campaigns.campaign_name",
		"username":                      "identity_review_users.username",
		"user_identifier":               "identity_review_users.user_identifier",
		"last_login":                    "identity_review_users.last_login",
		"current_role":                  "identity_review_users.current_role",
		"department":                    "departments.title_en",
		"department_title_en":           "departments.title_en",
		"department_title_ar":           "departments.title_ar",
		"privilege":                     "privileged_accounts.title_en",
		"privilege_title_en":            "privileged_accounts.title_en",
		"privilege_title_ar":            "privileged_accounts.title_ar",
		"account_type":                  "account_types.title_en",
		"account_type_title_en":         "account_types.title_en",
		"account_type_title_ar":         "account_types.title_ar",
		"created_at":                    "identity_review_users.created_at",
		"updated_at":                    "identity_review_users.updated_at",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_by' parameter"))
		return
	}

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

	allowedParams := []string{
		"identity_review_campaign_ids",
		"department_ids",
		"privilege_ids",
		"account_type_ids",
		"account_decision_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	response, err := iruc.service.GetAll(&dtos.IdentityReviewUserListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  c.Query("keyword"),
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
	})
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (iruc *IdentityReviewUserController) GetIdentityReviewUserDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	identityReviewUser, err := iruc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(nil, c, err.Error())
		}
		return
	}

	if campaignIDParam := c.Query("identity_review_campaign_id"); campaignIDParam != "" {
		campaignID, err := strconv.ParseUint(campaignIDParam, 10, 32)
		if err != nil || campaignID == 0 {
			utility.SendInvalidErrResp(c, errors.New("invalid 'identity_review_campaign_id' parameter"))
			return
		}
		if identityReviewUser.IdentityReviewCampaignID != uint(campaignID) {
			utility.SendNotFoundErrResp(c)
			return
		}
	}

	c.JSON(http.StatusOK, identityReviewUser)
}

func (iruc *IdentityReviewUserController) CreateIdentityReviewUser(c *gin.Context) {
	var formDTO dtos.IdentityReviewUserFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	identityReviewUser, err := iruc.service.Create(currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, identityReviewUser)
}

func (iruc *IdentityReviewUserController) UpdateIdentityReviewUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.IdentityReviewUserFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	identityReviewUser, err := iruc.service.Update(uint(id), currentUserID, &formDTO)
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

	c.JSON(http.StatusOK, identityReviewUser)
}

func (iruc *IdentityReviewUserController) UpdateIdentityReviewUserAccountDecision(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO dtos.IdentityReviewUserAccountDecisionUpdateFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	identityReviewUser, err := iruc.service.UpdateAccountDecision(uint(id), currentUserID, &formDTO)
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

	c.JSON(http.StatusOK, identityReviewUser)
}

func (iruc *IdentityReviewUserController) GetIdentityReviewUserAccountDecisionHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	history, err := iruc.service.GetAccountDecisionHistory(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
		}
		return
	}

	c.JSON(http.StatusOK, history)
}

func (iruc *IdentityReviewUserController) DeleteIdentityReviewUsers(c *gin.Context) {
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

	failedIDs, err := iruc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (iruc *IdentityReviewUserController) GenerateIdentityReviewUserTemplate(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	file, err := iruc.service.GenerateIdentityReviewUserTemplate(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template"})
		return
	}
	defer file.Close()

	c.Header("Content-Disposition", "attachment; filename=IdentityReviewUserTemplate.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(file.Name())
}

func (iruc *IdentityReviewUserController) ValidateIdentityReviewUserTemplate(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer openedFile.Close()

	campaignID, err := strconv.ParseUint(c.Query("identity_review_campaign_id"), 10, 32)
	if err != nil || campaignID == 0 {
		utility.SendInvalidErrResp(c, errors.New("invalid 'identity_review_campaign_id' parameter"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	validatedData, err := iruc.service.ValidateIdentityReviewUserTemplate(openedFile, uint(campaignID), currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (iruc *IdentityReviewUserController) ValidateImportIdentityReviewUsers(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	campaignID, err := strconv.ParseUint(c.Query("identity_review_campaign_id"), 10, 32)
	if err != nil || campaignID == 0 {
		utility.SendInvalidErrResp(c, errors.New("invalid 'identity_review_campaign_id' parameter"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	validatedData, err := iruc.service.ValidateImportIdentityReviewUsers(input, uint(campaignID), currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (iruc *IdentityReviewUserController) ImportIdentityReviewUsers(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	campaignID, err := strconv.ParseUint(c.Query("identity_review_campaign_id"), 10, 32)
	if err != nil || campaignID == 0 {
		utility.SendInvalidErrResp(c, errors.New("invalid 'identity_review_campaign_id' parameter"))
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	createdIDs, err := iruc.service.ImportIdentityReviewUsers(input, uint(campaignID), currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_ids": createdIDs})
}
