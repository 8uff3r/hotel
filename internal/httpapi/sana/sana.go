package sana

import (
	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	h "github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type SanaModule struct {
	db *gorm.DB
}

func New(db *gorm.DB) SanaModule {
	return SanaModule{db: db}
}

func (m SanaModule) RegisterRoutes(api *httpapi.API, s *fuego.Server) {
	h.Get(s, "/travel-reasons", m.getTravelReasons)
	h.Get(s, "/family-relationships", m.getFamilyRelationships)
	h.Get(s, "/nationalities", m.getNationalities)
	h.Get(s, "/countries", m.getCountries)
	h.Get(s, "/occupations", m.getOccupations)
}

func (m SanaModule) getTravelReasons(c h.ContextNoBody) ([]models.TravelReason, error) {
	var reasons []models.TravelReason
	err := m.db.Find(&reasons).Error
	return reasons, err
}

func (m SanaModule) getFamilyRelationships(c h.ContextNoBody) ([]models.FamilyRelationship, error) {
	var relationships []models.FamilyRelationship
	err := m.db.Find(&relationships).Error
	return relationships, err
}

func (m SanaModule) getNationalities(c h.ContextNoBody) ([]models.Nationality, error) {
	var nationalities []models.Nationality
	err := m.db.Find(&nationalities).Error
	return nationalities, err
}

func (m SanaModule) getCountries(c h.ContextNoBody) ([]models.Country, error) {
	var countries []models.Country
	err := m.db.Find(&countries).Error
	return countries, err
}

func (m SanaModule) getOccupations(c h.ContextNoBody) ([]models.Occupation, error) {
	var occupations []models.Occupation
	err := m.db.Find(&occupations).Error
	return occupations, err
}
