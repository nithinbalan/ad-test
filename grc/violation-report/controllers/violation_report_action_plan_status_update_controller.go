package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	commonDTOs "jethur/grc/common/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ViolationReportActionPlanStatusUpdateController struct {
	service services.ViolationReportActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewViolationReportActionPlanStatusUpdateController(service services.ViolationReportActionPlanStatusUpdateServiceInterface) *ViolationReportActionPlanStatusUpdateController {
	return &ViolationReportActionPlanStatusUpdateController{
		service: service,
	}
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) GetAllViolationReportActionPlanStatusUpdates(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	statusUpdates, err := vrasuc.service.GetAll(uint(violationReportID), uint(actionPlanID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, statusUpdates)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) CreateViolationReportActionPlanStatusUpdate(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
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
		utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		return
	}

	createdID, err := vrasuc.service.Create(uint(violationReportID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) UpdateViolationReportActionPlanStatusUpdate(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO

	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("statusUpdateId")

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
		utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		return
	}

	err = vrasuc.service.Update(uint(violationReportID), uint(actionPlanID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) GetLatestViolationReportActionPlanStatusUpdate(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdate, err := vrasuc.service.GetLatestViolationReportActionPlanStatusUpdate(uint(violationReportID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, statusUpdate)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) GetViolationReportActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("statusUpdateId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	filePath, err := vrasuc.service.GetDocumentPath(uint(violationReportID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		}
		return
	}

	c.File(filePath)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) DownloadViolationReportActionPlanStatusUpdateDocument(c *gin.Context) {
	violationReportIDStr := c.Param("id")

	violationReportID, err := strconv.ParseUint(violationReportIDStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanIDStr := c.Param("actionPlanId")

	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	idStr := c.Param("statusUpdateId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	filePath, err := vrasuc.service.GetDocumentPath(uint(violationReportID), uint(actionPlanID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(filePath)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) ApproveViolationReportActionPlan(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO

	violationReportID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
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
		utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		return
	}

	err = vrasuc.service.Approve(uint(violationReportID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (vrasuc *ViolationReportActionPlanStatusUpdateController) RevertViolationReportActionPlan(c *gin.Context) {
	var formDTO violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO

	violationReportID, err := strconv.Atoi(c.Param("violationReportID"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	actionPlanID, err := strconv.Atoi(c.Param("actionPlanId"))
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
		utility.SendInternalErrResp(vrasuc.db, c, err.Error())
		return
	}

	err = vrasuc.service.Revert(uint(violationReportID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
		default:
			utility.SendInternalErrResp(vrasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
