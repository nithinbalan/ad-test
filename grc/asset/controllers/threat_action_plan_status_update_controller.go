package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThreatActionPlanStatusUpdateController struct {
	service services.ThreatActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewThreatActionPlanStatusUpdateController(service services.ThreatActionPlanStatusUpdateServiceInterface) *ThreatActionPlanStatusUpdateController {
	return &ThreatActionPlanStatusUpdateController{
		service: service,
	}
}

func (tasuc *ThreatActionPlanStatusUpdateController) GetAllThreatActionPlanStatusUpdates(c *gin.Context) {
	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 32)
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

	threatStatuses, err := tasuc.service.GetAll(uint(threatID), uint(actionPlanID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, threatStatuses)
}

func (tasuc *ThreatActionPlanStatusUpdateController) CreateThreatActionPlanStatusUpdate(c *gin.Context) {
	var threatStatusUpdateFormDTO threatDTOs.ThreatActionPlanStatusUpdateFormDTO

	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&threatStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tasuc.db, c, err.Error())
		return
	}

	createdID, err := tasuc.service.Create(uint(threatID), uint(actionPlanID), currentUserID, &threatStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (tasuc *ThreatActionPlanStatusUpdateController) UpdateThreatActionPlanStatusUpdate(c *gin.Context) {
	var threatStatusUpdateFormDTO threatDTOs.ThreatActionPlanStatusUpdateFormDTO

	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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

	idStr := c.Param("actionPlanStatusUpdateId")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&threatStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tasuc.db, c, err.Error())
		return
	}

	err = tasuc.service.Update(uint(threatID), uint(actionPlanID), uint(id), currentUserID, &threatStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(tasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (tasuc *ThreatActionPlanStatusUpdateController) GetLatestThreatActionPlanStatusUpdate(c *gin.Context) {
	threatIDStr := c.Param("id")

	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
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

	threatActionPlanStatusUpdate, err := tasuc.service.GetLatestThreatActionPlanStatusUpdate(uint(threatID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(tasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, threatActionPlanStatusUpdate)
}

func (tasuc *ThreatActionPlanStatusUpdateController) GetDocumentPreview(c *gin.Context) {
	threatIDStr := c.Param("id")
	actionPlanIDStr := c.Param("actionPlanId")
	statusUpdateIDStr := c.Param("statusUpdateId")
	token := c.Query("token")

	fmt.Printf("Parameters received - ThreatID: %s, ActionPlanID: %s, StatusUpdateID: %s, Token: %s\n", threatIDStr, actionPlanIDStr, statusUpdateIDStr, token)

	threatID, err := strconv.ParseUint(threatIDStr, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse threatID: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanID, err := strconv.ParseUint(actionPlanIDStr, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse actionPlanID: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}
	statusUpdateID, err := strconv.ParseUint(statusUpdateIDStr, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse statusUpdateID: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := tasuc.service.GetDocumentPath(uint(threatID), uint(actionPlanID), uint(statusUpdateID), token)
	if err != nil {
		fmt.Printf("Error in service call for document path: %v\n", err)
		utility.SendNotFoundErrResp(c)
		return
	}

	c.File(documentPath)
}

func (tasuc *ThreatActionPlanStatusUpdateController) DownloadThreatActionPlanDocument(c *gin.Context) {
	threatID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	actionPlanID, _ := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	statusUpdateID, _ := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	token := c.Query("token")

	documentPath, err := tasuc.service.GetDocumentPath(uint(threatID), uint(actionPlanID), uint(statusUpdateID), token)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileName := filepath.Base(documentPath)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	c.File(documentPath)
}
