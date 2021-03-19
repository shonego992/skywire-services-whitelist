package app

import (
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
	"github.com/SkycoinPro/skywire-services-whitelist/src/whitelist"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var userService whitelist.UserService
var whitelistService whitelist.Service
var minerService whitelist.MinerService

type Handler int

func (h *Handler) GetUserAuthorization(req *authorization.GetRequest, resp *authorization.GetResponse) error {
	log.Debug("Received the call for user authorization", req)

	canReviewWhitelist := authorization.Right{Name: "review_whitelist", Label: "Review Whitelist", Value: false}
	canFlagVIP := authorization.Right{Name: "flag_vip", Label: "Flag VIP", Value: false}

	usr, err := userService.FindBy(req.Username)
	if err != nil {
		if err != whitelist.ErrCannotFindUser {
			log.Error("Unable to collect authorization data for user", req.Username)
			return err
		}
	} else {
		canReviewWhitelist.Value = usr.CanReviewWhitelsit()
		canFlagVIP.Value = usr.CanFlagUserAsVIP()
	}

	resp.Rights = []authorization.Right{canReviewWhitelist, canFlagVIP}
	return nil
}

func (h *Handler) SetUserAuthorization(req *authorization.SetRequest, resp *authorization.SetResponse) error {
	log.Debug("Received the call for user authorization", req)

	dbUser, err := userService.FindBy(req.Username)
	if err != nil {
		if err != whitelist.ErrCannotFindUser {
			log.Errorf("User with given username %v could not be accessed", req.Username)
			return err
		}
		if err := userService.Create(req.Username); err != nil {
			log.Errorf("User with given username %v could not be persisted", req.Username)
			return err
		}

		dbUser, err = userService.FindBy(req.Username)
		if err != nil {
			log.Errorf("User with given username %v could not be accessed", req.Username)
			return err
		}

	}

	if err := userService.UpdateRights(&dbUser, req.Rights); err != nil {
		log.Error("Unable to persist provided access rights", err)
		return err
	}

	return nil
}

func (h *Handler) DeleteNodesForUser(req *authorization.SetRequest, resp *authorization.SetResponse) error {
	log.Info("Received the call for deleting user nodes", req)

	dbUser, err := userService.FindBy(req.Username)
	if err != nil {
		if err == whitelist.ErrCannotFindUser {
			log.Errorf("User with given username %v does't exist", req.Username)
		} else {
			log.Errorf("User with given username %v could not be accessed", req.Username)
		}
		return err
	}

	miners, err := minerService.GetMinersForUser(dbUser.Username)
	if err != nil {
		log.Errorf("Loading miners for user: %v failed due to error: %v", dbUser.Username, err)
		return err
	}
	var currTime = time.Now()
	for _, miner := range miners {
		for _, node := range miner.Nodes {
			if err := whitelistService.RemoveNode(&node, currTime); err != nil {
				log.Errorf("Deleting nodes for miner: %v failed due to error: %v", miner.Username, err)
			}
		}
	}

	return nil
}

func RunRPCServer(us whitelist.UserService, ws whitelist.Service, ms whitelist.MinerService) {
	userService = us
	whitelistService = ws
	minerService = ms
	a := new(Handler)
	rpc.Register(a)
	rpc.HandleHTTP()

	s, err := net.Listen(viper.GetString("rpc.protocol"), viper.GetString("rpc.host"))
	if err != nil {
		log.Fatal("Can't initialize RPC server due to error ", err)
	}
	go http.Serve(s, nil)
	log.Infof("Listening for RPC requests on %v", viper.GetString("rpc.host"))
}
