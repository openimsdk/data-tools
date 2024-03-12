package internal

import (
	"github.com/openimsdk/data-tools/chat/internal/mongo/model/admin"
	"github.com/openimsdk/data-tools/chat/internal/mongo/model/chat"
	"github.com/openimsdk/data-tools/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const (
	versionTable = "chatver"
	versionKey   = "data_version"
	versionValue = 1
)

func buildTask(mysqlDB *gorm.DB, mongoDB *mongo.Database) []func() error {
	var (
		cc convertChat
		ca convertAdmin
	)
	return []func() error{
		// chat
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewAccount, cc.Account) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewAttribute, cc.Attribute) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewLogs, cc.Log) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewRegister, cc.Register) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewUserLoginRecord, cc.UserLoginRecord) },
		// admin
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewAdmin, ca.Admin) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewApplet, ca.Applet) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewClientConfig, ca.ClientConfig) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewForbiddenAccount, ca.ForbiddenAccount) },
		func() error {
			return utils.NewTask(mysqlDB, mongoDB, admin.NewInvitationRegister, ca.InvitationRegister)
		},
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewIPForbidden, ca.IPForbidden) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewLimitUserLoginIP, ca.LimitUserLoginIP) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewRegisterAddFriend, ca.RegisterAddFriend) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, admin.NewRegisterAddGroup, ca.RegisterAddGroup) },
	}
}

func Main(path string) error {
	conf, err := utils.ParseConfig[Config](path)
	if err != nil {
		return err
	}
	mysqlDB, err := utils.NewMySQL(conf.MySQL)
	if err != nil {
		return err
	}
	mongoDB, err := utils.NewMongo(conf.Mongo)
	if err != nil {
		return err
	}
	if err := utils.CheckVersion(mongoDB.Collection(versionTable), versionKey, versionValue); err != nil {
		return err
	}
	tasks := buildTask(mysqlDB, mongoDB)
	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}
	if err := utils.SetVersion(mongoDB.Collection(versionTable), versionKey, versionValue); err != nil {
		return err
	}
	return nil
}
