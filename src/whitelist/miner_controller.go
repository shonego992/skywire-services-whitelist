package whitelist

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-whitelist/src/api"
	"github.com/SkycoinPro/skywire-services-whitelist/src/template"
)

// MinerController is handling reguests regarding Model
type MinerController struct {
	minerService     MinerService
	userService      UserService
	whitelistService Service
	mailService      template.Service
}

func DefaultMinerController() MinerController {
	return NewMinerController(DefaultMinerService(), DefaultUserService(), DefaultService(), template.DefaultService())
}

func NewMinerController(ms MinerService, us UserService, ws Service, ts template.Service) MinerController {
	return MinerController{
		minerService:     ms,
		userService:      us,
		whitelistService: ws,
		mailService:      ts,
	}
}

func (ctrl MinerController) RegisterAPIs(public *gin.RouterGroup, closed *gin.RouterGroup) {
	whitelistGroup := closed.Group("/miners")

	whitelistGroup.GET("/miners", ctrl.getUserMiners)
	whitelistGroup.GET("/miner", ctrl.getSpecificMiner)
	whitelistGroup.POST("/miner", ctrl.isNotBlacklisted, ctrl.updateMiner)
	whitelistGroup.POST("/transferMiner", ctrl.transferMiner)

	whitelistGroup.GET("/minersForUser", ctrl.isAdminMiddleware, ctrl.getMinersForUser)
	whitelistGroup.GET("/allMiners", ctrl.isAdminMiddleware, ctrl.getAllMiners)
	whitelistGroup.GET("/minerForAdmin", ctrl.isAdminMiddleware, ctrl.getSpecificMinerForAdmin)
	whitelistGroup.DELETE("/miner/:id", ctrl.isAdminMiddleware, ctrl.deleteMinerByID)
	whitelistGroup.GET("miner/:id/activate", ctrl.isAdminMiddleware, ctrl.activateMinerByID)

	// import
	whitelistGroup.GET("/import", ctrl.getImportData)
	whitelistGroup.POST("/import", ctrl.updateImportData)
	whitelistGroup.POST("/import/process", ctrl.processImportData)
	whitelistGroup.POST("/uploadUserList", ctrl.uploadUserList)

	whitelistGroup.POST("/exportMiners", ctrl.isAdminMiddleware, ctrl.exportMiners)
	whitelistGroup.POST("/exportMinersNoLimitations", ctrl.isAdminMiddleware, ctrl.exportMinersNoLimitations)
}

func (ctrl MinerController) isNotBlacklisted(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || usr.BlockedFromSubmitingWhitelist() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errUserCannotSubmitNewApplications.Error()})
		return
	}
}

func (ctrl MinerController) transferMiner(c *gin.Context) {
	var transferReq transferMinerReq
	if err := c.BindJSON(&transferReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	_, err := ctrl.userService.FindBy(transferReq.TransferTo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errCannotFindTransferUser.Error()})
		return
	}

	err = ctrl.minerService.transferMiner(transferReq, currentUser(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	ctrl.notifyUserAfterMinerTransfer(transferReq.TransferTo, currentUser(c))

}

func (ctrl MinerController) exportMinersNoLimitations(c *gin.Context) {
	var importedData exportMinersReq
	if err := c.BindJSON(&importedData); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	startDate := convertDate(importedData.StartDate)
	endDate := convertDate(importedData.EndDate)
	var startDateUnix, endDateUnix int64
	if startDate.After(endDate) {
		startDateUnix = 0
		endDateUnix = 0
	} else {
		startDateUnix = startDate.Unix()
		endDateUnix = endDate.Unix()
	}

	allExports, err := ctrl.minerService.GetAllUptimeRecords(startDateUnix, endDateUnix)
	if err != nil {
		log.Debug("Error while fetching data for export")
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	var exportList []minerExport

	for _, export := range allExports {
		miner, _ := ctrl.minerService.GetUserByNodeKey(export.Key)
		address, _ := ctrl.userService.FindPayoutAddress(miner.Username)
		exportList = append(exportList, minerExport{
			Address:    address,
			Type:       mapType(miner.Type),
			Mail:       miner.Username,
			Key:        export.Key,
			BatchLabel: miner.BatchLabel,
			Uptime:     fmt.Sprintf("%f", export.Percentage),
		})
	}
	file, err := os.Create("resultsAll.csv")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	var headline = []string{"mail", "address", "key", "type", "batch_label", "uptime"}
	err = writer.Write(headline)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	for _, value := range exportList {
		var stringContent = []string{value.Mail, value.Address, value.Key, value.Type, value.BatchLabel, value.Uptime}
		err = writer.Write(stringContent)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}
	}
	writer.Flush()

	c.File("resultsAll.csv")
}

func (ctrl MinerController) exportMiners(c *gin.Context) {
	var importedData exportMinersReq
	if err := c.BindJSON(&importedData); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	startDate := convertDate(importedData.StartDate)
	endDate := convertDate(importedData.EndDate)
	//TODO: Check this call for date conversion and how to
	//miners, err := ctrl.minerService.exportMiners(importedData, startDate, endDate)
	users, err := ctrl.userService.GetUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: "miner controller: Cannot access database to load users"})
		return
	}

	// this buffered channel will block at the concurrency limit
	usersChan := make(chan struct{})
	// this channel will not block and collect the http request results
	resultsChan := make(chan *[]minerExport)
	// make sure we close these channels when we're done with them
	defer func() {
		close(resultsChan)
		close(usersChan)
	}()
	usrCnt := len(users)
	pUsrCnt := 0
	done := false
	var exportList []minerExport
	threads := viper.GetInt("export.threads")
	semaphoreChan := make(chan struct{}, threads)
	defer func() {
		close(semaphoreChan)
	}()
	now := time.Now()
	go ctrl.makeRequests(&semaphoreChan, &usersChan, &resultsChan, users, startDate, endDate)
	// start listening for any results over the resultsChan
	// once we get a result append it to the result slice
	for {
		select {
		case exports := <-resultsChan:
			numberOfOfficial := 0
			numberOfDiy := 0
			exportList = append(exportList, *exports...)
			address := ""
			mail := ""
			for _, export := range *exports {
				if len(address) == 0 {
					address = export.Address
					mail = export.Mail
				}
				if export.Type == "OFFICIAL" {
					numberOfOfficial++
				} else {
					numberOfDiy++
				}
			}
			if numberOfDiy > 0 {
				diyRecord := ExportRecord{
					Username:      mail,
					MinerType:     DIY,
					PayoutAddress: address,
					TimeOfExport:  time.Now(),
					NumberOfNodes: numberOfDiy,
				}
				if err = ctrl.userService.CreateExportRecord(diyRecord); err != nil {
					log.Error("Error creating export record for user", mail)
					c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
					return
				}
			}
			if numberOfOfficial > 0 {
				officialRecord := ExportRecord{
					Username:      mail,
					MinerType:     OFFICIAL,
					PayoutAddress: address,
					TimeOfExport:  time.Now(),
					NumberOfNodes: numberOfOfficial,
				}
				if err = ctrl.userService.CreateExportRecord(officialRecord); err != nil {
					log.Error("Error creating export record for user", mail)
					c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
					return
				}
			}

		case <-usersChan:
			pUsrCnt++
			if usrCnt == pUsrCnt {
				timer := time.NewTimer(20 * time.Second) //TODO improve this
				log.Info("Timer start")
				go func() {
					<-timer.C
					log.Info("timer ticked 20 secs")
					done = true
					usersChan <- struct{}{}
				}()
			}
		}

		log.Debugf("Checking is it done with %v total and %v processed", usrCnt, pUsrCnt)
		if done {
			log.Info("done")
			break
		}
	}
	log.Info("done in ", time.Since(now))

	file, err := os.Create("result.csv")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	var headline = []string{"mail", "address", "key", "type", "batch_label", "uptime"}
	err = writer.Write(headline)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	for _, value := range exportList {
		var stringContent = []string{value.Mail, value.Address, value.Key, value.Type, value.BatchLabel, value.Uptime}
		err = writer.Write(stringContent)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}
	}
	writer.Flush()

	c.File("result.csv")
}

