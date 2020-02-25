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
