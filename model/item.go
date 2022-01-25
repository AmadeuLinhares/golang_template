package model

import (
	"net/http"

	"github.com/tradersclub/TCUtils/tcerr"
)

// Item de exemplo
type Item struct {
	ID             string            `json:"id" db:"Id"`
	Token          string            `json:"token" db:"Token"`
	CreateAt       int64             `json:"create_at" db:"CreateAt"`
	ExpiresAt      int64             `json:"expires_at" db:"ExpiresAt"`
	LastActivityAt int64             `json:"last_activity_at" db:"LastActivityAt"`
	UserID         string            `json:"user_id" db:"UserId"`
	DeviceID       string            `json:"device_id" db:"DeviceId"`
	Roles          string            `json:"roles" db:"Roles"`
	IsOAuth        bool              `json:"is_oauth" db:"IsOAuth"`
	Props          map[string]string `json:"props" db:"Props"`
}

// ToItem converte uma interface{} para *Item
func ToItem(data interface{}) (*Item, error) {
	value, ok := data.(*Item)
	if !ok {
		return nil, tcerr.New(http.StatusInternalServerError, "não foi possível converter interface{} para *Item", nil)
	}
	return value, nil
}
