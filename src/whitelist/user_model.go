package whitelist

import (
	"time"

	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

type ExportType uint8

// User for User entity.
type User struct {
	Status         uint8                 `json:"status"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"-"`
	DeletedAt      *time.Time            `json:"-"`
	Username       string                `json:"username" gorm:"primary_key" example:"someone@mail.com"`
	Rights         []authorization.Right `json:"rights,omitempty" gorm:"-" sql:"-"`
	Applications   []Application         `json:"applications" gorm:"foreignkey:Username"`
	Miners         []Miner               `json:"miners,omitempty" gorm:"foreignkey:Username; PRELOAD:false"`
	ApiKeys        []ApiKey              `json:"apiKeys,omitempty" gorm:"foreignkey:Username; PRELOAD:false"`
	Addresses      []Address             `json:"addressHistory,omitempty" gorm:"foreignkey:Username; PRELOAD:false"`
	ExportRecords  []ExportRecord 		 `json:"exportRecords,omitempty" gorm:"foreignkey:Username; PRELOAD:false"`
}

type Address struct {
	ID             uint       `json:"id"`
	Username       string     `json:"username"`
	SkycoinAddress string     `json:"skycoinAddress" example:"2mj62xkVrB24z1rbVwBfLtGPe4FkFP1eKaW"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `json:"-"`
}

type ApiKey struct {
	ID        uint       `gorm:"primary_key" json:"id" example:"1"`
	Key       string     `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Username  string     `json:"username"`
}

type ExportRecord struct {
	ID            uint 			`gorm:"primary_key" json:"id"`
	PayoutAddress string		`json:"payoutAddress"`
	MinerType     MinerType 	`json:"minerType"`
	OfficialTx    string		`json:"officialTx"`
	DiyTx         string		`json:"diyTx"`
	CorrectionTx  string		`json:"correctionTx"`
	TimeOfExport  time.Time		`json:"timeOfExport"`
	CreatedAt     time.Time  	`json:"createdAt"`
	UpdatedAt     time.Time  	`json:"-"`
	DeletedAt     *time.Time 	`json:"-"`
	Username      string     	`json:"userId"`
	NumberOfNodes  int			`json:"numberOfNodes"`
}
const (
	vipUserMask             uint8 = 1 // 0000 0001
	canSumbmitWhitelistMask uint8 = 2 // 0000 0010

	flagVIPUserMask     uint8 = 128 // 1000 0000
	reviewWhitelistMask uint8 = 64  // 0100 0000

	isAdminMask uint8 = 248 // 1111 0000
)

func (m *User) BlockedFromSubmitingWhitelist() bool {
	return (m.Status & canSumbmitWhitelistMask) > 0
}

func (m *User) IsAdmin() bool {
	return (m.Status & isAdminMask) > 0
}

func (m *User) CanFlagUserAsVIP() bool {
	return (m.Status & flagVIPUserMask) > 0
}

func (m *User) SetCanFlagUserAsVIP(canFlagUserAsVIP bool) {
	m.setRole(flagVIPUserMask, canFlagUserAsVIP)
}

func (m *User) IsUserVIP() bool {
	return (m.Status & vipUserMask) > 0
}

func (m *User) SetIsUserVIP(isUserVIP bool) {
	m.setRole(vipUserMask, isUserVIP)
}

func (m *User) CanReviewWhitelsit() bool {
	return (m.Status & reviewWhitelistMask) > 0
}

func (m *User) SetCanReviewWhitelsit(canReviewWhitelist bool) {
	m.setRole(reviewWhitelistMask, canReviewWhitelist)
}

func (m *User) setRole(roleMask uint8, isActive bool) {
	if isActive {
		m.Status |= roleMask
	} else {
		m.Status &= ^roleMask
	}
}

func (User) TableName() string {
	return "users"
}
