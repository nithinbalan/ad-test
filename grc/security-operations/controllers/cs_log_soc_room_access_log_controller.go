package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CSLogSocRoomAccessLogController struct {
	service services.CSLogSocRoomAccessLogServiceInterface
	db      *gorm.DB
}

func NewCSLogSocRoomAccessLogController(service services.CSLogSocRoomAccessLogServiceInterface) *CSLogSocRoomAccessLogController {
	return &CSLogSocRoomAccessLogController{
		service: service,
	}
}

func (csralc *CSLogSocRoomAccessLogController) GetAllCSLogSocRoomAccessLogs(c *gin.Context) {
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
			utility.SendInternalErrResp(csralc.db, c, "Invalid export parameter")
			return
		}
	}

	allowedParams := []string{
		"department_ids",
		"authorizing_officer_ids",
	}
	validSortFields := map[string]string{
		"id":                                "soc.id",
		"reference_number":                  "soc.reference_number",
		"visitor_full_name":                 "soc.visitor_full_name",
		"entry_time":                        "soc.entry_time",
		"exit_time":                         "soc.exit_time",
		"visitor_id_number":                 "soc.visitor_id_number",
		"external_party_name":               "soc.external_party_name",
		"department_title_en":               "department.title_en",
		"department_title_ar":               "department.title_ar",
		"authorizing_officer_first_name_en": "authorizing_officer.first_name_en",
		"authorizing_officer_first_name_ar": "authorizing_officer.first_name_ar",
		"authorizing_officer_last_name_en":  "authorizing_officer.last_name_en",
		"authorizing_officer_last_name_ar":  "authorizing_officer.last_name_ar",
		"created_at":                        "soc.created_at",
		"updated_at":                        "soc.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
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
	queryParams := socDtos.CSLogSocRoomAccessLogListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortField,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
		Filters:  filters,
		Export:   export,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	entries, err := csralc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(csralc.db, c, err.Error())
		return
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(csralc.db, c, err.Error())
			return
		}

		excelFile, err := csralc.service.ExportCSLogSocRoomAccessLog(entries.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(csralc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="SOC Room Access Log.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, entries)
}

func (csralc *CSLogSocRoomAccessLogController) GetCSLogSocRoomAccessLogDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	entry, err := csralc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(csralc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, entry)
}

func (csralc *CSLogSocRoomAccessLogController) CreateCSLogSocRoomAccessLog(c *gin.Context) {
	var formDTO socDtos.CSLogSocRoomAccessLogFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csralc.db, c, err.Error())
		return
	}

	id, err := csralc.service.Create(currentUserID, &formDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csralc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (csralc *CSLogSocRoomAccessLogController) UpdateCSLogSocRoomAccessLog(c *gin.Context) {
	var formDTO socDtos.CSLogSocRoomAccessLogFormDTO

	entryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csralc.db, c, err.Error())
		return
	}

	err = csralc.service.Update(uint(entryID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(csralc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (csralc *CSLogSocRoomAccessLogController) DeleteCSLogSocRoomAccessLogs(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(csralc.db, c, err.Error())
		return
	}

	failedIDs, err := csralc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(csralc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (csralc *CSLogSocRoomAccessLogController) GetCSLogSocRoomAccessLogDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := csralc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(csralc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (csralc *CSLogSocRoomAccessLogController) DownloadCSLogSocRoomAccessLogDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := csralc.service.GetDocumentPath(uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(csralc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}

func (csralc *CSLogSocRoomAccessLogController) GetCSLogSocRoomAccessLogExport(c *gin.Context) {
	c.Set("export", true)

	csralc.GetAllCSLogSocRoomAccessLogs(c)
}
