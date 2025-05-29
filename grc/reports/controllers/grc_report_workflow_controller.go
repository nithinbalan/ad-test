package controllers

import (
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/services"
	"jethur/utility"
	"net/http"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GRCReportWorkflowController struct {
	service services.GRCReportWorkflowServiceInterface
	db      *gorm.DB
}

func NewGRCReportWorkflowController(service services.GRCReportWorkflowServiceInterface) *GRCReportWorkflowController {
	return &GRCReportWorkflowController{
		service: service,
	}
}

func (grwc *GRCReportWorkflowController) GetGRCReportWorkflowDetails(c *gin.Context) {
	grcReportWorkflowID, err := strconv.Atoi(c.Param("grcReportId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	details, err := grwc.service.GetFullDetailsByID(uint(grcReportWorkflowID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(grwc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, details)
}

func (grwc *GRCReportWorkflowController) UpdateGRCReportSectionWorkflowReviewUser(c *gin.Context) {
	var request dtos.GRCReportWorkflowFormDTO

	reportID, err := strconv.Atoi(c.Param("grcReportId"))
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
		utility.SendInternalErrResp(grwc.db, c, err.Error())
		return
	}

	grcReportWorkflowID, err := grwc.service.UpdateGRCReportSectionWorkflowReviewUser(uint(reportID), currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: grcReportWorkflowID})
}

func (grwc *GRCReportWorkflowController) UpdateGRCReportSectionWorkflowApproveUser(c *gin.Context) {
	var request dtos.GRCReportWorkflowFormDTO

	reportID, err := strconv.Atoi(c.Param("grcReportId"))
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
		utility.SendInternalErrResp(grwc.db, c, err.Error())
		return
	}

	grcReportWorkflowID, err := grwc.service.UpdateGRCReportSectionWorkflowApproveUser(uint(reportID), currentUserID, &request)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(grwc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: grcReportWorkflowID})
}

func (grwc *GRCReportWorkflowController) GRCReportSectionSubmit(c *gin.Context) {
	var formDTO dtos.GRCReportSectionReviewApproveFormFormDTO

	reportID, err := strconv.Atoi(c.Param("grcReportId"))
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
		utility.SendInternalErrResp(grwc.db, c, err.Error())
		return
	}

	err = grwc.service.GRCReportSectionSubmit(uint(reportID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			// Return a 403 Forbidden response with the error message
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(grwc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
