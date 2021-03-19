Skywire Whitelisting System API
=======================
This is the Skywire Whitelisting and Miner service. It also provides access to some User System APIs in order to properly associate miners with existing users.

**Version:** 1.0


# Skywire User System API
This is a Skywire User System service.

## Version: 1.0



### /admin/users

#### GET
##### Summary:

List all users

##### Description:

Method for admins to get list of all users

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [whitelist.User](#whitelist.user) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /admin/users/

#### GET
##### Summary:

Enable user to submit whitelist applications

##### Description:

Method for admins to enable user to submit whitelist

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | query | User email | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.User](#whitelist.user) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |

### /auth/info

#### GET
##### Summary:

Retrieve signed in User's info

##### Description:

Information about currently signed in user is collected and returned as response.

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.User](#whitelist.user) |
| 401 | Unauthorized | [api.ErrorResponse](#api.errorresponse) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/allMiners

#### GET
##### Summary:

List all miners

##### Description:

Method for admins to get list of all miners

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [whitelist.Miner](#whitelist.miner) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/import

#### GET
##### Summary:

Gets import data

##### Description:

If available, returns miner import data

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 404 | Not Found | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### POST
##### Summary:

Updates import data

##### Description:

Updates miner import data with information from request

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| importedData | body | Request containing data for update | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | string |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |

### /miners/import/process

#### POST
##### Summary:

Process import data

##### Description:

Import users and miners from  import data request

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| importedData | body | Request containing data for importing | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | string |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/miner

#### GET
##### Summary:

Get specific miner for admin

##### Description:

Returns miner for given miner id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Id | query | Miner's Id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.Miner](#whitelist.miner) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### POST
##### Summary:

Update miner data

##### Description:

Update specific miner according to request data

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| updateMinerReq | body | Request for updating miner | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/miner/

#### DELETE
##### Summary:

Removes miner

##### Description:

Removes miner for given id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | query | ID of miner to be removed | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/miners

#### GET
##### Summary:

List user's miners

##### Description:

Returns a list of miners under current user

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [whitelist.Miner](#whitelist.miner) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /miners/uploadUserList

#### POST
##### Summary:

Uploads user list

##### Description:

Exports the user list to csv file and returns number of exported users

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | integer |
| 400 | Bad Request | string |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/address

#### PATCH
##### Summary:

Update Users's Skycoin address

##### Description:

Collect, validate and store User's new Skycoin address.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| newAddress | body | New User | Yes | [whitelist.AddressUpdateReq](#whitelist.addressupdatereq) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.User](#whitelist.user) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /users/keys

#### DELETE
##### Summary:

Remove User's API key

##### Description:

Match provided API key and remove it if exists

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| keyToBeRemoved | body | User's API key to be removed | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### GET
##### Summary:

List User's API keys

##### Description:

Return collection of User's generated API keys

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ string ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### POST
##### Summary:

Generate User's API key

##### Description:

Method that is going to generate, persist and return User's new API key

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created | string |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /whitelist/application

#### GET
##### Summary:

Gets the application for curent user

##### Description:

Gets the application that is currently in progress for current user

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.Application](#whitelist.application) |
| 404 | Not Found | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

#### POST
##### Summary:

Create a new application in system for current user

##### Description:

Collect provided Application attributes from the body and create new Application in the system

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| newUser | body | New User | Yes | object |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /whitelist/linkNodes

#### POST
##### Summary:

Link nodes

##### Description:

Finds User by provided API key and create a new Miner with provided Nodes

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| linkNodesReq | body | Nodes to be linked | Yes | [whitelist.linkNodesReq](#whitelist.linknodesreq) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 422 | Unprocessable Entity | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /whitelist/updateApplication

#### POST
##### Summary:

Update an existing application in system for current user, without change in images

##### Description:

Collect provided Application attributes from the body and create new change history record for

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 |  |  |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /whitelist/whitelist

#### GET
##### Summary:

Gets whitelisted application

##### Description:

Returns an application for given application id

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | query | Whitelist application id | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [whitelist.Application](#whitelist.application) |
| 400 | Bad Request | [api.ErrorResponse](#api.errorresponse) |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### /whitelist/whitelists

#### GET
##### Summary:

Lists whitelisted applications

##### Description:

Returns an array of the whitelisted applications

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [whitelist.Application](#whitelist.application) ] |
| 500 | Internal Server Error | [api.ErrorResponse](#api.errorresponse) |

### Models


#### api.ErrorResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

#### whitelist.Address

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer |  | No |
| skycoinAddress | string |  | No |
| username | string |  | No |

#### whitelist.AddressUpdateReq

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |

#### whitelist.ApiKey

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | integer |  | No |
| username | string |  | No |

#### whitelist.Application

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| changeHistory | [ [whitelist.ChangeHistory](#whitelist.changehistory) ] |  | No |
| createdAt | string |  | No |
| currentStatus | [whitelist.ApplicationStatus](#whitelist.applicationstatus) |  | No |
| id | integer |  | No |
| miner | [whitelist.Miner](#whitelist.miner) | This connection is used to preserve current miner for app. It should not be preloaded. If miner is needed use GetMinerForApplication method in service. | No |
| userId | string |  | No |

#### whitelist.ApplicationStatus

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| whitelist.ApplicationStatus | object |  |  |

#### whitelist.ChangeHistory

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| adminComment | string |  | No |
| createdAt | string |  | No |
| description | string |  | No |
| id | integer |  | No |
| images | [ [whitelist.Image](#whitelist.image) ] |  | No |
| location | string |  | No |
| nodes | [ [whitelist.Node](#whitelist.node) ] |  | No |
| status | [whitelist.ApplicationStatus](#whitelist.applicationstatus) |  | No |
| userComment | string |  | No |

#### whitelist.ExportRecord

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| correctionTx | string |  | No |
| createdAt | string |  | No |
| diyTx | string |  | No |
| id | integer |  | No |
| minerType | [whitelist.MinerType](#whitelist.minertype) |  | No |
| numberOfNodes | integer |  | No |
| officialTx | string |  | No |
| payoutAddress | string |  | No |
| timeOfExport | string |  | No |
| userId | string |  | No |

#### whitelist.Image

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| id | integer |  | No |
| imgHash | string |  | No |
| minerId | integer |  | No |
| path | string |  | No |

#### whitelist.Miner

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| applicationId | integer |  | No |
| applications | [ [whitelist.Application](#whitelist.application) ] |  | No |
| approvedNodeCount | integer |  | No |
| batchLabel | string |  | No |
| createdAt | string |  | No |
| disabled | string |  | No |
| gifted | boolean |  | No |
| id | integer |  | No |
| images | [ [whitelist.Image](#whitelist.image) ] |  | No |
| minerTransfers | [ [whitelist.MinerTransfer](#whitelist.minertransfer) ] |  | No |
| nodes | [ [whitelist.Node](#whitelist.node) ] |  | No |
| type | [whitelist.MinerType](#whitelist.minertype) |  | No |
| username | string |  | No |

#### whitelist.MinerTransfer

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| minerId | integer |  | No |
| newUsername | string |  | No |
| oldUsername | string |  | No |

#### whitelist.MinerType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| whitelist.MinerType | object |  |  |

#### whitelist.Node

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| id | integer |  | No |
| key | string |  | No |
| minerId | integer |  | No |
| uptime | [whitelist.NodeUptimeResponse](#whitelist.nodeuptimeresponse) |  | No |

#### whitelist.NodeUptimeResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| downtime | number |  | No |
| key | string |  | No |
| online | boolean |  | No |
| percentage | number |  | No |
| uptime | number |  | No |

#### whitelist.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| addressHistory | [ [whitelist.Address](#whitelist.address) ] |  | No |
| apiKeys | [ [whitelist.ApiKey](#whitelist.apikey) ] |  | No |
| applications | [ [whitelist.Application](#whitelist.application) ] |  | No |
| createdAt | string |  | No |
| exportRecords | [ [whitelist.ExportRecord](#whitelist.exportrecord) ] |  | No |
| miners | [ [whitelist.Miner](#whitelist.miner) ] |  | No |
| rights | string |  | No |
| status | integer |  | No |
| username | string |  | No |

#### whitelist.linkNodesReq

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| key | string |  | No |
| nodeKeys | [ string ] |  | No |