func (ctrl MinerController) makeRequests(semaphoreChan *chan struct{}, usersChan *chan struct{}, resultsChan *chan *[]minerExport, users []User,
	startDate time.Time, endDate time.Time) {

	count := 0
	for _, user := range users {
		count++
		if count%20 == 0 {
			log.Infof("Processed %v users for export at %v", count, time.Now())
		}

		if len(user.Miners) == 0 {
			*usersChan <- struct{}{}
			continue
		}

		*semaphoreChan <- struct{}{}
		go ctrl.processExportForUser(semaphoreChan, usersChan, resultsChan, user, startDate, endDate)

		*usersChan <- struct{}{}
	}
}

func (ctrl MinerController) processExportForUser(semaphoreChan *chan struct{}, usersChan *chan struct{}, resultsChan *chan *[]minerExport,
	user User, startDate time.Time, endDate time.Time) {

	exportThreshold := viper.GetFloat64("export.reward-percentage") * 100
	var startDateUnix, endDateUnix int64
	if startDate.After(endDate) {
		startDateUnix = 0
		endDateUnix = 0
	} else {
		startDateUnix = startDate.Unix()
		endDateUnix = endDate.Unix()
	}
	var resp []minerExport
	address, err := ctrl.userService.FindPayoutAddress(user.Username)
	if err != nil {
		log.Debug("Can't find payout address for user ", user.Username)
		// c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		// return
	}

	isBlacklisted := user.Status == 2
	var userDiyNodes []Node
	strt := time.Now()
	allUserMiners, err := ctrl.minerService.getUserMiners(user.Username)
	if err != nil {
		//TODO: check what to do in this case
		log.Info("Not able to find miners for user ", user.Username)
		*usersChan <- struct{}{}
		*semaphoreChan <- struct{}{}
		return
	}
	log.Info("Miner fetch took ", time.Since(strt))
	for _, miner := range allUserMiners {
		if isBlacklisted && miner.Type == DIY {
			continue
		}

		if miner.Type == OFFICIAL {
			if len(miner.Nodes) > 0 {
				uptimes, err := ctrl.minerService.GetUptime(miner.Nodes, true, startDateUnix, endDateUnix)
				if err != nil && err != errNoKeysToGetDataFor {
					//TODO: check what to do in this case
					log.Errorf("Cannot fetch uptimes for nodes %v", miner.Nodes)
					// c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
					return

				}

				for _, node := range miner.Nodes {
					uptimeValue := ""
					for _, uptime := range uptimes {
						if uptime.Key == node.Key {
							uptimeValue = fmt.Sprintf("%f", uptime.Percentage)
							break
						}
					}
					resp = append(resp, minerExport{
						Address:    address,
						Type:       mapType(miner.Type),
						Mail:       miner.Username,
						Key:        node.Key,
						BatchLabel: miner.BatchLabel,
						Uptime:     uptimeValue,
					})
				}
			}
		} else {
			userDiyNodes = append(userDiyNodes, miner.Nodes...)
		}
	}

	if len(userDiyNodes) > 0 {
		uptimes, err := ctrl.minerService.GetUptime(userDiyNodes, true, startDateUnix, endDateUnix)
		if err != nil && err != errNoKeysToGetDataFor {
			//TODO: check what to do in this case
			log.Errorf("Cannot fetch uptimes for nodes %v", userDiyNodes)
			// c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			// return
		}
		sort.Slice(uptimes, func(i, j int) bool {
			return uptimes[i].Percentage > uptimes[j].Percentage
		})

		for i := 0; i < len(uptimes); i++ {
			if uptimes[i].Percentage < exportThreshold || i == 8 {
				break
			}
			uptime := fmt.Sprintf("%f", uptimes[i].Percentage)
			resp = append(resp, minerExport{
				Address:    address,
				Type:       "DIY",
				Mail:       user.Username,
				Key:        uptimes[i].Key,
				BatchLabel: "",
				Uptime:     uptime,
			})
		}
	}
	if len(resp) > 0 {
		*resultsChan <- &resp
	}
	<-*semaphoreChan
}

