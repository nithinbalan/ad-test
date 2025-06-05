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

type ReminderRuleController struct {
	service services.ReminderRuleServiceInterface
	db      *gorm.DB
}

func NewReminderRuleController(service services.ReminderRuleServiceInterface) *ReminderRuleController {
	return &ReminderRuleController{
		service: service,
	}
}

func (rrc *ReminderRuleController) GetAllReminderRules(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title_en")
	keyword := c.Query("keyword")
	status := c.Query("status")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"reminder_ids",
		"trigger_condition_ids",
	}

	validSortFields := map[string]string{
		"title_en":                   "reminder_rules.title_en",
		"title_ar":                   "reminder_rules.title_ar",
		"reminder_title_en":          "reminders.title_en",
		"reminder_title_ar":          "reminders.title_ar",
		"trigger_condition_title_en": "trigger_conditions.title_en",
		"trigger_condition_title_ar": "trigger_conditions.title_ar",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := settingDTOs.ReminderRuleListQueryParamsDTO{
		Sort:    sort,
		SortBy:  dbField,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}

	reminderRules, err := rrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, reminderRules)
}

func (rrc *ReminderRuleController) GetReminderRuleDetails(c *gin.Context) {
	idStr := c.Param("reminderRuleId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reminderRule, err := rrc.service.GetFullDetailsByID(uint(id))
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
	c.JSON(http.StatusOK, reminderRule)
}

func (rrc *ReminderRuleController) CreateReminderRule(c *gin.Context) {
	var reminderRuleFormDTO settingDTOs.ReminderRuleFormDTO

	if err := c.ShouldBindJSON(&reminderRuleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	id, err := rrc.service.Create(currentUserID, &reminderRuleFormDTO)

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

func (rrc *ReminderRuleController) UpdateReminderRule(c *gin.Context) {
	reminderRuleID, err := strconv.Atoi(c.Param("reminderRuleId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var reminderRuleFormDTO settingDTOs.ReminderRuleFormDTO
	if err := c.ShouldBindJSON(&reminderRuleFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rrc.db, c, err.Error())
		return
	}

	err = rrc.service.Update(uint(reminderRuleID), currentUserID, &reminderRuleFormDTO)
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

func (rrc *ReminderRuleController) DeleteReminderRules(c *gin.Context) {
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
