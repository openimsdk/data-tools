package organization

import (
	"context"
	"time"
)

type Department struct {
	DepartmentID       string    `bson:"department_id"`
	FaceURL            string    `bson:"face_url"`
	Name               string    `bson:"name"`
	ParentDepartmentID string    `bson:"parent_department_id"`
	Order              int32     `bson:"order"`
	CreateTime         time.Time `bson:"create_time"`
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
