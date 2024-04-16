package controllers

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserServiceInterface
	db      *gorm.DB
}

func NewUSerController(service services.UserServiceInterface) *UserController {
	return &UserController{
		service: service,
	}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	sort := c.Query("sort")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"designation_ids",
		"organization_unit_1_ids",
		"organization_unit_2_ids",
		"organization_unit_3_ids",
		"organization_unit_4_ids",
		"organization_unit_5_ids",
		"user_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := organizationDTOs.UserListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		Filters:     filters,
		QuickFilter: quickFilter,
	}

	users, err := uc.service.GetAll(&queryParams)

	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) CreateUser(c *gin.Context) {

	var userFormDTO organizationDTOs.UserFormDTO

	if err := c.ShouldBindJSON(&userFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	id, err := uc.service.Create(currentUserID, &userFormDTO)

	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusCreated, commonDTOs.IDResponse{ID: id})
}

func (uc *UserController) UpdateUser(c *gin.Context) {

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var userFormDTO organizationDTOs.UserFormDTO
	if err := c.ShouldBindJSON(&userFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	err = uc.service.Update(uint(userID), currentUserID, &userFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) DeleteUsers(c *gin.Context) {

	var request struct {
		IDs []uint `json:"ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	failedIDs, err := uc.service.Delete(request.IDs, currentUserID)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"failed_ids": failedIDs})
}

func (uc *UserController) GetSetting(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	userSetting, err := uc.service.GetUserSetting(uint(userID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, userSetting)
}

func (uc *UserController) UpdateSetting(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	var userSettingFormDTO organizationDTOs.UserSettingFormDTO
	if err := c.ShouldBindJSON(&userSettingFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	err = uc.service.UpdateUserSetting(uint(userID), currentUserID, &userSettingFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) EnableMFA(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = uc.service.EnableMFA(uint(userID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) DisableMFA(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = uc.service.DisableMFA(uint(userID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) EnableNotification(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = uc.service.EnableNotification(uint(userID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (uc *UserController) DisableNotification(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = uc.service.DisableNotification(uint(userID))

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) NewPassword(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = uc.service.ResetPassword(uint(userID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
