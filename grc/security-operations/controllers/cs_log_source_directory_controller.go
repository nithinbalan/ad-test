package controllers

import (
	"errors"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSLogSourceDirectoryController struct {
	service services.CSLogSourceDirectoryServiceInterface
	db      *gorm.DB
}

func NewCSLogSourceDirectoryController(service services.CSLogSourceDirectoryServiceInterface) *CSLogSourceDirectoryController {
	return &CSLogSourceDirectoryController{
		service: service,
	}
}

func (csdc *CSLogSourceDirectoryController) GetAllCSLogSourceDirectories(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "desc")
	sortBy := c.DefaultQuery("sort_by", "id")
	keyword := c.Query("keyword")

	exportFromContext, exists := c.Get("export")
	var export bool
	var err error

	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err = strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(csdc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"data_entity_type_ids",
		"responsible_entity_ids",
		"beneficiary_entity_ids",
	}
	validSortFields := map[string]string{
		"id":                          "sd.id",
		"reference_number":            "sd.reference_number",
		"entity_name":                 "sd.entity_name",
		"created_at":                  "sd.created_at",
		"updated_at":                  "sd.updated_at",
		"data_entity_type_title_en":   "data_entity_types.title_en",
		"data_entity_type_title_ar":   "data_entity_types.title_ar",
		"responsible_entity_title_en": "responsible_entity.title_en",
		"responsible_entity_title_ar": "responsible_entity.title_ar",
		"beneficiary_entity_title_en": "beneficiary_entity.title_en",
		"beneficiary_entity_title_ar": "beneficiary_entity.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page parameter"))
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid page_size parameter"))
		return
	}
	if page < 1 || pageSize < 1 {
		utility.SendInvalidErrResp(c, errors.New("page and page_size must be greater than zero"))
		return
	}
	if sort != "asc" && sort != "desc" {
		utility.SendInvalidErrResp(c, errors.New("invalid 'sort_order' parameter"))
		return
	}
	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}
	queryParams := socDtos.CSLogSourceDirectoryListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
	}
	entries, err := csdc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(csdc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(csdc.db, c, err.Error())
			return
		}

		excelFile, err := csdc.service.ExportCSLogSourceDirectory(entries.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(csdc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="Source Directory.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, entries)
}

func (csdc *CSLogSourceDirectoryController) GetCSLogSourceDirectoryExport(c *gin.Context) {
	c.Set("export", true)

	csdc.GetAllCSLogSourceDirectories(c)
}

func (csdc *CSLogSourceDirectoryController) GetCSLogSourceDirectoryDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	entry, err := csdc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csdc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, entry)
}

func (csdc *CSLogSourceDirectoryController) CreateCSLogSourceDirectory(c *gin.Context) {
	var sourceDirectoryFormDTO socDtos.CSLogSourceDirectoryFormDTO

	if err := c.ShouldBindJSON(&sourceDirectoryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, c, err.Error())
		return
	}

	id, err := csdc.service.Create(currentUserID, &sourceDirectoryFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csdc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (csdc *CSLogSourceDirectoryController) UpdateCSLogSourceDirectory(c *gin.Context) {
	var sourceDirectoryFormDTO socDtos.CSLogSourceDirectoryFormDTO

	sourceDirectoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&sourceDirectoryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, c, err.Error())
		return
	}

	err = csdc.service.Update(uint(sourceDirectoryID), currentUserID, &sourceDirectoryFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csdc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (csdc *CSLogSourceDirectoryController) DeleteCSLogSourceDirectories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, c, err.Error())
		return
	}

	failedIDs, err := csdc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(csdc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
