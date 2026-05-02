package sana

import (
	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type SanaModule struct {
	db *gorm.DB
}

func New(db *gorm.DB) SanaModule {
	return SanaModule{db: db}
}

func (m SanaModule) RegisterRoutes(api *httpapi.API, s *fuego.Server) {
	fuego.Get(s, "/travel-reasons", m.getTravelReasons)
	fuego.Get(s, "/family-relationships", m.getFamilyRelationships)
	fuego.Get(s, "/nationalities", m.getNationalities)
	fuego.Get(s, "/countries", m.getCountries)
	fuego.Get(s, "/occupations", m.getOccupations)
}

func (m SanaModule) getTravelReasons(c fuego.ContextNoBody) ([]models.TravelReason, error) {
	var reasons []models.TravelReason
	err := m.db.Find(&reasons).Error
	return reasons, err
}

func (m SanaModule) getFamilyRelationships(c fuego.ContextNoBody) ([]models.FamilyRelationship, error) {
	var relationships []models.FamilyRelationship
	err := m.db.Find(&relationships).Error
	return relationships, err
}

func (m SanaModule) getNationalities(c fuego.ContextNoBody) ([]models.Nationality, error) {
	var nationalities []models.Nationality
	err := m.db.Find(&nationalities).Error
	return nationalities, err
}

func (m SanaModule) getCountries(c fuego.ContextNoBody) ([]models.Country, error) {
	var countries []models.Country
	err := m.db.Find(&countries).Error
	return countries, err
}

func (m SanaModule) getOccupations(c fuego.ContextNoBody) ([]models.Occupation, error) {
	var occupations []models.Occupation
	err := m.db.Find(&occupations).Error
	return occupations, err
}
