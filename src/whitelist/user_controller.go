package whitelist

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/SkycoinPro/skywire-services-whitelist/src/api"
	"github.com/SkycoinPro/skywire-services-whitelist/src/template"
)

// UserController is handling reguests regarding Model
type UserController struct {
	userService UserService
	mailService template.Service
}

func DefaultUserController() UserController {
	return NewUserController(DefaultUserService(), template.DefaultService())
}

func NewUserController(ws UserService, ms template.Service) UserController {
	return UserController{
		userService: ws,
		mailService: ms,
	}
}

func (ctrl UserController) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	closedUserGroup := closed.Group("/users")
	adminGroup := closed.Group("/admin")

	adminGroup.GET("/users", ctrl.canManipulateUsersMiddleware, ctrl.getAllUsers)
	adminGroup.GET("/users/:username", ctrl.canManipulateUsersMiddleware, ctrl.getByUsername)

	closedUserGroup.GET("/keys", ctrl.listKeys)
	closedUserGroup.POST("/keys", ctrl.addKey)
	closedUserGroup.DELETE("/keys", ctrl.removeKey)

	closedUserGroup.PATCH("/address", ctrl.updateAddress)

	adminGroup.GET("/disableUser", ctrl.canManipulateUsersMiddleware, ctrl.disableUserFromSubmittingWhitelist)
	adminGroup.GET("/enableUser", ctrl.canManipulateUsersMiddleware, ctrl.enableUserToSubmitWhitelist)

	closed.GET("/info", ctrl.info)

}

// @Summary List User's API keys
// @Description Return collection of User's generated API keys
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} string
// @Failure 500 {object} api.ErrorResponse
// @Router /users/keys [get]
func (ctrl UserController) listKeys(c *gin.Context) {
	keys, err := ctrl.userService.GetKeys(currentUser(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	} else if len(keys) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// @Summary Generate User's API key
// @Description Method that is going to generate, persist and return User's new API key
// @Tags users
// @Accept  json
// @Produce  json
// @Success 201 {string} string
// @Failure 500 {object} api.ErrorResponse
// @Router /users/keys [post]
func (ctrl UserController) addKey(c *gin.Context) {
	key, err := ctrl.userService.AddKey(currentUser(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, key)
}

// @Summary Remove User's API key
// @Description Match provided API key and remove it if exists
// @Tags users
// @Accept  json
// @Produce  json
// @Param keyToBeRemoved body apiKey true "User's API key to be removed"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/keys [delete]
func (ctrl UserController) removeKey(c *gin.Context) {
	var keyToBeRemoved apiKey
	if err := c.BindJSON(&keyToBeRemoved); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	if err := ctrl.userService.RemoveKey(currentUser(c), keyToBeRemoved.Key); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary List all users
// @Description Method for admins to get list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} whitelist.User
// @Failure 500 {object} api.ErrorResponse
// @Router /admin/users [get]
func (ctrl UserController) getAllUsers(c *gin.Context) {
	users, err := ctrl.userService.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Get specific user
// @Description Method for admins to get specific user by username
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User email"
// @Success 200 {object} whitelist.User
// @Failure 400 {object} api.ErrorResponse
// @Router /admin/users/:username [get]
func (ctrl UserController) getByUsername(c *gin.Context) {
	email := c.Param("username")

	usr, err := ctrl.userService.FindBy(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

// @Summary Disable a user by username
// @Description Method for admins to prevent user from submitting whitelist applications
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User email"
// @Success 200 {object} whitelist.User
// @Failure 400 {object} api.ErrorResponse
// @Router /admin/users/:username [get]
func (ctrl UserController) disableUserFromSubmittingWhitelist(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Username]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectWhitelistIDSent.Error()})
		return
	}
	email := params[Username][0]

	err := ctrl.userService.ChangeUserWhitelistSubmissionPrivilege(email, false)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Enable user to submit whitelist applications
// @Description Method for admins to enable user to submit whitelist
// @Tags users
// @Accept json
// @Produce json
// @Param id query string true "User email"
// @Success 200 {object} whitelist.User
// @Failure 400 {object} api.ErrorResponse
// @Router /admin/users/:username [get]
func (ctrl UserController) enableUserToSubmitWhitelist(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Username]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectWhitelistIDSent.Error()})
		return
	}
	email := params[Username][0]

	err := ctrl.userService.ChangeUserWhitelistSubmissionPrivilege(email, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Retrieve signed in User's info
// @Description Information about currently signed in user is collected and returned as response.
// @Tags authorization
// @Accept  json
// @Produce  json
// @Success 200 {object} whitelist.User
// @Failure 401 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /auth/info [get]
func (ctrl *UserController) info(c *gin.Context) {
	userEmail := currentUser(c)

	if len(userEmail) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	usr, err := ctrl.userService.FindUserInfo(userEmail)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if msg := err.Error(); strings.Contains(msg, errMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, usr)
}

// @Summary Update Users's Skycoin address
// @Description Collect, validate and store User's new Skycoin address.
// @Tags users
// @Accept  json
// @Produce  json
// @Param newAddress body whitelist.AddressUpdateReq true "New User"
// @Success 200 {object} whitelist.User
// @Failure 400 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /users/address [patch]
func (ctrl UserController) updateAddress(c *gin.Context) {
	var newAddress AddressUpdateReq
	if err := c.BindJSON(&newAddress); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	updatedUser, err := ctrl.userService.UpdateAddress(currentUser(c), newAddress.Address)
	if err != nil {
		statusCode := http.StatusBadRequest
		if msg := err.Error(); strings.Contains(msg, errMissingMandatoryFields.Error()) {
			statusCode = http.StatusUnprocessableEntity
		} else if msg := err.Error(); strings.Contains(msg, errUnableToSave.Error()) {
			statusCode = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(statusCode, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, updatedUser)
}

func (ctrl UserController) canManipulateUsersMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !(usr.CanReviewWhitelsit() || usr.CanFlagUserAsVIP()) {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errNoAdminPrivlages.Error()})
		return
	}
}

type apiKey struct {
	Key string
}

type AddressUpdateReq struct {
	Address string
}
