package auth

import (
	"errors"
	"net/http"

	"github.com/SkycoinPro/skywire-services-whitelist/src/api"
	"github.com/SkycoinPro/skywire-services-whitelist/src/whitelist"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ErrUnauthorized is the error returned when user can't be found in JWT
var ErrUnauthorized = errors.New("auth controller: user is not recognized")

// Controller is handling reguests regarding Model
type Controller struct {
	whitelistService whitelist.Service
}

func DefaultController() Controller {
	return NewController(whitelist.DefaultService())
}

func NewController(us whitelist.Service) Controller {
	return Controller{
		whitelistService: us,
	}
}

func (ctrl Controller) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	authorization := ctrl.initJWT()
	closed.Use(authorization.MiddlewareFunc())
}

func (ctrl *Controller) initJWT() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      viper.GetString("jwt.realm"),
		Key:        []byte(viper.GetString("jwt.key")),
		Timeout:    viper.GetDuration("jwt.timeout"),
		MaxRefresh: viper.GetDuration("jwt.max-refresh"),
		Authenticator: func(email string, password string, c *gin.Context) (string, bool) {
			return email, false // this should not be invoked
		},
		Authorizator: func(userID string, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			if unconfirmed := claims["missing_confirmation"]; unconfirmed != nil && unconfirmed.(bool) {
				return false // don't allow access to not confirmed users
			}
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			if c.Request.URL.Path == "/api/v1/auth/login" {
				code = http.StatusBadRequest
			}
			c.JSON(code, api.ErrorResponse{Error: "auth service: Username and/or Password do not match any user"})
		},
		PayloadFunc: func(userID string) map[string]interface{} {
			return make(map[string]interface{}) // no claims should be defined here
		},
		//TODO (security) consider customizign some of parameters
	}
}