// @Summary Removes miner
// @Description Removes miner for given id
// @Tags miner
// @Accept json
// @Produce json
// @Param id query string true "ID of miner to be removed"
// @Success 200
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miner/:id [delete]
func (ctrl MinerController) deleteMinerByID(c *gin.Context) {
	id := c.Param("id")
	minerToBeDeleted, err1 := ctrl.minerService.getSpecificMinerForAdmin(id)
	username := minerToBeDeleted.Username
	var currTime = time.Now()
	for _, node := range minerToBeDeleted.Nodes {
		if node.DeletedAt == nil {
			if err := ctrl.whitelistService.RemoveNode(&node, currTime); err != nil {
				log.Errorf("Unable to delete node with key %v due to error: %v", node.Key, err)
			}
		}
	}

	if err := ctrl.minerService.RemoveMiner(id, currTime); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if err1 == nil {
		ctrl.notifyUserAfterMinerDeletion(username)
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary Reenables miner
// @Description Reenables miner for given id
// @Tags miner
// @Accept json
// @Produce json
// @Param id query string true "ID of miner to be reenabled"
// @Sucess 200
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miner/:id/activate [get]
func (ctrl MinerController) activateMinerByID(c *gin.Context) {
	id := c.Param("id")
	minerWithDeletedAt, err := ctrl.minerService.getSpecificMinerForAdmin(id)
	if err != nil {
		if err == errCannotFindMiner {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: err.Error()})
			return

		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return

	}
	var deletionTime = minerWithDeletedAt.DeletedAt
	if err := ctrl.minerService.ActivateMiner(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	activatedMiner, err := ctrl.minerService.getSpecificDisabledMinerForAdmin(id, deletionTime)

	err = ctrl.minerService.ReactivateNodes(activatedMiner.Nodes)
	if err != nil {
		log.Error("Could not reenable nodes due to error ", err)
	}

	if err == nil {
		ctrl.notifyUserAfterMinerReenabling(activatedMiner.Username)
	}
	c.Writer.WriteHeader(http.StatusOK)
}

// @Summary List user's miners
// @Description Returns a list of miners under current user
// @Tags miners
// @Accept json
// @Produce json
// @Success 200 {array} whitelist.Miner
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miners [get]
func (ctrl MinerController) getUserMiners(c *gin.Context) {
	miners, err := ctrl.minerService.getUserMiners(currentUser(c))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, miners)
}

// @Summary List user's miners
// @Description Returns a list of miners under a user with given username
// @Tags miners
// @Accept json
// @Produce json
// @Param Username query string true "User's username"
// @Success 200 {array} whitelist.Miner
// @Failure 400 {object} api.ErrorResponse
// @Router /miners/miner [get]
func (ctrl MinerController) getMinersForUser(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Username]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectUsernameSent.Error()})
		return
	}
	username := params[Username][0]

	miners, err := ctrl.minerService.GetMinersForUser(username)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, miners)
}

func (ctrl MinerController) persistMinerDataFromRequest(c *gin.Context) ([]MinerImport, error) {
	var importedData updateImportDataReq
	if err := c.BindJSON(&importedData); err != nil {
		return nil, err
	}

	return ctrl.persistMinerData(importedData.Data)
}

func (ctrl MinerController) persistMinerData(data []MinerImport) ([]MinerImport, error) {
	return ctrl.minerService.SaveImportData(data)
}

// @Summary Uploads user list
// @Description Exports the user list to csv file and returns number of exported users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {integer} int
// @Failure 400 {string} string
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/uploadUserList [post]
func (ctrl MinerController) uploadUserList(c *gin.Context) {
	file, err := c.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	isTitle := true
	csvFile, _ := file.Open()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var users []MinerImport
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Error(error)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
			return
		}
		if isTitle {
			isTitle = false
			continue
		}
		noOfMiners, err := strconv.Atoi(strings.TrimSpace(line[2]))
		if err != nil {
			noOfMiners = 0
		}

		//removing spaces and brackets in entire record

		users = append(users, MinerImport{
			Username:       strings.TrimSpace(line[0]),
			NumberOfMiners: noOfMiners,
		})
	}

	if len(users) > 0 {
		if _, err := ctrl.persistMinerData(users); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}
	} else {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errNoMinersFoundForImport.Error()})
		return
	}

	c.JSON(http.StatusOK, len(users))
}

