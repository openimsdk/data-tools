package organization

import (
	"context"
	"time"
)

type Department struct {
	DepartmentID       string    `gorm:"column:department_id;primary_key;size:64"`
	FaceURL            string    `gorm:"column:face_url;size:255"`
	Name               string    `gorm:"column:name;size:256;uniqueIndex:idx_name_parent_department_id"`
	ParentDepartmentID string    `gorm:"column:parent_department_id;size:64;uniqueIndex:idx_name_parent_department_id"`
	Order              int32     `gorm:"column:order"`
	CreateTime         time.Time `gorm:"column:create_time"`
}

func (Department) TableName() string {
	return "departments"
}

type DepartmentInterface interface {
	IncrOrder(ctx context.Context, parentDepartmentID string, order int32) error
	Create(ctx context.Context, department ...*Department) error
	FindOne(ctx context.Context, departmentID string) (*Department, error)
	Update(ctx context.Context, departmentID string, data map[string]any) error
	GetParent(ctx context.Context, id string) ([]*Department, error)
	GetList(ctx context.Context, departmentIdList []string) ([]*Department, error)
	UpdateParentID(ctx context.Context, oldParentID, newParentID string) error
	Delete(ctx context.Context, departmentIDList []string) error
	GetDepartment(ctx context.Context, departmentID string) (*Department, error)
	GetMaxOrder(ctx context.Context, parentID string) (int32, error)
	GetByName(ctx context.Context, name string, id string) (*Department, error)
	InitUngroupedName(ctx context.Context, id string, name string) error
	Search(ctx context.Context, keyword string) ([]string, error)
}
