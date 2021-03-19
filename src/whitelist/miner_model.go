package whitelist

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MinerType uint8

const (
	OFFICIAL MinerType = 0
	DIY      MinerType = 1
)

type Miner struct {
	ID                 uint            `json:"id"`
	Username           string          `json:"username"`
	Nodes              []Node          `json:"nodes" gorm:"foreignkey:MinerID;"`
	Images             []Image         `json:"images" gorm:"foreignkey:MinerID;"`
	MinerTransfers     []MinerTransfer `json:"minerTransfers" gorm:"foreignkey:MinerID;"`
	Type               MinerType       `json:"type"`
	Gifted             bool            `json:"gifted,omitempty"`
	BatchLabel         string          `json:"batchLabel,omitempty"`
	ApprovedNodesCount int             `json:"approvedNodeCount,omitempty"`
	ApplicationID      uint            `json:"applicationId,omitempty" sql:"default: null"`
	Applications       []*Application  `json:"applications" gorm:"many2many:miner_applications;association_autoupdate:false;association_autocreate:false"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
	DeletedAt          *time.Time      `json:"disabled,omitempty"`
}

type MinerImport struct {
	ID             uint       `json:"id"`
	Username       string     `json:"username"`
	NumberOfMiners int        `json:"numberOfMiners"`
	Nodes          []Node     `json:"-" gorm:"-" sql:"-"`
	BatchLabel     string     `json:"-" gorm:"-" sql:"-"`
	Gifted         bool       `json:"-" gorm:"-" sql:"-"`
	CreatedAt      time.Time  `json:"-"`
	UpdatedAt      time.Time  `json:"-"`
	DeletedAt      *time.Time `json:"-"`
}

type MinerTransfer struct {
	gorm.Model

	OldUsername string `json:"oldUsername"`
	NewUsername string `json:"newUsername"`
	MinerID     uint   `json:"minerId,omitempty" sql:"default: null"`
}

type ShopData struct {
	ID        uint64
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Status    string
}