// @Summary List all miners
// @Description Method for admins to get list of all miners
// @Tags miners
// @Accept json
// @Produce json
// @Success 200 {array} whitelist.Miner
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/allMiners [get]
func (ctrl MinerController) getAllMiners(c *gin.Context) {
	miners, err := ctrl.minerService.GetAllMiners()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, miners)
}

// @Summary Get specific miner
// @Description Returns miner for given miner id
// @Tags miners
// @Accept json
// @Produce json
// @Param Id query string true "Miner's Id"
// @Success 200 {object} whitelist.Miner
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miner [get]
func (ctrl MinerController) getSpecificMiner(c *gin.Context) {
	params := c.Request.URL.Query()
	currentUser := currentUser(c)
	if len(params[Id]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectMinerIDSent.Error()})
		return
	}
	id := params[Id][0]
	miner, err := ctrl.minerService.getSpecificMiner(id, currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errCannotFindMiner {
			status = http.StatusNotFound
		} else if err == errMinerNotFoundForUser {
			status = http.StatusForbidden
		}

		c.JSON(status, api.ErrorResponse{Error: err.Error()})
		return
	}
	if miner.DeletedAt != nil {
		miner, err = ctrl.minerService.getSpecificDisabledMiner(id, currentUser, miner.DeletedAt)
	}
	if err != nil {
		status := http.StatusInternalServerError
		if err == errCannotFindMiner {
			status = http.StatusNotFound
		} else if err == errMinerNotFoundForUser {
			status = http.StatusForbidden
		}

		c.JSON(status, api.ErrorResponse{Error: err.Error()})
		return
	}
	if miner.DeletedAt == nil {
		uptimes, err := ctrl.minerService.GetUptimeNoDate(miner.Nodes, false)
		if err == nil {
			for i := 0; i < len(miner.Nodes); i++ {
				addUptimeToNode(&miner.Nodes[i], uptimes)
			}
		}
	}

	c.JSON(http.StatusOK, miner)
}

// @Summary Get specific miner for admin
// @Description Returns miner for given miner id
// @Tags miners
// @Accept json
// @Produce json
// @Param Id query string true "Miner's Id"
// @Success 200 {object} whitelist.Miner
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miner [get]
func (ctrl MinerController) getSpecificMinerForAdmin(c *gin.Context) {
	params := c.Request.URL.Query()
	if len(params[Id]) <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse{Error: errIncorrectMinerIDSent.Error()})
		return
	}
	id := params[Id][0]
	miner, err := ctrl.minerService.getSpecificMinerForAdmin(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	if miner.DeletedAt != nil { //TODO check why miner is loaded if ignored with miner.DeletedAt == nil
		miner, err = ctrl.minerService.getSpecificDisabledMinerForAdmin(id, miner.DeletedAt)
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}
	if miner.DeletedAt == nil {
		uptimes, err := ctrl.minerService.GetUptimeNoDate(miner.Nodes, false)

		if err == nil {
			for i := 0; i < len(miner.Nodes); i++ {
				addUptimeToNode(&miner.Nodes[i], uptimes)
			}
		}
	}
	c.JSON(http.StatusOK, miner)
}

