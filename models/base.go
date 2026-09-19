package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model struct
type Model struct {
	Id        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;not null"`
	Seq       int64          `json:"seq" gorm:"index"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// BeforeCreate hook table
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.Id == uuid.Nil {
		m.Id = uuid.New()
	}

	m.Seq = time.Now().UnixNano()

	return nil
}

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	if value == nil {
		*a = make(JSONB)
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for JSONB")
	}

	return json.Unmarshal(b, a)
}
