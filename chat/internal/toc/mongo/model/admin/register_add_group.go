// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admin

import (
	"context"
	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/openimsdk/data-tools/chat/internal/toc/mongo/table/admin"
)

func NewRegisterAddGroup(db *mongo.Database) (admin.RegisterAddGroupInterface, error) {
	coll := db.Collection("register_add_group")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "group_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &RegisterAddGroup{
		coll: coll,
	}, nil
}

type RegisterAddGroup struct {
	coll *mongo.Collection
}

func (o *RegisterAddGroup) Add(ctx context.Context, registerAddGroups []*admin.RegisterAddGroup) error {
	return mgoutil.InsertMany(ctx, o.coll, registerAddGroups)
}

func (o *RegisterAddGroup) Del(ctx context.Context, groupIDs []string) error {
	if len(groupIDs) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"group_id": bson.M{"$in": groupIDs}})
}

func (o *RegisterAddGroup) FindGroupID(ctx context.Context, groupIDs []string) ([]string, error) {
	filter := bson.M{}
	if len(groupIDs) > 0 {
		filter["group_id"] = bson.M{"$in": groupIDs}
	}
	return mgoutil.Find[string](ctx, o.coll, filter, options.Find().SetProjection(bson.M{"_id": 0, "group_id": 1}))
}

func (o *RegisterAddGroup) Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*admin.RegisterAddGroup, error) {
	filter := bson.M{"group_id": bson.M{"$regex": keyword, "$options": "i"}}
	return mgoutil.FindPage[*admin.RegisterAddGroup](ctx, o.coll, filter, pagination)
}
