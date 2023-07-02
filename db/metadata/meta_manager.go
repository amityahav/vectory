package metadata

import (
	"Vectory/db/metadata/ent"
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

	err = c.Schema.Create(context.Background())
	if err != nil {
		return nil, err
	}

	return &MetaManager{db: c}, nil
}

func (m *MetaManager) CreateCollection(ctx context.Context) error {
	_, err := m.db.Collection.Create().
		SetName("cat_images").
		SetIndex("disk_ann").
		SetDataType("jpeg").
		SetEmbedder("word2vec").
		Save(ctx)

	return err
}
