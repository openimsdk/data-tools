package organization

import (
	"context"
	"time"
)

type DepartmentMember struct {
	UserID          string     `gorm:"column:user_id;primary_key;size:64"`
	DepartmentID    string     `gorm:"column:department_id;primary_key;size:64"`
	Position        string     `gorm:"column:position;size:256"`
	Station         string     `gorm:"column:station;size:256"`
	Order           int32      `gorm:"column:order"`
	EntryTime       time.Time  `gorm:"column:entry_time"`       // 入职时间
	TerminationTime *time.Time `gorm:"column:termination_time"` // 离职时间
	CreateTime      time.Time  `gorm:"column:create_time"`
}

func (DepartmentMember) TableName() string {
	return "department_members"
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
