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
	"github.com/openimsdk/data-tools/chat/internal/constant"
	"github.com/openimsdk/data-tools/chat/internal/mongo/table/admin"
)

func NewIPForbidden(db *mongo.Database) (admin.IPForbiddenInterface, error) {
	coll := db.Collection("ip_forbidden")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "ip", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &IPForbidden{
		coll: coll,
	}, nil
}

type IPForbidden struct {
	coll *mongo.Collection
}

func (o *IPForbidden) Take(ctx context.Context, ip string) (*admin.IPForbidden, error) {
	return mgoutil.FindOne[*admin.IPForbidden](ctx, o.coll, bson.M{"ip": ip})

}

func (o *IPForbidden) Find(ctx context.Context, ips []string) ([]*admin.IPForbidden, error) {
	return mgoutil.Find[*admin.IPForbidden](ctx, o.coll, bson.M{"ip": bson.M{"$in": ips}})

}

func (o *IPForbidden) Search(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*admin.IPForbidden, error) {
	filter := bson.M{}

	switch state {
	case constant.LimitNil:
		// 不添加额外的过滤条件
	case constant.LimitEmpty:
		filter = bson.M{"limit_register": 0, "limit_login": 0}
	case constant.LimitOnlyRegisterIP:
		filter = bson.M{"limit_register": 1, "limit_login": 0}
	case constant.LimitOnlyLoginIP:
		filter = bson.M{"limit_register": 0, "limit_login": 1}
	case constant.LimitRegisterIP:
		filter = bson.M{"limit_register": 1}
	case constant.LimitLoginIP:
		filter = bson.M{"limit_login": 1}
	case constant.LimitLoginRegisterIP:
		filter = bson.M{"limit_register": 1, "limit_login": 1}
	}

	if keyword != "" {
		filter["$or"] = []bson.M{
			{"ip": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mgoutil.FindPage[*admin.IPForbidden](ctx, o.coll, filter, pagination)
}

func (o *IPForbidden) Create(ctx context.Context, ms []*admin.IPForbidden) error {
	return mgoutil.InsertMany(ctx, o.coll, ms)
}

func (o *IPForbidden) Delete(ctx context.Context, ips []string) error {
	if len(ips) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"ip": bson.M{"$in": ips}})
}
