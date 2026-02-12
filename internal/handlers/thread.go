package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Thread struct {
	ID       int       `json:"id"`
	MatchID  int       `json:"matchId"`
	Title    string    `json:"title"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// In-memory stub for threads
var threads = []Thread{
	{
		ID: 1, MatchID: 1, Title: "Arsenal vs Chelsea Match Thread",
		Comments: []Comment{
			{User: "fan1", Message: "What a goal!", Time: "2026-02-10T18:10:00Z"},
			{User: "fan2", Message: "VAR check incoming...", Time: "2026-02-10T18:12:00Z"},
		},
	},
}

// ListMatchThreads returns all match threads
func ListMatchThreads(c *gin.Context) {
	c.JSON(http.StatusOK, threads)
}

// PostComment adds a comment to a thread (stub, no auth)
func PostComment(c *gin.Context) {
	var req struct {
		ThreadID int    `json:"threadId"`
		User     string `json:"user"`
		Message  string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	for i := range threads {
		if threads[i].ID == req.ThreadID {
			threads[i].Comments = append(threads[i].Comments, Comment{User: req.User, Message: req.Message, Time: "2026-02-10T18:20:00Z"})
			c.JSON(http.StatusOK, threads[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
}
