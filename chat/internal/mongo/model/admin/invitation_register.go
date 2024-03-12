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

func NewInvitationRegister(db *mongo.Database) (admin.InvitationRegisterInterface, error) {
	coll := db.Collection("invitation_register")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "invitation_code", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &InvitationRegister{
		coll: coll,
	}, nil
}

type InvitationRegister struct {
	coll *mongo.Collection
}

func (o *InvitationRegister) Find(ctx context.Context, codes []string) ([]*admin.InvitationRegister, error) {
	return mgoutil.Find[*admin.InvitationRegister](ctx, o.coll, bson.M{"invitation_code": bson.M{"$in": codes}})
}

func (o *InvitationRegister) Del(ctx context.Context, codes []string) error {
	if len(codes) == 0 {
		return nil
	}
	return mgoutil.DeleteMany(ctx, o.coll, bson.M{"invitation_code": bson.M{"$in": codes}})
}

func (o *InvitationRegister) Create(ctx context.Context, v []*admin.InvitationRegister) error {
	return mgoutil.InsertMany(ctx, o.coll, v)
}

func (o *InvitationRegister) Take(ctx context.Context, code string) (*admin.InvitationRegister, error) {
	return mgoutil.FindOne[*admin.InvitationRegister](ctx, o.coll, bson.M{"code": code})
}

func (o *InvitationRegister) Update(ctx context.Context, code string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return mgoutil.UpdateOne(ctx, o.coll, bson.M{"invitation_code": code}, bson.M{"$set": data}, false)
}

func (o *InvitationRegister) Search(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, pagination pagination.Pagination) (int64, []*admin.InvitationRegister, error) {
	filter := bson.M{}
	switch state {
	case constant.InvitationCodeUsed:
		filter = bson.M{"user_id": bson.M{"$ne": ""}}
	case constant.InvitationCodeUnused:
		filter = bson.M{"user_id": ""}
	}

	if len(userIDs) > 0 {
		filter["user_id"] = bson.M{"$in": userIDs}
	}
	if len(codes) > 0 {
		filter["invitation_code"] = bson.M{"$in": codes}
	}
	if keyword != "" {
		filter["$or"] = []bson.M{
			{"invitation_code": bson.M{"$regex": keyword, "$options": "i"}},
			{"user_id": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}
	return mgoutil.FindPage[*admin.InvitationRegister](ctx, o.coll, filter, pagination)

}
