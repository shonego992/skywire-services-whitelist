package whitelist

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/SkycoinPro/skywire-services-whitelist/src/api"
	"github.com/SkycoinPro/skywire-services-whitelist/src/template"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

const Token = "token"
const Email = "email"
const Id = "id"
const Username = "username"

// Controller is handling reguests regarding Model
type Controller struct {
	whitelistService Service
	userService      UserService
	mailService      template.Service
	minerService     MinerService
}

func DefaultController() Controller {
	return NewController(DefaultService(), template.DefaultService(), DefaultUserService(), DefaultMinerService())
}

func NewController(ws Service, ms template.Service, us UserService, mis MinerService) Controller {
	return Controller{
		whitelistService: ws,
		userService:      us,
		mailService:      ms,
		minerService:     mis,
	}
}

func (ctrl Controller) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	publicWhitelistGroup := public.Group("/whitelist")
	whitelistGroup := closed.Group("/whitelist")

	publicWhitelistGroup.POST("/linkNodes", ctrl.linkNodes)

	whitelistGroup.POST("/application", ctrl.canSubmitWhitelists, ctrl.createApplication)
	whitelistGroup.GET("/application", ctrl.getApplication)
	whitelistGroup.POST("/updateApplication", ctrl.canSubmitWhitelists, ctrl.updateApplication)

	whitelistGroup.GET("/whitelists", ctrl.canReviewWhitelistsMiddleware, ctrl.getWhitelists)
	whitelistGroup.GET("/whitelist", ctrl.canReviewWhitelistsMiddleware, ctrl.getWhitelist)
	whitelistGroup.POST("/whitelist", ctrl.canReviewWhitelistsMiddleware, ctrl.changeApplicationStatus)
}

func (ctrl Controller) canReviewWhitelistsMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !usr.CanReviewWhitelsit() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errNoAdminPrivlages.Error()})
		return
	}
}

func (ctrl Controller) isAdminMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !usr.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errNoAdminPrivlages.Error()})
		return
	}
}

func (ctrl Controller) canSubmitWhitelists(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || usr.BlockedFromSubmitingWhitelist() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errUserCannotSubmitNewApplications.Error()})
		return
	}
}

// @Summary Lists whitelisted applications
// @Description  Returns an array of the whitelisted applications
// @Tags application
// @Accept json
// @Produce json
// @Success 200 {array} whitelist.Application
// @Failure 500 {object} api.ErrorResponse
// @Router /whitelist/whitelists [get]
func (ctrl Controller) getWhitelists(c *gin.Context) {
	list, err := ctrl.whitelistService.GetAllWhitelists()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// @Summary Link nodes
// @Description  Finds User by provided API key and create a new Miner with provided Nodes
// @Tags application
// @Accept json
// @Produce json
// @Param linkNodesReq body whitelist.linkNodesReq true "Nodes to be linked"
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Failure 422 {object} api.ErrorResponse
// @Router /whitelist/linkNodes [post]
func (ctrl Controller) linkNodes(c *gin.Context) {
	var linkNodesReq linkNodesReq
	if err := c.BindJSON(&linkNodesReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	username, err := ctrl.userService.FindUserByApiKey(linkNodesReq.Key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	err = ctrl.minerService.addMinerToUser(username, linkNodesReq.NodeKeys)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
}

// @Summary Gets whitelisted application
// @Description  Returns an application for given application id
// @Tags application
// @Accept json
// @Produce json
// @Param id query string true "Whitelist application id"
// @Success 200 {object} whitelist.Application
// @Failure 500 {object} api.ErrorResponse
// @Failure 400 {object} api.ErrorResponse
// @Router /whitelist/whitelist [get]
func (ctrl Controller) getWhitelist(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Id]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectWhitelistIDSent.Error()})
		return
	}
	id := params[Id][0]
	application, err := ctrl.whitelistService.GetWhitelist(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, application)
}

// @Summary Change application status
// @Description Change application status according to application change status request
// @Tags application
// @Accept json
// @Produce json
// @Param changeApplicationStatusReq body ChangeApplicationStatus true "Application change request"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /whitelist/whitelist [post]

