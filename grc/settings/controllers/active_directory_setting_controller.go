package controllers

import (
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActiveDirectorySettingController struct {
	service services.ActiveDirectorySettingServiceInterface
	db      *gorm.DB
}

func NewActiveDirectorySettingController(service services.ActiveDirectorySettingServiceInterface) *ActiveDirectorySettingController {
	return &ActiveDirectorySettingController{
		service: service,
	}
}
func (adsc *ActiveDirectorySettingController) GetActiveDirectorySettings(c *gin.Context) {
	ad, err := adsc.service.GetFirst()
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, ad)
}

func (adsc *ActiveDirectorySettingController) CreateActiveDirectorySetting(c *gin.Context) {
	var req dtos.ActiveDirectorySettingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	err = adsc.service.Create(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(c, err.Error())
			return
		case *utility.ActiveDirectoryConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, e.Error())
			return
		default:
			utility.SendInternalErrResp(adsc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)

}

func (adsc *ActiveDirectorySettingController) CheckActiveDirectoryConnection(ctx *gin.Context) {
	var adDTO dtos.ActiveDirectorySettingDTO
	if err := ctx.ShouldBindJSON(&adDTO); err != nil {
		utility.SendInvalidErrResp(ctx, err)
		return
	}
	err := adsc.service.CheckADConnection(&adDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(ctx, err.Error())
			return
		case *utility.ActiveDirectoryConnectionError:
			utility.SendCustomErrResp(ctx, http.StatusUnprocessableEntity, e.Error())
			return
		default:
			utility.SendInternalErrResp(adsc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Active Directory connection successful"})
}

func (adsc *ActiveDirectorySettingController) CheckSavedADConnection(ctx *gin.Context) {
	err := adsc.service.GetADConnection()
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(ctx)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(ctx, e)
			return
		case *utility.ForbiddenError:
			utility.SendForbiddenErrResp(ctx, err.Error())
			return
		case *utility.ActiveDirectoryConnectionError:
			utility.SendCustomErrResp(ctx, http.StatusUnprocessableEntity, e.Error())
			return
		default:
			utility.SendInternalErrResp(adsc.db, ctx, e.Error())
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Active Directory connection (from saved config) successful"})
}



func (adsc *ActiveDirectorySettingController) DeleteActiveDirectorySetting(c *gin.Context) {

    currentUserID, err := utility.GetCurrentUserID(c)
    if err != nil {
        utility.SendInternalErrResp(adsc.db, c, err.Error())
        return
    }

    err = adsc.service.Delete(currentUserID)
    if err != nil {
        utility.SendInternalErrResp(adsc.db, c, err.Error())
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Active Directory setting deleted successfully",
    })
}


// ==============================
// Entra ID (OAuth2 / SSO / Graph)
// ==============================

func (adsc *ActiveDirectorySettingController) TestOAuthConnection(c *gin.Context) {
	var req dtos.EntraOAuthConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	token, err := adsc.service.TestOAuthConnection(&req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ActiveDirectoryConnectionError:
			utility.SendCustomErrResp(c, http.StatusUnprocessableEntity, e.Error())
			return
		default:
			utility.SendInternalErrResp(adsc.db, c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message":"OAuth connection successful","access_token": token})
}

func (adsc *ActiveDirectorySettingController) GetOAuthConfiguration(c *gin.Context) {
	config, err := adsc.service.GetOAuthConfiguration()
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, config)
}

func (adsc *ActiveDirectorySettingController) SaveOAuthConfiguration(c *gin.Context) {
	var req dtos.EntraOAuthConfigDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}
	err = adsc.service.SaveOAuthConfiguration(currentUserID, req)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(adsc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OAuth configuration saved successfully"})
}


func (adsc *ActiveDirectorySettingController) DeleteOAuthConfiguration(c *gin.Context) {

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}

	err = adsc.service.DeleteOAuthConfiguration(currentUserID)
	if err != nil {
		utility.SendInternalErrResp(adsc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OAuth configuration deleted successfully",
	})
}


func (adsc *ActiveDirectorySettingController) GetAzureUserGroups(c *gin.Context) {
    userID := c.Param("id")
    if userID == "" {
        utility.SendInvalidErrResp(c)
        return
    }

    groups, err := adsc.service.GetAzureUserGroups(userID)
    if err != nil {
        utility.SendInternalErrResp(adsc.db, c, err.Error())
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "User groups fetched successfully",
        "count":   len(groups),
        "data":    groups,
    })
}