// @Summary Update miner data
// @Description Update specific miner according to request data
// @Tags miners
// @Accept json
// @Produce json
// @Param updateMinerReq body updateMinerReq true "Request for updating miner"
// @Success 200
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/miner [post]
func (ctrl MinerController) updateMiner(c *gin.Context) {
	var updateMinerReq updateMinerReq
	if err := c.BindJSON(&updateMinerReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	currentUser := currentUser(c)
	var (
		miner      Miner
		err        error
		minerNodes []Node
		keepNodes  []*Node
		appError   ApplicationError
	)

	if miner, err = ctrl.minerService.getSpecificMinerWithApplications(updateMinerReq.Id, currentUser); err != nil {
		log.Errorf("Error %v while trying to find a miner record to be updated using input %v", err, updateMinerReq)
		c.AbortWithStatusJSON(http.StatusNotFound, api.ErrorResponse{Error: err.Error()})
		return
	}

	if miner.Type == OFFICIAL {
		if len(updateMinerReq.Nodes) > 8 {
			log.Errorf("Error while updating official miner with input %v. There are more then 8 node keys. Only first 8 will be saved.", updateMinerReq)
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: "miner controller: too many nodes"})
			return
		}

		if len(updateMinerReq.Nodes) > 0 { //TODO double-check but this if is probably not needed
			for _, n := range updateMinerReq.Nodes {
				if len(minerNodes) == 8 {
					break
				}
				keepNodes = append(keepNodes, &Node{ID: n.Id, CreatedAt: n.CreatedAt, Key: n.Key, MinerID: miner.ID})
				minerNodes = append(minerNodes, Node{ID: n.Id, CreatedAt: n.CreatedAt, Key: n.Key, MinerID: miner.ID})
			}
		}

		ctrl.whitelistService.ValidateNodeKeys(minerNodes, 0, miner.ID, &appError)
		if appError.HasErrors() {
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

			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: appError.Error.Error()})
			return
		}

		for i := len(minerNodes); i < 8; i++ {
			keepNodes = append(keepNodes, &Node{Key: "", MinerID: miner.ID})
			minerNodes = append(minerNodes, Node{Key: "", MinerID: miner.ID})
		}
		miner.Nodes = minerNodes

		if err := ctrl.minerService.UpdateMiner(&miner, keepNodes); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}

		return

	}
	// handle removal of nodes
	removedCount := 0
	for i, originalLength := 0, len(miner.Nodes); i < originalLength; i++ {
		j := i - removedCount
		shouldBeKept := false
		oldNode := miner.Nodes[j]
		for _, newNode := range updateMinerReq.Nodes {
			if oldNode.ID == newNode.Id && oldNode.Key == newNode.Key {
				keepNodes = append(keepNodes, &Node{ID: oldNode.ID, Key: oldNode.Key, MinerID: oldNode.MinerID})
				shouldBeKept = true
				break
			}
		}
		if !shouldBeKept {
			miner.Nodes = append(miner.Nodes[:j], miner.Nodes[j+1:]...)
			removedCount++
		}
	}

	var newNodes []Node
	//TODO: check the usage for this, and how we can improve
	var nodesForCheck []Node

	for _, node := range updateMinerReq.Nodes {
		nodesForCheck = append(nodesForCheck, Node{Key: node.Key})
		if node.Id == 0 {
			newNodes = append(newNodes, Node{Key: node.Key})
		}

		for _, oldNode := range miner.Nodes {
			if oldNode.ID == node.Id {
				if oldNode.Key != node.Key {
					newNodes = append(newNodes, Node{Key: node.Key})
				}
				break
			}
		}
	}

	ctrl.whitelistService.ValidateNodeKeys(nodesForCheck, 0, miner.ID, &appError)
	if appError.HasErrors() {
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
	if len(newNodes) > 0 || removedCount > 0 {
		shouldUpdateMiner := removedCount > 0

		//TODO (update improvements) cover the edge case if someone removes all DIY miner nodes and wants to add new ones
		// 	in that case we need to fetch any node related to the miner (including deleted ones)
		//	since there's very slim chance DIY Miner was approved without any node
		if len(newNodes) > 0 {
			for i, orgLength := 0, len(newNodes); i < orgLength; i++ {
				j := i - (orgLength - len(newNodes))
				if len(keepNodes) == miner.ApprovedNodesCount {
					break
				}
				keepNodes = append(keepNodes, &Node{Key: newNodes[j].Key, MinerID: newNodes[j].MinerID})
				newNodes = append(newNodes[:j], newNodes[j+1:]...)
				shouldUpdateMiner = true
			}
		}

		if shouldUpdateMiner {
			if err := ctrl.minerService.UpdateMiner(&miner, keepNodes); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			}
		}
	}

	if len(newNodes) == 0 {
		return
	}

	if miner.Type != DIY {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToAddNodesToMiner.Error()})
		return
	}

	if miner.ApplicationID > 0 {
		minerOrigin, err := ctrl.whitelistService.GetWhitelist(fmt.Sprint(miner.ApplicationID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, api.ErrorResponse{Error: err.Error()})
			return
		}

		if len(minerOrigin.ChangeHistory) == 0 {
			log.Error("Unable to find application used to approve the miner with change history records")
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: errUnableToProcessChanges.Error()})
			return
		}
		minerOriginLatestCH := minerOrigin.GetLatestChangeHistory() // minerOrigin.ChangeHistory[0] //TODO check here
		pendingApp, err := ctrl.whitelistService.db.findPendingApplication(currentUser)

		if err == nil && minerOrigin.ID != pendingApp.ID {
			//automatically decline User's application that was pending, if any, that's not the current Miner update process
			var ChangeApplicationStatus = ChangeApplicationStatus{
				ApplicationId: pendingApp.ID,
				Status:        AUTO_DISABLED,
				AdminComment:  autoDeclineMessage,
			}

			if _, err := ctrl.whitelistService.UpdateApplicationStatus(ChangeApplicationStatus); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
				return
			}
		}
		if minerOrigin.ID == pendingApp.ID {
			//TODO preload nodes for ChangeHistory
			//this is the case when User is adding Nodes on Miner edit page while already restaring Application for that Miner
			// newChangeHistory = minerOriginLatestCH

			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errAppInProcessForMiner.Error()})
			return
		}

		if err := ctrl.whitelistService.RestartApplication(newNodes, minerOrigin, minerOriginLatestCH.Location, minerOriginLatestCH.Description); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
			return
		}
	} else {
		appReq := ApplicationReq{
			Nodes: newNodes,
		}

		if pendingApp, err := ctrl.whitelistService.db.findPendingApplication(currentUser); err == nil {
			var ChangeApplicationStatus = ChangeApplicationStatus{
				ApplicationId: pendingApp.ID,
				Status:        AUTO_DISABLED,
				AdminComment:  autoDeclineMessage,
			}

			if _, err := ctrl.whitelistService.UpdateApplicationStatus(ChangeApplicationStatus); err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
				return
			}
		}

		app, appErr := ctrl.whitelistService.CreateApplication(appReq, miner.Username)
		if appErr.HasErrors() {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
			return
		}

		minerOrigin := miner.Applications[0]
		if minerOrigin.ID == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, api.ErrorResponse{Error: err.Error()})
			return
		}

		changeHistory := app.GetLatestChangeHistory()
		changeHistory.AdminComment = fmt.Sprint("Originated from application with id: ", minerOrigin.ID)
		if err := ctrl.whitelistService.db.updateChangeHistory(&changeHistory); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
			return
		}

		miner.ApplicationID = app.ID
		miner.Applications = append(miner.Applications, &app)
		if err := ctrl.minerService.UpdateMiner(&miner, []*Node{}); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: err.Error()})
			return
		}
	}

	ctrl.notifyUserAfterMinerUpdate(currentUser)
}

