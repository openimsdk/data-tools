package organization

import (
	"context"
	"github.com/OpenIMSDK/tools/mgoutil"
	table "github.com/openimsdk/data-tools/chat/internal/tob/mongo/table/organization"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func NewOrganization(db *mongo.Database) (table.OrganizationInterface, error) {
	coll := db.Collection("organization")
	return &Organization{coll: coll}, nil
}

type Organization struct {
	coll *mongo.Collection
}

func (o *Organization) Init(ctx context.Context) error {
	if err := mgoutil.InsertMany(ctx, o.coll, []*table.Organization{{ID: table.OrganizationID, Name: "xxx", CreateTime: time.Now()}}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (o *Organization) Set(ctx context.Context, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"_id": table.OrganizationID}, bson.M{"$set": update}, true)
}

func (o *Organization) Get(ctx context.Context) (*table.Organization, error) {
	return mgoutil.FindOne[*table.Organization](ctx, o.coll, bson.M{"_id": table.OrganizationID})
}
