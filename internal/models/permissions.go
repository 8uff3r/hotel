package models

type PermissionAction string

const (
	PermissionActionRead   PermissionAction = "read"
	PermissionActionCreate PermissionAction = "create"
	PermissionActionUpdate PermissionAction = "update"
	PermissionActionDelete PermissionAction = "delete"
	PermissionActionExport PermissionAction = "export"
)

type Permission struct {
	Base
	TranslateBase
	Page     string           `gorm:"not null;index" json:"page"`
	Action   PermissionAction `gorm:"not null" json:"action"`
	HotelID  *string          `json:"hotelId"`
	Category string           `gorm:"not null" json:"category"`
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
