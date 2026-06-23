package controllers

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemRegisterServiceImportanceController struct {
	service services.SystemRegisterServiceImportanceServiceInterface
	db      *gorm.DB
}

func NewSystemRegisterServiceImportanceController(service services.SystemRegisterServiceImportanceServiceInterface) *SystemRegisterServiceImportanceController {
	return &SystemRegisterServiceImportanceController{
		service: service,
	}
}

func (srrtc *SystemRegisterServiceImportanceController) GetAllSystemRegisterServiceImportances(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := assetDTOs.SystemRegisterServiceImportanceListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	tools, err := srrtc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, tools)
}

func (srrtc *SystemRegisterServiceImportanceController) GetSystemRegisterServiceImportanceDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	tool, err := srrtc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, tool)
}

func (srrtc *SystemRegisterServiceImportanceController) CreateSystemRegisterServiceImportance(c *gin.Context) {
	var tool *assetDTOs.SystemRegisterServiceImportanceFormDTO

	if err := c.ShouldBindJSON(&tool); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}
	id, err := srrtc.service.Create(currentUserID, tool)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(srrtc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (srrtc *SystemRegisterServiceImportanceController) UpdateSystemRegisterServiceImportance(c *gin.Context) {
	toolID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedTool *assetDTOs.SystemRegisterServiceImportanceFormDTO
	if err := c.ShouldBindJSON(&updatedTool); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	err = srrtc.service.Update(uint(toolID), currentUserID, updatedTool)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(srrtc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (srrtc *SystemRegisterServiceImportanceController) DeleteSystemRegisterServiceImportances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	failedIDs, err := srrtc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (srrtc *SystemRegisterServiceImportanceController) ActivateSystemRegisterServiceImportances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	err = srrtc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (srrtc *SystemRegisterServiceImportanceController) DeactivateSystemRegisterServiceImportances(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	err = srrtc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srrtc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
