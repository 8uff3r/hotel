package seed

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"hotel/internal/models"

	"gorm.io/gorm"
)

type (
	permissionActions map[string]string
	permissionPages   map[string]permissionActions
	permissionsJSON   map[string]permissionPages
)

//go:embed permissions.json
var permissionsFile []byte

var Permissions permissionsJSON

func init() {
	err := json.Unmarshal(permissionsFile, &Permissions)
	if err != nil {
		log.Fatalf("failed to load permissions: %v", err)
	}

	distFile := "frontend/app/utils/permissions.gen.ts"
	content := "export const PERMISSIONS = " + strings.TrimSpace(string(permissionsFile)) + ";"

	dir := filepath.Dir(distFile)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(distFile, []byte(content), 0o644); err != nil {
		log.Fatalf("failed to write permissions gen file: %v", err)
	}

	cmd := exec.Command("bun", "fmt", strings.Replace(distFile, "frontend/", "", 1))
	cmd.Dir = "frontend"
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("warning: bun fmt failed: %v, output: %s", err, out)
	}
}

func seedPermissions(db *gorm.DB) {
	t := Translations["permission"]
	tc := Translations["permission-category"]

	var categories []models.PermissionCategory
	for categorySlug := range Permissions {
		categories = append(categories, models.PermissionCategory{
			TranslateBase: models.TranslateBase{
				Slug:        categorySlug,
				Translation: tc[categorySlug],
			},
		})
	}
	seed(db, categories)

	var seededCategories []models.PermissionCategory
	db.Find(&seededCategories)
	categoryMap := make(map[string]uint) // slug -> ID
	for _, c := range seededCategories {
		categoryMap[c.Slug] = c.ID
	}

	var perms []models.Permission
	for categorySlug, pages := range Permissions {
		for _, actions := range pages {
			for action, pageAction := range actions {
				parts := strings.SplitN(pageAction, ":", 2)
				page := parts[0]
				perms = append(perms, models.Permission{
					Translation: t[pageAction],
					Resource:    page,
					Action:      models.PermissionAction(action),
					CategoryID:  categoryMap[categorySlug],
				})
			}
		}
	}
	seed(db, perms)
}

func seedPermissionTemplates(db *gorm.DB) {
	t := Translations["permission-template"]

	var allPerms []models.Permission
	db.Find(&allPerms)

	permMap := make(map[string]uint)
	for _, p := range allPerms {
		key := p.Resource + ":" + string(p.Action)
		permMap[key] = p.ID
	}

	var adminTemplate models.PermissionTemplate
	if err := db.Where("slug = ?", "admin").First(&adminTemplate).Error; err == gorm.ErrRecordNotFound {
		adminTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Slug: "admin", Translation: t["admin"]},
		}
		db.Create(&adminTemplate)
	}
	db.Model(&adminTemplate).Association("Permissions").Append(&allPerms)

	var managerTemplate models.PermissionTemplate
	if err := db.Where("slug = ?", "manager").First(&managerTemplate).Error; err == gorm.ErrRecordNotFound {
		managerTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Slug: "manager", Translation: t["manager"]},
		}
		db.Create(&managerTemplate)
	}
	var managerPerms []models.Permission
	for _, p := range allPerms {
		if p.Resource != "users" {
			managerPerms = append(managerPerms, p)
		}
	}
	db.Model(&managerTemplate).Association("Permissions").Replace(&managerPerms)

	var receptionistTemplate models.PermissionTemplate
	if err := db.Where("slug = ?", "receptionist").First(&receptionistTemplate).Error; err == gorm.ErrRecordNotFound {
		receptionistTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Slug: "receptionist", Translation: t["receptionist"]},
		}
		db.Create(&receptionistTemplate)
	}
	var receptionistPerms []models.Permission
	for _, p := range allPerms {
		if p.Resource == "index" ||
			p.Resource == "guests" || p.Resource == "guests/settle" ||
			p.Resource == "rooms" || p.Resource == "rooms/rack" ||
			p.Resource == "reservations" ||
			p.Resource == "parking" || p.Resource == "parking/transactions" || p.Resource == "parking/transactions/check-in" ||
			p.Resource == "parking/vehicles" {
			receptionistPerms = append(receptionistPerms, p)
		}
	}
	db.Model(&receptionistTemplate).Association("Permissions").Replace(&receptionistPerms)

	var staffTemplate models.PermissionTemplate
	if err := db.Where("slug = ?", "staff").First(&staffTemplate).Error; err == gorm.ErrRecordNotFound {
		staffTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Slug: "staff", Translation: t["staff"]},
		}
		db.Create(&staffTemplate)
	}
	var staffPerms []models.Permission
	for _, p := range allPerms {
		if p.Action == models.PermissionActionRead {
			staffPerms = append(staffPerms, p)
		}
	}
	db.Model(&staffTemplate).Association("Permissions").Replace(&staffPerms)

	var housekeeperTemplate models.PermissionTemplate
	if err := db.Where("slug = ?", "housekeeper").First(&housekeeperTemplate).Error; err == gorm.ErrRecordNotFound {
		housekeeperTemplate = models.PermissionTemplate{
			TranslateBase: models.TranslateBase{Slug: "housekeeper", Translation: t["housekeeper"]},
		}
		db.Create(&housekeeperTemplate)
	}
	var housekeeperPerms []models.Permission
	for _, p := range allPerms {
		if p.Resource == "rooms" || p.Resource == "rooms/rack" {
			housekeeperPerms = append(housekeeperPerms, p)
		}
	}
	db.Model(&housekeeperTemplate).Association("Permissions").Replace(&housekeeperPerms)
}
