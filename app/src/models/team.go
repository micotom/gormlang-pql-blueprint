package models

type Team struct {
	Slug    string   `json:"slug" gorm:"primaryKey"`
	Players []Player `json:"players" gorm:"many2many:player_teams;"`
}

func TeamTotalRaise(t Team) int {
	var sum = 0
	for _, p := range t.Players {
		sum += PlayerGetCurrentRaiseDiff(p)
	}
	return sum
}

func TeamCurrentValue(t Team) int {
	var sum = 0
	for _, p := range t.Players {
		sum += PlayerGetCurrentValue(p)
	}
	return sum
}
