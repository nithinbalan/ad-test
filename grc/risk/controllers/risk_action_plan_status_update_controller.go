package controllers

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RiskActionPlanStatusUpdateController struct {
	service services.RiskActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewRiskActionPlanStatusUpdateController(service services.RiskActionPlanStatusUpdateServiceInterface) *RiskActionPlanStatusUpdateController {
	return &RiskActionPlanStatusUpdateController{
		service: service,
	}
}

func (rasuc *RiskActionPlanStatusUpdateController) GetAllRiskActionPlanStatusUpdates(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	riskStatuses, err := rasuc.service.GetAll(uint(riskID), uint(actionPlanID), sort, sortBy, keyword)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, riskStatuses)
}

func (rasuc *RiskActionPlanStatusUpdateController) CreateRiskActionPlanStatusUpdate(c *gin.Context) {
	var riskStatusUpdateFormDTO riskDTOs.RiskActionPlanStatusUpdateFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
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

	if err := c.ShouldBindJSON(&riskStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasuc.db, c, err.Error())
		return
	}

	createdID, err := rasuc.service.Create(uint(riskID), uint(actionPlanID), currentUserID, &riskStatusUpdateFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (rasuc *RiskActionPlanStatusUpdateController) UpdateRiskActionPlanStatusUpdate(c *gin.Context) {
	var riskStatusUpdateFormDTO riskDTOs.RiskActionPlanStatusUpdateFormDTO

	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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

	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&riskStatusUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(rasuc.db, c, err.Error())
		return
	}

	err = rasuc.service.Update(uint(riskID), uint(actionPlanID), uint(id), currentUserID, &riskStatusUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(rasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (rasuc *RiskActionPlanStatusUpdateController) GetLatestRiskActionPlanStatusUpdate(c *gin.Context) {
	riskIDStr := c.Param("riskId")

	riskID, err := strconv.ParseUint(riskIDStr, 10, 64)
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

	riskActionPlanStatusUpdate, err := rasuc.service.GetLatestRiskActionPlanStatusUpdate(uint(riskID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(rasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, riskActionPlanStatusUpdate)
}

func (rasuc *RiskActionPlanStatusUpdateController) GetDocumentPreview(c *gin.Context) {
    // Parse parameters
    riskIDStr := c.Param("riskId")
    actionPlanIDStr := c.Param("actionPlanId")
    statusUpdateIDStr := c.Param("statusUpdateId")
    token := c.Query("token")

    fmt.Printf("Parameters received - RiskID: %s, ActionPlanID: %s, StatusUpdateID: %s, Token: %s\n", riskIDStr, actionPlanIDStr, statusUpdateIDStr, token)

    // Convert params to uint
    riskID, err := strconv.ParseUint(riskIDStr, 10, 32)
    if err != nil {
        fmt.Printf("Failed to parse riskID: %v\n", err)
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

    // Call service to get the document path
    documentPath, err := rasuc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(statusUpdateID), token)
    if err != nil {
        fmt.Printf("Error in service call for document path: %v\n", err)
        utility.SendNotFoundErrResp(c)
        return
    }

    // Send the document path in the response
    c.File(documentPath)
}


// Document Download
func (rasuc *RiskActionPlanStatusUpdateController) DownloadRiskActionPlanDocument(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("riskId"), 10, 32)
	actionPlanID, _ := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	statusUpdateID, _ := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	token := c.Query("token")

	documentPath, err := rasuc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(statusUpdateID), token)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	fileName := filepath.Base(documentPath)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(documentPath)
}