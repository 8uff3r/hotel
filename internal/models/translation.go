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
	Slug        string      `gorm:"uniqueIndex" json:"slug"`
	Label       string      `gorm:"-" json:"label"`
	Translation Translation `gorm:"type:jsonb" json:"-"`
}

func (t TranslateBase) UniqueCondition() any {
	return map[string]any{"slug": t.Slug}
}

type Translatable interface {
	GetTranslation() Translation
	SetLabel(string)
	ClearTranslation()
}

func (t *TranslateBase) GetTranslation() Translation {
	return t.Translation
}

func (t *TranslateBase) SetLabel(name string) {
	t.Label = name
}

func (t *TranslateBase) ClearTranslation() {
	t.Translation = nil
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

func ApplyTranslationOnTranslatable(t Translatable, lang string) {
	translations := t.GetTranslation()
	if translated, exists := translations[lang]; exists {
		t.SetLabel(translated)
	}
	t.ClearTranslation()
}

func ApplyTranslations[T any](items *[]T, lang string) {
	for i := range *items {
		item := &(*items)[i]
		if translatable, ok := any(item).(Translatable); ok {
			ApplyTranslationOnTranslatable(translatable, lang)
		}
	}
}

func ApplyFieldTranslations[T any](items *[]T, lang string) {
	for i := range *items {
		item := &(*items)[i]

		if container, ok := any(item).(HasTranslatables); ok {
			translatables := container.GetTranslatables()

			for _, tr := range translatables {
				ApplyTranslationOnTranslatable(tr, lang)
			}
		}
	}
}
