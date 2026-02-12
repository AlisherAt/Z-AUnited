package services

import (
	"errors"
	"strings"
	"time"

	"project/internal/cache"
	"project/internal/database"
	"project/internal/middleware"
	"project/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret string
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" || name == "" {
		return nil, errors.New("missing fields")
	}
	var exists models.User
	if err := s.DB.Where("email = ?", email).First(&exists).Error; err == nil {
		return nil, errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	var u models.User
	if err := s.DB.Where("email = ?", strings.ToLower(email)).First(&u).Error; err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token, err := middleware.GenerateToken(s.JWTSecret, u.ID, u.Role, 24*time.Hour)
	return token, &u, err
}

type TeamService struct{ DB *gorm.DB }

func (s *TeamService) List() ([]models.Team, error) {
	var teams []models.Team
	err := s.DB.Find(&teams).Error
	return teams, err
}

func (s *TeamService) Upsert(t *models.Team) error {
	return s.DB.Save(t).Error
}

type PlayerService struct{ DB *gorm.DB }

func (s *PlayerService) List(teamID uint) ([]models.Player, error) {
	var p []models.Player
	q := s.DB
	if teamID != 0 {
		q = q.Where("team_id = ?", teamID)
	}
	err := q.Preload("Stats").Find(&p).Error
	return p, err
}

func (s *PlayerService) Upsert(p *models.Player) error {
	return s.DB.Save(p).Error
}

type MatchService struct{ DB *gorm.DB }

func (s *MatchService) List() ([]models.Match, error) {
	var m []models.Match
	err := s.DB.Preload("HomeTeam").Preload("AwayTeam").Find(&m).Error
	return m, err
}

func (s *MatchService) UpdateResult(id uint, home, away int, status string) error {
	return s.DB.Model(&models.Match{}).Where("id = ?", id).
		Updates(map[string]interface{}{"home_score": home, "away_score": away, "status": status}).Error
}

type TableRow struct {
	TeamID   uint   `json:"team_id"`
	Team     string `json:"team"`
	Played   int    `json:"played"`
	Points   int    `json:"points"`
	GoalDiff int    `json:"gd"`
}

type TableService struct{ DB *gorm.DB }

var tableCache = cache.New()

func (s *TableService) Compute() ([]TableRow, error) {
	var rows []TableRow
	// Try Redis first
	found, err := cache.GetRedis("league_table", &rows)
	if err == nil && found {
		return rows, nil
	}
	// Fallback to in-memory
	if v, ok := tableCache.Get("league_table"); ok {
		if cachedRows, ok2 := v.([]TableRow); ok2 {
			return cachedRows, nil
		}
	}
	var teams []models.Team
	if err := s.DB.Order("points desc, goal_diff desc").Find(&teams).Error; err != nil {
		return nil, err
	}
	rows = make([]TableRow, 0, len(teams))
	for _, t := range teams {
		rows = append(rows, TableRow{
			TeamID:   t.ID,
			Team:     t.Name,
			Played:   t.MatchesPlayed,
			Points:   t.Points,
			GoalDiff: t.GoalDiff,
		})
	}
	// Set both Redis and in-memory cache
	_ = cache.SetRedis("league_table", rows, 30*time.Second)
	tableCache.Set("league_table", rows, 30*time.Second)
	return rows, nil
}

func DB() *gorm.DB { return database.DB }
