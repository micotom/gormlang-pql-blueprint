package models

import (
	"errors"
	"time"

	"funglejunk.com/kick-api/src/util"
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

func ValueEntryAt(ves []ValueEntry, y int, m time.Month, d int) (*ValueEntry, error) {
	return util.FindBy(ves, func(v ValueEntry) bool {
		vY, vM, vD := time.Time(v.Day).Date()
		return vY == y && vM == m && vD == d
	})
}

func ValueEntriesGetNewest(ves []ValueEntry) (*ValueEntry, error) {
	if len(ves) == 0 {
		return nil, errors.New("Values are empty")
	} else {
		var max = ves[0]
		for i, ve := range ves {
			if i != 0 {
				maxY, maxM, maxD := time.Time(max.Day).Date()
				y, m, d := time.Time(ve.Day).Date()
				if y > maxY {
					max = ve
				} else if y == maxY && m > maxM {
					max = ve
				} else if y == maxY && m == maxM && d > maxD {
					max = ve
				}
			}
		}
		return &max, nil
	}
}
