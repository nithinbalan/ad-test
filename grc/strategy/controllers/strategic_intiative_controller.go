package controllers

import (
	"fmt"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicIntiativeController struct {
	service services.StrategicIntiativeServiceInterface
	db      *gorm.DB
}

func NewStrategicIntiativeController(service services.StrategicIntiativeServiceInterface) *StrategicIntiativeController {
	return &StrategicIntiativeController{
		service: service,
	}
}

func (soc *StrategicIntiativeController) GetAllStrategicIntiatives(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Convert the export parameter to a boolean
	exportFromContext, exists := c.Get("export")
	var export bool
	if exists {
		export = exportFromContext.(bool)
	} else {
		export, err := strconv.ParseBool(c.DefaultQuery("export", "false"))
		if err != nil {
			utility.SendInternalErrResp(soc.db, c, "Invalid export parameter")
			return
		}
		_ = export
	}

	allowedParams := []string{
		"responsible_user_ids",
		"responsible_unit_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, "Invalid page parameter")
		return
	}
	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, "Invalid page_size parameter")
		return
	}

	validSortFields := map[string]string{
		"id":         "strategic_intiatives.id",
		"title_en":   "strategic_intiatives.title_en",
		"title_ar":   "strategic_intiatives.title_ar",
		"start_date": "strategic_intiatives.start_date",
		"end_date":   "strategic_intiatives.end_date",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.StrategicIntiativeListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
		Export:      export,
	}

	list, err := soc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	// If export is true, generate and return the Excel file
	if export {
		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendInternalErrResp(soc.db, c, err.Error())
			return
		}

		excelFile, err := soc.service.ExportStrategicInitiatives(list.Items, currentUserID)
		if err != nil {
			utility.SendInternalErrResp(soc.db, c, err.Error())
			return
		}
		defer excelFile.Close()

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", `attachment; filename="GovernanceInitiatives.xlsx"`)
		c.Header("Content-Type", "application/octet-stream")
		c.File(excelFile.Name())

		return
	}

	c.JSON(http.StatusOK, list)
}

func (sic *StrategicIntiativeController) GetStrategicIntiativeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicIntiative, err := sic.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, strategicIntiative)
}

func (sic *StrategicIntiativeController) CreateStrategicIntiative(c *gin.Context) {
	var strategicIntiativeFormDTO *dtos.StrategicIntiativeFormDTO

	if err := c.ShouldBindJSON(&strategicIntiativeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sic.db, c, err.Error())
		return
	}
	createdStrategicIntiativeID, err := sic.service.Create(currentUserID, strategicIntiativeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sic.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategicIntiativeID)
}

func (sic *StrategicIntiativeController) UpdateStrategicIntiative(c *gin.Context) {
	strategicIntiativeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var strategicIntiativeFormDTO *dtos.StrategicIntiativeFormDTO
	if err := c.ShouldBindJSON(&strategicIntiativeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sic.db, c, err.Error())
		return
	}

	err = sic.service.Update(uint(strategicIntiativeID), currentUserID, strategicIntiativeFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sic *StrategicIntiativeController) DeleteStrategicIntiatives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sic.db, c, err.Error())
		return
	}

	failedIDs, err := sic.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sic.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sic *StrategicIntiativeController) GetStrategicInitiativesExport(c *gin.Context) {
	c.Set("export", true)

	sic.GetAllStrategicIntiatives(c)
}

func (sic *StrategicIntiativeController) UpdateStrategicIntiativeStatusAndCompletion(c *gin.Context) {
	strategicIntiativeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var dto *dtos.StrategicIntiativeStatusCompletionUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sic.db, c, err.Error())
		return
	}

	err = sic.service.UpdateStatusAndCompletion(uint(strategicIntiativeID), currentUserID, dto)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sic.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
