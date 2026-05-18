package controllers

import (
	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BCMProcessTypeController struct {
	service services.BCMProcessTypeServiceInterface
	db      *gorm.DB
}

func NewBCMProcessTypeController(service services.BCMProcessTypeServiceInterface) *BCMProcessTypeController {
	return &BCMProcessTypeController{
		service: service,
	}
}

func (controller *BCMProcessTypeController) GetAllBCMProcessTypes(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	queryParams := bcmDTOs.BCMProcessTypeListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
		Status:  status,
	}
	tools, err := controller.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, tools)
}

func (controller *BCMProcessTypeController) GetBCMProcessTypeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	tool, err := controller.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, tool)
}

func (controller *BCMProcessTypeController) CreateBCMProcessType(c *gin.Context) {
	var formDTO *bcmDTOs.BCMProcessTypeFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	id, err := controller.service.Create(currentUserID, formDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(controller.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusCreated, id)
}

func (controller *BCMProcessTypeController) UpdateBCMProcessType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO *bcmDTOs.BCMProcessTypeFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	err = controller.service.Update(uint(id), currentUserID, formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(controller.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (controller *BCMProcessTypeController) DeleteBCMProcessTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	failedIDs, err := controller.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (controller *BCMProcessTypeController) ActivateBCMProcessTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	err = controller.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (controller *BCMProcessTypeController) DeactivateBCMProcessTypes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	err = controller.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
