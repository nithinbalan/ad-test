package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ControlReviewController struct {
	service services.ControlReviewServiceInterface
	db      *gorm.DB
}

func NewControlReviewController(service services.ControlReviewServiceInterface) *ControlReviewController {
	return &ControlReviewController{
		service: service,
	}
}

func (crc *ControlReviewController) GetAllControlReviews(c *gin.Context) {
	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	controlReviews, err := crc.service.GetAll(uint(controlID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, controlReviews)
}

func (crc *ControlReviewController) CreateControlReview(c *gin.Context) {
	var controlReviewFormDTO controlDTOs.ControlReviewFormDTO

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	createdID, err := crc.service.Create(uint(controlID), currentUserID, &controlReviewFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (crc *ControlReviewController) UpdateControlReview(c *gin.Context) {
	var controlReviewFormDTO controlDTOs.ControlReviewFormDTO

	controlIDStr := c.Param("controlId")

	controlID, err := strconv.ParseUint(controlIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	reviewIDStr := c.Param("id")

	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&controlReviewFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(crc.db, c, err.Error())
		return
	}

	err = crc.service.Update(uint(controlID), uint(reviewID), currentUserID, &controlReviewFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(crc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
