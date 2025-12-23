package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
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

	var themeFormDTO dtos.ThemeFormDTO

	if err := c.ShouldBindJSON(&themeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	id, err := tc.service.Create(currentUserID, &themeFormDTO)
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

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (tc *ThemeController) UpdateTheme(c *gin.Context) {

	var themeFormDTO *dtos.ThemeFormDTO

	idStr := c.Param("id")

	themeID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	if err := c.ShouldBindJSON(&themeFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	err = tc.service.Update(uint(themeID), currentUserID, themeFormDTO)
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

func (tc *ThemeController) DeleteTheme(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(tc.db, c, err.Error())
		return
	}

	failedID, err := tc.service.Delete(uint(id), currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, e.Error())
			return
		default:
			utility.SendInternalErrResp(tc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"failed_id": failedID})
}
