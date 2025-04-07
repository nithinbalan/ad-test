package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StakeholderActionPlanController struct {
	service services.StakeholderActionPlanServiceInterface
	db      *gorm.DB
}

func NewStakeholderActionPlanController(service services.StakeholderActionPlanServiceInterface) *StakeholderActionPlanController {
	return &StakeholderActionPlanController{
		service: service,
	}
}

func (sapc *StakeholderActionPlanController) GetAllStakeholderActionPlans(c *gin.Context) {
	idStr := c.Param("stakeholderId")

	stakeholderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "title")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	allowedParams := []string{
		"bcp_phase_ids",
		"process_ids",
		"responsible_user_ids",
		"severity_level_ids",
		"action_plan_status_id",
	}
	validSortFields := map[string]string{
		"title":                          "stakeholder_action_plans.title",
		"responsible_user_first_name_en": "responsible_user.first_name_en",
		"responsible_user_first_name_ar": "responsible_user.first_name_ar",
		"responsible_user_last_name_en":  "responsible_user.last_name_en",
		"responsible_user_last_name_ar":  "responsible_user.last_name_ar",
		"responsible_user_email":         "responsible_user.email",
		"responsible_unit_title_en":      "organization_units.title_en",
		"responsible_unit_title_ar":      "organization_units.title_ar",
		"severity_level_title_en":        "severity_levels.title_en",
		"severity_level_title_ar":        "severity_levels.title_ar",
		"action_plan_status_title_en":    "action_plan_statuses.title_en",
		"action_plan_status_title_ar":    "action_plan_statuses.title_ar",
		"created_user_first_name_en":     "created_user.first_name_en",
		"created_user_last_name_en":      "created_user.last_name_en",
		"updated_user_first_name_en":     "updated_user.first_name_en",
		"updated_user_last_name_en":      "updated_user.last_name_en",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	sortField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationalContextDTOs.StakeholderActionPlanQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortField,
		Keyword:     keyword,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	actionPlans, err := sapc.service.GetAll(uint(stakeholderID), &queryParams)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, actionPlans)
}

func (sapc *StakeholderActionPlanController) GetByID(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := sapc.service.GetByID(uint(stakeholderID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (sapc *StakeholderActionPlanController) GetStakeholderActionPlanDetails(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("actionPlanId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := sapc.service.GetFullDetailsByID(uint(stakeholderID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, actionPlan)
}

func (sapc *StakeholderActionPlanController) CreateStakeholderActionPlan(c *gin.Context) {
	idStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO organizationalContextDTOs.StakeholderActionPlanFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	id, err := sapc.service.Create(uint(stakeholderID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (sapc *StakeholderActionPlanController) UpdateStakeholderActionPlan(c *gin.Context) {
	idStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO organizationalContextDTOs.StakeholderActionPlanFormDTO

	actionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
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
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	err = sapc.service.Update(uint(stakeholderID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sapc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (sapc *StakeholderActionPlanController) DeleteStakeholderActionPlans(c *gin.Context) {
	idStr := c.Param("stakeholderId")
	stakeholderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	failedIDs, err := sapc.service.Delete(uint(stakeholderID), request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sapc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sapc *StakeholderActionPlanController) GetDocumentPreview(c *gin.Context) {
	stakeholderIDStr := c.Param("stakeholderId")
	actionPlanIDStr := c.Param("actionPlanId")
	token := c.Query("token")

	stakeholderID, err := strconv.ParseUint(stakeholderIDStr, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse stakeholderID: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse actionPlanID: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := sapc.service.GetDocumentPath(uint(stakeholderID), uint(actionPlanID), token)
	if err != nil {
		fmt.Printf("Error in service call for document path: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}

	c.File(documentPath)
}

func (sapc *StakeholderActionPlanController) DownloadStakeholderActionPlanDocument(c *gin.Context) {
	stakeholderID, _ := strconv.ParseUint(c.Param("stakeholderId"), 10, 32)
	actionPlanID, _ := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	token := c.Query("token")

	documentPath, err := sapc.service.GetDocumentPath(uint(stakeholderID), uint(actionPlanID), token)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileName := filepath.Base(documentPath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
