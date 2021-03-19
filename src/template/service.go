package template

import (
	"github.com/mattbaird/gochimp"
	"github.com/spf13/viper"
)

type Service struct {
	api *gochimp.MandrillAPI
}

// NewService prepares new instance of Service
func NewService(api *gochimp.MandrillAPI) Service {
	return Service{
		api: api,
	}
}

func DefaultService() Service {
	return NewService(api)
}

// AccountUpdatedAfterMinerImport sends an email to the user with new miners added after import
func (mailService *Service) AccountUpdatedAfterMinerImport(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-miners-page")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_account": link}

	return baseSend(mailService.api, receiver, "miners-imported-for-existing-account", attributes)
}

// AccountCreatedAfterMinerImport sends an email to the user with new miners added after import
func (mailService *Service) AccountCreatedAfterMinerImport(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_whitelist_login_page": link}

	return baseSend(mailService.api, receiver, "miners-imported-for-new-account", attributes)
}

// WhitelistApplicationCreated is sent when user successfully created his whitelisting application
func (mailService *Service) WhitelistApplicationCreated(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_login": link}

	return baseSend(mailService.api, receiver, "application-created", attributes)
}

// WhitelistApplicationUpdated is sent when Administrator changes the status of active whitelisting application
func (mailService *Service) WhitelistApplicationUpdated(receiver string, status string, comment string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{
		"link_to_login":         link,
		"status_of_application": status,
		"comment":               comment,
	}

	return baseSend(mailService.api, receiver, "application-updated", attributes)
}

// MinerNodesUpdates is sent when user sucessfuly adds new nodes to miner and application is resubmitted
func (mailService *Service) MinerNodesUpdates(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_login": link}

	return baseSend(mailService.api, receiver, "miner-nodes-added", attributes)
}

// MinerTransferred is sent when user tranfers his miner to some other user
func (mailService *Service) MinerTransferred(receiver, sender string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{
		"link_to_login":       link,
		"link_of_user_sender": sender,
	}

	return baseSend(mailService.api, receiver, "transfer-miner", attributes)
}

// MinerDeleted is sent when admin deletes miner, so the user gets notified
func (mailService *Service) MinerDeleted(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_login": link}

	return baseSend(mailService.api, receiver, "deleted-miner", attributes)
}

//MinerReenabled is sent when admin reenables miner, so the user gets notified about that change
func (mailService *Service) MinerReenabled(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_login": link}

	return baseSend(mailService.api, receiver, "reenabled-miner", attributes)
}

// AccountCreatedShopifyImport sends an email to the user with new miners added after import
func (mailService *Service) AccountCreatedShopifyImport(receiver string) error {
	url := viper.GetString("server.frontend-endpoint")
	attributes := map[string]string{"base_url": url}

	return baseSend(mailService.api, receiver, "shopify-import", attributes)
}

func (mailService *Service) WelcomeEmailThirdBatch(receiver string) error {
	url := viper.GetString("server.frontend-endpoint")
	attributes := map[string]string{"base_url": url}

	return baseSend(mailService.api, receiver, "third-batch-email", attributes)
}
func (mailService *Service) NotifyUserAboutNoUptime(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-miners-page")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_account": link}

	return baseSend(mailService.api, receiver, "notify-user-about-uptime", attributes)
}

func (mailService *Service) RemindUserAboutRewardsMissingAddress(receiver string) error {
	url := viper.GetString("server.frontend-endpoint") + viper.GetString("server.route-login")
	link := "<a href='" + url + "'>" + url + "</a>"
	attributes := map[string]string{"link_to_login": link}

	return baseSend(mailService.api, receiver, "remind-user-about-address", attributes)
}
