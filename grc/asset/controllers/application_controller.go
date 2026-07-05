package controllers

import (
	"jethur/grc/asset/services"
	"jethur/grc/models"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationController struct {
	service services.ApplicationServiceInterface
	db      *gorm.DB
}

func NewApplicationController(service services.ApplicationServiceInterface) *ApplicationController {
	return &ApplicationController{
		service: service,
	}
}

func (ac *ApplicationController) GetAllApplications(ctx *gin.Context) {
	sort := ctx.Query("sort_order")
	sortBy := ctx.Query("sort_by")
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")

	applications, err := ac.service.GetAll(sort, sortBy, keyword, status)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, applications)
}

func (ac *ApplicationController) GetApplicationDetails(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	application, err := ac.service.GetByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	ctx.JSON(http.StatusOK, application)
}

func (ac *ApplicationController) CreateApplication(ctx *gin.Context) {
	var application *models.Application

	if err := ctx.ShouldBindJSON(&application); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	createdApplication, err := ac.service.Create(currentUserID, application)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(ac.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusCreated, createdApplication)
}

func (ac *ApplicationController) UpdateApplication(ctx *gin.Context) {
	var updatedApplication *models.Application

	applicationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&updatedApplication); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	updatedApplication, err = ac.service.Update(uint(applicationID), currentUserID, updatedApplication)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		default:
			utility.SendInternalErrResp(ac.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, updatedApplication)
}

func (ac *ApplicationController) DeleteApplications(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	failedIDs, err := ac.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (ac *ApplicationController) ActivateApplications(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	err = ac.service.Activate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (ac *ApplicationController) DeactivateApplications(ctx *gin.Context) {
	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(ctx)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	err = ac.service.Deactivate(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(ac.db, ctx, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}