package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/utils"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mysql/table/organization"
	"gorm.io/gorm"
)

func NewDepartmentMember(db *gorm.DB) *DepartmentMember {
	return &DepartmentMember{
		db: db,
	}
}

type DepartmentMember struct {
	db *gorm.DB
}

func (o *DepartmentMember) FindByDepartmentID(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error) {
	if len(departmentIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id in ?", departmentIDList).Find(&ms).Error, "")
}

func (o *DepartmentMember) DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error {
	return utils.Wrap(o.db.WithContext(ctx).Where("department_id in ?", departmentIDList).Delete(&table.DepartmentMember{}).Error, "")
}

func (o *DepartmentMember) Create(ctx context.Context, m *table.DepartmentMember) error {
	return utils.Wrap(o.db.WithContext(ctx).Create(m).Error, "")
}

func (o *DepartmentMember) Creates(ctx context.Context, m []*table.DepartmentMember) error {
	return utils.Wrap(o.db.WithContext(ctx).Create(&m).Error, "")
}

func (o *DepartmentMember) Get(ctx context.Context, userID string) ([]*table.DepartmentMember, error) {
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Find(&ms).Error, "")
}

func (o *DepartmentMember) DeleteByKey(ctx context.Context, userID, departmentID string) error {
	return utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? AND department_id = ?", userID, departmentID).Delete(&table.DepartmentMember{}).Error, "")
}

func (o *DepartmentMember) Update(ctx context.Context, departmentID string, userID string, update map[string]any) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.DepartmentMember{}).Where("user_id = ? AND department_id = ?", userID, departmentID).Updates(update).Error, "")
}

func (o *DepartmentMember) FindByUserID(ctx context.Context, userIDList []string) ([]*table.DepartmentMember, error) {
	if len(userIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id in ?", userIDList).Find(&ms).Error, "")
}

func (o *DepartmentMember) GetByDepartmentID(ctx context.Context, departmentID string) ([]*table.DepartmentMember, error) {
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&ms).Error, "")
}

func (o *DepartmentMember) GetByKey(ctx context.Context, userID, departmentID string) (*table.DepartmentMember, error) {
	var ms table.DepartmentMember
	return &ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? and department_id = ?", userID, departmentID).First(&ms).Error, "")
}

func (o *DepartmentMember) Move(ctx context.Context, userID string, oldDepartmentID string, newDepartmentID string) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.DepartmentMember{}).Where("user_id = ? AND department_id = ?", userID, oldDepartmentID).Updates(map[string]any{"department_id": newDepartmentID}).Error)
}

func (o *DepartmentMember) Search(ctx context.Context, keyword string, departmentIDs []string) ([]string, error) {
	db := o.db.WithContext(ctx).Model(&[]table.DepartmentMember{}).Select("user_id")
	if len(departmentIDs) > 0 {
		db = db.Where("department_id in ? OR position like concat('%',?,'%') OR station like concat('%',?,'%')", departmentIDs, keyword, keyword)
	} else {
		db = db.Where("position like concat('%',?,'%') OR station like concat('%',?,'%')", keyword, keyword)
	}
	var userIDs []string
	if err := db.Find(&userIDs).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return departmentIDs, nil
}

func (o *DepartmentMember) GetMaxOrder(ctx context.Context, departmentID string) (int32, error) {
	var order int32
	return order, errs.Wrap(o.db.WithContext(ctx).Model(&table.DepartmentMember{}).Select("IFNULL(MAX(`order`), 0)").Where("department_id = ?", departmentID).Scan(&order).Error)
}

func (o *DepartmentMember) IncrOrder(ctx context.Context, departmentID string, order int32) error {
	return errs.Wrap(o.db.WithContext(ctx).Model(&table.DepartmentMember{}).Where("department_id = ? and `order` >= ?", departmentID, order).Updates(map[string]any{
		"`order`": gorm.Expr("`order` + ?", 1),
	}).Error)
}

func (o *DepartmentMember) GetNum(ctx context.Context, departmentIDs []string) (int64, error) {
	var count int64
	db := o.db.WithContext(ctx).Model(&[]table.DepartmentMember{}).Distinct("user_id")
	if len(departmentIDs) > 0 {
		db = db.Where("department_id in ?", departmentIDs)
	}
	return count, errs.Wrap(db.Count(&count).Error)
}
