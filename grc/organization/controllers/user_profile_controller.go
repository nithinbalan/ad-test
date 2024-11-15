package controllers

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/services"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/utility"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProfileController struct {
	service     services.UserSettingServiceInterface
	taskService services.UserTaskServiceInterface
	userService services.UserServiceInterface
	db          *gorm.DB
}

func NewUserProfileController(service services.UserSettingServiceInterface, taskService services.UserTaskServiceInterface, userService services.UserServiceInterface) *UserProfileController {
	return &UserProfileController{
		service:     service,
		taskService: taskService,
		userService: userService,
	}
}

func (upc *UserProfileController) GetProfileSetting(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	userProfile, err := upc.service.GetUserSetting(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, userProfile)
}

func (upc *UserProfileController) UpdateProfileSetting(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	var userProfileFormDTO organizationDTOs.UserSettingFormDTO
	if err := c.ShouldBindJSON(&userProfileFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err = upc.service.UpdateUserSetting(currentUserID, currentUserID, &userProfileFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) EnableMainNotification(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	err = upc.service.EnableMainNotification(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}
	c.Status(http.StatusOK)
}

func (upc *UserProfileController) DisableMainNotification(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	err = upc.service.DisableMainNotification(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) EnableMFA(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	err = upc.service.EnableMFA(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) DisableMFA(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	err = upc.service.DisableMFA(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) NewPassword(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	err = upc.service.ResetPassword(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) GetNotificationSettings(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	notifications, err := upc.service.GetNotificationSettings(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, notifications)
}

func (upc *UserProfileController) EnableNotification(c *gin.Context) {
	upc.toggleNotification(c, true)
}

func (upc *UserProfileController) DisableNotification(c *gin.Context) {
	upc.toggleNotification(c, false)
}

func (upc *UserProfileController) toggleNotification(c *gin.Context, enable bool) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	notificationMasterID, err := strconv.Atoi(c.Param("notificationMasterId"))
	if err != nil {
		utility.SendNotFoundErrResp(c)
		return
	}

	err = upc.service.UpdateNotificationSetting(currentUserID, uint(notificationMasterID), enable)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (upc *UserProfileController) GetAllUserTasks(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	status := c.Query("status")
	quickFilter := c.QueryArray("quick_filter")

	// Define an array of allowed parameter names
	allowedParams := []string{
		"compliance_action_plan_ids",
		"sla_action_plan_ids",
		"incident_action_plan_ids",
		"control_action_plan_ids",
		"severity_level_ids",
		"module_ids",
		"task_status_ids",
	}

	filters := utility.ParseQueryFilters(allowedParams, c.Request.URL.Query())
	// Filter for responsible_user_id=currentUserID
	filters["responsible_user_id"] = []string{strconv.Itoa(int(currentUserID))}

	// Parse target_date_from, target_date_to
	targetDateFrom := c.Query("target_date_from")
	targetDateTo := c.Query("target_date_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, "Invalid page_size parameter")
		return
	}

	queryParams := taskDTOs.TaskListQueryParamsDTO{
		Sort:           sort,
		SortBy:         sortBy,
		Keyword:        keyword,
		Status:         status,
		Page:           page,
		PageSize:       pageSize,
		Filters:        filters,
		QuickFilter:    quickFilter,
		TargetDateFrom: targetDateFrom,
		TargetDateTo:   targetDateTo,
	}
	userTasks, err := upc.taskService.GetAll(currentUserID, &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, userTasks)
}

func (upc *UserProfileController) GetUserTaskSummary(c *gin.Context) {
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	userTaskSummary, err := upc.taskService.GetUserTaskSummary(currentUserID)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, userTaskSummary)
}

func (upc *UserProfileController) GetUserLoginHistory(c *gin.Context) {
	sort := c.Query("sort_order")
	sortBy := c.Query("sort_by")
	keyword := c.Query("keyword")
	loginAtFrom := c.Query("login_at_from")
	loginAtTo := c.Query("login_at_to")

	page, err := utility.ParseQueryParamToInt(c, "page", 1)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, "Invalid page parameter")
		return
	}

	pageSize, err := utility.ParseQueryParamToInt(c, "page_size", 10)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, "Invalid page_size parameter")
		return
	}

	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
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
	userLoginHistory, err := upc.userService.GetUserLoginHistory(uint(currentUserID), &queryParams)

	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}
	c.JSON(http.StatusOK, userLoginHistory)
}

func (upc *UserProfileController) ChangePassword(c *gin.Context) {
	var changePasswordDTO organizationDTOs.ChangePasswordDTO

	// Bind the incoming JSON to the DTO
	if err := c.ShouldBindJSON(&changePasswordDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Get the current user's ID
	currentUserID, err := utility.GetCurrentUserID(c)
	if err != nil {
		utility.SendInternalErrResp(upc.db, c, err.Error())
		return
	}

	// Validate and process the password change
	err = upc.service.ChangePassword(currentUserID, changePasswordDTO.OldPassword, changePasswordDTO.NewPassword)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		default:
			utility.SendInternalErrResp(upc.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
