package controllers

import (
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

type KRIValueUpdateController struct {
	service services.KRIValueUpdateServiceInterface
	db      *gorm.DB
}

func NewKRIValueUpdateController(service services.KRIValueUpdateServiceInterface) *KRIValueUpdateController {
	return &KRIValueUpdateController{
		service: service,
	}
}

func (kvuc *KRIValueUpdateController) GetAllKRIValueUpdates(c *gin.Context) {
	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")

	reviewedDateFrom := c.Query("reviewed_date_from")
	reviewedDateTO := c.Query("reviewed_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := riskDTOs.KRIValueUpdateListQueryParamsDTO{
		Sort:             sort,
		SortBy:           sortBy,
		Keyword:          keyword,
		Status:           status,
		Page:             page,
		PageSize:         pageSize,
		ReviewedDateFrom: reviewedDateFrom,
		ReviewedDateTo:   reviewedDateTO,
	}
	kriRegisterUpdates, err := kvuc.service.GetAll(uint(kriRegisterID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, kriRegisterUpdates)
}

func (kvuc *KRIValueUpdateController) GetKRIValueUpdateDetails(c *gin.Context) {
	idStr := c.Param("updateId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriRegisterUpdates, err := kvuc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, kriRegisterUpdates)
}

func (kvuc *KRIValueUpdateController) CreateKRIValueUpdate(c *gin.Context) {
	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kriValueUpdateFormDTO riskDTOs.KRIValueUpdateFormDTO
	if err := c.ShouldBindJSON(&kriValueUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	id, err := kvuc.service.Create(uint(kriRegisterID), currentUserID, &kriValueUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (kvuc *KRIValueUpdateController) UpdateKRIValueUpdate(c *gin.Context) {
	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriValueUpdateID, err := strconv.Atoi(c.Param("updateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var kriValueUpdateFormDTO riskDTOs.KRIValueUpdateFormDTO
	if err := c.ShouldBindJSON(&kriValueUpdateFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	err = kvuc.service.Update(uint(kriRegisterID), uint(kriValueUpdateID), currentUserID, &kriValueUpdateFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (kvuc *KRIValueUpdateController) DeleteKRIValueUpdates(c *gin.Context) {
	kriRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(kvuc.db, c, err.Error())
		return
	}

	failedIDs, err := kvuc.service.Delete(uint(kriRegisterID), request.IDs, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(kvuc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (kvuc *KRIValueUpdateController) GetKRIValueUpdateDocumentPreview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriValueUpdateID, err := strconv.Atoi(c.Param("updateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := kvuc.service.GetDocumentPath(uint(id), uint(kriValueUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
		}
		return
	}
	c.File(documentPath)
}

func (kvuc *KRIValueUpdateController) DownloadKRIValueUpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	kriValueUpdateID, err := strconv.Atoi(c.Param("updateId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	token := c.Query("token")
	if token == "" {
		utility.SendNotFoundErrResp(c)
		return
	}

	documentPath, err := kvuc.service.GetDocumentPath(uint(id), uint(kriValueUpdateID), token)
	if err != nil {
		switch err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
		default:
			utility.SendInternalErrResp(kvuc.db, c, err.Error())
		}
		return
	}
	filename := filepath.Base(documentPath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	// Serve the file
	c.File(documentPath)
}
