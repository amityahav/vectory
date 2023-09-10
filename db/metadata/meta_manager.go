package metadata

import (
	collectionent "Vectory/entities/collection"
	"Vectory/gen/ent"
	"Vectory/gen/ent/collection"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/xiaoqidun/entps"
	"os"
)

// MetaManager is responsible for managing all of Vectory's metadata of collections, etc..
type MetaManager struct {
	db        *ent.Client
	filesPath string
}

func NewMetaManager(filesPath string) (*MetaManager, error) {
	err := os.Mkdir(filesPath, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	c, err := ent.Open("sqlite3", "file:"+filesPath+"/metadata.db")
	if err != nil {
		return nil, err
	}

	// schemas auto migration
	err = c.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}

	return &MetaManager{db: c, filesPath: filesPath}, nil
}

func (m *MetaManager) CreateCollection(ctx context.Context, cfg *collectionent.Collection) (int, error) {
	b, err := json.Marshal(cfg.IndexParams)
	if err != nil {
		return 0, nil
	}

	var params map[string]interface{}
	err = json.Unmarshal(b, &params)
	if err != nil {
		return 0, err
	}

	b, err = json.Marshal(cfg.EmbedderConfig)
	if err != nil {
		return 0, nil
	}

	var config map[string]interface{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		return 0, err
	}

	b, err = json.Marshal(cfg.EmbedderConfig)
	if err != nil {
		return 0, nil
	}

	c, err := m.db.Collection.Create().
		SetName(cfg.Name).
		SetIndexType(cfg.IndexType).
		SetDataType(cfg.DataType).
		SetEmbedderType(cfg.EmbedderType).
		SetEmbedderConfig(config).
		SetIndexParams(params).
		SetMappings(cfg.Mappings).
		Save(ctx)

	if err != nil {
		return 0, err
	}

	// creating collection data's dir
	err = os.Mkdir(fmt.Sprintf("%s/%s", m.filesPath, cfg.Name), 0750)
	if err != nil && !os.IsExist(err) {
		return 0, err
	}

	return c.ID, err
}

func (m *MetaManager) DeleteCollection(ctx context.Context, name string) error {
	n, err := m.db.Collection.Delete().Where(collection.Name(name)).Exec(ctx)
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrCollectionDoesntExist
	}

	// delete all collection's files
	return os.RemoveAll(fmt.Sprintf("%s/%s", m.filesPath, name))
}

func (m *MetaManager) GetCollection(ctx context.Context, name string) (*ent.Collection, error) {
	return m.db.Collection.Query().Where(collection.Name(name)).Only(ctx)
}

func (m *MetaManager) GetCollections(ctx context.Context) ([]*ent.Collection, error) {
	return m.db.Collection.Query().All(ctx)
}

func (m *MetaManager) Close() error {
	return m.db.Close()
}
