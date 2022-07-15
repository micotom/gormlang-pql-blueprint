package models

import (
	"gorm.io/datatypes"
)

type ValueEntry struct {
	ID         uint           `gorm:"primaryKey"`
	Day        datatypes.Date `json:"day"` // "2020-07-17 00:00:00"
	Value      int            `json:"value"`
	RaisePerc  float32        `json:"raise_perc"`
	RaiseDiff  int            `json:"raise_diff"`
	PlayerSlug string
}
