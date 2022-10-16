package models

import (
	"errors"
	"time"
)

type Player struct {
	Slug         string `json:"slug" gorm:"primaryKey"`
	Name         string `json:"name"`
	Club         string `json:"club"`
	Position     string `json:"position"`
	TotalPoints  int    `json:"total_points"`
	ValueEntries []ValueEntry
}

func GetPlayerPointsPerPrice(p Player) float32 {
	v := PlayerGetCurrentValue(p)
	ps := p.TotalPoints
	if ps == 0 {
		return -1
	}
	return float32(v) / float32(ps)
}

func PlayerGetCurrentRaisePerc(p Player) float32 {
	y, m, d := time.Now().Date()
	if v, e := PlayerGetValueAt(p, y, m, d); e == nil {
		return v.RaisePerc
	} else {
		return 0
	}
}

func PlayerGetCurrentRaiseDiff(p Player) int {
	y, m, d := time.Now().Date()
	if v, e := PlayerGetValueAt(p, y, m, d); e == nil {
		return v.RaiseDiff
	} else {
		return 0
	}
}

func PlayerGetCurrentValue(p Player) int {
	y, m, d := time.Now().Date()
	if v, e := PlayerGetValueAt(p, y, m, d); e == nil {
		return v.Value
	} else {
		// no current value, search for the last one
		if v, e := ValueEntriesGetNewest(p.ValueEntries); e == nil {
			return (*v).Value
		}
	}
	return -1
}

func PlayerGetValueAt(p Player, y int, m time.Month, d int) (*ValueEntry, error) {
	for _, entry := range p.ValueEntries {
		eY, eM, eD := time.Time(entry.Day).Date()
		if eY == y && eM == m && eD == d {
			return &entry, nil
		}
	}
	return nil, errors.New("No entry for day")
}

func PlayerValueTrend(p Player) string {
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	if tE, e := PlayerGetValueAt(p, today.Year(), today.Month(), today.Day()); e == nil {
		if yE, e := PlayerGetValueAt(p, yesterday.Year(), yesterday.Month(), yesterday.Day()); e == nil {
			if tE.RaiseDiff > yE.RaiseDiff {
				return "+"
			} else if tE.RaiseDiff < yE.RaiseDiff {
				return "-"
			} else {
				return "="
			}
		}
	}
	return "="
}
