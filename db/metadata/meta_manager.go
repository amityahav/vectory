package metadata

import (
	"Vectory/gen/api/models"
	"Vectory/gen/ent"
	"context"
	_ "github.com/xiaoqidun/entps"
)

type MetaManager struct {
	db *ent.Client
}

func NewMetaManager(filesPath string) (*MetaManager, error) {
	c, err := ent.Open("sqlite3", "file:"+filesPath+"/metadata.db")
	if err != nil {
		return nil, err
	}

	// schemas auto migration
	err = c.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}

	return &MetaManager{db: c}, nil
}

func (m *MetaManager) CreateCollection(ctx context.Context, cfg *models.Collection) (int, error) {
	params := cfg.IndexParams.(map[string]interface{}) // TODO: change this

	c, err := m.db.Collection.Create().
		SetName(cfg.Name).
		SetIndexType(cfg.IndexType).
		SetDataType(cfg.DataType).
		SetEmbedder(cfg.Embedder).
		SetIndexParams(params).
		Save(ctx)

	return c.ID, err
}

func (m *MetaManager) DeleteCollection(ctx context.Context, name string) error {
	return nil
}

func (m *MetaManager) GetCollections(ctx context.Context) ([]*ent.Collection, error) {
	return m.db.Collection.Query().All(ctx)
}
