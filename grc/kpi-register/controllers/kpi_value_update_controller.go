package controllers

import (
	"errors"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	"jethur/grc/kpi-register/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KPIValueUpdateController struct {
	service services.KPIValueUpdateServiceInterface
	db      *gorm.DB
}

func NewKPIValueUpdateController(service services.KPIValueUpdateServiceInterface) *KPIValueUpdateController {
	return &KPIValueUpdateController{
		service: service,
	}
}

func (kvuc *KPIValueUpdateController) GetAllKPIValueUpdates(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reviewed_date")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	var export bool
	if val, exists := c.Get("export"); exists {
		export = val.(bool)
	} else {
		export, _ = strconv.ParseBool(c.DefaultQuery("export", "false"))
	}

	allowedParams := []string{
		"kpi_register_ids",
		"kpi_category_ids",
		"review_user_ids",
		"kpi_review_status_ids",
		"created_by_ids",
		"updated_by_ids",
		"responsible_user_ids",
	}

	validSortFields := map[string]string{
		"id":                         "kpi_value_updates.id",
		"reviewed_date":              "kpi_value_updates.reviewed_date",
		"value":                      "kpi_value_updates.value",
		"start_date":                 "kpi_value_updates.start_date",
		"end_date":                   "kpi_value_updates.end_date",
		"period_label":               "kpi_value_updates.period_label",
		"remarks":                    "kpi_value_updates.remarks",
		"kpi_register_title":         "kpi_register.title",
		"kpi_category_title_en":      "kpi_categories.title_en",
		"kpi_category_title_ar":      "kpi_categories.title_ar",
		"review_user_first_name_en":  "review_user.first_name_en",
		"review_user_first_name_ar":  "review_user.first_name_ar",
		"review_user_last_name_en":   "review_user.last_name_en",
		"review_user_last_name_ar":   "review_user.last_name_ar",
		"kpi_review_status_title_en": "kpi_review_statuses.title_en",
		"kpi_review_status_title_ar": "kpi_review_statuses.title_ar",
		"created_at":                 "kpi_value_updates.created_at",
		"updated_at":                 "kpi_value_updates.updated_at",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page_size parameter")
		return
	}

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := kpiDTOs.KPIValueUpdateListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		Export:      export,
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	kpiValueUpdates, err := kvuc.service.GetAll(&queryParams, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
			return
		}

		excelFile, err := kvuc.service.ExportKPIValueUpdate(kpiValueUpdates.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="KPI value updates.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, kpiValueUpdates)
}

func (kvuc *KPIValueUpdateController) GetKPIValueUpdateExport(c *gin.Context) {
	c.Set("export", true)
	kvuc.GetAllKPIValueUpdates(c)
}
func (kvuc *KPIValueUpdateController) GetKPIValueUpdateDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kpiRegisterUpdates, err := kvuc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kpiRegisterUpdates)
}

func (kvuc *KPIValueUpdateController) UpdateKPIValueUpdate(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kpiValueUpdateFormDTO kpiDTOs.KPIValueUpdateFormDTO
	if err := c.ShouldBindJSON(&kpiValueUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	err = kvuc.service.Update(uint(kpiValueUpdateID), currentUserID, &kpiValueUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (kvuc *KPIValueUpdateController) GetKPIValueUpdateDocumentPreview(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := kvuc.service.GetDocumentPath(uint(kpiValueUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (kvuc *KPIValueUpdateController) DownloadKPIValueUpdateDocument(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := kvuc.service.GetDocumentPath(uint(kpiValueUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}

func (kvuc *KPIValueUpdateController) EvaluateKPIValueUpdateInputs(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kpiValueUpdateEvaluationFormDTO kpiDTOs.KPIValueUpdateEvaluationFormDTO
	if err := c.ShouldBindJSON(&kpiValueUpdateEvaluationFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	response, err := kvuc.service.EvaluateKPIValueUpdateInputs(uint(kpiValueUpdateID), &kpiValueUpdateEvaluationFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, response)
}

func (kvuc *KPIValueUpdateController) SubmitKPIValueUpdate(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	if err := kvuc.service.Submit(uint(kpiValueUpdateID), currentUserID); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (kvuc *KPIValueUpdateController) ApproveKPIValueUpdate(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form kpiDTOs.KPIValueUpdateCommentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	if err := kvuc.service.Approve(uint(kpiValueUpdateID), currentUserID, &form); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (kvuc *KPIValueUpdateController) RevertKPIValueUpdate(c *gin.Context) {
	kpiValueUpdateID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var form kpiDTOs.KPIValueUpdateCommentFormDTO
	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	if err := kvuc.service.Revert(uint(kpiValueUpdateID), currentUserID, &form); err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