// @Summary Gets import data
// @Description If available, returns miner import data
// @Tags miners
// @Accept json
// @Produce json
// Success 200 {array}  whitelist.MinerImport
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/import [get]
func (ctrl MinerController) getImportData(c *gin.Context) {
	importData, err := ctrl.minerService.GetImportData()

	if err != nil {
		if err == errCannotFindMinerImportData {
			c.AbortWithStatusJSON(http.StatusNotFound, api.ErrorResponse{Error: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, importData)
}

// @Summary Updates import data
// @Description Updates miner import data with information from request
// @Tags miners
// @Accept json
// @Produce json
// @Param  importedData body updateImportDataReq true "Request containing data for update"
// @Success 200 {string}  string
// @Failure 422 {object} api.ErrorResponse
// @Router /miners/import [post]
func (ctrl MinerController) updateImportData(c *gin.Context) {
	if _, err := ctrl.persistMinerDataFromRequest(c); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}
	c.JSON(http.StatusOK, "Updated data successfully")
}

// @Summary Process import data
// @Description Import users and miners from  import data request
// @Tags miners
// @Accept json
// @Produce json
// @Param  importedData body updateImportDataReq true "Request containing data for importing"
// @Success 200 {string}  string
// @Failure 422 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /miners/import/process [post]
func (ctrl MinerController) processImportData(c *gin.Context) {
	importedData, errOnImport := ctrl.persistMinerDataFromRequest(c)
	if errOnImport != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, api.ErrorResponse{Error: errUnableToProcessRequest.Error()})
		return
	}

	for _, data := range importedData {
		created, err := ctrl.userService.ImportUser(data.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}

		if err := ctrl.minerService.ImportMiners(data); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrorResponse{Error: err.Error()})
			return
		}

		if viper.GetBool("official-import.email-notification-sending") {
			if created {
				ctrl.notifyCreatedUserAfterImport(data.Username)
			} else {
				ctrl.notifyUpdatedUserAfterImport(data.Username)
			}
		}
	}

	c.JSON(http.StatusOK, "Imported data successfully")
}

func (ctrl MinerController) UpdateUsersFromShopify() error {
	productId := viper.GetString("shopify.product-id")
	shopifyOrders, err := getShopifyData()
	if err != nil {
		return err
	}
	for _, order := range shopifyOrders.Orders {
		totalMiners := 0
		databaseRecord, err := ctrl.minerService.getShopInfo(order.Id)
		if err != errCannotFindShopRecord && err != nil {
			return errUpdatingShopOrders
		}
		if databaseRecord.Status != "paid" {
			if order.FinancialStatus == "paid" {
				for _, lineItem := range order.LineItems {
					if lineItem.Sku == productId {
						totalMiners += lineItem.Quantity
					}
				}
				ctrl.minerService.updateStoreInfo(order, databaseRecord)
				created, err := ctrl.userService.ImportUser(order.Email)
				if err != nil {
					log.Error("Error creating user ", err)
					continue
				}
				//check if user  had prior third batch miners before updating with new ones
				hadThird, err := ctrl.userService.CheckIfUserHasThirdBatchMiners(order.Email)
				if err != nil {
					log.Error("Error while checking for prior third batch miners due to ", err)
				}

				data := MinerImport{
					Username:       order.Email,
					NumberOfMiners: totalMiners,
					BatchLabel:     "Third",
				}
				if err := ctrl.minerService.ImportMinersNoRemoval(data); err != nil {
					log.Error("Error creating miner ", err)
					continue
				}
				if viper.GetBool("shopify.email-notification-sending") {
					if created {
						ctrl.notifyCreatedUserAfterImport(order.Email)
					} else {
						ctrl.notifyUpdatedUserAfterImport(order.Email)
						if hadThird {
							ctrl.sendWelcomeEmail(order.Email)
						}
					}
				}
			}
		}

	}
	return nil
}
func (ctrl MinerController) NotifyUsersAboutNoUptime() error {
	exportThreshold := viper.GetFloat64("export.reward-percentage") * 100

	var activityStart, activityEnd time.Time
	var currentYear int
	var thisMonth time.Month
	currentLocation := activityStart.Location()

	currentTime := time.Now()
	currentYear, thisMonth, _ = currentTime.Date()
	var previousMonth = thisMonth - 1
	var previousYear = currentYear
	if thisMonth == 0 {
		previousMonth = 12
		previousYear = currentYear - 1
	}
	activityStart = time.Date(previousYear, time.Month(previousMonth), 1, 0, 0, 0, 0, currentLocation)
	activityEnd = time.Date(currentYear, thisMonth, 1, 0, 0, 0, 0, currentLocation)

	difference := activityStart.Sub(activityEnd).Seconds()
	quarter := int64(difference * (1 - viper.GetFloat64("export.reward-percentage")))
	lastAcceptedCreation := activityStart.Add(time.Second * time.Duration(quarter))

	users, err := ctrl.userService.GetUsers()
	if err != nil {
		log.Error("Unable to get list of users to perform node uptime check due to ", err)
		return err
	}
	for _, user := range users {
		var shouldSend bool = true
		miners, err := ctrl.minerService.GetMinersForUser(user.Username)
		if err != nil {
			log.Errorf("Error while fetching miners for user %v, unable to check reward eligibility", user.Username)
			continue
		}
		var nodesForUptimeService []Node
		for _, miner := range miners {
			if miner.CreatedAt.After(lastAcceptedCreation) {
				break
			}
			for _, node := range miner.Nodes {
				nodesForUptimeService = append(nodesForUptimeService, node)
			}
		}
		if len(nodesForUptimeService) == 0 {
			continue
		}
		uptimes, err := ctrl.minerService.GetUptimeNoDate(nodesForUptimeService, true)
		if err != nil {
			log.Errorf("Error while fetching uptimes for nodes of user %v, unable to check reward eligibility due to error %v", user.Username, err)
			continue
		}
		for _, uptime := range uptimes {
			if uptime.Percentage >= exportThreshold {
				shouldSend = false
				log.Infof("Not sending email to %v because uptime is %v", user.Username, uptime.Percentage)
				break
			}
		}
		log.Infof("Preparing to send email to %v", user.Username)
		if shouldSend {
			ctrl.notifyUserAboutNotEnoughUptime(user.Username)
		}
	}
	return nil

}

