package models

import "time"

type CheckStatus string

const (
	CheckStatusUp   CheckStatus = "up"
	CheckStatusDown CheckStatus = "down"
)

type MonitorCheck struct {
	ID             uint        `gorm:"primaryKey" json:"id"`
	MonitorID      uint        `gorm:"not null;index" json:"monitor_id"`
	Status         CheckStatus `gorm:"not null" json:"status"`
	StatusCode     *int        `json:"status_code,omitempty"`
	ResponseTimeMs *int64      `json:"response_time_ms,omitempty"`
	ErrorMessage   *string     `json:"error_message,omitempty"`
	CheckedAt      time.Time   `gorm:"index" json:"checked_at"`

	Monitor Monitor `gorm:"foreignKey:MonitorID" json:"-"`
}

type Incident struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	MonitorID       uint       `gorm:"not null;index" json:"monitor_id"`
	StartedAt       time.Time  `json:"started_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
	DurationSeconds *int64     `json:"duration_seconds,omitempty"`

	Monitor Monitor `gorm:"foreignKey:MonitorID" json:"-"`
}