func (ctrl Controller) changeApplicationStatus(c *gin.Context) {
	var changeApplicationStatusReq ChangeApplicationStatus
	if err := c.BindJSON(&changeApplicationStatusReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	app, err := ctrl.whitelistService.GetWhitelist(fmt.Sprint(changeApplicationStatusReq.ApplicationId))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	if app.CurrentStatus != PENDING {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errNoApplicationInProgressForUser.Error()})
		return
	}

	whitelistingApplicant, err := ctrl.whitelistService.UpdateApplicationStatus(changeApplicationStatusReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	if changeApplicationStatusReq.Status == APPROVED {
		miner := ctrl.whitelistService.GetMinerForApplication(app)
		latestChangeHistory := app.GetLatestChangeHistory()
		if miner.ID == 0 {
			//there's no existing miner, create new one
			if _, err := ctrl.minerService.CreateMinerEntitiesForUser(&latestChangeHistory, app); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
				return
			}

		} else {
			// append nodes to existing miner
			nodes := latestChangeHistory.Nodes
			for _, node := range miner.Nodes {
				//do not delete existing miner nodes when appending new ones
				newNode := Node{Key: node.Key}
				if node.ID > 0 {
					newNode.ID = node.ID
					newNode.CreatedAt = node.CreatedAt
				}
				nodes = append(nodes, &newNode)
			}
			miner.ApprovedNodesCount += len(latestChangeHistory.Nodes)

			if err := ctrl.minerService.UpdateMiner(&miner, nodes); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
				return
			}
		}
		//add images to created/existing miner
		if appMiner := ctrl.whitelistService.GetMinerForApplication(app); appMiner.ID > 0 {
			err := ctrl.minerService.AddImagesToMiner(appMiner.ID, latestChangeHistory.Images)
			if err != nil {
				log.Errorf("Miner saved but images are not added to miner profile")
			}

		}
	}

	if err := ctrl.mailService.WhitelistApplicationUpdated(whitelistingApplicant, statusMap[changeApplicationStatusReq.Status], changeApplicationStatusReq.UserComment); err != nil {
		log.Error("Was not able to send notification about whitelist application status update to user ", whitelistingApplicant)
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errAppReviewedEmailFailed.Error()})
		return
	}
}

// @Summary Create a new application in system for current user
// @Description Collect provided Application attributes from the body and create new Application in the system
// containing an initial change history record for that application
// @Tags application
// @Accept  json
// @Produce  json
// @Param newUser body Model true "New User"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /whitelist/application [post]
func (ctrl Controller) createApplication(c *gin.Context) {
	whitelistingApplicant := currentUser(c)
	appReq, err := extractDataFromApplicationRequest(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if len(appReq.Files)+len(appReq.OldImages) < 3 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errNotEnoughImages.Error()})
		return
	}
	imageIDs, err := ctrl.whitelistService.CheckForDuplicates(&appReq, true, currentUser(c))
	if len(imageIDs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errDuplicateDatabaseImages.Error(), imageIDs)})
		return
	}

	_, appError := ctrl.whitelistService.CreateApplication(appReq, whitelistingApplicant)
	if appError.Error != nil {
		if appError.Error == errWrongNodeKeys {
			if len(appError.AlreadyTakenKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errAlreadyTakenKeys.Error(), appError.AlreadyTakenKeys)})
				return
			} else if len(appError.WrongKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errWrongNodeKeys.Error(), appError.WrongKeys)})
				return
			} else if len(appError.DuplicateKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errDuplicateKeys.Error(), appError.DuplicateKeys)})
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: appError.Error.Error()})
		return
	}

	if err := ctrl.mailService.WhitelistApplicationCreated(whitelistingApplicant); err != nil {
		if len(appError.FailedImages) != 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errAppCreatedImagesAndEmailFailed.Error()})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errAppCreatedEmailFailed.Error()})
			return
		}
	}
	if len(appError.FailedImages) != 0 {
		c.JSON(http.StatusOK, "Application sucessfully created but some images were not saved")
	} else {
		c.JSON(http.StatusOK, "Application sucessfully created")
	}
}

