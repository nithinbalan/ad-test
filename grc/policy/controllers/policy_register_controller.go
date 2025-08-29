package controllers

import (
	"errors"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyRegisterController struct {
	service services.PolicyRegisterServiceInterface
	db      *gorm.DB
}

func NewPolicyRegisterController(service services.PolicyRegisterServiceInterface) *PolicyRegisterController {
	return &PolicyRegisterController{
		service: service,
	}
}

func (prc *PolicyRegisterController) GetAllDocuments(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid export parameter")
		return
	}

	// Define an array of allowed parameter names
	allowedParams := []string{
		"document_category_ids",
		"document_type_ids",
		"document_access_type_ids",
		"responsible_user_ids",
		"responsible_unit_ids",
		"document_status_ids",
		"policy_folder_ids",
	}

	validSortFields := map[string]string{
		"type_ar":             "document_types.title_ar",
		"type_en":             "document_types.title_en",
		"title_ar":            "documents.title_ar",
		"title_en":            "documents.title_en",
		"responsible_user_en": "responsible_user.first_name_en",
		"responsible_user_ar": "responsible_user.first_name_ar",
		"category_en":         "document_categories.title_en",
		"category_ar":         "document_categories.title_ar",
		"status_en":           "document_statuses.title_ar",
		"status_ar":           "document_statuses.title_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := policyDTOs.PolicyRegisterListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	documents, err := prc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}

		excelFile, err := prc.service.ExportDocuments(documents.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(prc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Document.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}
	c.JSON(http.StatusOK, documents)
}

func (prc *PolicyRegisterController) GetDocumentDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	document, err := prc.service.GetFullDetailsByID(uint(id), &currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, document)
}

func (prc *PolicyRegisterController) CreateDocument(c *gin.Context) {
	var documentFormDTO policyDTOs.PolicyRegisterFormDTO
	if err := c.ShouldBindJSON(&documentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	createdDocument, err := prc.service.Create(currentUserID, &documentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdDocument)
}

func (prc *PolicyRegisterController) UpdateDocument(c *gin.Context) {
	documentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var documentFormDTO policyDTOs.PolicyRegisterFormDTO
	if err := c.ShouldBindJSON(&documentFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	err = prc.service.Update(uint(documentID), currentUserID, &documentFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		default:
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (prc *PolicyRegisterController) DeleteDocuments(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(prc.db, c, err.Error())
		return
	}

	failedIDs, err := prc.service.Delete(request.IDs, currentUserID)
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
			utility.SendInternalErrResp(prc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
