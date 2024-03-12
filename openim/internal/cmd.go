package internal

import (
	"github.com/openimsdk/data-tools/openim/internal/im/mgo"
	rtcmgo "github.com/openimsdk/data-tools/openim/internal/rtc/mgo"
	"github.com/openimsdk/data-tools/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const (
	versionTable = "dataver"
	versionKey   = "data_version"
	versionValue = 35
)

func buildTask(mysqlDB *gorm.DB, mongoDB *mongo.Database, enable string) []func() error {
	var c convert
	return []func() error{
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewUserMongo, c.User) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewFriendMongo, c.Friend) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewFriendRequestMongo, c.FriendRequest) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewBlackMongo, c.Black) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewGroupMongo, c.Group) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewGroupMember, c.GroupMember) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewGroupRequestMgo, c.GroupRequest) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewConversationMongo, c.Conversation) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewS3Mongo, c.Object(enable)) },
		func() error { return utils.NewTask(mysqlDB, mongoDB, mgo.NewLogMongo, c.Log) },

		func() error { return utils.NewTask(mysqlDB, mongoDB, rtcmgo.NewSignal, c.SignalModel) },
		func() error {
			return utils.NewTask(mysqlDB, mongoDB, rtcmgo.NewSignalInvitation, c.SignalInvitationModel)
		},
		func() error { return utils.NewTask(mysqlDB, mongoDB, rtcmgo.NewMeeting, c.Meeting) },
		func() error {
			return utils.NewTask(mysqlDB, mongoDB, rtcmgo.NewMeetingInvitation, c.MeetingInvitationInfo)
		},
		func() error { return utils.NewTask(mysqlDB, mongoDB, rtcmgo.NewMeetingRecord, c.MeetingVideoRecord) },
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
	tasks := buildTask(mysqlDB, mongoDB, conf.S3Engine)
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