// @Summary Update an existing application in system for current user
// @Description Collect provided Application attributes from the body and create new change history record for
// users current application
// @Tags application
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /whitelist/updateApplication [post]
func (ctrl Controller) updateApplication(c *gin.Context) {
	appReq, err := extractDataFromApplicationRequest(c.Request)
	if err != nil {
		log.Info("Unable to parse request due to error ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if len(appReq.Files)+len(appReq.OldImages) < 3 {
		log.Info("Update request made with not enough images")
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errNotEnoughImages.Error()})
		return
	}
	imageIDs, err := ctrl.whitelistService.CheckForDuplicates(&appReq, false, currentUser(c))
	if err != nil {
		log.Info("Unable to check for duplicate images due to error ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if len(imageIDs) > 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errDuplicateDatabaseImages.Error(), imageIDs)})
		return
	}

	appError := ctrl.whitelistService.UpdateApplication(appReq, currentUser(c))
	if appError.Error != nil {
		if appError.Error == errWrongNodeKeys {
			if len(appError.AlreadyTakenKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errAlreadyTakenKeys.Error(), appError.AlreadyTakenKeys)})
			} else if len(appError.WrongKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errWrongNodeKeys.Error(), appError.WrongKeys)})
			} else if len(appError.DuplicateKeys) != 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: fmt.Sprintf("%s - %v", errDuplicateKeys.Error(), appError.DuplicateKeys)})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: "whitelist controller: unexpected state with wrong node keys"})
			}
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: appError.Error.Error()})
		return
	}

	if len(appError.FailedImages) != 0 {
		c.JSON(http.StatusOK, "Application sucessfully updated but some images were not saved")
	} else {
		c.JSON(http.StatusOK, "Application sucessfully updated")
	}
}

// @Summary Gets the application for curent user
// @Description Gets the application that is currently in progress for current user
// @Tags application
// @Accept json
// @Produce json
// @Success 200 {object} whitelist.Application
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /whitelist/application [get]
func (ctrl Controller) getApplication(c *gin.Context) {
	application, err := ctrl.whitelistService.getActiveApplication(currentUser(c))
	if err != nil {
		if err == ErrCannotFindUser{
			c.AbortWithStatusJSON(http.StatusNotFound, api.ErrorResponse{Error: err.Error()})
			return
		}
		if err == errNoApplicationInProgressForUser {
			c.JSON(http.StatusOK, Application{})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, application)
}

// TODO: check validation of request fields
func extractDataFromApplicationRequest(req *http.Request) (ApplicationReq, error) {
	var (
		//status int
		err error
	)
	// parse request
	const _24K = (1 << 20) * 24
	if err = req.ParseMultipartForm(_24K); nil != err {
		//status = http.StatusInternalServerError
		return ApplicationReq{}, err
	}
	form := req.MultipartForm
	var keys []Node
	var images []Image
	json.Unmarshal([]byte(form.Value["nodes"][0]), &keys)
	json.Unmarshal([]byte(form.Value["oldImages"][0]), &images)
	ApplicationReq := ApplicationReq{
		Location:    form.Value["location"][0],
		Description: form.Value["description"][0],
		Nodes:       keys,
		OldImages:   images,
		Files:       extractBytesFromFile(form.File),
	}
	return ApplicationReq, nil
}

func extractBytesFromFile(files map[string][]*multipart.FileHeader) (response []ApplicationImage) {
	for _, fheaders := range files {
		for _, hdr := range fheaders {
			if infile, err := hdr.Open(); nil != err {
				log.Error("Unable to read file from the request ", err)
			} else {
				buffer := make([]byte, hdr.Size)
				infile.Read(buffer)
				response = append(response, ApplicationImage{Name: hdr.Filename, File: buffer})
				infile.Close()
			}
		}
	}
	return
}

// returns username(email address) of current user
func currentUser(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(string)
}

type ChangeApplicationStatus struct {
	ApplicationId uint
	Status        uint8
	UserComment   string
	AdminComment  string
}

type ApplicationReq struct {
	Description string
	Location    string
	Nodes       []Node
	OldImages   []Image
	Files       []ApplicationImage
}

type ApplicationImage struct {
	Name   string
	File   []byte
	Path   string
	Hashed string
}

type linkNodesReq struct {
	Key      string
	NodeKeys []string
}
