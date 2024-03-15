package tob

import (
	"github.com/openimsdk/data-tools/chat/internal/tob/mongo/table/chat"
	org "github.com/openimsdk/data-tools/chat/internal/tob/mongo/table/organization"
	oldchat "github.com/openimsdk/data-tools/chat/internal/tob/mysql/table/chat"
	oldorg "github.com/openimsdk/data-tools/chat/internal/tob/mysql/table/organization"
)

type ConvertOrganization struct{}

func (ConvertOrganization) Attribute(v oldchat.Attribute) chat.Attribute {
	return chat.Attribute{
		UserID:           v.UserID,
		Account:          v.Account,
		PhoneNumber:      v.PhoneNumber,
		AreaCode:         v.AreaCode,
		Email:            v.Email,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Gender:           v.Gender,
		CreateTime:       v.CreateTime,
		ChangeTime:       v.ChangeTime,
		BirthTime:        v.BirthTime,
		Level:            v.Level,
		AllowVibration:   v.AllowVibration,
		AllowBeep:        v.AllowBeep,
		AllowAddFriend:   v.AllowAddFriend,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
		RegisterType:     v.RegisterType,
		EnglishName:      v.EnglishName,
		Station:          v.Station,
		Telephone:        v.Telephone,
		Status:           v.Status,
	}
}

func (ConvertOrganization) Organization(v oldorg.Organization) org.Organization {
	return org.Organization{
		ID:           org.OrganizationID,
		LogoURL:      v.LogoURL,
		Name:         v.Name,
		Homepage:     v.Homepage,
		Introduction: v.Introduction,
		CreateTime:   v.CreateTime,
	}
}

func (ConvertOrganization) Department(v oldorg.Department) org.Department {
	return org.Department{
		DepartmentID:       v.DepartmentID,
		FaceURL:            v.FaceURL,
		Name:               v.Name,
		ParentDepartmentID: v.ParentDepartmentID,
		Order:              v.Order,
		CreateTime:         v.CreateTime,
	}
}

func (ConvertOrganization) DepartmentMember(v oldorg.DepartmentMember) org.DepartmentMember {
	return org.DepartmentMember{
		UserID:          v.UserID,
		DepartmentID:    v.DepartmentID,
		Position:        v.Position,
		Station:         v.Station,
		Order:           v.Order,
		EntryTime:       v.EntryTime,
		TerminationTime: v.TerminationTime,
		CreateTime:      v.CreateTime,
	}
}
