package controllers

import (
	"jethur/grc/awareness-reporting/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AwarenessAssessmentThemeController struct {
	service services.AwarenessAssessmentThemeServiceInterface
	db      *gorm.DB
}

func NewAwarenessAssessmentThemeController(service services.AwarenessAssessmentThemeServiceInterface) *AwarenessAssessmentThemeController {
	return &AwarenessAssessmentThemeController{
		service: service,
	}
}

func (aatc *AwarenessAssessmentThemeController) GetAllAwarenessAssessmentThemes(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	themes, err := aatc.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, themes)
}

func (aatc *AwarenessAssessmentThemeController) GetAwarenessAssessmentThemeDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	theme, err := aatc.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, theme)
}

func (aatc *AwarenessAssessmentThemeController) CreateAwarenessAssessmentTheme(ctx *gin.Context) {
	var theme *models.AwarenessAssessmentTheme

	if err := ctx.ShouldBindJSON(&theme); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	createdTheme, err := aatc.service.Create(currentUserID, theme)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(aatc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, createdTheme)
}

func (aatc *AwarenessAssessmentThemeController) UpdateAwarenessAssessmentTheme(ctx *gin.Context) {
	var updatedTheme *models.AwarenessAssessmentTheme

	themeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&updatedTheme); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	updatedTheme, err = aatc.service.Update(uint(themeID), currentUserID, updatedTheme)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(aatc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updatedTheme)
}

func (aatc *AwarenessAssessmentThemeController) DeleteAwarenessAssessmentThemes(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	failedIDs, err := aatc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (aatc *AwarenessAssessmentThemeController) ActivateAwarenessAssessmentThemes(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	err = aatc.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (aatc *AwarenessAssessmentThemeController) DeactivateAwarenessAssessmentThemes(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	err = aatc.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(aatc.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}