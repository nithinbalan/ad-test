package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetRiskActionPlanStatusUpdateController struct {
	service services.AssetRiskActionPlanStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewAssetRiskActionPlanStatusUpdateController(service services.AssetRiskActionPlanStatusUpdateServiceInterface) *AssetRiskActionPlanStatusUpdateController {
	return &AssetRiskActionPlanStatusUpdateController{
		service: service,
	}
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) GetAllAssetRiskActionPlanStatusUpdates(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")

	statusUpdates, err := arasuc.service.GetAll(uint(riskID), uint(actionPlanID), sort, sortBy, keyword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, statusUpdates)
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) CreateAssetRiskActionPlanStatusUpdate(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO

	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 32)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
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
		utility.SendInternalErrResp(arasuc.db, c, err.Error())
		return
	}

	createdID, err := arasuc.service.Create(uint(riskID), uint(actionPlanID), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: createdID})
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) UpdateAssetRiskActionPlanStatusUpdate(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO

	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("statusUpdateId"), 10, 64)
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
		utility.SendInternalErrResp(arasuc.db, c, err.Error())
		return
	}

	err = arasuc.service.Update(uint(riskID), uint(actionPlanID), uint(id), currentUserID, &formDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) GetLatestAssetRiskActionPlanStatusUpdate(c *gin.Context) {
	riskIDStr := c.Param("id")

	var riskID uint64
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.ParseUint(riskIDStr, 10, 64)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
	}

	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdate, err := arasuc.service.GetLatestRiskActionPlanStatusUpdate(uint(riskID), uint(actionPlanID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, statusUpdate)
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) GetAssetRiskActionPlanStatusUpdateDocumentPreview(c *gin.Context) {
	riskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	actionPlanID, err := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	statusUpdateID, err := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	token := c.Query("token")

	documentPath, err := arasuc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(statusUpdateID), token)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.File(documentPath)
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) DownloadAssetRiskActionPlanDocument(c *gin.Context) {
	riskID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	actionPlanID, _ := strconv.ParseUint(c.Param("actionPlanId"), 10, 32)
	statusUpdateID, _ := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	token := c.Query("token")

	documentPath, err := arasuc.service.GetDocumentPath(uint(riskID), uint(actionPlanID), uint(statusUpdateID), token)
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

func (arasuc *AssetRiskActionPlanStatusUpdateController) ApproveAssetRiskActionPlan(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO

	riskIDStr := c.Param("id")

	var riskID int
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.Atoi(riskIDStr)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
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
		utility.SendInternalErrResp(arasuc.db, c, err.Error())
		return
	}

	err = arasuc.service.Approve(uint(riskID), uint(actionPlanID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arasuc *AssetRiskActionPlanStatusUpdateController) RevertAssetRiskActionPlan(c *gin.Context) {
	var formDTO assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO

	riskIDStr := c.Param("id")

	var riskID int
	var err error

	if riskIDStr != "" {
		riskID, err = strconv.Atoi(riskIDStr)
		if err != nil {
			utility.SendNotFoundErrResp(c)
			return
		}
	} else {
		riskID = 0
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
		utility.SendInternalErrResp(arasuc.db, c, err.Error())
		return
	}

	err = arasuc.service.Revert(uint(riskID), uint(actionPlanID), currentUserID, &formDTO)
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
			utility.SendInternalErrResp(arasuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
