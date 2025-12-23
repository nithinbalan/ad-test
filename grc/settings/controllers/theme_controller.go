package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ThemeController struct {
	service services.ThemeServiceInterface
	db      *gorm.DB
}

func NewThemeController(service services.ThemeServiceInterface) *ThemeController {
	return &ThemeController{
		service: service,
	}
}

func (tc *ThemeController) GetAllThemes(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	keyword := c.Query("keyword")

	themes, err := tc.service.GetAll(keyword, sortBy, sort)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, themes)
}

func (tc *ThemeController) GetThemeDetails(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	theme, err := tc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	c.JSON(http.StatusOK, theme)
}

func (tc *ThemeController) CreateTheme(c *gin.Context) {

	var form dtos.ThemeFormDTO

	if err := c.ShouldBindJSON(&form); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	themeID, err := tc.service.Create(currentUserID, &form)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": themeID,
	})
}

func (tc *ThemeController) UpdateTheme(c *gin.Context) {

	var updateDTO *dtos.ThemeFormDTO

	idStr := c.Param("id")

	themeID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	err = tc.service.Update(uint(themeID), currentUserID, updateDTO)
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
			return
		default:
			utility.SendInternalErrResp(nil, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (tc *ThemeController) DeleteThemes(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	failedIDs, err := tc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"failed_ids": failedIDs,
	})
}
