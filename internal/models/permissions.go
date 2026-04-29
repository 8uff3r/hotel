package models

type PermissionAction string

const (
	PermissionActionRead   PermissionAction = "read"
	PermissionActionCreate PermissionAction = "create"
	PermissionActionUpdate PermissionAction = "update"
	PermissionActionDelete PermissionAction = "delete"
	PermissionActionExport PermissionAction = "export"
)

type PermissionCategory struct {
	Base
	TranslateBase
}

type Permission struct {
	Base
	Translation Translation        `gorm:"type:jsonb" json:"translation,omitempty"`
	Resource    string             `gorm:"not null;index;uniqueIndex:idx_permissions_resource_action" json:"resource"`
	Action      PermissionAction   `gorm:"not null;uniqueIndex:idx_permissions_resource_action" json:"action"`
	CategoryID  uint               `gorm:"not null" json:"categoryId"`
	Category    PermissionCategory `gorm:"not null;foreignKey:CategoryID" json:"category"`
}

type PermissionTemplate struct {
	Base
	TranslateBase
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:template_permissions;" json:"permissions"`
}

type UserPermission struct {
	Base
	UserID       uint       `gorm:"not null;index" json:"userId"`
	HotelID      *string    `json:"hotelId"`
	PermissionID uint       `gorm:"not null;index" json:"permissionId"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
	Granted      bool       `gorm:"not null;default:true" json:"granted"`
}

type UserTemplate struct {
	Base
	UserID     uint               `gorm:"not null;index" json:"userId"`
	TemplateID uint               `gorm:"not null;index" json:"templateId"`
	Template   PermissionTemplate `gorm:"foreignKey:TemplateID" json:"template"`
}
