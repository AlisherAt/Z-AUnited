package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatsHandler returns player and team stats (stub)
func StatsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"topScorers": []gin.H{
			{"player": "Erling Haaland", "team": "Manchester City", "goals": 18},
			{"player": "Mohamed Salah", "team": "Liverpool", "goals": 15},
		},
		"topAssisters": []gin.H{
			{"player": "Kevin De Bruyne", "team": "Manchester City", "assists": 12},
			{"player": "Martin Ødegaard", "team": "Arsenal", "assists": 10},
		},
		"cleanSheets": []gin.H{
			{"player": "Alisson Becker", "team": "Liverpool", "cleanSheets": 11},
			{"player": "Ederson", "team": "Manchester City", "cleanSheets": 10},
		},
		"teamStandings": []gin.H{
			{"team": "Arsenal", "points": 55},
			{"team": "Manchester City", "points": 54},
		},
	})
}

// HistoricalDataHandler returns historical EPL data (stub)
func HistoricalDataHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"seasons": []gin.H{
			{"season": "2024/25", "winner": "Manchester City"},
			{"season": "2023/24", "winner": "Arsenal"},
		},
	})
}
