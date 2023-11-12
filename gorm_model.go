package reback

import (
	"time"

	"gorm.io/gorm"
)

type base struct {
	ID        string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	CreatedAt time.Time      `json:"created_at,omitempty" gorm:"default:now()"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" gorm:"default:now()"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	DeletedBy string `json:"deleted_by"`
}

type Role struct {
	base
	Name        string `json:"name"`
	Description string `json:"description"`

	Permissions []Permission `json:"permissions" gorm:"foreignKey:RoleId"`
}

type Permission struct {
	base

	RoleId string `json:"role_id"`
	Role   Role   `json:"role" gorm:"foreignKey:RoleId"`

	Resource string `json:"resource_name,omitempty"`

	IsCreate *bool `json:"is_create" gorm:"default:false"`
	IsRead   *bool `json:"is_read" gorm:"default:false"`
	IsUpdate *bool `json:"is_update" gorm:"default:false"`
	IsDelete *bool `json:"is_delete" gorm:"default:false"`
}
