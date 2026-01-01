package controllers

import (
	"errors"
	"fmt"
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

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(tpqlc.db, c, "Invalid export parameter")
			return
		}
	}

	fmt.Println(export)

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
		"reference_number":                       "third_party_question_libraries.reference_number",
		"question_title_en":                      "third_party_question_libraries.question_title_en",
		"question_title_ar":                      "third_party_question_libraries.question_title_ar",
		"third_party_question_category_title_en": "third_party_question_categories.title_en",
		"third_party_question_category_title_ar": "third_party_question_categories.title_ar",
		"third_party_question_type_title_en":     "third_party_question_types.title_en",
		"third_party_question_type_title_ar":     "third_party_question_types.title_ar",
		"created_at":                             "third_party_question_libraries.created_at",
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
		Export:      export,
	}

	result, err := tpqlc.service.GetAll(&query)
	if err != nil {
		utility.SendInternalErrResp(tpqlc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(tpqlc.db, c, err.Error())
			return
		}

		excelFile, err := tpqlc.service.ExportThirdPartyQuestionLibraries(result.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(tpqlc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Third Party Question Library.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

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
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
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

func (tpqlc *ThirdPartyQuestionLibraryController) GetThirdPartyQuestionLibraryExport(c *gin.Context) {
	c.Set("export", true)

	tpqlc.GetAllThirdPartyQuestionLibraries(c)
}

func (rc *ThirdPartyQuestionLibraryController) GenerateQuestionLibraryTemplate(c *gin.Context) {
	file, err := rc.service.GenerateQuestionLibraryTemplate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate template"})
		return
	}
	defer file.Close()

	c.Header("Content-Disposition", "attachment; filename=QuestionLibraryTemplate.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(file.Name())
}

func (rc *ThirdPartyQuestionLibraryController) ValidateQuestionLibraryTemplate(c *gin.Context) {
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

	validatedData, err := rc.service.ValidateQuestionLibraryTemplate(openedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (rc *ThirdPartyQuestionLibraryController) ValidateImportQuestionLibraries(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	validatedData, err := rc.service.ValidateImportQuestionLibraries(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, validatedData)
}

func (rc *ThirdPartyQuestionLibraryController) ImportQuestionLibraries(c *gin.Context) {
	var input []map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rc.db, c, err.Error())
		return
	}

	createdIDs, err := rc.service.ImportQuestionLibraries(input, uint(currentUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_ids": createdIDs})
}
