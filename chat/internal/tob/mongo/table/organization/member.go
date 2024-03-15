package organization

import (
	"context"
	"time"
)

type DepartmentMember struct {
	UserID          string     `bson:"user_id"`
	DepartmentID    string     `bson:"department_id"`
	Position        string     `bson:"position"`
	Station         string     `bson:"station"`
	Order           int32      `bson:"order"`
	EntryTime       time.Time  `bson:"entry_time"`
	TerminationTime *time.Time `bson:"termination_time"`
	CreateTime      time.Time  `bson:"create_time"`
}

type DepartmentMemberInterface interface {
	FindByDepartmentID(ctx context.Context, departmentIDList []string) ([]*DepartmentMember, error)
	DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error
	Create(ctx context.Context, m *DepartmentMember) error
	Creates(ctx context.Context, m []*DepartmentMember) error
	Get(ctx context.Context, userID string) ([]*DepartmentMember, error)
	DeleteByKey(ctx context.Context, userID string, departmentID string) error
	Update(ctx context.Context, departmentID string, userID string, update map[string]any) error
	FindByUserID(ctx context.Context, userIDList []string) ([]*DepartmentMember, error)
	GetByDepartmentID(ctx context.Context, departmentID string) ([]*DepartmentMember, error)
	GetByKey(ctx context.Context, userID, departmentID string) (*DepartmentMember, error)
	Move(ctx context.Context, userID string, oldDepartmentID string, newDepartmentID string) error
	Search(ctx context.Context, keyword string, departmentIDs []string) ([]string, error)
	GetMaxOrder(ctx context.Context, departmentID string) (int32, error)
	IncrOrder(ctx context.Context, departmentID string, order int32) error
	GetNum(ctx context.Context, departmentIDs []string) (int64, error)
}
