package migrations

import (
	"time"

	"project/internal/config"
	"project/internal/database"
	"project/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AutoMigrateAndSeed(cfg config.Config) error {
	db := database.DB
	if db == nil {
		return gorm.ErrInvalidDB
	}
	if err := db.AutoMigrate(&models.User{}, &models.Team{}, &models.Player{}, &models.PlayerStat{}, &models.Match{}); err != nil {
		return err
	}
	seedTop6(db)
	ensureAdmin(db, cfg.AdminEmail)
	seedMatches(db)
	return nil
}

func seedTop6(db *gorm.DB) {
	var count int64
	db.Model(&models.Team{}).Count(&count)
	if count > 0 {
		return
	}
	teams := []models.Team{
		{Name: "Manchester City", ShortName: "MCI", PrimaryColor: "#6CABDD", SecondaryColor: "#1C2C5B", LogoURL: "/static/logos/mci.png"},
		{Name: "Arsenal", ShortName: "ARS", PrimaryColor: "#EF0107", SecondaryColor: "#9C824A", LogoURL: "/static/logos/ars.png"},
		{Name: "Liverpool", ShortName: "LIV", PrimaryColor: "#C8102E", SecondaryColor: "#00A398", LogoURL: "/static/logos/liv.png"},
		{Name: "Manchester United", ShortName: "MUN", PrimaryColor: "#DA291C", SecondaryColor: "#FBE122", LogoURL: "/static/logos/mun.png"},
		{Name: "Chelsea", ShortName: "CHE", PrimaryColor: "#034694", SecondaryColor: "#DBA111", LogoURL: "/static/logos/che.png"},
		{Name: "Tottenham", ShortName: "TOT", PrimaryColor: "#132257", SecondaryColor: "#FFFFFF", LogoURL: "/static/logos/tot.png"},
	}
	for _, t := range teams {
		db.Create(&t)
	}
}

func ensureAdmin(db *gorm.DB, email string) {
	var u models.User
	if err := db.Where("email = ?", email).First(&u).Error; err == nil {
		return
	}
	// Default admin credentials: zhalgasandalisher@gmail.com / UnitedNom1!
	pw, _ := bcrypt.GenerateFromPassword([]byte("UnitedNom1!"), bcrypt.DefaultCost)
	admin := models.User{
		Name:         "Admin",
		Email:        "zhalgasandalisher@gmail.com",
		PasswordHash: string(pw),
		Role:         "admin",
	}
	db.Create(&admin)
}

func seedMatches(db *gorm.DB) {
	var teams []models.Team
	db.Find(&teams)
	if len(teams) < 2 {
		return
	}
	now := time.Now().Unix()
	for i := 0; i < len(teams)-1; i++ {
		m := models.Match{
			HomeTeamID: teams[i].ID,
			AwayTeamID: teams[i+1].ID,
			Date:       now + int64(i+1)*86400,
			Stadium:    "Stadium",
			Status:     "upcoming",
		}
		db.Create(&m)
	}
}
