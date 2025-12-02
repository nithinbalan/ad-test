package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"

	"jethur/grc/documents/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PolicyAcknowledgeUserController struct {
	service services.DocumentAcknowledgeUserServiceInterface
	db      *gorm.DB
}

func NewPolicyAcknowledgeUserController(service services.DocumentAcknowledgeUserServiceInterface) *PolicyAcknowledgeUserController {
	return &PolicyAcknowledgeUserController{
		service: service,
	}
}

func (pauc *PolicyAcknowledgeUserController) CreatePolicyAcknowledgeUser(c *gin.Context) {

	var documentAcknowledgeFormDTO documentDTOs.DocumentAcknowledgeFormDTO
	if err := c.ShouldBindJSON(&documentAcknowledgeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, err.Error())
		return
	}

	id, err := pauc.service.Create(uint(documentId), uint(documentVersionId), currentUserID, &documentAcknowledgeFormDTO)
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
			utility.SendInternalErrResp(pauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (pauc *PolicyAcknowledgeUserController) DeletePolicyAcknowledgeUser(c *gin.Context) {
	documentIdStr := c.Param("id")

	documentId, err := strconv.ParseUint(documentIdStr, 10, 64)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	versionIdStr := c.Param("versionId")
	documentVersionId, err := strconv.ParseUint(versionIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(pauc.db, c, err.Error())
		return
	}

	failedID, err := pauc.service.Delete(uint(documentId), uint(documentVersionId), currentUserID)
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
			utility.SendInternalErrResp(pauc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
