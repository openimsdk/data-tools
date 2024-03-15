package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mysql/table/organization"
	"gorm.io/gorm"
	"time"
)

func NewDepartment(db *gorm.DB) *Department {
	return &Department{
		db: db,
	}
}

type Department struct {
	db *gorm.DB
}

func (o *Department) IncrOrder(ctx context.Context, parentDepartmentID string, order int32) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("parent_department_id = ? and `order` >= ?", parentDepartmentID, order).Updates(map[string]any{
		"`order`": gorm.Expr("`order` + ?", 1),
	}).Error)
}

func (o *Department) GetParent(ctx context.Context, parentID string) ([]*table.Department, error) {
	var ms []*table.Department
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("parent_department_id = ?", parentID).Order("`order` ASC, `create_time` ASC").Find(&ms).Error, "")
}

func (o *Department) Update(ctx context.Context, departmentID string, data map[string]any) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("department_id = ?", departmentID).Updates(data).Error, "")
}

func (o *Department) Create(ctx context.Context, departments ...*table.Department) error {
	return errs.Wrap(o.db.WithContext(ctx).Create(departments).Error)
}

func (o *Department) FindOne(ctx context.Context, departmentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ?", departmentID).First(&m).Error, "")
}

func (o *Department) GetList(ctx context.Context, departmentIdList []string) ([]*table.Department, error) {
	if len(departmentIdList) == 0 {
		return []*table.Department{}, nil
	}
	var ms []*table.Department
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIdList).Order("`order` ASC, `create_time` ASC").Find(&ms).Error, "")
}

func (o *Department) UpdateParentID(ctx context.Context, oldParentID, newParentID string) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("parent_department_id = ?", oldParentID).Update("parent_department_id", newParentID).Error, "")
}

func (o *Department) Delete(ctx context.Context, departmentIDList []string) error {
	if len(departmentIDList) == 0 {
		return nil
	}
	return utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIDList).Delete(&table.Department{}).Error, "")
}

func (o *Department) GetDepartment(ctx context.Context, departmentId string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ?", departmentId).First(&m).Error, "")
}

func (o *Department) GetMaxOrder(ctx context.Context, parentID string) (int32, error) {
	var order int32
	return order, utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Select("IFNULL(MAX(`order`), 0)").Where("parent_department_id = ?", parentID).Scan(&order).Error, "")
}

func (o *Department) GetByName(ctx context.Context, name, parentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.WithContext(ctx).Where("name = ? AND parent_department_id = ?", name, parentID).First(&m).Error, "")
}

func (o *Department) Search(ctx context.Context, keyword string) ([]string, error) {
	var departmentIDs []string
	if err := o.db.WithContext(ctx).Model(&[]table.Department{}).Select("department_id").Where("name like concat('%',?,'%')", keyword).Scan(&departmentIDs).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return departmentIDs, nil
}

func (o *Department) InitUngroupedName(ctx context.Context, id string, name string) error {
	var m table.Department
	if err := o.db.WithContext(ctx).Where("department_id = ?", id).First(&m).Error; err == nil {
		//if m.Name != name {
		//	return o.Update(ctx, "", map[string]any{"name": name, "change_time": time.Now()})
		//}
		return nil
	} else if err == gorm.ErrRecordNotFound {
		return errs.Wrap(o.db.WithContext(ctx).Create(&table.Department{
			DepartmentID:       id,
			FaceURL:            "",
			Name:               name,
			ParentDepartmentID: "",
			Order:              -1,
			CreateTime:         time.Now(),
		}).Error)
	} else {
		return errs.Wrap(err)
	}
}
