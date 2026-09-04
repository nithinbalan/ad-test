package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/services"
	"jethur/utility"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PenetrationTestFindingStatusUpdateController struct {
	service services.PenetrationTestFindingStatusUpdateServiceInterface
	db      *gorm.DB
}

func NewPenetrationTestFindingStatusUpdateController(service services.PenetrationTestFindingStatusUpdateServiceInterface) *PenetrationTestFindingStatusUpdateController {
	return &PenetrationTestFindingStatusUpdateController{service: service}
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) GetAllPenetrationTestFindingStatusUpdates(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	values, err := ptfsuc.service.GetAll(uint(findingID), c.Query("sort_order"), c.Query("sort_by"), c.Query("keyword"))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, values)
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) GetLatestPenetrationTestFindingStatusUpdate(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	value, err := ptfsuc.service.GetLatest(uint(findingID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, value)
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) CreatePenetrationTestFindingStatusUpdate(c *gin.Context) {
	var form penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO

	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptfsuc.db, c, err.Error())
		return
	}

	id, err := ptfsuc.service.Create(uint(findingID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) UpdatePenetrationTestFindingStatusUpdate(c *gin.Context) {
	var form penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO

	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	statusUpdateID, err := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(ptfsuc.db, c, err.Error())
		return
	}

	err = ptfsuc.service.Update(uint(findingID), uint(statusUpdateID), currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) GetPenetrationTestFindingStatusUpdateDocumentPreview(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	path, err := ptfsuc.service.GetDocumentPath(uint(findingID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, err.Error())
		}
		return
	}

	c.File(path)
}

func (ptfsuc *PenetrationTestFindingStatusUpdateController) DownloadPenetrationTestFindingStatusUpdateDocument(c *gin.Context) {
	findingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("statusUpdateId"), 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	path, err := ptfsuc.service.GetDocumentPath(uint(findingID), uint(id), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(ptfsuc.db, c, err.Error())
		}
		return
	}

	filename := filepath.Base(path)

	// Set the content type and disposition headers for the file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(path)
}
