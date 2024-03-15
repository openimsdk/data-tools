package organization

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var OrganizationID primitive.ObjectID

func init() {
	for i := range OrganizationID {
		OrganizationID[i] = byte(i + 1)
	}
}

type Organization struct {
	ID           primitive.ObjectID `bson:"_id"`
	LogoURL      string             `bson:"logo_url"`
	Name         string             `bson:"name"`
	Homepage     string             `bson:"homepage"`
	Introduction string             `bson:"introduction"`
	CreateTime   time.Time          `bson:"create_time"`
}

type OrganizationInterface interface {
	Set(ctx context.Context, update map[string]any) error
	Get(ctx context.Context) (*Organization, error)
	Init(ctx context.Context) error
}
