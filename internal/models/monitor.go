package models

import (
	"time"
)

type MonitorStatus string

const (
	StatusUp      MonitorStatus = "up"
	StatusDown    MonitorStatus = "down"
	StatusPending MonitorStatus = "pending"
)

type Monitor struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	UserID        uint          `gorm:"not null;index" json:"user_id"`
	Name          string        `gorm:"not null" json:"name"`
	URL           string        `gorm:"not null" json:"url"`
	Interval      int           `gorm:"not null;default:5" json:"interval"` // минуты
	IsActive      bool          `gorm:"default:true" json:"is_active"`
	CurrentStatus MonitorStatus `gorm:"default:'pending'" json:"current_status"`
	LastCheckedAt *time.Time    `json:"last_checked_at,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`

	User   User           `gorm:"foreignKey:UserID" json:"-"`
	Checks []MonitorCheck `gorm:"foreignKey:MonitorID" json:"-"`
}