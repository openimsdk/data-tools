package organization

import (
	"context"
	"time"
)

type Organization struct {
	LogoURL      string    `gorm:"column:logo_url;size:255" json:"logoURL"`
	Name         string    `gorm:"column:name;size:256" json:"name" binding:"required"`
	Homepage     string    `gorm:"column:homepage" json:"homepage" `
	Introduction string    `gorm:"column:introduction;size:255" json:"introduction"`
	CreateTime   time.Time `gorm:"column:create_time" json:"createTime"`
}

func (Organization) TableName() string {
	return "organizations"
}

type OrganizationInterface interface {
	Set(ctx context.Context, update map[string]any) error
	Get(ctx context.Context) (*Organization, error)
	Init(ctx context.Context) error
}