func (ctrl MinerController) RemindUsersToUpdateAddress() error {
	log.Info("Checking for and reminding users eligible for rewards that forgot to update their skycoin address process started")
	users, err := ctrl.userService.GetUsersWithAddressesAndMiners()
	if err != nil {
		return err
	}
	timeUntilEndOfMonth := viper.GetFloat64("reminder.days-until-end-of-month")
	remainingDaysUptime := timeUntilEndOfMonth * 24 * time.Hour.Seconds()
	daysInMonth := float64(daysInMonth(time.Now().Month(), time.Now().Year()))
	totalUptimePossible := daysInMonth * 24 * time.Hour.Seconds()
	rewardLevel := viper.GetFloat64("export.reward-percentage")
	uptimeRequired := rewardLevel * totalUptimePossible

	for _, user := range users {
		if len(user.Addresses) > 0 {
			continue
		}
		var shouldSend bool
		miners, err := ctrl.minerService.GetMinersForUser(user.Username)
		if err != nil {
			if err != ErrCannotFindUser {
				log.Errorf("Error while fetching miners for user %v, unable to check reward eligibility due to %v", user.Username, err)
				continue
			}
			continue
		}
		var nodesForUptimeService []Node
		for _, miner := range miners {
			if len(miner.Nodes) > 0 {
				for _, node := range miner.Nodes {
					if len(node.Key) > 0 {
						nodesForUptimeService = append(nodesForUptimeService, node)
					}
				}
			}
		}
		uptimes, err := ctrl.minerService.GetUptimeNoDate(nodesForUptimeService, false)
		if err != nil {
			log.Errorf("Error while fetching uptimes for nodes of user %v, unable to check reward eligibility due to error %v", user.Username, err)
			continue
		}
		for _, uptime := range uptimes {
			potentialUptime := uptime.Uptime + remainingDaysUptime
			if potentialUptime >= uptimeRequired {
				shouldSend = true
				break
			}
		}
		if shouldSend {
			ctrl.notifyUserAboutMissingAddress(user.Username)
		}
	}
	log.Info("Finished sending out reminder emails to users that are eligible for rewards to update their skycoin address")
	return nil

}

func (ctrl MinerController) notifyUserAfterMinerDeletion(email string) {
	if mailErr := ctrl.mailService.MinerDeleted(email); mailErr != nil {
		log.Errorf("Unable to notify user %v about miner removed from his account", email)
	} else {
		log.Debugf("Successfully notified user %v about removing miner from his account", email)
	}
}
func (ctrl MinerController) notifyUserAfterMinerReenabling(email string) {
	if mailErr := ctrl.mailService.MinerReenabled(email); mailErr != nil {
		log.Errorf("Unable to notify user %v about reenabled miner", email)
	} else {
		log.Debugf("Successfully notified user %v about reenabling miner", email)
	}
}

func (ctrl MinerController) notifyUpdatedUserAfterImport(email string) {
	if mailErr := ctrl.mailService.AccountCreatedShopifyImport(email); mailErr != nil {
		log.Errorf("Unable to notify user %v about miners added to his account", email)
	} else {
		log.Debugf("Successfully notified user %v about adding miners to his account", email)
	}
}

func (ctrl MinerController) notifyCreatedUserAfterImport(email string) {
	if mailErr := ctrl.mailService.AccountCreatedShopifyImport(email); mailErr != nil {
		log.Errorf("Unable to notify user %v that new account with his miners has been created", email)
	} else {
		log.Debugf("Successfully invited user %v to his account with pre imported miners", email)
	}
}
func (ctrl MinerController) sendWelcomeEmail(email string) {
	if mailErr := ctrl.mailService.WelcomeEmailThirdBatch(email); mailErr != nil {
		log.Error("Unable to send welcome email to user ", email)
	} else {
		log.Debug("Successfully sent welcome email to user ", email)
	}
}
func (ctrl MinerController) notifyUserAfterMinerUpdate(email string) {
	if mailErr := ctrl.mailService.MinerNodesUpdates(email); mailErr != nil {
		log.Errorf("Unable to notify user %v that the application has been resubmitted", email)
	} else {
		log.Debugf("Successfully notified user %v about application resubmitting", email)
	}
}
func (ctrl MinerController) notifyUserAfterMinerTransfer(emailTo, emailFrom string) {
	if mailErr := ctrl.mailService.MinerTransferred(emailTo, emailFrom); mailErr != nil {
		log.Errorf("Unable to notify user %v about miner transfer from user %v", emailTo, emailFrom)
	} else {
		log.Debugf("Successfully notified user %v about miner transfer from user %v", emailTo, emailFrom)
	}
}
func (ctrl MinerController) notifyUserAboutNotEnoughUptime(email string) {
	if mailErr := ctrl.mailService.NotifyUserAboutNoUptime(email); mailErr != nil {
		log.Errorf("Unable to notify user %v about not enough uptime on nodes", email)
	} else {
		log.Debugf("Successfully notified user %v about not enough uptime on nodes", email)
	}
}

