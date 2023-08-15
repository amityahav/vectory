package collection

import (
	"Vectory/db/embeddings"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/index"
	"errors"
)

func Validate(cfg *Collection) error {
	if cfg.Name == "" {
		return ErrCollectionNameEmpty
	}

	var err error

	switch cfg.IndexType {
	case index.Hnsw:
		err = index.ValidateHnswParams(cfg.IndexParams)
	case index.DiskAnn:
	default:
		return ErrIndexTypeUnsupported
	}

	if err != nil {
		return err
	}

	switch cfg.EmbedderType {
	case "": // no use of embedder, user should provide his own vectors
	case text2vec.Text2VecHuggingFace: // TODO: validate embedder config
	case embeddings.FakeEmbedder:
	default:
		return ErrEmbedderTypeUnsupported
	}

	switch cfg.DataType {
	case TextDataType:
	default:
		return ErrDataTypeUnsupported
	}

	return nil
}

var (
	ErrCollectionNameEmpty     = errors.New("collection name field is empty")
	ErrIndexTypeUnsupported    = errors.New("index_type inserted is not supported")
	ErrEmbedderTypeUnsupported = errors.New("embedder_type inserted is not supported")
	ErrDataTypeUnsupported     = errors.New("data_type inserted is not supported")
)
