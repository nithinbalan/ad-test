package controllers

import (
	"errors"
	"jethur/grc/policy/dtos"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/grc/policy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyFolderController struct {
	service services.PolicyFolderServiceInterface
	db      *gorm.DB
}

func NewPolicyFolderController(service services.PolicyFolderServiceInterface) *PolicyFolderController {
	return &PolicyFolderController{
		service: service,
	}
}

func (pfc *PolicyFolderController) GetAllPolicyFolders(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, "Invalid export parameter")
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
		"parent_document_ids",
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
		utility.SendInternalErrResp(pfc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 20)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, "Invalid page_size parameter")
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

	documents, err := pfc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (pfc *PolicyFolderController) GetPolicyFolderDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	folder, err := pfc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, folder)
}

func (pfc *PolicyFolderController) CreatePolicyFolder(c *gin.Context) {
	var folder *dtos.PolicyFolderFormDTO

	if err := c.ShouldBindJSON(&folder); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, err.Error())
		return
	}

	createdFolder, err := pfc.service.Create(currentUserID, folder)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdFolder)
}

func (pfc *PolicyFolderController) UpdatePolicyFolder(c *gin.Context) {
	var updatedFolder *dtos.PolicyFolderFormDTO

	folderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updatedFolder); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, err.Error())
		return
	}

	err = pfc.service.Update(uint(folderID), currentUserID, updatedFolder)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(pfc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedFolder)
}

func (pfc *PolicyFolderController) DeletePolicyFolders(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, err.Error())
		return
	}

	failedIDs, err := pfc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(pfc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
