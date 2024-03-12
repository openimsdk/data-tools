// Copyright © 2023 OpenIM. All rights reserved.
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

package utils

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Uri         string   `yaml:"uri"`
	Address     []string `yaml:"address"`
	Database    string   `yaml:"database"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	MaxPoolSize int      `yaml:"maxPoolSize"`
}

func (c MongoConfig) buildURI() string {
	if c.Username != "" && c.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d", c.Username, c.Password, c.Address, c.Database, c.MaxPoolSize)
	}
	return fmt.Sprintf("mongodb://%s/%s?maxPoolSize=%d", c.Address, c.Database, c.MaxPoolSize)
}

func NewMongo(conf MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.buildURI()))
	if err != nil {
		return nil, err
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return mongoClient.Database(conf.Database), nil
}
