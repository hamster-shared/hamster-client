package keystorage

import (
	"time"
)

type Client interface {
	Get(key string) string
	Set(key, value string)
	Delete(key string)
	Err() error
	SetTableName(string)
}

type KeyStorage struct {
	ID        uint      `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`

	Key   string `gorm:"uniqueIndex" json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}
