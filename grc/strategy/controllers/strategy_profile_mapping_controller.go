package controllers

import (
	"jethur/grc/strategy/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrategyProfileMappingController struct {
	service services.StrategyProfileMappingServiceInterface
	db      *gorm.DB
}

func NewStrategyProfileMappingController(service services.StrategyProfileMappingServiceInterface) *StrategyProfileMappingController {
	return &StrategyProfileMappingController{
		service: service,
	}
}

func (spmc *StrategyProfileMappingController) GetStrategyProfile(cxt *gin.Context) {
	idParam := cxt.Param("strategyProfileId")
	profileID, err := strconv.Atoi(idParam)
	if err != nil {
		utility.SendInvalidErrResp(cxt, err)
		return
	}

	result, err := spmc.service.GetStrategyProfile(uint(profileID))
	if err != nil {
		utility.SendInternalErrResp(nil, cxt, err.Error())
		return
	}

	cxt.JSON(http.StatusOK, result)
}
