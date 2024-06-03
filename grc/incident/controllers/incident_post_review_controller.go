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

type IncidentPostReviewController struct {
	service services.IncidentPostReviewServiceInterface
	db      *gorm.DB
}

func NewIncidentPostReviewController(service services.IncidentPostReviewServiceInterface) *IncidentPostReviewController {
	return &IncidentPostReviewController{
		service: service,
	}
}

func (iprc *IncidentPostReviewController) GetIncidentPostReviewDetails(c *gin.Context) {
	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	incidentPostReview, err := iprc.service.GetFullDetailsByID(uint(incidentID), uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(iprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, incidentPostReview)
}

func (iprc *IncidentPostReviewController) CreateIncidentPostReview(c *gin.Context) {
	var incidentPostReviewFormDTO incidentDTOs.IncidentPostReviewFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentPostReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iprc.db, c, err.Error())
		return
	}

	id, err := iprc.service.Create(uint(incidentID), currentUserID, &incidentPostReviewFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (iprc *IncidentPostReviewController) UpdateIncidentPostReview(c *gin.Context) {
	var incidentPostReviewFormDTO incidentDTOs.IncidentPostReviewFormDTO

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&incidentPostReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iprc.db, c, err.Error())
		return
	}

	err = iprc.service.Update(uint(incidentID), uint(id), currentUserID, &incidentPostReviewFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iprc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (iprc *IncidentPostReviewController) DeleteIncidentPostReviews(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	incidentIDStr := c.Param("incidentID")

	incidentID, err := strconv.ParseUint(incidentIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(iprc.db, c, err.Error())
		return
	}

	failedIDs, err := iprc.service.Delete(uint(incidentID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(iprc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}
