package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IncidentClosureController struct {
	service services.IncidentClosureServiceInterface
	db      *gorm.DB
}

func NewIncidentClosureController(service services.IncidentClosureServiceInterface) *IncidentClosureController {
	return &IncidentClosureController{
		service: service,
	}
}

func (icc *IncidentClosureController) GetIncidentClosureDetails(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentClosure, err := icc.service.GetFullDetailsByID(uint(incidentID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, incidentClosure)
}

func (icc *IncidentClosureController) CreateIncidentClosure(c *gin.Context) {
	var IncidentClosureFormDTO incidentDTOs.IncidentClosureFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentClosureFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	id, err := icc.service.Create(uint(incidentID), currentUserID, &IncidentClosureFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (icc *IncidentClosureController) UpdateIncidentClosure(c *gin.Context) {
	var IncidentClosureFormDTO incidentDTOs.IncidentClosureFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&IncidentClosureFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	err = icc.service.Update(uint(incidentID), currentUserID, &IncidentClosureFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (icc *IncidentClosureController) DeleteIncidentClosures(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(icc.db, c, err.Error())
		return
	}

	failedID, err := icc.service.Delete(uint(incidentID), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(icc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
