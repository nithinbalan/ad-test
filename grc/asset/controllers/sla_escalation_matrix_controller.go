package controllers

import (
	"jethur/grc/models"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SLAEscalationMatrixController struct {
	service services.SLAEscalationMatrixServiceInterface
	db      *gorm.DB
}

func NewSLAEscalationMatrixController(service services.SLAEscalationMatrixServiceInterface) *SLAEscalationMatrixController {
	return &SLAEscalationMatrixController{
		service: service,
	}
}

func (sems *SLAEscalationMatrixController) GetAllMatrices(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	matrices, err := sems.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, matrices)
}

func (sems *SLAEscalationMatrixController) GetMatrixDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	matrix, err := sems.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, matrix)
}

func (sems *SLAEscalationMatrixController) CreateMatrix(c *gin.Context) {
	var matrix *models.SLAEscalationMatrix
	if err := c.ShouldBindJSON(&matrix); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	createdMatrix, err := sems.service.Create(currentUserID, matrix)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sems.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, createdMatrix)
}

func (sems *SLAEscalationMatrixController) UpdateMatrix(c *gin.Context) {
	matrixIDStr := c.Param("id")
	matrixID, err := strconv.ParseUint(matrixIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var updatedMatrix *models.SLAEscalationMatrix
	if err := c.ShouldBindJSON(&updatedMatrix); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	updatedMatrix, err = sems.service.Update(uint(matrixID), currentUserID, updatedMatrix)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(sems.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updatedMatrix)
}

func (sems *SLAEscalationMatrixController) DeleteMatrices(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	failedIDs, err := sems.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (sems *SLAEscalationMatrixController) ActivateMatrices(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	err = sems.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (sems *SLAEscalationMatrixController) DeactivateMatrices(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	err = sems.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(sems.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}