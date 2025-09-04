package models

import (
	"time"

	"github.com/google/uuid"
)

type DeviceInfo struct {
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	Language   string    `json:"language,omitempty" db:"language"`
	Platform   string    `json:"platform,omitempty" db:"platform"`
	Browser    string    `json:"browser,omitempty" db:"browser"`
	OS         string    `json:"os,omitempty" db:"os"`
	DeviceType string    `json:"device_type,omitempty" db:"device_type"`
	AppVersion string    `json:"app_version,omitempty" db:"app_version"`
	Country    string    `json:"country,omitempty" db:"country"`
	LastSeenAt time.Time `json:"last_seen_at" db:"last_seen_at"`
}

type Session struct {
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	TokenHash  string     `json:"-" db:"token_hash"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	DeviceInfo DeviceInfo `json:"device_info" db:"device_info"`
}