func (ctrl MinerController) notifyUserAboutMissingAddress(email string) {
	if mailErr := ctrl.mailService.RemindUserAboutRewardsMissingAddress(email); mailErr != nil {
		log.Errorf("Unable to notify user %v about his missing address", email)
	} else {
		log.Debugf("Successfully notified user %v about missing address", email)
	}
}

func (ctrl MinerController) isAdminMiddleware(c *gin.Context) {
	usr, err := ctrl.userService.FindBy(currentUser(c))
	if err != nil || !usr.IsAdmin() {
		c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse{Error: errNoAdminPrivlages.Error()})
		return
	}
}

type updateImportDataReq struct {
	Data []MinerImport `json:"data"`
}

type updateMinerReq struct {
	Id    string
	Nodes []updatedMinerNode
}

type updatedMinerNode struct {
	Id        uint
	MinerId   uint
	Key       string
	CreatedAt time.Time
}

type exportMinersReq struct {
	StartDate string
	EndDate   string
}

type minerExport struct {
	Mail       string
	Key        string
	Uptime     string
	Address    string
	Type       string
	BatchLabel string
}

type transferMinerReq struct {
	MinerId    string
	TransferTo string
}

func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func addUptimeToNode(node *Node, uptimeResponses []NodeUptimeResponse) {
	for _, resp := range uptimeResponses {
		if node.Key == resp.Key {
			node.Uptime = resp
			break
		}
	}

}

func mapType(nodeType MinerType) string {
	if nodeType == 0 {
		return "OFFICIAL"
	}
	return "DIY"
}

func convertDate(value string) time.Time {
	i, _ := strconv.ParseInt(value, 10, 64)
	tm := time.Unix(i, 0)
	fmt.Print(tm)
	return tm
}

func (ctrl MinerController) RunningRoutineForShopify() {
	diff := viper.GetDuration("shopify.refresh-interval")
	jobTicker := &jobTicker{}
	ctrl.UpdateUsersFromShopify()
	jobTicker.updateTimer(diff)
	for {
		<-jobTicker.timer.C
		log.Info("Scheduler triggered", time.Now())
		jobTicker.updateTimer(diff)
		ctrl.UpdateUsersFromShopify()
	}
}

func (ctrl MinerController) RunningRoutineForUptimeNotification() {
	jobTicker := &jobTicker{}
	currentDate := time.Now()
	y, m, _ := currentDate.Date()
	if viper.GetBool("reminder.uptime-notification-at-startup") {
		ctrl.NotifyUsersAboutNoUptime()
	}
	startOfNextMonth := time.Date(y, m+1, 1, 0, 0, 0, 0, currentDate.Location())
	jobTicker.updateTimer(startOfNextMonth.Sub(currentDate))
	for {
		<-jobTicker.timer.C
		log.Info("Scheduler for notifying about no uptime triggered", time.Now())
		startOfNextMonth = startOfNextMonth.AddDate(0, 1, 0)
		currentDate := time.Now()
		jobTicker.updateTimer(startOfNextMonth.Sub(currentDate))
		ctrl.NotifyUsersAboutNoUptime()
	}
}

func (ctrl MinerController) RunningRoutineForRewardsAddressReminderEmails() {
	daysUntilEndOfMonth := viper.GetFloat64("reminder.days-until-end-of-month")
	timeUntilEndOfMonth := time.Duration(daysUntilEndOfMonth*24) * time.Hour
	jobTicker := &jobTicker{}
	currentDate := time.Now()
	y, m, _ := currentDate.Date()
	endOfMonth := time.Date(y, m, 1, 0, 0, 0, 0, currentDate.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond)
	nextReminderDate := endOfMonth.Add(-timeUntilEndOfMonth)
	if viper.GetBool("reminder.run-uptime-reward-add-adddress-reminder-at-startup") {
		ctrl.RemindUsersToUpdateAddress()
	}
	endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month()+1, 1, 0, 0, 0, 0, endOfMonth.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond)
	nextReminderDate = endOfMonth.Add(-timeUntilEndOfMonth)
	jobTicker.updateTimer(nextReminderDate.Sub(currentDate))
	log.Info("Scheduler for missing address reminder will trigger on ", nextReminderDate)
	for {
		<-jobTicker.timer.C
		endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month()+1, 1, 0, 0, 0, 0, endOfMonth.Location()).AddDate(0, 1, 0).Add(-time.Nanosecond)
		currentDate = time.Now()
		nextReminderDate = endOfMonth.Add(-timeUntilEndOfMonth)
		log.Info("Next reminder scheduled for: ", nextReminderDate)
		jobTicker.updateTimer(nextReminderDate.Sub(currentDate))
		ctrl.RemindUsersToUpdateAddress()
	}
}
