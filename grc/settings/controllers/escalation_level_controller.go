package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	settingDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EscalationLevelController struct {
	service services.EscalationLevelServiceInterface
	db      *gorm.DB
}

func NewEscalationLevelController(service services.EscalationLevelServiceInterface) *EscalationLevelController {
	return &EscalationLevelController{
		service: service,
	}
}

func (rrc *EscalationLevelController) GetAllEscalationLevels(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "day_offset")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"escalation_ids",
		"day_offset",
	}

	validSortFields := map[string]string{
		"escalation_title_en": "escalations.title_en",
		"escalation_title_ar": "escalations.title_ar",
		"day_offset":          "escalation_levels.day_offset",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := settingDTOs.EscalationLevelListQueryParamsDTO{
		Sort:    sort,
		SortBy:  dbField,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	escalationLevels, err := rrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, escalationLevels)
}

func (rrc *EscalationLevelController) GetEscalationLevelDetails(c *gin.Context) {
	idStr := c.Param("escalationLevelId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	escalationLevel, err := rrc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, escalationLevel)
}

func (rrc *EscalationLevelController) CreateEscalationLevel(c *gin.Context) {
	var escalationLevelFormDTO settingDTOs.EscalationLevelFormDTO

	if err := c.ShouldBindJSON(&escalationLevelFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	id, err := rrc.service.Create(currentUserID, &escalationLevelFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (rrc *EscalationLevelController) UpdateEscalationLevel(c *gin.Context) {
	escalationLevelID, err := strconv.Atoi(c.Param("escalationLevelId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var escalationLevelFormDTO settingDTOs.EscalationLevelFormDTO
	if err := c.ShouldBindJSON(&escalationLevelFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Update(uint(escalationLevelID), currentUserID, &escalationLevelFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rrc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (rrc *EscalationLevelController) DeleteEscalationLevels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	failedIDs, err := rrc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
