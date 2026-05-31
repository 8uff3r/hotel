package permissions

import (
	"strconv"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type PermissionsModule struct {
	*h.API
}

type permissionsResponse struct {
	Data []models.Permission `json:"data"`
}

func (pm PermissionsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	pm = PermissionsModule{api}

	fuego.Get(s, "/", pm.permissionsList)
	fuego.Get(
		s,
		"/templates",
		h.ListModel[models.PermissionTemplate](
			api.Db,
			h.WithTranslation[models.PermissionTemplate](),
		),
	)
	fuego.Get(s, "/user/{userId}", pm.userPermissions)
	fuego.Post(s, "/user/{userId}/template/{templateId}", pm.applyTemplate)
	fuego.Post(s, "/user/{userId}/{permissionId}", pm.setUserPermission)
	fuego.Post(s, "/user/{userId}/grant-role", pm.grantPermissionsOfTemplateToUser)
	fuego.Post(s, "/user/{userId}/grant-all", pm.grantAllPermissionsToUser)
	fuego.Delete(s, "/user/{userId}/{permissionId}", pm.setUserPermission)
}

func (pm *PermissionsModule) permissionsList(c fuego.ContextNoBody) (permissionsResponse, error) {
	var perms []models.Permission
	var zero permissionsResponse
	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}

	sortColumn := c.QueryParam("sort")
	order := "category_id, resource, action"
	switch sortColumn {
	case "page":
		order = "resource, category_id, action"
	case "action":
		order = "action, resource, category_id"
	}

	if err := pm.Db.WithContext(c).Preload("Category").Order(order).Find(&perms).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	for i := range perms {
		perms[i].Resource = perms[i].Translation[lang]
		perms[i].Category.Label = perms[i].Category.Translation[lang]
	}

	return permissionsResponse{Data: perms}, nil
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
		result = append(result, models.UserPermissionInfo{
			PermissionID: up.PermissionID,
			Page:         up.Permission.Resource,
			Action:       up.Permission.Action,
			Label:        up.Permission.Translation[lang],
			Category:     up.Permission.Category,
			Granted:      up.Granted,
		})
	}

	return userPermissionsResponse{UserID: uint(uid), Permissions: result}, nil
}

type okResponse struct {
	Ok bool `json:"ok"`
}

func (pm *PermissionsModule) setUserPermission(c fuego.ContextNoBody) (*okResponse, error) {
	userID := c.PathParam("userId")
	permissionID := c.PathParam("permissionId")
	shouldAdd := c.Request().Method == "POST"

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_user_id"}
	}
	pid, err := strconv.ParseUint(permissionID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_permission_id"}
	}

	if err != nil {
		return nil, fuego.BadRequestError{}
	}

	var existing models.UserPermission
	result := pm.Db.Where("user_id = ? AND permission_id = ?", uid, pid).First(&existing)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		if shouldAdd {
			newPerm := models.UserPermission{
				UserID:       uint(uid),
				PermissionID: uint(pid),
				Granted:      true,
			}
			if err := pm.Db.WithContext(c).Create(&newPerm).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "create_failed"}
			}
		}
	case nil:
		if !shouldAdd {
			if err := pm.Db.Delete(&existing).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "delete_failed"}
			}
		} else if !existing.Granted {
			existing.Granted = true
			if err := pm.Db.Save(&existing).Error; err != nil {
				return nil, fuego.InternalServerError{Title: "update_failed"}
			}
		}
	default:
		return nil, fuego.InternalServerError{Title: "query_failed"}
	}

	return &okResponse{Ok: true}, nil
}

func (pm *PermissionsModule) grantAllPermissionsToUser(c fuego.ContextNoBody) (*okResponse, error) {
	userID := c.PathParam("userId")

	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fuego.BadRequestError{Title: "invalid_user_id"}
	}

	var allPermissions []models.Permission
	if err := pm.Db.WithContext(c).Model(models.Permission{}).Find(&allPermissions).Error; err != nil {
		return nil, fuego.InternalServerError{Err: err}
	}

	for _, v := range allPermissions {
		newPerm := models.UserPermission{
			UserID:       uint(uid),
			PermissionID: v.ID,
			Granted:      true,
		}

		if err := pm.Db.WithContext(c).Create(&newPerm).Error; err != nil {
			return nil, fuego.InternalServerError{Title: "create_failed"}
		}
	}

	return &okResponse{Ok: true}, nil
}

type GrantPermissionsOfTemplateToUserDto struct {
	Roles []uint `json:"roleIds"`
}

func (pm *PermissionsModule) grantPermissionsOfTemplateToUser(c fuego.ContextWithBody[GrantPermissionsOfTemplateToUserDto]) (*okResponse, error) {
	var zero *okResponse
	body, err := c.Body()
	if err != nil {
		return zero, err
	}

	userIDInt, err := c.PathParamIntErr("userId")
	if err != nil {
		return nil, err
	}
	userID := uint(userIDInt)

	// 1. Replace UserTemplates
	userTemplates := make([]models.UserTemplate, len(body.Roles))
	for i, templateID := range body.Roles {
		userTemplates[i] = models.UserTemplate{
			UserID:     userID,
			TemplateID: templateID,
		}
	}

	err = pm.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(c).Unscoped().
			Model(&models.User{Base: models.Base{ID: userID}}).
			Association("Roles").
			Replace(userTemplates); err != nil {
			return err
		}

		// 2. Fetch all permissions from the given templates
		var templates []models.PermissionTemplate
		if err := tx.WithContext(c).
			Preload("Permissions").
			Where("id IN ?", body.Roles).
			Find(&templates).Error; err != nil {
			return err
		}

		// 3. Flatten permissions into UserPermissions
		seen := make(map[uint]bool)
		var userPermissions []models.UserPermission
		for _, template := range templates {
			for _, perm := range template.Permissions {
				if seen[perm.ID] {
					continue
				}
				seen[perm.ID] = true
				userPermissions = append(userPermissions, models.UserPermission{
					UserID:       userID,
					PermissionID: perm.ID,
					Granted:      true,
				})
			}
		}

		// 4. Replace UserPermissions
		if err := tx.WithContext(c).Unscoped().
			Where("user_id = ?", userID).
			Delete(&models.UserPermission{}).Error; err != nil {
			return err
		}

		if len(userPermissions) > 0 {
			if err := tx.WithContext(c).Create(&userPermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &okResponse{Ok: true}, nil
}

func (pm *PermissionsModule) applyTemplate(c fuego.ContextNoBody) (*okResponse, error) {
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

	return &okResponse{Ok: true}, nil
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
