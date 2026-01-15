package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jethur/grc/masters/dtos"
	"jethur/grc/masters/services"
	"jethur/utility"
	"net/http"
	"strconv"
)

type ActionPlanMasterController struct {
	service services.ActionPlanMasterServiceInterface
	db      *gorm.DB
}

func NewActionPlanMasterController(
	service services.ActionPlanMasterServiceInterface,
) *ActionPlanMasterController {
	return &ActionPlanMasterController{
		service: service,
	}
}

func (apmc *ActionPlanMasterController) GetAllActionPlanMasters(c *gin.Context) {

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := &dtos.ActionPlanMasterListQueryParamsDTO{
		Sort:     sort,
		SortBy:   sortBy,
		Keyword:  keyword,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	}

	response, err := apmc.service.GetAll(queryParams)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}


func (apmc *ActionPlanMasterController) GetActionPlanMasterFullDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlan, err := apmc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, actionPlan)
}

func (apc *ActionPlanMasterController) CreateActionPlanMaster(c *gin.Context) {
	var form *dtos.ActionPlanMasterFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apc.db, c, err.Error())
		return
	}

	id, err := apc.service.Create(uint(currentUserID), form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, id)
}

func (apmc *ActionPlanMasterController) UpdateActionPlanMaster(c *gin.Context) {
	var form dtos.ActionPlanMasterFormDTO

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	err = apmc.service.Update(uint(id), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(apmc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (apmc *ActionPlanMasterController) DeleteActionPlanMasters(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	failedIDs, err := apmc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}

func (apmc *ActionPlanMasterController) ActivateActionPlanMasters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	err = apmc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (apmc *ActionPlanMasterController) DeactivateActionPlanMasters(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	err = apmc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(apmc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
