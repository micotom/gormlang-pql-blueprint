package models

type Player struct {
	Slug         string `json:"slug" gorm:"primaryKey"`
	Name         string `json:"name"`
	Club         string `json:"club"`
	Position     string `json:"position"`
	TotalPoints  int    `json:"total_points"`
	ValueEntries []ValueEntry
}
