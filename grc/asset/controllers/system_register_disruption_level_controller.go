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

type SystemRegisterDisruptionLevelController struct {
	service services.SystemRegisterDisruptionLevelServiceInterface
	db      *gorm.DB
}

func NewSystemRegisterDisruptionLevelController(service services.SystemRegisterDisruptionLevelServiceInterface) *SystemRegisterDisruptionLevelController {
	return &SystemRegisterDisruptionLevelController{
		service: service,
	}
}

func (srdlc *SystemRegisterDisruptionLevelController) GetAllSystemRegisterDisruptionLevels(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	allowedParams := []string{}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	queryParams := assetDTOs.SystemRegisterDisruptionLevelListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
		Filters: filters,
	}
	levels, err := srdlc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, levels)
}

func (srdlc *SystemRegisterDisruptionLevelController) GetSystemRegisterDisruptionLevelDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	level, err := srdlc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, level)
}

func (srdlc *SystemRegisterDisruptionLevelController) CreateSystemRegisterDisruptionLevel(c *gin.Context) {
	var level *assetDTOs.SystemRegisterDisruptionLevelFormDTO

	if err := c.ShouldBindJSON(&level); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}
	id, err := srdlc.service.Create(currentUserID, level)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(srdlc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (srdlc *SystemRegisterDisruptionLevelController) UpdateSystemRegisterDisruptionLevel(c *gin.Context) {
	levelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedLevel *assetDTOs.SystemRegisterDisruptionLevelFormDTO
	if err := c.ShouldBindJSON(&updatedLevel); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	err = srdlc.service.Update(uint(levelID), currentUserID, updatedLevel)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(srdlc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (srdlc *SystemRegisterDisruptionLevelController) DeleteSystemRegisterDisruptionLevels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	failedIDs, err := srdlc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (srdlc *SystemRegisterDisruptionLevelController) ActivateSystemRegisterDisruptionLevels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	err = srdlc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (srdlc *SystemRegisterDisruptionLevelController) DeactivateSystemRegisterDisruptionLevels(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	err = srdlc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(srdlc.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
