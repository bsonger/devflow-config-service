package service

import (
	"context"
	"time"

	"github.com/bsonger/devflow-common/client/logging"
	"github.com/bsonger/devflow-common/client/mongo"
	"github.com/bsonger/devflow-config-service/pkg/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var ConfigurationService = NewConfigurationService()

type configurationService struct{}

func NewConfigurationService() *configurationService {
	return &configurationService{}
}

func (s *configurationService) Create(ctx context.Context, cfg *model.Configuration) (uuid.UUID, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "create_configuration"),
	)
	doc, err := configToDoc(cfg)
	if err != nil {
		log.Error("prepare configuration doc failed", zap.Error(err))
		return uuid.Nil, err
	}

	if err := mongo.Repo.Create(ctx, doc); err != nil {
		log.Error("create configuration failed", zap.Error(err))
		return uuid.Nil, err
	}
	cfg.ID = bridgeObjectIDToUUID(doc.ID)

	log.Info("configuration created", zap.String("configuration_id", cfg.GetID().String()))
	return cfg.GetID(), nil
}

func (s *configurationService) Get(ctx context.Context, id uuid.UUID) (*model.Configuration, error) {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return nil, err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "get_configuration"),
		zap.String("configuration_id", id.String()),
	)

	doc := &configurationDoc{}
	if err := mongo.Repo.FindByID(ctx, doc, oid); err != nil {
		log.Error("get configuration failed", zap.Error(err))
		return nil, err
	}
	if doc.DeletedAt != nil {
		log.Warn("configuration already deleted")
		return nil, mongoDriver.ErrNoDocuments
	}

	cfg := configFromDoc(doc)
	log.Debug("configuration fetched", zap.String("configuration_name", cfg.Name))
	return &cfg, nil
}

func (s *configurationService) Update(ctx context.Context, cfg *model.Configuration) error {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "update_configuration"),
		zap.String("configuration_id", cfg.GetID().String()),
	)
	cfgOID, err := bridgeUUIDToObjectID(cfg.GetID())
	if err != nil {
		return err
	}

	currentDoc := &configurationDoc{}
	if err := mongo.Repo.FindByID(ctx, currentDoc, cfgOID); err != nil {
		log.Error("load configuration failed", zap.Error(err))
		return err
	}
	if currentDoc.DeletedAt != nil {
		log.Warn("update skipped for deleted configuration")
		return mongoDriver.ErrNoDocuments
	}

	current := configFromDoc(currentDoc)
	cfg.CreatedAt = current.CreatedAt
	cfg.DeletedAt = current.DeletedAt
	cfg.WithUpdateDefault()

	doc, err := configToDoc(cfg)
	if err != nil {
		return err
	}
	doc.ID = cfgOID

	if err := mongo.Repo.Update(ctx, doc); err != nil {
		log.Error("update configuration failed", zap.Error(err))
		return err
	}

	log.Debug("configuration updated", zap.String("configuration_name", cfg.Name))
	return nil
}

func (s *configurationService) Delete(ctx context.Context, id uuid.UUID) error {
	oid, err := bridgeUUIDToObjectID(id)
	if err != nil {
		return err
	}
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "delete_configuration"),
		zap.String("configuration_id", id.String()),
	)

	now := time.Now()
	update := primitive.M{
		"$set": primitive.M{
			"deleted_at": now,
			"updated_at": now,
		},
	}

	if err := mongo.Repo.UpdateByID(ctx, &configurationDoc{}, oid, update); err != nil {
		log.Error("delete configuration failed", zap.Error(err))
		return err
	}

	log.Info("configuration deleted")
	return nil
}

func (s *configurationService) List(ctx context.Context, filter primitive.M) ([]model.Configuration, error) {
	log := logging.LoggerWithContext(ctx).With(
		zap.String("operation", "list_configurations"),
		zap.Any("filter", filter),
	)

	var docs []configurationDoc
	if err := mongo.Repo.List(ctx, &configurationDoc{}, filter, &docs); err != nil {
		log.Error("list configurations failed", zap.Error(err))
		return nil, err
	}

	cfgs := make([]model.Configuration, 0, len(docs))
	for i := range docs {
		cfgs = append(cfgs, configFromDoc(&docs[i]))
	}

	log.Debug("configurations listed", zap.Int("count", len(cfgs)))
	return cfgs, nil
}
