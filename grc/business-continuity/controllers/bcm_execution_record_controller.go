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

type BCMExecutionRecordController struct {
	service services.BCMExecutionRecordServiceInterface
	db      *gorm.DB
}

func NewBCMExecutionRecordController(service services.BCMExecutionRecordServiceInterface) *BCMExecutionRecordController {
	return &BCMExecutionRecordController{
		service: service,
	}
}

func (controller *BCMExecutionRecordController) GetAllBCMExecutionRecords(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	// Validate sort_order to prevent SQL injection
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	queryParams := bcmDTOs.BCMExecutionRecordListQueryParamsDTO{
		Sort:    sort,
		SortBy:  sortBy,
		Keyword: keyword,
	}
	records, err := controller.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(controller.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, records)
}

func (controller *BCMExecutionRecordController) GetBCMExecutionRecordFullDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	record, err := controller.service.GetFullDetailsByID(uint(id))
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

	c.JSON(http.StatusOK, record)
}

func (controller *BCMExecutionRecordController) CreateBCMExecutionRecord(c *gin.Context) {
	var formDTO *bcmDTOs.BCMExecutionRecordFormDTO

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

func (controller *BCMExecutionRecordController) UpdateBCMExecutionRecord(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var formDTO *bcmDTOs.BCMExecutionRecordFormDTO
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

func (controller *BCMExecutionRecordController) DeleteBCMExecutionRecords(c *gin.Context) {
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
