package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestClosureController struct {
	service services.PenetrationTestClosureServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestClosureController(service services.PenetrationTestClosureServiceInterface) *PenetrationTestClosureController {
	return &PenetrationTestClosureController{
		service: service,
	}
}

func (ptcc *PenetrationTestClosureController) GetPenetrationTestClosureDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	closure, err := ptcc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, closure)
}

func (ptcc *PenetrationTestClosureController) CreatePenetrationTestClosure(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestClosureFormDTO

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
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
		utility.SendInternalErrResp(ptcc.db, c, err.Error())
		return
	}

	createdID, err := ptcc.service.Create(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (ptcc *PenetrationTestClosureController) UpdatePenetrationTestClosure(c *gin.Context) {
	var formDTO penetrationTestDTOs.PenetrationTestClosureFormDTO

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
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
		utility.SendInternalErrResp(ptcc.db, c, err.Error())
		return
	}

	err = ptcc.service.Update(uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptcc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptcc *PenetrationTestClosureController) DeletePenetrationTestClosures(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptcc.db, c, err.Error())
		return
	}

	failedID, err := ptcc.service.Delete(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptcc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
