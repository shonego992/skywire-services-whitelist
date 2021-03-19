package rpc_client

import (
	"net/rpc"
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func FetchCreatedAt(username string) (time.Time, error) {
	client, err := rpc.DialHTTP(viper.GetString("rpc.user.protocol"), viper.GetString("rpc.user.address"))
	if err != nil {
		log.Error("dialing:", err)
		return time.Time{}, err
	}

	args := &authorization.GetRequest{Username: username}
	var reply time.Time
	err = client.Call("Handler.GetCreatedAt", args, &reply)
	if err != nil {
		log.Error("authorization access rights fetch error: ", err)
		return time.Time{}, err
	} else {
		log.Debugf("Authorization rights fetched for %v successfully", args.Username)
	}

	return reply, nil
}
