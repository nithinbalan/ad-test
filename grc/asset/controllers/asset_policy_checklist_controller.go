package controllers

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetPolicyChecklistController struct {
	service services.AssetPolicyChecklistServiceInterface
	db      *gorm.DB
}

func NewAssetPolicyChecklistController(service services.AssetPolicyChecklistServiceInterface) *AssetPolicyChecklistController {
	return &AssetPolicyChecklistController{service: service}
}

func (apcc *AssetPolicyChecklistController) GetAssetPolicyChecklistDetails(c *gin.Context) {
	idStr := c.Param("checklistId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetChecklist, err := apcc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, assetChecklist)
}

func (apcc *AssetPolicyChecklistController) CreateAssetPolicyChecklist(c *gin.Context) {
	var formDTO dtos.AssetPolicyChecklistFormDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	assetIdStr := c.Param("id")
	_, err := strconv.ParseUint(assetIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	id, err := apcc.service.Create(&formDTO, currentUserID, uint(assetPolicyId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ids": id})
}

func (arc *AssetPolicyChecklistController) UpdateAssetPolicyChecklist(c *gin.Context) {
	idStr := c.Param("checklistId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var assetChecklistForm dtos.AssetPolicyChecklistUpdateDTO
	if err := c.ShouldBindJSON(&assetChecklistForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	err = arc.service.Update(uint(id), &assetChecklistForm, currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(arc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (arc *AssetPolicyChecklistController) DeleteAssetPolicyChecklists(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	failedIDs, err := arc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(arc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (apcc *AssetPolicyChecklistController) AnswerAssetPolicyChecklist(c *gin.Context) {
	var policyForm dtos.AssetPolicyChecklistAnswerFormDTO
	if err := c.ShouldBindJSON(&policyForm); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	assetIdStr := c.Param("id")
	assetId, err := strconv.ParseUint(assetIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	checklistIdStr := c.Param("checklistId")
	checklistId, err := strconv.ParseUint(checklistIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = apcc.service.AnswerAssetPolicyChecklist(uint(checklistId), &policyForm, currentUserID, uint(assetId), uint(assetPolicyId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.Status(http.StatusOK)
}

func (apcc *AssetPolicyChecklistController) DeleteAnswerAssetPolicyChecklist(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	assetIdStr := c.Param("id")
	assetId, err := strconv.ParseUint(assetIdStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	assetPolicyIdStr := c.Param("policyId")
	assetPolicyId, err := strconv.ParseUint(assetPolicyIdStr, 10, 32)
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
		utility.SendNotFoundErrResp(c)
		return
	}

	failedIDs, err := apcc.service.DeleteAnswerAssetPolicyChecklist(request.IDs, currentUserID, uint(assetId), uint(assetPolicyId))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (apcc *AssetPolicyChecklistController) GetAssetPolicyChecklistSummary(c *gin.Context) {
	assetRegisterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}
	policyId, err := strconv.Atoi(c.Param("policyId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	summary, err := apcc.service.GetAssetPolicyChecklistSummary(uint(assetRegisterID), uint(policyId))
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, summary)
}
