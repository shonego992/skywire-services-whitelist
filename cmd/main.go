package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-whitelist/src/app"
	"github.com/SkycoinPro/skywire-services-whitelist/src/auth"
	"github.com/SkycoinPro/skywire-services-whitelist/src/config"
	"github.com/SkycoinPro/skywire-services-whitelist/src/database/postgres"
	"github.com/SkycoinPro/skywire-services-whitelist/src/rpc_client"
	"github.com/SkycoinPro/skywire-services-whitelist/src/template"
	"github.com/SkycoinPro/skywire-services-whitelist/src/whitelist"
)

// @title Skywire User System API
// @version 1.0
// @description This is a Skywire User System service.

// @host localhost:8080
// @BasePath /api/v1
func main() {
	config.Init("whitelist-config")
	level, err := log.ParseLevel(viper.GetString("server.log-level"))
	if err != nil {
		log.Info("Unable to use configured log level. Using Info instead")
		level = log.InfoLevel
	}
	log.SetLevel(level)
	template.Init()

	tearDown := postgres.Init()
	defer tearDown()

	ws := whitelist.DefaultService() //TODO consider reusing these two services later in this file
	ms := whitelist.DefaultMinerService()
	us := whitelist.DefaultUserService()
	mc := whitelist.DefaultMinerController()

	if !viper.GetBool("shopify.disable-shopify") {
		go mc.RunningRoutineForShopify()
	} else {
		log.Info("Routine for shopify import not enabled in configuration")
	}

	if viper.GetBool("reminder.enable-uptime-reward-add-address-reminder") {
		go mc.RunningRoutineForRewardsAddressReminderEmails()
	} else {
		log.Info("Routine for reminder emails about missing address not enabled in configuration")
	}

	if viper.GetBool("fixup.import-created-at") {
		go importCreatedAt(&us)
	}
	if viper.GetBool("reminder.schedule-uptime-notification") {
		go mc.RunningRoutineForUptimeNotification()
	} else {
		log.Info("Routine for checking of uptime in last month not enabled in configuration")
	}
	if viper.GetBool("fixup.change-created-date-of-miners") {
		go ms.InsertCreatedAtForOldMiners()
	}

	app.RunRPCServer(us, ws, ms)

	// register all of the controllers here
	app.NewServer(
		auth.DefaultController(),
		whitelist.DefaultController(),
		whitelist.DefaultUserController(),
		mc,
	).Run()
}

func importCreatedAt(us *whitelist.UserService) {
	log.Info("Importing createdAt field from user service")
	users, err := us.GetUsers()
	if err != nil {
		log.Error("Error during fetching users ", err)
		return
	}
	count := 0
	for _, user := range users {
		count++
		if count%10 == 0 { //TODO consider step and pause time here and for admins
			log.Infof("Finished %v users", count)
			time.Sleep(5 * time.Second)
		}
		createdAt, err := rpc_client.FetchCreatedAt(user.Username)
		if err != nil {
			log.Errorf("Error while fetching created_at for username %v due to %v", user.Username, err)
			continue
		}
		err = us.UpdateCreatedAtForUser(user.Username, createdAt)
		if err != nil {
			log.Errorf("Error while updating created_at for username %v due to %v ", user.Username, err)
		}
	}
	log.Info("Finished importing created at for users, starting for admins")
	admins, err := us.GetAdmins()
	if err != nil {
		log.Error("Error during fetching admins ", err)
		return
	}
	count = 0
	for _, admin := range admins {
		count++
		if count%10 == 0 {
			log.Infof("Finished %v admins", count)
			time.Sleep(5 * time.Second)
		}
		createdAt, err := rpc_client.FetchCreatedAt(admin.Username)
		if err != nil {
			log.Errorf("Error while fetching created_at for username %v due to %v", admin.Username, err)
			continue
		}
		err = us.UpdateCreatedAtForUser(admin.Username, createdAt)
		if err != nil {
			log.Errorf("Error while updating created_at for username %v due to %v ", admin.Username, err)
		}
	}
	log.Info("Finished importing created at for admins")
}
