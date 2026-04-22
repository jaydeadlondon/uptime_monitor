package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string    `gorm:"not null" json:"-"`
	TelegramChatID *int64    `json:"telegram_chat_id,omitempty"`
	LinkToken      *string   `gorm:"index" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Monitors []Monitor `gorm:"foreignKey:UserID" json:"-"`
}