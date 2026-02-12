package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LiveMatchTrackerHandler streams real-time match data (stub for now)
func LiveMatchTrackerHandler(c *gin.Context) {
	// TODO: Replace with real-time data source (WebSocket or SSE)
	c.JSON(http.StatusOK, gin.H{
		"matchId":  1,
		"homeTeam": "Arsenal",
		"awayTeam": "Chelsea",
		"score":    "2-1",
		"minute":   67,
		"commentary": []string{
			"67' GOAL! Arsenal take the lead!",
			"65' Substitution: Smith Rowe on for Saka.",
			"60' Yellow card for Chelsea.",
		},
		"playerStats": []gin.H{
			{"name": "Bukayo Saka", "goals": 1, "assists": 0, "minutes": 67},
			{"name": "Kai Havertz", "goals": 1, "assists": 0, "minutes": 67},
		},
	})
}
