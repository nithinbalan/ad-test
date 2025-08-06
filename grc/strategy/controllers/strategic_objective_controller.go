package controllers

import (
	"errors"
	"fmt"
	"jethur/grc/strategy/dtos"
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategicObjectiveController struct {
	service services.StrategicObjectiveServiceInterface
	db      *gorm.DB
}

func NewStrategicObjectiveController(service services.StrategicObjectiveServiceInterface) *StrategicObjectiveController {
	return &StrategicObjectiveController{
		service: service,
	}
}

func (soc *StrategicObjectiveController) GetAllStrategicObjectives(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "reference_number")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	allowedParams := []string{
		"responsible_user_ids",
		"responsible_unit_ids",
		"strategic_objective_category_ids",
		"review_user_ids",
		"review_frequency_ids",
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
		"id":               "strategic_objectives.id",
		"reference_number": "strategic_objectives.reference_number",
		"title_en":         "strategic_objectives.title_en",
		"title_ar":         "strategic_objectives.title_ar",
		"start_date":       "strategic_objectives.start_date",
		"end_date":         "strategic_objectives.end_date",
	}

	sortField, valid := validSortFields[sortBy]
	if !valid {
		utility.SendInvalidErrResp(c, fmt.Errorf("invalid 'sort_by' parameter"))
		return
	}

	queryParams := dtos.StrategicObjectiveListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	list, err := soc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (soc *StrategicObjectiveController) GetStrategicObjectiveDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	strategicObjective, err := soc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, strategicObjective)
}

func (soc *StrategicObjectiveController) CreateStrategicObjective(c *gin.Context) {
	var strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO

	if err := c.ShouldBindJSON(&strategicObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}
	createdStrategicObjectiveID, err := soc.service.Create(currentUserID, strategicObjectiveFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdStrategicObjectiveID)
}

func (soc *StrategicObjectiveController) UpdateStrategicObjective(c *gin.Context) {
	strategicObjectiveID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var strategicObjectiveFormDTO *dtos.StrategicObjectiveFormDTO
	if err := c.ShouldBindJSON(&strategicObjectiveFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Update(uint(strategicObjectiveID), currentUserID, strategicObjectiveFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(soc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicObjectiveController) DeleteStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	failedIDs, err := soc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (soc *StrategicObjectiveController) ActivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicObjectiveController) DeactivateStrategicObjectives(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	err = soc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (soc *StrategicObjectiveController) GetFocusAreasByStrategicObjectiveID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	focusAreas, err := soc.service.GetFocusAreasByStrategicObjectiveID(uint(id))
	if err != nil {
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, focusAreas)
}

func (soc *StrategicObjectiveController) GetRiskSummary(c *gin.Context) {
	// Parse strategic objective ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	// Get risk summary from service
	summary, err := soc.service.GetRiskSummary(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.SendNotFoundErrResp(c)
			return
		}
		utility.SendInternalErrResp(soc.db, c, err.Error())
		return
	}

	// Return the risk summary
	c.JSON(http.StatusOK, summary)
}
