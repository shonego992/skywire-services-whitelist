package whitelist

import (
	"time"
)

type ApplicationStatus uint8

const (
	PENDING  = 0
	APPROVED = 1
	DECLINED = 2 //user can resubmit
	DISABLED = 3 //user can not resubmit
	CANCELED = 4
	AUTO_DISABLED = 5
)

var statusMap = map[uint8]string{
	0: "PENDING",
	1: "APPROVED",
	2: "DECLINED",
	3: "DISABLED",
	4: "CANCELED",
}

// Type for applications
type Application struct {
	ID            uint              `gorm:"primary_key" json:"id"`
	CurrentStatus ApplicationStatus `json:"currentStatus"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"-"`
	DeletedAt     *time.Time        `json:"-"`
	Username      string            `json:"userId"`
	ChangeHistory []ChangeHistory   `json:"changeHistory" gorm:"association_autoupdate:false"`
	// This connection is used to preserve current miner for app. It should not be preloaded. If miner is needed use GetMinerForApplication method in service.
	Miner Miner 					`json:"miner" gorm:"foreignkey:ApplicationID;PRELOAD:false"`
}

func (a Application) GetLatestChangeHistory() ChangeHistory {
	if len(a.ChangeHistory) == 0 {
		return ChangeHistory{}
	}
	latest := a.ChangeHistory[0]
	for _, ch := range a.ChangeHistory {
		if ch.ID > latest.ID {
			latest = ch
		}
	}
	return latest
}

// Type for application change history
type ChangeHistory struct {
	ID            uint              `gorm:"primary_key" json:"id"`
	Status        ApplicationStatus `json:"status"`
	Description   string            `json:"description"`
	Location      string            `json:"location"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"-"`
	DeletedAt     *time.Time        `json:"-"`
	AdminComment  string            `json:"adminComment" gorm:"column:comment"`
	UserComment   string            `json:"userComment"`
	ApplicationId uint              `json:"-"`
	Nodes         []*Node           `json:"nodes" gorm:"many2many:change_nodes;"`
	Images        []*Image          `json:"images" gorm:"many2many:change_images;"`
}

type Node struct {
	ID        uint               `json:"id"`
	Key       string             `json:"key"`
	MinerID   uint               `json:"minerId,omitempty" sql:"default: null"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"-"`
	DeletedAt *time.Time         `json:"disabled,omitempty"`
	Uptime    NodeUptimeResponse `json:"uptime" gorm:"-" sql:"-"`
	BatchLabel  string			 `json:"-" gorm:"-" `
}

type NodeUptimeResponse struct {
	Key        string  `json:"key"`
	Uptime     float64 `json:"uptime"`
	Downtime   float64 `json:"downtime"`
	Percentage float64 `json:"percentage"`
	Online     bool    `json:"online"`
}

type Image struct {
	ID        uint       `json:"id"`
	Path      string     `json:"path"`
	MinerID   uint       `json:"minerId,omitempty" sql:"default: null"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	ImgHash   string     `json:"imgHash"`
}

func (ChangeHistory) TableName() string {
	return "change_history"
}

type ApplicationError struct {
	AlreadyTakenKeys []string
	DuplicateKeys    []string
	WrongKeys        []string
	FailedImages     []string
	Error            error
}

func (ae ApplicationError) HasErrors() bool {
	return ae.Error != nil ||
		len(ae.AlreadyTakenKeys) > 0 ||
		len(ae.DuplicateKeys) > 0 ||
		len(ae.WrongKeys) > 0 ||
		len(ae.FailedImages) > 0
}


