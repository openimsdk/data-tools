package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"reflect"
)

// NewTask A mysql table B mongodb model C mongodb table
func NewTask[A interface{ TableName() string }, B any, C any](gormDB *gorm.DB, mongoDB *mongo.Database, mongoDBInit func(db *mongo.Database) (B, error), convert func(v A) C) error {
	var zero A
	tableName := zero.TableName()
	obj, err := mongoDBInit(mongoDB)
	if err != nil {
		return fmt.Errorf("init mongo table %s failed, err: %w", tableName, err)
	}
	coll, err := getColl(obj)
	if err != nil {
		return fmt.Errorf("get mongo collection %s failed, err: %w", tableName, err)
	}
	var count int
	defer func() {
		fmt.Printf("completed convert chat %s total %d\n", tableName, count)
	}()
	const batch = 100
	for page := 0; ; page++ {
		res := make([]A, 0, batch)
		if err := gormDB.Limit(batch).Offset(page * batch).Find(&res).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1146 {
				return nil // table not exist
			}
			return fmt.Errorf("find mysql table %s failed, err: %w", tableName, err)
		}
		if len(res) == 0 {
			return nil
		}
		temp := make([]any, len(res))
		for i := range res {
			temp[i] = convert(res[i])
		}
		if err := insertMany(coll, temp); err != nil {
			return fmt.Errorf("insert mongo table %s failed, err: %w", tableName, err)
		}
		count += len(res)
		if len(res) < batch {
			return nil
		}
		fmt.Printf("current convert chat %s completed %d\n", tableName, count)
	}
}

func insertMany(coll *mongo.Collection, objs []any) error {
	if _, err := coll.InsertMany(context.Background(), objs); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}
	for i := range objs {
		_, err := coll.InsertOne(context.Background(), objs[i])
		switch {
		case err == nil:
		case mongo.IsDuplicateKeyError(err):
		default:
			return err
		}
	}
	return nil
}

func getColl(obj any) (_ *mongo.Collection, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("not found %+v", e)
		}
	}()
	stu := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(&mongo.Collection{}).String()
	for i := 0; i < stu.NumField(); i++ {
		field := stu.Field(i)
		if field.Type().String() == typ {
			return (*mongo.Collection)(field.UnsafePointer()), nil
		}
	}
	return nil, errors.New("not found model collection")
}
