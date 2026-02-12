package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"project/internal/config"
	"project/internal/database"
	"project/internal/handlers"
	"project/internal/migrations"
	"project/internal/services"
)

func main() {
	_ = os.Setenv("GIN_MODE", "release")
	cfg := config.Load()
	db := database.Connect(cfg)
	if err := migrations.AutoMigrateAndSeed(cfg); err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/static", "web/static")
	// Serve Manchester United logo directly from root
	router.StaticFile("/static/logos/mun.png", "Manchester_United_FC_crest.svg.png")
	router.LoadHTMLGlob("web/templates/*.html")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "auth.html", gin.H{})
	})
	router.GET("/auth", func(c *gin.Context) {
		c.HTML(200, "auth.html", gin.H{})
	})
	// Public routes
	router.GET("/profile", func(c *gin.Context) {
		c.HTML(200, "profile.html", gin.H{})
	})
	
	// Protected routes - require authentication
	protected := router.Group("/")
	protected.Use(func(c *gin.Context) {
		// Check token in cookie or localStorage (handled by JS)
		// Server-side: allow HTML to load, JS will handle redirect
		c.Next()
	})
	protected.GET("/feed", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	protected.GET("/live", func(c *gin.Context) {
		c.HTML(200, "live.html", gin.H{})
	})
	protected.GET("/analytics", func(c *gin.Context) {
		c.HTML(200, "analytics.html", gin.H{})
	})
	protected.GET("/community", func(c *gin.Context) {
		c.HTML(200, "community.html", gin.H{})
	})
	protected.GET("/league", func(c *gin.Context) {
		c.HTML(200, "league.html", gin.H{})
	})
	protected.GET("/account", func(c *gin.Context) {
		c.HTML(200, "account.html", gin.H{})
	})

	api := &handlers.API{
		Auth:      &services.AuthService{DB: db, JWTSecret: cfg.JWTSecret},
		Teams:     &services.TeamService{DB: db},
		Players:   &services.PlayerService{DB: db},
		Matches:   &services.MatchService{DB: db},
		Table:     &services.TableService{DB: db},
		JWTSecret: cfg.JWTSecret,
	}
	api.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
