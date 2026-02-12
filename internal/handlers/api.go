package handlers

import (
	"net/http"
	"strconv"

	"project/internal/middleware"
	"project/internal/models"
	"project/internal/services"

	"github.com/gin-gonic/gin"
)

type API struct {
	Auth      *services.AuthService
	Teams     *services.TeamService
	Players   *services.PlayerService
	Matches   *services.MatchService
	Table     *services.TableService
	JWTSecret string
}

func (a *API) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/table", a.getTable)
	api.GET("/teams", a.getTeams)
	api.GET("/players", a.getPlayers)
	api.GET("/matches", a.getMatches)
	api.GET("/matchtracker", LiveMatchTrackerHandler)
	api.GET("/threads", ListMatchThreads)
	api.POST("/threads/comment", PostComment)
	api.GET("/stats", StatsHandler)
	api.GET("/historical", HistoricalDataHandler)

	api.POST("/auth/register", a.register)
	api.POST("/auth/login", a.login)
	api.POST("/auth/logout", a.logout)

	auth := api.Group("/")
	auth.Use(middleware.Auth(a.JWTSecret))
	auth.POST("/profile/favorite", a.setFavoriteTeam)
	auth.GET("/profile/me", a.me)
	auth.GET("/feed", PersonalizedFeedHandler)

	admin := auth.Group("/admin")
	admin.Use(middleware.RequireAdmin())
	admin.POST("/teams", a.upsertTeam)
	admin.POST("/players", a.upsertPlayer)
	admin.POST("/matches/:id/result", a.updateMatchResult)
}

func (a *API) register(c *gin.Context) {
	var body struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		FavoriteTeam *uint  `json:"favoriteTeam"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	u, err := a.Auth.Register(body.Name, body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Set favorite team if provided during registration
	if body.FavoriteTeam != nil && *body.FavoriteTeam > 0 {
		services.DB().Model(u).Update("favorite_team_id", *body.FavoriteTeam)
	}
	c.JSON(http.StatusOK, gin.H{"user": u})
}

func (a *API) login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	token, u, err := a.Auth.Login(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// Set token in HTTP-only cookie for better security
	c.SetCookie("auth_token", token, 86400, "/", "", false, true) // 24 hours, httpOnly
	c.JSON(http.StatusOK, gin.H{"token": token, "user": u})
}

func (a *API) logout(c *gin.Context) {
	// Clear auth cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) getTeams(c *gin.Context) {
	list, err := a.Teams.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) upsertTeam(c *gin.Context) {
	var t models.Team
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := a.Teams.Upsert(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (a *API) getPlayers(c *gin.Context) {
	teamStr := c.Query("teamId")
	var teamID uint
	if teamStr != "" {
		if v, err := strconv.Atoi(teamStr); err == nil {
			teamID = uint(v)
		}
	}
	list, err := a.Players.List(teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) upsertPlayer(c *gin.Context) {
	var p models.Player
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := a.Players.Upsert(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (a *API) getMatches(c *gin.Context) {
	list, err := a.Matches.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (a *API) updateMatchResult(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Home   int    `json:"home"`
		Away   int    `json:"away"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	if err := a.Matches.UpdateResult(uint(id64), body.Home, body.Away, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) getTable(c *gin.Context) {
	rows, err := a.Table.Compute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (a *API) setFavoriteTeam(c *gin.Context) {
	var body struct {
		TeamID uint `json:"teamId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	uidVal, _ := c.Get("uid")
	uid := uidVal.(uint)
	if err := services.DB().Model(&models.User{}).Where("id = ?", uid).Update("favorite_team_id", body.TeamID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *API) me(c *gin.Context) {
	uidVal, _ := c.Get("uid")
	uid := uidVal.(uint)
	var u models.User
	if err := services.DB().Preload("FavoriteTeam").First(&u, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}
