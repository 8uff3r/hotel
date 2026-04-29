package permissions

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"strconv"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type PermissionsModule struct {
	*h.API
}

type permissionsResponse struct {
	Data []models.Permission `json:"data"`
}

type templatesResponse struct {
	Data []models.PermissionTemplate `json:"data"`
}

func (m PermissionsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	pm := PermissionsModule{api}

	fuego.Get(s, "/", pm.permissionsList)
	fuego.Get(s, "/templates", pm.templatesList)
	fuego.Get(s, "/user/{userId}", pm.userPermissions)
	fuego.Post(s, "/user/{userId}/{permissionId}", pm.setUserPermission)
	fuego.Post(s, "/user/{userId}/template/{templateId}", pm.applyTemplate)
	fuego.Delete(s, "/user/{userId}/permission/{permissionId}", pm.removeUserPermission)
}

func (pm *PermissionsModule) permissionsList(c fuego.ContextNoBody) (permissionsResponse, error) {
	var perms []models.Permission
	var zero permissionsResponse

	sortColumn := c.QueryParam("sort")
	order := "category, page, action"
	if sortColumn == "page" {
		order = "page, category, action"
	} else if sortColumn == "action" {
		order = "action, page, category"
	}

	if err := pm.Db.WithContext(c).Order(order).Find(&perms).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	return permissionsResponse{Data: perms}, nil
}

func (pm *PermissionsModule) templatesList(c fuego.ContextNoBody) (templatesResponse, error) {
	var tpls []models.PermissionTemplate
	var zero templatesResponse

	if err := pm.Db.WithContext(c).Preload("Permissions").Find(&tpls).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	return templatesResponse{Data: tpls}, nil
}

type userPermissionsResponse struct {
	UserID      uint                        `json:"userId"`
	Permissions []models.UserPermissionInfo `json:"permissions"`
}

func (pm *PermissionsModule) userPermissions(c fuego.ContextNoBody) (userPermissionsResponse, error) {
	userID := c.PathParam("userId")

	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	var zero userPermissionsResponse

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_user_id"}
	}

	var userPerms []models.UserPermission
	if err := pm.Db.WithContext(c).Preload("Permission").Where("user_id = ?", uid).Find(&userPerms).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	result := make([]models.UserPermissionInfo, 0, len(userPerms))
	for _, up := range userPerms {
		models.ApplyTranslationOnTranslatable(&up.Permission, lang)
		result = append(result, models.UserPermissionInfo{
			PermissionID: up.PermissionID,
			Page:         up.Permission.Page,
			Action:       up.Permission.Action,
			Label:        up.Permission.Name,
			Category:     up.Permission.Category,
			Granted:      up.Granted,
		})
	}

	return userPermissionsResponse{UserID: uint(uid), Permissions: result}, nil
}

type setPermissionDto struct {
	Granted bool `json:"granted"`
}

func (pm *PermissionsModule) setUserPermission(c fuego.ContextWithBody[setPermissionDto]) (map[string]string, error) {
	userID := c.PathParam("userId")
	permissionID := c.PathParam("permissionId")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_user_id"}
	}
	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_permission_id"}
	}

	body, err := c.Body()
	if err != nil {
		return nil, fuego.BadRequestError{}
	}

	var existing models.UserPermission
	result := pm.Db.Where("user_id = ? AND permission_id = ?", uid, pid).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		if body.Granted {
			newPerm := models.UserPermission{
				UserID:       uint(uid),
				PermissionID: uint(pid),
				Granted:      true,
			}
			if err := pm.Db.WithContext(c).Create(&newPerm).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "create_failed"}
			}
		}
	} else if result.Error == nil {
		if !body.Granted {
			if err := pm.Db.Delete(&existing).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "delete_failed"}
			}
		} else if !existing.Granted {
			existing.Granted = true
			if err := pm.Db.Save(&existing).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "update_failed"}
			}
		}
	} else {
		return nil, fuego.InternalServerError{Title: "query_failed"}
	}

	return map[string]string{"status": "ok"}, nil
}

func (pm *PermissionsModule) applyTemplate(c fuego.ContextNoBody) (map[string]string, error) {
	userID := c.PathParam("userId")
	templateID := c.PathParam("templateId")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_user_id"}
	}
	tid, err := strconv.ParseUint(templateID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_template_id"}
	}

	var template models.PermissionTemplate
	if err := pm.Db.WithContext(c).Preload("Permissions").First(&template, tid).Error; err != nil {
		return nil, fuego.BadRequestError{Title: "template_not_found"}
	}

	var existing []models.UserPermission
	pm.Db.Where("user_id = ?", uid).Find(&existing)
	existingMap := make(map[uint]bool)
	for _, e := range existing {
		existingMap[e.PermissionID] = e.Granted
	}

	for _, perm := range template.Permissions {
		if _, exists := existingMap[perm.ID]; !exists {
			newUP := models.UserPermission{
				UserID:       uint(uid),
				PermissionID: perm.ID,
				Granted:      true,
			}
			if err := pm.Db.Create(&newUP).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "create_failed"}
			}
		} else if !existingMap[perm.ID] {
			pm.Db.Model(&models.UserPermission{}).Where("user_id = ? AND permission_id = ?", uid, perm.ID).Update("granted", true)
		}
	}

	return map[string]string{"status": "ok"}, nil
}

func (pm *PermissionsModule) removeUserPermission(c fuego.ContextNoBody) (map[string]string, error) {
	userID := c.PathParam("userId")
	permissionID := c.PathParam("permissionId")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_user_id"}
	}
	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_permission_id"}
	}

	if err := pm.Db.Where("user_id = ? AND permission_id = ?", uid, pid).Delete(&models.UserPermission{}).Error; err != nil {
		return nil, fuego.InternalServerError{Title: "delete_failed"}
	}

	return map[string]string{"status": "ok"}, nil
}

type allPermissionsResponse struct {
	Permissions []models.Permission         `json:"permissions"`
	Templates   []models.PermissionTemplate `json:"templates"`
}

func (pm *PermissionsModule) allPermissionsAndTemplates(c fuego.ContextNoBody) (allPermissionsResponse, error) {
	var zero allPermissionsResponse

	var perms []models.Permission
	if err := pm.Db.Order("category, page, action").Find(&perms).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	var tpls []models.PermissionTemplate
	if err := pm.Db.Preload("Permissions").Find(&tpls).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	return allPermissionsResponse{Permissions: perms, Templates: tpls}, nil
}
