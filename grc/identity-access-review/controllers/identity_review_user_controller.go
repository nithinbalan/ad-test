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
		"username":              "identity_review_users.username",
		"department_title_en":   "departments.title_en",
		"department_title_ar":   "departments.title_ar",
		"privilege_title_en":    "privileged_accounts.title_en",
		"privilege_title_ar":    "privileged_accounts.title_ar",
		"account_type_title_en": "account_types.title_en",
		"account_type_title_ar": "account_types.title_ar",
		"created_at":            "identity_review_users.created_at",
		"updated_at":            "identity_review_users.updated_at",
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

	allowedFilters := []string{"department_ids", "privilege_ids", "account_type_ids"}
	filters := utility.ParseQueryFilters(allowedFilters, c.Request.URL.Query())

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
