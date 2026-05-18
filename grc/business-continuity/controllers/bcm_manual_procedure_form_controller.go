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

type BCMManualProcedureFormController struct {
	service services.BCMManualProcedureFormServiceInterface
	db      *gorm.DB
}

func NewBCMManualProcedureFormController(service services.BCMManualProcedureFormServiceInterface) *BCMManualProcedureFormController {
	return &BCMManualProcedureFormController{
		service: service,
	}
}

func (controller *BCMManualProcedureFormController) GetAllBCMManualProcedureForms(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	queryParams := bcmDTOs.BCMManualProcedureFormListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	forms, err := controller.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, forms)
}

func (controller *BCMManualProcedureFormController) GetBCMManualProcedureFormFullDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	form, err := controller.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(controller.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, form)
}

func (controller *BCMManualProcedureFormController) CreateBCMManualProcedureForm(c *gin.Context) {
	var formDTO *bcmDTOs.BCMManualProcedureFormFormDTO

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

func (controller *BCMManualProcedureFormController) UpdateBCMManualProcedureForm(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO *bcmDTOs.BCMManualProcedureFormFormDTO
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

func (controller *BCMManualProcedureFormController) DeleteBCMManualProcedureForms(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required"`
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
