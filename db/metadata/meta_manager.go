package metadata

import (
	"Vectory/entities"
	"Vectory/gen/ent"
	"Vectory/gen/ent/collection"
	"context"
	"fmt"
	_ "github.com/xiaoqidun/entps"
	"os"
)

// MetaManager is responsible for managing all of Vectory's metadata about collections, etc..
type MetaManager struct {
	db        *ent.Client
	filesPath string
}

func NewMetaManager(filesPath string) (*MetaManager, error) {
	c, err := ent.Open("sqlite3", "file:"+filesPath+"/metadata.db")
	if err != nil {
		return nil, err
	}

	if stat, err := os.Stat(filesPath); err != nil || !stat.IsDir() {
		return nil, ErrPathNotDirectory
	}

	// schemas auto migration
	err = c.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}

	return &MetaManager{db: c, filesPath: filesPath}, nil
}

func (m *MetaManager) CreateCollection(ctx context.Context, cfg *entities.Collection) (int, error) {
	params := cfg.IndexParams.(map[string]interface{}) // TODO: change this

	c, err := m.db.Collection.Create().
		SetName(cfg.Name).
		SetIndexType(cfg.IndexType).
		SetDataType(cfg.DataType).
		SetEmbedder(cfg.Embedder).
		SetIndexParams(params).
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
