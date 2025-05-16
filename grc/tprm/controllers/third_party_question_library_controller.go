package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"jethur/grc/tprm/dtos"
	"jethur/grc/tprm/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThirdPartyQuestionLibraryController struct {
	service services.ThirdPartyQuestionLibraryServiceInterface
	db      *gorm.DB
}

func NewThirdPartyQuestionLibraryController(service services.ThirdPartyQuestionLibraryServiceInterface) *ThirdPartyQuestionLibraryController {
	return &ThirdPartyQuestionLibraryController{service: service}
}

func (tpqlc *ThirdPartyQuestionLibraryController) GetAllThirdPartyQuestionLibraries(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "question_title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"third_party_question_category_ids",
		"third_party_question_type_ids",
	}
	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("Invalid page_size parameter"))
		return
	}

	validSortFields := map[string]string{
		"question_title_en":                      "third_party_question_libraries.question_title_en",
		"question_title_ar":                      "third_party_question_libraries.question_title_ar",
		"third_party_question_category_title_en": "third_party_question_categories.title_en",
		"third_party_question_category_title_ar": "third_party_question_categories.title_ar",
		"third_party_question_type_title_en":     "third_party_question_types.title_en",
		"third_party_question_type_title_ar":     "third_party_question_types.title_ar",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, errors.New("Invalid sort_by parameter"))
		return
	}

	query := dtos.ThirdPartyQuestionLibraryListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	result, err := tpqlc.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (tpqlc *ThirdPartyQuestionLibraryController) GetThirdPartyQuestionLibraryDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	result, err := tpqlc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tpqlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (tpqlc *ThirdPartyQuestionLibraryController) CreateThirdPartyQuestionLibrary(c *gin.Context) {
	var form dtos.ThirdPartyQuestionLibraryFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}

	id, err := tpqlc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpqlc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (tpqlc *ThirdPartyQuestionLibraryController) UpdateThirdPartyQuestionLibrary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form dtos.ThirdPartyQuestionLibraryFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}

	err = tpqlc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tpqlc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tpqlc *ThirdPartyQuestionLibraryController) DeleteThirdPartyQuestionLibraries(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}

	failedIDs, err := tpqlc.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
