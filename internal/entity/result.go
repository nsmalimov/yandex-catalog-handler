package entity

import (
	"time"

	"github.com/jmoiron/sqlx/types"
)

type Result struct {
	ID        int            `db:"id"`
	Cause     string         `db:"cause"`
	Results   types.JSONText `db:"results" json:"results"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type ResultByFile struct {
	FileName   string `json:"filename"`
	Was        int    `json:"was"`
	Now        int    `json:"now"`
	WasRemove  int    `json:"was_remove"`
	ErrorCause string `json:"error_cause"`
}
