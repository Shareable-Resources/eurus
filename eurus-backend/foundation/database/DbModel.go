package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

type DbModel struct {
	CreatedDate      time.Time `gorm:"column:created_date;autoCreateTime" json:"createdDate"`
	LastModifiedDate time.Time `gorm:"column:last_modified_date;autoUpdateTime" json:"lastModifiedDate"`
}

func (me *DbModel) InitDate() {
	now := time.Now()
	me.CreatedDate = now
	me.LastModifiedDate = now
}

type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(`null`), nil
}

type NullTime struct {
	sql.NullTime
}

func (s NullTime) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.Time)
	}
	return []byte(`null`), nil
}
