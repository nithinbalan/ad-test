package controllers

import (
	"jethur/grc/models"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestExecutorController struct {
	service services.PenetrationTestExecutorServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestExecutorController(service services.PenetrationTestExecutorServiceInterface) *PenetrationTestExecutorController {
	return &PenetrationTestExecutorController{
		service: service,
	}
}

func (ptec *PenetrationTestExecutorController) GetAllPenetrationTestExecutors(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	penetrationTestExecutors, err := ptec.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, penetrationTestExecutors)
}

func (ptec *PenetrationTestExecutorController) GetPenetrationTestExecutorDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	penetrationTestExecutor, err := ptec.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, penetrationTestExecutor)
}

func (ptec *PenetrationTestExecutorController) CreatePenetrationTestExecutor(c *gin.Context) {
	var penetrationTestExecutor *models.PenetrationTestExecutor

	if err := c.ShouldBindJSON(&penetrationTestExecutor); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	created, err := ptec.service.Create(currentUserID, penetrationTestExecutor)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, created)
}

func (ptec *PenetrationTestExecutorController) UpdatePenetrationTestExecutor(c *gin.Context) {
	var updated *models.PenetrationTestExecutor

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updated); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	updated, err = ptec.service.Update(uint(id), currentUserID, updated)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptec.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, updated)
}

func (ptec *PenetrationTestExecutorController) DeletePenetrationTestExecutors(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	failedIDs, err := ptec.service.Delete(req.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ptec *PenetrationTestExecutorController) ActivatePenetrationTestExecutors(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	if err := ptec.service.Activate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (ptec *PenetrationTestExecutorController) DeactivatePenetrationTestExecutors(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	if err := ptec.service.Deactivate(req.IDs, currentUserID); err != nil {
		utility.SendInternalErrResp(ptec.db, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
