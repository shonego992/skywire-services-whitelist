package whitelist

import "errors"

// Whitelist service related errors
var (
	//whitelist controller errors
	errUnableToProcessRequest         = errors.New("whitelist controller: unable to process fields from the request")
	errForbidden                      = errors.New("whitelist controller: forbidden from completing the request")
	errNoApplicationInProgressForUser = errors.New("whitelist controller: no application in progress to show for user")
	ErrNoMinerForApplication          = errors.New("whitelist controller: no miner assigned to given application")
	errNoApplicationForNode           = errors.New("whitelist controller: no application assigned to given node")
	errCannotLoadWhitelists           = errors.New("whitelist controller: cannot load whitelist applications")
	ErrCannotFindWhitelist            = errors.New("whitelist controller: cannot find whitelist application with id")
	errCannotUpdateApplication        = errors.New("whitelist controller: cannot update whitelist application")
	errIncorrectWhitelistIDSent       = errors.New("whitelist controller: incorrect value for whitelist id sent")
	errNotEnoughImages                = errors.New("whitelist controller: at least 3 images required")
	errDuplicateDatabaseImages        = errors.New("whitelist controller: Images  already exists in database: ")
	errUploadingImages                = errors.New("whitelist controller: Error while uploading images to server")

	//user controller errors
	errCannotLoadUsers         = errors.New("user controller: cannot load users")
	errCannotFindMiner         = errors.New("user controller: Cannot find miner by ID")
	errMinerNotFoundForUser    = errors.New("user controller: Miner with specified ID does not exist for user")
	errCannotUpdateMiner       = errors.New("user controller: Cannot update miner")
	ErrCannotLoadMiners        = errors.New("user controller: Cannot load miners")
	errCannotCreateMiners      = errors.New("user controller: Cannot create miners")
	errCannotLoadPayoutAddress = errors.New("user controller: Cannot find payout address")
	errNoAddressSet            = errors.New("user controller: User has no address set")
	errCannotTransferMiner     = errors.New("user controller: Cannot transfer miner")
	errNotOwnerOfMiner         = errors.New("user controller: Cannot transfer miner not in your ownership")
	errWrongMinerType          = errors.New("user controller: Update is possible only for DIY type of miners")
	errCannotFindUserByKey     = errors.New("user controller: Cannot find user by api key")
	errCreatingExportRecord    = errors.New("user controller: Error creating export record in DB")

	//miner controller errors
	errIncorrectUsernameSent   = errors.New("miner controller: incorrect value for username sent")
	errNoMinersFoundForImport  = errors.New("miner controller: no miners found for import in provided file")
	errIncorrectMinerIDSent    = errors.New("miner controller: incorrect value for miner id sent")
	errUnableToAddNodesToMiner = errors.New("miner controller: unable to add nodes to official miner")
	errUnableToProcessChanges  = errors.New("miner controller: unable to process required changes")
	errAppInProcessForMiner    = errors.New("miner controller: there is already application in process connected to this miner")

	//whitelist service errors
	errUnableToSave               = errors.New("whitelist service: unable to persist provided data")
	errTechnicalError             = errors.New("whitelist service: technical error occured")
	errApplicationInProgress      = errors.New("whitelist service: active application already exists")
	errIdenticalAsPrevious        = errors.New("whitelist service: new input is identical as the previous one")
	errCannotFindMinerImportData  = errors.New("whitelist service: cannot find imported miner data")
	errCannotLoadUserMiners       = errors.New("whitelist service: cannot load user miners")
	errCannotFindActiveNode       = errors.New("whitelist service: cannot find active node by key")
	errCannotFindPendingApp       = errors.New("whitelist service: cannot find pending app for user")
	errCannotUpdateDeclined       = errors.New("whitelist service: cannot change previous autodeclined to pending")
	errAlreadyTakenKeys           = errors.New("whitelist service: already taken keys detected. Public keys are already taken")
	errWrongNodeKeys              = errors.New("whitelist service: wrong keys detected. Wrong keys -")
	errDuplicateKeys              = errors.New("whitelist service: duplicate keys detected. Duplicate keys -")
	ErrUnableToRestartApplication = errors.New("whitelist service: unable to restart whitelist")
	errWrongDecodedLength         = errors.New("whitelist service: internal error occured while validating node key")
	errCannotImageByHash          = errors.New("whitelist service: unable to find image by that hash")
	errCannotLoadImages           = errors.New("whitelist service: cannot load images")
	errCannotLoadImagesForUser    = errors.New("whitelist service: cannot load images for user")

	//user service errors
	errCannotFindAPIKey                = errors.New("user service: cannot find API key for provided user")
	errSkycoinAddressNotValid          = errors.New("user service: provided Skycoin address has no valid format")
	errCannotFindTransferUser          = errors.New("user service: transfer user does not exists")
	errNoAdminPrivlages                = errors.New("user service: no admin priviledges")
	errUserCannotSubmitNewApplications = errors.New("user service: User is denied to submit new applications")
	errUnableToUpdateCreatedAt         = errors.New("user service: unable to update created at for user")
	errCannotFindActiveAddress         = errors.New("user service: cannot find active address in database")
	errAddressAlreadyTaken             = errors.New("user service: address already taken")

	//miner service errors
	errCannotGetShopifyData = errors.New("miner service: Cannot load shopify data")
	errCannotFindShopRecord = errors.New("miner service: Cannot find shop record")
	errUpdatingShopOrders   = errors.New("miner service: Cannot update shop records in the database")
	errCannotGetUpTimeData  = errors.New("miner service: Cannot get uptime data for nodes")
	errNoKeysToGetDataFor   = errors.New("miner service: no node keys to make request for")

	//mail service errors
	errAppReviewedEmailFailed         = errors.New("mail service: application successfully reviewed but notification email sending failed")
	errAppCreatedEmailFailed          = errors.New("mail service: application successfully created but confirmation email sending failed")
	errAppCreatedImagesAndEmailFailed = errors.New("mail service: application successfully reviewed but notification email sending failed")

	//common service errors
	errUnableToRead           = errors.New("common service: unable to query persisted data")
	errMissingMandatoryFields = errors.New("common service: missing some mandatory fields")
	ErrCannotFindUser         = errors.New("common service: cannot find user by email")
)
