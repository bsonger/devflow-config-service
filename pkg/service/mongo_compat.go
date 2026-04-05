package service

import (
	"errors"

	commonmodel "github.com/bsonger/devflow-common/model"
	"github.com/bsonger/devflow-config-service/pkg/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	errInvalidUUIDBridge = errors.New("invalid bridged uuid")
	bridgeUUIDPrefix     = [4]byte{'d', 'f', 'l', 'w'}
)

type configurationDoc struct {
	commonmodel.BaseModel `bson:",inline"`

	ApplicationID    *primitive.ObjectID `bson:"application_id,omitempty"`
	Name             string              `bson:"name"`
	Env              string              `bson:"env,omitempty"`
	Status           string              `bson:"status,omitempty"`
	LatestRevisionNo int                 `bson:"latest_revision_no,omitempty"`
	LatestRevisionID *primitive.ObjectID `bson:"latest_revision_id,omitempty"`
}

func (configurationDoc) CollectionName() string { return "configuration" }

func bridgeObjectIDToUUID(id primitive.ObjectID) uuid.UUID {
	var raw [16]byte
	copy(raw[:4], bridgeUUIDPrefix[:])
	copy(raw[4:], id[:])
	return uuid.UUID(raw)
}

func bridgeUUIDToObjectID(id uuid.UUID) (primitive.ObjectID, error) {
	raw := [16]byte(id)
	if raw[0] != bridgeUUIDPrefix[0] || raw[1] != bridgeUUIDPrefix[1] || raw[2] != bridgeUUIDPrefix[2] || raw[3] != bridgeUUIDPrefix[3] {
		return primitive.NilObjectID, errInvalidUUIDBridge
	}
	var oid primitive.ObjectID
	copy(oid[:], raw[4:])
	return oid, nil
}

func configFromDoc(doc *configurationDoc) model.Configuration {
	cfg := model.Configuration{
		BaseModel: model.BaseModel{
			ID:        bridgeObjectIDToUUID(doc.ID),
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
			DeletedAt: doc.DeletedAt,
		},
		Name:             doc.Name,
		Env:              doc.Env,
		Status:           doc.Status,
		LatestRevisionNo: doc.LatestRevisionNo,
	}
	if doc.ApplicationID != nil && !doc.ApplicationID.IsZero() {
		cfg.ApplicationID = bridgeObjectIDToUUID(*doc.ApplicationID)
	}
	if doc.LatestRevisionID != nil && !doc.LatestRevisionID.IsZero() {
		id := bridgeObjectIDToUUID(*doc.LatestRevisionID)
		cfg.LatestRevisionID = &id
	}
	return cfg
}

func configToDoc(cfg *model.Configuration) (*configurationDoc, error) {
	id := primitive.NewObjectID()
	if cfg.ID != uuid.Nil {
		bridgedID, err := bridgeUUIDToObjectID(cfg.ID)
		if err == nil {
			id = bridgedID
		}
	}

	var appID *primitive.ObjectID
	if cfg.ApplicationID != uuid.Nil {
		oid, err := bridgeUUIDToObjectID(cfg.ApplicationID)
		if err != nil {
			return nil, err
		}
		appID = &oid
	}

	var latestRevisionID *primitive.ObjectID
	if cfg.LatestRevisionID != nil && *cfg.LatestRevisionID != uuid.Nil {
		oid, err := bridgeUUIDToObjectID(*cfg.LatestRevisionID)
		if err != nil {
			return nil, err
		}
		latestRevisionID = &oid
	}

	return &configurationDoc{
		BaseModel: commonmodel.BaseModel{
			ID:        id,
			CreatedAt: cfg.CreatedAt,
			UpdatedAt: cfg.UpdatedAt,
			DeletedAt: cfg.DeletedAt,
		},
		ApplicationID:    appID,
		Name:             cfg.Name,
		Env:              cfg.Env,
		Status:           cfg.Status,
		LatestRevisionNo: cfg.LatestRevisionNo,
		LatestRevisionID: latestRevisionID,
	}, nil
}
