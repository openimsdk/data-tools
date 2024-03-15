package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mgoutil"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mongo/table/organization"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewDepartment(db *mongo.Database) (table.DepartmentInterface, error) {
	coll := db.Collection("department")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "department_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "order", Value: 1},
			},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Department{coll: coll}, nil
}

type Department struct {
	coll *mongo.Collection
}

func (o *Department) isNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}

func (o *Department) IncrOrder(ctx context.Context, parentDepartmentID string, order int32) error {
	filter := bson.M{
		"department_id": parentDepartmentID,
		"order": bson.M{
			"$gte": order,
		},
	}
	update := bson.M{
		"$inc": bson.M{
			"order": 1,
		},
	}
	return mgoutil.UpdateOne(ctx, o.coll, filter, update, false)
}

func (o *Department) GetParent(ctx context.Context, parentID string) ([]*table.Department, error) {
	opt := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "create_time", Value: 1}})
	return mgoutil.Find[*table.Department](ctx, o.coll, bson.M{"parent_department_id": parentID}, opt)
}

func (o *Department) Update(ctx context.Context, departmentID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"department_id": departmentID}, bson.M{"$set": data}, false)
}

func (o *Department) Create(ctx context.Context, departments ...*table.Department) error {
	return mgoutil.InsertMany(ctx, o.coll, departments)
}

func (o *Department) FindOne(ctx context.Context, departmentID string) (*table.Department, error) {
	return mgoutil.FindOne[*table.Department](ctx, o.coll, bson.M{"department_id": departmentID})
}

func (o *Department) GetList(ctx context.Context, departmentIdList []string) ([]*table.Department, error) {
	if len(departmentIdList) == 0 {
		return []*table.Department{}, nil
	}
	opt := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "create_time", Value: 1}})
	return mgoutil.Find[*table.Department](ctx, o.coll, bson.M{"department_id": bson.M{"$in": departmentIdList}}, opt)
}

func (o *Department) UpdateParentID(ctx context.Context, oldParentID, newParentID string) error {
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"parent_department_id": oldParentID}, bson.M{"$set": bson.M{"parent_department_id": newParentID}}, false)
}

func (o *Department) Delete(ctx context.Context, departmentIDList []string) error {
	if len(departmentIDList) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"department_id": bson.M{"$in": departmentIDList}})
}

func (o *Department) GetDepartment(ctx context.Context, departmentId string) (*table.Department, error) {
	return mgoutil.FindOne[*table.Department](ctx, o.coll, bson.M{"department": departmentId})
}

func (o *Department) GetByName(ctx context.Context, name, parentID string) (*table.Department, error) {
	return mgoutil.FindOne[*table.Department](ctx, o.coll, bson.M{"name": name, "parent_department_id": parentID})
}

func (o *Department) InitUngroupedName(ctx context.Context, id string, name string) error {
	if _, err := o.FindOne(ctx, id); err == nil {
		return nil
	} else if !o.isNotFound(err) {
		return err
	}
	return o.Create(ctx, &table.Department{
		DepartmentID:       id,
		FaceURL:            "",
		Name:               name,
		ParentDepartmentID: "",
		Order:              -1,
		CreateTime:         time.Now(),
	})
}

func (o *Department) GetMaxOrder(ctx context.Context, parentID string) (int32, error) {
	opt := options.FindOne().SetSort(bson.M{"order": -1}).SetProjection(bson.M{"order": 1})
	res, err := mgoutil.FindOne[*table.Department](ctx, o.coll, bson.M{"parent_department_id": parentID}, opt)
	if err == nil {
		return res.Order, nil
	} else if o.isNotFound(err) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (o *Department) Search(ctx context.Context, keyword string) ([]string, error) {
	opt := options.Find().SetProjection(bson.M{"_id": 0, "department_id": 1})
	filter := bson.M{
		"name": bson.M{
			"$regex": keyword,
		},
	}
	return mgoutil.Find[string](ctx, o.coll, filter, opt)
}
