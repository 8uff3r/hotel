package seed

import (
	"fmt"
	"hotel/internal/config"
	"hotel/internal/models"
	"hotel/internal/sana"

	"gorm.io/gorm"
)

func seedSanaReferenceData(db *gorm.DB, cfg config.Config) {
	if cfg.Sana.KelidVahed == "" || cfg.Sana.KelidPeimankar == "" || cfg.Sana.CodeVahed == 0 {
		fmt.Println("Sana config not provided, skipping Sana reference data seeding")
		return
	}

	client := sana.NewClient(sana.Config{
		KelidVahed:     cfg.Sana.KelidVahed,
		KelidPeimankar: cfg.Sana.KelidPeimankar,
		CodeVahed:      cfg.Sana.CodeVahed,
	})

	seedTravelReasons(db, client)
	seedFamilyRelationships(db, client)
	seedNationalities(db, client)
	seedIranianCities(db, client)
	seedForeignCities(db, client)
	seedOccupations(db, client)
}

func seedTravelReasons(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinElaatSafar()
	if err != nil {
		fmt.Printf("Failed to fetch travel reasons: %v\n", err)
		return
	}

	var reasons []models.TravelReason
	for _, item := range items {
		reasons = append(reasons, models.TravelReason{
			SanaID:   item.ID,
			SanaName: item.Name,
		})
	}

	if len(reasons) > 0 {
		db.Where("1=1").Delete(&models.TravelReason{})
		db.CreateInBatches(reasons, 100)
	}
}

func seedFamilyRelationships(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinNesbat()
	if err != nil {
		fmt.Printf("Failed to fetch family relationships: %v\n", err)
		return
	}

	var relationships []models.FamilyRelationship
	for _, item := range items {
		relationships = append(relationships, models.FamilyRelationship{
			SanaID:   item.ID,
			SanaName: item.Name,
		})
	}

	if len(relationships) > 0 {
		db.Where("1=1").Delete(&models.FamilyRelationship{})
		db.CreateInBatches(relationships, 100)
	}
}

func seedNationalities(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinMeliat()
	if err != nil {
		fmt.Printf("Failed to fetch nationalities: %v\n", err)
		return
	}

	var nationalities []models.Nationality
	for _, item := range items {
		nationalities = append(nationalities, models.Nationality{
			SanaID:   item.ID,
			SanaName: item.Name,
		})
	}

	if len(nationalities) > 0 {
		db.Where("1=1").Delete(&models.Nationality{})
		db.CreateInBatches(nationalities, 100)
	}
}

func seedIranianCities(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinShahrhayeIran()
	if err != nil {
		fmt.Printf("Failed to fetch Iranian cities: %v\n", err)
		return
	}

	var cities []models.Country
	for _, item := range items {
		cities = append(cities, models.Country{
			SanaID:   item.ID,
			SanaName: item.Name,
			IsIran:   true,
		})
	}

	if len(cities) > 0 {
		db.Where("is_iran = ?", true).Delete(&models.Country{})
		db.CreateInBatches(cities, 100)
	}
}

func seedForeignCities(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinShahrhayeKhareji()
	if err != nil {
		fmt.Printf("Failed to fetch foreign cities: %v\n", err)
		return
	}

	var cities []models.Country
	for _, item := range items {
		cities = append(cities, models.Country{
			SanaID:   item.ID,
			SanaName: item.Name,
			IsIran:   false,
		})
	}

	if len(cities) > 0 {
		db.Where("is_iran = ?", false).Delete(&models.Country{})
		db.CreateInBatches(cities, 100)
	}
}

func seedOccupations(db *gorm.DB, client *sana.Client) {
	items, err := client.GereftanAnavinShoghl()
	if err != nil {
		fmt.Printf("Failed to fetch occupations: %v\n", err)
		return
	}

	var occupations []models.Occupation
	for _, item := range items {
		occupations = append(occupations, models.Occupation{
			SanaID:   item.ID,
			SanaName: item.Name,
		})
	}

	if len(occupations) > 0 {
		db.Where("1=1").Delete(&models.Occupation{})
		db.CreateInBatches(occupations, 100)
	}
}
