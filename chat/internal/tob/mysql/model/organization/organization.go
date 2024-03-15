package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mysql/table/organization"
	"gorm.io/gorm"
	"time"
)

func NewOrganization(db *gorm.DB) *Organization {
	return &Organization{
		db: db,
	}
}

type Organization struct {
	db *gorm.DB
}

func (o *Organization) Set(ctx context.Context, update map[string]any) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.Organization{}).Where("1 = 1").Updates(update).Error)
}

func (o *Organization) Get(ctx context.Context) (*table.Organization, error) {
	var m table.Organization
	if err := o.db.WithContext(ctx).First(&m).Error; err == gorm.ErrRecordNotFound {
		m.CreateTime = time.UnixMilli(0)
	} else if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return &m, nil
}

func (o *Organization) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	return o.db.WithContext(ctx).Begin(), nil
}

func (o *Organization) Init(ctx context.Context) error {
	var m table.Organization
	if err := o.db.WithContext(ctx).First(&m).Error; err == nil {
		return nil
	} else if err == gorm.ErrRecordNotFound {
		return errs.Wrap(o.db.Create(&table.Organization{
			Name:       "xxx",
			CreateTime: time.Now(),
		}).Error)
	} else {
		return err
	}
}
