package internal

import (
	"fmt"
	"github.com/openimsdk/data-tools/chat/internal/tob"
	bchat "github.com/openimsdk/data-tools/chat/internal/tob/mongo/model/chat"
	"github.com/openimsdk/data-tools/chat/internal/tob/mongo/model/organization"
	"github.com/openimsdk/data-tools/chat/internal/toc"
	"github.com/openimsdk/data-tools/chat/internal/toc/mongo/model/admin"
	"github.com/openimsdk/data-tools/chat/internal/toc/mongo/model/chat"
	oldchat "github.com/openimsdk/data-tools/chat/internal/toc/mysql/table/chat"
	"github.com/openimsdk/data-tools/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const (
	versionTable = "chatver"
	versionKey   = "data_version"
	versionValue = 1
)

func buildTask(mysqlDB *gorm.DB, mongoDB *mongo.Database, isToc bool) []func() error {
	var (
		cc toc.ConvertChat
		ca toc.ConvertAdmin
	)
	fns := []func() error{
		// chat
		func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewAccount, cc.Account) },
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
	if isToc {
		fns = append(fns, func() error { return utils.NewTask(mysqlDB, mongoDB, chat.NewAttribute, cc.Attribute) })
	} else {
		var co tob.ConvertOrganization
		fns = append(fns,
			func() error { return utils.NewTask(mysqlDB, mongoDB, bchat.NewAttribute, co.Attribute) },
			func() error { return utils.NewTask(mysqlDB, mongoDB, organization.NewOrganization, co.Organization) },
			func() error { return utils.NewTask(mysqlDB, mongoDB, organization.NewDepartment, co.Department) },
			func() error {
				return utils.NewTask(mysqlDB, mongoDB, organization.NewDepartmentMember, co.DepartmentMember)
			})
	}
	return fns
}

func isToc(mysqlDB *gorm.DB, mode string, dbname string) (bool, error) {
	switch mode {
	case "toc":
		return true, nil
	case "tob":
		return false, nil
	case "auto":
	default:
		return false, fmt.Errorf("unknow mode %s", mode)
	}
	var attribute oldchat.Attribute
	var names []string
	err := mysqlDB.Table("information_schema.columns").Select("column_name").Where("table_schema = ? and table_name = ?", dbname, attribute.TableName()).Scan(&names).Error
	if err != nil {
		return false, err
	}
	fields := []string{"english_name", "station", "telephone"}
	var count int
	for _, name := range names {
		for _, filed := range fields {
			if name == filed {
				count++
				break
			}
		}
	}
	return len(fields) == count, nil
}

func Main(path string, mode string) error {
	conf, err := utils.ParseConfig[Config](path)
	if err != nil {
		return err
	}
	mysqlDB, err := utils.NewMySQL(conf.MySQL)
	if err != nil {
		return err
	}
	it, err := isToc(mysqlDB, mode, conf.MySQL.Database)
	if err != nil {
		return fmt.Errorf("mode recognition failed %w", err)
	}
	mongoDB, err := utils.NewMongo(conf.Mongo)
	if err != nil {
		return err
	}
	if err := utils.CheckVersion(mongoDB.Collection(versionTable), versionKey, versionValue); err != nil {
		return err
	}
	tasks := buildTask(mysqlDB, mongoDB, it)
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
