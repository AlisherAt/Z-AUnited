package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"size:100"`
	Email          string `gorm:"size:180;uniqueIndex"`
	PasswordHash   string `gorm:"size:255"`
	Role           string `gorm:"size:20;default:user"`
	FavoriteTeamID *uint
	FavoriteTeam   *Team
}

type Team struct {
	gorm.Model
	Name          string `gorm:"size:100;uniqueIndex"`
	ShortName     string `gorm:"size:20"`
	LogoURL       string `gorm:"size:255"`
	PrimaryColor  string `gorm:"size:20"`
	SecondaryColor string `gorm:"size:20"`
	Points        int  `gorm:"default:0"`
	MatchesPlayed int  `gorm:"default:0"`
	GoalDiff      int  `gorm:"default:0"`
	Players       []Player
}

type Player struct {
	gorm.Model
	Name    string `gorm:"size:120"`
	TeamID  uint
	Team    Team
	Position string `gorm:"size:30"`
	Stats   []PlayerStat
}

type PlayerStat struct {
	gorm.Model
	PlayerID     uint
	Player       Player
	Season       string `gorm:"size:10"`
	Goals        int    `gorm:"default:0"`
	Assists      int    `gorm:"default:0"`
	CleanSheets  int    `gorm:"default:0"`
	MinutesPlayed int   `gorm:"default:0"`
}

type Match struct {
	gorm.Model
	HomeTeamID uint
	AwayTeamID uint
	HomeTeam   Team
	AwayTeam   Team
	HomeScore  *int
	AwayScore  *int
	Date       int64
	Stadium    string `gorm:"size:120"`
	Status     string `gorm:"size:20"` // upcoming | finished | live
}
