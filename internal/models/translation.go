package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Translation map[string]string
type HasTranslatables interface {
	GetTranslatables() []Translatable
}
type TranslateBase struct {
	Name        string      `json:"name"`
	Translation Translation `gorm:"type:jsonb" json:"translation,omitempty"`
}

type Translatable interface {
	GetTranslation() Translation
	SetName(string)
	ClearTranslation()
}

func (t *TranslateBase) GetTranslation() Translation {
	return t.Translation
}

func (t *TranslateBase) SetName(name string) {
	t.Name = name
}

func (t *TranslateBase) ClearTranslation() {
	t.Translation = nil // or zero value, depending on your Translation type
	// If Translation is a map, set to nil or make empty
}

func (t *Translation) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal Translation: %v", value)
	}
	return json.Unmarshal(bytes, t)
}

func (t Translation) Value() (driver.Value, error) {
	return json.Marshal(t)
}
