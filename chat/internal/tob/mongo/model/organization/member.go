package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mgoutil"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mongo/table/organization"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDepartmentMember(db *mongo.Database) (table.DepartmentMemberInterface, error) {
	coll := db.Collection("department_member")
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "department_id", Value: 1},
				{Key: "user_id", Value: 1},
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
	return &DepartmentMember{coll: coll}, nil
}

type DepartmentMember struct {
	coll *mongo.Collection
}

func (o *DepartmentMember) isNotFound(err error) bool {
	return errs.Unwrap(err) == mongo.ErrNoDocuments
}

func (o *DepartmentMember) FindByDepartmentID(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error) {
	if len(departmentIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	return mgoutil.Find[*table.DepartmentMember](ctx, o.coll, bson.M{"department_id": bson.M{"$in": departmentIDList}})
}

func (o *DepartmentMember) DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error {
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"department_id": bson.M{"$in": departmentIDList}})
}

func (o *DepartmentMember) Create(ctx context.Context, m *table.DepartmentMember) error {
	return mgoutil.InsertMany(ctx, o.coll, []*table.DepartmentMember{m})
}

func (o *DepartmentMember) Creates(ctx context.Context, m []*table.DepartmentMember) error {
	return mgoutil.InsertMany(ctx, o.coll, m)
}

func (o *DepartmentMember) Get(ctx context.Context, userID string) ([]*table.DepartmentMember, error) {
	return mgoutil.Find[*table.DepartmentMember](ctx, o.coll, bson.M{"user_id": userID})
}

func (o *DepartmentMember) DeleteByKey(ctx context.Context, userID, departmentID string) error {
	return mgoutil.DeleteOne(ctx, o.coll, bson.M{"user_id": userID, "department_id": departmentID})
}

func (o *DepartmentMember) Update(ctx context.Context, departmentID string, userID string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"user_id": userID, "department_id": departmentID}, bson.M{"$set": update}, false)
}

func (o *DepartmentMember) FindByUserID(ctx context.Context, userIDList []string) ([]*table.DepartmentMember, error) {
	if len(userIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	return mgoutil.Find[*table.DepartmentMember](ctx, o.coll, bson.M{"user_id": bson.M{"$in": userIDList}})
}

func (o *DepartmentMember) GetByDepartmentID(ctx context.Context, departmentID string) ([]*table.DepartmentMember, error) {
	return mgoutil.Find[*table.DepartmentMember](ctx, o.coll, bson.M{"department_id": departmentID})
}

func (o *DepartmentMember) GetByKey(ctx context.Context, userID, departmentID string) (*table.DepartmentMember, error) {
	return mgoutil.FindOne[*table.DepartmentMember](ctx, o.coll, bson.M{"user_id": userID, "department_id": departmentID})
}

func (o *DepartmentMember) Move(ctx context.Context, userID string, oldDepartmentID string, newDepartmentID string) error {
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"department_id": oldDepartmentID, "user_id": userID}, bson.M{"$set": bson.M{"department": newDepartmentID}}, false)
}

func (o *DepartmentMember) Search(ctx context.Context, keyword string, departmentIDs []string) ([]string, error) {
	filter := bson.M{}
	if len(departmentIDs) > 0 {
		filter["department_id"] = bson.M{
			"$in": departmentIDs,
		}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{
				"position": bson.M{
					"$regex": keyword,
				},
			},
			{
				"station": bson.M{
					"$regex": keyword,
				},
			},
		}
	}
	opt := options.Find().SetProjection(bson.M{"_id": 0, "user_id": 1})
	return mgoutil.Find[string](ctx, o.coll, filter, opt)
}

func (o *DepartmentMember) GetMaxOrder(ctx context.Context, departmentID string) (int32, error) {
	opt := options.FindOne().SetSort(bson.M{"order": -1}).SetProjection(bson.M{"order": 1})
	res, err := mgoutil.FindOne[*table.DepartmentMember](ctx, o.coll, bson.M{"department_id": departmentID}, opt)
	if err == nil {
		return res.Order, nil
	} else if o.isNotFound(err) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (o *DepartmentMember) IncrOrder(ctx context.Context, departmentID string, order int32) error {
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"department_id": departmentID, "order": bson.M{"$gte": order}}, bson.M{"$inc": bson.M{"order": 1}}, false)
}

func (o *DepartmentMember) GetNum(ctx context.Context, departmentIDs []string) (int64, error) {
	pipeline := make([]bson.M, 0, 4)
	if len(departmentIDs) > 0 {
		pipeline = append(pipeline,
			bson.M{
				"$match": bson.M{
					"department_id": bson.M{
						"$in": departmentIDs,
					},
				},
			},
		)
	}
	pipeline = append(pipeline,
		bson.M{
			"$project": bson.M{
				"_id":     0,
				"user_id": 1,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$user_id",
			},
		},
		bson.M{
			"$count": "count",
		},
	)
	res, err := mgoutil.Aggregate[int64](ctx, o.coll, pipeline)
	if err != nil {
		return 0, err
	}
	if len(res) == 0 {
		return 0, nil
	}
	return res[0], nil
}
