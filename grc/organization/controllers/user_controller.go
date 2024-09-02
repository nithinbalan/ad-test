package controllers

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	"jethur/utility"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserServiceInterface
	db      *gorm.DB
}

func NewUserController(service services.UserServiceInterface) *UserController {
	return &UserController{
		service: service,
	}
}

func (uc *UserController) GetAllUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort_order", "asc")
	sortBy := c.DefaultQuery("sort_by", "name")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")
	userIdsStr := c.Query("user_ids")

	userIds, err := utility.ParseCommaSeparatedStringToIntArray(userIdsStr)
	if err != nil {
		utility.SendInvalidErrResp(c, errors.New("invalid 'user_ids' parameter"))
		return
	}

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

	validSortFields := map[string]string{
		"id":            "users.id",
		"name":          "users.first_name_en",
		"email":         "users.email",
		"mobile_number": "users.mobile_number",
		"designation":   "designations.title_en",
		"role":          "roles_agg.roles_en",
		"user_status":   "user_statuses.title_en",
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

	dbField, valid := validSortFields[sortBy]
	if !valid {
		err := errors.New("invalid 'sort_by' parameter")
		utility.SendInvalidErrResp(c, err)
		return
	}

	queryParams := organizationDTOs.UserListQueryParamsDTO{
		Sort:        sort,
		SortBy:      dbField,
		Keyword:     keyword,
		Status:      status,
		Page:        page,
		PageSize:    pageSize,
		UserIDs:     userIds,
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

func (uc *UserController) GetUserDetails(c *gin.Context) {
	idStr := c.Param("userId")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	user, err := uc.service.GetFullDetailsByID(uint(id))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, user)
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
	userID, err := strconv.Atoi(c.Param("userId"))
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

func (uc *UserController) GetUserLoginHistory(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	loginAtFrom := c.Query("login_at_from")
	loginAtTo := c.Query("login_at_to")

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

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	queryParams := organizationDTOs.UserLoginHistoryListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		Page:        page,
		PageSize:    pageSize,
		LoginAtFrom: loginAtFrom,
		LoginAtTo:   loginAtTo,
	}
	userLoginHistory, err := uc.service.GetUserLoginHistory(uint(userID), &queryParams)

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
	c.JSON(http.StatusOK, userLoginHistory)
}

func (uc *UserController) GetCurrentUserDetails(c *gin.Context) {

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	user, err := uc.service.GetCurrentUserDetails(uint(currentUserID))
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(uc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserNotifications(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	quickFilter := c.QueryArray("quick_filter")

	moduleIDsStr := c.Query("module_ids")
	var moduleIDs []uint
	if moduleIDsStr != "" {
		for _, idStr := range strings.Split(moduleIDsStr, ",") {
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				utility.SendInternalErrResp(uc.db, c, "Invalid module_id parameter")
				return
			}
			moduleIDs = append(moduleIDs, uint(id))
		}
	}

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

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	// Prepare query parameters
	queryParams := organizationDTOs.UserNotificationListQueryParamsDTO{
		Sort:        sort,
		SortBy:      sortBy,
		Keyword:     keyword,
		ModuleIDs:   moduleIDs,
		Page:        page,
		PageSize:    pageSize,
		QuickFilter: quickFilter,
	}

	userNotifications, err := uc.service.GetUserNotifications(uint(currentUserID), &queryParams)
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

	c.JSON(http.StatusOK, userNotifications)
}

func (uc *UserController) GetUserNewNotificationCount(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	count, err := uc.service.GetUserNewNotificationCount(currentUserID)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (uc *UserController) MarkwNotificationAsRead(c *gin.Context) {
	var formDTO organizationDTOs.MarkNotificationAsReadFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	_, err = uc.service.MarkNotificationsAsRead(currentUserID, &formDTO)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (uc *UserController) MarkwNotificationAsUnread(c *gin.Context) {
	var formDTO organizationDTOs.MarkNotificationAsReadFormDTO

	if err := c.ShouldBindJSON(&formDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(uc.db, c, err.Error())
		return
	}

	_, err = uc.service.MarkNotificationsAsUnread(currentUserID, &formDTO)
	if err != nil {
		utility.SendInternalErrResp(nil, c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
